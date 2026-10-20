package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("usage: go run . <release-name> <phone>")
	}
	c, err := NewInfraiClient()
	if err != nil {
		log.Fatal(err)
	}
	if err := PublishRelease(c, Release{Name: os.Args[1], Phone: os.Args[2]}); err != nil {
		log.Fatal(err)
	}
	fmt.Println("release SMS published")
}
