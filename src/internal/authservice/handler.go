package authservice

import (
	"backend/src/internal/emailservice"
	"encoding/json"
	"log"
	"net/http"
	"strings"

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

	payload, err := jwtGen.VerifyToken(req.RefreshToken, true)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	_, err = DB.GetUserSession(payload.UserID, payload.TokenUUID)

	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	ip := strings.Split(r.RemoteAddr, ":")[0]

	go func() {
		if payload.IPAddress != ip {
			user, err := DB.GetUser(payload.UserID)
			if err != nil {
				log.Println("Failed to get user data...")
			}

			emailservice.SendWarning(user.Email)
		}
	}()

	tokens, err := jwtGen.CreateTokenPair(payload.UserID, ip)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	err = DB.UpdateSession(payload.TokenUUID, tokens.TokenUUID, tokens.Pair.RefreshToken, ip)
	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	Response(w, tokens.Pair)
}

func IssueTokens(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := uuid.Parse(params["user_id"])

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ip := strings.Split(r.RemoteAddr, ":")[0]

	tokens, err := jwtGen.CreateTokenPair(userID, ip)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	err = DB.SaveSession(tokens.Pair.RefreshToken, tokens.TokenUUID, userID, ip)

	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	Response(w, tokens.Pair)
}

func Response(w http.ResponseWriter, data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		http.Error(w, "JSON Marshal failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
