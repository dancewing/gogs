package board

import (
	"net/http"

	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/pkg/context"
	"github.com/gogits/gogs/routes/api/board/form"
)

// ListMilestones gets a list of milestone on board accessible by the authenticated user.
func ListMilestones(ctx *context.APIContext) {

	milestones, err := models.GetMilestonesByRepoID(ctx.Repo.Repository.ID)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, &form.ResponseError{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	results := make([]*form.Milestone, len(milestones))

	for i := range milestones {
		results[i] = form.MapMilestoneFromGogs(milestones[i])
	}
	ctx.JSON(200, &form.Response{
		Data: &results,
	})
}

// CreateMilestone creates a new board milestone.
func CreateMilestone(ctx *context.APIContext, form form.MilestoneRequest) {
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
