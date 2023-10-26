package main

import (
	"log"
)

func main() {
	app := NewApp()
	err := app.StartApp()
	if err != nil {
		log.Fatalln(err)
	}
}
