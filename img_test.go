package main

import (
	"testing"
)

func TestGImg(t *testing.T) {
	config := Config{CoverHight: 1200, CoverWidth: 860}
	_, _, err := config.generateCover("examples/cover_example.jpg")
	if err != nil {
		t.Fatal(err)
	}
}
