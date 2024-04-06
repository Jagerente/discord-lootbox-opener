package gui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"sort"
	"strings"
	"time"
)

type Log struct {
	CreatedAt time.Time
	Content   string
}

type Logger interface {
	Log(content string)
	GetLogs() []Log
}

type GUI struct {
	logger Logger

	app         *tview.Application
	statsWindow *tview.TextView
	logWindow   *tview.TextView
	tokenInput  *tview.InputField
	startButton *tview.Button
	stopButton  *tview.Button

	onStartHandler func(token string)
	onStopHandler  func()

	selectors            []tview.Primitive
	currentSelectorIndex int
}

func New(logger Logger) *GUI {
	g := &GUI{
		logger: logger,
	}

	g.app = tview.NewApplication()
	g.app.SetInputCapture(g.getHotKeysHandler())

	g.logWindow = tview.NewTextView()

	g.statsWindow = tview.NewTextView()

	g.tokenInput = tview.NewInputField().SetMaskCharacter('*')
	g.tokenInput.
		SetLabel("Token: ").
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				token := strings.TrimSpace(g.tokenInput.GetText())
				g.app.SetFocus(g.startButton)
				g.currentSelectorIndex = 1

				if g.onStartHandler != nil {
					g.onStartHandler(token)
				}
			}
		})

	g.startButton = tview.NewButton("Start")
	g.startButton.SetSelectedFunc(func() {
		if g.onStartHandler != nil {
			g.onStartHandler(g.tokenInput.GetText())
		}
	})

	g.stopButton = tview.NewButton("Stop")
	g.stopButton.SetSelectedFunc(func() {
		if g.onStopHandler != nil {
			g.onStopHandler()
		}
	})

	g.selectors = []tview.Primitive{
		g.tokenInput,
		g.startButton,
		g.stopButton,
	}

	return g
}

func (g *GUI) Draw() error {
	menuBox := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(g.buildHeader("Menu"), 1, 1, false).
		AddItem(nil, 1, 1, false).
		AddItem(g.tokenInput, 1, 1, false).
		AddItem(nil, 1, 1, false).
		AddItem(g.startButton, 3, 1, false).
		AddItem(nil, 1, 1, false).
		AddItem(g.stopButton, 3, 1, false)

	menuBox.SetBorder(true)

	statsBox := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(g.buildHeader("Stats"), 1, 1, false).
		AddItem(nil, 1, 1, false).
		AddItem(g.statsWindow, 0, 3, false)
	statsBox.SetBorder(true)

	leftMenu := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(menuBox, 13, 1, false).
		AddItem(statsBox, 0, 2, false)

	rightMenu := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(g.buildHeader("Logs"), 1, 1, false).
		AddItem(nil, 1, 1, false).
		AddItem(g.logWindow, 0, 10, false)
	rightMenu.SetBorder(true)

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(tview.NewFlex().
		AddItem(leftMenu, 0, 1, false).
		AddItem(rightMenu, 0, 1, false), 0, 1, false)

	if err := g.app.SetRoot(flex, true).SetFocus(g.tokenInput).EnableMouse(true).Run(); err != nil {
		return err
	}

	return nil
}

func (g *GUI) RegisterOnStartHandler(handler func(token string)) *GUI {
	g.onStartHandler = handler
	return g
}
func (g *GUI) RegisterOnStopHandler(handler func()) *GUI {
	g.onStopHandler = handler
	return g
}

func (g *GUI) UpdateStats(stats map[string]int) {
	g.redrawStatsWindow(getStatsText(stats))
}

func (g *GUI) Log(content string) {
	g.logger.Log(content)
	go g.redrawLogWindow()
}

func getStatsText(stats map[string]int) string {
	var builder strings.Builder
	keys := make([]string, 0, len(stats))
	for key := range stats {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		builder.WriteString(fmt.Sprintf("%s: %d\n", key, stats[key]))
	}
	return builder.String()
}

func (g *GUI) redrawLogWindow() {
	var builder strings.Builder
	logs := g.logger.GetLogs()
	for _, l := range logs {
		builder.WriteString(l.CreatedAt.Format("2006-01-02 15:04:05 ") + l.Content + "\n")
	}

	g.app.QueueUpdateDraw(func() {
		g.logWindow.SetText(builder.String())
	})
}

func (g *GUI) redrawStatsWindow(text string) {
	g.app.QueueUpdateDraw(func() {
		g.statsWindow.SetText(text)
	})
}

func (g *GUI) buildHeader(text string) tview.Primitive {
	return tview.NewTextView().SetText(text + ":").SetTextAlign(tview.AlignCenter)
}

func (g *GUI) nextSelection() {
	g.currentSelectorIndex++
	if g.currentSelectorIndex >= len(g.selectors) {
		g.currentSelectorIndex = 0
	}

	g.app.SetFocus(g.selectors[g.currentSelectorIndex])
}

func (g *GUI) previousSelection() {
	g.currentSelectorIndex--
	if g.currentSelectorIndex < 0 {
		g.currentSelectorIndex = len(g.selectors) - 1
	}

	g.app.SetFocus(g.selectors[g.currentSelectorIndex])
}

func (g *GUI) getHotKeysHandler() func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			g.app.SetFocus(nil)
			break
		case tcell.KeyDown:
			g.nextSelection()
			break
		case tcell.KeyUp:
			g.previousSelection()
			break
		}

		return event
	}
}
