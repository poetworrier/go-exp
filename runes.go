package main

import (
	"bufio"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

var ErrNoRunes = errors.New("no runes given")
var ErrNoEmo = errors.New("no emojis left")

//go:embed emoji-sequences.txt
var emojiSeq string

type RangeSpec int

const (
	Empty = RangeSpec(iota)
	Single
	Range
)

func NewEmojis(size int) []rune {
	rs, err := GenerateEmojiRunes()
	if err != nil {
		log.Fatal(err)
	}
	emojis := make([]rune, size)
	for i := 0; i < size; i++ {
		emojis[i], rs = rs[0], rs[1:]
	}
	return emojis
}


func GenerateEmojiRunes() ([]rune, error) {
	emojis := make([]rune, 0)
	seqSc := bufio.NewScanner(strings.NewReader(emojiSeq))
	for seqSc.Scan() {
		spec := seqSc.Text()
		if len(spec) == 0 || spec[0] == '#' {
			continue
		}
		codePoints, _, ok := strings.Cut(seqSc.Text(), " ")
		if !ok {
			return nil, fmt.Errorf("invalid format: %q", spec)
		}

		cpRange := strings.Split(codePoints, "..")
		switch RangeSpec(len(cpRange)) {
		case Empty:
			return nil, fmt.Errorf("invalid format: %q", spec)
		case Single:
			var cp0 int64
			if err := parseCodePoint(cpRange[0], &cp0); err != nil {
				return nil, err
			}
			emojis = append(emojis, rune(cp0))
		case Range:
			var cp0, cp1 int64
			if err := parseCodePoint(cpRange[0], &cp0); err != nil {
				return nil, err
			}
			if err := parseCodePoint(cpRange[1], &cp1); err != nil {
				return nil, err
			}
			for j := cp0; j <= cp1; j++ {
				emojis = append(emojis, rune(j))
			}
		}
	}

	rand.Shuffle(len(emojis), func(i, j int) {
		emojis[i], emojis[j] = emojis[j], emojis[i]
	})
	return emojis, nil
}

func parseCodePoint(spec string, out *int64) error {
	var err error
	if *out, err = strconv.ParseInt(spec, 16, 32); err != nil {
		return err
	}
	if *out == 0 {
		return fmt.Errorf("zero value codepoint: %q", spec)
	}
	return nil
}
