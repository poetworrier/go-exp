package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"unicode/utf8"

	"go.uber.org/zap"
)

var (
	decryptF = flag.Bool("decrypt", false, "decrypt input from crypt fileset")
	keypath  = flag.String("keypath", "crypt.key", "crypt key file path")
	pubpath  = flag.String("pubpath", "crypt.pub", "crypt pub file path")
	debug    = flag.Bool("debug", false, "enabled debug logging")
)

const (
	PemTypePub  = "OTP PUBLIC"
	PemTypePriv = "OTP KEY"
)

var logger *zap.Logger

var newBuf, newScan = bytes.NewBuffer, bufio.NewScanner

func main() {
	flag.Parse()
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	if *debug {
		logger, err = zap.NewDevelopment()
		if err != nil {
			log.Fatal(err)
		}
	}

	if *decryptF {
		keyfile, err := os.Open(*keypath)
		if err != nil {
			log.Fatal(err)
		}

		pubfile, err := os.Open(*pubpath)
		if err != nil {
			log.Fatal(err)
		}

		if err := decrypter(os.Stdout, pubfile, keyfile); err != nil {
			log.Fatal(err)
		}
		return
	}

	keyfile, err := os.Create(*keypath)
	if err != nil {
		log.Fatal(err)
	}
	pubfile, err := os.Create(*pubpath)
	if err != nil {
		log.Fatal(err)
	}
	err = encrypter(pubfile, keyfile, os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
}

func debugPemBlocks(pub, priv *pem.Block) {
	if logger.Level() == zap.DebugLevel {
		logger.Sugar().Debugf("public block: %+v\nprivate block: %+v\n", pub, priv)
	}
}

func encrypter(pubfile, privfile io.Writer, r io.Reader) error {
	pub, priv := &pem.Block{Type: PemTypePub}, &pem.Block{Type: PemTypePriv}
	if err := encryptReader(pub, priv, r); err != nil {
		return err
	}
	debugPemBlocks(pub, priv)

	if err := pem.Encode(pubfile, pub); err != nil {
		return err
	}
	if err := pem.Encode(privfile, priv); err != nil {
		return err
	}
	return nil
}

func scanCrypt(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	return 4, data[0:4], nil
}

func decrypter(out io.Writer, pub, priv io.Reader) error {
	pubBlck, err := readPem(pub)
	if err != nil {
		return err
	}
	privBlck, err := readPem(priv)
	if err != nil {
		return err
	}
	debugPemBlocks(pubBlck, privBlck)

	scPub := newScan(newBuf(pubBlck.Bytes))
	scPriv := newScan(newBuf(privBlck.Bytes))
	scPub.Split(scanCrypt)
	scPriv.Split(scanCrypt)

	for scPub.Scan() && scPriv.Scan() {
		if err := errors.Join(scPub.Err(), scPriv.Err()); err != nil {
			return err
		}
		fmt.Fprintf(out, "%c", decrypt(scPub.Bytes(), scPriv.Bytes()))
	}

	return nil
}

func encryptReader(pub, priv *pem.Block, reader io.Reader) error {
	sc := newScan(reader)
	sc.Split(bufio.ScanRunes)
	for sc.Scan() {
		if sc.Err() != nil {
			return sc.Err()
		}
		ru, size := utf8.DecodeRune(sc.Bytes())
		if size == 0 {
			return errors.New("empty rune in reader")
		}
		pubBytes, privByte, err := encrypt(ru)
		if err != nil {
			return err
		}
		pub.Bytes = append(pub.Bytes, pubBytes...)
		priv.Bytes = append(priv.Bytes, privByte...)
	}

	return nil
}

func encrypt(r rune) (encBytes, keyBytes []byte, err error) {
	keyBytes = make([]byte, 4)
	_, err = rand.Read(keyBytes)
	if err != nil {
		return
	}
	key := binary.BigEndian.Uint32(keyBytes)
	enc := uint32(r) ^ key
	encBytes = make([]byte, 4)
	binary.BigEndian.PutUint32(encBytes, enc)
	return
}

func decrypt(encR []byte, keyBytes []byte) rune {
	ruInt := binary.BigEndian.Uint32(encR)
	key := binary.BigEndian.Uint32(keyBytes)
	return rune(key ^ ruInt)
}

func readPem(tmp io.Reader) (*pem.Block, error) {
	buf, err := io.ReadAll(tmp)
	if err != nil {
		return nil, err
	}
	blck, _ := pem.Decode(buf)
	return blck, nil
}
