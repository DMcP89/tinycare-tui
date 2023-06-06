package main

import (
    "github.com/rivo/tview"
)

func main() {
    box := tview.NewBox().SetBorder(true).SetTitle("Hello World")
    if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
        panic(err)
    }
}
