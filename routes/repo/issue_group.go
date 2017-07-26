package repo

import (


	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/pkg/context"
)

const (
	GROUP_LABELS = "repo/issue/group_labels"
)

func RetrieveLabelGroups(c *context.Context) {
	labels, err := models.GetLabelGroupsByRepoID(c.Repo.Repository.ID)
	if err != nil {
		c.Handle(500, "RetrieveLabelGroups", err)
		return
	}
	groups := make([]string,len(labels))

	for i := range labels {
		groups[i] = labels[i].LabelGroup
	}

	c.Data["LabelGroups"] = groups
}

func RetrieveLabelsByGroup(c *context.Context) {

	labels, err := models.GetGroupedLabelsByRepoID(c.Repo.Repository.ID, c.Params(":group"))

	if err != nil {
		c.Handle(500, "RetrieveLabels.GetLabels", err)
		return
	}
	for _, l := range labels {
		l.CalOpenIssues()
	}
	c.Data["Labels"] = labels
	c.Data["NumLabels"] = len(labels)
}

func LabelsByGroup(c *context.Context) {
	c.Data["Title"] = c.Tr("repo.labels")
	c.Data["PageIsIssueList"] = true
	c.Data["PageIsLabels"] = true
	c.Data["RequireMinicolors"] = true
	c.Data["LabelTemplates"] = models.LabelTemplates
	c.HTML(200, GROUP_LABELS)
}
