package main

import (
    "github.com/rivo/tview"
    "github.com/DMcP89/tinycare-tui/internal/utils"
)

func main() {
   app := tview.NewApplication()

   daily_commits, err := utils.GetDailyCommits("/home/dave/workspace/projects")

    if err != nil {
        panic(err)
    }

   textView := tview.NewTextView().
        SetText(daily_commits).
        SetTextAlign(tview.AlignLeft).
        SetDynamicColors(true)

   textView.SetBorder(true).SetTitle("Daily Commits")

   app.SetRoot(textView, true)

   if err := app.Run(); err != nil{
        panic(err)
   }
}
