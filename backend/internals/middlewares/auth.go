package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func (dep *dependencies.dependencies) AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		cookie, err := r.Cookie("session_id")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		sessionID := cookie.Value
		userID, err := dep.Redis.Get(ctx, "session:"+sessionID).Result()
		if err == nil {

			ctx = context.WithValue(ctx, "user_uuid", userID)
			ctx = context.WithValue(ctx, "session_id", sessionID)

			dep.Redis.Expire(ctx, "session:"+sessionID, 24*time.Hour)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		session, err := dep.DB.GetSession(sessionID)
		if err != nil || session.ExpiresAt.Before(time.Now()) {
			dep.clearSessionCookie(w)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		err = dep.Redis.Set(ctx, "session:"+sessionID, session.UserID, 24*time.Hour).Err()
		if err != nil {
			fmt.Errorf("failed to cache session in Redis", "error", err)
		}

		
		ctx = context.WithValue(ctx, "user_uuid", session.UserID)
		ctx = context.WithValue(ctx, "session_id", sessionID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
