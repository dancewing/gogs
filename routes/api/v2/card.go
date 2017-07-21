package v2

import (
	"net/http"

	"github.com/gogits/gogs/pkg/context"
	"github.com/gogits/gogs/routes/api/v2/gitlab"
)

// ListCards gets a list of card on board accessible by the authenticated user.
func ListCards(ctx *context.APIContext) {
	//pr := ctx.Query("project_id")
	//
	//proj, _, err := ctx.DataSource.ListConnectBoard(pr)
	//if err != nil {
	//	ctx.JSON(http.StatusNotFound, &gitlab.ResponseError{
	//		Success: false,
	//		Message: err.Error(),
	//	})
	//	return
	//}
	//
	//current, err := ctx.DataSource.ItemBoard(pr)
	//
	//if err != nil {
	//	ctx.JSON(http.StatusNotFound, &gitlab.ResponseError{
	//		Success: false,
	//		Message: err.Error(),
	//	})
	//	return
	//}
	//
	//proj = append(proj, current)
	//
	//var cards, c []*gitlab.Card
	//
	//for _, p := range proj {
	//	c, err = ctx.DataSource.ListCards(p)
	//	cards = append(cards, c...)
	//}
	//
	//if err != nil {
	//	ctx.JSON(http.StatusUnauthorized, &gitlab.ResponseError{
	//		Success: false,
	//		Message: err.Error(),
	//	})
	//	return
	//}

	ctx.JSON(http.StatusOK, &gitlab.Response{
	//Data: cards,
	})
}

// CreateCard creates a new board card.
func CreateCard(ctx *context.APIContext, form gitlab.CardRequest) {
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
func UpdateCard(ctx *context.APIContext, form gitlab.CardRequest) {
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
func DeleteCard(ctx *context.APIContext, form gitlab.CardRequest) {
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
func MoveToCard(ctx *context.APIContext, form gitlab.CardRequest) {
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
	//	Event: "card.move",
	//})
	//
	//source := gitlab.ParseLabelToStage(form.Stage["source"])
	//dest := gitlab.ParseLabelToStage(form.Stage["dest"])
	//
	////TODO
	////autocomments := viper.GetBool("auto.comments")
	//autocomments := false
	//if source.Name != dest.Name && autocomments {
	//	com := gitlab.CommentRequest{
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
func ChangeProjectForCard(ctx *context.APIContext, form gitlab.CardRequest) {
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
