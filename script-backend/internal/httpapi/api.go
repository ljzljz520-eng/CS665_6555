package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"scriptstudio/script-backend/internal/catalog"
	"scriptstudio/script-backend/internal/preview"
	"scriptstudio/script-backend/internal/query"
	"scriptstudio/script-backend/internal/service"
)

type API struct {
	commands *service.Service
	queries  *query.Query
	mux      *http.ServeMux
}

func New(commands *service.Service, queries *query.Query) *API {
	a := &API{commands: commands, queries: queries, mux: http.NewServeMux()}
	a.routes()
	return a
}

func (a *API) Handler() http.Handler { return a.mux }

func (a *API) routes() {
	a.mux.HandleFunc("GET /api/scripts", a.listScripts)
	a.mux.HandleFunc("POST /api/scripts", a.createScript)
	a.mux.HandleFunc("GET /api/scripts/{id}", a.workspace)
	a.mux.HandleFunc("POST /api/scripts/{id}/scenes", a.addScene)
	a.mux.HandleFunc("POST /api/scripts/{id}/characters", a.addCharacter)
	a.mux.HandleFunc("POST /api/scripts/{id}/publish", a.publish)
	a.mux.HandleFunc("GET /api/scripts/{id}/preview", a.renderPreview)
	a.mux.HandleFunc("GET /api/catalog/blueprints", a.blueprints)
}

func write(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(value)
}
func problem(w http.ResponseWriter, err error) {
	write(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
}

func (a *API) listScripts(w http.ResponseWriter, r *http.Request) {
	items, err := a.queries.Scripts(query.ScriptFilter{Text: r.URL.Query().Get("q"), Genre: r.URL.Query().Get("genre")})
	if err != nil {
		problem(w, err)
		return
	}
	write(w, http.StatusOK, items)
}

func (a *API) createScript(w http.ResponseWriter, r *http.Request) {
	var in service.CreateScriptInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		problem(w, err)
		return
	}
	item, err := a.commands.CreateScript(in)
	if err != nil {
		problem(w, err)
		return
	}
	write(w, http.StatusCreated, item)
}

func (a *API) workspace(w http.ResponseWriter, r *http.Request) {
	item, err := a.commands.Workspace(r.PathValue("id"))
	if err != nil {
		problem(w, err)
		return
	}
	write(w, http.StatusOK, item)
}

func (a *API) addScene(w http.ResponseWriter, r *http.Request) {
	var in service.AddSceneInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		problem(w, err)
		return
	}
	in.ScriptID = r.PathValue("id")
	item, err := a.commands.AddScene(in)
	if err != nil {
		problem(w, err)
		return
	}
	write(w, http.StatusCreated, item)
}

func (a *API) addCharacter(w http.ResponseWriter, r *http.Request) {
	var raw struct{ Name, Bio, Objective string }
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		problem(w, err)
		return
	}
	item, err := a.commands.AddCharacter(service.AddCharacterInput{ScriptID: r.PathValue("id"), Name: raw.Name, Bio: raw.Bio, Objective: raw.Objective})
	if err != nil {
		problem(w, err)
		return
	}
	write(w, http.StatusCreated, item)
}

func (a *API) publish(w http.ResponseWriter, r *http.Request) {
	item, err := a.commands.Publish(r.PathValue("id"), "published by author")
	if err != nil {
		problem(w, err)
		return
	}
	write(w, http.StatusOK, item)
}

func (a *API) renderPreview(w http.ResponseWriter, r *http.Request) {
	item, err := a.commands.Workspace(r.PathValue("id"))
	if err != nil {
		problem(w, err)
		return
	}
	write(w, http.StatusOK, preview.Render(item))
}

func (a *API) blueprints(w http.ResponseWriter, r *http.Request) {
	term := strings.TrimSpace(r.URL.Query().Get("q"))
	genre := strings.TrimSpace(r.URL.Query().Get("genre"))
	if term != "" {
		write(w, http.StatusOK, catalog.Search(term))
		return
	}
	if genre != "" {
		write(w, http.StatusOK, catalog.ForGenre(genre))
		return
	}
	write(w, http.StatusOK, catalog.Blueprints())
}
