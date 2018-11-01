package ui

import (
	"fmt"
	"github.com/nsf/termbox-go"
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
	modeDone
)

var (
	lines  []line
	bottom = []segment{{text: "(q)", fg: colorDefault | bold}, {text: "uit"}}
	write  = make(chan line, 1)
	input  = make(chan rune)
	resize = make(chan interface{}, 1)
	change = make(chan modeType, 1)
)

func Init(title string) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.HideCursor()

	lines = []line{
		[]segment{{text: title}},
	}
	draw()
}

func Run() {
	go func() {
		for {
			evt := termbox.PollEvent()
			switch evt.Type {
			case termbox.EventKey:
				input <- evt.Ch
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
			draw()
		case ch := <-input:
			if ch == 'q' || mode == modeDone {
				termbox.Close()
				loop = false
			}
		case <-resize:
			draw()
		case m := <-change:
			mode = m
		}
	}
}

func draw() {
	termbox.Clear(colorDefault, colorDefault)

	w, h := termbox.Size()
	drawLine(w, 0, lines[0])
	for x := 0; x < w; x++ {
		termbox.SetCell(x, 1, '─', colorDefault, colorDefault)
	}
	for i, ln := range lines[1:] {
		drawLine(w, i+2, ln)
	}
	for x := 0; x < w; x++ {
		termbox.SetCell(x, h-2, '─', colorDefault, colorDefault)
	}
	drawLine(w, h-1, bottom)

	termbox.Flush()
}

func drawLine(w, y int, ln line) {
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

func Done() {
	write <- []segment{{text: "press any key to exit.",
		fg: colorDefault | bold}}
	change <- modeDone
}
