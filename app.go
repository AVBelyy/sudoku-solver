package main

import (
    "ui"
)

func main() {
    ui.Term = make(chan bool)
    ui.Init()

    go ui.Event_loop()
    <-ui.Term
}
