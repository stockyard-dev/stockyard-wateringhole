package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-wateringhole/internal/server";"github.com/stockyard-dev/stockyard-wateringhole/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="9050"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./wateringhole-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("wateringhole: %v",err)};defer db.Close();srv:=server.New(db)
fmt.Printf("\n  Watering Hole — community forum\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n\n",port,port)
log.Printf("wateringhole: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
