package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PushService struct {
	DB              *pgxpool.Pool
	VAPIDPublicKey  string
	VAPIDPrivateKey string
	VAPIDContact    string // mailto: or URL
}

func NewPushService(db *pgxpool.Pool, publicKey, privateKey, contact string) *PushService {
	return &PushService{
		DB:              db,
		VAPIDPublicKey:  publicKey,
		VAPIDPrivateKey: privateKey,
		VAPIDContact:    contact,
	}
}

type PushPayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Icon  string `json:"icon,omitempty"`
	URL   string `json:"url,omitempty"`
}

// SaveSubscription stores a push subscription for a user
func (ps *PushService) SaveSubscription(ctx context.Context, userID, endpoint, p256dh, auth string) error {
	_, err := ps.DB.Exec(ctx,
		`INSERT INTO push_subscriptions (user_id, endpoint, p256dh_key, auth_key)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (endpoint) DO UPDATE SET user_id = $1, p256dh_key = $3, auth_key = $4`,
		userID, endpoint, p256dh, auth)
	return err
}

// RemoveSubscription removes a push subscription by endpoint
func (ps *PushService) RemoveSubscription(ctx context.Context, userID, endpoint string) error {
	_, err := ps.DB.Exec(ctx,
		"DELETE FROM push_subscriptions WHERE user_id = $1 AND endpoint = $2",
		userID, endpoint)
	return err
}

// SendToUser sends a push notification to all subscriptions for a user
func (ps *PushService) SendToUser(ctx context.Context, userID string, payload PushPayload) {
	rows, err := ps.DB.Query(ctx,
		"SELECT endpoint, p256dh_key, auth_key FROM push_subscriptions WHERE user_id = $1",
		userID)
	if err != nil {
		log.Printf("Push: failed to query subscriptions for user %s: %v", userID, err)
		return
	}
	defer rows.Close()

	payloadBytes, _ := json.Marshal(payload)

	for rows.Next() {
		var endpoint, p256dh, authKey string
		if err := rows.Scan(&endpoint, &p256dh, &authKey); err != nil {
			continue
		}

		sub := &webpush.Subscription{
			Endpoint: endpoint,
			Keys: webpush.Keys{
				P256dh: p256dh,
				Auth:   authKey,
			},
		}

		resp, err := webpush.SendNotification(payloadBytes, sub, &webpush.Options{
			Subscriber:      ps.VAPIDContact,
			VAPIDPublicKey:  ps.VAPIDPublicKey,
			VAPIDPrivateKey: ps.VAPIDPrivateKey,
			TTL:             60,
		})
		if err != nil {
			log.Printf("Push: failed to send to %s: %v", endpoint[:40], err)
			// If subscription expired (410 Gone), remove it
			if resp != nil && resp.StatusCode == 410 {
				ps.RemoveSubscription(ctx, userID, endpoint)
			}
			continue
		}
		resp.Body.Close()

		fmt.Printf("Push: sent to user %s (status %d)\n", userID[:8], resp.StatusCode)
	}
}
