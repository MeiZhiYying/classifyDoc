package auth

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// Session 表示登录会话
type Session struct {
	Username  string
	CreatedAt time.Time
	ExpiresAt time.Time
}

var (
	sessions     = make(map[string]Session)
	sessionsLock sync.RWMutex
	// 默认会话有效期 24 小时
	defaultTTL = 24 * time.Hour
)

// CreateSession 创建会话并返回 sessionID
func CreateSession(username string) (string, Session) {
	sid := generateSessionID()
	now := time.Now()
	sess := Session{Username: username, CreatedAt: now, ExpiresAt: now.Add(defaultTTL)}

	sessionsLock.Lock()
	sessions[sid] = sess
	sessionsLock.Unlock()

	return sid, sess
}

// GetSession 根据 sessionID 获取会话
func GetSession(sessionID string) (Session, bool) {
	sessionsLock.RLock()
	sess, ok := sessions[sessionID]
	sessionsLock.RUnlock()
	if !ok {
		return Session{}, false
	}
	if time.Now().After(sess.ExpiresAt) {
		// 过期自动清理
		DeleteSession(sessionID)
		return Session{}, false
	}
	return sess, true
}

// DeleteSession 删除会话
func DeleteSession(sessionID string) {
	sessionsLock.Lock()
	delete(sessions, sessionID)
	sessionsLock.Unlock()
}

func generateSessionID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
