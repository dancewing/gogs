package v2

import (
	"net/http"

	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/pkg/context"
	"github.com/gogits/gogs/routes/api/v2/gitlab"
	"github.com/gogits/gogs/routes/api/v2/repo"
)

var (
	defaultKBStages = []string{
		"KB[stage][10][Backlog]",
		"KB[stage][20][Development]",
		"KB[stage][30][Testing]",
		"KB[stage][40][Production]",
		"KB[stage][50][Ready]",
	}
)

// ListBoards gets a list of board accessible by the authenticated user.
func ListBoards(ctx *context.APIContext) {
	repo.ListMyRepos(ctx)
}

// ListStarredBoards gets a list starred boards
func ListStarredBoards(ctx *context.APIContext) {
	repo.ListMyRepos(ctx)
}

// ItemBoard gets a specific board, identified by project ID or
// NAMESPACE/BOARD_NAME, which is owned by the authenticated user.
func ItemBoard(ctx *context.APIContext) {

	//TODO check privileges

	fullPath := ctx.Query("project_id")

	repo, err := models.GetRepositoryByRef(fullPath)

	//repo, err := models.GetRepositoryByID(ctx.QueryInt64("project_id"))

	if err != nil {
		if err, ok := err.(gitlab.ReceivedDataErr); ok {
			ctx.JSON(err.StatusCode, &gitlab.ResponseError{
				Success: false,
				Message: err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, &gitlab.ResponseError{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	board := gitlab.MapBoardFromGogs(repo)

	ctx.JSON(200, &gitlab.Response{
		Data: board,
	})

}

func GetItemBoardByID(BoardID int64) (*gitlab.Board, error) {
	repo, err := models.GetRepositoryByID(BoardID)

	if err != nil {
		return nil, err
	}
	board := gitlab.MapBoardFromGogs(repo)
	return board, nil
}

// Configure configure gitlab repository for usage board
func Configure(c *context.APIContext, form gitlab.BoardRequest) {

	//if !c.Repo.IsWriter() {
	//	c.Status(403)
	//	return
	//}

	labels := make([]*models.Label, len(defaultKBStages))

	for i := range defaultKBStages {
		labels[i] = &models.Label{
			Name:   defaultKBStages[i],
			Color:  "#ee0701",
			RepoID: c.Repo.Repository.ID,
		}
	}

	for i := range labels {
		if err := models.NewLabels(labels[i]); err != nil {
			c.Error(500, "NewLabel", err)
			return
		}
	}

	apiLabels := make([]*gitlab.Label, len(labels))
	for i := range labels {
		apiLabels[i] = gitlab.MapLabelFromGogs(labels[i])
	}
	c.JSON(200, &gitlab.Response{
		Data: &apiLabels,
	})

}

// CreateConnectBoard add other repository to current board
func CreateConnectBoard(ctx *context.APIContext, form gitlab.BoardRequest) {

	status, err := CreateConnectBoardToCache(ctx, ctx.ParamsInt64(":board"), form.BoardId)

	//status, err := ctx.DataSource.CreateConnectBoard(ctx.Params(":board"), form.BoardId)
	//
	if err != nil {
		ctx.JSON(status, &gitlab.ResponseError{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &gitlab.Response{})
}

// ListConnectBoard gets list connected boards
func ListConnectBoard(ctx *context.APIContext) {

	boards, status, err := ListConnectBoardFromCache(ctx, ctx.ParamsInt64(":board"))

	//boards, status, err := ctx.DataSource.ListConnectBoard(ctx.Params(":board"))
	//
	if err != nil {
		ctx.JSON(status, &gitlab.ResponseError{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &gitlab.Response{
		Data: boards,
	})
}

// DeleteConnectBoard deletes board from connected board
func DeleteConnectBoard(ctx *context.APIContext) {

	status, err := DeleteConnectBoardFromCache(ctx, ctx.Params(":board"), ctx.Query("board_id"))

	//status, err := ctx.DataSource.DeleteConnectBoard(ctx.Params(":board"), ctx.Query("board_id"))
	//
	if err != nil {
		ctx.JSON(status, &gitlab.ResponseError{
			Success: false,
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, &gitlab.Response{})
}
