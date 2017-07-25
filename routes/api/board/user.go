package board

import (
	"net/http"

	"github.com/gogits/gogs/pkg/context"
	"github.com/gogits/gogs/routes/api/board/gitlab"
)

// ListMembers gets a list of member on board accessible by the authenticated user.
func ListMembers(ctx *context.APIContext) {

	users, err := ctx.Repo.Repository.GetWriters()

	//members, err := ctx.DataSource.ListMembers(ctx.Query("project_id"))
	//
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, &gitlab.ResponseError{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	results := make([]*gitlab.User, len(users))

	for i := range users {
		results[i] = gitlab.MapUserFromGogs(users[i])
	}
	ctx.JSON(200, &gitlab.Response{
		Data: &results,
	})
	//
	//ctx.JSON(http.StatusOK, members)
}
