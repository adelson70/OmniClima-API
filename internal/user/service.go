package user

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type CreateUserInput struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type UserOutput struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Token     string
	Admin     bool
}

type Claims struct {
	UserID   uuid.UUID `json:"userID"`
	Email    string    `json:"email"`
	FullName string    `json:"fullName"`
	Admin    bool      `json:"admin"`

	jwt.RegisteredClaims
}

func (s *Service) CreateUser(in CreateUserInput) (UserOutput, error) {
	passwordHashed, err := hashPassword(in.Password)

	if err != nil {
		return UserOutput{}, err
	}

	user := &User{
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Password:  passwordHashed,
		Email:     in.Email,
		Admin:     false,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return UserOutput{}, err
	}

	token, err := GenerateToken(user.ID, user.FirstName, user.LastName, user.Email, user.Admin)

	if err != nil {
		return UserOutput{}, err
	}

	return UserOutput{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Token:     token,
		Admin:     user.Admin,
	}, nil

}

func hashPassword(password string) (string, error) {
	costStr := os.Getenv("BCRYPT_COST")

	cost, err := strconv.Atoi(costStr)
	if err != nil {
		cost = bcrypt.DefaultCost
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		cost,
	)

	return string(hash), err
}

func GenerateToken(userID uuid.UUID, firstName, lastName, email string, admin bool) (string, error) {
	claims := Claims{
		UserID:   userID,
		FullName: firstName + " " + lastName,
		Email:    email,
		Admin:    admin,

		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
