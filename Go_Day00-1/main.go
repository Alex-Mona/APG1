package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv" 
)

func InputArray() []int {
	var arr []int
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		txt := sc.Bytes()
		if len(txt) != 0 {
			if i, err := strconv.Atoi(string(txt)); err != nil || i < -100000 || i > 100000 {
				fmt.Println("Enter valid values in the range from -100.000 to 100.000")
			} else {
				arr = append(arr, i)
			}
		} else if len(txt) == 0 {
			break
		}
	}
	sort.Ints(arr)
	return arr
}

func ParseArguments(Args *map[string]bool) {
	q := len(os.Args[1:])
	if q != 0 && q < 5 {
		for _, e := range os.Args[1:] {
			switch e {
			case "Mean":
				(*Args)["Mean"] = true
			case "Median":
				(*Args)["Median"] = true
			case "Mode":
				(*Args)["Mode"] = true
			case "SD":
				(*Args)["SD"] = true
			default:
				fmt.Println("Invalid argument. Valid arguments: Mean, Median, Mode, SD")
				os.Exit(1)
			}
		}
	} else if q > 4 {
		fmt.Println("Incorrect number of arguments. Enter up to 4 Arguments")
		os.Exit(1)
	} else {
		(*Args)["Mean"] = true
		(*Args)["Median"] = true
		(*Args)["Mode"] = true
		(*Args)["SD"] = true
	}
}

func Mean(arr []int) float64 {
	var res float64
	count := 0
	for _, e := range arr {
		res += float64(e)
		count++
	}
	return res / float64(count)
}

func Median(arr []int) {
	if len(arr)%2 == 1 {
		fmt.Printf("Median: %d\n", arr[len(arr)/2])
	} else {
		fmt.Printf("Median: %.02f\n", (float64(arr[len(arr)/2])+float64(arr[len(arr)/2-1]))/2.0)
	}
}

func Mode(arr []int) {
	numbers := make(map[int]int)
	max := 0
	min := 100001
	for _, e := range arr {
		numbers[e]++
	}
	for i, e := range numbers {
		if e > max {
			max = e
			min = i
		} else if e == max && i < min {
			min = i
		}
	}
	fmt.Printf("Mode: %d\n", min)
}

func SD(arr []int) {
	if len(arr) == 1 {
		fmt.Printf("SD: NaN\n")
	} else {
		Mean := Mean(arr)
		var q, res float64
		for _, e := range arr {
			q = (float64(e) - Mean) * (float64(e) - Mean)
			res += q
		}
		res = math.Sqrt(res / float64(len(arr)-1))
		fmt.Printf("SD: %.02f\n", res)
	}
}

func Calculate(arr []int, Args map[string]bool) {
	for i, e := range Args {
		if e {
			switch i {
			case "Mean":
				fmt.Printf("Mean: %.02f\n", Mean(arr))
			case "Median":
				Median(arr)
			case "Mode":
				Mode(arr)
			case "SD":
				SD(arr)
			}
		}
	}
}

func main() {
	Args := map[string]bool{
		"Mean":   false,
		"Median": false,
		"Mode":   false,
		"SD":     false,
	}
	ParseArguments(&Args)
	Sequence := InputArray()
	if len(Sequence) != 0 {
		Calculate(Sequence, Args)
	} else {
		fmt.Println("Empty sequence")
	}
}
