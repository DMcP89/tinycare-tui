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

   dailyView := tview.NewTextView().
        SetText(daily_commits).
        SetTextAlign(tview.AlignLeft).
        SetDynamicColors(true)

   dailyView.SetBorder(true).SetTitle("Daily Commits")

   weekly_commits, err := utils.GetWeeklyCommits("/home/dave/workspace/projects")
   if err != nil {
        panic(err)
   }

   weeklyView := tview.NewTextView().
        SetText(weekly_commits).
        SetTextAlign(tview.AlignLeft).
        SetDynamicColors(true)

   weeklyView.SetBorder(true).SetTitle("Weekly Commits")

   flex := tview.NewFlex().
            AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
            AddItem(dailyView,0,1,false).
            AddItem(weeklyView,0,1,false),0,1,false)

   app.SetRoot(flex, true).SetFocus(flex)

   if err := app.Run(); err != nil{
        panic(err)
   }
}
