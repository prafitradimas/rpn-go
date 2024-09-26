package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	expr := "2 + 3 * 2 / (1 - 5)"

	tokens := tokenize(strings.TrimSpace(expr))
	fmt.Printf("tokens: %+v\n", tokens)

	tok := parse(tokens)
	fmt.Printf("tok: %v\n", tok)

	fmt.Printf("result: %d\n", eval(tok))
}

type Token = string

func prec(s string) int {
	switch s {
	case "+":
		return 2
	case "-":
		return 2
	case "*":
		return 3
	case "/":
		return 3
	default:
		return 0
	}
}

func isOp(s byte) bool {
	switch s {
	case '+':
		return true
	case '-':
		return true
	case '*':
		return true
	case '/':
		return true
	default:
		return false
	}
}

func tokenize(expr string) []Token {
	t := []Token{}

	pos := 0
	for len(expr) > pos {
		if isOp(expr[pos]) {
			t = append(t, string(expr[pos]))
			pos += 1
			continue
		}

		if unicode.IsDigit(rune(expr[pos])) {
			ptr := pos
			for unicode.IsDigit(rune(expr[ptr])) {
				ptr += 1
			}

			t = append(t, expr[pos:ptr])
			pos = ptr
		}

		t = append(t, string(expr[pos]))
		pos += 1
	}

	return t
}

func parse(tokens []Token) []Token {
	out := []Token{}
	stack := []Token{}

	for _, token := range tokens {
		if isOp(token[0]) {
			for len(stack) > 0 && prec(stack[len(stack)-1]) >= prec(token) {
				out = append(out, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}

			stack = append(stack, token)
		} else if token == "(" {
			stack = append(stack, token)
		} else if token == ")" {
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				out = append(out, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) > 0 && stack[len(stack)-1] == "(" {
				stack = stack[:len(stack)-1]
			}
		} else {
			out = append(out, token)
		}
	}

	for len(stack) > 0 {
		out = append(out, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return out
}

func eval(tokens []Token) int {
	stack := []int{}
	for _, token := range tokens {
		if isOp(token[0]) {
			a := stack[len(stack)-2]
			b := stack[len(stack)-1]

			var res int
			switch token {
			case "+":
				res = a + b
			case "-":
				res = a - b
			case "*":
				res = a * b
			case "/":
				if b > a {
					log.Fatalf("ERR: #DIV/0\n")
				}
				res = a / b
			default:
				res = 0
			}

			stack = stack[:len(stack)-1]
			stack[len(stack)-1] = res
		} else {
			num, err := strconv.Atoi(token)
			if err == nil {
				stack = append(stack, num)
			}
		}
	}

	return stack[0]
}
