package v2

import (
	"github.com/gogits/gogs/pkg/context"
	"github.com/gogits/gogs/routes/api/v2/gitlab"
)

// ListComments gets a list of comment on board and card
// accessible by the authenticated user.
func ListComments(ctx *context.APIContext) {
	//boards, err := ctx.DataSource.ListComments(ctx.Query("project_id"), ctx.Query("issue_id"))
	//
	//if err != nil {
	//	ctx.JSON(http.StatusUnauthorized, &gitlab.ResponseError{
	//		Success: false,
	//		Message: err.Error(),
	//	})
	//	return
	//}
	//
	//ctx.JSON(http.StatusOK, &gitlab.Response{
	//	Data: boards,
	//})
}

// CreateComment creates new kanban comment
func CreateComment(ctx *context.APIContext, form gitlab.CommentRequest) {
	//com, code, err := ctx.DataSource.CreateComment(&form)
	//
	//if err != nil {
	//	ctx.JSON(code, &gitlab.ResponseError{
	//		Success: false,
	//		Message: err.Error(),
	//	})
	//	return
	//}
	//
	//ctx.JSON(http.StatusOK, &gitlab.Response{
	//	Data: com,
	//})
}
