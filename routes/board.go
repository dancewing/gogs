// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package routes

import (
	"github.com/gogits/gogs/pkg/context"
	"github.com/gogits/gogs/pkg/setting"
)

const (
	BOARD = "board"
)

func Board(c *context.Context) {

	if !setting.ProdMode {

	}
	c.HTML(200, BOARD)
}
