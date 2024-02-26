package main

import (
	"fmt"
	"os"
	"time"

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

func RefreshText(app *tview.Application, view *TabTextView, textFunc func() (string, error)) {
	app.QueueUpdateDraw(func() {
		text, err := textFunc()
		if err != nil {
			handleRefreshError(err)
			return
		}
		view.SetText(text)
	})
}

func main() {
	app := tview.NewApplication()

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

	selfCareView := newTabTextView("", tview.AlignCenter, tasksView)
	selfCareView.SetBorder(true).SetTitle("Self Care 😁")

	weatherView := newTabTextView("", tview.AlignCenter, selfCareView)
	weatherView.SetBorder(true).SetTitle("Weather ⛅")

	weeklyView := newTabTextView("", tview.AlignLeft, weatherView)
	weeklyView.SetBorder(true).SetTitle("Weekly Commits 📦")

	dailyView := newTabTextView("", tview.AlignLeft, weeklyView)
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

	refresh := func() {
		go RefreshText(app, selfCareView, func() (string, error) {
			return utils.GetSelfCareAdvice(), nil
		})
		go RefreshText(app, tasksView, func() (string, error) {
			result, err := utils.GetTasks()
			if err != nil {
				return err.Error(), nil
			}
			return result, err
		})
		go RefreshText(app, weatherView, func() (string, error) {
			if POSTAL_CODE, ok := os.LookupEnv("TINYCARE_POSTAL_CODE"); ok {
				result, err := utils.GetWeather(POSTAL_CODE)
				if err != nil {
					return err.Error(), nil
				}
				return result, err
			} else {
				return "Please set TINYCARE_POSTAL_CODE environment variable", nil
			}
		})
		go RefreshText(app, weeklyView, func() (string, error) {
			if token, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
				result, err := utils.GetGitHubCommits(token, -7)
				if err != nil {
					return err.Error(), nil
				}
				return result, err
			} else {
				if TINYCARE_WORKSPACE, ok := os.LookupEnv("TINYCARE_WORKSPACE"); ok {
					result, err := utils.GetWeeklyCommits(TINYCARE_WORKSPACE)
					if err != nil {
						return err.Error(), nil
					}
					return result, err
				} else {
					return "Please set either the TINYCARE_WORKSPACE or GITHUB_TOKEN environment variables to retrive commits", nil
				}
			}
		})
		go RefreshText(app, dailyView, func() (string, error) {
			if token, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
				result, err := utils.GetGitHubCommits(token, -1)
				if err != nil {
					return err.Error(), nil
				}
				return result, err
			} else {
				if TINYCARE_WORKSPACE, ok := os.LookupEnv("TINYCARE_WORKSPACE"); ok {
					result, err := utils.GetDailyCommits(TINYCARE_WORKSPACE)
					if err != nil {
						return err.Error(), nil
					}
					return result, err
				} else {
					return "Please set either the TINYCARE_WORKSPACE or GITHUB_TOKEN environment variables to retrive commits", nil
				}
			}
		})
	}

	go func() {
		for {
			refresh()
			time.Sleep(30 * time.Second)
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
