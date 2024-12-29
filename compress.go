package main
import (
	"bufio"
	"fmt"
	"os"
	"io"
	"strings"
	"compress/gzip"
)

func DecompressFile(inputFile, outputFile string) error {
	inFile, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo comprimido: %w", err)
	}
	defer inFile.Close()

	gzipReader, err := gzip.NewReader(inFile)
	if err != nil {
		return fmt.Errorf("erro ao criar o leitor gzip: %w", err)
	}
	defer gzipReader.Close()

	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("erro ao criar o arquivo de saída: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, gzipReader)
	if err != nil {
		return fmt.Errorf("erro ao descomprimir o arquivo: %w", err)
	}

	fmt.Println("Arquivo descomprimido com sucesso!")
	return nil
}

func CompressFile(inputFile, outputFile string) error {

	inFile, err := os.Open(inputFile)
	if(err != nil){
		return fmt.Errorf("erro ao abrir o arquivo de entrada: %w", err)
	}
	defer inFile.Close()

	outFile, err := os.Create(outputFile)
	if(err != nil){
		return fmt.Errorf("erro ao criar arquivo compactado %w", err)
	}
	defer outFile.Close()

	gzipWriter, err := gzip.NewWriterLevel(outFile, gzip.BestCompression)
	if(err != nil){
		return fmt.Errorf("erro ao criar arquivo gzip %w", err)
	}
	defer gzipWriter.Close()

	_, err = io.Copy(gzipWriter, inFile)
	if(err != nil){
		return fmt.Errorf("erro ao copiar o arquivo comprimido %w", err)
	}

	compressedSize, err := getTamanho(outputFile)
	if(err != nil){
		return fmt.Errorf("erro ao calcular tamanho arquivo de saida %w", err)
	}

	originalSize, err := getTamanho(inputFile)
	if(err != nil){
		return fmt.Errorf("erro ao calcular tamanho arquivo de entrada %w", err)
	}

	reduction := 100 - (compressedSize/originalSize)*100
	fmt.Println("===========================")
	fmt.Printf("Arquivo Original: %.2f MB \n", originalSize)
	fmt.Printf("Arquivo Comprimido: %.2f MB \n", compressedSize)
	fmt.Printf("Redução: %.2f%% \n", reduction)
	fmt.Println("===========================")

	return nil
}

func getTamanho(filename string) (float64, error){
	fileInfo, err := os.Stat(filename)
	if(err != nil){
		return 0, fmt.Errorf("erro ao obter dados do arquivo %w", err)
	}

	return float64(fileInfo.Size()) / (1024 * 1024), nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("###### Informe a opção desejada ######")
		fmt.Println("1 - Comprimir")
		fmt.Println("2 - Descomprimir")
		fmt.Println("3 - Sair")

		var opcao int
		fmt.Scanln(&opcao) // ponteiro
		switch opcao {
		case 1:
			fmt.Println("Informe o nome do arquivo:")
			inputFile, _ := reader.ReadString('\n')
			inputFile = strings.TrimSpace(inputFile)

			fmt.Println("informe o nome do arquivo comprimido (.gz)")
			outputFile, _ := reader.ReadString('\n')
			outputFile = strings.TrimSpace(outputFile)

			if err := CompressFile(inputFile, outputFile); err != nil {
				fmt.Println("Error:", err)
			}
		case 2:
			fmt.Println("Informe o nome do arquivo de desecompressao:")
			inputFile, _ := reader.ReadString('\n')
			inputFile = strings.TrimSpace(inputFile)

			fmt.Println("informe o nome do arquivo descomprimido")
			outputFile, _ := reader.ReadString('\n')
			outputFile = strings.TrimSpace(outputFile)

			if err := DecompressFile(inputFile, outputFile); err != nil {
				fmt.Println("Error:", err)
			}
		case 3:
			fmt.Println("Sair do Programa")
			return
		default:
			fmt.Println("Opção Inválida!")
		}
	}
}