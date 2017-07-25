package board

import (
	"github.com/gogits/gogs/pkg/context"
	"github.com/gogits/gogs/routes/api/board/gitlab"
)

// UploadFile uploads file to datasource provider
func UploadFile(ctx *context.APIContext, f gitlab.UploadForm) {
	//res, err := ctx.DataSource.UploadFile(ctx.Params(":board"), f)
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
	//	Data: res,
	//})
}
