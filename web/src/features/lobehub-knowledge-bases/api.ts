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
import { api } from '@/lib/api'

import type {
  ApiResponse,
  KnowledgeBase,
  KnowledgeBaseChunk,
  KnowledgeBaseDocument,
  KnowledgeBaseFile,
  Page,
  WorkspaceFilter,
} from './types'

const base = '/api/lobehub/knowledge-bases'
const result = async <T>(request: Promise<{ data: ApiResponse<T> }>) =>
  (await request).data

export const getKnowledgeBases = (params: Record<string, unknown>) =>
  result<Page<KnowledgeBase>>(api.get(base, { params }))
export const getKnowledgeBaseFilters = () =>
  result<{ workspaces: WorkspaceFilter[] }>(api.get(`${base}/filters`))
export const getKnowledgeBase = (id: string) =>
  result<KnowledgeBase>(api.get(`${base}/${encodeURIComponent(id)}`))
export const updateKnowledgeBase = (
  id: string,
  payload: {
    name: string
    description: string | null
    avatar: string | null
    expected_updated_at: string
  }
) =>
  result<KnowledgeBase>(api.patch(`${base}/${encodeURIComponent(id)}`, payload))
export const getKnowledgeBaseFiles = (id: string, page: number) =>
  result<Page<KnowledgeBaseFile>>(
    api.get(`${base}/${encodeURIComponent(id)}/files`, {
      params: { page, page_size: 20 },
    })
  )
export const getKnowledgeBaseDocuments = (id: string, page: number) =>
  result<Page<KnowledgeBaseDocument>>(
    api.get(`${base}/${encodeURIComponent(id)}/documents`, {
      params: { page, page_size: 20 },
    })
  )
export const getKnowledgeBaseDocument = (id: string, documentId: string) =>
  result<KnowledgeBaseDocument>(
    api.get(
      `${base}/${encodeURIComponent(id)}/documents/${encodeURIComponent(documentId)}`
    )
  )
export const getKnowledgeBaseChunks = (
  id: string,
  fileId: string,
  page: number
) =>
  result<Page<KnowledgeBaseChunk>>(
    api.get(
      `${base}/${encodeURIComponent(id)}/files/${encodeURIComponent(fileId)}/chunks`,
      { params: { page, page_size: 20 } }
    )
  )
