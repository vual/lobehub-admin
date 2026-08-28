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

export type ApiResponse<T> = {
  success: boolean
  data?: T
  code?: string
  message?: string
}
export type JsonValue = unknown

export type KnowledgeBase = {
  id: string
  name: string
  description: string | null
  avatar: string | null
  type: string | null
  client_id: string | null
  is_public: boolean
  visibility: string
  settings: JsonValue
  owner: {
    id: string
    username: string | null
    email: string | null
    full_name: string | null
    avatar: string | null
  }
  workspace: { id: string; slug: string; name: string } | null
  stats: {
    file_count: number
    document_count: number
    chunk_count: number
    embedded_chunk_count: number
    total_size: number
  }
  rag_status: 'error' | 'processing' | 'empty' | 'ready' | 'unindexed'
  created_at: string
  updated_at: string
}

export type Page<T> = {
  items: T[]
  total: number
  page: number
  page_size: number
}
export type WorkspaceFilter = {
  id: string
  slug: string
  name: string
  knowledge_base_count: number
}

export type KnowledgeBaseFile = {
  id: string
  name: string
  file_type: string
  size: number
  url: string
  source: string | null
  metadata: JsonValue
  visibility: string
  chunk_task_id: string | null
  chunking_status: string | null
  chunking_error: JsonValue
  embedding_task_id: string | null
  embedding_status: string | null
  embedding_error: JsonValue
  chunk_count: number
  embedded_chunk_count: number
  document_count: number
  rag_status: string
  created_at: string
  updated_at: string
}

export type KnowledgeBaseDocument = {
  id: string
  title: string | null
  description: string | null
  file_type: string
  filename: string | null
  total_char_count: number
  total_line_count: number
  source_type: string
  source: string
  file_id: string | null
  parent_id: string | null
  visibility: string
  created_at: string
  updated_at: string
  content?: string | null
  metadata?: JsonValue
  pages?: JsonValue
  editor_data?: JsonValue
}

export type KnowledgeBaseChunk = {
  id: string
  index: number | null
  text: string | null
  abstract: string | null
  type: string | null
  metadata: JsonValue
  has_embedding: boolean
  model: string | null
  created_at: string
  updated_at: string
}
