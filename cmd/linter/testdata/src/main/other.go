package main

import "log"

func init() {
	log.Fatal("error in init") // want "prohibited use of log.Fatal or os.Exit outside main package main function"
	panic("panic in init")     // want "prohibited use of panic()"
}

func anotherFunc() {
	panic("should be reported") // want "prohibited use of panic()"
}
