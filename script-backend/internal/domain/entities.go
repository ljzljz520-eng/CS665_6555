package domain

type DraftStatus string

const (
	StatusIdea      DraftStatus = "idea"
	StatusDraft     DraftStatus = "draft"
	StatusReview    DraftStatus = "review"
	StatusPublished DraftStatus = "published"
)

type Script struct {
	ID       string      `json:"id"`
	Title    string      `json:"title"`
	Logline  string      `json:"logline"`
	Genre    string      `json:"genre"`
	Status   DraftStatus `json:"status"`
	Revision int         `json:"revision"`
}

type Scene struct {
	ID        string `json:"id"`
	ScriptID  string `json:"scriptId"`
	Heading   string `json:"heading"`
	Synopsis  string `json:"synopsis"`
	Location  string `json:"location"`
	TimeOfDay string `json:"timeOfDay"`
	Position  int    `json:"position"`
}

type Character struct {
	ID        string `json:"id"`
	ScriptID  string `json:"scriptId"`
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	Objective string `json:"objective"`
}

type Dialogue struct {
	ID          string `json:"id"`
	SceneID     string `json:"sceneId"`
	CharacterID string `json:"characterId"`
	Text        string `json:"text"`
	Direction   string `json:"direction"`
	Position    int    `json:"position"`
}

type DraftState struct {
	ScriptID string      `json:"scriptId"`
	State    DraftStatus `json:"state"`
	Note     string      `json:"note"`
	Version  int         `json:"version"`
}

type Workspace struct {
	Script     Script      `json:"script"`
	Scenes     []Scene     `json:"scenes"`
	Characters []Character `json:"characters"`
	Dialogues  []Dialogue  `json:"dialogues"`
	Draft      DraftState  `json:"draft"`
}
