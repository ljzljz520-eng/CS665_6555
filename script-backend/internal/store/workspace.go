package store

import (
	"go.etcd.io/bbolt"
	"scriptstudio/script-backend/internal/domain"
)

func (s *Store) Workspace(scriptID string) (domain.Workspace, error) {
	script, err := s.Script(scriptID)
	if err != nil {
		return domain.Workspace{}, err
	}
	draft, err := s.Draft(scriptID)
	if err != nil {
		return domain.Workspace{}, err
	}
	allScenes, err := s.Scenes()
	if err != nil {
		return domain.Workspace{}, err
	}
	allCharacters, err := s.Characters()
	if err != nil {
		return domain.Workspace{}, err
	}
	allDialogues, err := s.Dialogues()
	if err != nil {
		return domain.Workspace{}, err
	}
	w := domain.Workspace{Script: script, Draft: draft}
	sceneIDs := make(map[string]bool)
	for _, scene := range allScenes {
		if scene.ScriptID == scriptID {
			w.Scenes = append(w.Scenes, scene)
			sceneIDs[scene.ID] = true
		}
	}
	for _, character := range allCharacters {
		if character.ScriptID == scriptID {
			w.Characters = append(w.Characters, character)
		}
	}
	for _, line := range allDialogues {
		if sceneIDs[line.SceneID] {
			w.Dialogues = append(w.Dialogues, line)
		}
	}
	w.Scenes = domain.SortScenes(w.Scenes)
	w.Dialogues = domain.SortDialogues(w.Dialogues)
	return w, nil
}

func (s *Store) DeleteScript(id string) error {
	w, err := s.Workspace(id)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, scene := range w.Scenes {
			for _, line := range w.Dialogues {
				if line.SceneID == scene.ID {
					if err := tx.Bucket(bDialogues).Delete([]byte(line.ID)); err != nil {
						return err
					}
				}
			}
			if err := tx.Bucket(bScenes).Delete([]byte(scene.ID)); err != nil {
				return err
			}
		}
		for _, character := range w.Characters {
			if err := tx.Bucket(bCharacters).Delete([]byte(character.ID)); err != nil {
				return err
			}
		}
		if err := tx.Bucket(bDrafts).Delete([]byte(id)); err != nil {
			return err
		}
		return tx.Bucket(bScripts).Delete([]byte(id))
	})
}
