/* Date:10/4/23-2023
 * Author:fanzhao
 */
package http

import (
	"context"
	"simple-comment-sys/model"
)

func addComment(c context.Context) {
	var (
		r = new(model.AddCommentReq)
	)
	if err := c.BindWith(r, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if mid, ok := c.Get("mid"); ok {
		r.Mid = mid.(int64)
	}
	c.JSON(nil, svc.AddComment(c, r))
}

func delComment(c context.Context) {
	var (
		r = new(model.UpdateCommentReq)
	)
	if err := c.BindWith(r, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if mid, ok := c.Get("mid"); ok {
		r.Mid = mid.(int64)
	}
	c.JSON(nil, svc.DelComment(c, r))
}

func setCommentState(c context.Context) {
	var (
		r = new(model.UpdateCommentReq)
	)
	if err := c.BindWith(r, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svc.SetCommentState(c, r))
}

func listComments(c context.Context) {
	var (
		r = new(model.ListCommentsReq)
	)
	if err := c.BindWith(r, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svc.ListComments(c, r))
}

func exportComments(c context.Context) {
	var (
		r = new(export.CommentExportMsg)
	)
	if err := c.BindWith(r, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if mid, ok := c.Get("mid"); ok {
		r.Mid = mid.(int64)
	}
	c.JSON(nil, svc.ExportComments(c, r))
}

func getExportCommentsFile(c context.Context) {
	var (
		r = new(model.GetExportCommentReq)
	)
	if err := c.BindWith(r, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svc.GetExportCommentsFile(c, r))
}
