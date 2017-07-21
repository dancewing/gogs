package v2

import (
	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/pkg/context"
	"github.com/gogits/gogs/routes/api/v2/gitlab"
)

// ListLabels gets a list of label on board accessible by the authenticated user.
func ListLabels(ctx *context.APIContext) {

	labels, err := models.GetLabelsByRepoID(ctx.ParamsInt64(":project"))
	if err != nil {
		ctx.Error(500, "GetLabelsByRepoID", err)
		return
	}

	apiLabels := make([]*gitlab.Label, len(labels))
	for i := range labels {
		apiLabels[i] = gitlab.MapLabelFromGogs(labels[i])
	}
	ctx.JSON(200, &gitlab.Response{
		Data: &apiLabels,
	})
}

// EditLabel updates existing project label
func EditLabel(ctx *context.APIContext, form gitlab.LabelRequest) {
	//log.Printf("GOT LABEL req %+v", form)
	//label, err := ctx.DataSource.EditLabel(ctx.Params(":project"), &form)
	//if err != nil {
	//	ctx.JSON(http.StatusUnauthorized, &gitlab.ResponseError{
	//		Success: false,
	//		Message: err.Error(),
	//	})
	//}
	//
	//ctx.JSON(http.StatusOK, &gitlab.Response{
	//	Data: label,
	//})
}

// CreateLabel creates new label
func CreateLabel(ctx *context.APIContext, form gitlab.LabelRequest) {
	//label, err := ctx.DataSource.CreateLabel(ctx.Params(":project"), &form)
	//
	//if err != nil {
	//	ctx.JSON(http.StatusUnauthorized, &gitlab.ResponseError{
	//		Success: false,
	//		Message: err.Error(),
	//	})
	//}
	//
	//ctx.JSON(http.StatusOK, &gitlab.Response{Data: label})
}

// DeleteLabel removes existing project label
func DeleteLabel(ctx *context.APIContext) {
	//label, err := ctx.DataSource.DeleteLabel(ctx.Params(":project"), ctx.Params(":label"))
	//
	//if err != nil {
	//	ctx.JSON(http.StatusUnauthorized, &gitlab.ResponseError{
	//		Success: false,
	//		Message: err.Error(),
	//	})
	//}
	//
	//ctx.JSON(http.StatusOK, &gitlab.Response{
	//	Data: label,
	//})
}
