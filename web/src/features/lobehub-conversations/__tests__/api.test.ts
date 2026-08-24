import { beforeEach, describe, expect, test, vi } from 'vitest'

import {
  getLobeHubConversationFilters,
  getLobeHubConversationMessages,
  getLobeHubConversations,
} from '../api'

const { get } = vi.hoisted(() => ({ get: vi.fn() }))

vi.mock('@/lib/api', () => ({ api: { get } }))

describe('LobeHub conversation API contract', () => {
  beforeEach(() => {
    get.mockResolvedValue({ data: { success: true, data: {} } })
  })

  test('passes list filters, pagination and sorting to the server', async () => {
    await getLobeHubConversations({
      page: 2,
      page_size: 20,
      q: 'alice',
      type: 'group',
      status: 'completed',
      model: 'gpt-4o',
      provider: 'openai',
      sort_by: 'updated_at',
      sort_order: 'desc',
    })

    expect(get).toHaveBeenCalledWith('/api/lobehub/conversations', {
      params: {
        page: 2,
        page_size: 20,
        q: 'alice',
        type: 'group',
        status: 'completed',
        model: 'gpt-4o',
        provider: 'openai',
        sort_by: 'updated_at',
        sort_order: 'desc',
      },
    })
  })

  test('requests filters from the dedicated endpoint', async () => {
    await getLobeHubConversationFilters()

    expect(get).toHaveBeenCalledWith('/api/lobehub/conversations/filters')
  })

  test('encodes topic IDs and sends the older-message cursor', async () => {
    await getLobeHubConversationMessages('topic/one', 'cursor/value')

    expect(get).toHaveBeenCalledWith(
      '/api/lobehub/conversations/topic%2Fone/messages',
      { params: { cursor: 'cursor/value' } }
    )
  })
})
