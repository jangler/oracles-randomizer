package rom

//go:generate esc -o embed.go -pkg rom -prefix .. ../lgbtasm/lgbtasm.lua

import (
	"github.com/yuin/gopher-lua"
)

// wraps a lua state used for converting gb assembly code to machine code.
type assembler struct {
	ls      *lua.LState
	lgbtasm *lua.LTable
}

// returns a new assembler object, or an error if the source lua code cannot be
// loaded.
func newAssembler() (*assembler, error) {
	ls := lua.NewState()

	mod, err := ls.LoadString(FSMustString(false, "/lgbtasm/lgbtasm.lua"))
	if err != nil {
		return nil, err
	}

	env := ls.Get(lua.EnvironIndex)
	pkg := ls.GetField(env, "package")
	preload := ls.GetField(pkg, "preload")
	ls.SetField(preload, "lgbtasm", mod)
	ls.DoString(`lgbtasm = require "lgbtasm"`)

	return &assembler{
		ls:      ls,
		lgbtasm: ls.GetGlobal("lgbtasm").(*lua.LTable),
	}, nil
}

// compileBlock wraps `lgbtasm.compile_block()`.
func (asm *assembler) compileBlock(s, delim string) (string, error) {
	if err := asm.ls.CallByParam(lua.P{
		Fn:      asm.lgbtasm.RawGetString("compile_block"),
		NRet:    1,
		Protect: true,
	}, lua.LString(s), lua.LString(delim)); err != nil {
		return "", err
	}
	ret := asm.ls.Get(-1)
	asm.ls.Pop(1)

	return ret.(lua.LString).String(), nil
}
