package main

import (
	"os"
	"strings"

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

func Refresh(app *tview.Application, selfCareView *TabTextView, tasksView *TabTextView, weatherView *TabTextView, weeklyView *TabTextView, dailyView *TabTextView) {
	app.QueueUpdateDraw(func() {
		selfCareView.SetText(utils.GetSelfCareAdvice())
		tasks, weather, weekly_commits, daily_commits := GetTextForViews()
		tasksView.SetText(tasks)
		weatherView.SetText(weather)
		weeklyView.SetText(weekly_commits)
		dailyView.SetText(daily_commits)
	})
}

func GetTextForViews() (string, string, string, string) {

	var daily_commits string
	if strings.Contains(os.Getenv("TINYCARE_WORKSPACE"), ",") {
		for _, path := range strings.Split(os.Getenv("TINYCARE_WORKSPACE"), ",") {
			commits, err := utils.GetDailyCommits(path)
			if err != nil {
				panic(err)
			}

			daily_commits = daily_commits + commits
		}
	} else {
		commits, err := utils.GetDailyCommits(os.Getenv("TINYCARE_WORKSPACE"))

		if err != nil {
			panic(err)
		}

		daily_commits = commits
	}

	var weekly_commits string
	if strings.Contains(os.Getenv("TINYCARE_WORKSPACE"), ",") {
		for _, path := range strings.Split(os.Getenv("TINYCARE_WORKSPACE"), ",") {
			commits, err := utils.GetWeeklyCommits(path)
			if err != nil {
				panic(err)
			}

			weekly_commits = daily_commits + commits
		}
	} else {
		commits, err := utils.GetWeeklyCommits(os.Getenv("TINYCARE_WORKSPACE"))

		if err != nil {
			panic(err)
		}

		weekly_commits = commits
	}

	weather, err := utils.GetWeather(os.Getenv("TINYCARE_POSTAL_CODE"))
	if err != nil {
		panic(err)
	}

	tasks, err := utils.GetTodaysTasks()
	if err != nil {
		panic(err)
	}
	return tasks, weather, weekly_commits, daily_commits
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

	tasks, weather, weekly_commits, daily_commits := GetTextForViews()

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
			go Refresh(app, selfCareView, tasksView, weatherView, weeklyView, dailyView)
		}
		return event
	})

	app.SetRoot(flex, true).SetFocus(flex)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
