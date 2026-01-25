package users

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	repo "github.com/pick-cee/go-ecommerce-api/internal/adapters/postgresql/sqlc"
	"github.com/pick-cee/go-ecommerce-api/internal/env"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte(env.GetString("JWT_SECRET_KEY", ""))

type createUserDto struct {
	Name string `json:"name"`
	loginUserDto
}

type loginUserDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User  repo.User `json:"user"`
  Token string    `json:"token"`
}


type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	jwt.RegisteredClaims
}

func (l loginUserDto) Validate() bool {
	return l.Username != "" && l.Password != ""
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	 if err != nil {
		log.Println("Error hashing password:", "error", err)
	 }

	return  string(bytes), nil
}

func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}

func CreateToken(userId int, username, name string) (string, error) {
	claims := Claims{
		UserID: userId,
		Username: username,
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 2)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
      return nil, err
   }
  
  if claims, ok := token.Claims.(*Claims); ok && token.Valid {
    return claims, nil
  }
  
  return nil, errors.New("invalid token")
}