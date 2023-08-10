package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

const (
	version = "0.0.4"
)

var (
	reComment = regexp.MustCompile(`//.*$`)
	reLoadout = regexp.MustCompile(`(?i).*setunitloadout *(\[[^;]+);.*`)
	reQitem   = regexp.MustCompile(`.*("[^"]+").*`)

	versioninfo = fmt.Sprintf("QtLola v%s\n© 2021 Tobias Klausmann\n\nQtLola converts simple assignGear loadouts to ACE limited arsenals.\nhttps://github.com/klausman/qtlola", version)

	showv = flag.Bool("v", false, "Show version number and exit")

	window *widgets.QMainWindow
)

func main() {
	flag.Parse()

	if *showv {
		fmt.Println(versioninfo)
		os.Exit(0)
	}

	// Create application
	app := widgets.NewQApplication(len(os.Args), os.Args)

	// Create main window
	window = widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("QTLola")
	window.SetMinimumSize2(400, 400)

	// Create main layout
	layout := widgets.NewQVBoxLayout()

	// Create main widget and set the layout
	mainWidget := widgets.NewQWidget(nil, 0)
	mainWidget.SetLayout(layout)

	fileMenu := window.MenuBar().AddMenu2("File")
	quit := fileMenu.AddAction("&Quit")
	quit.ConnectTriggered(func(checked bool) { window.Close() })
	quit.SetShortcut(gui.NewQKeySequence2("Ctrl+Q", gui.QKeySequence__NativeText))

	helpMenu := window.MenuBar().AddMenu2("Help")
	action1 := helpMenu.AddAction("About QtLola")
	action1.ConnectTriggered(func(checked bool) { aboutDialog() })
	action2 := helpMenu.AddAction("About Qt")
	action2.ConnectTriggered(func(checked bool) { app.AboutQt() })

	// Create a text edit box and add it to the layout
	input := widgets.NewQTextEdit(nil)
	input.SetPlaceholderText("Paste simple assignGear contents here")
	input.SetAcceptRichText(false)
	input.SetFontFamily("Monospace")
	layout.AddWidget(input, 0, 0)

	output := widgets.NewQTextEdit(nil)
	output.SetPlaceholderText("Limited Arsenal SQF will appear here")
	input.SetAcceptRichText(false)
	input.SetFontFamily("Monospace")
	layout.AddWidget(output, 0, 0)

	// Create a button and add it to the layout
	button := widgets.NewQPushButton2("Convert", nil)
	layout.AddWidget(button, 0, 0)

	// Connect event for button
	button.ConnectClicked(func(checked bool) {
		out := getLAfromLO(input.Document().ToPlainText())
		output.SetPlainText(out)
	})

	// Set main widget as the central widget of the window
	window.SetCentralWidget(mainWidget)

	// Show the window
	window.Show()

	// Execute app
	app.Exec()
}

func getLAfromLO(s string) string {
	items := make(map[string]bool)
	for _, line := range strings.Split(s, "\n") {
		line = strings.Trim(line, " \n\r\t")
		line = reComment.ReplaceAllString(line, "")
		lo := reLoadout.FindStringSubmatch(line)
		if len(lo) == 0 {
			continue
		}
		for _, tok := range strings.Split(lo[1], ",") {
			stripped := reQitem.FindStringSubmatch(tok)
			if len(stripped) == 0 {
				continue
			}
			items[stripped[1]] = true
		}
	}
	itemlist := make([]string, 0, len(items))
	for k := range items {
		itemlist = append(itemlist, k)
	}
	sort.Strings(itemlist)
	return fmt.Sprintf("[this, [\n    %s]\n] call ace_arsenal_fnc_initBox;\n", strings.Join(itemlist, ",\n    "))
}

func aboutDialog() {
	widgets.QMessageBox_About(window, "About", versioninfo)
}
