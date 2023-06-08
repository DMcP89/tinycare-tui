package main

import (
	"github.com/DMcP89/tinycare-tui/internal/utils"
	"github.com/rivo/tview"
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

	weather, err := utils.GetWeather("07070")
	if err != nil {
		panic(err)
	}
	weatherView := tview.NewTextView().
		SetText(weather).
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	weatherView.SetBorder(true).SetTitle("Weather")

	tinyCareBotTweet, err := utils.GetLatestTweet("tinycarebot")

	if err != nil {
		panic(err)
	}

	tinycarebotView := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetText(tinyCareBotTweet)

	tinycarebotView.SetBorder(true).SetTitle("Tinycarebot")

	selfCareBotTweet, err := utils.GetLatestTweet("selfcare_bot")
	if err != nil {
		panic(err)
	}

	selfcarebotView := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetText(selfCareBotTweet)

	selfcarebotView.SetBorder(true).SetTitle("Selfcarebot")

	tasks, err := utils.GetTodaysTasks()
	if err != nil {
		panic(err)
	}

	tasksView := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true).
		SetText(tasks)

	tasksView.SetBorder(true).SetTitle("Today's Tasks")

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(dailyView, 0, 1, false).
			AddItem(weeklyView, 0, 1, false), 0, 2, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(weatherView, 0, 1, false).
			AddItem(tinycarebotView, 0, 1, false).
			AddItem(selfcarebotView, 0, 1, false).
			AddItem(tasksView, 0, 6, false), 0, 1, false)

	app.SetRoot(flex, true).SetFocus(flex)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
