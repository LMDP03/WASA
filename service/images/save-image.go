package images

import (
	"image/jpeg"
	"os"

	"github.com/nfnt/resize"
)

func SaveImage(filename string, w uint, h uint) error {

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() { err = file.Close() }()

	img, err := jpeg.Decode(file)

	if err != nil {
		return err
	}

	resizedImg := resize.Resize(w, h, img, resize.NearestNeighbor)

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() { err = out.Close() }()

	err = jpeg.Encode(out, resizedImg, nil)
	if err != nil {
		return err
	}

	return err

}
