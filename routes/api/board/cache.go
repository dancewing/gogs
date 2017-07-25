package board

import (
	"fmt"

	"strings"

	"github.com/gogits/gogs/pkg/context"
	"github.com/gogits/gogs/pkg/tool"
	"github.com/gogits/gogs/routes/api/board/gitlab"
)

// CreateConnectBoard connects other board to current
// for show all cards from other boards
func CreateConnectBoardToCache(ctx *context.APIContext, BoardID, ConnectBoardID int64) (int, error) {

	current, err := GetItemBoardByID(BoardID)

	if err != nil {
		return 0, err
	}

	con, err := GetItemBoardByID(ConnectBoardID)

	if err != nil {
		return 0, err
	}

	//_, err = ds.db.LPush(fmt.Sprintf("boards:%d:connect", current.ID), fmt.Sprintf("%d", con.ID)).Result()
	key := fmt.Sprintf("boards:%d:connect", current.Id)
	val := fmt.Sprintf("%d", con.Id)
	var vids []string
	if ctx.Cache.IsExist(key) {
		vids = strings.Split(ctx.Cache.Get(key).(string), "_")

		if notExist(vids, val) {
			vids = append(vids, val)
		}

		err = ctx.Cache.Put(key, strings.Join(vids, "_"), 18000)

		if err != nil {
			return 0, err
		}
	}

	return 0, nil
}
func notExist(arr []string, tar string) bool {

	for i := range arr {
		if arr[i] == tar {
			return true
		}
	}
	return false

}

// ListConnectBoard return list connect board for current board
func ListConnectBoardFromCache(ctx *context.APIContext, boardID int64) ([]*gitlab.Board, int, error) {
	b := []*gitlab.Board{}

	key := fmt.Sprintf("boards:%s:connect", boardID)

	if ctx.Cache.IsExist(key) {
		boards := strings.Split(ctx.Cache.Get(key).(string), "_")
		boardids := tool.StringsToInt64s(boards)

		for _, board := range boardids {
			item, _ := GetItemBoardByID(board)
			b = append(b, item)
		}
	}

	return b, 0, nil
}

// DeleteConnectBoard deletes from connected board list board
func DeleteConnectBoardFromCache(ctx *context.APIContext, boardID, ConnectBoardID string) (int, error) {

	key := fmt.Sprintf("boards:%s:connect", boardID)

	err := ctx.Cache.Delete(key)
	if err != nil {
		return 0, err
	}

	return 0, nil
}
