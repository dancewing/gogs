package board

import (
	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/pkg/context"
	"github.com/gogits/gogs/routes/api/board/form"
)

// ListComments gets a list of comment on board and card
// accessible by the authenticated user.
func ListComments(c *context.APIContext) {

	// comments,err:=models.GetCommentsByIssueIDSince(, since)
	issue, err := models.GetRawIssueByIndex(c.Repo.Repository.ID, c.ParamsInt64(":index"))
	if err != nil {
		c.Error(500, "GetRawIssueByIndex", err)
		return
	}

	comments, err := models.GetCommentsByIssueID(issue.ID)

	if err != nil {
		c.Error(500, "GetCommentsByIssueID", err)
		return
	}

	apiComments := make([]*form.Comment, len(comments))
	for i := range comments {
		apiComments[i] = form.MapCommentFromGogs(comments[i])
	}
	c.JSON(200, &form.Response{
		Data: apiComments,
	})
}

// CreateComment creates new kanban comment
func CreateComment(ctx *context.APIContext, form form.CommentRequest) {
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
