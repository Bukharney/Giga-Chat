package repositories

import (
	"fmt"
	"os"
	"time"

	"github.com/bukharney/giga-chat/modules/entities"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AuthRepo struct {
	Db *sqlx.DB
}

func NewAuthRepo(db *sqlx.DB) entities.AuthRepository {
	return &AuthRepo{Db: db}
}

func (a *AuthRepo) SignUsersAccessToken(req *entities.UsersPassport) (string, error) {
	claims := &entities.UsersClaims{
		Id:       req.Id,
		Username: req.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "users_access_token",
			Issuer:    "gigachat_auth",
			ID:        uuid.New().String(),
		},
	}

	mySigningKey := os.Getenv("JWT_SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(mySigningKey))
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return ss, nil
}
