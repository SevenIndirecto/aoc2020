#!/usr/bin/env bash

set -euo

# import .env
set -o allexport; source .env; set +o allexport
if [[ -z $YEAR || -z $SESSION ]]; then
  echo "Set YEAR and SESSION in an .env file"
  exit
fi

if [[ ! -d 'template' ]]; then
  echo "Make sure you're in the root directory and template dir exists".
  exit
fi

# get previous day dir name
PREV_DAY=$(ls -d */ | grep -E "[0-9]{2}" | sed 's#/##' | tail -n1)

if [[ -z $PREV_DAY ]]; then
  NEXT_NUM=1
else
  # 10# to force base 10, since our format can be "08" and interpreted as octal
  NEXT_NUM=$((10#${PREV_DAY}+1))
fi

NEXT=$(printf "%02d" "$NEXT_NUM")
mkdir "$NEXT"
NEW_GO_MAIN_FILE="$NEXT/aoc$NEXT.go"
cp template/main.go "$NEW_GO_MAIN_FILE"
cp template/main_test.go "$NEXT/aoc${NEXT}_test.go"

INPUT_FILE="aoc${NEXT}.txt"
sed -i "s/INPUT_FILE/${INPUT_FILE}/g" "$NEW_GO_MAIN_FILE"

curl "https://adventofcode.com/$YEAR/day/$NEXT_NUM/input" -H "cookie: session=$SESSION" > "$NEXT/$INPUT_FILE"

echo "Day $NEXT_NUM ready, instructions at https://adventofcode.com/$YEAR/day/$NEXT_NUM"

