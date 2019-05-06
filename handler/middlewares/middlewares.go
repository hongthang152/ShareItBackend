package middlewares

import (
	"github.com/hongthang152/ShareItBackend/utils"
	"log"
	"net/http"
	"os"
)

func SetCORSPolicy(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		h.ServeHTTP(w, r)
	})
}

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, "Error from server")
				log.Println(err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
