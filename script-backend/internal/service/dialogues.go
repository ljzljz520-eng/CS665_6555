package service

import (
	"errors"
	"scriptstudio/script-backend/internal/domain"
	"strings"
)

type AddCharacterInput struct{ ScriptID, Name, Bio, Objective string }

func (s *Service) AddCharacter(in AddCharacterInput) (domain.Character, error) {
	if err := s.requireScript(in.ScriptID); err != nil {
		return domain.Character{}, err
	}
	items, err := s.repo.Characters()
	if err != nil {
		return domain.Character{}, err
	}
	for _, item := range items {
		if item.ScriptID == in.ScriptID && strings.EqualFold(item.Name, in.Name) {
			return domain.Character{}, errors.New("character name already exists")
		}
	}
	v := domain.Character{ID: nextID("character", len(items)), ScriptID: in.ScriptID, Name: strings.TrimSpace(in.Name), Bio: strings.TrimSpace(in.Bio), Objective: strings.TrimSpace(in.Objective)}
	if err := domain.ValidateCharacter(v); err != nil {
		return domain.Character{}, err
	}
	return v, s.repo.SaveCharacter(v)
}

type AddDialogueInput struct{ SceneID, CharacterID, Text, Direction string }

func (s *Service) AddDialogue(in AddDialogueInput) (domain.Dialogue, error) {
	scenes, err := s.repo.Scenes()
	if err != nil {
		return domain.Dialogue{}, err
	}
	characters, err := s.repo.Characters()
	if err != nil {
		return domain.Dialogue{}, err
	}
	sceneFound, characterFound := false, false
	for _, item := range scenes {
		if item.ID == in.SceneID {
			sceneFound = true
		}
	}
	for _, item := range characters {
		if item.ID == in.CharacterID {
			characterFound = true
		}
	}
	if !sceneFound || !characterFound {
		return domain.Dialogue{}, errors.New("scene and character are required")
	}
	items, err := s.repo.Dialogues()
	if err != nil {
		return domain.Dialogue{}, err
	}
	owned := make([]domain.Dialogue, 0)
	for _, item := range items {
		if item.SceneID == in.SceneID {
			owned = append(owned, item)
		}
	}
	v := domain.Dialogue{ID: nextID("dialogue", len(items)), SceneID: in.SceneID, CharacterID: in.CharacterID, Text: strings.TrimSpace(in.Text), Direction: strings.TrimSpace(in.Direction), Position: domain.NextDialoguePosition(owned)}
	if err := domain.ValidateDialogue(v); err != nil {
		return domain.Dialogue{}, err
	}
	return v, s.repo.SaveDialogue(v)
}

func (s *Service) DialogueList(sceneID string) ([]domain.Dialogue, error) {
	items, err := s.repo.Dialogues()
	if err != nil {
		return nil, err
	}
	out := make([]domain.Dialogue, 0)
	for _, item := range items {
		if item.SceneID == sceneID {
			out = append(out, item)
		}
	}
	return domain.SortDialogues(out), nil
}
