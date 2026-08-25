package service

import (
	"scriptstudio/script-backend/internal/domain"
	"strings"
)

type AddSceneInput struct {
	ScriptID  string `json:"scriptId"`
	Heading   string `json:"heading"`
	Synopsis  string `json:"synopsis"`
	Location  string `json:"location"`
	TimeOfDay string `json:"timeOfDay"`
}

func (s *Service) AddScene(in AddSceneInput) (domain.Scene, error) {
	if err := s.requireScript(in.ScriptID); err != nil {
		return domain.Scene{}, err
	}
	items, err := s.repo.Scenes()
	if err != nil {
		return domain.Scene{}, err
	}
	owned := make([]domain.Scene, 0)
	for _, item := range items {
		if item.ScriptID == in.ScriptID {
			owned = append(owned, item)
		}
	}
	v := domain.Scene{ID: nextID("scene", len(items)), ScriptID: in.ScriptID, Heading: strings.TrimSpace(in.Heading), Synopsis: strings.TrimSpace(in.Synopsis), Location: strings.TrimSpace(in.Location), TimeOfDay: strings.ToUpper(strings.TrimSpace(in.TimeOfDay)), Position: domain.NextScenePosition(owned)}
	if err := domain.ValidateScene(v); err != nil {
		return domain.Scene{}, err
	}
	return v, s.repo.SaveScene(v)
}

func (s *Service) MoveScene(scriptID, sceneID string, position int) error {
	if position < 1 {
		return domain.ValidateScene(domain.Scene{})
	}
	items, err := s.repo.Scenes()
	if err != nil {
		return err
	}
	found := false
	for _, item := range items {
		if item.ScriptID != scriptID {
			continue
		}
		if item.ID == sceneID {
			item.Position = position
			found = true
		} else if item.Position >= position {
			item.Position++
		}
		if err := s.repo.SaveScene(item); err != nil {
			return err
		}
	}
	if !found {
		return s.requireScript("missing-scene")
	}
	return nil
}

func (s *Service) SceneList(scriptID string) ([]domain.Scene, error) {
	items, err := s.repo.Scenes()
	if err != nil {
		return nil, err
	}
	out := make([]domain.Scene, 0)
	for _, item := range items {
		if item.ScriptID == scriptID {
			out = append(out, item)
		}
	}
	return domain.SortScenes(out), nil
}
