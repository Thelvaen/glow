package ui

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"glow/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Run starts the Fyne application and builds the UI from the provided config
func Run(cfg config.Config) {
	a := app.New()
	w := a.NewWindow(cfg.Title)
	w.Resize(fyne.NewSize(1000, 600))

	// Tabs area
	tabs := container.NewAppTabs()
	tabs.SetTabLocation(container.TabLocationTop)

	// Keep track of tab sessions and tab containers
	var mu sync.Mutex
	tabSessions := map[string]*Session{}
	tabsMap := map[string]*fyne.Container{}

	getOrCreateTab := func(tabName string) *fyne.Container {
		mu.Lock()
		defer mu.Unlock()
		if cont, ok := tabsMap[tabName]; ok {
			return cont
		}
		// Create new tab with placeholder label
		termCont := container.NewBorder(nil, nil, nil, nil, widget.NewLabel("No session"))
		tabs.Append(container.NewTabItem(tabName, termCont))
		tabsMap[tabName] = termCont
		return termCont
	}

	startSessionOnTab := func(tabName, cmd string) error {
		cont := getOrCreateTab(tabName)
		// create session
		sess, err := NewSession(tabName, cmd, "", nil)
		if err != nil {
			return err
		}

		mu.Lock()
		if prev, ok := tabSessions[tabName]; ok {
			prev.Close()
		}
		tabSessions[tabName] = sess
		mu.Unlock()

		// Build container with control buttons and terminal
		stopBtn := widget.NewButton("Stop", func() {
			mu.Lock()
			if cur, ok := tabSessions[tabName]; ok {
				cur.Close()
				delete(tabSessions, tabName)
			}
			mu.Unlock()
		})
		reloadBtn := widget.NewButton("Restart", func() {
			_ = startSessionOnTab(tabName, cmd)
		})
		top := container.NewHBox(widget.NewLabel(fmt.Sprintf("Tab: %s", tabName)), stopBtn, reloadBtn)
		termCont := container.NewBorder(top, nil, nil, nil, sess.Term)

		// Replace content
		cont.Objects = []fyne.CanvasObject{termCont}
		cont.Refresh()
		return nil
	}

	runScenario := func(s config.Scenario) {
		for _, a := range s.Actions {
			if a.Wait > 0 {
				// wait
				time.Sleep(time.Duration(a.Wait) * time.Second)
				continue
			}
			if a.Tab == "" {
				// default tab name = scenario
				a.Tab = s.Name
			}
			if len(a.Cmds) == 0 {
				continue
			}
			cmdStr := strings.Join(a.Cmds, " ; ")
			if err := startSessionOnTab(a.Tab, cmdStr); err != nil {
				dialog.ShowError(err, w)
				return
			}
			// brief small delay so sequential actions settle
			time.Sleep(250 * time.Millisecond)
		}
	}

	// Left side: scenario buttons
	buttons := container.NewVBox()
	for _, s := range cfg.Scenarios {
		s := s
		btn := widget.NewButton(s.Name, func() {
			go runScenario(s)
		})
		buttons.Add(btn)
	}

	head := container.NewVBox(widget.NewLabel(cfg.Title), widget.NewSeparator(), buttons)

	split := container.NewHSplit(head, tabs)
	split.Offset = 0.25

	w.SetContent(split)
	w.ShowAndRun()
}
