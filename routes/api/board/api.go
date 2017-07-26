// Copyright 2015 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package board

import (
	"strings"

	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"

	"net/http"

	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/models/errors"
	"github.com/gogits/gogs/pkg/context"
	"github.com/gogits/gogs/routes/api/board/form"
)

func repoAssignment() macaron.Handler {
	return func(c *context.APIContext) {
		userName := c.Params(":username")
		repoName := c.Params(":reponame")

		var (
			owner *models.User
			err   error
		)

		// Check if the user is the same as the repository owner.
		if c.IsLogged && c.User.LowerName == strings.ToLower(userName) {
			owner = c.User
		} else {
			owner, err = models.GetUserByName(userName)
			if err != nil {
				if errors.IsUserNotExist(err) {
					c.Status(404)
				} else {
					c.Error(500, "GetUserByName", err)
				}
				return
			}
		}
		c.Repo.Owner = owner

		// Get repository.
		repo, err := models.GetRepositoryByName(owner.ID, repoName)
		if err != nil {
			if errors.IsRepoNotExist(err) {
				c.Status(404)
			} else {
				c.Error(500, "GetRepositoryByName", err)
			}
			return
		} else if err = repo.GetOwner(); err != nil {
			c.Error(500, "GetOwner", err)
			return
		}

		if c.IsLogged && c.User.IsAdmin {
			c.Repo.AccessMode = models.ACCESS_MODE_OWNER
		} else {

			if c.IsLogged {
				mode, err := models.AccessLevel(c.User.ID, repo)
				if err != nil {
					c.Error(500, "AccessLevel", err)
					return
				}
				c.Repo.AccessMode = mode
			} else {
				c.Repo.AccessMode = models.ACCESS_MODE_READ
			}

		}

		if !c.Repo.HasAccess() {
			c.Status(404)
			return
		}

		c.Repo.Repository = repo
	}
}

// Contexter middleware already checks token for user sign in process.
func reqToken() macaron.Handler {
	return func(c *context.Context) {
		if !c.IsLogged {
			//c.Error(401)
			c.JSON(http.StatusUnauthorized, &form.ResponseError{
				Success: false,
				Message: "Unauthorized",
			})
			return
		}
	}
}

func reqBasicAuth() macaron.Handler {
	return func(c *context.Context) {
		if !c.IsBasicAuth {
			//c.Error(401)
			c.JSON(http.StatusUnauthorized, &form.ResponseError{
				Success: false,
				Message: "Unauthorized",
			})
			return
		}
	}
}

func reqAdmin() macaron.Handler {
	return func(c *context.Context) {
		if !c.IsLogged || !c.User.IsAdmin {
			//c.Error(403)
			c.JSON(http.StatusForbidden, &form.ResponseError{
				Success: false,
				Message: "Forbidden",
			})
			return
		}
	}
}

func reqRepoWriter() macaron.Handler {
	return func(c *context.Context) {
		if !c.Repo.IsWriter() {
			//c.Error(403)
			c.JSON(http.StatusForbidden, &form.ResponseError{
				Success: false,
				Message: "Forbidden",
			})
			return
		}
	}
}

func orgAssignment(args ...bool) macaron.Handler {
	var (
		assignOrg  bool
		assignTeam bool
	)
	if len(args) > 0 {
		assignOrg = args[0]
	}
	if len(args) > 1 {
		assignTeam = args[1]
	}
	return func(c *context.APIContext) {
		c.Org = new(context.APIOrganization)

		var err error
		if assignOrg {
			c.Org.Organization, err = models.GetUserByName(c.Params(":orgname"))
			if err != nil {
				if errors.IsUserNotExist(err) {
					c.Status(404)
				} else {
					c.Error(500, "GetUserByName", err)
				}
				return
			}
		}

		if assignTeam {
			c.Org.Team, err = models.GetTeamByID(c.ParamsInt64(":teamid"))
			if err != nil {
				if errors.IsUserNotExist(err) {
					c.Status(404)
				} else {
					c.Error(500, "GetTeamById", err)
				}
				return
			}
		}
	}
}

func mustEnableIssues(c *context.APIContext) {
	if !c.Repo.Repository.EnableIssues || c.Repo.Repository.EnableExternalTracker {
		c.Status(404)
		return
	}
}

func RegisterBoardRoutes(m *macaron.Macaron) {

	m.Group("/boards", func() {

		//m.Post("/configure", binding.Json(gitlab.BoardRequest{}), Configure)

		m.Group("/:username/:reponame", func() {

			m.Get("", ItemBoard)

			m.Get("/current", GetAuthenticatedUser)

			m.Post("/configure", binding.Json(form.BoardRequest{}), Configure)

			m.Combo("/labels").
				Get(ListLabels).
				Put(binding.Json(form.LabelRequest{}), EditLabel).
				Delete(DeleteLabel).
				Post(binding.Json(form.LabelRequest{}), CreateLabel)

			m.Combo("/connect").
				Get(ListConnectBoard).
				Post(binding.Json(form.BoardRequest{}), CreateConnectBoard).
				Delete(DeleteConnectBoard)

			m.Group("/cards", func() {
				m.Get("", ListCards)
				m.Combo("/:card").
					Post(binding.Json(form.CardRequest{}), CreateCard).
					Put(binding.Json(form.CardRequest{}), UpdateCard).
					Delete(binding.Json(form.CardRequest{}), DeleteCard)

				m.Group("/:index", func() {
					m.Combo("/comments").
						Get(ListComments).
						Post(binding.Json(form.CommentRequest{}), CreateComment)
				})
			})

			m.Combo("/milestones").
				Get(ListMilestones).
				Post(binding.Json(form.MilestoneRequest{}), CreateMilestone)

			m.Get("/users", ListMembers)

			m.Combo("/move").
				Put(reqToken(), binding.Json(form.CardRequest{}), MoveToCard).
				Post(reqToken(), binding.Json(form.CardRequest{}), ChangeProjectForCard)

			m.Post("/upload", binding.MultipartForm(form.UploadForm{}), UploadFile)
		}, repoAssignment())

	}, context.APIContexter())

	m.Group("/card/:board", func() {
		m.Combo("").
			Post(binding.Json(form.CardRequest{}), CreateCard).
			Put(binding.Json(form.CardRequest{}), UpdateCard).
			Delete(binding.Json(form.CardRequest{}), DeleteCard)

		m.Put("/move", binding.Json(form.CardRequest{}), MoveToCard)
		m.Post("/move/:projectId", binding.Json(form.CardRequest{}), ChangeProjectForCard)

	})

}
