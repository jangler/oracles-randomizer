package ui

import (
	"fmt"
	"github.com/nsf/termbox-go"
)

const colorDefault = termbox.ColorDefault // going to be typing this a lot

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
			_ = ch
			if mode == modeDone {
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

	w, _ := termbox.Size()
	drawLine(w, 0, lines[0])
	for x := 0; x < w; x++ {
		termbox.SetCell(x, 1, '─', colorDefault, colorDefault)
	}
	for i, ln := range lines[1:] {
		drawLine(w, i+2, ln)
	}

	termbox.Flush()
}

func drawLine(w, y int, ln line) {
	lineWidth := 0
	for _, seg := range ln {
		lineWidth += len(seg.text)
	}

	if lineWidth < w-1 {
		x := 0
		for _, seg := range ln {
			for _, ch := range seg.text {
				termbox.SetCell(x, y, ch, seg.fg, seg.bg)
				x++
			}
		}
	}
}

func Printf(format string, a ...interface{}) {
	write <- []segment{{text: fmt.Sprintf(format, a...)}}
}

func Done() {
	write <- []segment{{text: "press any key to exit."}}
	change <- modeDone
}
