package store

import (
	"encoding/json"
	"errors"
	"fmt"

	"go.etcd.io/bbolt"
	"scriptstudio/script-backend/internal/domain"
)

var (
	bScripts    = []byte("scripts")
	bScenes     = []byte("scenes")
	bCharacters = []byte("characters")
	bDialogues  = []byte("dialogues")
	bDrafts     = []byte("drafts")
)

type Store struct{ db *bbolt.DB }

func Open(path string) (*Store, error) {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		for _, name := range [][]byte{bScripts, bScenes, bCharacters, bDialogues, bDrafts} {
			if _, e := tx.CreateBucketIfNotExists(name); e != nil {
				return e
			}
		}
		return nil
	})
	if err != nil {
		db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

func put[T any](s *Store, bucket []byte, id string, value T) error {
	if id == "" {
		return errors.New("empty persistence key")
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucket).Put([]byte(id), data) })
}

func get[T any](s *Store, bucket []byte, id string) (T, error) {
	var out T
	err := s.db.View(func(tx *bbolt.Tx) error {
		data := tx.Bucket(bucket).Get([]byte(id))
		if data == nil {
			return fmt.Errorf("record %s not found", id)
		}
		return json.Unmarshal(data, &out)
	})
	return out, err
}

func list[T any](s *Store, bucket []byte) ([]T, error) {
	out := make([]T, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucket).ForEach(func(_, data []byte) error {
			var item T
			if err := json.Unmarshal(data, &item); err != nil {
				return err
			}
			out = append(out, item)
			return nil
		})
	})
	return out, err
}

func (s *Store) SaveScript(v domain.Script) error        { return put(s, bScripts, v.ID, v) }
func (s *Store) Script(id string) (domain.Script, error) { return get[domain.Script](s, bScripts, id) }
func (s *Store) Scripts() ([]domain.Script, error)       { return list[domain.Script](s, bScripts) }
func (s *Store) SaveScene(v domain.Scene) error          { return put(s, bScenes, v.ID, v) }
func (s *Store) Scenes() ([]domain.Scene, error)         { return list[domain.Scene](s, bScenes) }
func (s *Store) SaveCharacter(v domain.Character) error  { return put(s, bCharacters, v.ID, v) }
func (s *Store) Characters() ([]domain.Character, error) {
	return list[domain.Character](s, bCharacters)
}
func (s *Store) SaveDialogue(v domain.Dialogue) error  { return put(s, bDialogues, v.ID, v) }
func (s *Store) Dialogues() ([]domain.Dialogue, error) { return list[domain.Dialogue](s, bDialogues) }
func (s *Store) SaveDraft(v domain.DraftState) error   { return put(s, bDrafts, v.ScriptID, v) }
func (s *Store) Draft(id string) (domain.DraftState, error) {
	return get[domain.DraftState](s, bDrafts, id)
}
