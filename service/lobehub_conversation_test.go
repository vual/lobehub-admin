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
	"testing"
	"time"

	"github.com/QuantumNous/new-api/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListLobeHubConversationsRejectsUnsupportedSort(t *testing.T) {
	_, err := ListLobeHubConversations(model.LobeHubConversationListParams{
		Page:     1,
		PageSize: 20,
		SortBy:   "content",
	})
	assert.ErrorIs(t, err, ErrInvalidLobeHubConversationInput)
}

func TestListLobeHubConversationsRejectsInvalidDateRange(t *testing.T) {
	from := time.Date(2026, time.August, 24, 0, 0, 0, 0, time.UTC)
	to := from.Add(-time.Minute)
	_, err := ListLobeHubConversations(model.LobeHubConversationListParams{
		Page:        1,
		PageSize:    20,
		UpdatedFrom: &from,
		UpdatedTo:   &to,
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidLobeHubConversationInput)
}

func TestListLobeHubConversationMessagesRejectsBlankID(t *testing.T) {
	_, err := ListLobeHubConversationMessages(" ", "")
	assert.ErrorIs(t, err, ErrInvalidLobeHubConversationInput)
}
