// Copyright 2019 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

// +build !go1.12

package loggo

// logDefaultDepth holds a constant value for the default depth that the
// runtime.CallerFrames is called with. This has two values depending on the
// runtime version you're using; 1.12 and greater uses 1, because of a fix
// to a golang bug[1], where the stacktrace showed the wrong function location.
// Anything less than 1.12, will use 2.
//
// 1. https://golang.org/issue/26839
const logDefaultDepth = 2
