package main

import (
	// "fmt"
	// "io/ioutil"
	// "net/http"
	"os"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	// "github.com/diamondburned/gotk4/pkg/glib/v2"
	// "github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func main() {
	app := adw.NewApplication("com.github.diamondburned.gotk4-examples.gtk4.simple", gio.ApplicationFlagsNone)
	app.ConnectActivate(func() { activate(app) })

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *adw.Application) {
	bld := gtk.NewBuilderFromFile("./blueprint/main_window.xml")
	win := bld.GetObject("mainWindow").Cast().(*adw.Window)

	app.AddWindow(&win.Window)
	win.Show()
}
