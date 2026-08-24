export type LobeHubConversationUser = {
  id: string
  username: string | null
  email: string | null
  full_name: string | null
  avatar: string | null
}

export type LobeHubConversationSource = {
  id: string | null
  type: 'agent' | 'group' | 'unknown'
  name: string | null
  avatar: string | null
}

export type LobeHubConversation = {
  id: string
  title: string | null
  type: 'agent' | 'group' | 'unknown'
  status: string | null
  trigger: string | null
  mode: string | null
  user: LobeHubConversationUser
  source: LobeHubConversationSource
  message_count: number
  model: string | null
  provider: string | null
  total_cost: string | null
  total_tokens: number | null
  created_at: string
  updated_at: string
}

export type LobeHubConversationAttachment = {
  id: string
  name: string
  file_type: string
  size: number
  url: string
}

export type LobeHubConversationPlugin = {
  tool_call_id: string | null
  type: string | null
  api_name: string | null
  arguments: string | null
  identifier: string | null
  intervention: unknown
  state: unknown
  error: unknown
}

export type LobeHubConversationTTS = {
  content_md5: string | null
  file_id: string | null
  voice: string | null
}

export type LobeHubConversationTranslation = {
  content: string | null
  from: string | null
  to: string | null
}

export type LobeHubConversationActor = {
  id: string | null
  name: string | null
  avatar: string | null
  role: string
}

export type LobeHubConversationMessage = {
  id: string
  role: string
  content: string | null
  editor_data: unknown
  summary: string | null
  reasoning: unknown
  search: unknown
  metadata: unknown
  usage: unknown
  error: unknown
  tools: unknown
  model: string | null
  provider: string | null
  actor: LobeHubConversationActor
  plugin: LobeHubConversationPlugin | null
  translation: LobeHubConversationTranslation | null
  thread: unknown
  message_group: unknown
  queries: Array<{ id: string; rewrite_query: string | null; user_query: string | null }>
  tts: LobeHubConversationTTS | null
  attachments: LobeHubConversationAttachment[]
  thread_id: string | null
  parent_id: string | null
  message_group_id: string | null
  target_id: string | null
  created_at: string
  updated_at: string
}

export type LobeHubConversationList = {
  items: LobeHubConversation[]
  total: number
  page: number
  page_size: number
}

export type LobeHubConversationFilters = {
  statuses: string[]
  triggers: string[]
  models: string[]
  providers: string[]
}

export type LobeHubConversationMessages = {
  conversation: LobeHubConversation
  items: LobeHubConversationMessage[]
  has_more: boolean
  next_cursor?: string
}

export type LobeHubConversationListParams = {
  page?: number
  page_size?: number
  q?: string
  type?: string
  status?: string
  trigger?: string
  model?: string
  provider?: string
  updated_from?: string
  updated_to?: string
  sort_by?: 'updated_at' | 'created_at' | 'message_count' | 'total_tokens' | 'total_cost'
  sort_order?: 'asc' | 'desc'
}

export type LobeHubApiResponse<T> = {
  success: boolean
  code?: string
  message?: string
  data?: T
}
