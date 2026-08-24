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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLobeHubConversationCursorRoundTrip(t *testing.T) {
	original := LobeHubConversationCursor{
		CreatedAt: time.Date(2026, time.August, 23, 10, 11, 12, 13, time.UTC),
		ID:        "message-001",
	}

	encoded, err := EncodeLobeHubConversationCursor(original)
	require.NoError(t, err)
	decoded, err := DecodeLobeHubConversationCursor(encoded)
	require.NoError(t, err)
	assert.Equal(t, original, decoded)
}

func TestLobeHubConversationCursorRejectsInvalidValues(t *testing.T) {
	_, err := EncodeLobeHubConversationCursor(LobeHubConversationCursor{})
	assert.ErrorIs(t, err, ErrLobeHubConversationCursor)

	_, err = DecodeLobeHubConversationCursor("not-a-cursor")
	assert.ErrorIs(t, err, ErrLobeHubConversationCursor)
}

func TestLobeHubConversationSchemaRequiresReadTables(t *testing.T) {
	for _, table := range []string{"topics", "messages", "agents", "chat_groups", "message_plugins", "message_translates", "messages_files", "files"} {
		_, ok := lobeHubConversationRequiredColumns[table]
		assert.Truef(t, ok, "conversation schema must validate %s", table)
	}
}

func TestMapLobeHubMessageActor(t *testing.T) {
	userID := "user-1"
	userName := "Alice"
	userAvatar := "https://example.com/alice.png"
	agentID := "agent-1"
	agentName := "Assistant"
	agentAvatar := "https://example.com/assistant.png"

	tests := []struct {
		name string
		row  messageRow
		want LobeHubConversationActor
	}{
		{
			name: "user messages use the user even when an agent id is present",
			row: messageRow{
				Role:               "user",
				MessageUserID:      &userID,
				MessageUserName:    &userName,
				MessageUserAvatar:  &userAvatar,
				MessageAgentID:     &agentID,
				MessageAgentName:   &agentName,
				MessageAgentAvatar: &agentAvatar,
			},
			want: LobeHubConversationActor{ID: &userID, Name: &userName, Avatar: &userAvatar, Role: "user"},
		},
		{
			name: "non-user messages use the agent",
			row: messageRow{
				Role:             "assistant",
				MessageUserID:    &userID,
				MessageUserName:  &userName,
				MessageAgentID:   &agentID,
				MessageAgentName: &agentName,
			},
			want: LobeHubConversationActor{ID: &agentID, Name: &agentName, Role: "assistant"},
		},
		{
			name: "unknown roles without an agent are not attributed to the owner",
			row: messageRow{
				Role:            "system",
				MessageUserID:   &userID,
				MessageUserName: &userName,
			},
			want: LobeHubConversationActor{Role: "system"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, mapLobeHubMessageActor(test.row))
		})
	}
}
