
# text2pic  
  
## Description  
  
>Convert text into pictures, this is designed for posting long text message to weibo initially.  
  
1. paint text on picture.  
2. support jpg and png as well  
  
## Functions  

> Create new picture with Configure
 
`pic := text2pic.NewTextPicture(text2pic.Configure{Width: 720, })`

> Add text line to picture with font, color and padding 

`pic.AddTextLine(" The Turkish lira plunged as much as 11% against the dollar", 13, f, text2pic.ColorBlue, text2pic.Padding{Left: 20, Right: 20, Bottom: 30})`

> Add picture  io.reader is required and padding as well

`pic.AddPictureLine(file, text2pic.Padding{Bottom: 20})`

> Draw it on io.writer. TypePng and TypeJpeg are supported

`pic.Draw(writer, text2pic.TypeJpeg)`


## Example  
> Example see in the example directory

```
package main

import (
	"fmt"
	"github.com/golang/freetype"
	"github.com/hqbobo/text2pic"
	"io/ioutil"
	"log"
	"os"
	"bufio"
)

func main() {

	// Read the font data.
	fontBytes, err := ioutil.ReadFile("FZHTJW.TTF")
	if err != nil {
		log.Println(err)
		return
	}
	//produce the fonttype
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	//define New picture with given width in px
	//the height will be calucated before draw on picture
	//picture will be resize to 80% of the width you given
	pic := text2pic.NewTextPicture(text2pic.Configure{
		Width: 720,
	})

	//add chinese line
	pic.AddTextLine("1.这个是标题", 20, f, text2pic.ColorRed, text2pic.Padding{Left: 20, Top: 10, Bottom: 20})
	pic.AddTextLine("    北京铁路局今天凌晨2时16分发布消息称：8月12日23时04分，aaaa京沪高铁廊坊至北京aaaaa南间发生设备故障，导致部分列车晚点。铁路部门及时启动应急预案处置时16分发布消息称时16分发布消息称北京铁路局今天凌晨2时16分发布消息称：8月12日23时04分，aaaa京沪高铁廊坊至北京aaaaa南间发生设备故障，导致部分列车晚点。铁路部门及时启动应急预案处置时16分发布消息称时16分发布消息称北京铁路局今天凌晨2时16分发布消息称：8月12日23时04分，aaaa京沪高铁廊坊至北京aaaaa南间发生设备故障，导致部分列车晚点。铁路部门及时启动应急预案处置时16分发布消息称时16分发布消息称", 12, f, text2pic.ColorGreen, text2pic.Padding{Left: 20, Right: 20, Bottom: 30})
	//add picture
	file, err := os.Open("timg.jpg")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	
	pic.AddPictureLine(file, text2pic.Padding{Bottom: 20})
	
	//add full english text
	pic.AddTextLine("3.For English", 20, f, text2pic.ColorRed, text2pic.Padding{Bottom: 20})
	pic.AddTextLine(" The Turkish lira plunged as much as 11% against the dollar, hitting a record low, before recovering some of its losses in volatile trading. The lira had already plummeted more than 20% last week as a political clash with the United States intensified and investors fretted about the Turkish government's lack of action to tackle the problems plaguing its economy.  ", 13, f, text2pic.ColorBlue, text2pic.Padding{Left: 20, Right: 20, Bottom: 30})

	// Save the output to file
	outFile, err := os.Create("out.jpg")
	if err != nil {
		return
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	//produce the output
	pic.Draw(b, text2pic.TypeJpeg)
	e := b.Flush()
	if e!=nil {
		fmt.Println(e)
	}

}
```