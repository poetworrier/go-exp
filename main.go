package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
)

func main() {
	rs := GenerateRuneSet()
	emojis := make([]rune, 256)
	for i := 0; i < 256; i++ {
		emojis[i], rs = rs[0], rs[1:] 
	}

	fmt.Println("running with emoji map")
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			fmt.Printf("%c", emojis[(15*i)+j])
		}
		fmt.Print("\n")
	}
	fmt.Println()

	sc := bufio.NewScanner(os.Stdin)
	sc.Split(bufio.ScanRunes)
	for sc.Scan() {
		key, encrypted, err := encrypt(emojis, sc.Bytes()[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s %c\n", base64.StdEncoding.EncodeToString([]byte{key}), encrypted)
	}
	if sc.Err() != nil {
		log.Fatal(sc.Err())
	}
}

func encrypt(runes RuneSet, r byte) (byte, rune, error) {
	randomByte := make([]byte, 1)
	_, err := rand.Read(randomByte)
	if err != nil {
		return 0, 0, err
	}
	key := randomByte[0]
	encryptedRune := runes[(r^key)%255]
	return key, encryptedRune, nil
}
