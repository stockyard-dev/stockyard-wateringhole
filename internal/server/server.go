package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"github.com/stockyard-dev/stockyard-wateringhole/internal/store"
)

type Server struct { db *store.DB; mux *http.ServeMux; port int; limits Limits }

func New(db *store.DB, port int, limits Limits) *Server {
	s := &Server{db: db, mux: http.NewServeMux(), port: port, limits: limits}
	s.mux.HandleFunc("POST /api/profiles", s.hCreateProfile)
	s.mux.HandleFunc("GET /api/profiles", s.hListProfiles)
	s.mux.HandleFunc("GET /api/profiles/{id}", s.hGetProfile)
	s.mux.HandleFunc("DELETE /api/profiles/{id}", s.hDelProfile)
	s.mux.HandleFunc("POST /api/profiles/{id}/links", s.hCreateLink)
	s.mux.HandleFunc("GET /api/profiles/{id}/links", s.hListLinks)
	s.mux.HandleFunc("DELETE /api/links/{id}", s.hDelLink)
	s.mux.HandleFunc("GET /api/links/{id}/click", s.hClick)
	s.mux.HandleFunc("GET /p/{slug}", s.hPublicPage)
	s.mux.HandleFunc("GET /api/status", func(w http.ResponseWriter, r *http.Request) { wj(w, 200, s.db.Stats()) })
	s.mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) { wj(w, 200, map[string]string{"status": "ok"}) })
	s.mux.HandleFunc("GET /ui", s.handleUI)
	s.mux.HandleFunc("GET /api/version", func(w http.ResponseWriter, r *http.Request) { wj(w, 200, map[string]any{"product": "stockyard-wateringhole", "version": "0.1.0"}) })
	return s
}

func (s *Server) Start() error { log.Printf("[wateringhole] :%d", s.port); return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.mux) }

func (s *Server) hCreateProfile(w http.ResponseWriter, r *http.Request) {
	var req struct { Slug string `json:"slug"`; Name string `json:"name"`; Bio string `json:"bio"` }
	if json.NewDecoder(r.Body).Decode(&req) != nil || req.Slug == "" || req.Name == "" { wj(w, 400, map[string]string{"error": "slug and name required"}); return }
	p, err := s.db.CreateProfile(req.Slug, req.Name, req.Bio)
	if err != nil { wj(w, 500, map[string]string{"error": err.Error()}); return }
	wj(w, 201, map[string]any{"profile": p, "url": fmt.Sprintf("/p/%s", p.Slug)})
}
func (s *Server) hListProfiles(w http.ResponseWriter, r *http.Request) { ps, _ := s.db.ListProfiles(); if ps == nil { ps = []store.Profile{} }; wj(w, 200, map[string]any{"profiles": ps, "count": len(ps)}) }
func (s *Server) hGetProfile(w http.ResponseWriter, r *http.Request) {
	p, err := s.db.GetProfile(r.PathValue("id")); if err != nil { wj(w, 404, map[string]string{"error": "not found"}); return }
	links, _ := s.db.ListLinks(p.ID); if links == nil { links = []store.Link{} }
	wj(w, 200, map[string]any{"profile": p, "links": links})
}
func (s *Server) hDelProfile(w http.ResponseWriter, r *http.Request) { s.db.DeleteProfile(r.PathValue("id")); wj(w, 200, map[string]string{"status": "deleted"}) }

func (s *Server) hCreateLink(w http.ResponseWriter, r *http.Request) {
	pid := r.PathValue("id")
	var req struct { Title string `json:"title"`; URL string `json:"url"`; Icon string `json:"icon"` }
	if json.NewDecoder(r.Body).Decode(&req) != nil || req.Title == "" || req.URL == "" { wj(w, 400, map[string]string{"error": "title and url required"}); return }
	l, err := s.db.CreateLink(pid, req.Title, req.URL, req.Icon)
	if err != nil { wj(w, 500, map[string]string{"error": err.Error()}); return }
	wj(w, 201, map[string]any{"link": l})
}
func (s *Server) hListLinks(w http.ResponseWriter, r *http.Request) { ls, _ := s.db.ListLinks(r.PathValue("id")); if ls == nil { ls = []store.Link{} }; wj(w, 200, map[string]any{"links": ls, "count": len(ls)}) }
func (s *Server) hDelLink(w http.ResponseWriter, r *http.Request) { s.db.DeleteLink(r.PathValue("id")); wj(w, 200, map[string]string{"status": "deleted"}) }
func (s *Server) hClick(w http.ResponseWriter, r *http.Request) { s.db.ClickLink(r.PathValue("id")); wj(w, 200, map[string]string{"status": "clicked"}) }

func (s *Server) hPublicPage(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	p, err := s.db.GetProfileBySlug(slug)
	if err != nil { http.NotFound(w, r); return }
	links, _ := s.db.ListLinks(p.ID)
	var linksHTML strings.Builder
	for _, l := range links {
		linksHTML.WriteString(fmt.Sprintf(`<a href="%s" target="_blank" onclick="fetch('/api/links/%s/click')" style="display:block;background:#241e18;border:1px solid #2e261e;padding:1rem;margin-bottom:.6rem;text-decoration:none;color:#f0e6d3;font-family:'JetBrains Mono',monospace;font-size:.85rem;text-align:center;transition:border-color .2s" onmouseover="this.style.borderColor='#e8753a'" onmouseout="this.style.borderColor='#2e261e'">%s</a>`,
			he(l.URL), l.ID, he(l.Title)))
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>@%s</title>
<link href="https://fonts.googleapis.com/css2?family=Libre+Baskerville:wght@400;700&family=JetBrains+Mono:wght@400;600&display=swap" rel="stylesheet">
<style>body{background:#1a1410;color:#f0e6d3;font-family:'Libre Baskerville',serif;margin:0;min-height:100vh;display:flex;justify-content:center;padding:3rem 1rem}
.container{max-width:400px;width:100%%}.name{font-size:1.5rem;text-align:center;margin-bottom:.3rem}.bio{text-align:center;font-size:.85rem;color:#bfb5a3;margin-bottom:2rem}
.footer{text-align:center;margin-top:2rem;font-size:.5rem;color:#7a7060;font-family:'JetBrains Mono',monospace}
.footer a{color:#e8753a;text-decoration:none}</style></head><body>
<div class="container"><div class="name">%s</div><div class="bio">%s</div>%s
<div class="footer">Powered by <a href="https://stockyard.dev">Stockyard</a></div></div></body></html>`,
		he(p.Slug), he(p.Name), he(p.Bio), linksHTML.String())
}

func he(s string) string { s = strings.ReplaceAll(s, "&", "&amp;"); s = strings.ReplaceAll(s, "<", "&lt;"); s = strings.ReplaceAll(s, ">", "&gt;"); s = strings.ReplaceAll(s, "\"", "&quot;"); return s }
func wj(w http.ResponseWriter, code int, v any) { w.Header().Set("Content-Type", "application/json"); w.WriteHeader(code); json.NewEncoder(w).Encode(v) }
