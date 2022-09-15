package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	// with mandala
	"github.com/withmandala/go-log"
)

// open missing.txt
// for every missing file
//   read file from _src/DSC000.HIF
//   write file to _missing/DSC000.HIF

func main() {
	logger := log.New(os.Stderr).WithColor()
	if os.Getenv("DEBUG") != "" {
		logger = logger.WithDebug()
	}

	logger.Info("Opened missing.txt")
	fileBytes, errRead := os.Open("missing.txt")
	if errRead != nil {
		panic(errRead)
	}
	defer fileBytes.Close()
	logger.Debug("Opened missing.txt")

	scanner := bufio.NewScanner(fileBytes)
	damagedHeader := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "These files appear to be unsupported or damaged") {
			damagedHeader = true
			continue
		}
		if damagedHeader {
			parts := strings.Split(line, "/")
			fileName := parts[len(parts)-1]
			striped := strings.Replace(fileName, ".heic", "", 1)
			fileHif := striped + ".HIF"
			if errCopy := copyFile("_src/"+fileHif, "_missing/"+fileHif, logger); errCopy != nil {
				panic(errCopy)
			}
		}
	}
}

func copyFile(copyFrom string, copyTo string, logger *log.Logger) error {
	if _, err := os.Stat(copyFrom); err != nil {
		logger.Debug(fmt.Sprintf("File %s does not exist", copyFrom))
		// nil on purpose, so we keep processing
		return nil
	}

	fileBytes, errRead := os.ReadFile(copyFrom)
	if errRead != nil {
		return errRead
	}

	errWrite := os.WriteFile(copyTo, fileBytes, 0644)
	if errWrite != nil {
		return errWrite
	}

	return nil
}
