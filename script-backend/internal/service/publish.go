package service

import (
	"fmt"
	"scriptstudio/script-backend/internal/domain"
	"strings"
)

func (s *Service) TransitionDraft(scriptID string, to domain.DraftStatus, note string) (domain.DraftState, error) {
	current, err := s.repo.Draft(scriptID)
	if err != nil {
		return domain.DraftState{}, err
	}
	next, err := domain.Transition(current, to, strings.TrimSpace(note))
	if err != nil {
		return domain.DraftState{}, err
	}
	if to == domain.StatusReview {
		w, err := s.repo.Workspace(scriptID)
		if err != nil {
			return domain.DraftState{}, err
		}
		if issues := domain.PublicationIssues(w); len(issues) > 0 {
			return domain.DraftState{}, fmt.Errorf("cannot review: %s", strings.Join(issues, "; "))
		}
	}
	if err := s.repo.SaveDraft(next); err != nil {
		return domain.DraftState{}, err
	}
	script, err := s.repo.Script(scriptID)
	if err != nil {
		return domain.DraftState{}, err
	}
	script.Status = to
	script.Revision++
	if err := s.repo.SaveScript(script); err != nil {
		return domain.DraftState{}, err
	}
	return next, nil
}

func (s *Service) Publish(scriptID, note string) (domain.DraftState, error) {
	current, err := s.repo.Draft(scriptID)
	if err != nil {
		return domain.DraftState{}, err
	}
	if current.State == domain.StatusDraft {
		current, err = s.TransitionDraft(scriptID, domain.StatusReview, "ready for review")
		if err != nil {
			return domain.DraftState{}, err
		}
	}
	if current.State != domain.StatusReview {
		return domain.DraftState{}, fmt.Errorf("script is %s, not in review", current.State)
	}
	return s.TransitionDraft(scriptID, domain.StatusPublished, note)
}
