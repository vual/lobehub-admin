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

import { api } from '@/lib/api'

import type {
  LobeHubApiResponse,
  LobeHubPasswordResetResult,
  LobeHubUser,
  LobeHubUserActionResult,
  LobeHubUserDetail,
  LobeHubUserList,
  LobeHubUserListParams,
  LobeHubUserUpdatePayload,
} from './types'

export async function getLobeHubUsers(
  params: LobeHubUserListParams
): Promise<LobeHubApiResponse<LobeHubUserList>> {
  const response = await api.get('/api/lobehub/users', { params })
  return response.data
}

export async function getLobeHubUser(
  id: string
): Promise<LobeHubApiResponse<LobeHubUserDetail>> {
  const response = await api.get(`/api/lobehub/users/${encodeURIComponent(id)}`)
  return response.data
}

export async function updateLobeHubUser(
  id: string,
  payload: LobeHubUserUpdatePayload
): Promise<LobeHubApiResponse<LobeHubUser>> {
  const response = await api.patch(
    `/api/lobehub/users/${encodeURIComponent(id)}`,
    payload
  )
  return response.data
}

export async function banLobeHubUser(
  id: string,
  payload: { reason: string; expires_at: string | null }
): Promise<LobeHubApiResponse<LobeHubUserActionResult>> {
  const response = await api.post(
    `/api/lobehub/users/${encodeURIComponent(id)}/ban`,
    payload
  )
  return response.data
}

export async function unbanLobeHubUser(
  id: string
): Promise<LobeHubApiResponse<LobeHubUserActionResult>> {
  const response = await api.post(
    `/api/lobehub/users/${encodeURIComponent(id)}/unban`
  )
  return response.data
}

export async function changeLobeHubUserRole(
  id: string,
  role: 'user' | 'admin',
  confirmOverwriteCustomRole: boolean
): Promise<LobeHubApiResponse<LobeHubUserActionResult>> {
  const response = await api.patch(
    `/api/lobehub/users/${encodeURIComponent(id)}/role`,
    {
      role,
      confirm_overwrite_custom_role: confirmOverwriteCustomRole,
    }
  )
  return response.data
}

export async function resetLobeHubUserPassword(
  id: string
): Promise<LobeHubApiResponse<LobeHubPasswordResetResult>> {
  const response = await api.post(
    `/api/lobehub/users/${encodeURIComponent(id)}/reset-password`
  )
  return response.data
}
