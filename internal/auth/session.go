package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

var ErrInvalidSession = errors.New("invalid f1 session")

type Session struct {
	Token              string
	ExpiresAt          time.Time
	SubscriptionStatus string
}

func (s *Session) IsExpired(now time.Time) bool {
	return !s.ExpiresAt.After(now)
}

func (s *Session) HasActiveSubscription() bool {
	return strings.EqualFold(s.SubscriptionStatus, "active")
}

func sessionFromCookie(value string) (*Session, error) {
	decoded, err := url.PathUnescape(value)
	if err != nil {
		return nil, fmt.Errorf("%w: decoding cookie: %v", ErrInvalidSession, err)
	}

	var cookie struct {
		Data struct {
			SubscriptionToken string `json:"subscriptionToken"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(decoded), &cookie); err != nil {
		return nil, fmt.Errorf("%w: parsing cookie: %v", ErrInvalidSession, err)
	}
	if cookie.Data.SubscriptionToken == "" {
		return nil, fmt.Errorf("%w: cookie has no subscription token", ErrInvalidSession)
	}

	return sessionFromToken(cookie.Data.SubscriptionToken)
}

func sessionFromToken(token string) (*Session, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("%w: token should have 3 parts, got %d", ErrInvalidSession, len(parts))
	}

	payload, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(parts[1], "="))
	if err != nil {
		return nil, fmt.Errorf("%w: decoding token payload: %v", ErrInvalidSession, err)
	}

	var claims struct {
		SubscriptionStatus string `json:"SubscriptionStatus"`
		ExpiresAt          int64  `json:"exp"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("%w: parsing token payload: %v", ErrInvalidSession, err)
	}
	if claims.ExpiresAt == 0 {
		return nil, fmt.Errorf("%w: token has no expiry", ErrInvalidSession)
	}

	return &Session{
		Token:              token,
		ExpiresAt:          time.Unix(claims.ExpiresAt, 0),
		SubscriptionStatus: claims.SubscriptionStatus,
	}, nil
}
