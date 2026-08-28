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
package model

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrLobeHubKnowledgeBaseNotFound = errors.New("LobeHub knowledge base not found")
	ErrLobeHubKnowledgeBaseStale    = errors.New("the LobeHub knowledge base was changed by another administrator")
)

var lobeHubKnowledgeBaseRequiredColumns = map[string][]string{
	"knowledge_bases": {
		"id", "name", "description", "avatar", "type", "user_id", "client_id", "is_public",
		"settings", "workspace_id", "visibility", "created_at", "updated_at",
	},
	"users":                {"id", "username", "email", "full_name", "avatar"},
	"workspaces":           {"id", "slug", "name"},
	"knowledge_base_files": {"knowledge_base_id", "file_id"},
	"files": {
		"id", "user_id", "file_type", "name", "size", "url", "source", "metadata",
		"chunk_task_id", "embedding_task_id", "visibility", "created_at", "updated_at",
	},
	"documents": {
		"id", "title", "description", "content", "file_type", "filename", "total_char_count",
		"total_line_count", "metadata", "pages", "source_type", "source", "file_id",
		"knowledge_base_id", "parent_id", "user_id", "editor_data", "visibility", "created_at", "updated_at",
	},
	"async_tasks": {"id", "type", "status", "error", "duration", "created_at", "updated_at"},
	"file_chunks": {"file_id", "chunk_id"},
	"chunks":      {"id", "text", "abstract", "metadata", "index", "type", "created_at", "updated_at"},
	"embeddings":  {"chunk_id", "model"},
}

type LobeHubKnowledgeBaseOwner struct {
	ID       string  `json:"id"`
	Username *string `json:"username"`
	Email    *string `json:"email"`
	FullName *string `json:"full_name"`
	Avatar   *string `json:"avatar"`
}

type LobeHubKnowledgeBaseWorkspace struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type LobeHubKnowledgeBaseStats struct {
	FileCount          int64 `json:"file_count"`
	DocumentCount      int64 `json:"document_count"`
	ChunkCount         int64 `json:"chunk_count"`
	EmbeddedChunkCount int64 `json:"embedded_chunk_count"`
	TotalSize          int64 `json:"total_size"`
}

type LobeHubKnowledgeBase struct {
	ID          string                         `json:"id"`
	Name        string                         `json:"name"`
	Description *string                        `json:"description"`
	Avatar      *string                        `json:"avatar"`
	Type        *string                        `json:"type"`
	ClientID    *string                        `json:"client_id"`
	IsPublic    bool                           `json:"is_public"`
	Visibility  string                         `json:"visibility"`
	Settings    LobeHubJSON                    `json:"settings"`
	Owner       LobeHubKnowledgeBaseOwner      `json:"owner"`
	Workspace   *LobeHubKnowledgeBaseWorkspace `json:"workspace"`
	Stats       LobeHubKnowledgeBaseStats      `json:"stats"`
	RAGStatus   string                         `json:"rag_status"`
	CreatedAt   time.Time                      `json:"created_at"`
	UpdatedAt   time.Time                      `json:"updated_at"`
}

type lobeHubKnowledgeBaseRow struct {
	ID                 string      `gorm:"column:id"`
	Name               string      `gorm:"column:name"`
	Description        *string     `gorm:"column:description"`
	Avatar             *string     `gorm:"column:avatar"`
	Type               *string     `gorm:"column:type"`
	ClientID           *string     `gorm:"column:client_id"`
	IsPublic           bool        `gorm:"column:is_public"`
	Visibility         string      `gorm:"column:visibility"`
	Settings           LobeHubJSON `gorm:"column:settings"`
	OwnerID            string      `gorm:"column:owner_id"`
	OwnerUsername      *string     `gorm:"column:owner_username"`
	OwnerEmail         *string     `gorm:"column:owner_email"`
	OwnerFullName      *string     `gorm:"column:owner_full_name"`
	OwnerAvatar        *string     `gorm:"column:owner_avatar"`
	WorkspaceID        *string     `gorm:"column:workspace_id"`
	WorkspaceSlug      *string     `gorm:"column:workspace_slug"`
	WorkspaceName      *string     `gorm:"column:workspace_name"`
	FileCount          int64       `gorm:"column:file_count"`
	DocumentCount      int64       `gorm:"column:document_count"`
	ChunkCount         int64       `gorm:"column:chunk_count"`
	EmbeddedChunkCount int64       `gorm:"column:embedded_chunk_count"`
	TotalSize          int64       `gorm:"column:total_size"`
	RAGStatus          string      `gorm:"column:rag_status"`
	HasError           bool        `gorm:"column:has_error"`
	HasProcessing      bool        `gorm:"column:has_processing"`
	TasksReady         bool        `gorm:"column:tasks_ready"`
	CreatedAt          time.Time   `gorm:"column:created_at"`
	UpdatedAt          time.Time   `gorm:"column:updated_at"`
	TotalCount         int64       `gorm:"column:total_count"`
}

func (row lobeHubKnowledgeBaseRow) item() LobeHubKnowledgeBase {
	var workspace *LobeHubKnowledgeBaseWorkspace
	if row.WorkspaceID != nil {
		workspace = &LobeHubKnowledgeBaseWorkspace{ID: *row.WorkspaceID}
		if row.WorkspaceSlug != nil {
			workspace.Slug = *row.WorkspaceSlug
		}
		if row.WorkspaceName != nil {
			workspace.Name = *row.WorkspaceName
		}
	}
	return LobeHubKnowledgeBase{
		ID: row.ID, Name: row.Name, Description: row.Description, Avatar: row.Avatar, Type: row.Type,
		ClientID: row.ClientID, IsPublic: row.IsPublic, Visibility: row.Visibility, Settings: row.Settings,
		Owner:     LobeHubKnowledgeBaseOwner{ID: row.OwnerID, Username: row.OwnerUsername, Email: row.OwnerEmail, FullName: row.OwnerFullName, Avatar: row.OwnerAvatar},
		Workspace: workspace,
		Stats:     LobeHubKnowledgeBaseStats{FileCount: row.FileCount, DocumentCount: row.DocumentCount, ChunkCount: row.ChunkCount, EmbeddedChunkCount: row.EmbeddedChunkCount, TotalSize: row.TotalSize},
		RAGStatus: lobeHubRAGStatus(row.FileCount, row.ChunkCount, row.EmbeddedChunkCount, row.HasError, row.HasProcessing, row.TasksReady), CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	}
}

func lobeHubRAGStatus(fileCount int64, chunkCount int64, embeddedChunkCount int64, hasError bool, hasProcessing bool, tasksReady bool) string {
	if hasError {
		return "error"
	}
	if hasProcessing {
		return "processing"
	}
	if fileCount == 0 {
		return "empty"
	}
	if tasksReady && chunkCount > 0 && chunkCount == embeddedChunkCount {
		return "ready"
	}
	return "unindexed"
}

type LobeHubKnowledgeBaseListParams struct {
	Page        int
	PageSize    int
	Query       string
	Scope       string
	WorkspaceID string
	Visibility  string
	RAGStatus   string
	SortBy      string
	SortOrder   string
}

type LobeHubKnowledgeBaseList struct {
	Items    []LobeHubKnowledgeBase `json:"items"`
	Total    int64                  `json:"total"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
}

func lobeHubKnowledgeBaseSummarySQL(schema string) string {
	kb := lobeHubTable(schema, "knowledge_bases")
	users := lobeHubTable(schema, "users")
	workspaces := lobeHubTable(schema, "workspaces")
	kbFiles := lobeHubTable(schema, "knowledge_base_files")
	files := lobeHubTable(schema, "files")
	documents := lobeHubTable(schema, "documents")
	fileChunks := lobeHubTable(schema, "file_chunks")
	embeddings := lobeHubTable(schema, "embeddings")
	tasks := lobeHubTable(schema, "async_tasks")
	return fmt.Sprintf(`
		WITH file_base AS (
			SELECT kbf.knowledge_base_id, f.id, f.size,
				COALESCE(BOOL_OR(ct.status = 'error' OR et.status = 'error'), false) AS has_error,
				COALESCE(BOOL_OR(ct.status IN ('pending','processing') OR et.status IN ('pending','processing')), false) AS has_processing,
				COALESCE(BOOL_AND(ct.status = 'success' AND et.status = 'success'), false) AS tasks_ready
			FROM %s kbf
			JOIN %s f ON f.id = kbf.file_id
			LEFT JOIN %s ct ON ct.id = f.chunk_task_id
			LEFT JOIN %s et ON et.id = f.embedding_task_id
			GROUP BY kbf.knowledge_base_id, f.id, f.size
		), file_rollup AS (
			SELECT knowledge_base_id, COUNT(*)::bigint AS file_count,
				COALESCE(SUM(size), 0)::bigint AS total_size,
				BOOL_OR(has_error) AS has_error, BOOL_OR(has_processing) AS has_processing,
				BOOL_AND(tasks_ready) AS tasks_ready
			FROM file_base GROUP BY knowledge_base_id
		), chunk_rollup AS (
			SELECT kbf.knowledge_base_id, COUNT(DISTINCT fc.chunk_id)::bigint AS chunk_count,
				COUNT(DISTINCT e.chunk_id)::bigint AS embedded_chunk_count
			FROM %s kbf JOIN %s fc ON fc.file_id = kbf.file_id
			LEFT JOIN %s e ON e.chunk_id = fc.chunk_id GROUP BY kbf.knowledge_base_id
		), document_links AS (
			SELECT d.knowledge_base_id, d.id FROM %s d WHERE d.knowledge_base_id IS NOT NULL
			UNION
			SELECT kbf.knowledge_base_id, d.id FROM %s kbf JOIN %s d ON d.file_id = kbf.file_id
		), document_rollup AS (
			SELECT knowledge_base_id, COUNT(DISTINCT id)::bigint AS document_count
			FROM document_links GROUP BY knowledge_base_id
		), summary AS (
			SELECT kb.id, kb.name, kb.description, kb.avatar, kb.type, kb.client_id,
				COALESCE(kb.is_public, false) AS is_public, kb.visibility, kb.settings,
				u.id AS owner_id, u.username AS owner_username, u.email AS owner_email,
				u.full_name AS owner_full_name, u.avatar AS owner_avatar,
				kb.workspace_id, w.slug AS workspace_slug, w.name AS workspace_name,
				COALESCE(fr.file_count, 0)::bigint AS file_count,
				COALESCE(dr.document_count, 0)::bigint AS document_count,
				COALESCE(cr.chunk_count, 0)::bigint AS chunk_count,
				COALESCE(cr.embedded_chunk_count, 0)::bigint AS embedded_chunk_count,
				COALESCE(fr.total_size, 0)::bigint AS total_size,
				COALESCE(fr.has_error, false) AS has_error,
				COALESCE(fr.has_processing, false) AS has_processing,
				COALESCE(fr.tasks_ready, false) AS tasks_ready,
				CASE
					WHEN COALESCE(fr.has_error, false) THEN 'error'
					WHEN COALESCE(fr.has_processing, false) THEN 'processing'
					WHEN COALESCE(fr.file_count, 0) = 0 THEN 'empty'
					WHEN fr.tasks_ready AND cr.chunk_count > 0 AND cr.chunk_count = cr.embedded_chunk_count THEN 'ready'
					ELSE 'unindexed'
				END AS rag_status,
				kb.created_at, kb.updated_at
			FROM %s kb
			JOIN %s u ON u.id = kb.user_id
			LEFT JOIN %s w ON w.id = kb.workspace_id
			LEFT JOIN file_rollup fr ON fr.knowledge_base_id = kb.id
			LEFT JOIN chunk_rollup cr ON cr.knowledge_base_id = kb.id
			LEFT JOIN document_rollup dr ON dr.knowledge_base_id = kb.id
		)
	`, kbFiles, files, tasks, tasks, kbFiles, fileChunks, embeddings, documents, kbFiles, documents, kb, users, workspaces)
}

func ListLobeHubKnowledgeBases(params LobeHubKnowledgeBaseListParams) (*LobeHubKnowledgeBaseList, error) {
	schema, err := checkLobeHubSchema(lobeHubKnowledgeBaseRequiredColumns)
	if err != nil {
		return nil, err
	}
	where := []string{"1 = 1"}
	args := make([]interface{}, 0)
	if query := strings.TrimSpace(params.Query); query != "" {
		pattern := "%" + escapeLobeHubSearch(query) + "%"
		where = append(where, `(s.id ILIKE ? ESCAPE '\' OR s.name ILIKE ? ESCAPE '\' OR s.owner_id ILIKE ? ESCAPE '\' OR s.owner_username ILIKE ? ESCAPE '\' OR s.owner_email ILIKE ? ESCAPE '\' OR s.owner_full_name ILIKE ? ESCAPE '\' OR s.workspace_id ILIKE ? ESCAPE '\' OR s.workspace_name ILIKE ? ESCAPE '\')`)
		for range 8 {
			args = append(args, pattern)
		}
	}
	if params.Scope == "personal" {
		where = append(where, "s.workspace_id IS NULL")
	} else if params.Scope == "workspace" {
		where = append(where, "s.workspace_id IS NOT NULL")
	}
	if params.WorkspaceID != "" {
		where = append(where, "s.workspace_id = ?")
		args = append(args, params.WorkspaceID)
	}
	if params.Visibility != "" {
		where = append(where, "s.visibility = ?")
		args = append(args, params.Visibility)
	}
	if params.RAGStatus != "" {
		where = append(where, "s.rag_status = ?")
		args = append(args, params.RAGStatus)
	}
	sortColumns := map[string]string{
		"created_at": "s.created_at", "updated_at": "s.updated_at", "file_count": "s.file_count", "total_size": "s.total_size",
	}
	sortColumn := sortColumns[params.SortBy]
	if sortColumn == "" {
		sortColumn = "s.updated_at"
	}
	sortOrder := "DESC"
	if strings.EqualFold(params.SortOrder, "asc") {
		sortOrder = "ASC"
	}
	querySQL := lobeHubKnowledgeBaseSummarySQL(schema) + `
		SELECT s.*, COUNT(*) OVER()::bigint AS total_count
		FROM summary s WHERE ` + strings.Join(where, " AND ") + `
		ORDER BY ` + sortColumn + ` ` + sortOrder + `, s.id ASC LIMIT ? OFFSET ?`
	args = append(args, params.PageSize, (params.Page-1)*params.PageSize)
	var rows []lobeHubKnowledgeBaseRow
	if err := DB.Raw(querySQL, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	items := make([]LobeHubKnowledgeBase, 0, len(rows))
	for _, row := range rows {
		items = append(items, row.item())
	}
	var total int64
	if len(rows) > 0 {
		total = rows[0].TotalCount
	}
	return &LobeHubKnowledgeBaseList{Items: items, Total: total, Page: params.Page, PageSize: params.PageSize}, nil
}

func GetLobeHubKnowledgeBase(id string) (*LobeHubKnowledgeBase, error) {
	schema, err := checkLobeHubSchema(lobeHubKnowledgeBaseRequiredColumns)
	if err != nil {
		return nil, err
	}
	var row lobeHubKnowledgeBaseRow
	query := lobeHubKnowledgeBaseSummarySQL(schema) + `SELECT s.*, 1::bigint AS total_count FROM summary s WHERE s.id = ? LIMIT 1`
	result := DB.Raw(query, strings.TrimSpace(id)).Scan(&row)
	if result.Error != nil {
		return nil, result.Error
	}
	if row.ID == "" {
		return nil, ErrLobeHubKnowledgeBaseNotFound
	}
	item := row.item()
	return &item, nil
}

type LobeHubKnowledgeBaseWorkspaceFilter struct {
	ID                 string `json:"id" gorm:"column:id"`
	Slug               string `json:"slug" gorm:"column:slug"`
	Name               string `json:"name" gorm:"column:name"`
	KnowledgeBaseCount int64  `json:"knowledge_base_count" gorm:"column:knowledge_base_count"`
}

func ListLobeHubKnowledgeBaseWorkspaces() ([]LobeHubKnowledgeBaseWorkspaceFilter, error) {
	schema, err := checkLobeHubSchema(lobeHubKnowledgeBaseRequiredColumns)
	if err != nil {
		return nil, err
	}
	var items []LobeHubKnowledgeBaseWorkspaceFilter
	query := fmt.Sprintf(`SELECT w.id, w.slug, w.name, COUNT(kb.id)::bigint AS knowledge_base_count
		FROM %s w JOIN %s kb ON kb.workspace_id = w.id
		GROUP BY w.id, w.slug, w.name ORDER BY w.name ASC, w.id ASC`,
		lobeHubTable(schema, "workspaces"), lobeHubTable(schema, "knowledge_bases"))
	if err := DB.Raw(query).Scan(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

type LobeHubKnowledgeBaseFile struct {
	ID                 string      `json:"id" gorm:"column:id"`
	Name               string      `json:"name" gorm:"column:name"`
	FileType           string      `json:"file_type" gorm:"column:file_type"`
	Size               int64       `json:"size" gorm:"column:size"`
	URL                string      `json:"url" gorm:"column:url"`
	Source             *string     `json:"source" gorm:"column:source"`
	Metadata           LobeHubJSON `json:"metadata" gorm:"column:metadata"`
	Visibility         string      `json:"visibility" gorm:"column:visibility"`
	ChunkTaskID        *string     `json:"chunk_task_id" gorm:"column:chunk_task_id"`
	ChunkingStatus     *string     `json:"chunking_status" gorm:"column:chunking_status"`
	ChunkingError      LobeHubJSON `json:"chunking_error" gorm:"column:chunking_error"`
	EmbeddingTaskID    *string     `json:"embedding_task_id" gorm:"column:embedding_task_id"`
	EmbeddingStatus    *string     `json:"embedding_status" gorm:"column:embedding_status"`
	EmbeddingError     LobeHubJSON `json:"embedding_error" gorm:"column:embedding_error"`
	ChunkCount         int64       `json:"chunk_count" gorm:"column:chunk_count"`
	EmbeddedChunkCount int64       `json:"embedded_chunk_count" gorm:"column:embedded_chunk_count"`
	DocumentCount      int64       `json:"document_count" gorm:"column:document_count"`
	RAGStatus          string      `json:"rag_status" gorm:"column:rag_status"`
	CreatedAt          time.Time   `json:"created_at" gorm:"column:created_at"`
	UpdatedAt          time.Time   `json:"updated_at" gorm:"column:updated_at"`
	TotalCount         int64       `json:"-" gorm:"column:total_count"`
}

type LobeHubKnowledgeBaseFileList struct {
	Items    []LobeHubKnowledgeBaseFile `json:"items"`
	Total    int64                      `json:"total"`
	Page     int                        `json:"page"`
	PageSize int                        `json:"page_size"`
}

func ListLobeHubKnowledgeBaseFiles(id string, page int, pageSize int) (*LobeHubKnowledgeBaseFileList, error) {
	schema, err := checkLobeHubSchema(lobeHubKnowledgeBaseRequiredColumns)
	if err != nil {
		return nil, err
	}
	if _, err := GetLobeHubKnowledgeBase(id); err != nil {
		return nil, err
	}
	query := fmt.Sprintf(`
		WITH file_stats AS (
			SELECT fc.file_id, COUNT(DISTINCT fc.chunk_id)::bigint AS chunk_count,
				COUNT(DISTINCT e.chunk_id)::bigint AS embedded_chunk_count
			FROM %s fc LEFT JOIN %s e ON e.chunk_id = fc.chunk_id GROUP BY fc.file_id
		), document_stats AS (
			SELECT file_id, COUNT(*)::bigint AS document_count FROM %s WHERE file_id IS NOT NULL GROUP BY file_id
		)
		SELECT f.id, f.name, f.file_type, f.size::bigint, f.url, f.source, f.metadata, f.visibility,
			f.chunk_task_id, ct.status AS chunking_status, ct.error AS chunking_error,
			f.embedding_task_id, et.status AS embedding_status, et.error AS embedding_error,
			COALESCE(fs.chunk_count, 0)::bigint AS chunk_count,
			COALESCE(fs.embedded_chunk_count, 0)::bigint AS embedded_chunk_count,
			COALESCE(ds.document_count, 0)::bigint AS document_count,
			CASE WHEN ct.status = 'error' OR et.status = 'error' THEN 'error'
				WHEN ct.status IN ('pending','processing') OR et.status IN ('pending','processing') THEN 'processing'
				WHEN ct.status = 'success' AND et.status = 'success' AND COALESCE(fs.chunk_count, 0) > 0 AND fs.chunk_count = fs.embedded_chunk_count THEN 'ready'
				ELSE 'unindexed' END AS rag_status,
			f.created_at, f.updated_at, COUNT(*) OVER()::bigint AS total_count
		FROM %s kbf JOIN %s f ON f.id = kbf.file_id
		LEFT JOIN %s ct ON ct.id = f.chunk_task_id
		LEFT JOIN %s et ON et.id = f.embedding_task_id
		LEFT JOIN file_stats fs ON fs.file_id = f.id
		LEFT JOIN document_stats ds ON ds.file_id = f.id
		WHERE kbf.knowledge_base_id = ? ORDER BY f.updated_at DESC, f.id ASC LIMIT ? OFFSET ?`,
		lobeHubTable(schema, "file_chunks"), lobeHubTable(schema, "embeddings"), lobeHubTable(schema, "documents"),
		lobeHubTable(schema, "knowledge_base_files"), lobeHubTable(schema, "files"),
		lobeHubTable(schema, "async_tasks"), lobeHubTable(schema, "async_tasks"))
	items := make([]LobeHubKnowledgeBaseFile, 0)
	if err := DB.Raw(query, id, pageSize, (page-1)*pageSize).Scan(&items).Error; err != nil {
		return nil, err
	}
	var total int64
	if len(items) > 0 {
		total = items[0].TotalCount
	}
	return &LobeHubKnowledgeBaseFileList{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}

type LobeHubKnowledgeBaseDocument struct {
	ID             string      `json:"id" gorm:"column:id"`
	Title          *string     `json:"title" gorm:"column:title"`
	Description    *string     `json:"description" gorm:"column:description"`
	FileType       string      `json:"file_type" gorm:"column:file_type"`
	Filename       *string     `json:"filename" gorm:"column:filename"`
	TotalCharCount int64       `json:"total_char_count" gorm:"column:total_char_count"`
	TotalLineCount int64       `json:"total_line_count" gorm:"column:total_line_count"`
	SourceType     string      `json:"source_type" gorm:"column:source_type"`
	Source         string      `json:"source" gorm:"column:source"`
	FileID         *string     `json:"file_id" gorm:"column:file_id"`
	ParentID       *string     `json:"parent_id" gorm:"column:parent_id"`
	Visibility     string      `json:"visibility" gorm:"column:visibility"`
	CreatedAt      time.Time   `json:"created_at" gorm:"column:created_at"`
	UpdatedAt      time.Time   `json:"updated_at" gorm:"column:updated_at"`
	Content        *string     `json:"content,omitempty" gorm:"column:content"`
	Metadata       LobeHubJSON `json:"metadata,omitempty" gorm:"column:metadata"`
	Pages          LobeHubJSON `json:"pages,omitempty" gorm:"column:pages"`
	EditorData     LobeHubJSON `json:"editor_data,omitempty" gorm:"column:editor_data"`
	TotalCount     int64       `json:"-" gorm:"column:total_count"`
}

type LobeHubKnowledgeBaseDocumentList struct {
	Items    []LobeHubKnowledgeBaseDocument `json:"items"`
	Total    int64                          `json:"total"`
	Page     int                            `json:"page"`
	PageSize int                            `json:"page_size"`
}

func lobeHubKnowledgeBaseDocumentWhere(schema string) string {
	return fmt.Sprintf(`d.knowledge_base_id = ? OR EXISTS (SELECT 1 FROM %s kbf WHERE kbf.knowledge_base_id = ? AND kbf.file_id = d.file_id)`, lobeHubTable(schema, "knowledge_base_files"))
}

func ListLobeHubKnowledgeBaseDocuments(id string, page int, pageSize int) (*LobeHubKnowledgeBaseDocumentList, error) {
	schema, err := checkLobeHubSchema(lobeHubKnowledgeBaseRequiredColumns)
	if err != nil {
		return nil, err
	}
	if _, err := GetLobeHubKnowledgeBase(id); err != nil {
		return nil, err
	}
	query := fmt.Sprintf(`SELECT d.id, d.title, d.description, d.file_type, d.filename,
		d.total_char_count::bigint, d.total_line_count::bigint, d.source_type, d.source,
		d.file_id, d.parent_id, d.visibility, d.created_at, d.updated_at,
		COUNT(*) OVER()::bigint AS total_count FROM %s d WHERE %s
		ORDER BY d.updated_at DESC, d.id ASC LIMIT ? OFFSET ?`, lobeHubTable(schema, "documents"), lobeHubKnowledgeBaseDocumentWhere(schema))
	items := make([]LobeHubKnowledgeBaseDocument, 0)
	if err := DB.Raw(query, id, id, pageSize, (page-1)*pageSize).Scan(&items).Error; err != nil {
		return nil, err
	}
	var total int64
	if len(items) > 0 {
		total = items[0].TotalCount
	}
	return &LobeHubKnowledgeBaseDocumentList{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}

func GetLobeHubKnowledgeBaseDocument(id string, documentID string) (*LobeHubKnowledgeBaseDocument, error) {
	schema, err := checkLobeHubSchema(lobeHubKnowledgeBaseRequiredColumns)
	if err != nil {
		return nil, err
	}
	var item LobeHubKnowledgeBaseDocument
	query := fmt.Sprintf(`SELECT d.id, d.title, d.description, d.file_type, d.filename,
		d.total_char_count::bigint, d.total_line_count::bigint, d.source_type, d.source,
		d.file_id, d.parent_id, d.visibility, d.created_at, d.updated_at,
		d.content, d.metadata, d.pages, d.editor_data FROM %s d
		WHERE d.id = ? AND (%s) LIMIT 1`, lobeHubTable(schema, "documents"), lobeHubKnowledgeBaseDocumentWhere(schema))
	if err := DB.Raw(query, documentID, id, id).Scan(&item).Error; err != nil {
		return nil, err
	}
	if item.ID == "" {
		return nil, ErrLobeHubKnowledgeBaseNotFound
	}
	return &item, nil
}

type LobeHubKnowledgeBaseChunk struct {
	ID           string      `json:"id" gorm:"column:id"`
	Index        *int        `json:"index" gorm:"column:index"`
	Text         *string     `json:"text" gorm:"column:text"`
	Abstract     *string     `json:"abstract" gorm:"column:abstract"`
	Type         *string     `json:"type" gorm:"column:type"`
	Metadata     LobeHubJSON `json:"metadata" gorm:"column:metadata"`
	HasEmbedding bool        `json:"has_embedding" gorm:"column:has_embedding"`
	Model        *string     `json:"model" gorm:"column:model"`
	CreatedAt    time.Time   `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time   `json:"updated_at" gorm:"column:updated_at"`
	TotalCount   int64       `json:"-" gorm:"column:total_count"`
}

type LobeHubKnowledgeBaseChunkList struct {
	Items    []LobeHubKnowledgeBaseChunk `json:"items"`
	Total    int64                       `json:"total"`
	Page     int                         `json:"page"`
	PageSize int                         `json:"page_size"`
}

func ListLobeHubKnowledgeBaseChunks(id string, fileID string, page int, pageSize int) (*LobeHubKnowledgeBaseChunkList, error) {
	schema, err := checkLobeHubSchema(lobeHubKnowledgeBaseRequiredColumns)
	if err != nil {
		return nil, err
	}
	var linked bool
	linkQuery := fmt.Sprintf(`SELECT EXISTS (SELECT 1 FROM %s WHERE knowledge_base_id = ? AND file_id = ?)`, lobeHubTable(schema, "knowledge_base_files"))
	if err := DB.Raw(linkQuery, id, fileID).Scan(&linked).Error; err != nil {
		return nil, err
	}
	if !linked {
		return nil, ErrLobeHubKnowledgeBaseNotFound
	}
	query := fmt.Sprintf(`SELECT c.id::text AS id, c.index, c.text, c.abstract, c.type, c.metadata,
		(e.chunk_id IS NOT NULL) AS has_embedding, e.model, c.created_at, c.updated_at,
		COUNT(*) OVER()::bigint AS total_count
		FROM %s fc JOIN %s c ON c.id = fc.chunk_id
		LEFT JOIN %s e ON e.chunk_id = c.id
		WHERE fc.file_id = ? ORDER BY c.index ASC NULLS LAST, c.created_at ASC, c.id ASC LIMIT ? OFFSET ?`,
		lobeHubTable(schema, "file_chunks"), lobeHubTable(schema, "chunks"), lobeHubTable(schema, "embeddings"))
	items := make([]LobeHubKnowledgeBaseChunk, 0)
	if err := DB.Raw(query, fileID, pageSize, (page-1)*pageSize).Scan(&items).Error; err != nil {
		return nil, err
	}
	var total int64
	if len(items) > 0 {
		total = items[0].TotalCount
	}
	return &LobeHubKnowledgeBaseChunkList{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}

type LobeHubKnowledgeBaseUpdate struct {
	Name              string
	Description       *string
	Avatar            *string
	ExpectedUpdatedAt time.Time
}

func UpdateLobeHubKnowledgeBase(id string, update LobeHubKnowledgeBaseUpdate) (*LobeHubKnowledgeBase, error) {
	schema, err := checkLobeHubSchema(lobeHubKnowledgeBaseRequiredColumns)
	if err != nil {
		return nil, err
	}
	kb := lobeHubTable(schema, "knowledge_bases")
	var updatedID string
	query := `UPDATE ` + kb + ` SET name = ?, description = ?, avatar = ?, updated_at = NOW()
		WHERE id = ? AND updated_at = ? RETURNING id`
	if err := DB.Raw(query, update.Name, update.Description, update.Avatar, strings.TrimSpace(id), update.ExpectedUpdatedAt).Scan(&updatedID).Error; err != nil {
		return nil, err
	}
	if updatedID == "" {
		var exists bool
		if err := DB.Raw(`SELECT EXISTS (SELECT 1 FROM `+kb+` WHERE id = ?)`, strings.TrimSpace(id)).Scan(&exists).Error; err != nil {
			return nil, err
		}
		if !exists {
			return nil, ErrLobeHubKnowledgeBaseNotFound
		}
		return nil, ErrLobeHubKnowledgeBaseStale
	}
	return GetLobeHubKnowledgeBase(updatedID)
}
