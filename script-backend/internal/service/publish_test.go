package service

import (
	"path/filepath"
	"scriptstudio/script-backend/internal/domain"
	"scriptstudio/script-backend/internal/store"
	"testing"
)

func TestWorkflowPublishPreview(t *testing.T) {
	repo, err := store.Open(filepath.Join(t.TempDir(), "publish.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()
	app := New(repo)
	script, _ := app.CreateScript(CreateScriptInput{RequestKey: "publish-1", Title: "Signal", Genre: "thriller"})
	scene, _ := app.AddScene(AddSceneInput{ScriptID: script.ID, Heading: "INT. RADIO ROOM", TimeOfDay: "NIGHT"})
	character, _ := app.AddCharacter(AddCharacterInput{ScriptID: script.ID, Name: "Jo"})
	if _, err := app.AddDialogue(AddDialogueInput{SceneID: scene.ID, CharacterID: character.ID, Text: "The signal is ours."}); err != nil {
		t.Fatal(err)
	}
	if _, err := app.TransitionDraft(script.ID, domain.StatusDraft, "first draft"); err != nil {
		t.Fatal(err)
	}
	state, err := app.Publish(script.ID, "approved")
	if err != nil {
		t.Fatal(err)
	}
	w, err := app.Workspace(script.ID)
	if err != nil {
		t.Fatal(err)
	}
	if state.State != domain.StatusPublished || w.Script.Status != domain.StatusPublished {
		t.Fatalf("unexpected state: %#v", state)
	}
}
