package service

import (
	"path/filepath"
	"scriptstudio/script-backend/internal/store"
	"testing"
)

func TestWorkflowComposeDialogue(t *testing.T) {
	repo, err := store.Open(filepath.Join(t.TempDir(), "dialogue.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()
	app := New(repo)
	script, _ := app.CreateScript(CreateScriptInput{RequestKey: "dialogue-1", Title: "After Rain", Genre: "romance"})
	scene, err := app.AddScene(AddSceneInput{ScriptID: script.ID, Heading: "EXT. PLATFORM", TimeOfDay: "DAY"})
	if err != nil {
		t.Fatal(err)
	}
	character, err := app.AddCharacter(AddCharacterInput{ScriptID: script.ID, Name: "Lin", Objective: "Tell the truth"})
	if err != nil {
		t.Fatal(err)
	}
	line, err := app.AddDialogue(AddDialogueInput{SceneID: scene.ID, CharacterID: character.ID, Text: "I waited.", Direction: "quietly"})
	if err != nil {
		t.Fatal(err)
	}
	items, err := app.DialogueList(scene.ID)
	if err != nil {
		t.Fatal(err)
	}
	if line.Position != 1 || len(items) != 1 || items[0].Text != "I waited." {
		t.Fatalf("unexpected dialogues: %#v", items)
	}
}
