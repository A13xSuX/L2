package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Config struct {
	Keycolumn int  // сортировать по столбцу (колонке) №N
	Numeric   bool // сортировать по числовому значению (строки интерпретируются как числа).
	Reverse   bool // реверс
	Unique    bool // не выводить повторяющиеся строки
	// month     bool // сортировать по названию месяца (Jan, Feb, ... Dec), т.е. распознавать специфический формат дат.(Доделать)
	Space bool // хвостовые пробелы
	Check bool // были ли фильтрации
	// human     bool // человеческий вид(Доделать)
}

func getColumn(line string, column int) string {
	if line == "" {
		return ""
	}

	columns := strings.Fields(line)
	if column == 0 {
		fmt.Print("Ошибка: Колонки нумеруются с 1")
		return ""
	}
	if column < 1 || column > len(columns) {
		return "" //колонки нет
	}

	return columns[column-1]
}

func parseNumber(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return num
}

func removeDuplicates(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}

	result := []string{lines[0]}

	for i := 1; i < len(lines); i++ {
		if lines[i] != lines[i-1] {
			result = append(result, lines[i])
		}
	}
	return result
}

func checkSorted(lines []string, config Config) bool {
	for i := 1; i < len(lines); i++ {
		if !compareLines(lines[i-1], lines[i], config) {
			fmt.Printf("Disorder: %s\n", lines[i])
			return false
		}
	}
	return true
}

func compareLines(lineI, lineJ string, config Config) bool { //flag k
	var valueI, valueJ string

	if config.Keycolumn > 0 { //сравниваем только указанные колонки
		valueI = getColumn(lineI, config.Keycolumn)
		valueJ = getColumn(lineJ, config.Keycolumn)
		fmt.Printf("DEBUG: Comparing '%s' vs '%s'\n", valueI, valueJ)
	} else { //сравниваем полные строки
		valueI = lineI
		valueJ = lineJ
	}
	var result bool

	if config.Space {
		valueI = strings.TrimRight(valueI, " \t\n\r")
		valueJ = strings.TrimRight(valueJ, " \t\n\r")
	}

	if config.Numeric {
		numI := parseNumber(valueI)
		numJ := parseNumber(valueJ)
		result = numI < numJ
	}

	if config.Reverse {
		return !result
	}

	return result
}

func main() {
	//с файла
	var lines []string
	var filename string
	var config Config

	flag.IntVar(&config.Keycolumn, "k", 0, "sort by column")
	flag.BoolVar(&config.Numeric, "n", false, "numeric sort")
	flag.BoolVar(&config.Reverse, "r", false, "reverse")
	flag.BoolVar(&config.Unique, "u", false, "unique")
	// flag.BoolVar(&config.month, "M", false, "sort by months")
	flag.BoolVar(&config.Space, "b", false, "space")
	flag.BoolVar(&config.Check, "c", false, "check")
	// flag.BoolVar(&config.human, "h", false, "humanview")

	flag.Parse()

	if len(flag.Args()) > 0 {
		filename = flag.Args()[0] //./program -numeric data.txt          # numeric=true, filename="data.txt"
	} //./program -file=data.txt -numeric    # filename="data.txt", numeric=true

	var scanner *bufio.Scanner

	if filename != "" { //если есть файл
		file, err := os.Open(filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	} else {
		scanner = bufio.NewScanner(os.Stdin) //с консоли
	}

	for scanner.Scan() {
		line := scanner.Text()
		if filename == "" && line == "" { // остановка по пустой строке(только в stdin)
			break
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка чтения")
	}

	if config.Check {
		if checkSorted(lines, config) {
			fmt.Println("Данные отсортированы")
			os.Exit(0)
		} else {
			fmt.Println("Данные не отсортированы")
			os.Exit(1)
		}
	}

	sort.Slice(lines, func(i, j int) bool {
		return compareLines(lines[i], lines[j], config)
	})

	if config.Unique {
		lines = removeDuplicates(lines)
	}
	for _, line := range lines {
		fmt.Println(line)
	}
}
