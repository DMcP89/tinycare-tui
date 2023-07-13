package main

import (
	"github.com/DMcP89/tinycare-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func cycleFocus(app *tview.Application, elements []tview.Primitive, reverse bool) {
	for i, el := range elements {
		if !el.HasFocus() {
			continue
		}

		if reverse {
			i = i - 1
			if i < 0 {
				i = len(elements) - 1
			}
		} else {
			i = i + 1
			i = i % len(elements)
		}

		app.SetFocus(elements[i])
		return
	}
}

func main() {

	app := tview.NewApplication()

	newTextView := func(text string, text_alignment int) *tview.TextView {
		return tview.NewTextView().
			SetText(text).
			SetTextAlign(text_alignment).
			SetDynamicColors(true)
	}

	daily_commits, err := utils.GetDailyCommits("/home/dave/workspace/projects")

	if err != nil {
		panic(err)
	}

	dailyView := newTextView(daily_commits, tview.AlignLeft)
	dailyView.SetBorder(true).SetTitle("Daily Commits 📦")

	weekly_commits, err := utils.GetWeeklyCommits("/home/dave/workspace/projects")
	if err != nil {
		panic(err)
	}

	weeklyView := newTextView(weekly_commits, tview.AlignLeft)
	weeklyView.SetBorder(true).SetTitle("Weekly Commits 📦")

	weather, err := utils.GetWeather("07070")
	if err != nil {
		panic(err)
	}
	weatherView := newTextView(weather, tview.AlignCenter)
	weatherView.SetBorder(true).SetTitle("Weather ⛅")

	tasks, err := utils.GetTodaysTasks()
	if err != nil {
		panic(err)
	}

	tasksView := newTextView(tasks, tview.AlignLeft).SetWordWrap(false).SetWrap(false)
	tasksView.SetBorder(true).SetTitle("Today's Tasks 📋")

	selfCareView := newTextView(utils.GetSelfCareAdvice(), tview.AlignCenter)
	selfCareView.SetBorder(true).SetTitle("Self Care 😁")

	textViews := []tview.Primitive{
		dailyView,
		weeklyView,
		weatherView,
		selfCareView,
		tasksView,
	}

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(dailyView, 0, 1, true).
			AddItem(weeklyView, 0, 2, false), 0, 2, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(weatherView, 0, 1, false).
			AddItem(selfCareView, 0, 1, false).
			AddItem(tasksView, 0, 4, false), 0, 1, false)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			cycleFocus(app, textViews, false)
		} else if event.Key() == tcell.KeyBacktab {
			cycleFocus(app, textViews, true)
		} else if (event.Key() == tcell.KeyRune) && (event.Rune() == rune('q')) {
			app.Stop()
		}
		return event
	})

	app.SetRoot(flex, true).SetFocus(flex)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
