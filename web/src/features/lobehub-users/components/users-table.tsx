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

import { useQuery } from '@tanstack/react-query'
import { getRouteApi } from '@tanstack/react-router'
import type { OnChangeFn, SortingState } from '@tanstack/react-table'
import { useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import {
  DISABLED_ROW_DESKTOP,
  DISABLED_ROW_MOBILE,
  DataTablePage,
  useDataTable,
} from '@/components/data-table'
import { useMediaQuery } from '@/hooks'
import { useTableUrlState } from '@/hooks/use-table-url-state'

import { getLobeHubUsers } from '../api'
import type { LobeHubUser, LobeHubUserSortBy } from '../types'
import { useLobeHubUsersColumns } from './users-columns'

const route = getRouteApi('/_authenticated/lobehub/users/')
const SORTABLE_COLUMNS = new Set<LobeHubUserSortBy>([
  'created_at',
  'last_active_at',
  'username',
  'email',
])

type UsersTableProps = {
  onEdit: (user: LobeHubUser) => void
}

function getBannedRowClassName(user: LobeHubUser, isMobile: boolean) {
  if (!user.banned) return undefined
  return isMobile ? DISABLED_ROW_MOBILE : DISABLED_ROW_DESKTOP
}

export function LobeHubUsersTable({ onEdit }: UsersTableProps) {
  const { t } = useTranslation()
  const columns = useLobeHubUsersColumns({ onEdit })
  const isMobile = useMediaQuery('(max-width: 640px)')
  const [sorting, setSorting] = useState<SortingState>([])

  const {
    globalFilter,
    onGlobalFilterChange,
    columnFilters,
    onColumnFiltersChange,
    pagination,
    onPaginationChange,
    ensurePageInRange,
  } = useTableUrlState({
    search: route.useSearch(),
    navigate: route.useNavigate(),
    pagination: { defaultPage: 1, defaultPageSize: isMobile ? 10 : 20 },
    globalFilter: { enabled: true, key: 'filter' },
    columnFilters: [
      { columnId: 'status', searchKey: 'status', type: 'array' },
      { columnId: 'role', searchKey: 'role', type: 'array' },
      {
        columnId: 'email_verified',
        searchKey: 'emailVerified',
        type: 'array',
      },
      {
        columnId: 'two_factor_enabled',
        searchKey: 'twoFactor',
        type: 'array',
      },
    ],
  })

  const status =
    ((columnFilters.find((filter) => filter.id === 'status')?.value as
      | string[]
      | undefined) ?? [])[0] ?? ''
  const role =
    ((columnFilters.find((filter) => filter.id === 'role')?.value as
      | string[]
      | undefined) ?? [])[0] ?? ''
  const emailVerified =
    ((columnFilters.find((filter) => filter.id === 'email_verified')?.value as
      | string[]
      | undefined) ?? [])[0] ?? ''
  const twoFactor =
    ((columnFilters.find((filter) => filter.id === 'two_factor_enabled')
      ?.value as string[] | undefined) ?? [])[0] ?? ''

  const sortParams = useMemo(() => {
    const activeSort = sorting[0]
    if (
      !activeSort ||
      !SORTABLE_COLUMNS.has(activeSort.id as LobeHubUserSortBy)
    ) {
      return {}
    }
    return {
      sort_by: activeSort.id as LobeHubUserSortBy,
      sort_order: activeSort.desc ? ('desc' as const) : ('asc' as const),
    }
  }, [sorting])

  const handleSortingChange: OnChangeFn<SortingState> = (updater) => {
    setSorting(updater)
    if (pagination.pageIndex > 0) {
      onPaginationChange({ ...pagination, pageIndex: 0 })
    }
  }

  const { data, isLoading, isFetching } = useQuery({
    queryKey: [
      'lobehub-users',
      pagination.pageIndex + 1,
      pagination.pageSize,
      globalFilter,
      status,
      role,
      emailVerified,
      twoFactor,
      sortParams.sort_by,
      sortParams.sort_order,
    ],
    queryFn: async () => {
      const result = await getLobeHubUsers({
        page: pagination.pageIndex + 1,
        page_size: pagination.pageSize,
        q: globalFilter?.trim() || undefined,
        status: status || undefined,
        role: role || undefined,
        email_verified:
          emailVerified === '' ? undefined : emailVerified === 'true',
        two_factor: twoFactor === '' ? undefined : twoFactor === 'true',
        ...sortParams,
      })
      if (!result.success || !result.data) {
        toast.error(result.message || t('LobeHub Failed to load users'))
        return { items: [], total: 0 }
      }
      return { items: result.data.items, total: result.data.total }
    },
    placeholderData: (previousData) => previousData,
  })

  const { table } = useDataTable({
    data: data?.items ?? [],
    columns,
    columnFilters,
    globalFilter,
    pagination,
    sorting,
    initialColumnVisibility: {
      email_verified: false,
      two_factor_enabled: false,
    },
    onPaginationChange,
    onGlobalFilterChange,
    onColumnFiltersChange,
    onSortingChange: handleSortingChange,
    manualPagination: true,
    manualFiltering: true,
    manualSorting: true,
    totalCount: data?.total ?? 0,
    ensurePageInRange,
  })

  return (
    <DataTablePage
      table={table}
      columns={columns}
      isLoading={isLoading}
      isFetching={isFetching}
      emptyTitle={t('LobeHub No Users Found')}
      emptyDescription={t(
        'LobeHub Check the schema or adjust your search and filters.'
      )}
      skeletonKeyPrefix='lobehub-users-skeleton'
      applyHeaderSize
      toolbarProps={{
        searchPlaceholder: t(
          'LobeHub Search by ID, username, email, name or phone...'
        ),
        searchDebounceMs: 500,
        filters: [
          {
            columnId: 'status',
            title: t('Status'),
            singleSelect: true,
            options: [
              { label: t('Active'), value: 'active' },
              { label: t('LobeHub Banned'), value: 'banned' },
              { label: t('LobeHub Ban Expired'), value: 'expired' },
            ],
          },
          {
            columnId: 'role',
            title: t('LobeHub Role'),
            singleSelect: true,
            options: [
              { label: t('LobeHub User'), value: 'user' },
              { label: t('LobeHub Admin'), value: 'admin' },
              { label: t('LobeHub Other'), value: 'other' },
            ],
          },
          {
            columnId: 'email_verified',
            title: t('LobeHub Email Verified'),
            singleSelect: true,
            options: [
              { label: t('LobeHub Verified'), value: 'true' },
              { label: t('LobeHub Unverified'), value: 'false' },
            ],
          },
          {
            columnId: 'two_factor_enabled',
            title: t('LobeHub 2FA Enabled'),
            singleSelect: true,
            options: [
              { label: t('LobeHub Enabled'), value: 'true' },
              { label: t('LobeHub Disabled'), value: 'false' },
            ],
          },
        ],
      }}
      getRowClassName={(row, { isMobile: mobile }) =>
        getBannedRowClassName(row.original, mobile)
      }
    />
  )
}
