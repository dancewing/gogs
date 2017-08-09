// Copyright 2015 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package v4

import (
	"strings"

	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"

	api "github.com/gogits/go-gogs-client"

	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/models/errors"
	"github.com/gogits/gogs/pkg/context"
	"github.com/gogits/gogs/pkg/form"
	"github.com/gogits/gogs/routes/api/v1/repo"
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
			mode, err := models.AccessLevel(c.User.ID, repo)
			if err != nil {
				c.Error(500, "AccessLevel", err)
				return
			}
			c.Repo.AccessMode = mode
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
			c.Error(401)
			return
		}
	}
}

func reqBasicAuth() macaron.Handler {
	return func(c *context.Context) {
		if !c.IsBasicAuth {
			c.Error(401)
			return
		}
	}
}

func reqAdmin() macaron.Handler {
	return func(c *context.Context) {
		if !c.IsLogged || !c.User.IsAdmin {
			c.Error(403)
			return
		}
	}
}

func reqRepoWriter() macaron.Handler {
	return func(c *context.Context) {
		if !c.Repo.IsWriter() {
			c.Error(403)
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

// RegisterRoutes registers all v1 APIs routes to web application.
// FIXME: custom form error response
func RegisterRoutes(m *macaron.Macaron) {
	bind := binding.Bind

	m.Group("/v4", func() {
		// Handle preflight OPTIONS request
		m.Options("/*", func() {})

		m.Group("/repos", func() {
			m.Post("/migrate", bind(form.MigrateRepo{}), repo.Migrate)
			m.Combo("/:username/:reponame", repoAssignment()).Get(repo.Get).
				Delete(repo.Delete)

			m.Group("/:username/:reponame", func() {

				m.Post("/pipeline", bind(form.PipelineCallback{}), PipelineCallback)

				m.Group("/hooks", func() {
					m.Combo("").Get(repo.ListHooks).
						Post(bind(api.CreateHookOption{}), repo.CreateHook)
					m.Combo("/:id").Patch(bind(api.EditHookOption{}), repo.EditHook).
						Delete(repo.DeleteHook)
				})
				m.Group("/collaborators", func() {
					m.Get("", repo.ListCollaborators)
					m.Combo("/:collaborator").Get(repo.IsCollaborator).Put(bind(api.AddCollaboratorOption{}), repo.AddCollaborator).
						Delete(repo.DeleteCollaborator)
				})
				m.Get("/raw/*", context.RepoRef(), repo.GetRawFile)
				m.Get("/archive/*", repo.GetArchive)
				m.Get("/forks", repo.ListForks)
				m.Group("/branches", func() {
					m.Get("", repo.ListBranches)
					m.Get("/*", repo.GetBranch)
				})
				m.Group("/keys", func() {
					m.Combo("").Get(repo.ListDeployKeys).
						Post(bind(api.CreateKeyOption{}), repo.CreateDeployKey)
					m.Combo("/:id").Get(repo.GetDeployKey).
						Delete(repo.DeleteDeploykey)
				})
				m.Group("/issues", func() {
					m.Combo("").Get(repo.ListIssues).Post(bind(api.CreateIssueOption{}), repo.CreateIssue)
					m.Group("/comments", func() {
						m.Get("", repo.ListRepoIssueComments)
						m.Combo("/:id").Patch(bind(api.EditIssueCommentOption{}), repo.EditIssueComment)
					})
					m.Group("/:index", func() {
						m.Combo("").Get(repo.GetIssue).Patch(bind(api.EditIssueOption{}), repo.EditIssue)

						m.Group("/comments", func() {
							m.Combo("").Get(repo.ListIssueComments).Post(bind(api.CreateIssueCommentOption{}), repo.CreateIssueComment)
							m.Combo("/:id").Patch(bind(api.EditIssueCommentOption{}), repo.EditIssueComment).
								Delete(repo.DeleteIssueComment)
						})

						m.Group("/labels", func() {
							m.Combo("").Get(repo.ListIssueLabels).
								Post(bind(api.IssueLabelsOption{}), repo.AddIssueLabels).
								Put(bind(api.IssueLabelsOption{}), repo.ReplaceIssueLabels).
								Delete(repo.ClearIssueLabels)
							m.Delete("/:id", repo.DeleteIssueLabel)
						})

					})
				}, mustEnableIssues)
				m.Group("/labels", func() {
					m.Combo("").Get(repo.ListLabels).
						Post(bind(api.CreateLabelOption{}), repo.CreateLabel)
					m.Combo("/:id").Get(repo.GetLabel).Patch(bind(api.EditLabelOption{}), repo.EditLabel).
						Delete(repo.DeleteLabel)
				})
				m.Group("/milestones", func() {
					m.Combo("").Get(repo.ListMilestones).
						Post(reqRepoWriter(), bind(api.CreateMilestoneOption{}), repo.CreateMilestone)
					m.Combo("/:id").Get(repo.GetMilestone).
						Patch(reqRepoWriter(), bind(api.EditMilestoneOption{}), repo.EditMilestone).
						Delete(reqRepoWriter(), repo.DeleteMilestone)
				})
				m.Post("/mirror-sync", repo.MirrorSync)
				m.Get("/editorconfig/:filename", context.RepoRef(), repo.GetEditorconfig)
			}, repoAssignment())
		}, reqToken())

	}, context.APIContexter())
}
