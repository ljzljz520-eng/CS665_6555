package preview

import (
	"scriptstudio/script-backend/internal/domain"
	"strings"
	"testing"
)

func TestRenderReadingPreview(t *testing.T) {
	w := domain.Workspace{Script: domain.Script{Title: "Open Water"}, Draft: domain.DraftState{State: domain.StatusDraft}, Scenes: []domain.Scene{{ID: "s1", Heading: "EXT. BOAT", Location: "SEA", TimeOfDay: "DAY", Position: 1}}, Characters: []domain.Character{{ID: "c1", Name: "Ari"}}, Dialogues: []domain.Dialogue{{SceneID: "s1", CharacterID: "c1", Text: "Hold the line.", Position: 1}}}
	doc := Render(w)
	if doc.SceneCount != 1 || !strings.Contains(doc.Text, "ARI") || !strings.Contains(doc.Text, "Hold the line.") {
		t.Fatalf("unexpected preview: %#v", doc)
	}
}
