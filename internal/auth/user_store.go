package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
)

// 为避免引入外部依赖，这里使用 sha256(salt+password) 简易哈希。
// 后续上线可替换为 bcrypt/argon2，并持久化到外部存储。

type userRecord struct {
	Salt string
	Hash string
}

var (
	users     = make(map[string]userRecord)
	usersLock sync.RWMutex
)

var (
	ErrUserExists    = errors.New("用户已存在")
	ErrUserNotFound  = errors.New("用户不存在")
	ErrInvalidSecret = errors.New("密码错误")
)

func RegisterUser(username, password string) error {
	usersLock.Lock()
	defer usersLock.Unlock()
	if _, ok := users[username]; ok {
		return ErrUserExists
	}
	salt := randomHex(16)
	hash := hashPassword(salt, password)
	users[username] = userRecord{Salt: salt, Hash: hash}
	return nil
}

func ValidateUser(username, password string) error {
	usersLock.RLock()
	rec, ok := users[username]
	usersLock.RUnlock()
	if !ok {
		return ErrUserNotFound
	}
	if rec.Hash != hashPassword(rec.Salt, password) {
		return ErrInvalidSecret
	}
	return nil
}

func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func hashPassword(salt, password string) string {
	h := sha256.Sum256([]byte(salt + ":" + password))
	return hex.EncodeToString(h[:])
}
