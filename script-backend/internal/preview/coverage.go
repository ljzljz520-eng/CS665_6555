package preview

import "scriptstudio/script-backend/internal/domain"

type Coverage struct {
	CharacterID string `json:"characterId"`
	Name        string `json:"name"`
	Lines       int    `json:"lines"`
	Scenes      int    `json:"scenes"`
}

func CharacterCoverage(w domain.Workspace) []Coverage {
	out := make([]Coverage, 0, len(w.Characters))
	for _, character := range w.Characters {
		scenes := make(map[string]bool)
		count := 0
		for _, line := range w.Dialogues {
			if line.CharacterID == character.ID {
				count++
				scenes[line.SceneID] = true
			}
		}
		out = append(out, Coverage{CharacterID: character.ID, Name: character.Name, Lines: count, Scenes: len(scenes)})
	}
	return out
}
