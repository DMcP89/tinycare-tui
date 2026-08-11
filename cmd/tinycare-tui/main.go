package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"time"

	"github.com/DMcP89/tinycare-tui/internal/apis"
	"github.com/DMcP89/tinycare-tui/internal/local"
	"github.com/DMcP89/tinycare-tui/internal/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func GetTextForView(f func(string) (string, error), envVar string, missingEnvErrorMessage string) string {
	if token, ok := os.LookupEnv(envVar); ok {
		result, err := f(token)
		if err != nil {
			return err.Error()
		}
		return result
	} else {
		return missingEnvErrorMessage
	}
}

func GetRefreshInterval() time.Duration {
	const defaultInterval = 300 // 300 seconds = 5 minutes
	if intervalStr, ok := os.LookupEnv("TINYCARE_REFRESH_INTERVAL"); ok {
		if interval, err := strconv.Atoi(intervalStr); err == nil && interval > 0 {
			return time.Duration(interval) * time.Second
		}
	}
	return defaultInterval * time.Second
}

func runProfiledRefresh(runs int) {
	if ws, ok := os.LookupEnv("TINYCARE_WORKSPACE"); ok {
		for i := 0; i < runs; i++ {
			_, _, err := local.GetCommits(ws)
			if err != nil {
				fmt.Println("GetCommits error:", err)
			}
		}
	}
	for i := 0; i < runs; i++ {
		local.GetSelfCareAdvice()
	}
}

func main() {
	profileFile := flag.String("profile", "", "write CPU profile to file and exit")
	memFile := flag.String("memprofile", "", "write heap profile to file and exit")
	runs := flag.Int("profile-runs", 5, "number of data-fetch iterations while profiling")
	flag.Parse()

	if *profileFile != "" || *memFile != "" {
		if *profileFile != "" {
			f, err := os.Create(*profileFile)
			if err != nil {
				panic(err)
			}
			if err := pprof.StartCPUProfile(f); err != nil {
				panic(err)
			}
			runProfiledRefresh(*runs)
			pprof.StopCPUProfile()
			f.Close()
			fmt.Println("CPU profile written to", *profileFile)
		}
		if *memFile != "" {
			f, err := os.Create(*memFile)
			if err != nil {
				panic(err)
			}
			runProfiledRefresh(*runs)
			runtime.GC()
			if err := pprof.WriteHeapProfile(f); err != nil {
				panic(err)
			}
			f.Close()
			fmt.Println("heap profile written to", *memFile)
		}
		return
	}

	app := tview.NewApplication()
	changeFunc := func() { app.Draw() }
	newTabTextView := func(text string, textAlignment int, next *ui.TabTextView) *ui.TabTextView {
		view := ui.NewTabTextView(next)
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
			text := GetTextForView(apis.GetTodaysTasks, "TODOIST_TOKEN", "")
			if text == "" {
				text = GetTextForView(func(_ string) (string, error) { return local.GetTaskwarriorTasks() }, "TASKWARRIOR", "")
			}
			if text == "" {
				text = GetTextForView(local.GetLocalTasks, "TODO_FILE", "Please set either the TODOIST_TOKEN, TASKWARRIOR, or TODO_FILE environment variable")
			}
			tasksView.SetText(text)
		}()

		go func() {
			text := GetTextForView(apis.GetWeather, "TINYCARE_POSTAL_CODE", "Please set TINYCARE_POSTAL_CODE environment variable")
			weatherView.SetText(text)
		}()

		go func() {
			if token, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
				dayResult, weekResult, err := apis.GetGitHubCommits(token)
				if err != nil {
					dailyView.SetText(err.Error())
					weeklyView.SetText(err.Error())
				}
				dailyView.SetText(dayResult)
				weeklyView.SetText(weekResult)
			} else {
				if TINYCARE_WORKSPACE, ok := os.LookupEnv("TINYCARE_WORKSPACE"); ok {
					dayResult, weekResult, err := local.GetCommits(TINYCARE_WORKSPACE)
					if err != nil {
						dailyView.SetText(err.Error())
						weeklyView.SetText(err.Error())
					} else {
						dailyView.SetText(dayResult)
						weeklyView.SetText(weekResult)
					}
				} else {
					weeklyView.SetText("Please set either the TINYCARE_WORKSPACE or GITHUB_TOKEN environment variables to retrive commits")
				}
			}
		}()
	}

	refreshInterval := GetRefreshInterval()
	go func() {
		for {
			refresh()
			time.Sleep(refreshInterval)
		}
	}()

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			app.SetFocus(app.GetFocus().(*ui.TabTextView).GetNext())
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
