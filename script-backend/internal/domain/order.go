package domain

import "sort"

func SortScenes(items []Scene) []Scene {
	out := append([]Scene(nil), items...)
	sort.SliceStable(out, func(i, j int) bool { return out[i].Position < out[j].Position })
	return out
}

func SortDialogues(items []Dialogue) []Dialogue {
	out := append([]Dialogue(nil), items...)
	sort.SliceStable(out, func(i, j int) bool { return out[i].Position < out[j].Position })
	return out
}

func CharacterIndex(items []Character) map[string]Character {
	out := make(map[string]Character, len(items))
	for _, item := range items {
		out[item.ID] = item
	}
	return out
}

func NextScenePosition(items []Scene) int {
	max := 0
	for _, item := range items {
		if item.Position > max {
			max = item.Position
		}
	}
	return max + 1
}

func NextDialoguePosition(items []Dialogue) int {
	max := 0
	for _, item := range items {
		if item.Position > max {
			max = item.Position
		}
	}
	return max + 1
}
