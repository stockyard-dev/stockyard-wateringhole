package store

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"
	_ "modernc.org/sqlite"
)

type DB struct{ conn *sql.DB }

func Open(dataDir string) (*DB, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil { return nil, fmt.Errorf("create data dir: %w", err) }
	conn, err := sql.Open("sqlite", filepath.Join(dataDir, "wateringhole.db"))
	if err != nil { return nil, err }
	conn.Exec("PRAGMA journal_mode=WAL")
	conn.Exec("PRAGMA busy_timeout=5000")
	conn.SetMaxOpenConns(4)
	db := &DB{conn: conn}
	return db, db.migrate()
}

func (db *DB) Close() error { return db.conn.Close() }

func (db *DB) migrate() error {
	_, err := db.conn.Exec(`
CREATE TABLE IF NOT EXISTS profiles (
    id TEXT PRIMARY KEY, slug TEXT NOT NULL UNIQUE, name TEXT NOT NULL,
    bio TEXT DEFAULT '', avatar_url TEXT DEFAULT '', theme TEXT DEFAULT 'dark',
    created_at TEXT DEFAULT (datetime('now'))
);
CREATE TABLE IF NOT EXISTS links (
    id TEXT PRIMARY KEY, profile_id TEXT NOT NULL, title TEXT NOT NULL,
    url TEXT NOT NULL, icon TEXT DEFAULT '', sort_order INTEGER DEFAULT 0,
    enabled INTEGER DEFAULT 1, clicks INTEGER DEFAULT 0,
    created_at TEXT DEFAULT (datetime('now'))
);
CREATE INDEX IF NOT EXISTS idx_links_profile ON links(profile_id);
`)
	return err
}

type Profile struct {
	ID        string `json:"id"`
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
	Theme     string `json:"theme"`
	CreatedAt string `json:"created_at"`
}

func (db *DB) CreateProfile(slug, name, bio string) (*Profile, error) {
	id := "prof_" + genID(6)
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := db.conn.Exec("INSERT INTO profiles (id,slug,name,bio,created_at) VALUES (?,?,?,?,?)", id, slug, name, bio, now)
	if err != nil { return nil, err }
	return &Profile{ID: id, Slug: slug, Name: name, Bio: bio, Theme: "dark", CreatedAt: now}, nil
}

func (db *DB) ListProfiles() ([]Profile, error) {
	rows, err := db.conn.Query("SELECT id,slug,name,bio,avatar_url,theme,created_at FROM profiles ORDER BY created_at DESC")
	if err != nil { return nil, err }
	defer rows.Close()
	var out []Profile
	for rows.Next() { var p Profile; rows.Scan(&p.ID, &p.Slug, &p.Name, &p.Bio, &p.AvatarURL, &p.Theme, &p.CreatedAt); out = append(out, p) }
	return out, rows.Err()
}

func (db *DB) GetProfile(id string) (*Profile, error) {
	var p Profile
	err := db.conn.QueryRow("SELECT id,slug,name,bio,avatar_url,theme,created_at FROM profiles WHERE id=?", id).
		Scan(&p.ID, &p.Slug, &p.Name, &p.Bio, &p.AvatarURL, &p.Theme, &p.CreatedAt)
	return &p, err
}

func (db *DB) GetProfileBySlug(slug string) (*Profile, error) {
	var p Profile
	err := db.conn.QueryRow("SELECT id,slug,name,bio,avatar_url,theme,created_at FROM profiles WHERE slug=?", slug).
		Scan(&p.ID, &p.Slug, &p.Name, &p.Bio, &p.AvatarURL, &p.Theme, &p.CreatedAt)
	return &p, err
}

func (db *DB) DeleteProfile(id string) { db.conn.Exec("DELETE FROM links WHERE profile_id=?", id); db.conn.Exec("DELETE FROM profiles WHERE id=?", id) }

type Link struct {
	ID        string `json:"id"`
	ProfileID string `json:"profile_id"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	Icon      string `json:"icon"`
	SortOrder int    `json:"sort_order"`
	Enabled   bool   `json:"enabled"`
	Clicks    int    `json:"clicks"`
	CreatedAt string `json:"created_at"`
}

func (db *DB) CreateLink(profileID, title, url, icon string) (*Link, error) {
	id := "lnk_" + genID(6)
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := db.conn.Exec("INSERT INTO links (id,profile_id,title,url,icon,created_at) VALUES (?,?,?,?,?,?)", id, profileID, title, url, icon, now)
	if err != nil { return nil, err }
	return &Link{ID: id, ProfileID: profileID, Title: title, URL: url, Icon: icon, Enabled: true, CreatedAt: now}, nil
}

func (db *DB) ListLinks(profileID string) ([]Link, error) {
	rows, err := db.conn.Query("SELECT id,profile_id,title,url,icon,sort_order,enabled,clicks,created_at FROM links WHERE profile_id=? AND enabled=1 ORDER BY sort_order, created_at", profileID)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []Link
	for rows.Next() { var l Link; var en int; rows.Scan(&l.ID, &l.ProfileID, &l.Title, &l.URL, &l.Icon, &l.SortOrder, &en, &l.Clicks, &l.CreatedAt); l.Enabled = en == 1; out = append(out, l) }
	return out, rows.Err()
}

func (db *DB) ClickLink(id string) { db.conn.Exec("UPDATE links SET clicks=clicks+1 WHERE id=?", id) }
func (db *DB) DeleteLink(id string) { db.conn.Exec("DELETE FROM links WHERE id=?", id) }

func (db *DB) Stats() map[string]any {
	var profiles, links, clicks int
	db.conn.QueryRow("SELECT COUNT(*) FROM profiles").Scan(&profiles)
	db.conn.QueryRow("SELECT COUNT(*) FROM links").Scan(&links)
	db.conn.QueryRow("SELECT COALESCE(SUM(clicks),0) FROM links").Scan(&clicks)
	return map[string]any{"profiles": profiles, "links": links, "total_clicks": clicks}
}

func genID(n int) string { b := make([]byte, n); rand.Read(b); return hex.EncodeToString(b) }
