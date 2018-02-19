package main

import (
	"github.com/yuin/gopher-lua"
	"strconv"
	"golang.org/x/image/colornames"
	"fmt"
	"github.com/faiface/pixel"
)

func mainLoop() {

	initiate()

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
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", windowTitlePrefix, frames))
			frames = 0
		default:
		}

	}
}
