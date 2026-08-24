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
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
)

var (
	ErrLobeHubConversationNotFound = errors.New("LobeHub conversation not found")
	ErrLobeHubConversationCursor   = errors.New("invalid LobeHub conversation cursor")
)

var lobeHubConversationRequiredColumns = map[string][]string{
	"users": {
		"id", "username", "email", "full_name", "avatar",
	},
	"topics": {
		"id", "title", "session_id", "agent_id", "group_id", "user_id", "trigger", "mode", "status",
		"total_cost", "total_tokens", "model", "provider", "workspace_id", "created_at", "updated_at",
	},
	"messages": {
		"id", "role", "content", "editor_data", "summary", "reasoning", "search", "metadata", "usage", "error", "tools",
		"model", "provider", "user_id", "topic_id", "thread_id", "parent_id", "agent_id", "group_id",
		"target_id", "message_group_id", "workspace_id", "created_at", "updated_at",
	},
	"threads":            {"id", "title", "content", "editor_data", "type", "status", "topic_id", "source_message_id", "parent_thread_id", "agent_id", "group_id", "metadata", "user_id", "created_at", "updated_at"},
	"message_groups":     {"id", "topic_id", "user_id", "parent_group_id", "parent_message_id", "title", "description", "type", "content", "editor_data", "metadata", "created_at", "updated_at"},
	"agents":             {"id", "title", "name", "avatar"},
	"chat_groups":        {"id", "title", "avatar"},
	"message_plugins":    {"id", "tool_call_id", "type", "api_name", "arguments", "identifier", "intervention", "state", "error"},
	"message_translates": {"id", "content", "from", "to"},
	"message_queries":    {"id", "message_id", "rewrite_query", "user_query"},
	"message_tts":        {"id", "content_md5", "file_id", "voice"},
	"messages_files":     {"message_id", "file_id"},
	"files":              {"id", "name", "file_type", "size", "url", "created_at"},
}

// LobeHubJSON keeps JSONB values opaque so the admin can display the exact
// persisted payload without interpreting provider-specific fields.
type LobeHubJSON json.RawMessage

func (j *LobeHubJSON) Scan(value interface{}) error {
	switch v := value.(type) {
	case nil:
		*j = nil
		return nil
	case []byte:
		*j = append((*j)[:0], v...)
		return nil
	case string:
		*j = LobeHubJSON([]byte(v))
		return nil
	default:
		return fmt.Errorf("unsupported JSON value type %T", value)
	}
}

func (j LobeHubJSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return []byte(j), nil
}

func (j LobeHubJSON) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return []byte(j), nil
}

type LobeHubConversationUser struct {
	ID       string  `json:"id"`
	Username *string `json:"username"`
	Email    *string `json:"email"`
	FullName *string `json:"full_name"`
	Avatar   *string `json:"avatar"`
}

type LobeHubConversationSource struct {
	ID     *string `json:"id"`
	Type   string  `json:"type"`
	Name   *string `json:"name"`
	Avatar *string `json:"avatar"`
}

type LobeHubConversationItem struct {
	ID           string                    `json:"id"`
	Title        *string                   `json:"title"`
	Type         string                    `json:"type"`
	Status       *string                   `json:"status"`
	Trigger      *string                   `json:"trigger"`
	Mode         *string                   `json:"mode"`
	User         LobeHubConversationUser   `json:"user"`
	Source       LobeHubConversationSource `json:"source"`
	MessageCount int64                     `json:"message_count"`
	Model        *string                   `json:"model"`
	Provider     *string                   `json:"provider"`
	TotalCost    *string                   `json:"total_cost"`
	TotalTokens  *int64                    `json:"total_tokens"`
	CreatedAt    time.Time                 `json:"created_at"`
	UpdatedAt    time.Time                 `json:"updated_at"`
}

type LobeHubConversationAttachment struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	FileType string `json:"file_type"`
	Size     int64  `json:"size"`
	URL      string `json:"url"`
}

type LobeHubConversationActor struct {
	ID     *string `json:"id"`
	Name   *string `json:"name"`
	Avatar *string `json:"avatar"`
	Role   string  `json:"role"`
}

type LobeHubConversationPlugin struct {
	ToolCallID   *string     `json:"tool_call_id"`
	Type         *string     `json:"type"`
	APIName      *string     `json:"api_name"`
	Arguments    *string     `json:"arguments"`
	Identifier   *string     `json:"identifier"`
	Intervention LobeHubJSON `json:"intervention"`
	State        LobeHubJSON `json:"state"`
	Error        LobeHubJSON `json:"error"`
}

type LobeHubConversationThread struct {
	ID              string      `json:"id"`
	Title           *string     `json:"title"`
	Content         *string     `json:"content"`
	EditorData      LobeHubJSON `json:"editor_data"`
	Type            string      `json:"type"`
	Status          *string     `json:"status"`
	SourceMessageID *string     `json:"source_message_id"`
	ParentThreadID  *string     `json:"parent_thread_id"`
	Metadata        LobeHubJSON `json:"metadata"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type LobeHubConversationMessageGroup struct {
	ID              string      `json:"id"`
	TopicID         *string     `json:"topic_id"`
	ParentGroupID   *string     `json:"parent_group_id"`
	ParentMessageID *string     `json:"parent_message_id"`
	Title           *string     `json:"title"`
	Description     *string     `json:"description"`
	Type            *string     `json:"type"`
	Content         *string     `json:"content"`
	EditorData      LobeHubJSON `json:"editor_data"`
	Metadata        LobeHubJSON `json:"metadata"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type LobeHubConversationQuery struct {
	ID           string  `json:"id"`
	RewriteQuery *string `json:"rewrite_query"`
	UserQuery    *string `json:"user_query"`
}

type LobeHubConversationTTS struct {
	ContentMD5 *string `json:"content_md5"`
	FileID     *string `json:"file_id"`
	Voice      *string `json:"voice"`
}

type LobeHubConversationTranslation struct {
	Content *string `json:"content"`
	From    *string `json:"from"`
	To      *string `json:"to"`
}

type LobeHubConversationMessage struct {
	ID             string                           `json:"id"`
	Role           string                           `json:"role"`
	Content        *string                          `json:"content"`
	EditorData     LobeHubJSON                      `json:"editor_data"`
	Summary        *string                          `json:"summary"`
	Reasoning      LobeHubJSON                      `json:"reasoning"`
	Search         LobeHubJSON                      `json:"search"`
	Metadata       LobeHubJSON                      `json:"metadata"`
	Usage          LobeHubJSON                      `json:"usage"`
	Error          LobeHubJSON                      `json:"error"`
	Tools          LobeHubJSON                      `json:"tools"`
	Model          *string                          `json:"model"`
	Provider       *string                          `json:"provider"`
	Actor          LobeHubConversationActor         `json:"actor"`
	Plugin         *LobeHubConversationPlugin       `json:"plugin"`
	Translation    *LobeHubConversationTranslation  `json:"translation"`
	Thread         *LobeHubConversationThread       `json:"thread"`
	MessageGroup   *LobeHubConversationMessageGroup `json:"message_group"`
	Queries        []LobeHubConversationQuery       `json:"queries"`
	TTS            *LobeHubConversationTTS          `json:"tts"`
	Attachments    []LobeHubConversationAttachment  `json:"attachments"`
	ThreadID       *string                          `json:"thread_id"`
	ParentID       *string                          `json:"parent_id"`
	MessageGroupID *string                          `json:"message_group_id"`
	TargetID       *string                          `json:"target_id"`
	CreatedAt      time.Time                        `json:"created_at"`
	UpdatedAt      time.Time                        `json:"updated_at"`
}

type LobeHubConversationListParams struct {
	Page        int
	PageSize    int
	Query       string
	Type        string
	Status      string
	Trigger     string
	Model       string
	Provider    string
	UpdatedFrom *time.Time
	UpdatedTo   *time.Time
	SortBy      string
	SortOrder   string
}

type LobeHubConversationList struct {
	Items    []LobeHubConversationItem `json:"items"`
	Total    int64                     `json:"total"`
	Page     int                       `json:"page"`
	PageSize int                       `json:"page_size"`
}

type LobeHubConversationFilters struct {
	Statuses  []string `json:"statuses"`
	Triggers  []string `json:"triggers"`
	Models    []string `json:"models"`
	Providers []string `json:"providers"`
}

type LobeHubConversationMessages struct {
	Conversation LobeHubConversationItem      `json:"conversation"`
	Items        []LobeHubConversationMessage `json:"items"`
	HasMore      bool                         `json:"has_more"`
	NextCursor   string                       `json:"next_cursor,omitempty"`
}

type conversationRow struct {
	ID            string    `gorm:"column:id"`
	Title         *string   `gorm:"column:title"`
	Type          string    `gorm:"column:type"`
	Status        *string   `gorm:"column:status"`
	Trigger       *string   `gorm:"column:trigger"`
	Mode          *string   `gorm:"column:mode"`
	OwnerID       string    `gorm:"column:owner_id"`
	OwnerUsername *string   `gorm:"column:owner_username"`
	OwnerEmail    *string   `gorm:"column:owner_email"`
	OwnerFullName *string   `gorm:"column:owner_full_name"`
	OwnerAvatar   *string   `gorm:"column:owner_avatar"`
	SourceID      *string   `gorm:"column:source_id"`
	SourceName    *string   `gorm:"column:source_name"`
	SourceAvatar  *string   `gorm:"column:source_avatar"`
	MessageCount  int64     `gorm:"column:message_count"`
	Model         *string   `gorm:"column:model"`
	Provider      *string   `gorm:"column:provider"`
	TotalCost     *string   `gorm:"column:total_cost"`
	TotalTokens   *int64    `gorm:"column:total_tokens"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

type messageRow struct {
	ID                          string      `gorm:"column:id"`
	Role                        string      `gorm:"column:role"`
	Content                     *string     `gorm:"column:content"`
	EditorData                  LobeHubJSON `gorm:"column:editor_data"`
	Summary                     *string     `gorm:"column:summary"`
	Reasoning                   LobeHubJSON `gorm:"column:reasoning"`
	Search                      LobeHubJSON `gorm:"column:search"`
	Metadata                    LobeHubJSON `gorm:"column:metadata"`
	Usage                       LobeHubJSON `gorm:"column:usage"`
	Error                       LobeHubJSON `gorm:"column:error"`
	Tools                       LobeHubJSON `gorm:"column:tools"`
	Model                       *string     `gorm:"column:model"`
	Provider                    *string     `gorm:"column:provider"`
	MessageUserID               *string     `gorm:"column:message_user_id"`
	MessageUserName             *string     `gorm:"column:message_user_name"`
	MessageUserAvatar           *string     `gorm:"column:message_user_avatar"`
	MessageAgentID              *string     `gorm:"column:message_agent_id"`
	MessageAgentName            *string     `gorm:"column:message_agent_name"`
	MessageAgentAvatar          *string     `gorm:"column:message_agent_avatar"`
	PluginID                    *string     `gorm:"column:plugin_id"`
	PluginToolCallID            *string     `gorm:"column:plugin_tool_call_id"`
	PluginType                  *string     `gorm:"column:plugin_type"`
	PluginAPIName               *string     `gorm:"column:plugin_api_name"`
	PluginArguments             *string     `gorm:"column:plugin_arguments"`
	PluginIdentifier            *string     `gorm:"column:plugin_identifier"`
	PluginIntervention          LobeHubJSON `gorm:"column:plugin_intervention"`
	PluginState                 LobeHubJSON `gorm:"column:plugin_state"`
	PluginError                 LobeHubJSON `gorm:"column:plugin_error"`
	TranslationID               *string     `gorm:"column:translation_id"`
	TranslationContent          *string     `gorm:"column:translation_content"`
	TranslationFrom             *string     `gorm:"column:translation_from"`
	TranslationTo               *string     `gorm:"column:translation_to"`
	ThreadID                    *string     `gorm:"column:thread_id"`
	ThreadTitle                 *string     `gorm:"column:thread_title"`
	ThreadContent               *string     `gorm:"column:thread_content"`
	ThreadEditorData            LobeHubJSON `gorm:"column:thread_editor_data"`
	ThreadType                  *string     `gorm:"column:thread_type"`
	ThreadStatus                *string     `gorm:"column:thread_status"`
	ThreadSourceMessageID       *string     `gorm:"column:thread_source_message_id"`
	ThreadParentID              *string     `gorm:"column:thread_parent_id"`
	ThreadMetadata              LobeHubJSON `gorm:"column:thread_metadata"`
	ThreadCreatedAt             *time.Time  `gorm:"column:thread_created_at"`
	ThreadUpdatedAt             *time.Time  `gorm:"column:thread_updated_at"`
	MessageGroupID              *string     `gorm:"column:message_group_id"`
	MessageGroupTopicID         *string     `gorm:"column:message_group_topic_id"`
	MessageGroupParentID        *string     `gorm:"column:message_group_parent_id"`
	MessageGroupParentMessageID *string     `gorm:"column:message_group_parent_message_id"`
	MessageGroupTitle           *string     `gorm:"column:message_group_title"`
	MessageGroupDescription     *string     `gorm:"column:message_group_description"`
	MessageGroupType            *string     `gorm:"column:message_group_type"`
	MessageGroupContent         *string     `gorm:"column:message_group_content"`
	MessageGroupEditorData      LobeHubJSON `gorm:"column:message_group_editor_data"`
	MessageGroupMetadata        LobeHubJSON `gorm:"column:message_group_metadata"`
	MessageGroupCreatedAt       *time.Time  `gorm:"column:message_group_created_at"`
	MessageGroupUpdatedAt       *time.Time  `gorm:"column:message_group_updated_at"`
	TTSContentMD5               *string     `gorm:"column:tts_content_md5"`
	TTSFileID                   *string     `gorm:"column:tts_file_id"`
	TTSVoice                    *string     `gorm:"column:tts_voice"`
	ParentID                    *string     `gorm:"column:parent_id"`
	TargetID                    *string     `gorm:"column:target_id"`
	CreatedAt                   time.Time   `gorm:"column:created_at"`
	UpdatedAt                   time.Time   `gorm:"column:updated_at"`
}

type attachmentRow struct {
	MessageID string `gorm:"column:message_id"`
	ID        string `gorm:"column:id"`
	Name      string `gorm:"column:name"`
	FileType  string `gorm:"column:file_type"`
	Size      int64  `gorm:"column:size"`
	URL       string `gorm:"column:url"`
}

type queryRow struct {
	MessageID    string  `gorm:"column:message_id"`
	ID           string  `gorm:"column:id"`
	RewriteQuery *string `gorm:"column:rewrite_query"`
	UserQuery    *string `gorm:"column:user_query"`
}

func mapLobeHubMessageActor(row messageRow) LobeHubConversationActor {
	actor := LobeHubConversationActor{Role: row.Role}
	if row.Role == "user" {
		actor.ID = row.MessageUserID
		actor.Name = row.MessageUserName
		actor.Avatar = row.MessageUserAvatar
		return actor
	}
	if row.MessageAgentID != nil {
		actor.ID = row.MessageAgentID
		actor.Name = row.MessageAgentName
		actor.Avatar = row.MessageAgentAvatar
	}
	return actor
}

func conversationSelect(schema string) string {
	topics := lobeHubTable(schema, "topics")
	users := lobeHubTable(schema, "users")
	agents := lobeHubTable(schema, "agents")
	groups := lobeHubTable(schema, "chat_groups")
	messages := lobeHubTable(schema, "messages")
	return fmt.Sprintf(`
		t.id, t.title,
		CASE WHEN t.group_id IS NOT NULL THEN 'group' WHEN t.agent_id IS NOT NULL THEN 'agent' ELSE 'unknown' END AS type,
		t.status, t.trigger, t.mode,
		u.id AS owner_id, u.username AS owner_username, u.email AS owner_email,
		u.full_name AS owner_full_name, u.avatar AS owner_avatar,
		CASE WHEN t.group_id IS NOT NULL THEN cg.id ELSE a.id END AS source_id,
		CASE WHEN t.group_id IS NOT NULL THEN cg.title ELSE COALESCE(NULLIF(a.name, ''), a.title) END AS source_name,
		CASE WHEN t.group_id IS NOT NULL THEN cg.avatar ELSE a.avatar END AS source_avatar,
		(SELECT COUNT(*) FROM %s m WHERE m.topic_id = t.id) AS message_count,
		t.model, t.provider, t.total_cost, t.total_tokens, t.created_at, t.updated_at
		FROM %s t
		JOIN %s u ON u.id = t.user_id
		LEFT JOIN %s a ON a.id = t.agent_id
		LEFT JOIN %s cg ON cg.id = t.group_id`, messages, topics, users, agents, groups)
}

func mapConversation(row conversationRow) LobeHubConversationItem {
	return LobeHubConversationItem{
		ID: row.ID, Title: row.Title, Type: row.Type, Status: row.Status, Trigger: row.Trigger, Mode: row.Mode,
		User:         LobeHubConversationUser{ID: row.OwnerID, Username: row.OwnerUsername, Email: row.OwnerEmail, FullName: row.OwnerFullName, Avatar: row.OwnerAvatar},
		Source:       LobeHubConversationSource{ID: row.SourceID, Type: row.Type, Name: row.SourceName, Avatar: row.SourceAvatar},
		MessageCount: row.MessageCount, Model: row.Model, Provider: row.Provider, TotalCost: row.TotalCost, TotalTokens: row.TotalTokens,
		CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	}
}

func buildConversationWhere(params LobeHubConversationListParams) (string, []interface{}) {
	where := []string{"1 = 1"}
	args := make([]interface{}, 0)
	if query := strings.TrimSpace(params.Query); query != "" {
		pattern := "%" + escapeLobeHubSearch(query) + "%"
		where = append(where, `(t.id ILIKE ? ESCAPE '\' OR t.title ILIKE ? ESCAPE '\' OR u.id ILIKE ? ESCAPE '\' OR u.username ILIKE ? ESCAPE '\' OR u.email ILIKE ? ESCAPE '\' OR u.full_name ILIKE ? ESCAPE '\' OR a.id ILIKE ? ESCAPE '\' OR a.title ILIKE ? ESCAPE '\' OR a.name ILIKE ? ESCAPE '\' OR cg.id ILIKE ? ESCAPE '\' OR cg.title ILIKE ? ESCAPE '\')`)
		for range 11 {
			args = append(args, pattern)
		}
	}
	switch params.Type {
	case "agent":
		where = append(where, "t.agent_id IS NOT NULL AND t.group_id IS NULL")
	case "group":
		where = append(where, "t.group_id IS NOT NULL")
	case "unknown":
		where = append(where, "t.agent_id IS NULL AND t.group_id IS NULL")
	}
	if value := strings.TrimSpace(params.Status); value != "" {
		where = append(where, "t.status = ?")
		args = append(args, value)
	}
	if value := strings.TrimSpace(params.Trigger); value != "" {
		where = append(where, "t.trigger = ?")
		args = append(args, value)
	}
	if value := strings.TrimSpace(params.Model); value != "" {
		where = append(where, "t.model = ?")
		args = append(args, value)
	}
	if value := strings.TrimSpace(params.Provider); value != "" {
		where = append(where, "t.provider = ?")
		args = append(args, value)
	}
	if params.UpdatedFrom != nil {
		where = append(where, "t.updated_at >= ?")
		args = append(args, *params.UpdatedFrom)
	}
	if params.UpdatedTo != nil {
		where = append(where, "t.updated_at < ?")
		args = append(args, *params.UpdatedTo)
	}
	return strings.Join(where, " AND "), args
}

func ListLobeHubConversations(params LobeHubConversationListParams) (*LobeHubConversationList, error) {
	schema, err := checkLobeHubSchema(lobeHubConversationRequiredColumns)
	if err != nil {
		return nil, err
	}
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 20
	}
	whereSQL, args := buildConversationWhere(params)
	base := conversationSelect(schema)
	fromIndex := strings.LastIndex(base, "FROM ")
	if fromIndex < 0 {
		return nil, ErrLobeHubSchemaUnavailable
	}
	fromSQL := base[fromIndex:]
	var total int64
	if err := DB.Raw("SELECT COUNT(*) "+fromSQL+" WHERE "+whereSQL, args...).Scan(&total).Error; err != nil {
		return nil, err
	}
	sortColumns := map[string]string{
		"updated_at": "t.updated_at", "created_at": "t.created_at", "message_count": "message_count",
		"total_tokens": "t.total_tokens", "total_cost": "t.total_cost",
	}
	sortColumn := sortColumns[params.SortBy]
	if sortColumn == "" {
		sortColumn = "t.updated_at"
	}
	sortOrder := "DESC"
	if strings.EqualFold(params.SortOrder, "asc") {
		sortOrder = "ASC"
	}
	query := "SELECT " + base
	query += " WHERE " + whereSQL + " ORDER BY " + sortColumn + " " + sortOrder + " NULLS LAST, t.id ASC LIMIT ? OFFSET ?"
	queryArgs := append(append([]interface{}{}, args...), params.PageSize, (params.Page-1)*params.PageSize)
	rows := make([]conversationRow, 0)
	if err := DB.Raw(query, queryArgs...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	items := make([]LobeHubConversationItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapConversation(row))
	}
	return &LobeHubConversationList{Items: items, Total: total, Page: params.Page, PageSize: params.PageSize}, nil
}

func ListLobeHubConversationFilters() (*LobeHubConversationFilters, error) {
	schema, err := checkLobeHubSchema(lobeHubConversationRequiredColumns)
	if err != nil {
		return nil, err
	}
	table := lobeHubTable(schema, "topics")
	result := &LobeHubConversationFilters{Statuses: []string{}, Triggers: []string{}, Models: []string{}, Providers: []string{}}
	queries := []struct {
		column string
		target *[]string
	}{
		{"status", &result.Statuses}, {"trigger", &result.Triggers}, {"model", &result.Models}, {"provider", &result.Providers},
	}
	for _, query := range queries {
		if err := DB.Raw("SELECT DISTINCT " + query.column + " FROM " + table + " WHERE " + query.column + " IS NOT NULL AND BTRIM(" + query.column + ") <> '' ORDER BY " + query.column).Scan(query.target).Error; err != nil {
			return nil, err
		}
	}
	return result, nil
}

func getLobeHubConversation(id string) (*LobeHubConversationItem, error) {
	schema, err := checkLobeHubSchema(lobeHubConversationRequiredColumns)
	if err != nil {
		return nil, err
	}
	base := conversationSelect(schema)
	var row conversationRow
	result := DB.Raw("SELECT "+base+" WHERE t.id = ?", id).Scan(&row)
	if result.Error != nil {
		return nil, result.Error
	}
	if row.ID == "" {
		return nil, ErrLobeHubConversationNotFound
	}
	item := mapConversation(row)
	return &item, nil
}

type LobeHubConversationCursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
}

func EncodeLobeHubConversationCursor(cursor LobeHubConversationCursor) (string, error) {
	if cursor.ID == "" || cursor.CreatedAt.IsZero() {
		return "", ErrLobeHubConversationCursor
	}
	data, err := common.Marshal(cursor)
	if err != nil {
		return "", ErrLobeHubConversationCursor
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}

func DecodeLobeHubConversationCursor(value string) (LobeHubConversationCursor, error) {
	if strings.TrimSpace(value) == "" {
		return LobeHubConversationCursor{}, nil
	}
	data, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return LobeHubConversationCursor{}, ErrLobeHubConversationCursor
	}
	var cursor LobeHubConversationCursor
	if err := common.Unmarshal(data, &cursor); err != nil || cursor.ID == "" || cursor.CreatedAt.IsZero() {
		return LobeHubConversationCursor{}, ErrLobeHubConversationCursor
	}
	return cursor, nil
}

func messageSelect(schema string) string {
	messages := lobeHubTable(schema, "messages")
	agents := lobeHubTable(schema, "agents")
	users := lobeHubTable(schema, "users")
	threads := lobeHubTable(schema, "threads")
	messageGroups := lobeHubTable(schema, "message_groups")
	plugins := lobeHubTable(schema, "message_plugins")
	translates := lobeHubTable(schema, "message_translates")
	tts := lobeHubTable(schema, "message_tts")
	return fmt.Sprintf(`
		m.id, m.role, m.content, m.editor_data, m.summary, m.reasoning, m.search, m.metadata, m.usage, m.error, m.tools,
		m.model, m.provider, m.thread_id, m.parent_id, m.message_group_id, m.target_id, m.created_at, m.updated_at,
		u.id AS message_user_id, COALESCE(u.full_name, u.username, u.email) AS message_user_name,
		u.avatar AS message_user_avatar,
		a.id AS message_agent_id, COALESCE(NULLIF(a.name, ''), a.title) AS message_agent_name,
		a.avatar AS message_agent_avatar,
		p.id AS plugin_id, p.tool_call_id AS plugin_tool_call_id, p.type AS plugin_type, p.api_name AS plugin_api_name,
		p.arguments AS plugin_arguments, p.identifier AS plugin_identifier, p.intervention AS plugin_intervention,
		p.state AS plugin_state, p.error AS plugin_error,
		tr.id AS translation_id, tr.content AS translation_content, tr."from" AS translation_from, tr."to" AS translation_to,
		th.title AS thread_title, th.content AS thread_content, th.editor_data AS thread_editor_data,
		th.type AS thread_type, th.status AS thread_status, th.source_message_id AS thread_source_message_id,
		th.parent_thread_id AS thread_parent_id, th.metadata AS thread_metadata, th.created_at AS thread_created_at,
		th.updated_at AS thread_updated_at,
		mg.topic_id AS message_group_topic_id, mg.parent_group_id AS message_group_parent_id,
		mg.parent_message_id AS message_group_parent_message_id, mg.title AS message_group_title,
		mg.description AS message_group_description, mg.type AS message_group_type, mg.content AS message_group_content,
		mg.editor_data AS message_group_editor_data, mg.metadata AS message_group_metadata,
		mg.created_at AS message_group_created_at, mg.updated_at AS message_group_updated_at,
		tts.content_md5 AS tts_content_md5, tts.file_id AS tts_file_id, tts.voice AS tts_voice
		FROM %s m
		LEFT JOIN %s a ON a.id = m.agent_id
		LEFT JOIN %s u ON u.id = m.user_id
		LEFT JOIN %s th ON th.id = m.thread_id
		LEFT JOIN %s mg ON mg.id = m.message_group_id
		LEFT JOIN %s p ON p.id = m.id
		LEFT JOIN %s tr ON tr.id = m.id
		LEFT JOIN %s tts ON tts.id = m.id`, messages, agents, users, threads, messageGroups, plugins, translates, tts)
}

func ListLobeHubConversationMessages(id string, cursor LobeHubConversationCursor) (*LobeHubConversationMessages, error) {
	conversation, err := getLobeHubConversation(id)
	if err != nil {
		return nil, err
	}
	schema, err := checkLobeHubSchema(lobeHubConversationRequiredColumns)
	if err != nil {
		return nil, err
	}
	base := messageSelect(schema)
	query := "SELECT " + base
	args := []interface{}{id}
	where := "m.topic_id = ?"
	if !cursor.CreatedAt.IsZero() {
		where += " AND (m.created_at < ? OR (m.created_at = ? AND m.id < ?))"
		args = append(args, cursor.CreatedAt, cursor.CreatedAt, cursor.ID)
	}
	query += " WHERE " + where + " ORDER BY m.created_at DESC, m.id DESC LIMIT 21"
	rows := make([]messageRow, 0)
	if err := DB.Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	hasMore := len(rows) > 20
	if hasMore {
		rows = rows[:20]
	}
	for left, right := 0, len(rows)-1; left < right; left, right = left+1, right-1 {
		rows[left], rows[right] = rows[right], rows[left]
	}

	items := make([]LobeHubConversationMessage, 0, len(rows))
	messageIDs := make([]string, 0, len(rows))
	for _, row := range rows {
		messageIDs = append(messageIDs, row.ID)
		message := LobeHubConversationMessage{
			ID: row.ID, Role: row.Role, Content: row.Content, EditorData: row.EditorData, Summary: row.Summary, Reasoning: row.Reasoning,
			Search: row.Search, Metadata: row.Metadata, Usage: row.Usage, Error: row.Error, Tools: row.Tools,
			Model: row.Model, Provider: row.Provider,
			Actor:    mapLobeHubMessageActor(row),
			ThreadID: row.ThreadID, ParentID: row.ParentID, MessageGroupID: row.MessageGroupID, TargetID: row.TargetID,
			Attachments: []LobeHubConversationAttachment{}, Queries: []LobeHubConversationQuery{}, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
		}
		if row.PluginID != nil {
			message.Plugin = &LobeHubConversationPlugin{ToolCallID: row.PluginToolCallID, Type: row.PluginType, APIName: row.PluginAPIName, Arguments: row.PluginArguments, Identifier: row.PluginIdentifier, Intervention: row.PluginIntervention, State: row.PluginState, Error: row.PluginError}
		}
		if row.TranslationID != nil {
			message.Translation = &LobeHubConversationTranslation{Content: row.TranslationContent, From: row.TranslationFrom, To: row.TranslationTo}
		}
		if row.ThreadID != nil {
			thread := &LobeHubConversationThread{ID: *row.ThreadID, Title: row.ThreadTitle, Content: row.ThreadContent, EditorData: row.ThreadEditorData, Status: row.ThreadStatus, SourceMessageID: row.ThreadSourceMessageID, ParentThreadID: row.ThreadParentID, Metadata: row.ThreadMetadata}
			if row.ThreadType != nil {
				thread.Type = *row.ThreadType
			}
			if row.ThreadCreatedAt != nil {
				thread.CreatedAt = *row.ThreadCreatedAt
			}
			if row.ThreadUpdatedAt != nil {
				thread.UpdatedAt = *row.ThreadUpdatedAt
			}
			message.Thread = thread
		}
		if row.MessageGroupID != nil {
			group := &LobeHubConversationMessageGroup{ID: *row.MessageGroupID, TopicID: row.MessageGroupTopicID, ParentGroupID: row.MessageGroupParentID, ParentMessageID: row.MessageGroupParentMessageID, Title: row.MessageGroupTitle, Description: row.MessageGroupDescription, Type: row.MessageGroupType, Content: row.MessageGroupContent, EditorData: row.MessageGroupEditorData, Metadata: row.MessageGroupMetadata}
			if row.MessageGroupCreatedAt != nil {
				group.CreatedAt = *row.MessageGroupCreatedAt
			}
			if row.MessageGroupUpdatedAt != nil {
				group.UpdatedAt = *row.MessageGroupUpdatedAt
			}
			message.MessageGroup = group
		}
		if row.TTSContentMD5 != nil || row.TTSFileID != nil || row.TTSVoice != nil {
			message.TTS = &LobeHubConversationTTS{ContentMD5: row.TTSContentMD5, FileID: row.TTSFileID, Voice: row.TTSVoice}
		}
		items = append(items, message)
	}
	if len(messageIDs) > 0 {
		attachments := make([]attachmentRow, 0)
		attachmentsTable := lobeHubTable(schema, "messages_files")
		filesTable := lobeHubTable(schema, "files")
		if err := DB.Raw("SELECT mf.message_id, f.id, f.name, f.file_type, f.size, f.url FROM "+attachmentsTable+" mf JOIN "+filesTable+" f ON f.id = mf.file_id WHERE mf.message_id IN ? ORDER BY mf.message_id, f.created_at, f.id", messageIDs).Scan(&attachments).Error; err != nil {
			return nil, err
		}
		byMessage := make(map[string][]LobeHubConversationAttachment)
		for _, attachment := range attachments {
			byMessage[attachment.MessageID] = append(byMessage[attachment.MessageID], LobeHubConversationAttachment{ID: attachment.ID, Name: attachment.Name, FileType: attachment.FileType, Size: attachment.Size, URL: attachment.URL})
		}
		for index := range items {
			items[index].Attachments = byMessage[items[index].ID]
			if items[index].Attachments == nil {
				items[index].Attachments = []LobeHubConversationAttachment{}
			}
		}
	}
	if len(messageIDs) > 0 {
		queriesTable := lobeHubTable(schema, "message_queries")
		queries := make([]queryRow, 0)
		if err := DB.Raw("SELECT message_id, id, rewrite_query, user_query FROM "+queriesTable+" WHERE message_id IN ? ORDER BY message_id, id", messageIDs).Scan(&queries).Error; err != nil {
			return nil, err
		}
		byMessage := make(map[string][]LobeHubConversationQuery)
		for _, query := range queries {
			byMessage[query.MessageID] = append(byMessage[query.MessageID], LobeHubConversationQuery{ID: query.ID, RewriteQuery: query.RewriteQuery, UserQuery: query.UserQuery})
		}
		for index := range items {
			if values := byMessage[items[index].ID]; values != nil {
				items[index].Queries = values
			}
		}
	}
	result := &LobeHubConversationMessages{Conversation: *conversation, Items: items, HasMore: hasMore}
	if hasMore && len(items) > 0 {
		result.NextCursor, err = EncodeLobeHubConversationCursor(LobeHubConversationCursor{CreatedAt: items[0].CreatedAt, ID: items[0].ID})
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
