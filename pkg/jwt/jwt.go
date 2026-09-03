package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret []byte
	jwtExpiry time.Duration
)

// Init 初始化 JWT 包配置，应在程序启动时调用
func Init(secret string, expireHours int) {
	jwtSecret = []byte(secret)
	if expireHours <= 0 {
		expireHours = 24 // 默认 24 小时
	}
	jwtExpiry = time.Duration(expireHours) * time.Hour
}

// GetExpiry 获取 JWT 过期时间
func GetExpiry() time.Duration {
	return jwtExpiry
}

// Claims JWT 载荷
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	RoleID   int64  `json:"role_id"`
	IsSuper  int16  `json:"is_super"`
	JTI      string `json:"jti"`
	jwt.RegisteredClaims
}

// generateJTI 生成唯一的 JWT ID（16 字节 hex = 32 字符）
func generateJTI() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// GenerateToken 生成 JWT token
func GenerateToken(userID int64, username string, roleID int64, isSuper int16) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RoleID:   roleID,
		IsSuper:  isSuper,
		JTI:      generateJTI(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "gin-blog-server",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析 JWT token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
