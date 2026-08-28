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
import { beforeEach, describe, expect, test, vi } from 'vitest'

import {
  getKnowledgeBaseChunks,
  getKnowledgeBases,
  updateKnowledgeBase,
} from './api'

const { get, patch } = vi.hoisted(() => ({ get: vi.fn(), patch: vi.fn() }))
vi.mock('@/lib/api', () => ({ api: { get, patch } }))

describe('LobeHub knowledge base API contract', () => {
  beforeEach(() => {
    get.mockResolvedValue({ data: { success: true, data: { items: [] } } })
    patch.mockResolvedValue({ data: { success: true } })
  })

  test('passes server-side search, filters, sorting and pagination', async () => {
    const params = {
      page: 2,
      page_size: 50,
      q: 'team',
      scope: 'workspace',
      workspace_id: 'ws/1',
      visibility: 'private',
      rag_status: 'unindexed',
      sort_by: 'total_size',
      sort_order: 'desc',
    }
    await getKnowledgeBases(params)
    expect(get).toHaveBeenCalledWith('/api/lobehub/knowledge-bases', { params })
  })

  test('updates only safe metadata with the optimistic-lock timestamp', async () => {
    const payload = {
      name: 'Docs',
      description: null,
      avatar: null,
      expected_updated_at: '2026-08-28T00:00:00Z',
    }
    await updateKnowledgeBase('kb/one', payload)
    expect(patch).toHaveBeenCalledWith(
      '/api/lobehub/knowledge-bases/kb%2Fone',
      payload
    )
    expect(Object.keys(payload).sort()).toEqual([
      'avatar',
      'description',
      'expected_updated_at',
      'name',
    ])
  })

  test('uses a knowledge-base-scoped file chunks endpoint', async () => {
    await getKnowledgeBaseChunks('kb/one', 'file/two', 3)
    expect(get).toHaveBeenCalledWith(
      '/api/lobehub/knowledge-bases/kb%2Fone/files/file%2Ftwo/chunks',
      { params: { page: 3, page_size: 20 } }
    )
  })
})
