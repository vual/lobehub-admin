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
	"time"
	"unicode"

	"github.com/QuantumNous/new-api/model"
)

var ErrInvalidLobeHubKnowledgeBaseInput = errors.New("invalid LobeHub knowledge base input")

type LobeHubKnowledgeBaseUpdateInput struct {
	Name              string
	Description       *string
	Avatar            *string
	ExpectedUpdatedAt time.Time
}

type LobeHubKnowledgeBaseChange struct {
	Before model.LobeHubKnowledgeBase `json:"before"`
	After  model.LobeHubKnowledgeBase `json:"after"`
}

func UpdateLobeHubKnowledgeBase(id string, input LobeHubKnowledgeBaseUpdateInput) (*LobeHubKnowledgeBaseChange, error) {
	id = strings.TrimSpace(id)
	name := strings.TrimSpace(input.Name)
	if id == "" || name == "" || len(name) > 255 || input.ExpectedUpdatedAt.IsZero() || strings.IndexFunc(name, unicode.IsControl) >= 0 {
		return nil, ErrInvalidLobeHubKnowledgeBaseInput
	}
	description, err := normalizeLobeHubKnowledgeBaseText(input.Description, 4096)
	if err != nil {
		return nil, err
	}
	avatar, err := normalizeLobeHubKnowledgeBaseText(input.Avatar, 2048)
	if err != nil {
		return nil, err
	}
	before, err := model.GetLobeHubKnowledgeBase(id)
	if err != nil {
		return nil, err
	}
	after, err := model.UpdateLobeHubKnowledgeBase(id, model.LobeHubKnowledgeBaseUpdate{
		Name: name, Description: description, Avatar: avatar, ExpectedUpdatedAt: input.ExpectedUpdatedAt,
	})
	if err != nil {
		return nil, err
	}
	return &LobeHubKnowledgeBaseChange{Before: *before, After: *after}, nil
}

func normalizeLobeHubKnowledgeBaseText(value *string, maximum int) (*string, error) {
	if value == nil {
		return nil, nil
	}
	normalized := strings.TrimSpace(*value)
	if normalized == "" {
		return nil, nil
	}
	if len(normalized) > maximum || strings.IndexFunc(normalized, unicode.IsControl) >= 0 {
		return nil, ErrInvalidLobeHubKnowledgeBaseInput
	}
	return &normalized, nil
}
