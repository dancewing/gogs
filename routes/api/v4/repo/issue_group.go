package repo

import (
	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/pkg/context"
	"github.com/gogits/gogs/pkg/form"
)

const (
	GROUP_LABELS = "repo/issue/group_labels"
)

func RetrieveLabelGroups(c *context.RestContext) {
	labels, err := models.GetLabelGroupsByRepoID(c.Repo.Repository.ID)
	if err != nil {
		c.Handle(500, "RetrieveLabelGroups", err)
		return
	}
	groups := make([]string, len(labels))

	for i := range labels {
		groups[i] = labels[i].LabelGroup
	}

	c.Data["LabelGroups"] = groups

}

func RetrieveLabelsByGroup(c *context.RestContext) {

	currentGroup := c.Params(":group")

	labels, err := models.GetGroupedLabelsByRepoID(c.Repo.Repository.ID, currentGroup)

	if err != nil {
		c.Handle(500, "RetrieveLabels.GetLabels", err)
		return
	}
	//for _, l := range labels {
	//	l.CalOpenIssues()
	//}
	c.Data["CurrentGroup"] = currentGroup
	c.Data["Labels"] = labels
	c.Data["NumLabels"] = len(labels)
}

func LabelsByGroup(c *context.RestContext) {
	c.Data["Title"] = c.Tr("repo.labels")
	c.Data["PageIsIssueList"] = true
	c.Data["PageIsLabels"] = true
	c.Data["RequireMinicolors"] = true
	c.Data["LabelTemplates"] = models.LabelTemplates
	c.HTML(200, GROUP_LABELS)
}

func BatchUpdateLabelsByGroup(c *context.RestContext, f form.BatchUpdateLabel) {
	c.Data["Title"] = c.Tr("repo.labels")
	c.Data["PageIsIssueList"] = true
	c.Data["PageIsLabels"] = true
	c.Data["RequireMinicolors"] = true
	c.Data["LabelTemplates"] = models.LabelTemplates

	currentGroup := c.Params(":group")

	if f.Labels != nil {
		for _, l := range f.Labels {
			label, err := models.GetLabelByID(l.ID)
			if err == nil {
				label.Name = l.Title
				label.Color = l.Color
				label.LabelOrder = l.Order
				label.LabelGroup = l.Group
				err = models.UpdateLabel(label)

				if err != nil {

				}
			}
		}
	}

	c.Redirect(c.Repo.RepoLink+"/labels/group/"+currentGroup, 200)
}
