/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
package model

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

const defaultLobeHubSchema = "public"

var (
	ErrLobeHubPostgresRequired  = errors.New("LobeHub user management requires PostgreSQL")
	ErrLobeHubSchemaUnavailable = errors.New("the configured LobeHub database schema is unavailable or incompatible")
	ErrLobeHubUserNotFound      = errors.New("LobeHub user not found")
	ErrLobeHubConflict          = errors.New("LobeHub user data conflicts with an existing record")
	ErrLobeHubStaleUpdate       = errors.New("the LobeHub user was changed by another administrator")
	ErrLobeHubCredentialState   = errors.New("multiple credential accounts exist for this LobeHub user")
	ErrLobeHubRoleConfirmation  = errors.New("overwriting a nonstandard LobeHub role requires explicit confirmation")
)

type LobeHubUser struct {
	ID                  string     `json:"id" gorm:"column:id"`
	Username            *string    `json:"username" gorm:"column:username"`
	Email               *string    `json:"email" gorm:"column:email"`
	Avatar              *string    `json:"avatar" gorm:"column:avatar"`
	Phone               *string    `json:"phone" gorm:"column:phone"`
	FirstName           *string    `json:"first_name" gorm:"column:first_name"`
	LastName            *string    `json:"last_name" gorm:"column:last_name"`
	FullName            *string    `json:"full_name" gorm:"column:full_name"`
	EmailVerified       bool       `json:"email_verified" gorm:"column:email_verified"`
	PhoneNumberVerified bool       `json:"phone_number_verified" gorm:"column:phone_number_verified"`
	Role                *string    `json:"role" gorm:"column:role"`
	Banned              bool       `json:"banned" gorm:"column:banned"`
	BanReason           *string    `json:"ban_reason" gorm:"column:ban_reason"`
	BanExpires          *time.Time `json:"ban_expires" gorm:"column:ban_expires"`
	TwoFactorEnabled    bool       `json:"two_factor_enabled" gorm:"column:two_factor_enabled"`
	LastActiveAt        time.Time  `json:"last_active_at" gorm:"column:last_active_at"`
	CreatedAt           time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt           time.Time  `json:"updated_at" gorm:"column:updated_at"`
	PasswordSet         bool       `json:"password_set" gorm:"column:password_set"`
	SessionCount        int64      `json:"session_count" gorm:"column:session_count"`
}

type LobeHubLoginProvider struct {
	ProviderID string    `json:"provider_id" gorm:"column:provider_id"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
}

type LobeHubUserDetail struct {
	User      LobeHubUser            `json:"user"`
	Providers []LobeHubLoginProvider `json:"providers"`
}

type LobeHubUserListParams struct {
	Page          int
	PageSize      int
	Query         string
	Status        string
	Role          string
	EmailVerified *bool
	TwoFactor     *bool
	SortBy        string
	SortOrder     string
}

type LobeHubUserList struct {
	Items    []LobeHubUser `json:"items"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

type LobeHubUserUpdate struct {
	Username            *string
	Email               *string
	Avatar              *string
	Phone               *string
	FirstName           *string
	LastName            *string
	FullName            *string
	EmailVerified       bool
	PhoneNumberVerified bool
	ExpectedUpdatedAt   time.Time
}

type LobeHubRevokeResult struct {
	Sessions int64 `json:"sessions"`
	OIDC     int64 `json:"oidc"`
}

type LobeHubPasswordResetResult struct {
	Revoked LobeHubRevokeResult `json:"revoked"`
}

var lobeHubRequiredColumns = map[string][]string{
	"users": {
		"id", "username", "email", "normalized_email", "avatar", "phone", "first_name",
		"last_name", "full_name", "email_verified", "phone_number_verified", "role", "banned",
		"ban_reason", "ban_expires", "two_factor_enabled", "last_active_at", "created_at", "updated_at",
	},
	"accounts":                 {"id", "user_id", "account_id", "provider_id", "password", "created_at", "updated_at"},
	"auth_sessions":            {"id", "token", "user_id"},
	"oidc_access_tokens":       {"user_id"},
	"oidc_authorization_codes": {"user_id"},
	"oidc_refresh_tokens":      {"user_id"},
	"oidc_device_codes":        {"user_id"},
	"oidc_grants":              {"user_id"},
	"oidc_sessions":            {"user_id"},
}

func lobeHubSchema() (string, error) {
	schema := strings.TrimSpace(os.Getenv("LOBEHUB_DB_SCHEMA"))
	if schema == "" {
		schema = defaultLobeHubSchema
	}
	if len(schema) > 63 {
		return "", ErrLobeHubSchemaUnavailable
	}
	for index, char := range schema {
		if (char >= 'a' && char <= 'z') || char == '_' || (index > 0 && char >= '0' && char <= '9') {
			continue
		}
		return "", ErrLobeHubSchemaUnavailable
	}
	return schema, nil
}

func lobeHubTable(schema string, table string) string {
	return fmt.Sprintf(`"%s"."%s"`, schema, table)
}

func checkLobeHubSchema(requiredColumns map[string][]string) (string, error) {
	if !common.UsingMainDatabase(common.DatabaseTypePostgreSQL) || DB == nil || DB.Dialector.Name() != "postgres" {
		return "", ErrLobeHubPostgresRequired
	}
	schema, err := lobeHubSchema()
	if err != nil {
		return "", err
	}

	type columnRow struct {
		TableName  string `gorm:"column:table_name"`
		ColumnName string `gorm:"column:column_name"`
	}
	var rows []columnRow
	if err := DB.Raw(
		`SELECT table_name, column_name FROM information_schema.columns WHERE table_schema = ?`,
		schema,
	).Scan(&rows).Error; err != nil {
		return "", fmt.Errorf("%w: %v", ErrLobeHubSchemaUnavailable, err)
	}

	found := make(map[string]map[string]struct{})
	for _, row := range rows {
		if found[row.TableName] == nil {
			found[row.TableName] = make(map[string]struct{})
		}
		found[row.TableName][row.ColumnName] = struct{}{}
	}
	for table, columns := range requiredColumns {
		for _, column := range columns {
			if _, ok := found[table][column]; !ok {
				return "", fmt.Errorf("%w: missing %s.%s", ErrLobeHubSchemaUnavailable, table, column)
			}
		}
	}
	return schema, nil
}

func CheckLobeHubSchema() (string, error) {
	return checkLobeHubSchema(lobeHubRequiredColumns)
}

func lobeHubUserSelect(schema string) string {
	accounts := lobeHubTable(schema, "accounts")
	sessions := lobeHubTable(schema, "auth_sessions")
	return fmt.Sprintf(`
		u.id, u.username, u.email, u.avatar, u.phone, u.first_name, u.last_name, u.full_name,
		u.email_verified, COALESCE(u.phone_number_verified, false) AS phone_number_verified,
		u.role, COALESCE(u.banned, false) AS banned, u.ban_reason, u.ban_expires,
		COALESCE(u.two_factor_enabled, false) AS two_factor_enabled, u.last_active_at, u.created_at, u.updated_at,
		EXISTS (SELECT 1 FROM %s a WHERE a.user_id = u.id AND a.provider_id = 'credential' AND a.password IS NOT NULL) AS password_set,
		(SELECT COUNT(*) FROM %s s WHERE s.user_id = u.id) AS session_count`, accounts, sessions)
}

func escapeLobeHubSearch(value string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(value)
}

func ListLobeHubUsers(params LobeHubUserListParams) (*LobeHubUserList, error) {
	schema, err := CheckLobeHubSchema()
	if err != nil {
		return nil, err
	}

	usersTable := lobeHubTable(schema, "users")
	where := []string{"1 = 1"}
	args := make([]interface{}, 0)
	if query := strings.TrimSpace(params.Query); query != "" {
		pattern := "%" + escapeLobeHubSearch(query) + "%"
		where = append(where, `(u.id ILIKE ? ESCAPE '\' OR u.username ILIKE ? ESCAPE '\' OR u.email ILIKE ? ESCAPE '\' OR u.full_name ILIKE ? ESCAPE '\' OR u.phone ILIKE ? ESCAPE '\')`)
		args = append(args, pattern, pattern, pattern, pattern, pattern)
	}
	switch params.Status {
	case "active":
		where = append(where, `COALESCE(u.banned, false) = false`)
	case "banned":
		where = append(where, `COALESCE(u.banned, false) = true AND (u.ban_expires IS NULL OR u.ban_expires > NOW())`)
	case "expired":
		where = append(where, `COALESCE(u.banned, false) = true AND u.ban_expires <= NOW()`)
	}
	switch params.Role {
	case "user":
		where = append(where, `COALESCE(NULLIF(BTRIM(u.role), ''), 'user') = 'user'`)
	case "admin":
		where = append(where, `BTRIM(u.role) = 'admin'`)
	case "other":
		where = append(where, `COALESCE(NULLIF(BTRIM(u.role), ''), 'user') NOT IN ('user', 'admin')`)
	}
	if params.EmailVerified != nil {
		where = append(where, `u.email_verified = ?`)
		args = append(args, *params.EmailVerified)
	}
	if params.TwoFactor != nil {
		where = append(where, `COALESCE(u.two_factor_enabled, false) = ?`)
		args = append(args, *params.TwoFactor)
	}

	whereSQL := strings.Join(where, " AND ")
	var total int64
	if err := DB.Raw(`SELECT COUNT(*) FROM `+usersTable+` u WHERE `+whereSQL, args...).Scan(&total).Error; err != nil {
		return nil, err
	}

	sortColumns := map[string]string{
		"created_at":     "u.created_at",
		"last_active_at": "u.last_active_at",
		"username":       "u.username",
		"email":          "u.email",
	}
	sortColumn := sortColumns[params.SortBy]
	if sortColumn == "" {
		sortColumn = "u.created_at"
	}
	sortOrder := "DESC"
	if strings.EqualFold(params.SortOrder, "asc") {
		sortOrder = "ASC"
	}

	queryArgs := append(append([]interface{}{}, args...), params.PageSize, (params.Page-1)*params.PageSize)
	query := `SELECT ` + lobeHubUserSelect(schema) + ` FROM ` + usersTable + ` u WHERE ` + whereSQL +
		` ORDER BY ` + sortColumn + ` ` + sortOrder + ` NULLS LAST, u.id ASC LIMIT ? OFFSET ?`
	items := make([]LobeHubUser, 0)
	if err := DB.Raw(query, queryArgs...).Scan(&items).Error; err != nil {
		return nil, err
	}
	return &LobeHubUserList{Items: items, Total: total, Page: params.Page, PageSize: params.PageSize}, nil
}

func GetLobeHubUser(id string) (*LobeHubUserDetail, error) {
	schema, err := CheckLobeHubSchema()
	if err != nil {
		return nil, err
	}
	usersTable := lobeHubTable(schema, "users")
	var user LobeHubUser
	result := DB.Raw(`SELECT `+lobeHubUserSelect(schema)+` FROM `+usersTable+` u WHERE u.id = ?`, id).Scan(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if user.ID == "" {
		return nil, ErrLobeHubUserNotFound
	}

	providers := make([]LobeHubLoginProvider, 0)
	if err := DB.Raw(
		`SELECT provider_id, MIN(created_at) AS created_at FROM `+lobeHubTable(schema, "accounts")+` WHERE user_id = ? GROUP BY provider_id ORDER BY provider_id`,
		id,
	).Scan(&providers).Error; err != nil {
		return nil, err
	}
	return &LobeHubUserDetail{User: user, Providers: providers}, nil
}

func UpdateLobeHubUser(id string, update LobeHubUserUpdate) (*LobeHubUser, error) {
	schema, err := CheckLobeHubSchema()
	if err != nil {
		return nil, err
	}
	usersTable := lobeHubTable(schema, "users")
	values := map[string]interface{}{
		"username":              update.Username,
		"email":                 update.Email,
		"normalized_email":      update.Email,
		"avatar":                update.Avatar,
		"phone":                 update.Phone,
		"first_name":            update.FirstName,
		"last_name":             update.LastName,
		"full_name":             update.FullName,
		"email_verified":        update.EmailVerified,
		"phone_number_verified": update.PhoneNumberVerified,
		"updated_at":            time.Now(),
	}
	result := DB.Table(usersTable).Where("id = ? AND updated_at = ?", id, update.ExpectedUpdatedAt).Updates(values)
	if result.Error != nil {
		return nil, normalizeLobeHubWriteError(result.Error)
	}
	if result.RowsAffected == 0 {
		var count int64
		if err := DB.Table(usersTable).Where("id = ?", id).Count(&count).Error; err != nil {
			return nil, err
		}
		if count == 0 {
			return nil, ErrLobeHubUserNotFound
		}
		return nil, ErrLobeHubStaleUpdate
	}
	detail, err := GetLobeHubUser(id)
	if err != nil {
		return nil, err
	}
	return &detail.User, nil
}

func SetLobeHubUserBan(id string, reason *string, expires *time.Time, banned bool) (*LobeHubUser, LobeHubRevokeResult, error) {
	schema, err := CheckLobeHubSchema()
	if err != nil {
		return nil, LobeHubRevokeResult{}, err
	}
	var revoked LobeHubRevokeResult
	err = DB.Transaction(func(tx *gorm.DB) error {
		if err := lockLobeHubUser(tx, schema, id); err != nil {
			return err
		}
		values := map[string]interface{}{
			"banned":      banned,
			"ban_reason":  reason,
			"ban_expires": expires,
			"updated_at":  time.Now(),
		}
		if !banned {
			values["ban_reason"] = nil
			values["ban_expires"] = nil
		}
		if err := tx.Table(lobeHubTable(schema, "users")).Where("id = ?", id).Updates(values).Error; err != nil {
			return normalizeLobeHubWriteError(err)
		}
		if banned {
			var err error
			revoked, err = revokeLobeHubCredentials(tx, schema, id)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, LobeHubRevokeResult{}, err
	}
	detail, err := GetLobeHubUser(id)
	if err != nil {
		return nil, LobeHubRevokeResult{}, err
	}
	return &detail.User, revoked, nil
}

func SetLobeHubUserRole(id string, role string, confirmOverwriteCustomRole bool) (*LobeHubUser, LobeHubRevokeResult, error) {
	schema, err := CheckLobeHubSchema()
	if err != nil {
		return nil, LobeHubRevokeResult{}, err
	}
	var revoked LobeHubRevokeResult
	err = DB.Transaction(func(tx *gorm.DB) error {
		if err := lockLobeHubUser(tx, schema, id); err != nil {
			return err
		}
		var current struct {
			Role *string `gorm:"column:role"`
		}
		if err := tx.Table(lobeHubTable(schema, "users")).Select("role").Where("id = ?", id).Scan(&current).Error; err != nil {
			return err
		}
		if !confirmOverwriteCustomRole && (current.Role == nil || (*current.Role != "user" && *current.Role != "admin")) {
			return ErrLobeHubRoleConfirmation
		}
		if err := tx.Table(lobeHubTable(schema, "users")).Where("id = ?", id).Updates(map[string]interface{}{
			"role":       role,
			"updated_at": time.Now(),
		}).Error; err != nil {
			return normalizeLobeHubWriteError(err)
		}
		var err error
		revoked, err = revokeLobeHubCredentials(tx, schema, id)
		return err
	})
	if err != nil {
		return nil, LobeHubRevokeResult{}, err
	}
	detail, err := GetLobeHubUser(id)
	if err != nil {
		return nil, LobeHubRevokeResult{}, err
	}
	return &detail.User, revoked, nil
}

func ResetLobeHubUserPassword(id string, accountID string, passwordHash string) (*LobeHubPasswordResetResult, error) {
	schema, err := CheckLobeHubSchema()
	if err != nil {
		return nil, err
	}
	var revoked LobeHubRevokeResult
	err = DB.Transaction(func(tx *gorm.DB) error {
		if err := lockLobeHubUser(tx, schema, id); err != nil {
			return err
		}
		type credentialRow struct {
			ID string `gorm:"column:id"`
		}
		var credentials []credentialRow
		if err := tx.Raw(
			`SELECT id FROM `+lobeHubTable(schema, "accounts")+` WHERE user_id = ? AND provider_id = 'credential' FOR UPDATE`,
			id,
		).Scan(&credentials).Error; err != nil {
			return err
		}
		if len(credentials) > 1 {
			return ErrLobeHubCredentialState
		}
		now := time.Now()
		if len(credentials) == 1 {
			if err := tx.Table(lobeHubTable(schema, "accounts")).Where("id = ?", credentials[0].ID).Updates(map[string]interface{}{
				"password":   passwordHash,
				"updated_at": now,
			}).Error; err != nil {
				return normalizeLobeHubWriteError(err)
			}
		} else {
			if err := tx.Exec(
				`INSERT INTO `+lobeHubTable(schema, "accounts")+` (id, user_id, account_id, provider_id, password, created_at, updated_at) VALUES (?, ?, ?, 'credential', ?, ?, ?)`,
				accountID, id, id, passwordHash, now, now,
			).Error; err != nil {
				return normalizeLobeHubWriteError(err)
			}
		}
		var err error
		revoked, err = revokeLobeHubCredentials(tx, schema, id)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &LobeHubPasswordResetResult{Revoked: revoked}, nil
}

func lockLobeHubUser(tx *gorm.DB, schema string, id string) error {
	var user LobeHubUser
	result := lockForUpdate(tx).Table(lobeHubTable(schema, "users")).Select("id").Where("id = ?", id).Take(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ErrLobeHubUserNotFound
	}
	return result.Error
}

func revokeLobeHubCredentials(tx *gorm.DB, schema string, id string) (LobeHubRevokeResult, error) {
	result := LobeHubRevokeResult{}
	sessions := tx.Exec(`DELETE FROM `+lobeHubTable(schema, "auth_sessions")+` WHERE user_id = ?`, id)
	if sessions.Error != nil {
		return result, sessions.Error
	}
	result.Sessions = sessions.RowsAffected
	for _, table := range []string{
		"oidc_access_tokens",
		"oidc_authorization_codes",
		"oidc_refresh_tokens",
		"oidc_device_codes",
		"oidc_grants",
		"oidc_sessions",
	} {
		deleted := tx.Exec(`DELETE FROM `+lobeHubTable(schema, table)+` WHERE user_id = ?`, id)
		if deleted.Error != nil {
			return result, deleted.Error
		}
		result.OIDC += deleted.RowsAffected
	}
	return result, nil
}

func normalizeLobeHubWriteError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrLobeHubConflict
	}
	return err
}
