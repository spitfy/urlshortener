package main

import (
	"log"
	"os"
)

func main() {
	log.Fatal("error in main")
	os.Exit(1)
}

func otherFunc() {
	os.Exit(1) // want "prohibited use of log.Fatal or os.Exit outside main package main function"
}
