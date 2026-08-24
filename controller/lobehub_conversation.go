/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

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
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/service"
	"github.com/gin-gonic/gin"
)

func ListLobeHubConversations(c *gin.Context) {
	writeLobeHubNoStore(c)
	params, err := parseLobeHubConversationListParams(c)
	if err != nil {
		writeLobeHubError(c, err)
		return
	}
	result, err := service.ListLobeHubConversations(params)
	if err != nil {
		writeLobeHubError(c, err)
		return
	}
	common.ApiSuccess(c, result)
}

func ListLobeHubConversationFilters(c *gin.Context) {
	writeLobeHubNoStore(c)
	result, err := service.ListLobeHubConversationFilters()
	if err != nil {
		writeLobeHubError(c, err)
		return
	}
	common.ApiSuccess(c, result)
}

func ListLobeHubConversationMessages(c *gin.Context) {
	writeLobeHubNoStore(c)
	result, err := service.ListLobeHubConversationMessages(c.Param("id"), c.Query("cursor"))
	if err != nil {
		writeLobeHubError(c, err)
		return
	}
	common.ApiSuccess(c, result)
}

func parseLobeHubConversationListParams(c *gin.Context) (model.LobeHubConversationListParams, error) {
	params := model.LobeHubConversationListParams{
		Page:      parsePositiveQueryInt(c.Query("page"), 1),
		PageSize:  parsePositiveQueryInt(c.Query("page_size"), 20),
		Query:     c.Query("q"),
		Type:      strings.TrimSpace(c.Query("type")),
		Status:    strings.TrimSpace(c.Query("status")),
		Trigger:   strings.TrimSpace(c.Query("trigger")),
		Model:     strings.TrimSpace(c.Query("model")),
		Provider:  strings.TrimSpace(c.Query("provider")),
		SortBy:    strings.TrimSpace(c.Query("sort_by")),
		SortOrder: strings.ToLower(strings.TrimSpace(c.Query("sort_order"))),
	}
	if params.PageSize > 100 {
		return params, service.ErrInvalidLobeHubConversationInput
	}
	var err error
	if value := strings.TrimSpace(c.Query("updated_from")); value != "" {
		parsed, parseErr := time.Parse(time.RFC3339, value)
		if parseErr != nil {
			err = service.ErrInvalidLobeHubConversationInput
		} else {
			params.UpdatedFrom = &parsed
		}
	}
	if value := strings.TrimSpace(c.Query("updated_to")); value != "" {
		parsed, parseErr := time.Parse(time.RFC3339, value)
		if parseErr != nil {
			err = service.ErrInvalidLobeHubConversationInput
		} else {
			params.UpdatedTo = &parsed
		}
	}
	if err != nil {
		return params, err
	}
	return params, nil
}

func writeLobeHubNoStore(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	c.Header("Pragma", "no-cache")
}
