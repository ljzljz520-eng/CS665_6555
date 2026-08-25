package catalog

import "testing"

func TestBlueprintCatalogSupportsEditor(t *testing.T) {
	genres := Genres()
	drama := ForGenre("drama")
	results := Search("主人公")
	if len(genres) < 20 || len(drama) != 10 || len(results) == 0 {
		t.Fatalf("unexpected catalog sizes: %d %d %d", len(genres), len(drama), len(results))
	}
}
