package main

import (
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"time"
)

func main() {
	width := 300
	height := 300

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	r :=rand.New(rand.NewSource(time.Now().UnixNano()))
	

	var col = [3]color.RGBA{
		{255,255,255,0xff},
		{134,27,227,0xff},
		{68,235,153,0xff},
	}

	for x := 0; x < 5; x++ {
		for y := 0; y < 10; y++ {
			switch {
			case x == 0 || y == 0 || y == 9:
				paintReq(img,x*30,y*30,col[0])
				paintReq(img,270-x*30,y*30,col[0])
			default:
				c := col[r.Intn(3)]
				paintReq(img,x*30,y*30,c)
				paintReq(img,270-x*30,y*30,c)
			}
		}
	}
	
	f, _ := os.Create("image.png")
	png.Encode(f, img)
}

func paintReq(img *image.RGBA,x,y int, col color.Color) {
	for i:=x;i<x+30;i++{
		for j:=y;j<y+30;j++{
			img.Set(i,j,col)
		}
	}
}