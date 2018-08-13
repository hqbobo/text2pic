package text2pic

import (
	"image"
	"image/draw"
	"golang.org/x/image/math/fixed"
	"github.com/golang/freetype/truetype"
	"os"
	"image/png"
	"image/jpeg"
	"image/color"
	"bufio"
	"fmt"
	"io"
)

type Color image.Image

var (
	ColorRed = image.NewUniform(color.RGBA{0xFF, 0x00, 0x00, 0xff})
	ColorGreen = image.NewUniform(color.RGBA{0x00, 0xFF, 0x00, 0xff})
	ColorBlue = image.NewUniform(color.RGBA{0x00, 0x00, 0xFF, 0xff})
	ColorWhite = image.White
	ColorBlack = image.Black
)

type Configure struct {
	Width  int
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

func (this *TextPicture) AddTextLine(text string, fontSize float64, font *truetype.Font, color Color, padding int) {
	textline := new(textLine)
	textline.font = font
	textline.fontsize = fontSize
	textline.text = text
	textline.color = color
	textline.padding = padding
	this.lines = append(this.lines, textline)
}

func (this *TextPicture) AddPictureLine(reader io.Reader, padding int) {
	picline := new(pictureLine)
	picline.reader = reader
	picline.padding = padding
	this.lines = append(this.lines, picline)
}

type writer struct {
	buf []byte
}

func (this *writer) Write(p []byte) (n int, err error){
	this.buf = make([]byte, len(p))
	this.buf = p
	return len(this.buf), nil
}

func (this *writer) Get()[]byte  {
	return this.buf
}


func (this *TextPicture) Draw() error {
	// Initialize the context.
	bg := image.White

	height := 0
	width := this.conf.Width

	rgba := image.NewRGBA(image.Rect(0, 0, width, height))

	for _ , v :=range this.lines {
		height += v.getHeight(width ,rgba)
	}
	rgba = image.NewRGBA(image.Rect(0, 0, width, height))
	fmt.Println("bg size:", width, ":", height)
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)

	pt := fixed.Point26_6{X:fixed.Int26_6(0), Y:fixed.Int26_6(0)}
	for _ , v :=range this.lines {
		if e := v.draw(this.conf.Width, &pt , rgba); e != nil {
			fmt.Println("draw error :", e)
		}
	}


	// Save that RGBA image to disk.
	outFile, err := os.Create("out.png")
	if err != nil {
		return err
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	if err != nil {
		return err

	}

	err = b.Flush()
	if err != nil {
		return err

	}
	buf := new(writer)
	b1 := bufio.NewWriter(buf)
	err = jpeg.Encode(b1, rgba, nil)
	if err != nil {
		return err

	}
	err = b1.Flush()
	if err != nil {
		return err

	}
	//fmt.Println(buf.Get())
	fmt.Println("Wrote out.png OK.")
	return nil
}