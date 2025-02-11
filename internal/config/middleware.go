package config

import (
	"go_base/internal/repositories"
	"go_base/internal/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type Middleware struct {
	userRepo repositories.UserRepository
}

func NewMiddleware(userRepo repositories.UserRepository) *Middleware {
	return &Middleware{userRepo: userRepo}
}

func (m *Middleware) CheckAuth(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if auth == "" {
		return c.JSON(utils.ResponseApi(fiber.StatusUnauthorized, "Not Authorized", nil))
	}

	tokenString := strings.Split(auth, " ")
	if len(tokenString) != 2 {
		return c.JSON(utils.ResponseApi(fiber.StatusUnauthorized, "Invalid token", nil))
	}

	token, err := jwt.Parse(tokenString[1], func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, nil
		}

		return []byte(utils.GetEnv("JWT_SECRET_KEY", "secret")), nil
	})
	if err != nil {
		return c.JSON(utils.ResponseApi(fiber.StatusUnauthorized, err.Error(), nil))
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.JSON(utils.ResponseApi(fiber.StatusUnauthorized, "Invalid token", nil))
	}
	userID, ok := claims["user_id"]
	if !ok {
		return c.JSON(utils.ResponseApi(fiber.StatusUnauthorized, "Invalid token", nil))
	}

	user, err := m.userRepo.GetUserByID(userID.(string))
	if err != nil {
		return c.JSON(utils.ResponseApi(fiber.StatusUnauthorized, "Invalid token", nil))
	}

	c.Locals("user", user)

	return c.Next()
}
