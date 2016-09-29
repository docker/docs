package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	outputFile, err := os.Create("OSS-LICENSES.txt")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer outputFile.Close()
	outputWriter := bufio.NewWriter(outputFile)
	directories := os.Args[1:]
	if len(directories) == 0 {
		directories = append(directories, ".")
	}
	for _, directory := range directories {
		writeLicenceInDir(directory, outputWriter)
	}
}

func writeLicenceInDir(dir string, outputWriter *bufio.Writer) {
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		fullPath := filepath.Join(dir, f.Name())
		if f.IsDir() {
			writeLicenceInDir(fullPath, outputWriter)
			continue
		}
		if strings.HasPrefix(f.Name(), "LICENSE") && !strings.HasSuffix(f.Name(), ".skip") {
			fmt.Println("Outputing", fullPath)
			bytes, err := ioutil.ReadFile(fullPath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			outputWriter.WriteString(fmt.Sprintf("begin %s\n", fullPath))
			outputWriter.Write(bytes)
			outputWriter.WriteString(fmt.Sprintf("end %s\n\n", fullPath))
			outputWriter.Flush()
		}
	}
}
