package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/redis/rueidis"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

func Logging() Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			slog.Info(r.RemoteAddr, r.Method, r.URL.Path)

			next(w, r)
		}
	}
}

func Cache(s *Server) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			url, _ := url.Parse(s.target)
			reqUrl := fmt.Sprintf("%s:%s", url.Host, r.URL.RequestURI())
			data, err := s.db.GetString(reqUrl)

			if err == rueidis.Nil {
				next(w, r)
				return
			} else if err != nil {
				slog.Error("cache error :", "s.db.GetString()", err)
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}

			var res interface{}
			err = json.Unmarshal([]byte(data), &res)
			if err != nil {
				slog.Error("error parsing json", "json.Unmarshal()", err)
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Add("X-Cache", "HIT")
			json.NewEncoder(w).Encode(res)

		}
	}
}
