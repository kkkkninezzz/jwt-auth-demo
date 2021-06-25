package main

import (
    "gossodemo/internal/app/gossodemo"
    "log"
)

func main() {
    log.SetFlags(log.LstdFlags)
    gossodemo.Boot()
}
