package rom

import (
	"github.com/yuin/gopher-lua"
)

// wraps a lua state used for converting gb assembly code to machine code.
type assembler struct {
	ls            *lua.LState
	lgbtasm       *lua.LTable
	compileOpts   *lua.LTable
	decompileOpts *lua.LTable
	defs          *lua.LTable
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

	asm := &assembler{
		ls:            ls,
		lgbtasm:       ls.GetGlobal("lgbtasm").(*lua.LTable),
		compileOpts:   ls.NewTable(),
		decompileOpts: ls.NewTable(),
		defs:          ls.NewTable(),
	}

	asm.compileOpts.RawSetString("delims", lua.LString("\n;"))
	asm.compileOpts.RawSetString("defs", asm.defs)
	asm.decompileOpts.RawSetString("defs", asm.defs)

	return asm, nil
}

// compile wraps `lgbtasm.compile()`.
func (asm *assembler) compile(s string) (string, error) {
	if err := asm.ls.CallByParam(lua.P{
		Fn:      asm.lgbtasm.RawGetString("compile"),
		NRet:    1,
		Protect: true,
	}, lua.LString(s), asm.compileOpts); err != nil {
		return "", err
	}
	ret := asm.ls.Get(-1)
	asm.ls.Pop(1)

	return ret.(lua.LString).String(), nil
}

// decompile wraps `lgbtasm.decompile()`.
func (asm *assembler) decompile(s string) (string, error) {
	if err := asm.ls.CallByParam(lua.P{
		Fn:      asm.lgbtasm.RawGetString("decompile"),
		NRet:    1,
		Protect: true,
	}, lua.LString(s), asm.decompileOpts); err != nil {
		return "", err
	}
	ret := asm.ls.Get(-1)
	asm.ls.Pop(1)

	return ret.(lua.LString).String(), nil
}

// add a constant def as if `define symbol,string` were run.
func (asm *assembler) define(symbol string, value uint16) {
	asm.defs.RawSetString(symbol, lua.LNumber(value))
}
