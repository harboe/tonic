package logging

import (
	"net/http"
	"regexp"
)

var regExt = regexp.MustCompile("(.*/)?(favicon.([iI][cC][oO]|[gG][iI][fF])|apple-touch-icon)(/.*)?")

func IsFavicon(req *http.Request) bool {
	if m := regExt.FindString(req.URL.Path); len(m) > 0 {
		return true
	}

	return false
}
