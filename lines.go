package text2pic

import (
	"flag"
	"golang.org/x/image/math/fixed"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image"
	"image/draw"
	"fmt"
	"image/jpeg"
	"io"
	"github.com/nfnt/resize"
	"image/png"
	"errors"
	"bufio"
	"bytes"
)
var (
	dpi     = flag.Float64("dpi", 288, "screen resolution in Dots Per Inch")
	hinting = flag.String("hinting", "none", "none | full")
	size    = flag.Float64("size", 8, "font size in points")
	spacing = flag.Float64("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")
	textpadding = 25
)

type line interface {
	draw(width int,  pt *fixed.Point26_6, image draw.Image) error
	getHeight(width int, image draw.Image) int
}

type textLine struct {
	fontsize float64
	text     string
	font     *truetype.Font
	padding  Padding
	color 	 image.Image
	lines    int
}

func (this *textLine) draw(width int, pt *fixed.Point26_6, image draw.Image) error {
	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetFont(this.font)
	c.SetFontSize(float64(this.fontsize))
	c.SetClip(image.Bounds())
	c.SetDst(image)
	c.SetSrc(this.color)

	//switch *hinting {
	//default:
	c.SetHinting(font.HintingFull)
	//case "full":
	//	c.SetHinting(font.HintingFull)
	//}
	// Draw the guidelines.
	//for i := 0; i < 200; i++ {
	//	image.Set(10, 10+i, ruler)
	//	image.Set(10+i, 10, ruler)
	//}

	//line define adjust padding
	//ajust line space
	width -= this.padding.Right
	toppadding := freetype.Pt(0, this.padding.Top)
	pts := freetype.Pt(this.padding.Left, this.padding.LineSpace+int(c.PointToFixed(this.fontsize)>>6))
	pt.X = pts.X
	pt.Y = pt.Y + pts.Y + toppadding.Y //add paddingtop
	index := 0
	lastwidth := 0
	lc := 1
	for i := 1 ; i <= len(this.text); i++ {

		//test the text length
		//fmt.Println("i[", i ,"]len:",len(this.text),"-",string([]byte(this.text)[index:i]))
		//fmt.Println([]byte(this.text)[index:i])

		rpts , err := c.DrawString(string([]byte(this.text)[index:i]), *pt)
		if err != nil {
			return err
		}
		//fmt.Println(rpts.X.Floor() ,":", width,":",lastwidth, "[",[]byte(this.text)[i-1],"]")

		//english char
		if []byte(this.text)[i-1] >33 && []byte(this.text)[i-1] < 127 && (rpts.X.Floor() + int(c.PointToFixed(this.fontsize)>>6) > width) {
			index = i
			pt.Y = pt.Y + pts.Y
			lc++
			continue
		}
		//other encoding
		if (rpts.X.Floor() + int(c.PointToFixed(this.fontsize)>>6) > width) && (lastwidth > width) {
			_ , err := c.DrawString(string([]byte(this.text)[index:i-1]), *pt)
			if err != nil {
				return err
			}
			index = i
			pt.Y = pt.Y + pts.Y
			lc++
		}
		lastwidth = rpts.X.Floor()
	}
	//add extraline
	for ; this.lines - lc > 0; lc++ {
		pt.Y = pt.Y + pts.Y
	}

	//add buttom padding
	padding := freetype.Pt(0, 0+this.padding.Bottom + textpadding)
	pt.Y = pt.Y + padding.Y
	return nil
}

func (this *textLine) getHeight(width int, image draw.Image) int {
	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetFont(this.font)
	c.SetFontSize(float64(this.fontsize))
	c.SetClip(image.Bounds())
	c.SetDst(image)
	c.SetHinting(font.HintingFull)
	pt := freetype.Pt(0, 0)
	rpt , err := c.DrawString(this.text, pt)
	if err != nil {
		return 0
	}
	//calculate lines count
	lines := (rpt.X.Floor() / (width - this.padding.Left - this.padding.Right))
	//check the left text
	if rpt.X.Floor() %width > 0 {
		lines++
	}

	//add extraline avoid missing char
	if lines >=5 {
		lines++
	}
	this.lines = lines
	//add padding
	return int(c.PointToFixed(this.fontsize)>>6) *lines + lines *this.padding.LineSpace + this.padding.Bottom + this.padding.Top + textpadding
}

type pictureLine struct {
	reader  io.Reader
	padding Padding
	img	image.Image
}

func (this *pictureLine) draw(width int, pt *fixed.Point26_6, img draw.Image) error {
	if this.img == nil {
		return errors.New("err pic")
	}
	draw.Draw(img, this.img.Bounds().Add(image.Pt(width/10, pt.Y.Floor())) ,this.img, this.img.Bounds().Min  ,draw.Src)
	//add padding
	padding := freetype.Pt(0, 0+this.padding.Bottom + this.img.Bounds().Max.Y)
	pt.Y = pt.Y + padding.Y
	return nil
}

func (this *pictureLine) getHeight(width int, image draw.Image) int {
	//try png and jpeg to decode
	reader := bufio.NewReader(this.reader)
	buf := make ([]byte, 10 * 1024 * 1024)
	n , e := reader.Read(buf)
	img, e := jpeg.Decode(bytes.NewReader(buf[0:n]))
	if e != nil {
		img, e  = png.Decode(bytes.NewReader(buf[0:n]))
		if e != nil {
			fmt.Println("failed to parse image file")
			return 0
		}
	}
	// resize to 80% of the width using Lanczos resampling
	m := resize.Resize(uint(width/5*4), 0, img, resize.Lanczos3)
	this.img = m
	return m.Bounds().Size().Y + this.padding.Bottom + this.padding.Top
}