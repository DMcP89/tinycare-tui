package main

import (
	"github.com/DMcP89/tinycare-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TabTextView struct {
	*tview.TextView
	next *TabTextView
}

func (view *TabTextView) SetNext(next *TabTextView) *TabTextView {
	view.next = next
	return view
}

func (view *TabTextView) GetNext() *TabTextView {
	return view.next
}

func NewTabTextView(next *TabTextView) *TabTextView {
	return &TabTextView{
		TextView: tview.NewTextView(),
		next:     next,
	}
}

func refresh(app *tview.Application, selfCareView *TabTextView) {
	app.QueueUpdateDraw(func() {
		selfCareView.SetText(utils.GetSelfCareAdvice())
	})
}

func main() {

	app := tview.NewApplication()

	newTabTextView := func(text string, text_alignment int, next *TabTextView) *TabTextView {
		view := NewTabTextView(next)
		view.SetText(text).
			SetTextAlign(text_alignment).
			SetDynamicColors(true)
		return view

	}

	daily_commits, err := utils.GetDailyCommits("/home/dave/workspace/projects")

	if err != nil {
		panic(err)
	}

	weekly_commits, err := utils.GetWeeklyCommits("/home/dave/workspace/projects")
	if err != nil {
		panic(err)
	}

	weather, err := utils.GetWeather("07070")
	if err != nil {
		panic(err)
	}

	tasks, err := utils.GetTodaysTasks()
	if err != nil {
		panic(err)
	}

	tasksView := newTabTextView(tasks, tview.AlignLeft, nil)
	tasksView.SetWordWrap(false).SetWrap(false)
	tasksView.SetBorder(true).SetTitle("Today's Tasks 📋")

	selfCareView := newTabTextView(utils.GetSelfCareAdvice(), tview.AlignCenter, tasksView)
	selfCareView.SetBorder(true).SetTitle("Self Care 😁")

	weatherView := newTabTextView(weather, tview.AlignCenter, selfCareView)
	weatherView.SetBorder(true).SetTitle("Weather ⛅")

	weeklyView := newTabTextView(weekly_commits, tview.AlignLeft, weatherView)
	weeklyView.SetBorder(true).SetTitle("Weekly Commits 📦")

	dailyView := newTabTextView(daily_commits, tview.AlignLeft, weeklyView)
	dailyView.SetBorder(true).SetTitle("Daily Commits 📦")

	tasksView.SetNext(dailyView)

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(dailyView, 0, 1, true).
			AddItem(weeklyView, 0, 2, false), 0, 2, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(weatherView, 0, 2, false).
			AddItem(selfCareView, 0, 1, false).
			AddItem(tasksView, 0, 4, false), 0, 1, false)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			app.SetFocus(app.GetFocus().(*TabTextView).next)
		} else if (event.Key() == tcell.KeyRune) && (event.Rune() == rune('q')) {
			app.Stop()
		} else if (event.Key() == tcell.KeyRune) && (event.Rune() == rune('r')) {
			go refresh(app, selfCareView)
		}
		return event
	})

	app.SetRoot(flex, true).SetFocus(flex)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
