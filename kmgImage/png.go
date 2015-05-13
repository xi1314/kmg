package kmgImage

import (
	"image"
	"image/png"
	"os"
)

func PngDecodeConfigFromFile(path string) (conf image.Config, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	return png.DecodeConfig(file)
}
