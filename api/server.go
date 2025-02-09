package api

import (
	"db/database"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Server struct {
	db     *database.Database
	router *mux.Router
}

func NewServer(db *database.Database) *Server {
	s := &Server{
		db:     db,
		router: mux.NewRouter(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.HandleFunc("/tables", s.handleCreateTable).Methods("POST")
	s.router.HandleFunc("/tables/{name}/rows", s.handleInsert).Methods("POST")
	s.router.HandleFunc("/tables/{name}/rows", s.handleSelect).Methods("GET")
	s.router.HandleFunc("/tables/{name}/rows", s.handleUpdate).Methods("PUT")
	s.router.HandleFunc("/tables/{name}/rows", s.handleDelete).Methods("DELETE")
}

func (s *Server) Run(addr string) error {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type"},
	})

	handler := c.Handler(s.router)
	return http.ListenAndServe(addr, handler)
}

func (s *Server) handleCreateTable(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name    string            `json:"name"`
		Columns []database.Column `json:"columns"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.db.CreateTable(req.Name, req.Columns); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) handleInsert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tableName := vars["name"]

	table, err := s.db.GetTable(tableName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var row database.Row
	if err := json.NewDecoder(r.Body).Decode(&row); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := table.InsertRow(row); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) handleSelect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tableName := vars["name"]

	table, err := s.db.GetTable(tableName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	whereParam := r.URL.Query().Get("where")
	var conditions map[string]interface{}

	if whereParam != "" {
		if err := json.Unmarshal([]byte(whereParam), &conditions); err != nil {
			http.Error(w, "Invalid where clause", http.StatusBadRequest)
			return
		}
	}

	rows, err := table.Select(conditions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rows)
}

func (s *Server) handleUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tableName := vars["name"]

	table, err := s.db.GetTable(tableName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var req struct {
		Where   map[string]interface{} `json:"where"`
		Updates map[string]interface{} `json:"updates"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updated, err := table.Update(req.Where, req.Updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"updated": updated})
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tableName := vars["name"]

	table, err := s.db.GetTable(tableName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	whereParam := r.URL.Query().Get("where")
	var conditions map[string]interface{}

	if whereParam != "" {
		if err := json.Unmarshal([]byte(whereParam), &conditions); err != nil {
			http.Error(w, "Invalid where clause", http.StatusBadRequest)
			return
		}
	}

	deleted, err := table.Delete(conditions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"deleted": deleted})
}
