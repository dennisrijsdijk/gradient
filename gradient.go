package gradient

import (
	"errors"
	"github.com/disintegration/imaging"
	"github.com/mazznoer/colorgrad"
	"github.com/ojrac/opensimplex-go"
	"image"
)

type Options struct {
	Width      int
	Height     int
	Type       string // "basic", "tilted" or "noise"
	TiltAngle  float64
	NoiseSeed  int64
	HtmlColors []string
}

func basic(width int, height int, colors ...string) (image.Image, error) {
	grad, err := colorgrad.NewGradient().HtmlColors(colors...).Build()
	if err != nil {
		return nil, err
	}

	fWidth := float64(width)

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for x := 0; x < width; x++ {
		col := grad.At(float64(x) / fWidth)
		for y := 0; y < height; y++ {
			img.Set(x, y, col)
		}
	}

	return img, nil
}

func basicTilted(width int, height int, angle float64, colors ...string) (image.Image, error) {
	var size int
	if !(height < width) {
		size = int(float64(height) * 1.5)
	} else {
		size = int(float64(width) * 1.5)
	}
	img, err := basic(size, size, colors...)
	if err != nil {
		return nil, err
	}
	imgRot := imaging.Rotate(img, angle, image.Black)
	return imaging.CropCenter(imgRot, height, width), nil
}

func noise(width int, height int, seed int64, colors ...string) (image.Image, error) {
	scale := 0.02

	grad, err := colorgrad.NewGradient().HtmlColors(colors...).Build()
	if err != nil {
		return nil, err
	}
	grad = grad.Sharp(12, 0.2)
	noise := opensimplex.NewNormalized(seed)
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			t := noise.Eval2(float64(x)*scale, float64(y)*scale)
			img.Set(x, y, grad.At(t))
		}
	}
	return img, nil
}

func Draw(options Options) (image.Image, error) {
	var img image.Image
	var err error
	if options.Width == 0 || options.Height == 0 {
		return nil, errors.New("height and width must be non-0")
	} else if len(options.HtmlColors) == 0 {
		return nil, errors.New("colors must not be empty")
	}
	switch options.Type {
	case "basic":
		img, err = basic(options.Width, options.Height, options.HtmlColors...)
	case "tilted":
		img, err = basicTilted(options.Width, options.Height, options.TiltAngle, options.HtmlColors...)
	case "noise":
		img, err = noise(options.Width, options.Height, options.NoiseSeed, options.HtmlColors...)
	default:
		return nil, errors.New("type must be one of 'basic', 'tilted' or 'noise'")
	}
	return img, err
}
