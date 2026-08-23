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
package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net/mail"
	"net/url"
	"strings"
	"time"
	"unicode"

	"github.com/QuantumNous/new-api/model"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/text/unicode/norm"
)

var ErrInvalidLobeHubUserInput = errors.New("invalid LobeHub user input")

type LobeHubUserUpdateInput struct {
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

type LobeHubUserChange struct {
	Before model.LobeHubUser `json:"before"`
	After  model.LobeHubUser `json:"after"`
}

type LobeHubUserActionResult struct {
	User    model.LobeHubUser         `json:"user"`
	Revoked model.LobeHubRevokeResult `json:"revoked"`
}

type LobeHubPasswordReset struct {
	TemporaryPassword string                    `json:"temporary_password"`
	Revoked           model.LobeHubRevokeResult `json:"revoked"`
}

func UpdateLobeHubUser(id string, input LobeHubUserUpdateInput) (*LobeHubUserChange, error) {
	if strings.TrimSpace(id) == "" || input.ExpectedUpdatedAt.IsZero() {
		return nil, ErrInvalidLobeHubUserInput
	}
	detail, err := model.GetLobeHubUser(id)
	if err != nil {
		return nil, err
	}

	username, err := normalizeLobeHubText(input.Username)
	if err != nil {
		return nil, err
	}
	email, err := normalizeLobeHubEmail(input.Email)
	if err != nil {
		return nil, err
	}
	avatar, err := normalizeLobeHubAvatar(input.Avatar)
	if err != nil {
		return nil, err
	}
	phone, err := normalizeLobeHubText(input.Phone)
	if err != nil {
		return nil, err
	}
	firstName, err := normalizeLobeHubText(input.FirstName)
	if err != nil {
		return nil, err
	}
	lastName, err := normalizeLobeHubText(input.LastName)
	if err != nil {
		return nil, err
	}
	fullName, err := normalizeLobeHubText(input.FullName)
	if err != nil {
		return nil, err
	}

	emailVerified := input.EmailVerified
	if !equalNullableString(detail.User.Email, email) {
		emailVerified = false
	}
	phoneVerified := input.PhoneNumberVerified
	if !equalNullableString(detail.User.Phone, phone) {
		phoneVerified = false
	}
	after, err := model.UpdateLobeHubUser(id, model.LobeHubUserUpdate{
		Username:            username,
		Email:               email,
		Avatar:              avatar,
		Phone:               phone,
		FirstName:           firstName,
		LastName:            lastName,
		FullName:            fullName,
		EmailVerified:       emailVerified,
		PhoneNumberVerified: phoneVerified,
		ExpectedUpdatedAt:   input.ExpectedUpdatedAt,
	})
	if err != nil {
		return nil, err
	}
	return &LobeHubUserChange{Before: detail.User, After: *after}, nil
}

func BanLobeHubUser(id string, reason string, expires *time.Time) (*LobeHubUserActionResult, error) {
	reason = strings.TrimSpace(reason)
	if id == "" || reason == "" || len(reason) > 500 {
		return nil, ErrInvalidLobeHubUserInput
	}
	if expires != nil && !expires.After(time.Now()) {
		return nil, ErrInvalidLobeHubUserInput
	}
	user, revoked, err := model.SetLobeHubUserBan(id, &reason, expires, true)
	if err != nil {
		return nil, err
	}
	return &LobeHubUserActionResult{User: *user, Revoked: revoked}, nil
}

func UnbanLobeHubUser(id string) (*LobeHubUserActionResult, error) {
	if strings.TrimSpace(id) == "" {
		return nil, ErrInvalidLobeHubUserInput
	}
	user, revoked, err := model.SetLobeHubUserBan(id, nil, nil, false)
	if err != nil {
		return nil, err
	}
	return &LobeHubUserActionResult{User: *user, Revoked: revoked}, nil
}

func ChangeLobeHubUserRole(id string, role string, confirmOverwriteCustomRole bool) (*LobeHubUserActionResult, error) {
	role = strings.TrimSpace(role)
	if strings.TrimSpace(id) == "" || (role != "user" && role != "admin") {
		return nil, ErrInvalidLobeHubUserInput
	}
	user, revoked, err := model.SetLobeHubUserRole(id, role, confirmOverwriteCustomRole)
	if err != nil {
		return nil, err
	}
	return &LobeHubUserActionResult{User: *user, Revoked: revoked}, nil
}

func ResetLobeHubPassword(id string) (*LobeHubPasswordReset, error) {
	if strings.TrimSpace(id) == "" {
		return nil, ErrInvalidLobeHubUserInput
	}
	temporaryPassword, err := generateLobeHubTemporaryPassword()
	if err != nil {
		return nil, err
	}
	passwordHash, err := hashBetterAuthPassword(temporaryPassword)
	if err != nil {
		return nil, err
	}
	accountID, err := secureRandomString(12, "_-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	if err != nil {
		return nil, err
	}
	reset, err := model.ResetLobeHubUserPassword(id, accountID, passwordHash)
	if err != nil {
		return nil, err
	}
	return &LobeHubPasswordReset{TemporaryPassword: temporaryPassword, Revoked: reset.Revoked}, nil
}

func normalizeLobeHubText(value *string) (*string, error) {
	if value == nil {
		return nil, nil
	}
	normalized := strings.TrimSpace(*value)
	if normalized == "" {
		return nil, nil
	}
	if len(normalized) > 2048 || strings.IndexFunc(normalized, unicode.IsControl) >= 0 {
		return nil, ErrInvalidLobeHubUserInput
	}
	return &normalized, nil
}

func normalizeLobeHubEmail(value *string) (*string, error) {
	normalized, err := normalizeLobeHubText(value)
	if err != nil || normalized == nil {
		return normalized, err
	}
	email := strings.ToLower(*normalized)
	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Address != email || len(email) > 320 {
		return nil, ErrInvalidLobeHubUserInput
	}
	return &email, nil
}

func normalizeLobeHubAvatar(value *string) (*string, error) {
	normalized, err := normalizeLobeHubText(value)
	if err != nil || normalized == nil {
		return normalized, err
	}
	parsed, err := url.ParseRequestURI(*normalized)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return nil, ErrInvalidLobeHubUserInput
	}
	return normalized, nil
}

func equalNullableString(left *string, right *string) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}

func generateLobeHubTemporaryPassword() (string, error) {
	base, err := secureRandomString(20, "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789!@#$%")
	if err != nil {
		return "", err
	}
	characters := []rune(base + "Aa1!")
	for index := len(characters) - 1; index > 0; index-- {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(index+1)))
		if err != nil {
			return "", err
		}
		otherIndex := int(randomIndex.Int64())
		characters[index], characters[otherIndex] = characters[otherIndex], characters[index]
	}
	return string(characters), nil
}

func secureRandomString(length int, alphabet string) (string, error) {
	if length <= 0 || alphabet == "" {
		return "", ErrInvalidLobeHubUserInput
	}
	result := make([]byte, length)
	limit := big.NewInt(int64(len(alphabet)))
	for index := range result {
		value, err := rand.Int(rand.Reader, limit)
		if err != nil {
			return "", err
		}
		result[index] = alphabet[value.Int64()]
	}
	return string(result), nil
}

func hashBetterAuthPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	return hashBetterAuthPasswordWithSalt(password, salt)
}

func hashBetterAuthPasswordWithSalt(password string, salt []byte) (string, error) {
	if len(salt) != 16 {
		return "", ErrInvalidLobeHubUserInput
	}
	saltHex := hex.EncodeToString(salt)
	key, err := scrypt.Key([]byte(norm.NFKC.String(password)), []byte(saltHex), 16384, 16, 1, 64)
	if err != nil {
		return "", fmt.Errorf("hash BetterAuth password: %w", err)
	}
	return saltHex + ":" + hex.EncodeToString(key), nil
}
