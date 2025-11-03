package main

import (
	"log"
	"os"
)

func main() {
	// Эти вызовы должны быть разрешены
	log.Fatal("error in main")
	os.Exit(1)
}

func otherFunc() {
	// Этот вызов должен быть запрещен
	os.Exit(1) // want "prohibited use of log.Fatal or os.Exit outside main package main function"
}
