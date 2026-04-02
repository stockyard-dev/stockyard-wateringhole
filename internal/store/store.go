package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Forum struct{
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Category string `json:"category"`
	PostCount int `json:"post_count"`
	Visibility string `json:"visibility"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"wateringhole.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS forums(id TEXT PRIMARY KEY,name TEXT NOT NULL,description TEXT DEFAULT '',category TEXT DEFAULT '',post_count INTEGER DEFAULT 0,visibility TEXT DEFAULT 'public',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Forum)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO forums(id,name,description,category,post_count,visibility,created_at)VALUES(?,?,?,?,?,?,?)`,e.ID,e.Name,e.Description,e.Category,e.PostCount,e.Visibility,e.CreatedAt);return err}
func(d *DB)Get(id string)*Forum{var e Forum;if d.db.QueryRow(`SELECT id,name,description,category,post_count,visibility,created_at FROM forums WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.Description,&e.Category,&e.PostCount,&e.Visibility,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Forum{rows,_:=d.db.Query(`SELECT id,name,description,category,post_count,visibility,created_at FROM forums ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Forum;for rows.Next(){var e Forum;rows.Scan(&e.ID,&e.Name,&e.Description,&e.Category,&e.PostCount,&e.Visibility,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM forums WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM forums`).Scan(&n);return n}
