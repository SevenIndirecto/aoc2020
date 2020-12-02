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
