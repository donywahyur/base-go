package repositories

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"go_base/internal/models"
	"go_base/internal/utils"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/argon2"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByUsername(username string) (models.User, error)
	GetUserByID(userID string) (models.User, error)
	HashPassword(password string) (string, error)
	CompareHash(password, hashedPassword string) (bool, error)
	GenereateJWTToken(userID string) (string, error)
}

type userRepository struct {
	db    *gorm.DB
	param params
}
type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db,
		param: params{
			memory:      64 * 1024,
			iterations:  3,
			parallelism: 2,
			saltLength:  16,
			keyLength:   32,
		}}
}

func (r *userRepository) GetUserByUsername(username string) (models.User, error) {
	var user models.User

	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return user, nil
	}

	return user, nil
}

func (r *userRepository) GetUserByID(userID string) (models.User, error) {
	var user models.User

	err := r.db.Preload("Role").Where("id = ?", userID).First(&user).Error
	if err != nil {
		return user, nil
	}

	return user, nil
}

func (r *userRepository) HashPassword(password string) (string, error) {
	salt := make([]byte, r.param.saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, r.param.iterations, r.param.memory, r.param.parallelism, r.param.keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, r.param.memory, r.param.iterations, r.param.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func (r *userRepository) CompareHash(password, passwordHash string) (bool, error) {
	vals := strings.Split(passwordHash, "$")
	if len(vals) != 6 {
		return false, errors.New("invalid hash")
	}

	var memory, time uint32
	var parallelism uint8

	_, err := fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &memory, &time, &parallelism)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return false, err
	}

	decryptedHash, err := base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return false, err
	}

	var keyLen = uint32(len(decryptedHash))

	comparisonHash := argon2.IDKey([]byte(password), salt, time, memory, parallelism, keyLen)

	return subtle.ConstantTimeCompare(comparisonHash, decryptedHash) == 1, nil

}

func (r *userRepository) GenereateJWTToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtKey := []byte(utils.GetEnv("JWT_SECRET_KEY", "secret"))
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
