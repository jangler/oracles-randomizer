package ui

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"strings"
)

const (
	colorDefault = termbox.ColorDefault // going to be typing this a lot
	bold         = termbox.AttrBold
)

type line []segment

type segment struct {
	fg, bg termbox.Attribute
	text   string
	el     ellipsis
}

// defines whether a segment will opt in to being truncated at the left or
// right to fit onscreen. by default, lines are truncated at the right, but if
// (for example) a segment in the line has ellipsisLeft, that segment will be
// truncated at the left in order to make the line fit.
type ellipsis int

const (
	ellipsisNone ellipsis = iota
	ellipsisLeft
	ellipsisRight
)

type modeType int

const (
	modeWorking modeType = iota
	modePrompt
	modeDone
)

var (
	lines   []line
	bottom  = []segment{{text: "(q)", fg: colorDefault | bold}, {text: "uit"}}
	write   = make(chan line, 1)
	rewrite = make(chan line)
	input   = make(chan rune)
	prompt  = make(chan rune)
	resize  = make(chan interface{}, 1)
	change  = make(chan modeType, 1)
)

func Init(title string) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	lines = []line{
		[]segment{{text: title}},
	}
	draw(modeWorking)
}

func Run() {
	go func() {
		for {
			evt := termbox.PollEvent()
			switch evt.Type {
			case termbox.EventKey:
				switch evt.Key {
				case termbox.KeyCtrlC, '\x7f': // 7f == backspace
					input <- rune(evt.Key)
				default:
					input <- evt.Ch
				}
			case termbox.EventResize:
				resize <- 1
			}
		}
	}()

	mode := modeWorking
	loop := true
	for loop {
		select {
		case ln := <-write:
			lines = append(lines, ln)
			draw(mode)
		case ln := <-rewrite:
			lines[len(lines)-1] = ln
			draw(mode)
		case ch := <-input:
			if ch == 'q' || ch == '\x03' || mode == modeDone {
				termbox.Close()
				loop = false
			} else if mode == modePrompt {
				prompt <- ch
			}
		case <-resize:
			draw(mode)
		case m := <-change:
			mode = m
			draw(mode)
		}
	}
}

func draw(mode modeType) {
	termbox.Clear(colorDefault, colorDefault)

	w, h := termbox.Size()
	drawLine(w, 0, lines[0])
	var x int
	for x := 0; x < w; x++ {
		termbox.SetCell(x, 1, '─', colorDefault, colorDefault)
	}

	scroll := 0
	if len(lines) > h-3 {
		scroll = len(lines) - (h - 3)
	}
	for i, ln := range lines[scroll+1:] {
		x = drawLine(w, i+2, ln)
	}

	for x := 0; x < w; x++ {
		termbox.SetCell(x, h-2, '─', colorDefault, colorDefault)
	}
	drawLine(w, h-1, bottom)

	if mode == modePrompt {
		termbox.SetCursor(x, len(lines)-scroll)
	} else {
		termbox.HideCursor()
	}

	termbox.Flush()
}

func drawLine(w, y int, ln line) int {
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
				x = drawEllipsis(x, y, seg.fg, seg.bg)
				text = text[len(text)-truncLen:]
			} else {
				text = text[:truncLen]
			}
		}

		for _, ch := range text {
			termbox.SetCell(x, y, ch, seg.fg, seg.bg)
			x++
		}

		// text...
		if i == truncIndex && seg.el != ellipsisLeft {
			x = drawEllipsis(x, y, seg.fg, seg.bg)
		}
	}

	return x
}

func drawEllipsis(x, y int, fg, bg termbox.Attribute) int {
	for i := 0; i < 3; i++ {
		termbox.SetCell(x, y, '.', fg, bg)
		x++
	}
	return x
}

func Printf(format string, a ...interface{}) {
	write <- []segment{{text: fmt.Sprintf(format, a...)}}
}

func PrintPath(pre, path, post string) {
	write <- []segment{
		{text: pre},
		{text: path, el: ellipsisLeft},
		{text: post},
	}
}

// Prompt prints the given string and blocks until the user inputs one of the
// alphanumeric characters shown in parentheses. For example, a prompt
// containing "(y/n)" would accept either 'y' or 'n' as input. Multiple
// parentheticals can be included.
func Prompt(s string) rune {
	acceptedRunes := ""

	// show parentheticals in bold
	line := make([]segment, 0)
	pos := 0
	for pos < len(s) {
		open := strings.IndexRune(s[pos:], '(')
		if open == -1 {
			line = append(line, segment{text: s[pos:]})
			break
		} else {
			close := strings.IndexRune(s[pos+open:], ')')
			if close == -1 {
				line = append(line, segment{text: s[pos:]})
				break
			} else {
				line = append(line, segment{text: s[pos : pos+open]})
				line = append(line, segment{
					text: s[pos+open : pos+open+close+1],
					fg:   colorDefault | bold,
				})
				acceptedRunes += s[pos+open+1 : pos+open+close]
				pos += open + close + 1
			}
		}
	}

	// add space before cursor
	line = append(line, segment{text: " "})

	write <- line
	change <- modePrompt
	for {
		ch := <-prompt
		if strings.ContainsRune(acceptedRunes, ch) {
			change <- modeWorking
			return ch
		}
	}
}

// PromptSeed waits for the user to input 8 hex digits, then returns the
// string.
func PromptSeed(s string) string {
	acceptedRunes := "0123456789abcdef"
	line := []segment{{text: s + " "}, {text: ""}}

	write <- line
	change <- modePrompt
	for {
		ch := <-prompt
		if strings.ContainsRune(acceptedRunes, ch) {
			line[1].text += string(ch)
		} else if ch == '\x7f' && len(line[1].text) > 0 {
			line[1].text = line[1].text[:len(line[1].text)-1]
		}
		rewrite <- line

		if len(line[1].text) == 8 {
			change <- modeWorking
			return line[1].text
		}
	}
}

func Done() {
	write <- []segment{{}}
	write <- []segment{{text: "press any key to exit.",
		fg: colorDefault | bold}}
	change <- modeDone
}
