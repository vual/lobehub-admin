import type { ColumnDef } from '@tanstack/react-table'
import { Eye } from 'lucide-react'
import { useTranslation } from 'react-i18next'

import { LongText } from '@/components/long-text'
import { StatusBadge } from '@/components/status-badge'
import { TableId } from '@/components/table-id'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'

import type { LobeHubConversation } from '../types'

type ConversationsColumnsOptions = {
  onView: (conversation: LobeHubConversation) => void
}

function formatDate(value: string): string {
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? '-' : date.toLocaleString()
}

function getDisplayName(conversation: LobeHubConversation): string {
  return (
    conversation.title?.trim() ||
    conversation.source.name?.trim() ||
    conversation.id
  )
}

function getStatusVariant(status: string | null) {
  if (status === 'completed' || status === 'active') return 'success' as const
  if (status === 'failed') return 'danger' as const
  if (status === 'running') return 'info' as const
  return 'neutral' as const
}

function getSourceTypeLabel(
  type: LobeHubConversation['type'],
  translate: (key: string) => string
): string {
  if (type === 'group') return translate('LobeHub Group')
  if (type === 'agent') return translate('LobeHub Agent')
  return translate('LobeHub Unknown Source')
}

export function useLobeHubConversationsColumns(
  props: ConversationsColumnsOptions
): ColumnDef<LobeHubConversation>[] {
  const { t } = useTranslation()

  return [
    {
      accessorKey: 'title',
      header: t('LobeHub Conversation'),
      cell: ({ row }) => {
        const conversation = row.original
        const displayName = getDisplayName(conversation)
        return (
          <div className='flex min-w-[250px] items-center gap-3'>
            <Avatar>
              {conversation.source.avatar && (
                <AvatarImage src={conversation.source.avatar} alt={displayName} />
              )}
              <AvatarFallback>
                {displayName.slice(0, 2).toUpperCase()}
              </AvatarFallback>
            </Avatar>
            <div className='min-w-0'>
              <LongText className='max-w-[220px] font-medium'>
                {displayName}
              </LongText>
              <TableId value={conversation.id} className='max-w-[220px] truncate text-xs' />
            </div>
          </div>
        )
      },
      size: 300,
      meta: { mobileTitle: true },
    },
    {
      id: 'user',
      header: t('LobeHub User'),
      cell: ({ row }) => {
        const user = row.original.user
        const displayName = user.full_name || user.username || user.email || user.id
        return (
          <div className='min-w-[190px]'>
            <LongText className='max-w-[210px]'>{displayName}</LongText>
            <LongText className='text-muted-foreground max-w-[210px] text-xs'>
              {user.email || user.id}
            </LongText>
          </div>
        )
      },
      size: 230,
    },
    {
      id: 'source',
      header: t('LobeHub Agent or Group'),
      cell: ({ row }) => {
        const source = row.original.source
        const label = source.name || source.id || t('LobeHub Unknown Source')
        return (
          <div className='min-w-[180px] space-y-1'>
            <div className='font-medium'>{label}</div>
            <StatusBadge
              label={getSourceTypeLabel(source.type, t)}
              variant='neutral'
              copyable={false}
            />
          </div>
        )
      },
      size: 220,
    },
    {
      accessorKey: 'status',
      header: t('Status'),
      cell: ({ row }) => (
        <StatusBadge
          label={row.original.status || '-'}
          variant={getStatusVariant(row.original.status)}
          copyable={false}
        />
      ),
      size: 130,
      meta: { mobileBadge: true },
    },
    {
      id: 'messages',
      header: t('LobeHub Messages'),
      cell: ({ row }) => (
        <span className='text-muted-foreground text-sm'>
          {row.original.message_count}
        </span>
      ),
      size: 100,
      meta: { mobileOrder: 20 },
    },
    {
      id: 'model',
      header: t('Model'),
      cell: ({ row }) => (
        <div className='min-w-[160px] space-y-1'>
          <LongText className='max-w-[180px]'>{row.original.model || '-'}</LongText>
          <LongText className='text-muted-foreground max-w-[180px] text-xs'>
            {row.original.provider || '-'}
          </LongText>
        </div>
      ),
      size: 200,
      meta: { mobileHidden: true },
    },
    {
      accessorKey: 'updated_at',
      header: t('Updated At'),
      cell: ({ row }) => (
        <span className='text-muted-foreground text-sm whitespace-nowrap'>
          {formatDate(row.original.updated_at)}
        </span>
      ),
      size: 190,
      meta: { mobileHidden: true },
    },
    {
      id: 'actions',
      header: () => t('Actions'),
      cell: ({ row }) => (
        <Button
          aria-label={t('LobeHub View Conversation')}
          onClick={() => props.onView(row.original)}
          size='icon-sm'
          type='button'
          variant='ghost'
        >
          <Eye aria-hidden='true' className='size-4' />
        </Button>
      ),
      size: 56,
      enableSorting: false,
      enableHiding: false,
      meta: { pinned: 'right' },
    },
  ]
}
