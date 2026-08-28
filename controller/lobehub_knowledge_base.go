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

type lobeHubKnowledgeBaseUpdateRequest struct {
	Name              string    `json:"name"`
	Description       *string   `json:"description"`
	Avatar            *string   `json:"avatar"`
	ExpectedUpdatedAt time.Time `json:"expected_updated_at"`
}

func ListLobeHubKnowledgeBases(c *gin.Context) {
	page, pageSize := parseLobeHubPage(c)
	scope := strings.TrimSpace(c.Query("scope"))
	visibility := strings.TrimSpace(c.Query("visibility"))
	ragStatus := strings.TrimSpace(c.Query("rag_status"))
	sortBy := strings.TrimSpace(c.Query("sort_by"))
	sortOrder := strings.ToLower(strings.TrimSpace(c.Query("sort_order")))
	if !validLobeHubOption(scope, "", "personal", "workspace") ||
		!validLobeHubOption(visibility, "", "private", "public") ||
		!validLobeHubOption(ragStatus, "", "error", "processing", "empty", "ready", "unindexed") ||
		!validLobeHubOption(sortBy, "", "created_at", "updated_at", "file_count", "total_size") ||
		!validLobeHubOption(sortOrder, "", "asc", "desc") {
		writeLobeHubError(c, service.ErrInvalidLobeHubKnowledgeBaseInput)
		return
	}
	result, err := model.ListLobeHubKnowledgeBases(model.LobeHubKnowledgeBaseListParams{
		Page: page, PageSize: pageSize, Query: c.Query("q"), Scope: scope,
		WorkspaceID: strings.TrimSpace(c.Query("workspace_id")), Visibility: visibility,
		RAGStatus: ragStatus, SortBy: sortBy, SortOrder: sortOrder,
	})
	if err != nil {
		writeLobeHubError(c, err)
		return
	}
	common.ApiSuccess(c, result)
}

func ListLobeHubKnowledgeBaseFilters(c *gin.Context) {
	items, err := model.ListLobeHubKnowledgeBaseWorkspaces()
	if err != nil {
		writeLobeHubError(c, err)
		return
	}
	common.ApiSuccess(c, gin.H{"workspaces": items})
}

func GetLobeHubKnowledgeBase(c *gin.Context) {
	item, err := model.GetLobeHubKnowledgeBase(strings.TrimSpace(c.Param("id")))
	if err != nil {
		writeLobeHubError(c, err)
		return
	}
	common.ApiSuccess(c, item)
}

func UpdateLobeHubKnowledgeBase(c *gin.Context) {
	var request lobeHubKnowledgeBaseUpdateRequest
	if err := common.DecodeJson(c.Request.Body, &request); err != nil {
		writeLobeHubError(c, service.ErrInvalidLobeHubKnowledgeBaseInput)
		return
	}
	id := strings.TrimSpace(c.Param("id"))
	change, err := service.UpdateLobeHubKnowledgeBase(id, service.LobeHubKnowledgeBaseUpdateInput{
		Name: request.Name, Description: request.Description, Avatar: request.Avatar, ExpectedUpdatedAt: request.ExpectedUpdatedAt,
	})
	if err != nil {
		writeLobeHubError(c, err)
		return
	}
	changedFields := make([]string, 0, 3)
	if change.Before.Name != change.After.Name {
		changedFields = append(changedFields, "name")
	}
	if !sameStringPointer(change.Before.Description, change.After.Description) {
		changedFields = append(changedFields, "description")
	}
	if !sameStringPointer(change.Before.Avatar, change.After.Avatar) {
		changedFields = append(changedFields, "avatar")
	}
	recordManageAudit(c, "lobehub.knowledge_base.update", map[string]interface{}{
		"lobehub_knowledge_base_id": id, "changed_fields": changedFields,
	})
	common.ApiSuccess(c, change.After)
}

func ListLobeHubKnowledgeBaseFiles(c *gin.Context) {
	page, pageSize := parseLobeHubPage(c)
	items, err := model.ListLobeHubKnowledgeBaseFiles(strings.TrimSpace(c.Param("id")), page, pageSize)
	if err != nil {
		writeLobeHubError(c, err)
		return
	}
	common.ApiSuccess(c, items)
}

func ListLobeHubKnowledgeBaseDocuments(c *gin.Context) {
	page, pageSize := parseLobeHubPage(c)
	items, err := model.ListLobeHubKnowledgeBaseDocuments(strings.TrimSpace(c.Param("id")), page, pageSize)
	if err != nil {
		writeLobeHubError(c, err)
		return
	}
	common.ApiSuccess(c, items)
}

func GetLobeHubKnowledgeBaseDocument(c *gin.Context) {
	item, err := model.GetLobeHubKnowledgeBaseDocument(strings.TrimSpace(c.Param("id")), strings.TrimSpace(c.Param("documentId")))
	if err != nil {
		writeLobeHubError(c, err)
		return
	}
	common.ApiSuccess(c, item)
}

func ListLobeHubKnowledgeBaseChunks(c *gin.Context) {
	page, pageSize := parseLobeHubPage(c)
	items, err := model.ListLobeHubKnowledgeBaseChunks(strings.TrimSpace(c.Param("id")), strings.TrimSpace(c.Param("fileId")), page, pageSize)
	if err != nil {
		writeLobeHubError(c, err)
		return
	}
	common.ApiSuccess(c, items)
}

func parseLobeHubPage(c *gin.Context) (int, int) {
	page := parsePositiveQueryInt(c.Query("page"), 1)
	pageSize := parsePositiveQueryInt(c.Query("page_size"), 20)
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func validLobeHubOption(value string, options ...string) bool {
	for _, option := range options {
		if value == option {
			return true
		}
	}
	return false
}
