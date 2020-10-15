package randomizer

import (
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func runGUI() {
	app := app.New()

	w := app.NewWindow("oracles-randomizer " + version)
	w.SetContent(widget.NewLabel("hello, dev!"))

	w.ShowAndRun()
}
