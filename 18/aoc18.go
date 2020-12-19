package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	MUL    = 1
	ADD    = 2
	NUMBER = 3
)

type Node struct {
	Left  *Node
	Right *Node
	Value int
	Type  int
}

func (n *Node) String() string {
	var typeStr string
	switch n.Type {
	case MUL:
		typeStr = "MUL"
	case ADD:
		typeStr = "ADD"
	case NUMBER:
		typeStr = "LEAF"
	}
	return fmt.Sprintf("\nN[Type: %s, Value: %d, Left: %v, Right: %v]", typeStr, n.Value, n.Left, n.Right)
}

func (n *Node) Evaluate() int {
	if n.Type == NUMBER {
		return n.Value
	}

	if n.Type == MUL {
		return n.Left.Evaluate() * n.Right.Evaluate()
	} else if n.Type == ADD {
		return n.Left.Evaluate() + n.Right.Evaluate()
	} else {
		panic(fmt.Sprintf("Unexpected type %v", n))
	}
}

func StripRedundantParens(s string) string {
	prev := ""
	current := s

	for ; prev != current; {
		prev = current
		current = StripRedundantParensSinglePass(prev)
	}
	return current
}

func StripRedundantParensSinglePass(s string) string {
	if len(s) < 2 || string(s[0]) != "(" || string(s[len(s)-1]) != ")" {
		return s
	}
	// strip opening brace
	balance := 0
	for _, c := range s[1 : len(s)-1] {
		char := string(c)
		if char == "(" {
			balance++
		} else if char == ")" {
			balance--
		}
		if balance < 0 {
			return s
		}
	}

	return s[1 : len(s)-1]
}

func Parse(expr string) Node {
	leftExpr := ""
	rightExpr := ""
	nodeType := 0
	var parenthesesBalance = 0

	expr = StripRedundantParens(expr)
	// Read until OP found
	size := len(expr)
forLoop:
	for i := range expr {
		// Loop backwards
		index := size - i - 1
		char := string(expr[index])
		switch char {
		case "(":
			parenthesesBalance++
		case ")":
			parenthesesBalance--
		case "+", "*":
			if parenthesesBalance == 0 {
				if char == "+" {
					nodeType = ADD
				} else {
					nodeType = MUL
				}
				leftExpr = expr[:index]
				break forLoop
			}
		}
		rightExpr = char + rightExpr
	}

	// Leaf node
	if nodeType == 0 {
		val, _ := strconv.Atoi(rightExpr)
		return Node{Type: NUMBER, Value: val}
	}

	leftNode := Parse(leftExpr)
	rightNode := Parse(rightExpr)
	return Node{
		Type:  nodeType,
		Left:  &leftNode,
		Right: &rightNode,
	}
}

func AddPrecedenceParens(s string) string {
	prev := ""
	singlePass := s
	for ; singlePass != prev; {
		prev = singlePass
		singlePass = AddPrecedenceParensSinglePass(prev)
	}
	return singlePass
}

func AddPrecedenceParensSinglePass(s string) string {
	leftMulSignOnBalance := make(map[int]int)
	balance := 0

	for i, c := range s {
		char := string(c)
		switch char {
		case "(":
			balance++
			// Check if we need to unset mul sign tracking
			if _, exists := leftMulSignOnBalance[balance]; exists {
				delete(leftMulSignOnBalance, balance)
			}
		case ")":
			balance--
		case "*":
			leftMulSignOnBalance[balance] = i
		case "+":
			// Find first left mul on same level
			leftMulIndex, leftExists := leftMulSignOnBalance[balance]

			// Find first right mul on same level
			rightSearchBalance := balance
			rightMulIndex := -1

		rightSearchForLoop:
			for j := i + 1; j < len(s); j++ {
				searchChar := string(s[j])
				switch searchChar {
				case "(":
					rightSearchBalance++
					if rightSearchBalance == balance {
						// We exited and re-entered same balance level, fail to find
						break rightSearchForLoop
					}
				case ")":
					rightSearchBalance--
				case "*":
					if balance == rightSearchBalance {
						// Found candidate
						rightMulIndex = j
						break rightSearchForLoop
					}
				}
			}

			// Now see if we should add any parens
			if leftExists && rightMulIndex > -1 {
				// Found valid MUL on both sides
				return s[0:leftMulIndex+1] + "(" + s[leftMulIndex+1:rightMulIndex] + ")" + s[rightMulIndex:]
			} else if leftExists {
				// Only found valid MUL on left side
				rightPos := len(s)
				rightPosBalance := balance
				for j := i + 1; j < len(s); j++ {
					char := string(s[j])
					if char == "(" {
						rightPosBalance++
					} else if char == ")" {
						if rightPosBalance == balance {
							rightPos = j
							break
						} else {
							rightPosBalance--
						}
					} else if char == "+" && rightPosBalance == balance {
						rightPos = j
						break
					}
				}
				return s[0:leftMulIndex+1] + "(" + s[leftMulIndex+1:rightPos] + ")" + s[rightPos:]
			} else if rightMulIndex > -1 {
				// Only found valid MUL on right side

				// Find first opening param or + sign to the left on the same level
				leftPos := -1
				leftSearchBalance := balance
				for j := i - 1; j >= 0; j-- {
					char := string(s[j])
					if char == ")" {
						leftSearchBalance--
					} else if char == "(" {
						if leftSearchBalance == balance {
							leftPos = j
							break
						} else {
							leftSearchBalance++
						}
					} else if char == "+" && leftSearchBalance == balance {
						leftPos = j
						break
					}
				}
				return s[0:leftPos+1] + "(" + s[leftPos+1:rightMulIndex] + ")" + s[rightMulIndex:]
			}
		}
	}
	return s
}

func StripLine(line string) string {
	return strings.ReplaceAll(line, " ", "")
}

func main() {
	dat, err := ioutil.ReadFile("aoc18.txt")
	if err != nil {
		panic(err)
	}

	txt := string(dat)
	lines := strings.Split(txt, "\n")

	sum := 0
	for _, line := range lines {
		if len(line) < 2 {
			continue
		}

		node := Parse(StripLine(line))
		sum += node.Evaluate()
	}
	fmt.Println("Part one:", sum)

	sum = 0
	for _, line := range lines {
		if len(line) < 2 {
			continue
		}

		stripped := StripLine(line)
		withPrecedence := AddPrecedenceParens(stripped)
		node := Parse(withPrecedence)
		sum += node.Evaluate()
	}
	fmt.Println("Part two", sum)
}
