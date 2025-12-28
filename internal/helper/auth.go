package helper

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Mukam21/Go_E-Commerce_App/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	Secret string
}

func SetupAuth(s string) Auth {
	return Auth{
		Secret: s,
	}
}

func (a Auth) CreateHashedPassword(p string) (string, error) {

	if len(p) < 6 {
		return "", errors.New("password length should be at least 6 characters long")
	}

	hashP, err := bcrypt.GenerateFromPassword([]byte(p), 10)

	if err != nil {
		return "", errors.New("failed to hash password")
	}

	return string(hashP), nil
}

func (a Auth) GenerateToken(id uint, email string, role string) (string, error) {

	if id == 0 || email == "" || role == "" {
		return "", errors.New("required inputs are missing to generate token")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(), // token valid for 30 days
	})

	tokenStr, err := token.SignedString([]byte(a.Secret))

	if err != nil {
		return "", errors.New("unable to signed the token")
	}

	return tokenStr, nil
}

func (a Auth) VerifyPassword(pP string, hP string) error {

	if len(pP) < 6 {
		return errors.New("password length should be at least 6 characters long")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hP), []byte(pP))

	if err != nil {
		return errors.New("password does not match")
	}

	return nil
}

func (a Auth) VerifyToken(t string) (domain.User, error) {

	tokenArr := strings.Split(t, " ")
	if len(tokenArr) != 2 {
		return domain.User{}, errors.New("invalid authorization header")
	}

	if tokenArr[0] != "Bearer" {
		return domain.User{}, errors.New("invalid token type")
	}

	tokenStr := tokenArr[1]

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(a.Secret), nil
	})

	if err != nil {
		return domain.User{}, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return domain.User{}, errors.New("invalid token claims")
	}

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		return domain.User{}, errors.New("token expired")
	}

	user := domain.User{
		ID:       uint(claims["user_id"].(float64)),
		Email:    claims["email"].(string),
		UserType: claims["role"].(string),
	}

	return user, nil
}

func (a Auth) Authorize(ctx *fiber.Ctx) error {

	headers := ctx.GetReqHeaders()

	authHeaders, ok := headers["Authorization"]
	if !ok || len(authHeaders) == 0 {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "missing authorization header",
		})
	}

	authHeader := authHeaders[0] // ← string

	user, err := a.VerifyToken(authHeader)

	if err == nil && user.ID > 0 {
		ctx.Locals("user", user)
		return ctx.Next()
	}

	return ctx.Status(401).JSON(fiber.Map{
		"message": "authorization failed",
	})
}

func (a Auth) GetCurrentUser(ctx *fiber.Ctx) domain.User {

	user := ctx.Locals("user")

	return user.(domain.User)
}

func (a Auth) GenerateCode() (int, error) {
	return RandomNomber(6)
}
