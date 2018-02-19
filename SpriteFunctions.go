package main

import "github.com/yuin/gopher-lua"

var (
	x = float64(screenWidth / 2)
	y = float64(screenHeight / 2)
)

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