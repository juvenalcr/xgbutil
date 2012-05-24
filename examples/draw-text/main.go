// Example draw-text shows how to draw text to an xgraphics.Image type.
package main

import (
	"image"
	"log"
	"os"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
)

var (
	// The geometry of the canvas to draw text on.
	canvasWidth, canvasHeight = 600, 100

	// The background color of the canvas.
	bg = xgraphics.BGRA{B: 0xff, G: 0x66, R: 0x33, A: 0xff}

	// The path to the font used to draw text.
	fontPath = "/usr/share/fonts/TTF/FreeMonoBold.ttf"

	// The color of the text.
	fg = xgraphics.BGRA{B: 0xff, G: 0xff, R: 0xff, A: 0xff}

	// The size of the text.
	size = 20.0

	// The text to draw.
	msg = "This is some text drawn by xgraphics!"
)

func main() {
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	// Load some font. You may need to change the path depending upon your
	// system configuration.
	fontReader, err := os.Open(fontPath)
	if err != nil {
		log.Fatal(err)
	}

	// Now parse the font.
	font, err := xgraphics.ParseFont(fontReader)
	if err != nil {
		log.Fatal(err)
	}

	// Create some canvas.
	ximg := xgraphics.New(X, image.Rect(0, 0, canvasWidth, canvasHeight))
	ximg.For(func(x, y int) xgraphics.BGRA {
		return bg
	})

	// Now write the text. The x,y returned is the position at the end of
	// the text drawn.
	_, y, err := ximg.Text(10, 10, fg, size, font, msg)
	if err != nil {
		log.Fatal(err)
	}

	// Now show the image in its own window.
	win := ximg.XShowExtra("Drawing text using xgraphics", true)

	// Now draw some more text below the above and demonstrate how to update
	// only the region we've updated.
	x2, y2, err := ximg.Text(10, y+10, fg, size, font, "Some more text.")
	if err != nil {
		log.Fatal(err)
	}

	// Now repaint on the region that we drew text on. Then update the screen.
	ximg.SubImage(image.Rect(10, y+10, x2, y2)).XDraw()
	ximg.XPaint(win.Id)

	// All we really need to do is block, which could be achieved using
	// 'select{}'. Invoking the main event loop however, will emit error
	// message if anything went seriously wrong above.
	xevent.Main(X)
}
