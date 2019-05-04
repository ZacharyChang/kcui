package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ZacharyChang/kcui/pkg/log"
	"github.com/ZacharyChang/kcui/pkg/option"
	"github.com/ZacharyChang/kcui/pkg/page"
	"github.com/ZacharyChang/kcui/version"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"gopkg.in/urfave/cli.v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func main() {
	app := cli.NewApp()
	app.Name = "KCUI"
	app.Compiled = time.Now()
	app.Usage = "A simple tool to monitor the k8s pods and logs."
	app.Version = version.Version

	opts := option.NewOptions()
	opts.AddFlags(app)

	app.Action = func(c *cli.Context) error {
		if err := startView(opts); err != nil {
			return cli.NewExitError(fmt.Sprintf("application failed to start: %v\n", err), 1)
		}
		log.Info("application started...")
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Errorf("fail to run: %s", err.Error())
	}

}

type Slide func(options *option.Options, stopCh chan struct{}) (title string, content tview.Primitive)

// The application.
var app = tview.NewApplication()

func startView(opts *option.Options) error {
	// The presentation slides.
	slides := []Slide{
		page.Log,
	}

	// The bottom row has some info on where we are.
	info := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false)

	// Create the pages for all slides.
	currentSlide := 0
	info.Highlight(strconv.Itoa(currentSlide))
	pages := tview.NewPages()
	previousSlide := func() {
		currentSlide = (currentSlide - 1 + len(slides)) % len(slides)
		info.Highlight(strconv.Itoa(currentSlide)).
			ScrollToHighlight()
		pages.SwitchToPage(strconv.Itoa(currentSlide))
	}
	nextSlide := func() {
		currentSlide = (currentSlide + 1) % len(slides)
		info.Highlight(strconv.Itoa(currentSlide)).
			ScrollToHighlight()
		pages.SwitchToPage(strconv.Itoa(currentSlide))
	}
	stopCh := make(chan struct{})
	defer close(stopCh)

	for index, slide := range slides {
		title, primitive := slide(opts, stopCh)
		pages.AddPage(strconv.Itoa(index), primitive, true, index == currentSlide)
		_, _ = fmt.Fprintf(info, `%d ["%d"][darkcyan]%s[white][""]  `, index+1, index, title)
	}

	// Create the main layout.
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(pages, 0, 1, true).
		AddItem(info, 1, 1, false)

	// Shortcuts to navigate the slides.
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlN {
			nextSlide()
		} else if event.Key() == tcell.KeyCtrlP {
			previousSlide()
		}
		if event.Key() == tcell.KeyESC {
			app.Stop()
		}
		if event.Key() == tcell.KeyCtrlQ {
			app.Stop()
		}
		return event
	})

	// Refresh every 0.5 second
	go wait.Forever(func() {
		app.Draw()
	}, time.Millisecond*500)

	// Start the application.
	if err := app.SetRoot(layout, true).Run(); err != nil {
		panic(err)
	}
	return nil
}
