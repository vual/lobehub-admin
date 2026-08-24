import { api } from '@/lib/api'

import type {
  LobeHubApiResponse,
  LobeHubConversationFilters,
  LobeHubConversationList,
  LobeHubConversationListParams,
  LobeHubConversationMessages,
} from './types'

export async function getLobeHubConversations(
  params: LobeHubConversationListParams
): Promise<LobeHubApiResponse<LobeHubConversationList>> {
  const response = await api.get('/api/lobehub/conversations', { params })
  return response.data
}

export async function getLobeHubConversationFilters(): Promise<
  LobeHubApiResponse<LobeHubConversationFilters>
> {
  const response = await api.get('/api/lobehub/conversations/filters')
  return response.data
}

export async function getLobeHubConversationMessages(
  id: string,
  cursor?: string
): Promise<LobeHubApiResponse<LobeHubConversationMessages>> {
  const response = await api.get(
    `/api/lobehub/conversations/${encodeURIComponent(id)}/messages`,
    { params: cursor ? { cursor } : undefined }
  )
  return response.data
}
