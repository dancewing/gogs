package v2

import (
	"net/http"

	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/pkg/context"
	"github.com/gogits/gogs/routes/api/v2/gitlab"
)

// ListMilestones gets a list of milestone on board accessible by the authenticated user.
func ListMilestones(ctx *context.APIContext) {

	milestones, err := models.GetMilestonesByRepoID(ctx.Repo.Repository.ID)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, &gitlab.ResponseError{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	results := make([]*gitlab.Milestone, len(milestones))

	for i := range milestones {
		results[i] = gitlab.MapMilestoneFromGogs(milestones[i])
	}
	ctx.JSON(200, &gitlab.Response{
		Data: &results,
	})
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
