package models

func GetLabelGroupsByRepoID(repoID int64) ([]*Label, error) {
	labels := make([]*Label, 0, 10)
	return labels, x.Where("repo_id = ?", repoID).GroupBy("label_group").Find(&labels)
}

func GetGroupedLabelsByRepoID(repoID int64, group string) ([]*Label, error) {
	labels := make([]*Label, 0, 10)
	return labels, x.Where("repo_id = ? and label_group = ?", repoID, group).Find(&labels)
}
