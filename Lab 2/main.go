package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// ============================================================
// Задача 1: AVL-дерево
// ============================================================

type AVLNode struct {
	key    int
	height int
	left   *AVLNode
	right  *AVLNode
}

type AVLTree struct {
	root *AVLNode
}

func height(node *AVLNode) int {
	if node == nil {
		return 0
	}
	return node.height
}

func balanceFactor(node *AVLNode) int {
	if node == nil {
		return 0
	}
	return height(node.left) - height(node.right)
}

func updateHeight(node *AVLNode) {
	node.height = 1 + max(height(node.left), height(node.right))
}

func rightRotate(y *AVLNode) *AVLNode {
	x := y.left
	T2 := x.right

	x.right = y
	y.left = T2

	updateHeight(y)
	updateHeight(x)
	return x
}

func leftRotate(x *AVLNode) *AVLNode {
	y := x.right
	T2 := y.left

	y.left = x
	x.right = T2

	updateHeight(x)
	updateHeight(y)
	return y
}

func (t *AVLTree) Insert(key int) {
	t.root = insertNode(t.root, key)
}

func insertNode(node *AVLNode, key int) *AVLNode {
	if node == nil {
		return &AVLNode{key: key, height: 1}
	}

	if key < node.key {
		node.left = insertNode(node.left, key)
	} else if key > node.key {
		node.right = insertNode(node.right, key)
	} else {
		return node
	}

	updateHeight(node)
	balance := balanceFactor(node)

	if balance > 1 && key < node.left.key {
		return rightRotate(node)
	}
	if balance < -1 && key > node.right.key {
		return leftRotate(node)
	}
	if balance > 1 && key > node.left.key {
		node.left = leftRotate(node.left)
		return rightRotate(node)
	}
	if balance < -1 && key < node.right.key {
		node.right = rightRotate(node.right)
		return leftRotate(node)
	}

	return node
}

func (t *AVLTree) Search(key int) bool {
	return searchNode(t.root, key)
}

func searchNode(node *AVLNode, key int) bool {
	if node == nil {
		return false
	}
	if key == node.key {
		return true
	}
	if key < node.key {
		return searchNode(node.left, key)
	}
	return searchNode(node.right, key)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ============================================================
// Задача 1: Тестирование производительности
// ============================================================

func runPerformanceTest() {
	const N = 300000
	rangeLimit := 1000000

	fmt.Printf("\n=== Задача 1. Производительность (N = %d) ===\n", N)
	fmt.Printf("Диапазон чисел: [%d, %d]\n", -rangeLimit, rangeLimit)

	rand.Seed(time.Now().UnixNano())
	slice := make([]int, 0, N)
	avl := &AVLTree{}

	for i := 0; i < N; i++ {
		val := rand.Intn(2*rangeLimit+1) - rangeLimit
		slice = append(slice, val)
		avl.Insert(val)
	}

	secondBatch := make([]int, N)
	for i := 0; i < N; i++ {
		secondBatch[i] = rand.Intn(2*rangeLimit+1) - rangeLimit
	}

	// Линейный поиск в slice
	startSlice := time.Now()
	countSlice := 0
	for _, val := range secondBatch {
		for _, v := range slice {
			if v == val {
				countSlice++
				break
			}
		}
	}
	timeSlice := time.Since(startSlice)

	// Поиск в AVL
	startAvl := time.Now()
	countAvl := 0
	for _, val := range secondBatch {
		if avl.Search(val) {
			countAvl++
		}
	}
	timeAvl := time.Since(startAvl)

	fmt.Printf("\nРезультаты поиска %d чисел:\n", N)
	fmt.Printf("  Slice (линейный поиск): найдено = %d, время = %v\n", countSlice, timeSlice)
	fmt.Printf("  AVL-дерево:            найдено = %d, время = %v\n", countAvl, timeAvl)

	startAvlOnly := time.Now()
	for _, val := range secondBatch {
		avl.Search(val)
	}
	timeAvlOnly := time.Since(startAvlOnly)

	fmt.Printf("\nТолько поиск в AVL-дереве: %v\n", timeAvlOnly)
	fmt.Println("Вывод: AVL-дерево значительно быстрее линейного поиска в slice.")
}

// ============================================================
// Задача 2 и 3: Разбор выражения
// ============================================================

type token struct {
	typ string // "num", "op", "lparen", "rparen", "var"
	val int
	op  byte
}

func tokenize(expr string) []token {
	tokens := []token{}
	i := 0
	n := len(expr)

	for i < n {
		ch := expr[i]
		if unicode.IsSpace(rune(ch)) {
			i++
			continue
		}
		if unicode.IsDigit(rune(ch)) {
			start := i
			for i < n && unicode.IsDigit(rune(expr[i])) {
				i++
			}
			num, _ := strconv.Atoi(expr[start:i])
			tokens = append(tokens, token{typ: "num", val: num})
			continue
		}
		if ch == '(' {
			tokens = append(tokens, token{typ: "lparen"})
			i++
			continue
		}
		if ch == ')' {
			tokens = append(tokens, token{typ: "rparen"})
			i++
			continue
		}
		if ch == '+' || ch == '-' || ch == '*' || ch == '/' {
			tokens = append(tokens, token{typ: "op", op: ch})
			i++
			continue
		}
		if unicode.IsLetter(rune(ch)) {
			tokens = append(tokens, token{typ: "var", op: ch})
			i++
			continue
		}
		i++
	}
	return tokens
}

func precedence(op byte) int {
	switch op {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	}
	return 0
}

// Проверка деления на ноль
func applyOp(a, b int, op byte) int {
	switch op {
	case '+':
		return a + b
	case '-':
		return a - b
	case '*':
		return a * b
	case '/':
		if b == 0 {
			panic(fmt.Sprintf("Ошибка: деление на ноль (%d / %d)", a, b))
		}
		return a / b
	}
	return 0
}

func evaluateTokens(tokens []token, vars map[byte]int) int {
	values := []int{}
	ops := []byte{}

	for _, t := range tokens {
		switch t.typ {
		case "num":
			values = append(values, t.val)
		case "var":
			val, ok := vars[t.op]
			if !ok {
				panic(fmt.Sprintf("Переменная %c не определена", t.op))
			}
			values = append(values, val)
		case "lparen":
			ops = append(ops, '(')
		case "rparen":
			for len(ops) > 0 && ops[len(ops)-1] != '(' {
				op := ops[len(ops)-1]
				ops = ops[:len(ops)-1]
				b := values[len(values)-1]
				a := values[len(values)-2]
				values = values[:len(values)-2]
				values = append(values, applyOp(a, b, op))
			}
			ops = ops[:len(ops)-1]
		case "op":
			for len(ops) > 0 && ops[len(ops)-1] != '(' && precedence(ops[len(ops)-1]) >= precedence(t.op) {
				op := ops[len(ops)-1]
				ops = ops[:len(ops)-1]
				b := values[len(values)-1]
				a := values[len(values)-2]
				values = values[:len(values)-2]
				values = append(values, applyOp(a, b, op))
			}
			ops = append(ops, t.op)
		}
	}

	for len(ops) > 0 {
		op := ops[len(ops)-1]
		ops = ops[:len(ops)-1]
		b := values[len(values)-1]
		a := values[len(values)-2]
		values = values[:len(values)-2]
		values = append(values, applyOp(a, b, op))
	}

	return values[0]
}

func EvaluateExpression(expr string) int {
	tokens := tokenize(expr)
	for _, t := range tokens {
		if t.typ == "var" {
			panic("Обнаружена переменная. Используйте EvaluateWithVariables.")
		}
	}
	return evaluateTokens(tokens, nil)
}

func EvaluateWithVariables(expr string) int {
	tokens := tokenize(expr)

	varSet := make(map[byte]bool)
	for _, t := range tokens {
		if t.typ == "var" {
			varSet[t.op] = true
		}
	}

	vars := make(map[byte]int)
	scanner := bufio.NewScanner(os.Stdin)

	for v := range varSet {
		fmt.Printf("Введите значение переменной %c: ", v)
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())
		val, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("Ошибка: '%s' не является целым числом\n", input)
			for {
				fmt.Printf("Введите ЦЕЛОЕ число для переменной %c: ", v)
				scanner.Scan()
				input = strings.TrimSpace(scanner.Text())
				val, err = strconv.Atoi(input)
				if err == nil {
					break
				}
				fmt.Printf("Ошибка: '%s' не является целым числом\n", input)
			}
		}
		vars[v] = val
	}

	return evaluateTokens(tokens, vars)
}

// ============================================================
// main
// ============================================================

func main() {
	fmt.Println("Лабораторная работа №2")
	fmt.Println("Реализовано: AVL-дерево (сбалансированное дерево поиска)")

	runPerformanceTest()

	fmt.Println("\n=== Задача 2. Разбор выражения (без переменных) ===")
	exprs := []string{"(2+3)*4", "10 + 20 * 2", "(100 - 10) / 3"}
	for _, expr := range exprs {
		fmt.Printf("%s = %d\n", expr, EvaluateExpression(expr))
	}

	fmt.Println("\n=== Задача 3*. Разбор выражения с переменными ===")
	fmt.Println("Пример: a + b * 2")
	result := EvaluateWithVariables("a + b * 2")
	fmt.Printf("Результат: %d\n", result)

	fmt.Println("\nПример с делением (проверка на ноль): (x - y) * 3 / z")
	fmt.Println("При делении на ноль программа выдаст ошибку и завершится")
	result2 := EvaluateWithVariables("(x - y) * 3 / z")
	fmt.Printf("Результат: %d\n", result2)
}
