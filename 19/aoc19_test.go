package main

import (
	"testing"
)

const INPUT_1 = `0: 1 2
1: "a"
2: 1 3 | 3 1
3: "b"
`

const INPUT_2 = `0: 4 1 5
1: 2 3 | 3 2
2: 4 4 | 5 5
3: 4 5 | 5 4
4: "a"
5: "b"

abab
aba
`

type Fixture struct {
	RuleId int
	PatternId int
	Matches bool
}

func TestMatcher_PatternMatchesRule(t *testing.T) {
	input := `1: 2 3 | 3 2
0: 4 1 5
2: 4 4 | 5 5
3: 4 5 | 5 4
4: "a"
5: "b"

ababbb
bababa
abbbab
aaabbb
aaaabbb
`
	fixtures := []Fixture{
		{0, 0, true},
		{0, 1, false},
		{0, 2, true},
		{0, 3, false},
		{0, 4, false},
	}

	m := NewMatcher(input)

	for _, f := range fixtures {
		got := m.PatternMatchesRule(f.PatternId, f.RuleId)

		if got != f.Matches {
			t.Errorf("Mismatched pattern %s got %v expected %v", m.Patterns[f.PatternId], got, f.Matches)
		}
	}
}

func TestGetMatchCount(t *testing.T) {
	input := `42: 9 14 | 10 1
9: 14 27 | 1 26
10: 23 14 | 28 1
1: "a"
11: 42 31
5: 1 14 | 15 1
19: 14 1 | 14 14
12: 24 14 | 19 1
16: 15 1 | 14 14
31: 14 17 | 1 13
6: 14 14 | 1 14
2: 1 24 | 14 4
0: 8 11
13: 14 3 | 1 12
15: 1 | 14
17: 14 2 | 1 7
23: 25 1 | 22 14
28: 16 1
4: 1 1
20: 14 14 | 1 15
3: 5 14 | 16 1
27: 1 6 | 14 18
14: "b"
21: 14 1 | 1 14
25: 1 1 | 1 14
22: 14 14
8: 42
26: 14 22 | 1 20
18: 15 15
7: 14 5 | 1 21
24: 14 1

aaaaabbaabaaaaa
abbbbbabbbaaaababbaabbbbabababbbabbbbbbabaaaa
bbabbbbaabaabba
babbbbaabbbbbabbbbbbaabaaabaaa
aaabbbbbbaaaabaababaabababbabaaabbababababaaa
bbbbbbbaaaabbbbaaabbabaaa
bbbababbbbaaaaaaaabbababaaababaabab
ababaaaaaabaaab
ababaaaaabbbaba
baabbaaaabbaaaababbaababb
abbbbabbbbaaaababbbbbbaaaababb
aaaaabbaabaaaaababaa
aaaabbaaaabbaaa
aaaabbaabbaaaaaaabbbabbbaaabbaabaaa
babaaabbbaaabaababbaabababaaab
aabbbbbaabbbaaaaaabbbbbababaaaaabbaaabba
`

	got := GetMatchCount(input, true, 8)
	expected := 12

	if got != expected {
		t.Errorf("Failed PATCHED got %d expected %d", got, expected)
	}
}

// 1: 4 | 4 1
func TestSmallTail(t *testing.T) {
	input := `0: 1 2 5
1: 4
2: 4 4 | 5 5
4: "a"
5: "b"

aabbb
aaabbb
aaaaaaaaaabbb
aaaaaaaaaaaaaaaaaaaabbb
`
	got := GetMatchCount(input, true, 1)
	expected := 3

	if got != expected {
		t.Errorf("Failed got %d expected %d", got, expected)
	}
}

func TestSmallMid(t *testing.T) {
	input := `0: 1 2
1: 4 5 | 4 1 5
2: 4 4 | 5 5
4: "a"
5: "b"

abbb
aabbbb
aaabbbaa
aaabbba
`
	got := GetMatchCount(input, true, 1)
	expected := 3

	if got != expected {
		t.Errorf("Failed got %d expected %d", got, expected)
	}
}
