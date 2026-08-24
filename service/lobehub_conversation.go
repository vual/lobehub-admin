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
package service

import (
	"errors"
	"strings"

	"github.com/QuantumNous/new-api/model"
)

var ErrInvalidLobeHubConversationInput = errors.New("invalid LobeHub conversation input")

func ListLobeHubConversations(params model.LobeHubConversationListParams) (*model.LobeHubConversationList, error) {
	if params.Page < 1 || params.PageSize < 1 || params.PageSize > 100 {
		return nil, ErrInvalidLobeHubConversationInput
	}
	if params.Type != "" && params.Type != "agent" && params.Type != "group" && params.Type != "unknown" {
		return nil, ErrInvalidLobeHubConversationInput
	}
	if params.SortBy != "" && params.SortBy != "updated_at" && params.SortBy != "created_at" && params.SortBy != "message_count" && params.SortBy != "total_tokens" && params.SortBy != "total_cost" {
		return nil, ErrInvalidLobeHubConversationInput
	}
	if params.SortOrder != "" && !strings.EqualFold(params.SortOrder, "asc") && !strings.EqualFold(params.SortOrder, "desc") {
		return nil, ErrInvalidLobeHubConversationInput
	}
	if params.UpdatedFrom != nil && params.UpdatedTo != nil && !params.UpdatedFrom.Before(*params.UpdatedTo) {
		return nil, ErrInvalidLobeHubConversationInput
	}
	return model.ListLobeHubConversations(params)
}

func ListLobeHubConversationFilters() (*model.LobeHubConversationFilters, error) {
	return model.ListLobeHubConversationFilters()
}

func ListLobeHubConversationMessages(id string, cursor string) (*model.LobeHubConversationMessages, error) {
	if strings.TrimSpace(id) == "" {
		return nil, ErrInvalidLobeHubConversationInput
	}
	decoded, err := model.DecodeLobeHubConversationCursor(cursor)
	if err != nil {
		return nil, ErrInvalidLobeHubConversationInput
	}
	return model.ListLobeHubConversationMessages(strings.TrimSpace(id), decoded)
}
