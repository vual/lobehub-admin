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

	"github.com/QuantumNous/new-api/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLobeHubRAGStatusPriority(t *testing.T) {
	tests := []struct {
		name                                string
		files, chunks, embeddings           int64
		hasError, hasProcessing, tasksReady bool
		want                                string
	}{
		{name: "error wins", files: 1, chunks: 1, embeddings: 1, hasError: true, hasProcessing: true, tasksReady: true, want: "error"},
		{name: "processing wins over empty", hasProcessing: true, want: "processing"},
		{name: "empty", want: "empty"},
		{name: "ready", files: 1, chunks: 3, embeddings: 3, tasksReady: true, want: "ready"},
		{name: "unindexed", files: 1, chunks: 3, embeddings: 2, tasksReady: true, want: "unindexed"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, lobeHubRAGStatus(test.files, test.chunks, test.embeddings, test.hasError, test.hasProcessing, test.tasksReady))
		})
	}
}

func TestLobeHubKnowledgeBaseWorkspaceMapping(t *testing.T) {
	personal := lobeHubKnowledgeBaseRow{ID: "personal", OwnerID: "owner"}.item()
	assert.Nil(t, personal.Workspace)

	workspaceID, slug, name := "workspace-id", "team", "Team"
	workspace := lobeHubKnowledgeBaseRow{ID: "workspace", OwnerID: "owner", WorkspaceID: &workspaceID, WorkspaceSlug: &slug, WorkspaceName: &name}.item()
	require.NotNil(t, workspace.Workspace)
	assert.Equal(t, workspaceID, workspace.Workspace.ID)
	assert.Equal(t, name, workspace.Workspace.Name)
}

func TestLobeHubKnowledgeBaseChunkNeverSerializesVector(t *testing.T) {
	payload, err := common.Marshal(LobeHubKnowledgeBaseChunk{ID: "chunk", HasEmbedding: true})
	require.NoError(t, err)
	assert.NotContains(t, string(payload), "embeddings")
	assert.NotContains(t, string(payload), "vector")
}
