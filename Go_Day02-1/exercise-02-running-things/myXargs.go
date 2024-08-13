package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func myXargs() {
	if len(os.Args) < 2 {
		fmt.Println("SYNOPSIS\n       ./myXargs [options] [command [initial-arguments]]")
		return
	}
	command := os.Args[1]
	args := os.Args[2:]
	scanner := bufio.NewScanner(os.Stdin)
	var directories []string
	for scanner.Scan() {
		dir := scanner.Text()
		directories = append(directories, dir)
	}
	for _, directory := range directories {
		directory := strings.ReplaceAll(directory, " ", "")
		cmd := exec.Command(command, append(args, directory)...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error running %s: %s\n", command, err)
			return
		}
		fmt.Println(string(output))
	}
}

func main() {
	myXargs()
}
