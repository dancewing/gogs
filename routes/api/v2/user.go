package v2

import "github.com/gogits/gogs/pkg/context"

// ListMembers gets a list of member on board accessible by the authenticated user.
func ListMembers(ctx *context.APIContext) {
	//members, err := ctx.DataSource.ListMembers(ctx.Query("project_id"))
	//
	//if err != nil {
	//	ctx.JSON(http.StatusUnauthorized, &gitlab.ResponseError{
	//		Success: false,
	//		Message: err.Error(),
	//	})
	//	return
	//}
	//
	//ctx.JSON(http.StatusOK, members)
}
