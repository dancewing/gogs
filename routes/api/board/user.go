package board

import (
	"net/http"

	"github.com/gogits/gogs/pkg/context"
	"github.com/gogits/gogs/routes/api/board/form"
)

// ListMembers gets a list of member on board accessible by the authenticated user.
func ListMembers(ctx *context.APIContext) {

	users, err := ctx.Repo.Repository.GetWriters()

	//members, err := ctx.DataSource.ListMembers(ctx.Query("project_id"))
	//
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, &form.ResponseError{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	results := make([]*form.User, len(users))

	for i := range users {
		results[i] = form.MapUserFromGogs(users[i])
	}
	ctx.JSON(200, &form.Response{
		Data: &results,
	})
	//
	//ctx.JSON(http.StatusOK, members)
}


func GetAuthenticatedUser(c *context.Context) {
	c.JSONSuccess(&form.Response{
		Data: c.User.APIFormat(),
	})
}
