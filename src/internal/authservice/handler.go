package authservice

import (
	"backend/src/internal/emailservice"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// base64, bcrypt???
func RefreshTokens(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	payload, err := jwtGen.VerifyToken(req.AccessToken, true)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid access token", http.StatusUnauthorized)
		return
	}

	userSession, err := DB.GetUserSession(payload.UserID, payload.UserID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to get user session", http.StatusUnauthorized)
		return
	}

	go func() {
		if userSession.IPAddress != r.RemoteAddr {
			user, err := DB.GetUser(payload.UserID)
			if err != nil {
				log.Println("Failed to get user data...")
			}

			emailservice.SendWarning(user.Email)
		}
	}()

	// TODO: Добавить логику обновления токена в бд
}

func IssueTokens(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := uuid.Parse(params["user_id"])

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ip := r.RemoteAddr

	tokenPair, err := jwtGen.CreateTokenPair(userID, ip)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	err = DB.SaveSession(string(tokenPair.RefreshToken), userID, userID, ip)

	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(tokenPair)
	if err != nil {
		log.Println(err)
		http.Error(w, "JSON Marshal failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
