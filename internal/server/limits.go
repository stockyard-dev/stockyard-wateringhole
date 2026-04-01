package server

import "github.com/stockyard-dev/stockyard-wateringhole/internal/license"

type Limits struct { MaxProfiles int; MaxLinks int }
var freeLimits = Limits{MaxProfiles: 1, MaxLinks: 10}
var proLimits = Limits{MaxProfiles: 0, MaxLinks: 0}
func LimitsFor(info *license.Info) Limits { if info != nil && info.IsPro() { return proLimits }; return freeLimits }
func LimitReached(l, c int) bool { return l > 0 && c >= l }
