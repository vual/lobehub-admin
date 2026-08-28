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
package controller

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/service"
	"github.com/gin-gonic/gin"
)

type lobeHubUserUpdateRequest struct {
	Username            *string   `json:"username"`
	Email               *string   `json:"email"`
	Avatar              *string   `json:"avatar"`
	Phone               *string   `json:"phone"`
	FirstName           *string   `json:"first_name"`
	LastName            *string   `json:"last_name"`
	FullName            *string   `json:"full_name"`
	EmailVerified       bool      `json:"email_verified"`
	PhoneNumberVerified bool      `json:"phone_number_verified"`
	ExpectedUpdatedAt   time.Time `json:"expected_updated_at"`
}

type lobeHubBanRequest struct {
	Reason    string     `json:"reason"`
	ExpiresAt *time.Time `json:"expires_at"`
}

type lobeHubRoleRequest struct {
	Role                       string `json:"role"`
	ConfirmOverwriteCustomRole bool   `json:"confirm_overwrite_custom_role"`
}

func ListLobeHubUsers(c *gin.Context) {
	page := parsePositiveQueryInt(c.Query("page"), 1)
	pageSize := parsePositiveQueryInt(c.Query("page_size"), 20)
	if pageSize > 100 {
		pageSize = 100
	}

	emailVerified, err := parseOptionalBool(c.Query("email_verified"))
	if err != nil {
		writeLobeHubError(c, service.ErrInvalidLobeHubUserInput)
		return
	}
	twoFactor, err := parseOptionalBool(c.Query("two_factor"))
	if err != nil {
		writeLobeHubError(c, service.ErrInvalidLobeHubUserInput)
		return
	}
	status := strings.TrimSpace(c.Query("status"))
	if status != "" && status != "active" && status != "banned" && status != "expired" {
		writeLobeHubError(c, service.ErrInvalidLobeHubUserInput)
		return
	}
	role := strings.TrimSpace(c.Query("role"))
	if role != "" && role != "user" && role != "admin" && role != "other" {
		writeLobeHubError(c, service.ErrInvalidLobeHubUserInput)
		return
	}
	sortBy := strings.TrimSpace(c.Query("sort_by"))
	if sortBy != "" && sortBy != "created_at" && sortBy != "last_active_at" && sortBy != "username" && sortBy != "email" {
		writeLobeHubError(c, service.ErrInvalidLobeHubUserInput)
		return
	}
	sortOrder := strings.ToLower(strings.TrimSpace(c.Query("sort_order")))
	if sortOrder != "" && sortOrder != "asc" && sortOrder != "desc" {
		writeLobeHubError(c, service.ErrInvalidLobeHubUserInput)
		return
	}

	result, err := model.ListLobeHubUsers(model.LobeHubUserListParams{
		Page:          page,
		PageSize:      pageSize,
		Query:         c.Query("q"),
		Status:        status,
		Role:          role,
		EmailVerified: emailVerified,
		TwoFactor:     twoFactor,
		SortBy:        sortBy,
		SortOrder:     sortOrder,
	})
	if err != nil {
		writeLobeHubError(c, err)
		return
	}
	common.ApiSuccess(c, result)
}

func GetLobeHubUser(c *gin.Context) {
	result, err := model.GetLobeHubUser(strings.TrimSpace(c.Param("id")))
	if err != nil {
		writeLobeHubError(c, err)
		return
	}
	common.ApiSuccess(c, result)
}

func UpdateLobeHubUser(c *gin.Context) {
	var request lobeHubUserUpdateRequest
	if err := common.DecodeJson(c.Request.Body, &request); err != nil {
		writeLobeHubError(c, service.ErrInvalidLobeHubUserInput)
		return
	}
	id := strings.TrimSpace(c.Param("id"))
	change, err := service.UpdateLobeHubUser(id, service.LobeHubUserUpdateInput{
		Username:            request.Username,
		Email:               request.Email,
		Avatar:              request.Avatar,
		Phone:               request.Phone,
		FirstName:           request.FirstName,
		LastName:            request.LastName,
		FullName:            request.FullName,
		EmailVerified:       request.EmailVerified,
		PhoneNumberVerified: request.PhoneNumberVerified,
		ExpectedUpdatedAt:   request.ExpectedUpdatedAt,
	})
	if err != nil {
		writeLobeHubError(c, err)
		return
	}
	recordManageAudit(c, "lobehub.user.update", map[string]interface{}{
		"lobehub_user_id": id,
		"changed_fields":  lobeHubChangedFields(change.Before, change.After),
	})
	common.ApiSuccess(c, change.After)
}

func BanLobeHubUser(c *gin.Context) {
	var request lobeHubBanRequest
	if err := common.DecodeJson(c.Request.Body, &request); err != nil {
		writeLobeHubError(c, service.ErrInvalidLobeHubUserInput)
		return
	}
	id := strings.TrimSpace(c.Param("id"))
	result, err := service.BanLobeHubUser(id, request.Reason, request.ExpiresAt)
	if err != nil {
		writeLobeHubError(c, err)
		return
	}
	recordManageAudit(c, "lobehub.user.ban", map[string]interface{}{
		"lobehub_user_id": id,
		"reason":          strings.TrimSpace(request.Reason),
		"expires_at":      request.ExpiresAt,
		"revoked":         result.Revoked,
	})
	common.ApiSuccess(c, result)
}

func UnbanLobeHubUser(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	result, err := service.UnbanLobeHubUser(id)
	if err != nil {
		writeLobeHubError(c, err)
		return
	}
	recordManageAudit(c, "lobehub.user.unban", map[string]interface{}{
		"lobehub_user_id": id,
	})
	common.ApiSuccess(c, result)
}

func ChangeLobeHubUserRole(c *gin.Context) {
	var request lobeHubRoleRequest
	if err := common.DecodeJson(c.Request.Body, &request); err != nil {
		writeLobeHubError(c, service.ErrInvalidLobeHubUserInput)
		return
	}
	id := strings.TrimSpace(c.Param("id"))
	result, err := service.ChangeLobeHubUserRole(id, request.Role, request.ConfirmOverwriteCustomRole)
	if err != nil {
		writeLobeHubError(c, err)
		return
	}
	recordManageAudit(c, "lobehub.user.role_change", map[string]interface{}{
		"lobehub_user_id": id,
		"role":            request.Role,
		"revoked":         result.Revoked,
	})
	common.ApiSuccess(c, result)
}

func ResetLobeHubUserPassword(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	result, err := service.ResetLobeHubPassword(id)
	if err != nil {
		writeLobeHubError(c, err)
		return
	}
	c.Header("Cache-Control", "no-store")
	c.Header("Pragma", "no-cache")
	recordManageAudit(c, "lobehub.user.password_reset", map[string]interface{}{
		"lobehub_user_id": id,
		"revoked":         result.Revoked,
	})
	common.ApiSuccess(c, result)
}

func parsePositiveQueryInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func parseOptionalBool(value string) (*bool, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func writeLobeHubError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	code := "LOBEHUB_INTERNAL_ERROR"
	message := "Failed to manage LobeHub data"
	switch {
	case errors.Is(err, service.ErrInvalidLobeHubKnowledgeBaseInput):
		status, code, message = http.StatusBadRequest, "LOBEHUB_INVALID_INPUT", err.Error()
	case errors.Is(err, service.ErrInvalidLobeHubUserInput):
		status, code, message = http.StatusBadRequest, "LOBEHUB_INVALID_INPUT", err.Error()
	case errors.Is(err, service.ErrInvalidLobeHubConversationInput), errors.Is(err, model.ErrLobeHubConversationCursor):
		status, code, message = http.StatusBadRequest, "LOBEHUB_INVALID_INPUT", err.Error()
	case errors.Is(err, model.ErrLobeHubPostgresRequired):
		status, code, message = http.StatusNotImplemented, "LOBEHUB_POSTGRES_REQUIRED", err.Error()
	case errors.Is(err, model.ErrLobeHubSchemaUnavailable):
		status, code, message = http.StatusServiceUnavailable, "LOBEHUB_SCHEMA_UNAVAILABLE", err.Error()
	case errors.Is(err, model.ErrLobeHubUserNotFound):
		status, code, message = http.StatusNotFound, "LOBEHUB_USER_NOT_FOUND", err.Error()
	case errors.Is(err, model.ErrLobeHubConversationNotFound):
		status, code, message = http.StatusNotFound, "LOBEHUB_CONVERSATION_NOT_FOUND", err.Error()
	case errors.Is(err, model.ErrLobeHubKnowledgeBaseNotFound):
		status, code, message = http.StatusNotFound, "LOBEHUB_KNOWLEDGE_BASE_NOT_FOUND", err.Error()
	case errors.Is(err, model.ErrLobeHubKnowledgeBaseStale):
		status, code, message = http.StatusConflict, "LOBEHUB_KNOWLEDGE_BASE_CONFLICT", err.Error()
	case errors.Is(err, model.ErrLobeHubConflict), errors.Is(err, model.ErrLobeHubStaleUpdate), errors.Is(err, model.ErrLobeHubCredentialState), errors.Is(err, model.ErrLobeHubRoleConfirmation):
		status, code, message = http.StatusConflict, "LOBEHUB_USER_CONFLICT", err.Error()
	}
	c.JSON(status, gin.H{"success": false, "code": code, "message": message})
}

func lobeHubChangedFields(before model.LobeHubUser, after model.LobeHubUser) []string {
	changed := make([]string, 0)
	if !sameStringPointer(before.Username, after.Username) {
		changed = append(changed, "username")
	}
	if !sameStringPointer(before.Email, after.Email) {
		changed = append(changed, "email")
	}
	if !sameStringPointer(before.Avatar, after.Avatar) {
		changed = append(changed, "avatar")
	}
	if !sameStringPointer(before.Phone, after.Phone) {
		changed = append(changed, "phone")
	}
	if !sameStringPointer(before.FirstName, after.FirstName) {
		changed = append(changed, "first_name")
	}
	if !sameStringPointer(before.LastName, after.LastName) {
		changed = append(changed, "last_name")
	}
	if !sameStringPointer(before.FullName, after.FullName) {
		changed = append(changed, "full_name")
	}
	if before.EmailVerified != after.EmailVerified {
		changed = append(changed, "email_verified")
	}
	if before.PhoneNumberVerified != after.PhoneNumberVerified {
		changed = append(changed, "phone_number_verified")
	}
	return changed
}

func sameStringPointer(left *string, right *string) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}
