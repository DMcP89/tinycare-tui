package main

import (
	"fmt"
	"os"
	"time"

	"github.com/DMcP89/tinycare-tui/internal/apis"
	"github.com/DMcP89/tinycare-tui/internal/local"
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

func main() {
	app := tview.NewApplication()
	changeFunc := func() { app.Draw() }
	newTabTextView := func(text string, textAlignment int, next *TabTextView) *TabTextView {
		view := NewTabTextView(next)
		view.SetText(text).
			SetTextAlign(textAlignment).
			SetDynamicColors(true)
		return view
	}

	tasksView := newTabTextView("", tview.AlignLeft, nil)
	tasksView.SetWordWrap(true).SetWrap(true)
	tasksView.SetBorder(true).SetTitle("Today's Tasks 📋")
	tasksView.SetChangedFunc(changeFunc)

	selfCareView := newTabTextView("", tview.AlignCenter, tasksView)
	selfCareView.SetBorder(true).SetTitle("Self Care 😁")
	selfCareView.SetChangedFunc(changeFunc)

	jokeView := newTabTextView("", tview.AlignCenter, selfCareView)
	jokeView.SetBorder(true).SetTitle("Joke 🤣")
	jokeView.SetChangedFunc(changeFunc)

	weatherView := newTabTextView("", tview.AlignCenter, jokeView)
	weatherView.SetBorder(true).SetTitle("Weather ⛅")
	weatherView.SetChangedFunc(changeFunc)

	weeklyView := newTabTextView("", tview.AlignLeft, weatherView)
	weeklyView.SetBorder(true).SetTitle("Weekly Commits 📦")
	weeklyView.SetChangedFunc(changeFunc)

	dailyView := newTabTextView("", tview.AlignLeft, weeklyView)
	dailyView.SetBorder(true).SetTitle("Daily Commits 📦")
	dailyView.SetChangedFunc(changeFunc)

	tasksView.SetNext(dailyView)

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(dailyView, 0, 1, true).
			AddItem(weeklyView, 0, 2, false), 0, 2, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(weatherView, 0, 2, false).
			AddItem(jokeView, 0, 1, false).
			AddItem(selfCareView, 0, 1, false).
			AddItem(tasksView, 0, 4, false), 0, 1, false)
	refresh := func() {
		go func() {
			advice := local.GetSelfCareAdvice()
			selfCareView.SetText(advice)
		}()

		go func() {
			joke, err := apis.GetJoke()
			if err != nil {
				jokeView.SetText(err.Error())
			}
			jokeView.SetText(joke)
		}()

		go func() {
			var result string
			var err error
			if token, ok := os.LookupEnv("TODOIST_TOKEN"); ok {
				result, err = apis.GetTodaysTasks(token)
			} else {
				if todoFile, ok := os.LookupEnv("TODO_FILE"); ok {
					result, err = local.GetLocalTasks(todoFile)
				} else {
					tasksView.SetText("Please set either the TODOIST_TOKEN or TODO_FILE environment variable")
				}
			}
			if err != nil {
				tasksView.SetText(err.Error())
			}
			tasksView.SetText(result)
		}()

		go func() {
			if POSTAL_CODE, ok := os.LookupEnv("TINYCARE_POSTAL_CODE"); ok {
				result, err := apis.GetWeather(POSTAL_CODE)
				if err != nil {
					weatherView.SetText(err.Error())
				}
				weatherView.SetText(result)
			} else {
				weatherView.SetText("Please set TINYCARE_POSTAL_CODE environment variable")
			}
		}()

		go func() {
			if token, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
				result, err := apis.GetGitHubCommits(token, -7)
				if err != nil {
					weeklyView.SetText(err.Error())
				}
				weeklyView.SetText(result)
			} else {
				if TINYCARE_WORKSPACE, ok := os.LookupEnv("TINYCARE_WORKSPACE"); ok {
					result, err := local.GetCommits(TINYCARE_WORKSPACE, -7)
					if err != nil {
						weeklyView.SetText(err.Error())
					}
					weeklyView.SetText(result)
				} else {
					weeklyView.SetText("Please set either the TINYCARE_WORKSPACE or GITHUB_TOKEN environment variables to retrive commits")
				}
			}
		}()
		go func() {
			if token, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
				result, err := apis.GetGitHubCommits(token, -1)
				if err != nil {
					dailyView.SetText(err.Error())
				}
				dailyView.SetText(result)
			} else {
				if TINYCARE_WORKSPACE, ok := os.LookupEnv("TINYCARE_WORKSPACE"); ok {
					result, err := local.GetCommits(TINYCARE_WORKSPACE, -1)
					if err != nil {
						dailyView.SetText(err.Error())
					}
					dailyView.SetText(result)
				} else {
					dailyView.SetText("Please set either the TINYCARE_WORKSPACE or GITHUB_TOKEN environment variables to retrive commits")
				}
			}
		}()
	}

	go func() {
		for {
			refresh()
			time.Sleep(300 * time.Second)
		}
	}()

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			app.SetFocus(app.GetFocus().(*TabTextView).next)
		} else if (event.Key() == tcell.KeyRune) && (event.Rune() == rune('q')) {
			app.Stop()
		} else if (event.Key() == tcell.KeyRune) && (event.Rune() == rune('r')) {
			go refresh()
		}
		return event
	})

	app.SetRoot(flex, true).SetFocus(flex)

	if err := app.Run(); err != nil {
		panic(err)
	}

}

func handleRefreshError(err error) {
	fmt.Printf("Error: %v\n", err)
}
