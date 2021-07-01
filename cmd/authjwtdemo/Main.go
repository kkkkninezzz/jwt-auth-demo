package main

import (
	"authjwtdemo/internal/app/authjwtdemo"
	"flag"
	"log"
)

var ConfigPath = flag.String("config", "../../configs/app.yaml", "配置文件路径")

func main() {
	flag.Parse()

	log.SetFlags(log.LstdFlags)
	authjwtdemo.Boot(*ConfigPath)
}
