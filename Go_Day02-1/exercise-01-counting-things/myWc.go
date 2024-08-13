package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"unicode/utf8"
)

func ReadingFile(fileName string, countingLines bool, countingCharacters bool, countingWords bool, wg *sync.WaitGroup, numberFiles int) {
	defer wg.Done()
	file, err := os.Open(fileName)
	Error(err)
	defer file.Close()
	buf := make([]byte, 10000)
	countLines, countChars, countWords := 0, 0, 0
	newLine := false
	for {
		n, err := file.Read(buf)
		if err != nil {
			break
		}
		data := string(buf[:n])
		lines := strings.Split(data, "\n")
		for i, line := range lines {
			if i == len(lines)-1 {
				if newLine || len(line) > 0 {
					countLines++
				}
			} else {
				countLines++
				newLine = true
			}
			countChars += utf8.RuneCountInString(line)
			words := strings.Fields(line)
			countWords += len(words)
		}
	}
	if countingLines {
		fmt.Printf("%8d", countLines)
	}
	if countingWords {
		fmt.Printf("%8d", countWords)
	}
	if countingCharacters {
		fmt.Printf("%8d", countChars+countLines)
	}
	fmt.Printf(" %s\n", fileName)
}

func Error(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ErrorOptions(files []string) {
	if len(files) == 0 {
		fmt.Println("SYNOPSIS\n       ./myWc [OPTION]: [-l] or [-m] or [-w]... [FILE]...")
		os.Exit(1)
	}
}

func myWc() {
	countingLines := flag.Bool("l", false, "Counting Lines")
	countingCharacters := flag.Bool("m", false, "Counting Characters")
	countingWords := flag.Bool("w", false, "Counting Words")
	flag.Parse()
	files := flag.Args()
	numberFiles := len(files)
	ErrorOptions(files)
	var wg sync.WaitGroup
	for _, fileName := range files {
		wg.Add(1)
		go ReadingFile(fileName, *countingLines, *countingCharacters, *countingWords, &wg, numberFiles)
	}
	wg.Wait()
}

func main() {
	myWc()
}
