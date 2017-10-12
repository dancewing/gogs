// Copyright 2016 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package context

import (
	"fmt"
	"strings"

	"github.com/Unknwon/paginater"
	log "gopkg.in/clog.v1"
	"gopkg.in/macaron.v1"

	"github.com/gogits/gogs/pkg/setting"
)

type RestContext struct {
	*Context
}
// Error responses error message to client with given message.
// If status is 500, also it prints error to log.
func (c *RestContext) Error(status int, title string, obj interface{}) {
	var message string
	if err, ok := obj.(error); ok {
		message = err.Error()
	} else {
		message = obj.(string)
	}

	if status == 500 {
		log.Error(3, "%s: %s", title, message)
	}

	c.JSON(status, map[string]string{
		"message": message,
		"url":     DOC_URL,
	})
}

// SetLinkHeader sets pagination link header by given totol number and page size.
func (c *RestContext) SetLinkHeader(total, pageSize int) {
	page := paginater.New(total, pageSize, c.QueryInt("page"), 0)
	links := make([]string, 0, 4)
	if page.HasNext() {
		links = append(links, fmt.Sprintf("<%s%s?page=%d>; rel=\"next\"", setting.AppURL, c.Req.URL.Path[1:], page.Next()))
	}
	//if !page.IsLast() {
	//	links = append(links, fmt.Sprintf("<%s%s?page=%d>; rel=\"last\"", setting.AppURL, c.Req.URL.Path[1:], page.TotalPages()))
	//}
	//if !page.IsFirst() {
	//	links = append(links, fmt.Sprintf("<%s%s?page=1>; rel=\"first\"", setting.AppURL, c.Req.URL.Path[1:]))
	//}

			links = append(links, fmt.Sprintf("<%s%s?page=%d>; rel=\"last\"", setting.AppURL, c.Req.URL.Path[1:], page.TotalPages()))

			links = append(links, fmt.Sprintf("<%s%s?page=1>; rel=\"first\"", setting.AppURL, c.Req.URL.Path[1:]))

	if page.HasPrevious() {
		links = append(links, fmt.Sprintf("<%s%s?page=%d>; rel=\"prev\"", setting.AppURL, c.Req.URL.Path[1:], page.Previous()))
	}

	if len(links) > 0 {
		c.Header().Set("Link", strings.Join(links, ","))
	}
	c.Header().Set("X-Total-Count", fmt.Sprintf("%d", total))
}

func RestContexter() macaron.Handler {
	return func(ctx *Context) {
		c := &RestContext{
			Context: ctx,
		}
		ctx.Map(c)
	}
}
