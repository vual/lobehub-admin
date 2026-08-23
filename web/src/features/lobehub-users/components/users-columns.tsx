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

import type { ColumnDef } from '@tanstack/react-table'
import { useTranslation } from 'react-i18next'

import { LongText } from '@/components/long-text'
import { StatusBadge } from '@/components/status-badge'
import { TableId } from '@/components/table-id'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'

import type { LobeHubUser } from '../types'
import { UserRowActions } from './user-row-actions'

type UsersColumnsOptions = {
  onEdit: (user: LobeHubUser) => void
}

function formatDate(value: string): string {
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? '-' : date.toLocaleString()
}

function getBanStatus(user: LobeHubUser) {
  if (!user.banned) return 'active'
  if (user.ban_expires && new Date(user.ban_expires).getTime() <= Date.now()) {
    return 'expired'
  }
  return 'banned'
}

function getRoleVariant(role: string) {
  if (role === 'admin') return 'purple' as const
  if (role === 'user') return 'neutral' as const
  return 'warning' as const
}

export function useLobeHubUsersColumns({
  onEdit,
}: UsersColumnsOptions): ColumnDef<LobeHubUser>[] {
  const { t } = useTranslation()

  return [
    {
      accessorKey: 'username',
      header: t('LobeHub User'),
      cell: ({ row }) => {
        const user = row.original
        const displayName =
          user.full_name || user.username || user.email || user.id
        return (
          <div className='flex min-w-[220px] items-center gap-3'>
            <Avatar>
              {user.avatar && (
                <AvatarImage src={user.avatar} alt={displayName} />
              )}
              <AvatarFallback>
                {displayName.slice(0, 2).toUpperCase()}
              </AvatarFallback>
            </Avatar>
            <div className='min-w-0'>
              <LongText className='max-w-[190px] font-medium'>
                {displayName}
              </LongText>
              <TableId
                value={user.id}
                className='max-w-[190px] truncate text-xs'
              />
            </div>
          </div>
        )
      },
      size: 280,
      meta: { mobileTitle: true },
    },
    {
      accessorKey: 'email',
      header: t('LobeHub Contact'),
      cell: ({ row }) => (
        <div className='flex min-w-[190px] flex-col gap-1'>
          <LongText className='max-w-[220px]'>
            {row.original.email || '-'}
          </LongText>
          <LongText className='text-muted-foreground max-w-[220px] text-xs'>
            {row.original.phone || '-'}
          </LongText>
        </div>
      ),
      size: 240,
      meta: { mobileOrder: 20 },
    },
    {
      accessorKey: 'role',
      header: t('LobeHub Role'),
      cell: ({ row }) => {
        const role = row.original.role?.trim() || 'user'
        return (
          <StatusBadge
            label={role}
            variant={getRoleVariant(role)}
            copyable={false}
          />
        )
      },
      enableSorting: false,
      size: 110,
      meta: { mobileBadge: true },
    },
    {
      id: 'security',
      header: t('Security'),
      cell: ({ row }) => {
        const user = row.original
        const values = [
          {
            label: `email_verified: ${String(user.email_verified)}`,
            variant: user.email_verified
              ? ('success' as const)
              : ('neutral' as const),
          },
          {
            label: `two_factor_enabled: ${String(user.two_factor_enabled)}`,
            variant: user.two_factor_enabled
              ? ('info' as const)
              : ('neutral' as const),
          },
          {
            label: `password_set: ${String(user.password_set)}`,
            variant: user.password_set
              ? ('success' as const)
              : ('neutral' as const),
          },
          {
            label: `session_count: ${user.session_count}`,
            variant: 'neutral' as const,
          },
        ]

        return (
          <div className='flex min-w-[210px] flex-wrap gap-1'>
            {values.map((value) => (
              <StatusBadge
                key={value.label}
                label={value.label}
                variant={value.variant}
                copyable={false}
              />
            ))}
          </div>
        )
      },
      enableSorting: false,
      size: 260,
      meta: { mobileOrder: 30 },
    },
    {
      accessorKey: 'email_verified',
      header: t('LobeHub Email Verified'),
      cell: () => null,
      enableSorting: false,
    },
    {
      accessorKey: 'two_factor_enabled',
      header: t('LobeHub 2FA Enabled'),
      cell: () => null,
      enableSorting: false,
    },
    {
      id: 'status',
      header: t('Status'),
      cell: ({ row }) => {
        const status = getBanStatus(row.original)
        const config = {
          active: { label: t('Active'), variant: 'success' as const },
          banned: { label: t('LobeHub Banned'), variant: 'danger' as const },
          expired: {
            label: t('LobeHub Ban Expired'),
            variant: 'warning' as const,
          },
        }[status]
        return (
          <StatusBadge
            label={config.label}
            variant={config.variant}
            copyable={false}
          />
        )
      },
      enableSorting: false,
      size: 130,
      meta: { mobileBadge: true },
    },
    {
      accessorKey: 'last_active_at',
      header: t('LobeHub Last Active'),
      cell: ({ row }) => (
        <span className='text-muted-foreground text-sm whitespace-nowrap'>
          {formatDate(row.original.last_active_at)}
        </span>
      ),
      size: 190,
      meta: { mobileHidden: true },
    },
    {
      accessorKey: 'created_at',
      header: t('Created At'),
      cell: ({ row }) => (
        <span className='text-muted-foreground text-sm whitespace-nowrap'>
          {formatDate(row.original.created_at)}
        </span>
      ),
      size: 190,
      meta: { mobileHidden: true },
    },
    {
      id: 'actions',
      header: () => t('Actions'),
      cell: ({ row }) => <UserRowActions user={row.original} onEdit={onEdit} />,
      size: 56,
      enableSorting: false,
      enableHiding: false,
      meta: { pinned: 'right' },
    },
  ]
}
