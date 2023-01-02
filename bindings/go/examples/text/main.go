package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os/user"
	"time"

	"github.com/flynch/rpi-rgb-led-matrix/bindings/go/ledmatrix"
	"github.com/fogleman/gg"
)

var (
	rows                     = flag.Int("led-rows", 32, "number of rows supported")
	cols                     = flag.Int("led-cols", 64, "number of columns supported")
	parallel                 = flag.Int("led-parallel", 1, "number of daisy-chained panels")
	chain                    = flag.Int("led-chain", 2, "number of displays daisy-chained")
	brightness               = flag.Int("brightness", 50, "brightness (0-100)")
	hardware_mapping         = flag.String("led-gpio-mapping", "adafruit-hat", "Name of GPIO mapping used.")
	show_refresh             = flag.Bool("led-show-refresh", false, "Show refresh rate.")
	inverse_colors           = flag.Bool("led-inverse", false, "Switch if your matrix has inverse colors on.")
	disable_hardware_pulsing = flag.Bool("led-no-hardware-pulse", false, "Don't use hardware pin-pulse generation.")
	fontpath                 = flag.String("led-font", "./SpaceMono-Regular.ttf", "Path to font file")
	text                     = flag.String("led-text", "", "Text to display")
	fontsize                 = flag.Int("led-fontsize", 10, "Font size")
)

func isRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("[isRoot] Unable to get current user: %s", err)
	}
	return currentUser.Username == "root"
}

func main() {
	config := &ledmatrix.DefaultConfig
	config.Rows = *rows
	config.Cols = *cols
	config.Parallel = *parallel
	config.ChainLength = *chain
	config.Brightness = *brightness
	config.HardwareMapping = *hardware_mapping
	config.ShowRefreshRate = *show_refresh
	config.InverseColors = *inverse_colors
	config.DisableHardwarePulsing = *disable_hardware_pulsing
	config.PixelMapping = "V-mapper"
	m, err := ledmatrix.NewRGBLedMatrix(config)

	fatal(err)

	c := ledmatrix.NewCanvas(m)
	defer c.Close()

	fmt.Println("in main")

	// using the standard draw.Draw function we copy a white image onto the Canvas

	//put string on image
	path := *fontpath //"path/to/ttffile"
	img := gg.NewContext(c.Bounds().Max.X, c.Bounds().Max.Y)
	img.SetColor(color.RGBA{255, 255, 0, 255})
	if err = img.LoadFontFace(path, float64(*fontsize)); err != nil {
		fmt.Println("Wrong font!")
		panic(err)
	}

	fmt.Print(len(*text))

	switch len(*text) {
	case 1: //if character is single then displaying it in center
		img.DrawStringWrapped(
			*text,
			float64(c.Bounds().Max.X/2),
			float64(c.Bounds().Max.Y/2),
			0.5,
			0.5,
			float64(30),
			1.5,
			gg.AlignCenter,
		)
	default: // if it is word then displaying it in top row of LED
		img.DrawStringWrapped(
			*text,
			float64(c.Bounds().Min.X/2+10/2)-4,
			float64(c.Bounds().Min.Y/2+10/2)-4,
			0,
			0,
			float64(50),
			1.5,
			gg.AlignLeft,
		)
	}

	img.Fill()

	draw.Draw(c, c.Bounds(), img.Image(), image.ZP, draw.Over)

	// don't forget call Render to display the new led status
	c.Render()
	time.Sleep(5 * time.Second)
}

func init() {
	flag.Parse()
}

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}
