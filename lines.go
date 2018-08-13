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
)

var (
	dpi     = flag.Float64("dpi", 144, "screen resolution in Dots Per Inch")
	hinting = flag.String("hinting", "none", "none | full")
	size    = flag.Float64("size", 16, "font size in points")
	spacing = flag.Float64("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")
)

type line interface {
	draw(width int,  pt *fixed.Point26_6, image draw.Image) error
	getHeight(width int, image draw.Image) int
}

type textLine struct {
	fontsize float64
	text     string
	font     *truetype.Font
	padding  int
	color 	 image.Image
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

	//line define
	pts := freetype.Pt(0, 0+int(c.PointToFixed(this.fontsize)>>6))
	pt.X = pts.X
	pt.Y = pt.Y + pts.Y
	fmt.Println("pt:",pt, "- pts:", pts)
	index := 0
	lastwidth := 0
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
		}
		lastwidth = rpts.X.Floor()

	}
	//add padding
	padding := freetype.Pt(0, 0+this.padding)
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
	lines := (rpt.X.Floor() / width)
	//check the left text
	if rpt.X.Floor() %width > 0 {
		lines++
	}

	fmt.Println("------", lines)
	//add padding
	return int(c.PointToFixed(this.fontsize)>>6) *lines + 10 + this.padding
}



type reader struct {
	buf []byte
}

func (this *reader) Read(p []byte) (n int, err error){
	p = this.buf
	return len(p), nil
}



type pictureLine struct {
	reader  io.Reader
	padding int
	img	image.Image
}

func (this *pictureLine) draw(width int, pt *fixed.Point26_6, img draw.Image) error {
	fmt.Println("picture:",pt)
	draw.Draw(img, this.img.Bounds().Add(image.Pt(0, pt.Y.Floor())) ,this.img, this.img.Bounds().Min  ,draw.Src)
	//add padding
	padding := freetype.Pt(0, 0+this.padding + this.img.Bounds().Max.Y)
	pt.Y = pt.Y + padding.Y

	return nil
}

func (this *pictureLine) getHeight(width int, image draw.Image) int {
	img, e := jpeg.Decode(this.reader)

	//resize
	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(uint(width), 0, img, resize.Lanczos3)

	this.img = m
	fmt.Println(img.Bounds(), m.Bounds(), e)
	return m.Bounds().Size().Y + this.padding
}
