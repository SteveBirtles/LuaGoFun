package main

import (
	"time"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

var (
	windowTitlePrefix = "Go Pixel & Lua Test"
	frames            = 0
	second            = time.Tick(time.Second)
	win               *pixelgl.Window
	sprite            *pixel.Sprite
	spriteBatch       *pixel.Batch
	overlay           *imdraw.IMDraw
)

func initiate() {

	var initError error

	cfg := pixelgl.WindowConfig{
		Bounds: pixel.R(0, 0, screenWidth, screenHeight),
		VSync:  true,
	}

	win, initError = pixelgl.NewWindow(cfg)
	if initError != nil {
		panic(initError)
	}

	spriteImage, initError := loadImageFile("sonic.png")
	if initError != nil {
		panic(initError)
	}

	spritePic := pixel.PictureDataFromImage(spriteImage)

	sprite = pixel.NewSprite(spritePic, spritePic.Bounds())

	spriteBatch = pixel.NewBatch(&pixel.TrianglesData{}, spritePic)

	overlay = imdraw.New(nil)

	executeLua(L,
		`angle = 0
				function getAngle() return angle end`)

}