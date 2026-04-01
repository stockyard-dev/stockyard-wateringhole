package server

import "net/http"

func (s *Server) handleUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`<!DOCTYPE html><html><head><meta charset="UTF-8"><title>Watering Hole — Stockyard</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;600&display=swap" rel="stylesheet">
<style>*{margin:0;padding:0;box-sizing:border-box}body{background:#1a1410;color:#f0e6d3;font-family:'JetBrains Mono',monospace;padding:2rem}
.hdr{font-size:.7rem;color:#a0845c;letter-spacing:3px;text-transform:uppercase;margin-bottom:2rem;border-bottom:2px solid #8b3d1a;padding-bottom:.8rem}
.section{margin-bottom:2rem}.section h2{font-size:.65rem;letter-spacing:3px;text-transform:uppercase;color:#e8753a;margin-bottom:.8rem}
.item{background:#241e18;padding:.6rem .8rem;margin-bottom:.4rem;border:1px solid #2e261e;font-size:.72rem}
.empty{color:#7a7060;font-style:italic;padding:1rem;text-align:center}
</style></head><body>
<div class="hdr">Stockyard · Watering Hole</div>
<div class="section"><h2>Profiles</h2><div id="list"></div></div>
<script>
async function refresh(){const d=await(await fetch('/api/profiles')).json();const ps=d.profiles||[];
document.getElementById('list').innerHTML=ps.length?ps.map(p=>'<div class="item"><a href="/@'+p.slug+'" style="color:#e8753a;text-decoration:none">@'+p.slug+'</a> — '+p.name+'</div>').join(''):'<div class="empty">No profiles</div>';}
refresh();setInterval(refresh,8000);
</script></body></html>`))
}
