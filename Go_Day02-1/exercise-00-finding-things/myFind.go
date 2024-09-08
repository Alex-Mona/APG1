package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func FilePath(rootDirectory string, onlySymlink bool, onlyDirectory bool, onlyFiles bool, fileExtension string) ([]string, error) {
	var result []string
	err := filepath.Walk(rootDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				return nil
			}
			return err
		}
		if onlySymlink && info.Mode()&os.ModeSymlink != 0 {
			targetPath, err := os.Readlink(path)
			if err == nil {
				result = append(result, path+" -> "+targetPath)
			} else {
				result = append(result, path+" -> [broken]")
			}
			return nil
		}
		if onlyFiles && !info.IsDir() {
			if fileExtension == "" || filepath.Ext(info.Name()) == ("."+fileExtension) {
				result = append(result, path)
			}
		}
		if onlyDirectory && info.IsDir() {
			result = append(result, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func myFind() {
	onlySymlink := flag.Bool("sl", false, "Finding symbolic links")
	onlyDirectory := flag.Bool("d", false, "Directory search")
	onlyFiles := flag.Bool("f", false, "Search files")
	fileExtension := flag.String("ext", "", "File extension")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("SYNOPSIS\n       ./myFind [-sl] [-d] [-f] [-f -ext] [starting-point...]")
		return
	}
	if !*onlyFiles && !*onlyDirectory && !*onlySymlink && *fileExtension == "" {
		*onlyFiles = true
		*onlyDirectory = true
	}
	if !*onlyFiles && *fileExtension != "" {
		fmt.Println("./myFind: Works ONLY when -f is specified\nUsage: ./myFind [-f -ext] [starting-point...]")
		return
	}
	directory := args[len(args)-1]
	files, err := FilePath(directory, *onlyFiles, *onlyDirectory, *onlySymlink, *fileExtension)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, file := range files {
		fmt.Println(file)
	}
}

func main() {
	myFind()
}
