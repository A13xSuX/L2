package main

//for launch
//go mod init test
//go build -o cut.exe
//"1,2,3,4,5,6" | .\cut.exe -d "," -f "1,3-4,6" in PowerShell on Windows
//"john,25,london,engineer" | .\cut.exe -d "," -f "1,3" -s

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	fields    string
	delimiter string
	separated bool
}

func parseFields(fieldsStr string) ([]int, error) {
	result := make([]int, 0)
	if fieldsStr == "" { // check empty fields
		return result, errors.New("Fields is empty")
	}
	split := strings.Split(fieldsStr, ",")
	for _, field := range split {
		if strings.Contains(field, "-") { //range
			manyNum := strings.Split(field, "-")
			if len(manyNum) != 2 {
				return result, errors.New("Fields is invalid, len!= 2")
			}
			start, err := strconv.Atoi(strings.TrimSpace(manyNum[0]))
			if err != nil {
				return result, errors.New("Start is invalid")
			}
			end, err := strconv.Atoi(strings.TrimSpace(manyNum[len(manyNum)-1]))
			if err != nil {
				return result, errors.New("End is invalid")
			}
			if start > end {
				return result, errors.New("Start cannot be greater than End")
			}
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
		} else { //solo element
			num, err := strconv.Atoi(field)
			if err != nil {
				return result, err
			}
			result = append(result, num)
		}
	}
	return result, nil
}

func processLine(line string, fields []int, delimiter string, separated bool) string {
	if !strings.Contains(line, delimiter) && separated {
		return ""
	}

	column := strings.Split(line, delimiter)
	var resultFields []string

	for _, fieldNum := range fields { //fieldNum from user begin from 1
		if fieldNum-1 < len(column) && fieldNum > 0 {
			resultFields = append(resultFields, column[fieldNum-1])
		} //ignore fields for out
	}
	return strings.Join(resultFields, delimiter)
}

func main() {
	var config Config

	flag.StringVar(&config.fields, "f", "", "Fields to use in the input")
	flag.StringVar(&config.delimiter, "d", "\t", "Delimiter to use in the input")
	flag.BoolVar(&config.separated, "s", false, "Separated by newline")
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	fields, err := parseFields(config.fields)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for scanner.Scan() {
		line := scanner.Text()
		result := processLine(line, fields, config.delimiter, config.separated)
		if result != "" {
			fmt.Println(result)
		}

	}

}
