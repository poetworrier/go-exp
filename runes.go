package main

import (
	"bufio"
	_ "embed"
	"errors"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

type RuneSet []rune

var ErrNoRunes = errors.New("no runes given")
var ErrNoEmo = errors.New("no emojis left")

//go:embed emoji-sequences.txt
var emojiSeq string

func GenerateRuneSet() RuneSet {
	emojis := make([]rune, 0)
	seqSc := bufio.NewScanner(strings.NewReader(emojiSeq))
	for seqSc.Scan() {
		text := seqSc.Text()
		if len(text) == 0 {
			continue
		}
		if text[0] == '#' {
			continue
		}
		codePoints, _, ok := strings.Cut(seqSc.Text(), " ")
		if !ok {
			log.Fatal("invalid format")
		}
		cpRange := strings.Split(codePoints, "..")
		switch len(cpRange) {
		case 0:
			log.Fatal("invalid format")
		case 1:
			cp0, err := strconv.ParseInt(cpRange[0], 16, 32)
			if err != nil {
				log.Fatal(err)
			}
			if cp0 == 0 {
				continue
			}
			emojis = append(emojis, rune(cp0))
		case 2:
			cp0, err := strconv.ParseInt(cpRange[0], 16, 32)
			if err != nil {
				log.Fatal(err)
			}
			cp1, err := strconv.ParseInt(cpRange[1], 16, 32)
			if err != nil {
				log.Fatal(err)
			}
			if cp0 == 0 || cp1 == 0 {
				continue
			}
	
			for j := cp0; j <= cp1; j++ {
				emojis = append(emojis, rune(j))
			} 
		}	
	}

	rand.Shuffle(len(emojis), func(i, j int) {
		emojis[i], emojis[j] = emojis[j], emojis[i]
	})
	return emojis
}
