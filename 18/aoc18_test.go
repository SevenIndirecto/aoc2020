package main

import (
	"testing"
)

type Fixture struct {
	Expression string
	Expected   string
}

func TestStripRedundantParens(t *testing.T) {
	fixtures := []Fixture{
		{"(4 * (5 + 6))", "4 * (5 + 6)"},
		{"(2 * 3) + (4 * (5 + 6))", "(2 * 3) + (4 * (5 + 6))"},
		{"2 * 3", "2 * 3"},
		{"(2 * 3", "(2 * 3"},
	}

	for _, f := range fixtures {
		got := StripRedundantParens(f.Expression)

		if got != f.Expected {
			t.Errorf("Failed to strip got: %s, expected: %s", got, f.Expected)
		}
	}
}

type FixtureEval struct {
	Expression string
	Expected      int
}

func TestNode_Evaluate(t *testing.T) {
	fixtures := []FixtureEval{
		{"1 + 2 * 3", 9},
		{"1 + 2 * 3 + 4 * 5 + 6", 71},
		{"1 + (2 * 3) + (4 * (5 + 6))", 51},
		{"2 * 3 + (4 * 5)", 26},
		{"5 + (8 * 3 + 9 + 3 * 4 * 3)", 437},
		{"5 * 9 * (7 * 3 * 3 + 9 * 3 + (8 + 6 * 4))", 12240},
		{"((2 + 4 * 9) * (6 + 9 * 8 + 6) + 6) + 2 + 4 * 2", 13632},
		{"7 * (4 * 2 + 8 * (6 + 9) * 7 * 6) + (2 + 5 * 5 * 4 + 6 + 9) * 6", 424290},
		{"(1 + 2) * (3 + 4) * (5 + 6)", 231},
		{"2 * (3 + (4 * 5))", 46},
	}

	for _, f := range fixtures {
		stripped := StripLine(f.Expression)
		node := Parse(stripped)
		got := node.Evaluate()

		if got != f.Expected {
			t.Errorf("Failed to evaluate expression(%v) got %d, expected %d", f.Expression, got, f.Expected)
		}
	}
}

type FixturePP struct {
	Expression string
	Expected string
}

func TestAddPrecedenceParensSinglePass(t *testing.T) {
	fixtures := []FixturePP{
		{"1 + 2 * 3 + 4 * 5 + 6", "(1 + 2 )* 3 + 4 * 5 + 6"},
		{"(1 + 2) * 3 + 4 * 5 + 6", "(1 + 2) *( 3 + 4 )* 5 + 6"},
		{"2 * 3 + (4 * 5)", "2 *( 3 + (4 * 5))"},
	}

	for _, f := range fixtures {
		got := AddPrecedenceParensSinglePass(f.Expression)

		if got != f.Expected {
			t.Errorf("Invalid precedence prefill for [%s] got [%s] expected [%s]", f.Expression, got, f.Expected)
		}
	}
}

func TestAddPrecedenceParens(t *testing.T) {
	fixtures := []FixturePP{
		{"1 + 2 * 3 + 4 * 5 + 6", "(1 + 2 )*( 3 + 4 )*( 5 + 6)"},
		{"2 * 3 + (4 * 5)", "2 *( 3 + (4 * 5))"},
	}

	for _, f := range fixtures {
		got := AddPrecedenceParens(f.Expression)

		if got != f.Expected {
			t.Errorf("Invalid precedence prefill for [%s] got [%s] expected [%s]", f.Expression, got, f.Expected)
		}
	}
}

func TestNode_EvaluateUsingPrecedence(t *testing.T) {
	fixtures := []FixtureEval{
		{"1 + 2 * 3 + 4 * 5 + 6", 231},
		{"1 + (2 * 3) + (4 * (5 + 6))", 51},
		{"2 * 3 + (4 * 5)", 46},
		{"5 + (8 * 3 + 9 + 3 * 4 * 3)", 1445},
		{"5 * 9 * (7 * 3 * 3 + 9 * 3 + (8 + 6 * 4))", 669060},
		{"((2 + 4 * 9) * (6 + 9 * 8 + 6) + 6) + 2 + 4 * 2", 23340},
	}

	for _, f := range fixtures {
		stripped := StripLine(f.Expression)
		withPrecedence := AddPrecedenceParens(stripped)
		node := Parse(withPrecedence)
		got := node.Evaluate()

		if got != f.Expected {
			t.Errorf("Failed to evaluate expression[%v] got %d, expected %d", f.Expression, got, f.Expected)
		}
	}
}
