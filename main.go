package main

import (
	"github.com/yuin/gopher-lua"
	"fmt"
)

func square(L *lua.LState) int {
	i := L.ToInt(1)
	k := i * i
	L.Push(lua.LNumber(k))
	return 1
}

func add(L *lua.LState) int {
	i := L.ToInt(1)
	j := L.ToInt(2)
	k := i + j
	L.Push(lua.LNumber(k))
	return 1
}

func linkToLua(luaState *lua.LState, goFunction lua.LGFunction, goFunctionName string) {
	luaState.SetGlobal(goFunctionName, luaState.NewFunction(goFunction))
}

func executeLua(luaState *lua.LState, luaCode string) {
	if err := luaState.DoString(luaCode);
	err != nil {
		fmt.Println("Lua error: " + err.Error())
	}
}

func executeLuaFile(luaState *lua.LState, luaFile string) {
	if err := luaState.DoFile(luaFile);
	err != nil {
		fmt.Println("Lua error: " + err.Error())
	}
}

func executeLuaFunction(luaState *lua.LState, functionName string, functionArgs []lua.LValue) lua.LValue {

	if err := luaState.CallByParam(lua.P{
		Fn:      luaState.GetGlobal(functionName),
		NRet:    1,
		Protect: true,
	}, functionArgs...);
	err != nil {
		fmt.Println("Lua error: " + err.Error())
		return lua.LNil
	}

	if str, ok := luaState.Get(-1).(lua.LValue);
	ok {
		defer luaState.Pop(1)
		return str
	}

	return lua.LNil
}


func main() {

	L := lua.NewState()
	defer L.Close()

	executeLua(L, `function sayHello() 
			print("Hello Again") 
		end
		function concat(a, b)
			return a .. " + " .. b
		end`)

	executeLua(L, "sayHello()")

	linkToLua(L, square, "square")
	linkToLua(L, add, "add")

	executeLua(L,`print("4 squared is " .. square(4))`)
	executeLua(L,`print("1 plus 2 is " .. add(1, 2))`)

	executeLuaFile(L, "test.lua")

	result := executeLuaFunction(L, "concat", []lua.LValue{lua.LString("Go"), lua.LString("Lua")})

	fmt.Println(result.String())

}
