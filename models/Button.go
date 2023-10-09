package models

import (
	"fmt"
	"image"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

type Button struct {
	Bounds pixel.Rect
	Text   string
	Sprite *pixel.Sprite
}

func NewButton(text, imagePath string, x, y, width, height float64) *Button {
	pic, err := LoadPicture(imagePath)
	if err != nil {
		panic(fmt.Sprintf("Error al cargar %s: %v", imagePath, err))
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())

	return &Button{
		Bounds: pixel.R(x, y, x+width, y+height),
		Text:   text,
		Sprite: sprite,
	}
}

func LoadPicture(path string) (pixel.Picture, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return pixel.PictureDataFromImage(img), nil
}

func (b *Button) Draw(win *pixelgl.Window, atlas *text.Atlas, scale float64) {

	b.Sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, scale).Moved(b.Bounds.Center()))

	txt := text.New(b.Bounds.Center(), atlas)
	txt.Color = colornames.White //
	txt.Dot = txt.Dot.Sub(txt.BoundsOf(b.Text).Center())
	fmt.Fprint(txt, b.Text)
	txt.Draw(win, pixel.IM)
}

func (b *Button) Contains(point pixel.Vec) bool {
	return b.Bounds.Contains(point)
}
