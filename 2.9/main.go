package main

import (
	"fmt"
	"strconv"
	"strings"
)

func convertString(str string) (string, error) {
	if len(str) == 0 { // пустая строка
		return "", nil
	}

	var resultStr string
	var hasLetters bool // флаг для проверки наличия букв

	for i := 0; i < len(str); i++ {
		if str[i] >= 'a' && str[i] <= 'z' {
			hasLetters = true

			if i < len(str)-1 && str[i+1] >= '0' && str[i+1] <= '9' {
				// извлекаем последовательные цифры
				j := i + 1
				for j < len(str) && str[j] >= '0' && str[j] <= '9' {
					j++
				}

				count, err := strconv.Atoi(str[i+1 : j])
				if err != nil {
					return "", fmt.Errorf("ошибка преобразования числа: %v", err)
				}

				resultStr += strings.Repeat(string(str[i]), count)
				i = j - 1 // пропускаем цифры, которые уже обработали
			} else {
				// на случай если строка оканчивается буквой
				resultStr += string(str[i])
			}
		} else if str[i] >= '0' && str[i] <= '9' {
			// цифра идет первой или после другой цифры
			if i == 0 || !hasLetters {
				return "", fmt.Errorf("некорректная строка: цифра без предшествующей буквы")
			}
		} else {
			// Обработка других символов
			hasLetters = true

			// Проверяем, есть ли следующий символ и является ли он цифрой
			if i < len(str)-1 && str[i+1] >= '0' && str[i+1] <= '9' {
				j := i + 1
				for j < len(str) && str[j] >= '0' && str[j] <= '9' {
					j++
				}

				count, err := strconv.Atoi(str[i+1 : j])
				if err != nil {
					return "", fmt.Errorf("ошибка преобразования числа: %v", err)
				}

				resultStr += strings.Repeat(string(str[i]), count)
				i = j - 1
			} else {
				resultStr += string(str[i])
			}
		}
	}

	if !hasLetters {
		return "", fmt.Errorf("некорректная строка, т.к. в строке только цифры")
	}

	return resultStr, nil
}

func main() {
	testCases := []string{
		"a4bc2d5e",
		"abcd",
		"45",
		"",
		"a1b2c3",
		"a10b2",
	}

	for _, test := range testCases {
		result, err := convertString(test)
		if err != nil {
			fmt.Printf("Вход: %q\nВыход: Ошибка - %v\n\n", test, err)
		} else {
			fmt.Printf("Вход: %q\nВыход: %q\n\n", test, result)
		}
	}
}
