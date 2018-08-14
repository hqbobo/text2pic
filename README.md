# text2pic

##Description
convert text into picture


##example

`
package main

import (
	"flag"
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
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	pic := text2pic.NewTextPicture(text2pic.Configure{
		Width: 720,
	})

	pic.AddTextLine("1.这个是标题", 20, f, text2pic.ColorRed, text2pic.Padding{Left: 20, Top: 10, Bottom: 20})
	pic.AddTextLine("    北京铁路局今天凌晨2时16分发布消息称：8月12日23时04分，aaaa京沪高铁廊坊至北京aaaaa南间发生设备故障，导致部分列车晚点。铁路部门及时启动应急预案处置时16分发布消息称时16分发布消息称北京铁路局今天凌晨2时16分发布消息称：8月12日23时04分，aaaa京沪高铁廊坊至北京aaaaa南间发生设备故障，导致部分列车晚点。铁路部门及时启动应急预案处置时16分发布消息称时16分发布消息称北京铁路局今天凌晨2时16分发布消息称：8月12日23时04分，aaaa京沪高铁廊坊至北京aaaaa南间发生设备故障，导致部分列车晚点。铁路部门及时启动应急预案处置时16分发布消息称时16分发布消息称", 12, f, text2pic.ColorGreen, text2pic.Padding{Left: 20, Right: 20, Bottom: 30})

	file, err := os.Open("timg.jpg")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	pic.AddPictureLine(file, text2pic.Padding{Bottom: 20})

	pic.AddTextLine("2.这个是标题", 20, f, text2pic.ColorRed, text2pic.Padding{Bottom: 20})
	pic.AddTextLine("    北京铁路局今天凌晨2时16分发布消息称：8月12日23时04分，aaaa京沪高铁廊坊至北京aaaaa南间发生设备故障，导致部分列车晚点。铁路部门及时启动应急预案处置时16分发布消息称时16分发布消息称北京铁路局今天凌晨2时16分发布消息称：8月12日23时04分，aaaa京沪高铁廊坊至北京aaaaa南间发生设备故障，导致部分列车晚点。铁路部门及时启动应急预案处置时16分发布消息称时16分发布消息称北京铁路局今天凌晨2时16分发布消息称：8月12日23时04分，aaaa京沪高铁廊坊至北京aaaaa南间发生设备故障，导致部分列车晚点。铁路部门及时启动应急预案处置时16分发布消息称时16分发布消息称", 13, f, text2pic.ColorBlue, text2pic.Padding{Bottom: 30})

	file1, err := os.Open("timg1.jpg")
	if err != nil {
		fmt.Println(err)
	}
	defer file1.Close()
	pic.AddPictureLine(file1, text2pic.Padding{Bottom: 20})

	pic.AddTextLine("3.For English", 20, f, text2pic.ColorRed, text2pic.Padding{Bottom: 20})
	pic.AddTextLine(" The Turkish lira plunged as much as 11% against the dollar, hitting a record low, before recovering some of its losses in volatile trading. The lira had already plummeted more than 20% last week as a political clash with the United States intensified and investors fretted about the Turkish government's lack of action to tackle the problems plaguing its economy.  ", 13, f, text2pic.ColorBlue, text2pic.Padding{Left: 20, Right: 20, Bottom: 30})

	// Save that RGBA image to disk.
	outFile, err := os.Create("out.jpg")
	if err != nil {
		return
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)

	fmt.Println(pic.Draw(b, text2pic.TypeJpeg))
	e := b.Flush()
	if e!=nil {
		fmt.Println(e)
	}

}

`