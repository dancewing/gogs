// Copyright 2015 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package v2

import (
	"strings"

	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"

	api "github.com/gogits/go-gogs-client"

	"net/http"

	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/models/errors"
	"github.com/gogits/gogs/pkg/context"
	"github.com/gogits/gogs/pkg/form"
	"github.com/gogits/gogs/routes/api/v2/admin"
	"github.com/gogits/gogs/routes/api/v2/gitlab"
	"github.com/gogits/gogs/routes/api/v2/misc"
	"github.com/gogits/gogs/routes/api/v2/org"
	"github.com/gogits/gogs/routes/api/v2/repo"
	"github.com/gogits/gogs/routes/api/v2/user"
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

// RegisterRoutes registers all v1 APIs routes to web application.
// FIXME: custom form error response
func RegisterRoutes(m *macaron.Macaron) {
	bind := binding.Bind

	m.Group("/v2", func() {
		// Handle preflight OPTIONS request
		m.Options("/*", func() {})

		// Miscellaneous
		m.Post("/markdown", bind(api.MarkdownOption{}), misc.Markdown)
		m.Post("/markdown/raw", misc.MarkdownRaw)

		// Users
		m.Group("/users", func() {
			m.Get("/search", user.Search)

			m.Group("/:username", func() {
				m.Get("", user.GetInfo)

				m.Group("/tokens", func() {
					m.Combo("").Get(user.ListAccessTokens).
						Post(bind(api.CreateAccessTokenOption{}), user.CreateAccessToken)
				}, reqBasicAuth())
			})
		})

		m.Group("/users", func() {
			m.Group("/:username", func() {
				m.Get("/keys", user.ListPublicKeys)

				m.Get("/followers", user.ListFollowers)
				m.Group("/following", func() {
					m.Get("", user.ListFollowing)
					m.Get("/:target", user.CheckFollowing)
				})
			})
		}, reqToken())

		m.Group("/user", func() {
			m.Get("", user.GetAuthenticatedUser)
			m.Combo("/emails").Get(user.ListEmails).
				Post(bind(api.CreateEmailOption{}), user.AddEmail).
				Delete(bind(api.CreateEmailOption{}), user.DeleteEmail)

			m.Get("/followers", user.ListMyFollowers)
			m.Group("/following", func() {
				m.Get("", user.ListMyFollowing)
				m.Combo("/:username").Get(user.CheckMyFollowing).Put(user.Follow).Delete(user.Unfollow)
			})

			m.Group("/keys", func() {
				m.Combo("").Get(user.ListMyPublicKeys).
					Post(bind(api.CreateKeyOption{}), user.CreatePublicKey)
				m.Combo("/:id").Get(user.GetPublicKey).
					Delete(user.DeletePublicKey)
			})

			m.Combo("/issues").Get(repo.ListUserIssues)
		}, reqToken())

		// Repositories
		m.Get("/users/:username/repos", reqToken(), repo.ListUserRepositories)
		m.Get("/orgs/:org/repos", reqToken(), repo.ListOrgRepositories)

		m.Combo("/user/repos", reqToken()).Get(repo.ListMyRepos).
			Post(bind(api.CreateRepoOption{}), repo.Create)

		m.Post("/org/:org/repos", reqToken(), bind(api.CreateRepoOption{}), repo.CreateOrgRepo)

		m.Group("/repos", func() {
			m.Get("/search", repo.Search)
		})

		m.Group("/repos", func() {
			m.Post("/migrate", bind(form.MigrateRepo{}), repo.Migrate)
			m.Combo("/:username/:reponame", repoAssignment()).Get(repo.Get).
				Delete(repo.Delete)

			m.Group("/:username/:reponame", func() {
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
					m.Group("/:issue", func() {
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

				m.Combo("/connect").
					Get(repo.ListConnectBoard).
					Post(repo.CreateConnectBoard).
					Delete(repo.DeleteConnectBoard)

				//m.Put("/move", bind(gitlab.CardRequest{}), repo.MoveToCard)

				m.Group("/labels", func() {
					m.Combo("").Get(repo.ListLabels).
						Post(bind(api.CreateLabelOption{}), repo.CreateLabel)
					m.Combo("/:id").Get(repo.GetLabel).Patch(bind(api.EditLabelOption{}), repo.EditLabel).
						Delete(repo.DeleteLabel)

					m.Post("/createKBLabels", repo.CreateKBIssueLabels)
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

		m.Get("/issues", reqToken(), repo.ListUserIssues)

		// Organizations
		m.Get("/user/orgs", reqToken(), org.ListMyOrgs)
		m.Get("/users/:username/orgs", org.ListUserOrgs)
		m.Group("/orgs/:orgname", func() {
			m.Combo("").Get(org.Get).Patch(bind(api.EditOrgOption{}), org.Edit)
			m.Combo("/teams").Get(org.ListTeams)
		}, orgAssignment(true))

		m.Any("/*", func(c *context.Context) {
			c.Error(404)
		})

		m.Group("/admin", func() {
			m.Group("/users", func() {
				m.Post("", bind(api.CreateUserOption{}), admin.CreateUser)

				m.Group("/:username", func() {
					m.Combo("").Patch(bind(api.EditUserOption{}), admin.EditUser).
						Delete(admin.DeleteUser)
					m.Post("/keys", bind(api.CreateKeyOption{}), admin.CreatePublicKey)
					m.Post("/orgs", bind(api.CreateOrgOption{}), admin.CreateOrg)
					m.Post("/repos", bind(api.CreateRepoOption{}), admin.CreateRepo)
				})
			})

			m.Group("/orgs/:orgname", func() {
				m.Group("/teams", func() {
					m.Post("", orgAssignment(true), bind(api.CreateTeamOption{}), admin.CreateTeam)
				})
			})
			m.Group("/teams", func() {
				m.Group("/:teamid", func() {
					m.Combo("/members/:username").Put(admin.AddTeamMember).Delete(admin.RemoveTeamMember)
					m.Combo("/repos/:reponame").Put(admin.AddTeamRepository).Delete(admin.RemoveTeamRepository)
				}, orgAssignment(false, true))
			})
		}, reqAdmin())
	}, context.APIContexter())
}

func RegisterKBRoutes(m *macaron.Macaron) {

	m.Get("/labels/:project", ListLabels)
	m.Put("/labels/:project", binding.Json(gitlab.LabelRequest{}), EditLabel)
	m.Delete("/labels/:project/:label", DeleteLabel)
	m.Post("/labels/:project", binding.Json(gitlab.LabelRequest{}), CreateLabel)

	m.Group("/boards", func() {
		m.Get("", ListBoards)
		m.Get("/starred", ListStarredBoards)
		//m.Post("/configure", binding.Json(gitlab.BoardRequest{}), Configure)

		m.Group("/:username/:reponame", func() {

			m.Get("", ItemBoard)

			m.Post("/configure", binding.Json(gitlab.BoardRequest{}), Configure)

			m.Combo("/labels").
				Get(ListLabels).
				Put(binding.Json(gitlab.LabelRequest{}), EditLabel).
				Delete(DeleteLabel).
				Post(binding.Json(gitlab.LabelRequest{}), CreateLabel)

			m.Combo("/connect").
				Get(ListConnectBoard).
				Post(binding.Json(gitlab.BoardRequest{}), CreateConnectBoard).
				Delete(DeleteConnectBoard)

			m.Group("/cards", func() {
				m.Get("", ListCards)
				m.Combo("/:card").
					Post(binding.Json(gitlab.CardRequest{}), CreateCard).
					Put(binding.Json(gitlab.CardRequest{}), UpdateCard).
					Delete(binding.Json(gitlab.CardRequest{}), DeleteCard)
			})

			m.Combo("/milestones").
				Get(ListMilestones).
				Post(binding.Json(gitlab.MilestoneRequest{}), CreateMilestone)

			m.Get("/users", ListMembers)

			m.Combo("/move").
				Put(binding.Json(gitlab.CardRequest{}), MoveToCard).
				Post(binding.Json(gitlab.CardRequest{}), ChangeProjectForCard)

			//m.Put("/move", binding.Json(gitlab.CardRequest{}), MoveToCard)
			//m.Post("/move/:projectId", binding.Json(gitlab.CardRequest{}), ChangeProjectForCard)

			m.Combo("/comments").
				Get(ListComments).
				Post(binding.Json(gitlab.CommentRequest{}), CreateComment)

			m.Post("/upload", binding.MultipartForm(gitlab.UploadForm{}), UploadFile)
		}, repoAssignment())

	}, context.APIContexter())

	//m.Get("/board", ItemBoard)

	m.Get("/cards", ListCards)
	m.Combo("/milestones").
		Get(ListMilestones).
		Post(binding.Json(gitlab.MilestoneRequest{}), CreateMilestone)

	m.Get("/users", ListMembers)
	m.Combo("/comments").
		Get(ListComments).
		Post(binding.Json(gitlab.CommentRequest{}), CreateComment)

	m.Group("/card/:board", func() {
		m.Combo("").
			Post(binding.Json(gitlab.CardRequest{}), CreateCard).
			Put(binding.Json(gitlab.CardRequest{}), UpdateCard).
			Delete(binding.Json(gitlab.CardRequest{}), DeleteCard)

		m.Put("/move", binding.Json(gitlab.CardRequest{}), MoveToCard)
		m.Post("/move/:projectId", binding.Json(gitlab.CardRequest{}), ChangeProjectForCard)

	})

	m.Get("/current", user.GetAuthenticatedUser)

}
