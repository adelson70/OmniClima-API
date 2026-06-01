package user

import (
	"OmniClima/internal/platform/apperror"
	"net/http"
	"os"
	"strconv"
	"strings"
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

type UserInput struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type UpdateUserInput struct {
	UserID    uuid.UUID
	FirstName *string
	LastName  *string
	Email     *string
	Password  *string
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

func (s *Service) CreateUser(in UserInput) (UserOutput, error) {
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

func (s *Service) UpdateUser(in UpdateUserInput) error {

	if !in.hasUpdates() {
		return apperror.New(http.StatusBadRequest, "Nenhum campo para atualizar")
	}

	updates := map[string]interface{}{}

	if stringFieldSet(in.FirstName) {
		updates["first_name"] = strings.TrimSpace(*in.FirstName)
	}
	if stringFieldSet(in.LastName) {
		updates["last_name"] = strings.TrimSpace(*in.LastName)
	}
	if stringFieldSet(in.Email) {
		updates["email"] = strings.TrimSpace(*in.Email)
	}
	if stringFieldSet(in.Password) {
		passwordHashed, err := hashPassword(strings.TrimSpace(*in.Password))

		if err != nil {
			return err
		}

		updates["password"] = passwordHashed
	}

	if err := s.repo.UpdateUser(in.UserID, updates); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteUser(userID uuid.UUID) error {
	return s.repo.DeleteUser(userID)
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

func (in UpdateUserInput) hasUpdates() bool {
	return stringFieldSet(in.FirstName) ||
		stringFieldSet(in.LastName) ||
		stringFieldSet(in.Email) ||
		stringFieldSet(in.Password)
}

func stringFieldSet(s *string) bool {
	return s != nil && strings.TrimSpace(*s) != ""
}
