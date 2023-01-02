package main

import (
	"flag"
	"os"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/flynch/rpi-rgb-led-matrix/bindings/go/ledmatrix"
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
	img                      = flag.String("image", "", "image path")

	rotate = flag.Int("rotate", 0, "rotate angle, 90, 180, 270")
)

func main() {
	f, err := os.Open(*img)
	fatal(err)

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

	tk := ledmatrix.NewToolKit(m)
	defer tk.Close()

	switch *rotate {
	case 90:
		tk.Transform = imaging.Rotate90
	case 180:
		tk.Transform = imaging.Rotate180
	case 270:
		tk.Transform = imaging.Rotate270
	}

	if strings.HasSuffix(*img, "gif") {
		close, err := tk.PlayGIF(f)
		fatal(err)

		time.Sleep(time.Second * 30)
		close <- true
	} else if strings.HasSuffix(*img, "jpg") {
		err = tk.PlayJpeg(f)
		fatal(err)
	} else if strings.HasSuffix(*img, "png") {
		err = tk.PlayPng(f)
		fatal(err)
	}

}

func init() {
	flag.Parse()
}

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}
