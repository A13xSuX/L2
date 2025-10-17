//test STDIN:  "start", "PATTERN1", "middle", "PATTERN2", "end" | go run grep.go -A 1 "PATTERN"
//test File: go run grep.go -C 1 -n "PATTERN" test.txt

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Config struct {
	After      int
	Before     int
	Context    int
	Count      bool
	IgnoreCase bool
	Invert     bool
	Fixed      bool
	LineNum    bool
	Pattern    string
	FileName   string
}

// Line for correct -n with buffer
type Line struct {
	Number int
	Text   string
}

// scanner - source of lines, matcher - pattern matching function
func (config Config) processWithContext(scanner *bufio.Scanner, matcher func(string) bool) {
	beforeBuffer := make([]Line, 0, config.Before)
	afterCount := 0
	lineNum := 0
	cnt := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		// -A
		if config.After > 0 && afterCount > 0 {
			if config.LineNum { // -n
				fmt.Printf("%d\t%s\n", lineNum, line)
			} else {
				fmt.Println(line)
			}
			afterCount--
			beforeBuffer = beforeBuffer[:0]
			continue
		}

		matched := matcher(line)
		shouldPrint := matched
		// -v
		if config.Invert {
			shouldPrint = !shouldPrint
		}

		if shouldPrint {
			cnt++
			if config.Before > 0 { // -B
				start := len(beforeBuffer) - config.Before
				if start < 0 {
					start = 0
				}
				for _, bufLine := range beforeBuffer[start:] {
					if config.LineNum { // -n
						fmt.Printf("%d\t%s\n", bufLine.Number, bufLine.Text)
					} else {
						fmt.Println(bufLine.Text)
					}
				}
			}
			if config.LineNum { // -n
				fmt.Printf("%d\t%s\n", lineNum, line)
			} else {
				fmt.Println(line)
			}
			beforeBuffer = beforeBuffer[:0]
			afterCount = config.After

		} else {
			if afterCount == 0 {
				if len(beforeBuffer) > config.Before {
					beforeBuffer = beforeBuffer[1:] // for store last string in buffer
				}
				beforeBuffer = append(beforeBuffer, Line{Number: lineNum, Text: line})
			}
		}
	}
	if config.Count { // -c
		fmt.Println(cnt)
	}
}

func matcher(config Config) func(string) bool {
	if config.Fixed { // -F
		if config.IgnoreCase { // -i
			return func(line string) bool {
				return strings.Contains(strings.ToLower(line), strings.ToLower(config.Pattern))
			}
		} else {
			return func(line string) bool {
				return strings.Contains(line, config.Pattern)
			}
		}
	} else { // pattern
		var re *regexp.Regexp
		var err error

		if config.IgnoreCase { //-i
			re, err = regexp.Compile("(?i)" + config.Pattern)
		} else {
			re, err = regexp.Compile(config.Pattern)
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return func(line string) bool {
			return re.MatchString(line)
		}
	}

}

func main() {
	config := Config{}
	flag.IntVar(&config.After, "A", 0, "Вывести N строк после неё ")
	flag.IntVar(&config.Before, "B", 0, "Вывести N строк до неё ")
	flag.IntVar(&config.Context, "C", 0, "Вывести N строк после неё и до нее")
	flag.BoolVar(&config.Count, "c", false, "выводить только то количество строк, что совпадающих с шаблоном")
	flag.BoolVar(&config.IgnoreCase, "i", false, "игнорировать регистр")
	flag.BoolVar(&config.Invert, "v", false, "инвертировать фильтр: выводить строки, не содержащие шаблон")
	flag.BoolVar(&config.Fixed, "F", false, "воспринимать шаблон как фиксированную строку")
	flag.BoolVar(&config.LineNum, "n", false, "выводить номер строки перед каждой найденной строкой")

	flag.Parse()

	// -C
	if config.Context > 0 {
		config.After = config.Context
		config.Before = config.Context
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Не переданы аргументы")
		return
	} else if len(args) == 1 {
		config.Pattern = args[0]
	} else if len(args) == 2 {
		config.Pattern = args[0]
		config.FileName = args[1]
	} else if len(args) > 2 {
		fmt.Println("Слишком много аргументов")
		return
	}
	if config.After < 0 || config.Before < 0 || config.Context < 0 {
		fmt.Println("Кол-во строк не может быть отрицательным в -A -B -C")
		return
	}

	if config.FileName != "" { //FILE
		matcher := matcher(config)
		file, err := os.Open(config.FileName)
		if err != nil {
			fmt.Println("File not found")
			os.Exit(1)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		config.processWithContext(scanner, matcher)
	} else { //STDIN
		matcher := matcher(config)
		scanner := bufio.NewScanner(os.Stdin)

		config.processWithContext(scanner, matcher)

	}
}
