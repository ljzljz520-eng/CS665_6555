package query

import (
	"scriptstudio/script-backend/internal/domain"
	"sort"
	"strings"
)

type Source interface {
	Scripts() ([]domain.Script, error)
	Workspace(string) (domain.Workspace, error)
}

type Query struct{ source Source }

func New(source Source) *Query { return &Query{source: source} }

type ScriptFilter struct {
	Text, Genre string
	Status      domain.DraftStatus
}

func (q *Query) Scripts(filter ScriptFilter) ([]domain.Script, error) {
	items, err := q.source.Scripts()
	if err != nil {
		return nil, err
	}
	text := strings.ToLower(strings.TrimSpace(filter.Text))
	genre := strings.ToLower(strings.TrimSpace(filter.Genre))
	out := make([]domain.Script, 0)
	for _, item := range items {
		if text != "" && !strings.Contains(strings.ToLower(item.Title+" "+item.Logline), text) {
			continue
		}
		if genre != "" && strings.ToLower(item.Genre) != genre {
			continue
		}
		if filter.Status != "" && item.Status != filter.Status {
			continue
		}
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Title < out[j].Title })
	return out, nil
}

type Statistics struct {
	SceneCount            int `json:"sceneCount"`
	CharacterCount        int `json:"characterCount"`
	DialogueCount         int `json:"dialogueCount"`
	WordCount             int `json:"wordCount"`
	ScenesWithoutDialogue int `json:"scenesWithoutDialogue"`
}

func (q *Query) Statistics(scriptID string) (Statistics, error) {
	w, err := q.source.Workspace(scriptID)
	if err != nil {
		return Statistics{}, err
	}
	result := Statistics{SceneCount: len(w.Scenes), CharacterCount: len(w.Characters), DialogueCount: len(w.Dialogues)}
	used := make(map[string]bool)
	for _, line := range w.Dialogues {
		result.WordCount += len(strings.Fields(line.Text))
		used[line.SceneID] = true
	}
	for _, scene := range w.Scenes {
		if !used[scene.ID] {
			result.ScenesWithoutDialogue++
		}
	}
	return result, nil
}
