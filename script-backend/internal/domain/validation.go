package domain

import (
	"errors"
	"fmt"
	"strings"
)

func ValidateScript(v Script) error {
	if strings.TrimSpace(v.ID) == "" {
		return errors.New("script id is required")
	}
	if strings.TrimSpace(v.Title) == "" {
		return errors.New("script title is required")
	}
	if len([]rune(v.Title)) > 120 {
		return errors.New("script title is too long")
	}
	if v.Revision < 1 {
		return errors.New("revision must be positive")
	}
	switch v.Status {
	case StatusIdea, StatusDraft, StatusReview, StatusPublished:
		return nil
	default:
		return fmt.Errorf("unsupported status %q", v.Status)
	}
}

func ValidateScene(v Scene) error {
	if v.ID == "" || v.ScriptID == "" {
		return errors.New("scene identity is required")
	}
	if strings.TrimSpace(v.Heading) == "" {
		return errors.New("scene heading is required")
	}
	if v.Position < 1 {
		return errors.New("scene position must be positive")
	}
	if v.TimeOfDay != "" && v.TimeOfDay != "DAY" && v.TimeOfDay != "NIGHT" {
		return errors.New("time of day must be DAY or NIGHT")
	}
	return nil
}

func ValidateCharacter(v Character) error {
	if v.ID == "" || v.ScriptID == "" {
		return errors.New("character identity is required")
	}
	if strings.TrimSpace(v.Name) == "" {
		return errors.New("character name is required")
	}
	if len([]rune(v.Name)) > 80 {
		return errors.New("character name is too long")
	}
	return nil
}

func ValidateDialogue(v Dialogue) error {
	if v.ID == "" || v.SceneID == "" || v.CharacterID == "" {
		return errors.New("dialogue identity is required")
	}
	if strings.TrimSpace(v.Text) == "" {
		return errors.New("dialogue text is required")
	}
	if v.Position < 1 {
		return errors.New("dialogue position must be positive")
	}
	return nil
}

func ValidateDraft(v DraftState) error {
	if v.ScriptID == "" {
		return errors.New("draft script id is required")
	}
	if v.Version < 1 {
		return errors.New("draft version must be positive")
	}
	return nil
}
