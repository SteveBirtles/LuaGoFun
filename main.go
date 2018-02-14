package main

import (
	"os"
	"fmt"
	"time"
	"image"
	_ "image/png"
	"github.com/yuin/gopher-lua"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
	"strconv"
)

var L *lua.LState

// -----------------------------------------------------------------------------------

const screenWidth = 1024
const screenHeight = 768

var x = float64(screenWidth / 2)
var y = float64(screenHeight / 2)

func up(L *lua.LState) int {
	dy := L.ToNumber(1)
	y += float64(dy)
	L.Push(lua.LNumber(y))
	return 1
}

func down(L *lua.LState) int {
	dy := L.ToNumber(1)
	y -= float64(dy)
	L.Push(lua.LNumber(y))
	return 1
}

func left(L *lua.LState) int {
	dx := L.ToNumber(1)
	x -= float64(dx)
	L.Push(lua.LNumber(x))
	return 1
}

func right(L *lua.LState) int {
	dx := L.ToNumber(1)
	x += float64(dx)
	L.Push(lua.LNumber(x))
	return 1
}

// -----------------------------------------------------------------------------------

func linkToLua(luaState *lua.LState, goFunction lua.LGFunction, goFunctionName string) {
	luaState.SetGlobal(goFunctionName, luaState.NewFunction(goFunction))
}

func executeLua(luaState *lua.LState, luaCode string) {
	err := luaState.DoString(luaCode)
	if err != nil {
		fmt.Println("Lua error: " + err.Error())
	}
}

func executeLuaFile(luaState *lua.LState, luaFile string) {
	err := luaState.DoFile(luaFile)
	if err != nil {
		fmt.Println("Lua error: " + err.Error())
	}
}

func executeLuaFunction(luaState *lua.LState, functionName string, functionArgs []lua.LValue) lua.LValue {

	err := luaState.CallByParam(lua.P{
		Fn:      luaState.GetGlobal(functionName),
		NRet:    1,
		Protect: true,
	}, functionArgs...)

	if err != nil {
		fmt.Println("Lua error: " + err.Error())
		return lua.LNil
	}

	str, ok := luaState.Get(-1).(lua.LValue)

	if ok {
		defer luaState.Pop(1)
		return str
	}

	return lua.LNil
}

// -----------------------------------------------------------------------------------

func loadImageFile(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// -----------------------------------------------------------------------------------

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

// -----------------------------------------------------------------------------------

func main() {

	L = lua.NewState()
	defer L.Close()

	linkToLua(L, up, "up")
	linkToLua(L, down, "down")
	linkToLua(L, left, "left")
	linkToLua(L, right, "right")

	pixelgl.Run(run)

}