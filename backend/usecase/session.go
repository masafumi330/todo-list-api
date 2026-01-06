package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Session struct {
	SessionID string    `json:"session_id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type SessionUsecase struct {
	client    *redis.Client
	ttl       time.Duration
	keyPrefix string
}

func NewSessionUsecase(client *redis.Client, ttl time.Duration) *SessionUsecase {
	return &SessionUsecase{
		client:    client,
		ttl:       ttl,
		keyPrefix: "session:",
	}
}

func (uc *SessionUsecase) Create(ctx context.Context, userID int) (string, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return "", err
	}

	now := time.Now().UTC()
	session := Session{
		SessionID: sessionID,
		UserID:    userID,
		CreatedAt: now,
		ExpiresAt: now.Add(uc.ttl),
	}

	payload, err := json.Marshal(session)
	if err != nil {
		return "", err
	}

	if err := uc.client.Set(ctx, uc.key(sessionID), payload, uc.ttl).Err(); err != nil {
		return "", err
	}

	return sessionID, nil
}

func (uc *SessionUsecase) key(id string) string {
	return fmt.Sprintf("%s%s", uc.keyPrefix, id)
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
