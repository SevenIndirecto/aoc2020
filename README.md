# Advent of Code 2020 - golang

## Initial setup

First create your .env file by copying .env.example

```bash
cp .env.example .env
```

1. Login on https://adventofcode.com/ and grab your session cookie ID.
2. Input into .env file


## Start new day

```bash
./next.sh
```

This prepares a new dir, copies templates and downloads your input.

To run tests in a dir.

```bash
go test -run ''
```

To run program in a dir obviously:

```bash
got run aoc01.go // etc
```


## Advent of Code 2020 - Closing Thoughts

**Golang** was definitely a pleasure to write in. Sure I missed some quick utils languages like python offer,
but those are usually 3 lines long loops anyway.

**Unit testing** with Golang was a breeze and made this year's approach much less error prone.
Static typing also certainly helped compared to previous year's python and javascript endeavours.

Compared to last year I've been coding in **vim** (mostly using Intellij Idea + IdeaVim) for the last year so not much 
vim-fu was learned during AoC. That being said still very happy I took the time to learn vim and see myself continuing 
to use it daily.

As for **Advent of Code 2020** itself, I guess this year's aoc was likely intentionally easier than the prior 2 years, 
but I liked the fact it took less time in a busy work day. I'm usually trying to learn something when doing these,
be it last year's vim/python adventures, this year's golang or just doing AoC for the first time in javascript.

As for the repeatable "IntCode" style challenges from last year? Not sure if they're bad or good, as I do 
these challenges in sequence anyway, but isolation and not having to care about keeping compatibility with previous 
challenges is welcome.

I overcomplicated day 20 for sure, had to check for hints for day 23 (Crab Cups!) as I started looking for patterns
and loops way too soon, not realizing the numbers for part 2 are quite computable.

Thanks **[Eric Wastl](http://was.tl/)** for making [Advent of Code](https://adventofcode.com/2020/)
