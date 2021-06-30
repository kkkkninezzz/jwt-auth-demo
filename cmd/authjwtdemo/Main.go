package main

import (
	"authjwtdemo/internal/app/authjwtdemo"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags)
	authjwtdemo.Boot()
}
