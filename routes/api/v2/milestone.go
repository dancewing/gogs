package v2

import (
	"github.com/gogits/gogs/pkg/context"
	"github.com/gogits/gogs/routes/api/v2/gitlab"
)

// ListMilestones gets a list of milestone on board accessible by the authenticated user.
func ListMilestones(ctx *context.APIContext) {
	//labels, err := ctx.DataSource.ListMilestones(ctx.Query("project_id"))
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
	//	Data: labels,
	//})
}

// CreateMilestone creates a new board milestone.
func CreateMilestone(ctx *context.APIContext, form gitlab.MilestoneRequest) {
	//milestone, code, err := ctx.DataSource.CreateMilestone(&form)
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
	//	Data: milestone,
	//})
}
