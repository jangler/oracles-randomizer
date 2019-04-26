package rom

//go:generate esc -o embed.go -pkg rom -prefix .. ../lgbtasm/lgbtasm.lua ../asm/

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

	asm.compileOpts.RawSet(lua.LString("delims"), lua.LString("\n;"))
	asm.compileOpts.RawSet(lua.LString("defs"), asm.defs)
	asm.decompileOpts.RawSet(lua.LString("defs"), asm.defs)

	return asm, nil
}

// compile wraps `lgbtasm.compile()`.
func (asm *assembler) compile(s, delim string) (string, error) {
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
func (asm *assembler) decompile(s, delim string) (string, error) {
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
