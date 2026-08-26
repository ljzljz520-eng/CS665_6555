package preview

import (
	"fmt"
	"scriptstudio/script-backend/internal/domain"
	"strings"
)

type Document struct {
	Title      string             `json:"title"`
	Status     domain.DraftStatus `json:"status"`
	Text       string             `json:"text"`
	SceneCount int                `json:"sceneCount"`
}

func Render(w domain.Workspace) Document {
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n%s\n\n", strings.ToUpper(w.Script.Title), w.Script.Logline)
	characters := domain.CharacterIndex(w.Characters)
	linesByScene := make(map[string][]domain.Dialogue)
	for _, line := range w.Dialogues {
		linesByScene[line.SceneID] = append(linesByScene[line.SceneID], line)
	}
	for _, scene := range domain.SortScenes(w.Scenes) {
		fmt.Fprintf(&b, "%s - %s - %s\n", strings.ToUpper(scene.Heading), strings.ToUpper(scene.Location), scene.TimeOfDay)
		if scene.Synopsis != "" {
			fmt.Fprintf(&b, "[%s]\n", scene.Synopsis)
		}
		for _, line := range domain.SortDialogues(linesByScene[scene.ID]) {
			name := "UNKNOWN"
			if character, ok := characters[line.CharacterID]; ok {
				name = strings.ToUpper(character.Name)
			}
			fmt.Fprintf(&b, "\n%s\n", name)
			if line.Direction != "" {
				fmt.Fprintf(&b, "(%s)\n", line.Direction)
			}
			fmt.Fprintf(&b, "%s\n", line.Text)
		}
		b.WriteString("\n")
	}
	return Document{Title: w.Script.Title, Status: w.Draft.State, Text: b.String(), SceneCount: len(w.Scenes)}
}
