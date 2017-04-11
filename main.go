// Example fullscreen shows how to make a window showing the Go Gopher go
// fullscreen and back using keybindings and EWMH.
package main

import (
	"bytes"
	"image"
	_ "image/png"
	"log"
	"os"
	"strconv"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgb/screensaver"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/gopher"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"
)

//xproto.QueryTree(w.X.Conn(), w.Id).Reply()
//[]string{"/usr/lib/xscreensaver/moire", "-window-id", "81788931"}
func main() {

	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}
	//Basically this is where we get the window id
	screensaver.Init(X.Conn())
	if len(os.Args) > 1 {
		i, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		X.RootWinSet(xproto.Window(i))
	}

	keybind.Initialize(X) // call once before using keybind package

	// Read an example gopher image into a regular png image.
	img, _, err := image.Decode(bytes.NewBuffer(gopher.GopherPng()))
	if err != nil {
		log.Fatal(err)
	}

	// Now convert it into an X image.
	ximg := xgraphics.NewConvert(X, img)
	//	w, h := ximg.Rect.Dx(), ximg.Rect.Dy()
	ximg.XShow()

	// Now show it in a new window.
	// We set the window title and tell the program to quit gracefully when
	// the window is closed.
	// There is also a convenience method, XShow, that requires no parameters.
	win := xwindow.New(X, X.RootWin())
	win.WMGracefulClose(func(w *xwindow.Window) {
		xevent.Detach(w.X, w.Id)
		keybind.Detach(w.X, w.Id)
		mousebind.Detach(w.X, w.Id)
		w.Destroy()
		xevent.Quit(w.X)
	})
	// Listen for key press events.
	win.Listen(xproto.EventMaskKeyPress)
	win.Map()

	err = keybind.KeyPressFun(
		func(X *xgbutil.XUtil, ev xevent.KeyPressEvent) {
			println("fullscreen!")
			err := ewmh.WmStateReq(X, win.Id, ewmh.StateToggle,
				"_NET_WM_STATE_FULLSCREEN")
			if err != nil {
				log.Fatal(err)
			}
		}).Connect(X, win.Id, "f", false)
	if err != nil {
		log.Fatal(err)
	}

	err = keybind.KeyPressFun(
		func(X *xgbutil.XUtil, ev xevent.KeyPressEvent) {
			os.Exit(0)
		}).Connect(X, win.Id, "Escape", false)
	if err != nil {
		log.Fatal(err)
	}

	xevent.DestroyNotifyFun(
		func(xu *xgbutil.XUtil, event xevent.DestroyNotifyEvent) {
			os.Exit(0)
		}).Connect(X, win.Id)

	xevent.FocusOutFun(
		func(xu *xgbutil.XUtil, event xevent.FocusOutEvent) {
			os.Exit(0)
		}).Connect(X, win.Id)

	xevent.LeaveNotifyFun(
		func(xu *xgbutil.XUtil, event xevent.LeaveNotifyEvent) {
			os.Exit(0)
		}).Connect(X, win.Id)

	xevent.MotionNotifyFun(
		func(xu *xgbutil.XUtil, event xevent.MotionNotifyEvent) {
			os.Exit(0)
		}).Connect(X, win.Id)

	// If we don't block, the program will end and the window will disappear.
	// We could use a 'select{}' here, but xevent.Main will emit errors if
	// something went wrong, so use that instead.
	xevent.Main(X)
}
