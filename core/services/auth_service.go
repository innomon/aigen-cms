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
		Roles:        []string{descriptors.RoleUser},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	query, args, err := s.dao.GetBuilder().Insert(descriptors.UserTableName).
		Columns("email", "password_hash", "role", "created_at", "updated_at").
		Values(user.Email, user.PasswordHash, user.Roles[0], user.CreatedAt, user.UpdatedAt).ToSql()
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

	// Link role in junction table if RBAC is active
	// For now, we keep the legacy 'role' column updated too for compatibility
	_, _, _ = s.dao.GetBuilder().Insert("__user_roles").
		Columns("user_id", "role_id").
		Select(squirrel.Select(fmt.Sprintf("%d", user.Id), "id").From("__roles").Where(squirrel.Eq{"name": descriptors.RoleUser})).
		ToSql()
	// (Implementation of the above would need more robust error handling and potentially moving to a method)

	return user, nil
}

func (s *AuthService) getRoles(ctx context.Context, userId int64, legacyRole string) ([]string, []descriptors.Role, error) {
	query, args, err := s.dao.GetBuilder().Select("r.id, r.name, r.disabled, r.dashboard_page_id, r.menu_id").From("__user_roles ur").
		Join("__roles r ON ur.role_id = r.id").
		Where(squirrel.Eq{"ur.user_id": userId}).ToSql()

	if err != nil {
		return []string{legacyRole}, nil, nil // Fallback to legacy
	}

	rows, err := s.dao.GetDb().QueryContext(ctx, query, args...)
	if err != nil {
		return []string{legacyRole}, nil, nil // Fallback to legacy
	}
	defer rows.Close()

	var roles []string
	var rolesDetails []descriptors.Role
	for rows.Next() {
		var role descriptors.Role
		var dashboardPageId, menuId sql.NullString
		if err := rows.Scan(&role.Id, &role.Name, &role.Disabled, &dashboardPageId, &menuId); err == nil {
			role.DashboardPageId = dashboardPageId.String
			role.MenuId = menuId.String
			rolesDetails = append(rolesDetails, role)
			roles = append(roles, role.Name)
		}
	}

	if len(roles) == 0 {
		return []string{legacyRole}, nil, nil
	}

	return roles, rolesDetails, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	query, args, err := s.dao.GetBuilder().Select("id", "email", "password_hash", "role", "avatar_path", "default_role_id", "created_at", "updated_at").From(descriptors.UserTableName).
		Where(squirrel.Eq{"email": email}).Limit(1).ToSql()
	if err != nil {
		return "", err
	}

	row := s.dao.GetDb().QueryRowContext(ctx, query, args...)
	var id int64
	var e, pass, role string
	var avatar sql.NullString
	var defaultRoleId sql.NullInt64
	var createdAt, updatedAt time.Time
	if err := row.Scan(&id, &e, &pass, &role, &avatar, &defaultRoleId, &createdAt, &updatedAt); err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(pass), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid password")
	}

	roles, _, _ := s.getRoles(ctx, id, role)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": id,
		"roles":  roles,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(s.secret)
}

func (s *AuthService) Me(ctx context.Context, userId int64) (*descriptors.User, error) {
	query, args, err := s.dao.GetBuilder().Select("id", "email", "password_hash", "role", "avatar_path", "default_role_id", "created_at", "updated_at").From(descriptors.UserTableName).
		Where(squirrel.Eq{"id": userId}).Limit(1).ToSql()
	if err != nil {
		return nil, err
	}

	row := s.dao.GetDb().QueryRowContext(ctx, query, args...)
	var id int64
	var e, pass, role string
	var avatar sql.NullString
	var defaultRoleId sql.NullInt64
	var createdAt, updatedAt time.Time
	if err := row.Scan(&id, &e, &pass, &role, &avatar, &defaultRoleId, &createdAt, &updatedAt); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	roles, rolesDetails, _ := s.getRoles(ctx, id, role)
	
	var defaultRolePtr *int64
	if defaultRoleId.Valid {
		defaultRolePtr = &defaultRoleId.Int64
	}

	return &descriptors.User{
		Id:            id,
		Email:         e,
		PasswordHash:  pass,
		Roles:         roles,
		RolesDetails:  rolesDetails,
		DefaultRoleId: defaultRolePtr,
		AvatarPath:    avatar.String,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (int64, []string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil {
		return 0, nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := int64(claims["userId"].(float64))
		
		var roles []string
		if r, ok := claims["roles"].([]interface{}); ok {
			for _, role := range r {
				roles = append(roles, role.(string))
			}
		} else if r, ok := claims["role"].(string); ok {
			// Backward compatibility with single role tokens
			roles = []string{r}
		}

		return userId, roles, nil
	}

	return 0, nil, fmt.Errorf("invalid token")
}
