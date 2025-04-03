package main

import (
	"github.com/danilkompaniets/go-rss/internal/auth"
	"github.com/danilkompaniets/go-rss/internal/database"
	"net/http"
)

type authedHandler func(w http.ResponseWriter, r *http.Request, user database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		apiKey, err := auth.GetApiKey(request.Header)

		if err != nil {
			respondWithError(writer, 400, "Invalid API key "+err.Error())
			return
		}

		user, err := apiCfg.DB.GetUserByApiKey(request.Context(), apiKey)

		if err != nil {
			respondWithError(writer, 404, "User not found")
			return
		}

		handler(writer, request, user)
	}
}
