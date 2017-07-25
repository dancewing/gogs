package board

import (
	"net/http"

	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/pkg/context"
	"github.com/gogits/gogs/routes/api/board/form"
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

// ItemBoard gets a specific board, identified by project ID or
// NAMESPACE/BOARD_NAME, which is owned by the authenticated user.
func ItemBoard(ctx *context.APIContext) {

	//TODO check privileges

	err := ctx.Repo.Repository.LoadAttributes()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &form.ResponseError{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	board := form.MapBoardFromGogs(ctx.Repo.Repository)

	ctx.JSON(200, &form.Response{
		Data: board,
	})

}

func GetItemBoardByID(BoardID int64) (*form.Board, error) {
	repo, err := models.GetRepositoryByID(BoardID)

	if err != nil {
		return nil, err
	}
	board := form.MapBoardFromGogs(repo)
	return board, nil
}

// Configure configure gitlab repository for usage board
func Configure(c *context.APIContext, f form.BoardRequest) {

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

	apiLabels := make([]*form.Label, len(labels))
	for i := range labels {
		apiLabels[i] = form.MapLabelFromGogs(labels[i])
	}
	c.JSON(200, &form.Response{
		Data: &apiLabels,
	})

}

// CreateConnectBoard add other repository to current board
func CreateConnectBoard(ctx *context.APIContext, f form.BoardRequest) {

	status, err := CreateConnectBoardToCache(ctx, ctx.ParamsInt64(":board"), f.BoardId)

	//status, err := ctx.DataSource.CreateConnectBoard(ctx.Params(":board"), form.BoardId)
	//
	if err != nil {
		ctx.JSON(status, &form.ResponseError{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &form.Response{})
}

// ListConnectBoard gets list connected boards
func ListConnectBoard(ctx *context.APIContext) {

	boards, status, err := ListConnectBoardFromCache(ctx, ctx.Repo.Repository.ID)

	//boards, status, err := ctx.DataSource.ListConnectBoard(ctx.Params(":board"))
	//
	if err != nil {
		ctx.JSON(status, &form.ResponseError{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &form.Response{
		Data: boards,
	})
}

// DeleteConnectBoard deletes board from connected board
func DeleteConnectBoard(ctx *context.APIContext) {

	status, err := DeleteConnectBoardFromCache(ctx, ctx.Params(":board"), ctx.Query("board_id"))

	//status, err := ctx.DataSource.DeleteConnectBoard(ctx.Params(":board"), ctx.Query("board_id"))
	//
	if err != nil {
		ctx.JSON(status, &form.ResponseError{
			Success: false,
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, &form.Response{})
}
