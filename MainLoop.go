package main

import (
	"github.com/yuin/gopher-lua"
	"strconv"
	"time"
	"golang.org/x/image/colornames"
	"fmt"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

func run() {

	cfg := pixelgl.WindowConfig{
		Title:  "Go Pixel & Lua Test",
		Bounds: pixel.R(0, 0, screenWidth, screenHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	spriteImage, err := loadImageFile("sonic.png")
	if err != nil {
		panic(err)
	}

	spritePic := pixel.PictureDataFromImage(spriteImage)

	sprite := pixel.NewSprite(spritePic, spritePic.Bounds())

	spriteBatch := pixel.NewBatch(&pixel.TrianglesData{}, spritePic)
	overlay := imdraw.New(nil)

	executeLua(L,
		`angle = 0
				function getAngle() return angle end`)

	for !win.Closed() {

		spriteBatch.Clear()
		overlay.Clear()

		executeLua(L,
			`angle = angle + 0.05
					up(5 * math.sin(angle))
					right(5 * math.cos(angle))`)

		angleLua := executeLuaFunction(L, "getAngle", []lua.LValue{})

		angle, _ := strconv.ParseFloat(angleLua.String(), 64)

		matrix := pixel.IM.Rotated(pixel.ZV, angle).Scaled(pixel.ZV, 0.2).Moved(pixel.Vec{X: x, Y: y})

		sprite.Draw(spriteBatch, matrix)

		win.Clear(colornames.Black)

		win.SetComposeMethod(pixel.ComposeOver)
		spriteBatch.Draw(win)

		win.SetComposeMethod(pixel.ComposePlus)
		overlay.Draw(win)

		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}
