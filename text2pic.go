package text2pic

import (
	"fmt"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
)

type Color image.Image

var (
	ColorRed   = image.NewUniform(color.RGBA{0xFF, 0x00, 0x00, 0xff})
	ColorGreen = image.NewUniform(color.RGBA{0x00, 0xFF, 0x00, 0xff})
	ColorBlue  = image.NewUniform(color.RGBA{0x00, 0x00, 0xFF, 0xff})
	ColorWhite = image.White
	ColorBlack = image.Black
	TypePng = 1
	TypeJpeg = 2
)

type Padding struct {
	Left   int
	Right  int
	Bottom int
	Top    int
}

type Configure struct {
	Width int
	BgColor Color
}

func NewTextPicture(conf Configure) *TextPicture {
	pic := new(TextPicture)
	pic.conf = conf
	return pic
}

type TextPicture struct {
	text  string
	conf  Configure
	lines []line
}

func (this *TextPicture) AddTextLine(text string, fontSize float64, font *truetype.Font, color Color, padding Padding) {
	textline := new(textLine)
	textline.font = font
	textline.fontsize = fontSize
	textline.text = text
	textline.color = color
	textline.padding = padding
	this.lines = append(this.lines, textline)
}

func (this *TextPicture) AddPictureLine(reader io.Reader, padding Padding) {
	picline := new(pictureLine)
	picline.reader = reader
	picline.padding = padding
	this.lines = append(this.lines, picline)
}

func (this *TextPicture) Draw(writer io.Writer, filetype int) error {
	var err error
	// Initialize the context.
	height := 0
	width := this.conf.Width
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))

	for _, v := range this.lines {
		height += v.getHeight(width, rgba)
	}
	rgba = image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(rgba, rgba.Bounds(), this.conf.BgColor, image.ZP, draw.Src)
	pt := fixed.Point26_6{X: fixed.Int26_6(0), Y: fixed.Int26_6(0)}
	for _, v := range this.lines {
		if e := v.draw(this.conf.Width, &pt, rgba); e != nil {
			fmt.Println("draw error :", e)
		}
	}
	if filetype == TypePng {
		err = png.Encode(writer, rgba)
		if err != nil {
			return err

		}
	} else {
		err = jpeg.Encode(writer, rgba, nil)
		if err != nil {
			return err

		}
	}
	return nil
}
