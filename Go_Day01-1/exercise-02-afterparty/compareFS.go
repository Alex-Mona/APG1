package main

import (
	"bufio"
	"fmt"
	"os"
)

func CheckFileExtension(str string) {
	if str != ".txt" {
		fmt.Println("File missing")
		os.Exit(1)
	}
}

func ParsingArguments() map[string]string {
	if len(os.Args[1:]) != 4 {
		fmt.Println("Enter file names with --old or --new flags")
		os.Exit(1)
	}
	files := map[string]string{
		"old": "",
		"new": "",
	}
	for i, e := range os.Args[1:] {
		if i%2 == 0 && (e == "--old" || e == "--new") {
			if files[e[2:]] == "" {
				files[e[2:]] = os.Args[2+i]
			} else {
				fmt.Println("There cannot be two new or old flags")
				os.Exit(1)
			}
		} else if i%2 == 1 {
			continue
		} else {
			fmt.Println("Enter file names with --old or --new flags")
			os.Exit(1)
		}
	}
	len_old := len(files["old"])
	len_new := len(files["new"])
	if len_old < 5 || len_new < 5 {
		fmt.Println("Invalid value entered")
		os.Exit(1)
	}
	if files["old"] == files["new"] {
		fmt.Println("Same file entered")
		os.Exit(1)
	}
	CheckFileExtension(files["old"][len_old-4:])
	CheckFileExtension(files["new"][len_new-4:])
	return files
}

func Calculate() {
	files := ParsingArguments()
	flag_old, err := os.Open(files["old"])
	if err != nil {
		panic(err)
	}
	defer flag_old.Close()

	flag_new, err := os.Open(files["new"])
	if err != nil {
		panic(err)
	}
	defer flag_new.Close()
	old_scan := bufio.NewScanner(flag_old)
	new_scan := bufio.NewScanner(flag_new)
	var old_lines, new_lines []string
	for old_scan.Scan() {
		old_lines = append(old_lines, old_scan.Text())
	}
	for new_scan.Scan() {
		new_lines = append(new_lines, new_scan.Text())
	}
	for _, string_new := range new_lines {
		flag := 0
		for _, string_old := range old_lines {
			if string_new == string_old {
				flag = 1
				break
			}
		}
		if flag == 0 {
			fmt.Printf("ADDED %s\n", string_new)
		}
	}
	for _, string_old := range old_lines {
		flag := 0
		for _, string_new := range new_lines {
			if string_new == string_old {
				flag = 1
				break
			}
		}
		if flag == 0 {
			fmt.Printf("REMOVED %s\n", string_old)
		}
	}
}

func main() {
	Calculate()
}
