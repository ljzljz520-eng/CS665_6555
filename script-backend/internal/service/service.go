package service

import (
	"fmt"
	"strconv"
	"strings"

	"scriptstudio/script-backend/internal/domain"
)

type Repository interface {
	SaveScript(domain.Script) error
	Script(string) (domain.Script, error)
	Scripts() ([]domain.Script, error)
	SaveScene(domain.Scene) error
	Scenes() ([]domain.Scene, error)
	SaveCharacter(domain.Character) error
	Characters() ([]domain.Character, error)
	SaveDialogue(domain.Dialogue) error
	Dialogues() ([]domain.Dialogue, error)
	SaveDraft(domain.DraftState) error
	Draft(string) (domain.DraftState, error)
	Workspace(string) (domain.Workspace, error)
}

type Service struct{ repo Repository }

func New(repo Repository) *Service { return &Service{repo: repo} }

func nextID(prefix string, count int) string { return prefix + "-" + strconv.Itoa(count+1) }

type CreateScriptInput struct {
	RequestKey string `json:"requestKey"`
	Title      string `json:"title"`
	Logline    string `json:"logline"`
	Genre      string `json:"genre"`
}

func (s *Service) CreateScript(in CreateScriptInput) (domain.Script, error) {
	items, err := s.repo.Scripts()
	if err != nil {
		return domain.Script{}, err
	}
	v := domain.Script{ID: nextID("script", len(items)), Title: strings.TrimSpace(in.Title), Logline: strings.TrimSpace(in.Logline), Genre: strings.TrimSpace(in.Genre), Status: domain.StatusIdea, Revision: 1}
	if err := domain.ValidateScript(v); err != nil {
		return domain.Script{}, err
	}
	if err := s.repo.SaveScript(v); err != nil {
		return domain.Script{}, err
	}
	draft := domain.DraftState{ScriptID: v.ID, State: domain.StatusIdea, Note: "created", Version: 1}
	if err := s.repo.SaveDraft(draft); err != nil {
		return domain.Script{}, err
	}
	return v, nil
}

func (s *Service) ReviseScript(id, title, logline, genre string) (domain.Script, error) {
	v, err := s.repo.Script(id)
	if err != nil {
		return domain.Script{}, err
	}
	if strings.TrimSpace(title) != "" {
		v.Title = strings.TrimSpace(title)
	}
	if strings.TrimSpace(logline) != "" {
		v.Logline = strings.TrimSpace(logline)
	}
	if strings.TrimSpace(genre) != "" {
		v.Genre = strings.TrimSpace(genre)
	}
	v.Revision++
	if err := domain.ValidateScript(v); err != nil {
		return domain.Script{}, err
	}
	return v, s.repo.SaveScript(v)
}

func (s *Service) Workspace(id string) (domain.Workspace, error) { return s.repo.Workspace(id) }

func (s *Service) ListScripts() ([]domain.Script, error) { return s.repo.Scripts() }

func (s *Service) requireScript(id string) error {
	_, err := s.repo.Script(id)
	if err != nil {
		return fmt.Errorf("script %s: %w", id, err)
	}
	return nil
}
