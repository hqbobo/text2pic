//package main
//
//import (
//	"fmt"
//	"image"
//	"image/draw"
//	"image/jpeg"
//	"os"
//)
//
//func main() {
//	file, err := os.Create("dst.jpg")
//	if err != nil {
//		fmt.Println(err)
//	}
//	defer file.Close()
//
//	file1, err := os.Open("20.jpg")
//	if err != nil {
//		fmt.Println(err)
//	}
//	defer file1.Close()
//	img, _ := jpeg.Decode(file1)
//
//	jpg := image.NewRGBA(image.Rect(0, 0, 100, 100))
//	draw.Draw(jpg, img.Bounds().Add(image.Pt(10, 10)), img, img.Bounds().Min, draw.Src) //截取图片的一部分
//
//	jpeg.Encode(file, jpg, nil)
//
//}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var (
	dpi      = flag.Float64("dpi", 144, "screen resolution in Dots Per Inch")
	fontfile = flag.String("fontfile", "FZHTJW.TTF", "filename of the ttf font")
	hinting  = flag.String("hinting", "none", "none | full")
	size     = flag.Float64("size", 16, "font size in points")
	spacing  = flag.Float64("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")
	wonb     = flag.Bool("whiteonblack", false, "white text on a black background")
)

var text = []string{
	"’Twas brillig, and the slithy toves",
	"测试中文字测试中文字测试中文字测试中文字测试中文字测试中文字测试中文字测试字测",
	"Did gyre and gimble in the wabe;",
	"All mimsy were the borogoves,",
	"And the mome raths outgrabe.",
	"",
	"“Beware the Jabberwock, my son!",
	"The jaws that bite, the claws that catch!",
	"Beware the Jubjub bird, and shun",
	"The frumious Bandersnatch!”",
	"",
	"He took his vorpal sword in hand:",
	"Long time the manxome foe he sought—",
	"So rested he by the Tumtum tree,",
	"And stood awhile in thought.",
	"",
	"And as in uffish thought he stood,",
	"The Jabberwock, with eyes of flame,",
	"Came whiffling through the tulgey wood,",
	"And burbled as it came!",
	"",
	"One, two! One, two! and through and through",
	"The vorpal blade went snicker-snack!",
	"He left it dead, and with its head",
	"He went galumphing back.",
	"",
	"“And hast thou slain the Jabberwock?",
	"Come to my arms, my beamish boy!",
	"O frabjous day! Callooh! Callay!”",
	"He chortled in his joy.",
	"",
	"’Twas brillig, and the slithy toves",
	"Did gyre and gimble in the wabe;",
	"All mimsy were the borogoves,",
	"And the mome raths outgrabe.",
}

type line interface {
	draw(width, x, y int, image draw.Image, cc image.Image)
}

type textLine struct {
	Fontsize int
	Line     string
	font     *truetype.Font
}

func (this *textLine) draw(width, x, y int, image draw.Image, cc image.Image) error {
	var err error
	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetFont(this.font)
	c.SetFontSize(*size)
	c.SetClip(image.Bounds())
	c.SetDst(image)
	c.SetSrc(cc)
	switch *hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	// Draw the guidelines.
	//for i := 0; i < 200; i++ {
	//	image.Set(10, 10+i, ruler)
	//	image.Set(10+i, 10, ruler)
	//}

	// Draw the text.
	pt := freetype.Pt(10, 10+int(c.PointToFixed(*size)>>6))
	for _, s := range text {
		_, err = c.DrawString(s, pt)
		if err != nil {
			log.Println(err)
			return nil
		}
		pt.Y += c.PointToFixed(*size * *spacing)
	}
	return nil
}

type pictureLine struct {
	Fontsize int
	FontType string
}

func (this *pictureLine) draw(width, x, y int, image draw.Image, cc image.Image) error {
	return nil
}

type Configure struct {
	Width  int
	Height int
	bg     image.Image
}

func NewTextPicture(conf Configure) *TextPicture {
	pic := new(TextPicture)
	pic.conf = conf
	pic.rgba = image.NewRGBA(image.Rect(0, 0, conf.Width, conf.Height))
	draw.Draw(pic.rgba, pic.rgba.Bounds(), conf.bg, image.ZP, draw.Src)
	return pic
}

type TextPicture struct {
	text  string
	conf  Configure
	rgba  *image.RGBA
	lines []line
}

func (this *TextPicture) AddLine(l line) {
	this.lines = append(this.lines, l)
}

func main() {
	flag.Parse()

	// Read the font data.
	fontBytes, err := ioutil.ReadFile(*fontfile)
	if err != nil {
		log.Println(err)
		return
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	// Initialize the context.
	fg, bg := image.Black, image.White
	//ruler := color.RGBA{0xdd, 0xdd, 0xdd, 0xff}
	if *wonb {
		fg, bg = image.White, image.Black
		//ruler = color.RGBA{0xFF, 0x00, 0x00, 0xff}
	}
	rgba := image.NewRGBA(image.Rect(0, 0, 1280, 960))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetFont(f)
	c.SetFontSize(*size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	fmt.Println(fg)
	c.SetSrc(image.NewUniform(color.RGBA{0xFF, 0x00, 0x00, 0xff}))
	switch *hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	// Draw the guidelines.
	//for i := 0; i < 200; i++ {
	//	rgba.Set(10, 10+i, ruler)
	//	rgba.Set(10+i, 10, ruler)
	//}

	// Draw the text.
	pt := freetype.Pt(10, 10+int(c.PointToFixed(*size)>>6))
	for _, s := range text {
		_, err = c.DrawString(s, pt)
		if err != nil {
			log.Println(err)
			return
		}
		pt.Y += c.PointToFixed(*size * *spacing)
	}

	// Save that RGBA image to disk.
	outFile, err := os.Create("out.png")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("Wrote out.png OK.")
}
