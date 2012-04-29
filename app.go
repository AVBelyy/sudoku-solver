package main

import (
    "ui"
    "flag"
)

func main() {
    flag.Bool("9", true, "use 9x9 field")
    is6 := flag.Bool("6", false, "use 6x6 field")
    flag.Parse()

    size := uint(9)
    if *is6 { size = 6 }

    ui.Term = make(chan bool)
    ui.Init(size)

    go ui.Event_loop()
    <-ui.Term
}
