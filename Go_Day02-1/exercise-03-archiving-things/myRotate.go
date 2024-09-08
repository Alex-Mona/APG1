package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func Archive(fileName string, directory string, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	withoutLog := strings.TrimSuffix(fileName, ".log")
	currentTime := time.Now()
	timestamp := currentTime.Unix()
	timestampFormat := strconv.FormatInt(timestamp, 10)
	archiveName := directory + withoutLog + "_" + timestampFormat + ".tar.gz"
	archiveFile, err := os.Create(archiveName)
	Error(err)
	defer archiveFile.Close()
	zipWriter := zip.NewWriter(archiveFile)
	defer zipWriter.Close()
	fileToAdd, err := os.Open(fileName)
	Error(err)
	defer fileToAdd.Close()
	fileInArchive, err := zipWriter.Create(fileName)
	Error(err)
	_, err = io.Copy(fileInArchive, fileToAdd)
	Error(err)
}

func Error(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ErrorOptions(files []string) {
	if len(files) == 0 {
		fmt.Println("SYNOPSIS\n       ./myRotate [FILE]")
		os.Exit(1)
	}
}

func myRotate() {
	directory := flag.String("a", "", "Archive creation directory")
	flag.Parse()
	if *directory != "" {
		*directory = *directory + "/"
	}
	fileNames := flag.Args()
	ErrorOptions(fileNames)
	var waitGroup sync.WaitGroup
	for _, fileName := range fileNames {
		waitGroup.Add(1)
		go Archive(fileName, *directory, &waitGroup)
	}
	waitGroup.Wait()
}

func main() {
	myRotate()
}
