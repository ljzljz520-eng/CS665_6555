package service

import (
	"path/filepath"
	"scriptstudio/script-backend/internal/store"
	"testing"
)

func TestWorkflowCreateScript(t *testing.T) {
	repo, err := store.Open(filepath.Join(t.TempDir(), "author.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()
	app := New(repo)
	script, err := app.CreateScript(CreateScriptInput{RequestKey: "create-1", Title: "Harbor Light", Logline: "A keeper confronts an abandoned promise.", Genre: "drama"})
	if err != nil {
		t.Fatal(err)
	}
	scene, err := app.AddScene(AddSceneInput{ScriptID: script.ID, Heading: "EXT. LIGHTHOUSE", Synopsis: "A storm arrives", Location: "CLIFF", TimeOfDay: "NIGHT"})
	if err != nil {
		t.Fatal(err)
	}
	revised, err := app.ReviseScript(script.ID, "Harbor Lights", "", "")
	if err != nil {
		t.Fatal(err)
	}
	w, err := app.Workspace(script.ID)
	if err != nil {
		t.Fatal(err)
	}
	if scene.Position != 1 || revised.Revision != 2 || len(w.Scenes) != 1 {
		t.Fatalf("unexpected workflow result: %#v", w)
	}
}
