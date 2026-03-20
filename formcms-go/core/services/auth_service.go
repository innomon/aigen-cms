package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/formcms/formcms-go/core/descriptors"
	"github.com/formcms/formcms-go/infrastructure/relationdbdao"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	dao    relationdbdao.IPrimaryDao
	secret []byte
}

func NewAuthService(dao relationdbdao.IPrimaryDao, secret string) *AuthService {
	return &AuthService{
		dao:    dao,
		secret: []byte(secret),
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (*descriptors.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &descriptors.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         descriptors.RoleUser,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	query, args, err := s.dao.GetBuilder().Insert(descriptors.UserTableName).
		Columns("email", "password_hash", "role", "created_at", "updated_at").
		Values(user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt).ToSql()
	if err != nil {
		return nil, err
	}

	var newId int64
	if strings.Contains(query, "$1") {
		// Postgres
		err = s.dao.GetDb().QueryRowContext(ctx, query+" RETURNING id", args...).Scan(&newId)
	} else {
		// SQLite
		res, err := s.dao.GetDb().ExecContext(ctx, query, args...)
		if err != nil {
			return nil, err
		}
		newId, err = res.LastInsertId()
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	user.Id = newId
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	query, args, err := s.dao.GetBuilder().Select("*").From(descriptors.UserTableName).
		Where(squirrel.Eq{"email": email}).Limit(1).ToSql()
	if err != nil {
		return "", err
	}

	row := s.dao.GetDb().QueryRowContext(ctx, query, args...)
	var id int64
	var e, pass, role string
	var avatar sql.NullString
	var createdAt, updatedAt time.Time
	if err := row.Scan(&id, &e, &pass, &role, &avatar, &createdAt, &updatedAt); err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(pass), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": id,
		"role":   role,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(s.secret)
}

func (s *AuthService) Me(ctx context.Context, userId int64) (*descriptors.User, error) {
	query, args, err := s.dao.GetBuilder().Select("*").From(descriptors.UserTableName).
		Where(squirrel.Eq{"id": userId}).Limit(1).ToSql()
	if err != nil {
		return nil, err
	}

	row := s.dao.GetDb().QueryRowContext(ctx, query, args...)
	var id int64
	var e, pass, role string
	var avatar sql.NullString
	var createdAt, updatedAt time.Time
	if err := row.Scan(&id, &e, &pass, &role, &avatar, &createdAt, &updatedAt); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &descriptors.User{
		Id:           id,
		Email:        e,
		PasswordHash: pass,
		Role:         role,
		AvatarPath:   avatar.String,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (int64, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil {
		return 0, "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := int64(claims["userId"].(float64))
		role := claims["role"].(string)
		return userId, role, nil
	}

	return 0, "", fmt.Errorf("invalid token")
}
