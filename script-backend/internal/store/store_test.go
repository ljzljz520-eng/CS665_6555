package store

import (
	"path/filepath"
	"scriptstudio/script-backend/internal/domain"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "studio.db")
	repo, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	script := domain.Script{ID: "script-1", Title: "Northbound", Status: domain.StatusDraft, Revision: 2}
	if err := repo.SaveScript(script); err != nil {
		t.Fatal(err)
	}
	if err := repo.SaveScene(domain.Scene{ID: "scene-1", ScriptID: script.ID, Heading: "INT. TRAIN", Position: 1}); err != nil {
		t.Fatal(err)
	}
	if err := repo.SaveCharacter(domain.Character{ID: "character-1", ScriptID: script.ID, Name: "Mara"}); err != nil {
		t.Fatal(err)
	}
	if err := repo.SaveDialogue(domain.Dialogue{ID: "dialogue-1", SceneID: "scene-1", CharacterID: "character-1", Text: "We leave at dawn.", Position: 1}); err != nil {
		t.Fatal(err)
	}
	if err := repo.SaveDraft(domain.DraftState{ScriptID: script.ID, State: domain.StatusDraft, Version: 2}); err != nil {
		t.Fatal(err)
	}
	if err := repo.Close(); err != nil {
		t.Fatal(err)
	}
	repo, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()
	w, err := repo.Workspace(script.ID)
	if err != nil {
		t.Fatal(err)
	}
	if w.Script.Title != "Northbound" || len(w.Scenes) != 1 || len(w.Characters) != 1 || len(w.Dialogues) != 1 || w.Draft.Version != 2 {
		t.Fatalf("unexpected reopened workspace: %#v", w)
	}
}
