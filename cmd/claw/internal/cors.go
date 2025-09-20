package internal

import "net/http"

func corsDevMidddlware(enable bool) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		if !enable {
			return h
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if host := r.Header.Get("Origin"); host != "" {
				w.Header().Set("Access-Control-Allow-Origin", host)
			} else {
				w.Header().Set("Access-Control-Allow-Origin", "http://"+r.Host)
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			h.ServeHTTP(w, r)
		})
	}
}
