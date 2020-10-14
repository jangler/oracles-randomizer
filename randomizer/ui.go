package randomizer

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell"
)

// a uiSegment is a text string with formatting information.
type uiSegment struct {
	style tcell.Style
	text  string
	el    uiEllipsis
}

// a uiLine is just an array of segments.
type uiLine []uiSegment

// a uiEllipsis defines whether a segment will opt in to being truncated at the
// left or right to fit onscreen. by default, lines are truncated at the right,
// but if (for example) a segment in the line has ellipsisLeft, that segment
// will be truncated at the left in order to make the line fit.
type uiEllipsis int

// ellipsis constants
const (
	ellipsisNone uiEllipsis = iota
	ellipsisLeft
	ellipsisRight
)

// a uiMode defines the way the UI currently handles input and displays
// information.
type uiMode int

// uiMode constants
const (
	modeWorking uiMode = iota
	modePrompt
	modeDone
)

var uiBottom = []uiSegment{
	{text: "(q)", style: tcell.StyleDefault.Bold(true)},
	{text: "uit"},
}

type uiInstance struct {
	// this one's actually used as a constant, but can't be declared as one
	lines          []uiLine
	screen         tcell.Screen
	write, rewrite chan uiLine
	input, prompt  chan rune
	resize         chan interface{}
	change         chan uiMode
}

// creates and displays a blank TUI.
func newUI(title string) *uiInstance {
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err := screen.Init(); err != nil {
		_ = err
	}

	ui := &uiInstance{
		lines: []uiLine{
			[]uiSegment{{text: title}},
		},
		screen:  screen,
		write:   make(chan uiLine, 1),      // add a line
		rewrite: make(chan uiLine),         // rewrite the last line
		input:   make(chan rune),           // key input to main loop
		prompt:  make(chan rune),           // key input passed from main to prompt
		resize:  make(chan interface{}, 1), // send to update window size
		change:  make(chan uiMode, 1),      // change uiMode
	}

	ui.draw(modeWorking)

	return ui
}

// runs the TUI, waiting for input from other functions and displaying updated
// information as needed.
func (ui *uiInstance) run() {
	// run event processing in a different goroutine
	go func() {
		for {
			switch evt := ui.screen.PollEvent().(type) {
			case *tcell.EventKey:
				switch evt.Key() {
				case tcell.KeyCtrlC, tcell.KeyDEL: // del is backspace
					ui.input <- rune(evt.Key())
				case tcell.KeyRune:
					ui.input <- evt.Rune()
				}
			case *tcell.EventResize:
				ui.resize <- 1
			}
		}
	}()

	// continuously select from various channels
	mode := modeWorking
	loop := true
	for loop {
		select {
		case ln := <-ui.write:
			ui.lines = append(ui.lines, ln)
			ui.draw(mode)
		case ln := <-ui.rewrite:
			ui.lines[len(ui.lines)-1] = ln
			ui.draw(mode)
		case ch := <-ui.input:
			if ch == 'q' || ch == rune(tcell.KeyCtrlC) || mode == modeDone {
				ui.screen.Fini()
				loop = false
			} else if mode == modePrompt {
				ui.prompt <- ch
			}
		case <-ui.resize:
			ui.screen.Sync()
			ui.draw(mode)
		case m := <-ui.change:
			mode = m
			ui.draw(mode)
		}
	}
}

// (re)draws the entire display.
func (ui *uiInstance) draw(mode uiMode) {
	ui.screen.Clear()

	// draw title bar
	w, h := ui.screen.Size()
	ui.drawLine(w, 0, ui.lines[0])
	var x int
	for x := 0; x < w; x++ {
		ui.screen.SetContent(x, 1, '─', nil, tcell.StyleDefault)
	}

	// draw content lines
	scroll := 0
	if len(ui.lines) > h-3 {
		scroll = len(ui.lines) - (h - 3)
	}
	for i, ln := range ui.lines[scroll+1:] {
		x = ui.drawLine(w, i+2, ln)
	}

	// draw bottom bar
	for x := 0; x < w; x++ {
		ui.screen.SetContent(x, h-2, '─', nil, tcell.StyleDefault)
	}
	ui.drawLine(w, h-1, uiBottom)

	// draw cursor if applicable
	if mode == modePrompt {
		ui.screen.ShowCursor(x, len(ui.lines)-scroll)
	} else {
		ui.screen.HideCursor()
	}

	ui.screen.Show()
}

// draws a line of text on the display, truncating it as needed (not wrapping
// it).
func (ui *uiInstance) drawLine(w, y int, ln uiLine) int {
	// figure out whether the line needs to be shortened
	var truncLen int
	truncIndex := -1
	lineWidth := 0
	for _, seg := range ln {
		lineWidth += len(seg.text)
	}
	if lineWidth > w-2 {
		// figure out which segment to shorten
		truncIndex = len(ln) - 1
		for i, seg := range ln {
			if seg.el != ellipsisNone {
				truncIndex = i
				break
			}
		}

		// figure out by how much to shorten it
		truncLen = len(ln[truncIndex].text) - (lineWidth - (w - 2) + 3)
		if truncLen < 0 {
			truncLen = 0
		}
	}

	// draw characters
	x := 0
	for i, seg := range ln {
		text := seg.text

		// ...text
		if i == truncIndex {
			if seg.el == ellipsisLeft {
				x = ui.drawEllipsis(x, y, seg.style)
				text = text[len(text)-truncLen:]
			} else {
				text = text[:truncLen]
			}
		}

		for _, ch := range text {
			ui.screen.SetContent(x, y, ch, nil, seg.style)
			x++
		}

		// text...
		if i == truncIndex && seg.el != ellipsisLeft {
			x = ui.drawEllipsis(x, y, seg.style)
		}
	}

	return x
}

// drawEllipsis draws "..." at the given coords and returns the new x position.
func (ui *uiInstance) drawEllipsis(x, y int, style tcell.Style) int {
	for i := 0; i < 3; i++ {
		ui.screen.SetContent(x, y, '.', nil, style)
		x++
	}
	return x
}

// adds a line to the display, formatted by fmt.Sprintf.
func (ui *uiInstance) printf(format string, a ...interface{}) {
	ui.write <- []uiSegment{{text: fmt.Sprintf(format, a...)}}
}

// adds a line to the display, with a middle path segment that is truncated at
// the left rather than the right.
func (ui *uiInstance) printPath(pre, path, post string) {
	ui.write <- []uiSegment{
		{text: pre},
		{text: path, el: ellipsisLeft},
		{text: post},
	}
}

// prints the given string and blocks until the user inputs one of the
// alphanumeric characters shown in parentheses. for example, a prompt
// containing "(y/n)" would accept either 'y' or 'n' as input. multiple
// parentheticals can be included.
func (ui *uiInstance) doPrompt(s string) rune {
	acceptedRunes := ""

	// show parentheticals in bold
	line := make([]uiSegment, 0)
	pos := 0
	for pos < len(s) {
		open := strings.IndexRune(s[pos:], '(')
		if open == -1 {
			line = append(line, uiSegment{text: s[pos:]})
			break
		} else {
			close := strings.IndexRune(s[pos+open:], ')')
			if close == -1 {
				line = append(line, uiSegment{text: s[pos:]})
				break
			} else {
				line = append(line, uiSegment{text: s[pos : pos+open]})
				line = append(line, uiSegment{
					text:  s[pos+open : pos+open+close+1],
					style: tcell.StyleDefault.Bold(true),
				})
				acceptedRunes += s[pos+open+1 : pos+open+close]
				pos += open + close + 1
			}
		}
	}

	// add space before cursor
	line = append(line, uiSegment{text: " "})

	// wait for and return a valid rune
	ui.write <- line
	ui.change <- modePrompt
	for {
		ch := <-ui.prompt
		if strings.ContainsRune(acceptedRunes, ch) {
			ui.rewrite <- append(line, uiSegment{text: string(ch)})
			ui.change <- modeWorking
			return ch
		}
	}
}

// waits for the user to input 8 hex digits, then returns the string.
func (ui *uiInstance) promptSeed(s string) string {
	acceptedRunes := "0123456789abcdef"
	line := []uiSegment{{text: s + " "}, {text: ""}}

	ui.write <- line
	ui.change <- modePrompt
	for {
		ch := <-ui.prompt
		if strings.ContainsRune(acceptedRunes, ch) {
			line[1].text += string(ch)
		} else if ch == rune(tcell.KeyDEL) && len(line[1].text) > 0 {
			line[1].text = line[1].text[:len(line[1].text)-1]
		}
		ui.rewrite <- line

		if len(line[1].text) == 8 {
			ui.change <- modeWorking
			return line[1].text
		}
	}
}

// changes the mode to one where no action is taken, and any input closes the
// program.
func (ui *uiInstance) done() {
	ui.write <- []uiSegment{{}}
	ui.write <- []uiSegment{{text: "press any key to exit.",
		style: tcell.StyleDefault.Bold(true)}}
	ui.change <- modeDone
}
