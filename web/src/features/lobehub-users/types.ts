/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/

export type LobeHubUser = {
  id: string
  username: string | null
  email: string | null
  avatar: string | null
  phone: string | null
  first_name: string | null
  last_name: string | null
  full_name: string | null
  email_verified: boolean
  phone_number_verified: boolean
  role: string | null
  banned: boolean
  ban_reason: string | null
  ban_expires: string | null
  two_factor_enabled: boolean
  last_active_at: string
  created_at: string
  updated_at: string
  password_set: boolean
  session_count: number
}

export type LobeHubLoginProvider = {
  provider_id: string
  created_at: string
}

export type LobeHubUserDetail = {
  user: LobeHubUser
  providers: LobeHubLoginProvider[]
}

export type LobeHubApiResponse<T> = {
  success: boolean
  code?: string
  message?: string
  data?: T
}

export type LobeHubUserList = {
  items: LobeHubUser[]
  total: number
  page: number
  page_size: number
}

export type LobeHubUserSortBy =
  | 'created_at'
  | 'last_active_at'
  | 'username'
  | 'email'

export type LobeHubUserListParams = {
  page?: number
  page_size?: number
  q?: string
  status?: string
  role?: string
  email_verified?: boolean
  two_factor?: boolean
  sort_by?: LobeHubUserSortBy
  sort_order?: 'asc' | 'desc'
}

export type LobeHubUserUpdatePayload = {
  username: string | null
  email: string | null
  avatar: string | null
  phone: string | null
  first_name: string | null
  last_name: string | null
  full_name: string | null
  email_verified: boolean
  phone_number_verified: boolean
  expected_updated_at: string
}

export type LobeHubRevokedCredentials = {
  sessions: number
  oidc: number
}

export type LobeHubUserActionResult = {
  user: LobeHubUser
  revoked: LobeHubRevokedCredentials
}

export type LobeHubPasswordResetResult = {
  temporary_password: string
  revoked: LobeHubRevokedCredentials
}
