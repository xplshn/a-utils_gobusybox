// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bbmain is the package imported by all rewritten busybox
// command-packages to register themselves.
package bbmain

import (
	"errors"
	"fmt"
	"os"
	"sort"
	// There MUST NOT be any other dependencies here.
	//
	// It is preferred to copy minimal code necessary into this file, as
	// dependency management for this main file is... hard.
)

// ErrNotRegistered is returned by Run if the given command is not registered.
//var ErrNotRegistered = errors.New("command is not present in busybox")

var ErrNotRegistered = errors.New(`
 Copyright (c) 2024: xplshn, and contributors
 For more details refer to https://github.com/xplshn/a-utils

  Synopsis
    a-utils [program] <args>
  Description:
    a-utils multicall binary
  Notes:
    This binary is a multicall binary.
    It contains various commands inside.
    They can be accessed either by
    symlinking this multicall binary to a file
    with the name of the binary contained here
    that you wish to use or by calling this binary
    with it as the first argument.

`)

// Noop is a noop function.
var Noop = func() {}

// ListCmds returns all supported commands.
func ListCmds() []string {
	var cmds []string
	for c := range bbCmds {
		cmds = append(cmds, c)
	}
	sort.Strings(cmds)
	return cmds
}

type bbCmd struct {
	init, main func()
}

var bbCmds = map[string]bbCmd{}

var defaultCmd *bbCmd

// Register registers an init and main function for name.
func Register(name string, init, main func()) {
	if _, ok := bbCmds[name]; ok {
		panic(fmt.Sprintf("cannot register two commands with name %q", name))
	}
	bbCmds[name] = bbCmd{
		init: init,
		main: main,
	}
}

// RegisterDefault registers a default init and main function.
func RegisterDefault(init, main func()) {
	defaultCmd = &bbCmd{
		init: init,
		main: main,
	}
}

// Run runs the command with the given name.
//
// If the command's main exits without calling os.Exit, Run will exit with exit
// code 0.
func Run(name string) error {
	var cmd *bbCmd
	if c, ok := bbCmds[name]; ok {
		cmd = &c
	} else if defaultCmd != nil {
		cmd = defaultCmd
	} else {
		return fmt.Errorf("%w"+"\x1B[31merror\x1B[m: '%s' is not present in this multicall binary of a-utils", ErrNotRegistered, name)
	}
	cmd.init()
	cmd.main()
	os.Exit(0)
	// Unreachable.
	return nil
}
