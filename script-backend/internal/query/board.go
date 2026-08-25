package query

import "scriptstudio/script-backend/internal/domain"

type BoardColumn struct {
	Status  domain.DraftStatus `json:"status"`
	Scripts []domain.Script    `json:"scripts"`
}

func (q *Query) Board() ([]BoardColumn, error) {
	items, err := q.source.Scripts()
	if err != nil {
		return nil, err
	}
	columns := []BoardColumn{{Status: domain.StatusIdea}, {Status: domain.StatusDraft}, {Status: domain.StatusReview}, {Status: domain.StatusPublished}}
	for _, item := range items {
		for i := range columns {
			if columns[i].Status == item.Status {
				columns[i].Scripts = append(columns[i].Scripts, item)
				break
			}
		}
	}
	return columns, nil
}
