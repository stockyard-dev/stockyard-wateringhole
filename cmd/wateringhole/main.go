package main

import (
	"fmt"; "log"; "os"; "strconv"
	"github.com/stockyard-dev/stockyard-wateringhole/internal/license"
	"github.com/stockyard-dev/stockyard-wateringhole/internal/server"
	"github.com/stockyard-dev/stockyard-wateringhole/internal/store"
)
var version = "dev"
func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") { fmt.Printf("wateringhole %s\n", version); os.Exit(0) }
	if len(os.Args) > 1 && os.Args[1] == "--health" { fmt.Println("ok"); os.Exit(0) }
	log.SetFlags(log.Ltime | log.Lshortfile)
	port := 9000; if p := os.Getenv("PORT"); p != "" { if n, e := strconv.Atoi(p); e == nil { port = n } }
	dataDir := os.Getenv("DATA_DIR"); if dataDir == "" { dataDir = "./data" }
	lk := os.Getenv("WATERINGHOLE_LICENSE_KEY"); li, le := license.Validate(lk, "wateringhole")
	if lk != "" && le != nil { log.Printf("[license] %v", le); li = nil }
	limits := server.LimitsFor(li)
	db, err := store.Open(dataDir); if err != nil { log.Fatalf("db: %v", err) }; defer db.Close()
	log.Printf("  Stockyard Watering Hole %s — http://localhost:%d/ui", version, port)
	srv := server.New(db, port, limits); if err := srv.Start(); err != nil { log.Fatalf("server: %v", err) }
}
