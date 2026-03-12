package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Задание 1: Пересечение мультимножеств
func Intersect(a, b []int) []int {
	if len(a) > len(b) {
		a, b = b, a
	}

	// Подсчитываем частоту элементов в меньшем списке a
	counts := make(map[int]int, len(a))
	for _, v := range a {
		counts[v]++
	}

	// Проходим по большему списку b и собираем пересечение
	result := make([]int, 0, len(a))
	for _, v := range b {
		if cnt, exists := counts[v]; exists && cnt > 0 {
			result = append(result, v)
			counts[v]-- // уменьшаем счетчик, чтобы учесть кратность
		}
	}

	return result
}

// Функция для демонстрации Задания 1
func Task1Demo() {
	fmt.Println("=== Задание 1: Пересечение мультимножеств ===")

	A := []int{1, 2, 2, 2, 3, 4, 4, 4, 5}
	B := []int{2, 2, 4, 4, 6, 7, 8}

	C := Intersect(A, B)
	fmt.Println("A =", A)
	fmt.Println("B =", B)
	fmt.Println("Пересечение:", C)

	D := []int{3}
	E := []int{1, 1, 2, 2, 4}
	F := Intersect(D, E)
	fmt.Println("\nD =", D)
	fmt.Println("E =", E)
	fmt.Println("Пересечение:", F)
	fmt.Println()
}

// TextToWords преобразует текст в список слов
// Слова состоят только из букв, разделители - пробелы, знаки препинания и концы строк
func TextToWords(text string) []string {
	// Приводим текст к нижнему регистру
	text = strings.ToLower(text)

	// Заменяем разделители на пробелы
	separators := []rune{' ', ',', '.', ';', '!', '?', ':', '\n', '\r', '\t'}
	result := make([]string, 0)

	// Проходим по тексту и собираем слова
	var currentWord strings.Builder

	for _, ch := range text {
		isSeparator := false
		for _, sep := range separators {
			if ch == sep {
				isSeparator = true
				break
			}
		}

		if isSeparator {
			// Если накопили слово, добавляем его
			if currentWord.Len() > 0 {
				result = append(result, currentWord.String())
				currentWord.Reset()
			}
		} else {
			// Проверяем, является ли символ буквой
			if (ch >= 'a' && ch <= 'z') || (ch >= 'а' && ch <= 'я') {
				currentWord.WriteRune(ch)
			}
		}
	}

	// Добавляем последнее слово, если есть
	if currentWord.Len() > 0 {
		result = append(result, currentWord.String())
	}

	return result
}

// BuildBigramDictionary строит словарь биграмм из списка слов
func BuildBigramDictionary(words []string) map[string][]string {
	dict := make(map[string][]string)

	// Проходим по всем парам соседних слов
	for i := 0; i < len(words)-1; i++ {
		firstWord := words[i]
		secondWord := words[i+1]

		// Добавляем второе слово в список для первого слова
		dict[firstWord] = append(dict[firstWord], secondWord)
	}

	return dict
}

// Autocomplete выполняет автодополнение для введенного слова
func Autocomplete(startWord string, dict map[string][]string, rng *rand.Rand) {
	currentWord := startWord

	// Проверяем, есть ли слово в словаре
	if _, exists := dict[currentWord]; !exists {
		fmt.Printf("Слово '%s' не найдено в словаре\n", currentWord)
		return
	}

	// Строим предложение
	sentence := []string{currentWord}

	// Генерируем продолжения, пока они есть
	for i := 0; i < 3; i++ {
		possibleNext, exists := dict[currentWord]
		if !exists || len(possibleNext) == 0 {
			break // Дальше продолжить нельзя
		}

		// Используем локальный генератор для случайного выбора
		nextWord := possibleNext[rng.Intn(len(possibleNext))]
		sentence = append(sentence, nextWord)
		currentWord = nextWord
	}

	fmt.Println(strings.Join(sentence, " "))
}

// Task2Demo демонстрирует работу интеллектуального помощника ввода
func Task2Demo() {
	fmt.Println("=== Задание 2: Интеллектуальный помощник ввода ===")

	// Пример текста для демонстрации
	exampleText := "We study programming languages C++, C#, Go. We are programmers!"

	fmt.Println("Исходный текст:", exampleText)
	fmt.Println()

	// 2.1 Преобразование текста в список слов
	words := TextToWords(exampleText)
	fmt.Println("2.1 Список слов:", words)
	fmt.Println()

	// 2.2 Построение словаря биграмм
	bigramDict := BuildBigramDictionary(words)
	fmt.Println("2.2 Словарь биграмм:")
	for key, values := range bigramDict {
		fmt.Printf("  \"%s\" -> %v\n", key, values)
	}
	fmt.Println()

	// 2.3 Автодополнение
	fmt.Println("2.3 Автодополнение (демонстрация)")
	fmt.Println("Введите слово для начала (или 'exit' для выхода):")

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		input := strings.TrimSpace(strings.ToLower(scanner.Text()))

		if input == "exit" {
			break
		}

		Autocomplete(input, bigramDict, rng)
		fmt.Println()
	}
}

// HashTable представляет словарь с хеш-таблицей
type HashTable struct {
	buckets [][]string // Массив бакетов для разрешения коллизий методом цепочек
	size    int        // Размер таблицы (количество бакетов)
}

// NewHashTable создает новую хеш-таблицу заданного размера
func NewHashTable(size int) *HashTable {
	buckets := make([][]string, size)
	for i := range buckets {
		buckets[i] = make([]string, 0)
	}
	return &HashTable{
		buckets: buckets,
		size:    size,
	}
}

// HashFunction - пользовательская хеш-функция для слов на естественном языке
// Учитывает длину слова, позицию символов (первые важнее), гласные/согласные
func (ht *HashTable) HashFunction(word string) int {
	word = strings.ToLower(word)
	if len(word) == 0 {
		return 0
	}

	// 1. Учитываем длину слова
	lengthFactor := len(word) * 7

	// 2. Учитываем сумму кодов символов с весами
	sum := 0
	vowels := "aeiouyаеёиоуыэюя"

	for i, ch := range word {
		// Вес убывает с позицией: первые символы важнее
		posWeight := len(word) - i

		isVowel := false
		for _, v := range vowels {
			if ch == v {
				isVowel = true
				break
			}
		}

		// Гласные получают больший вес
		charWeight := 3
		if isVowel {
			charWeight = 5
		}

		sum += int(ch) * posWeight * charWeight
	}

	// 3. Учитываем первый и последний символы отдельно
	firstCharWeight := int(word[0]) * 11
	lastCharWeight := int(word[len(word)-1]) * 13

	// Комбинируем все факторы
	result := (sum*31 + lengthFactor*17 + firstCharWeight + lastCharWeight) % ht.size
	if result < 0 {
		result = -result
	}
	return result
}

// Add добавляет слово в словарь
func (ht *HashTable) Add(word string) {
	word = strings.ToLower(strings.TrimSpace(word))
	if word == "" {
		return
	}

	// Проверяем, есть ли уже такое слово
	if ht.Check(word) {
		return // Слово уже есть, не добавляем повторно
	}

	// Вычисляем хеш
	hash := ht.HashFunction(word)

	// Добавляем слово в бакет
	ht.buckets[hash] = append(ht.buckets[hash], word)
}

// Check проверяет наличие слова в словаре
func (ht *HashTable) Check(word string) bool {
	word = strings.ToLower(strings.TrimSpace(word))
	if word == "" {
		return false
	}

	hash := ht.HashFunction(word)

	// Ищем слово в бакете
	for _, w := range ht.buckets[hash] {
		if w == word {
			return true
		}
	}
	return false
}

// Task3 - функция для запуска задания 3.1
func Task3Demo() {
	fmt.Println("=== Задание 3.1: Интерактивный режим ===")
	fmt.Println("Команды: add <слово> или check <слово> (или 'exit' для выхода)")
	fmt.Println()

	// Создаем хеш-таблицу размером 100 бакетов
	ht := NewHashTable(100)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())

		if input == "exit" {
			break
		}

		parts := strings.Fields(input)
		if len(parts) < 2 {
			fmt.Println("Неверный формат команды. Используйте: add <слово> или check <слово>")
			continue
		}

		command := parts[0]
		word := parts[1]

		switch command {
		case "add":
			ht.Add(word)
			fmt.Printf("Слово '%s' добавлено\n", word)
		case "check":
			if ht.Check(word) {
				fmt.Println("yes")
			} else {
				fmt.Println("no")
			}
		default:
			fmt.Println("Неизвестная команда. Используйте add или check")
		}
	}
}

func main() {
	Task1Demo()
	Task2Demo()
	//  Слова с разным регистром считаются одинаковыми
	Task3Demo()
}
