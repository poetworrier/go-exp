package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"go.uber.org/zap"
)

func setup(t *testing.T) {
	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		t.Error(err)
	}

}

func TestEncryptDecrypt(t *testing.T) {
	setup(t)
	// TODO: Fuzz on this
	want := 'H'
	r, k, err := encrypt(want)
	if err != nil {
		t.Error(err)
	}
	got := decrypt(r, k)
	if got != want {
		t.Errorf("got=%c, want=%c", got, want)
	}
}

func TestEncrypter_Decrypter(t *testing.T) {
	setup(t)
	want := "Hello World"
	var pubBuf, privBuf bytes.Buffer
	in := strings.NewReader(want)
	if err := encrypter(&pubBuf, &privBuf, in); err != nil {
		t.Error(err)
		return
	}
	var got strings.Builder
	if err := decrypter(&got, &pubBuf, &privBuf); err != nil {
		t.Error(err)
		return
	}

	if want != got.String() {
		t.Errorf("want=%q, got=%q", want, got.String())
	}
	fmt.Println(got)
}
