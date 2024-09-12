package authservice

import (
	"backend/src/internal/db"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID    uuid.UUID `db:"id"`
	Email string    `db:"email"`
}

type Session struct {
	TokenHash     string    `db:"token_hash"`
	UserID        uuid.UUID `db:"user_id"`
	AccessTokenID uuid.UUID `db:"access_token_id"`
	IssuedAt      time.Time `db:"issued_at"`
	ExpiresAt     time.Time `db:"expires_at"`
	IPAddress     string    `db:"ip_address"`
}

type AuthDB struct {
	DBInstance db.DatabaseInstance
}

func (db *AuthDB) GetUser(userID uuid.UUID) (User, error) {
	var u User
	query := `SELECT * FROM users WHERE id = $1`
	err := db.DBInstance.Conn.Get(&u, query, userID)

	return u, err
}

func (db *AuthDB) GetUserSession(userID uuid.UUID, accessTokenID uuid.UUID) (Session, error) {
	var userSession Session
	query := `SELECT * FROM user_sessions WHERE user_id = $1 AND access_token_id = $2`
	err := db.DBInstance.Conn.Get(&userSession, query, userID, accessTokenID)

	return userSession, err
}

func (db *AuthDB) SaveSession(tokenHash string, accessTokenID uuid.UUID, userID uuid.UUID, ip string) error {
	query := `INSERT INTO user_sessions (token_hash, user_id, access_token_id, issued_at, expires_at, ip_address) 
              VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.DBInstance.Conn.Exec(query,
		tokenHash,
		userID,
		accessTokenID,
		time.Now(),
		time.Now().Add(jwtGen.RefreshExpires),
		ip,
	)
	return err
}

func (db *AuthDB) UpdateSession(oldTokenID uuid.UUID, newTokenID uuid.UUID, newTokenHash string, ip string) error {
	query := `
		UPDATE user_sessions
		SET token_hash = $1, issued_at = $2, expires_at = $3, ip_address = $4, access_token_id = $5
		WHERE access_token_id = $6`

	_, err := db.DBInstance.Conn.Exec(query,
		newTokenHash,
		time.Now(),
		time.Now().Add(jwtGen.RefreshExpires),
		ip,
		newTokenID,
		oldTokenID,
	)
	return err
}
