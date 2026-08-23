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

import { beforeEach, describe, expect, test, vi } from 'vitest'

import {
  changeLobeHubUserRole,
  getLobeHubUsers,
  resetLobeHubUserPassword,
  updateLobeHubUser,
} from './api'

const { get, patch, post } = vi.hoisted(() => ({
  get: vi.fn(),
  patch: vi.fn(),
  post: vi.fn(),
}))

vi.mock('@/lib/api', () => ({ api: { get, patch, post } }))

describe('LobeHub user API contract', () => {
  beforeEach(() => {
    get.mockResolvedValue({ data: { success: true, data: { items: [] } } })
    patch.mockResolvedValue({ data: { success: true } })
    post.mockResolvedValue({ data: { success: true } })
  })

  test('passes server-side pagination, filtering and sorting parameters', async () => {
    await getLobeHubUsers({
      page: 3,
      page_size: 50,
      q: 'alice',
      status: 'banned',
      role: 'admin',
      sort_by: 'last_active_at',
      sort_order: 'asc',
    })

    expect(get).toHaveBeenCalledWith('/api/lobehub/users', {
      params: {
        page: 3,
        page_size: 50,
        q: 'alice',
        status: 'banned',
        role: 'admin',
        sort_by: 'last_active_at',
        sort_order: 'asc',
      },
    })
  })

  test('sends the optimistic-lock timestamp when updating a user', async () => {
    const payload = {
      username: 'alice',
      email: 'alice@example.com',
      avatar: null,
      phone: null,
      first_name: null,
      last_name: null,
      full_name: 'Alice',
      email_verified: true,
      phone_number_verified: false,
      expected_updated_at: '2026-08-22T00:00:00Z',
    }

    await updateLobeHubUser('user/one', payload)

    expect(patch).toHaveBeenCalledWith('/api/lobehub/users/user%2Fone', payload)
  })

  test('uses the dedicated password reset endpoint without sending a password', async () => {
    await resetLobeHubUserPassword('user_1')

    expect(post).toHaveBeenCalledWith(
      '/api/lobehub/users/user_1/reset-password'
    )
  })

  test('explicitly confirms overwriting a nonstandard role', async () => {
    await changeLobeHubUserRole('user_1', 'admin', true)

    expect(patch).toHaveBeenCalledWith('/api/lobehub/users/user_1/role', {
      role: 'admin',
      confirm_overwrite_custom_role: true,
    })
  })
})
