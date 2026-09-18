package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Uso: go run extract_audio.go <video_entrada.mp4> <audio_saida.mp3>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

cmd := exec.Command("ffmpeg", "-i", inputFile, "-vn", "-b:a", "64k", "-ac", "1", outputFile)

	// Captura o Stderr onde o FFmpeg envia as informações de progresso
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("Erro ao criar pipe para o FFmpeg: %v", err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("Erro ao iniciar o FFmpeg: %v", err)
	}

	scanner := bufio.NewScanner(stderrPipe)
	// Configura o scanner para ler quebras de linha e retornos de carro (\r),
	// já que o FFmpeg atualiza a linha de progresso usando \r.
	scanner.Split(scanLinesOrCarriageReturns)

	var totalDuration float64 = 0
	durationRe := regexp.MustCompile(`Duration:\s*(\d{2}):(\d{2}):([\d.]+)`)
	timeRe := regexp.MustCompile(`time\s*=\s*(\d{2}):(\d{2}):([\d.]+)`)

	fmt.Printf("Extraindo áudio de '%s'...\n", inputFile)

	for scanner.Scan() {
		line := scanner.Text()

		// Tenta extrair a duração total do vídeo (aparece nos metadados iniciais)
		if totalDuration == 0 {
			if matches := durationRe.FindStringSubmatch(line); len(matches) == 4 {
				if d, err := parseDurationToSeconds(matches[1], matches[2], matches[3]); err == nil && d > 0 {
					totalDuration = d
				}
			}
		}

		// Tenta extrair o tempo atual processado para calcular a %
		if totalDuration > 0 {
			if timeMatches := timeRe.FindStringSubmatch(line); len(timeMatches) == 4 {
				if current, err := parseDurationToSeconds(timeMatches[1], timeMatches[2], timeMatches[3]); err == nil {
					pct := (current / totalDuration) * 100
					if pct > 100 {
						pct = 100
					}
					printProgressBar(pct)
				}
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		log.Fatalf("\nErro ao executar o FFmpeg: %v", err)
	}

	fmt.Printf("\n\nPronto! Áudio extraído e salvo em: '%s'\n", outputFile)
}

// Converte o formato HH:MM:SS.ss para segundos totais
func parseDurationToSeconds(h, m, s string) (float64, error) {
	hours, err := strconv.Atoi(h)
	if err != nil {
		return 0, err
	}
	mins, err := strconv.Atoi(m)
	if err != nil {
		return 0, err
	}
	secs, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return float64(hours*3600+mins*60) + secs, nil
}

// Função auxiliar para escanear tanto por quebra de linha (\n) quanto por \r
func scanLinesOrCarriageReturns(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	for i := 0; i < len(data); i++ {
		if data[i] == '\n' || data[i] == '\r' {
			return i + 1, data[0:i], nil
		}
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

// Desenha a barra de progresso no terminal reescrevendo a linha (\r)
func printProgressBar(pct float64) {
	width := 30
	completed := int(pct / 100 * float64(width))
	if completed > width {
		completed = width
	}
	bar := strings.Repeat("=", completed) + strings.Repeat(" ", width-completed)
	fmt.Printf("\r[ %s ] %.1f%%", bar, pct)
}