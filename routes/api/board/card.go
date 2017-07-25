package board

import (
	"net/http"

	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/pkg/context"
	"github.com/gogits/gogs/routes/api/board/form"

	"github.com/gogits/gogs/models/errors"
	log "gopkg.in/clog.v1"
)

// ListCards gets a list of card on board accessible by the authenticated user.
func ListCards(ctx *context.APIContext) {
	//pr := ctx.QueryInt64("project_id")

	//proj, _, err := ctx.DataSource.ListConnectBoard(pr)

	proj, _, err := ListConnectBoardFromCache(ctx, ctx.Repo.Repository.ID)

	if err != nil {
		ctx.JSON(http.StatusNotFound, &form.ResponseError{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	//current, err := ctx.DataSource.ItemBoard(pr)
	current, err := GetItemBoardByID(ctx.Repo.Repository.ID)

	if err != nil {
		ctx.JSON(http.StatusNotFound, &form.ResponseError{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	proj = append(proj, current)

	var cards, c []*form.Card

	for _, p := range proj {
		c, err = listCards(p.Id)
		cards = append(cards, c...)
	}

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, &form.ResponseError{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &form.Response{
		Data: cards,
	})
}

func listCards(projectId int64) ([]*form.Card, error) {
	issues, err := models.Issues(&models.IssuesOptions{
		RepoID: projectId,
	})
	apiIssues := make([]*form.Card, len(issues))
	for i := range issues {
		if err = issues[i].LoadAttributes(); err != nil {
			// c.Error(500, "LoadAttributes", err)
			//return
		}
		apiIssues[i] = form.MapCardFromGogs(issues[i])
	}

	return apiIssues, nil
}

// CreateCard creates a new board card.
func CreateCard(ctx *context.APIContext, form form.CardRequest) {
	//card, code, err := ctx.DataSource.CreateCard(&form)
	//
	//if err != nil {
	//	ctx.JSON(code, &gitlab.ResponseError{
	//		Success: false,
	//		Message: err.Error(),
	//	})
	//	return
	//}
	//
	//card.BoardID = ctx.Params(":board")
	//
	//ctx.JSON(http.StatusOK, &gitlab.Response{
	//	Data: card,
	//})
	//
	//ctx.Broadcast(card.RoutingKey(), &gitlab.Response{
	//	Data:  card,
	//	Event: "card.create",
	//})
}

// UpdateCard updates an existing board card.
func UpdateCard(ctx *context.APIContext, form form.CardRequest) {
	//card, code, err := ctx.DataSource.UpdateCard(&form)
	//
	//if err != nil {
	//	ctx.JSON(code, &gitlab.ResponseError{
	//		Success: false,
	//		Message: err.Error(),
	//	})
	//	return
	//}
	//
	//card.BoardID = ctx.Params(":board")
	//
	//ctx.JSON(http.StatusOK, &gitlab.Response{
	//	Data: card,
	//})
	//
	//ctx.Broadcast(card.RoutingKey(), &gitlab.Response{
	//	Data:  card,
	//	Event: "card.update",
	//})
}

// DeleteCard closed an existing board card.
func DeleteCard(ctx *context.APIContext, form form.CardRequest) {
	//card, code, err := ctx.DataSource.DeleteCard(&form)
	//
	//if err != nil {
	//	ctx.JSON(code, &gitlab.ResponseError{
	//		Success: false,
	//		Message: err.Error(),
	//	})
	//	return
	//}
	//card.BoardID = ctx.Params(":board")
	//
	//ctx.JSON(http.StatusOK, &gitlab.Response{
	//	Data: card,
	//})
	//
	//ctx.Broadcast(card.RoutingKey(), &gitlab.Response{
	//	Data:  card,
	//	Event: "card.delete",
	//})
}

// MoveToCard updates an existing board card.
func MoveToCard(ctx *context.APIContext, f form.CardRequest) {

	issue, err := models.GetIssueByID(f.CardId)
	if err != nil {
		if errors.IsIssueNotExist(err) {
			ctx.Status(404)
		} else {
			ctx.Error(500, "GetIssueByID", err)
		}
		return
	}

	if !issue.IsPoster(ctx.User.ID) && !ctx.Repo.IsWriter() {
		ctx.Status(403)
		return
	}

	source := f.Stage["source"]
	dest := f.Stage["dest"]

	//change label
	if dest != source {
		sl, err := models.GetLabelOfRepoByName(ctx.Repo.Repository.ID, source)

		if err != nil {

			ctx.Error(500, "GetLabelOfRepoByName", err)
			return
		}

		dl, err := models.GetLabelOfRepoByName(ctx.Repo.Repository.ID, dest)

		if err != nil {

			ctx.Error(500, "GetLabelOfRepoByName", err)
			return
		}

		log.Info("try to change %s 's status from %s to %s", issue.Title, sl.Name, dl.Name)

		issue.RemoveLabel(ctx.User, sl)

		issue.AddLabel(ctx.User, dl)
	}

	if issue.MilestoneID != f.MilestoneId {
		milestore , err := models.GetMilestoneByRepoID(ctx.Repo.Repository.ID, f.MilestoneId)

		if err != nil {
			ctx.Error(500, "GetMilestoneByRepoID", err)
			return
		}

		issue.MilestoneID =milestore.ID
		err = models.UpdateIssue(issue)

		if err != nil {
			ctx.Error(500, "UpdateIssue", err)
			return
		}
	}

	if issue.AssigneeID != f.AssigneeId {
		assignee , err := models.GetAssigneeByID(ctx.Repo.Repository, f.AssigneeId)
		if err != nil {
			ctx.Error(500, "GetAssigneeByID", err)
			return
		}

		issue.AssigneeID = assignee.ID

		err = models.UpdateIssue(issue)
		if err != nil {
			ctx.Error(500, "UpdateIssue", err)
			return
		}
	}

	err = issue.LoadAttributes()

	if err != nil {
		ctx.Error(500, "issue.LoadAttributes", err)
		return
	}

	card:= form.MapCardFromGogs(issue)

	ctx.JSON(http.StatusOK, &form.Response{
		Data: card,
	})

	ctx.Broadcast(card.RoutingKey(), &form.Response{
		Data:  card,
		Event: "card.move",
	})

	//if source.Name != dest.Name && viper.GetBool("auto.comments") {
	//	com := models.CommentRequest{
	//		CardId:    form.CardId,
	//		ProjectId: form.ProjectId,
	//		Body:      fmt.Sprintf("moved issue from **%s** to **%s**", source.Name, dest.Name),
	//	}
	//
	//	go func() {
	//		ctx.DataSource.CreateComment(&com)
	//	}()
	//}
}

// ChangeProjectForCard locate card to another project
func ChangeProjectForCard(ctx *context.APIContext, form form.CardRequest) {
	//oldCrard := gitlab.Card{
	//	Id:        form.CardId,
	//	ProjectId: form.ProjectId,
	//	BoardID:   ctx.Params(":board"),
	//}
	//
	//card, code, err := ctx.DataSource.ChangeProjectForCard(&form, ctx.Params(":projectId"))
	//
	//if err != nil {
	//	ctx.JSON(code, &gitlab.ResponseError{
	//		Success: false,
	//		Message: err.Error(),
	//	})
	//	return
	//}
	//
	//ctx.Broadcast(oldCrard.RoutingKey(), &gitlab.Response{
	//	Data:  oldCrard,
	//	Event: "card.delete",
	//})
	//
	//card.BoardID = ctx.Params(":board")
	//
	//ctx.Broadcast(card.RoutingKey(), &gitlab.Response{
	//	Data:  card,
	//	Event: "card.create",
	//})
	//
	//ctx.JSON(http.StatusOK, &gitlab.Response{
	//	Data: card,
	//})
}
