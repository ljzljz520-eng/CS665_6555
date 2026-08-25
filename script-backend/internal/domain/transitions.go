package domain

import "fmt"

func CanTransition(from, to DraftStatus) bool {
	if from == to {
		return true
	}
	switch from {
	case StatusIdea:
		return to == StatusDraft
	case StatusDraft:
		return to == StatusIdea || to == StatusReview
	case StatusReview:
		return to == StatusDraft || to == StatusPublished
	case StatusPublished:
		return to == StatusDraft
	default:
		return false
	}
}

func Transition(v DraftState, to DraftStatus, note string) (DraftState, error) {
	if !CanTransition(v.State, to) {
		return DraftState{}, fmt.Errorf("cannot transition from %s to %s", v.State, to)
	}
	v.State = to
	v.Note = note
	v.Version++
	return v, nil
}

func PublicationIssues(w Workspace) []string {
	issues := make([]string, 0)
	if len(w.Scenes) == 0 {
		issues = append(issues, "at least one scene is required")
	}
	if len(w.Characters) == 0 {
		issues = append(issues, "at least one character is required")
	}
	if len(w.Dialogues) == 0 {
		issues = append(issues, "at least one dialogue is required")
	}
	for _, scene := range w.Scenes {
		found := false
		for _, line := range w.Dialogues {
			if line.SceneID == scene.ID {
				found = true
				break
			}
		}
		if !found {
			issues = append(issues, "scene "+scene.ID+" has no dialogue")
		}
	}
	return issues
}
