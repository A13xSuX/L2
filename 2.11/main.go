package main

import (
	"fmt"
	"slices"
	"strings"
)

// sugar
func toLow(str string) string {
	return strings.ToLower(str)
}

func removeDuplicatesAndSort(s []string) []string {
	seen := make(map[string]bool)
	unique := []string{}
	for i := range s {
		if !seen[s[i]] {
			seen[s[i]] = true
			unique = append(unique, s[i])
		}
	}
	slices.Sort(unique)
	return unique
}

func searchAnagram(str []string) map[string][]string {
	result := make(map[string][]string)
	part := make(map[string][]string) //промежуточная
	for i := 0; i < len(str); i++ {
		forSort := []rune(toLow(str[i])) //способ создания ключа
		slices.Sort(forSort)
		key := string(forSort)
		part[key] = append(part[key], toLow(str[i])) //пока что ключ сортирован
	}
	for _, v := range part {
		if len(v) <= 1 {
			continue
		}
		uniqueSorted := removeDuplicatesAndSort(v)
		firstWord := uniqueSorted[0]
		result[firstWord] = uniqueSorted
	}
	return result
}

func main() {
	str := []string{"пятак", "пЯтка", "тяпка", "листок", "слиток", "столик", "стол"}
	fmt.Println(searchAnagram(str))
}
