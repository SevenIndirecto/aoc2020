package main

import (
	"strings"
	"testing"
)

type Fixture struct {
	path string
	expected Point
}

func TestGetTilePoint(t *testing.T) {
	fixtures := []Fixture{
		{"esew", Point{x: 1, y: -1}},
		{"nwwswee", Point{x: 0, y: 0}},
		{"wneneese", Point{x: 3, y: 1}},
	}

	for _, f := range fixtures {
		got := GetTilePoint(f.path, Point{0, 0})

		if got != f.expected {
			t.Errorf("Path(%s) got %v expected %v", f.path, got, f.expected)
		}
	}
}

const PRESET = `sesenwnenenewseeswwswswwnenewsewsw
neeenesenwnwwswnenewnwwsewnenwseswesw
seswneswswsenwwnwse
nwnwneseeswswnenewneswwnewseswneseene
swweswneswnenwsewnwneneseenw
eesenwseswswnenwswnwnwsewwnwsene
sewnenenenesenwsewnenwwwse
wenwwweseeeweswwwnwwe
wsweesenenewnwwnwsenewsenwwsesesenwne
neeswseenwwswnwswswnw
nenwswwsewswnenenewsenwsenwnesesenew
enewnwewneswsewnwswenweswnenwsenwsw
sweneswneswneneenwnewenewwneswswnese
swwesenesewenwneswnwwneseswwne
enesenwswwswneneswsenwnewswseenwsese
wnwnesenesenenwwnenwsewesewsesesew
nenewswnwewswnenesenwnesewesw
eneswnwswnwsenenwnwnwwseeswneewsenese
neswnwewnwnwseenwseesewsenwsweewe
wseweeenwnesenwwwswnew
`

func TestTileSet_Paint(t *testing.T) {
	expected := 10
	lines := strings.Split(PRESET, "\n")
	ts := NewTileSet()
	ts.Paint(lines)
	got := ts.GetBlackCount()

	if got != expected {
		t.Errorf("Failed to paint got %d, expected %d", got, expected)
	}
}

type FixtureDaily struct {
	days int
	expected int
}

func TestTileSet_ExecuteDailyPaints(t *testing.T) {
	fixtures := []FixtureDaily{
		{1, 15},
		{10, 37},
		{100, 2208},
	}

	lines := strings.Split(PRESET, "\n")

	for _, f := range fixtures {
		ts := NewTileSet()
		ts.Paint(lines)
		ts.ExecuteDailyPaints(f.days)
		got := ts.GetBlackCount()

		if got != f.expected {
			t.Errorf("Days(%d) got %d expected %d", f.days, got, f.expected)
		}
	}
}
