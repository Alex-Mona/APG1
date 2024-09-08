package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
)

type JSON struct {
	Cake []Cake
}

type Cake struct {
	Name        string
	Time        string
	Ingredients []Ingredients
}

type Ingredients struct {
	Ingredient_name  string
	Ingredient_count string
	Ingredient_unit  string
}

type XML struct {
	Cake []struct {
		Name        string `xml:"name"`
		Stovetime   string `xml:"stovetime"`
		Ingredients []struct {
			Item []struct {
				Itemname  string `xml:"itemname"`
				Itemcount string `xml:"itemcount"`
				Itemunit  string `xml:"itemunit"`
			} `xml:"item"`
		} `xml:"ingredients"`
	} `xml:"cake"`
}

func PrintXml(filename string) {
	f_xml, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f_xml.Close()
	xml_data, err := io.ReadAll(f_xml)
	if err != nil {
		log.Fatal(err)
	}
	var xml_result XML
	xml_err := xml.Unmarshal(xml_data, &xml_result)
	if xml_err != nil {
		log.Fatal(xml_err)
	}
	json_res, _ := json.MarshalIndent(xml_result, "", "    ")
	fmt.Printf("%s", json_res)
}

func PrintJson(filename string) {
	f_json, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f_json.Close()
	json_data, err := io.ReadAll(f_json)
	if err != nil {
		log.Fatal(err)
	}
	var json_result JSON
	json_err := json.Unmarshal(json_data, &json_result)
	if json_err != nil {
		log.Fatal(json_err)
	}
	xml_res, _ := xml.MarshalIndent(json_result, "", "    ")
	fmt.Printf("%s", xml_res)
}

func ParsingArguments() {
	count_args := len(os.Args[1:])
	if count_args == 0 {
		fmt.Println("There is no argument. Enter *.xml/*.json file")
		os.Exit(1)
	} else if count_args > 1 {
		fmt.Println("Enter only one argument - the file name")
		os.Exit(1)
	} else {
		len_arg := len(os.Args[1])
		if len_arg < 5 {
			fmt.Println("Invalid value entered")
			os.Exit(1)
		}
		type_file := os.Args[1][len(os.Args[1])-4:]
		switch type_file {
		case ".xml":
			PrintXml(os.Args[1])
		case "json":
			PrintJson(os.Args[1])
		default:
			fmt.Println("File missing")
			os.Exit(1)
		}
		fmt.Println(type_file)
	}
}

func main() {
	ParsingArguments()
}
