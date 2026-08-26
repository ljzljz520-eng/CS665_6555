package query

import (
	"path/filepath"
	"scriptstudio/script-backend/internal/service"
	"scriptstudio/script-backend/internal/store"
	"testing"
)

func TestQueryFiltersAndStatistics(t *testing.T) {
	repo, err := store.Open(filepath.Join(t.TempDir(), "query.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()
	app := service.New(repo)
	script, _ := app.CreateScript(service.CreateScriptInput{RequestKey: "q-1", Title: "Glass City", Logline: "An architect sees hidden streets.", Genre: "mystery"})
	if _, err := app.AddScene(service.AddSceneInput{ScriptID: script.ID, Heading: "EXT. SQUARE", TimeOfDay: "DAY"}); err != nil {
		t.Fatal(err)
	}
	q := New(repo)
	items, err := q.Scripts(ScriptFilter{Text: "hidden", Genre: "mystery"})
	if err != nil {
		t.Fatal(err)
	}
	stats, err := q.Statistics(script.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || stats.SceneCount != 1 || stats.ScenesWithoutDialogue != 1 {
		t.Fatalf("unexpected query result: %#v %#v", items, stats)
	}
}
