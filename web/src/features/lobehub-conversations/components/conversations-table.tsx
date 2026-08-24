import { useQuery } from '@tanstack/react-query'
import { getRouteApi } from '@tanstack/react-router'
import type { OnChangeFn, SortingState } from '@tanstack/react-table'
import { useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import {
  DataTablePage,
  useDataTable,
} from '@/components/data-table'
import { useMediaQuery } from '@/hooks'
import { useTableUrlState } from '@/hooks/use-table-url-state'
import { Input } from '@/components/ui/input'

import {
  getLobeHubConversationFilters,
  getLobeHubConversations,
} from '../api'
import type {
  LobeHubConversation,
  LobeHubConversationListParams,
} from '../types'
import { useLobeHubConversationsColumns } from './conversations-columns'

const route = getRouteApi('/_authenticated/lobehub/conversations/')

type ConversationsTableProps = {
  onView: (conversation: LobeHubConversation) => void
}

export function LobeHubConversationsTable(props: ConversationsTableProps) {
  const { t } = useTranslation()
  const columns = useLobeHubConversationsColumns(props)
  const isMobile = useMediaQuery('(max-width: 640px)')
  const [sorting, setSorting] = useState<SortingState>([])
  const search = route.useSearch()
  const navigate = route.useNavigate()

  const {
    globalFilter,
    onGlobalFilterChange,
    columnFilters,
    onColumnFiltersChange,
    pagination,
    onPaginationChange,
    ensurePageInRange,
  } = useTableUrlState({
    search,
    navigate,
    pagination: { defaultPage: 1, defaultPageSize: isMobile ? 10 : 20 },
    globalFilter: { enabled: true, key: 'filter' },
    columnFilters: [
      { columnId: 'type', searchKey: 'type', type: 'array' },
      { columnId: 'status', searchKey: 'status', type: 'array' },
      { columnId: 'trigger', searchKey: 'trigger', type: 'array' },
      { columnId: 'model', searchKey: 'model', type: 'array' },
      { columnId: 'provider', searchKey: 'provider', type: 'array' },
    ],
  })

  const updatedFrom = search.updatedFrom || ''
  const updatedTo = search.updatedTo || ''

  const updateDateFilter = (key: 'updatedFrom' | 'updatedTo', value: string) => {
    navigate({
      search: (previous) => ({
        ...previous,
        page: undefined,
        [key]: value || undefined,
      }),
    })
  }

  const getFilterValue = (id: string) =>
    ((columnFilters.find((filter) => filter.id === id)?.value as
      | string[]
      | undefined) ?? [])[0] ?? ''

  const type = getFilterValue('type')
  const status = getFilterValue('status')
  const trigger = getFilterValue('trigger')
  const model = getFilterValue('model')
  const provider = getFilterValue('provider')

  const getExclusiveDate = (value: string) => {
    if (!value) return undefined
    const date = new Date(`${value}T00:00:00.000Z`)
    date.setUTCDate(date.getUTCDate() + 1)
    return date.toISOString()
  }

  const sortParams = useMemo<LobeHubConversationListParams>(() => {
    const activeSort = sorting[0]
    if (!activeSort) return {}
    const validSorts = new Set<LobeHubConversationListParams['sort_by']>([
      'updated_at',
      'created_at',
      'message_count',
      'total_tokens',
      'total_cost',
    ])
    if (!validSorts.has(activeSort.id as LobeHubConversationListParams['sort_by'])) {
      return {}
    }
    return {
      sort_by: activeSort.id as LobeHubConversationListParams['sort_by'],
      sort_order: activeSort.desc ? 'desc' : 'asc',
    }
  }, [sorting])

  const { data: filters } = useQuery({
    queryKey: ['lobehub-conversation-filters'],
    queryFn: async () => {
      const result = await getLobeHubConversationFilters()
      if (!result.success || !result.data) {
        throw new Error(result.message || t('LobeHub Failed to load filters'))
      }
      return result.data
    },
    staleTime: 60_000,
  })

  const { data, isLoading, isFetching } = useQuery({
    queryKey: [
      'lobehub-conversations',
      pagination.pageIndex + 1,
      pagination.pageSize,
      globalFilter,
      type,
      status,
      trigger,
      model,
      provider,
      updatedFrom,
      updatedTo,
      sortParams.sort_by,
      sortParams.sort_order,
    ],
    queryFn: async () => {
      const result = await getLobeHubConversations({
        page: pagination.pageIndex + 1,
        page_size: pagination.pageSize,
        q: globalFilter?.trim() || undefined,
        type: type || undefined,
        status: status || undefined,
        trigger: trigger || undefined,
        model: model || undefined,
        provider: provider || undefined,
        updated_from: updatedFrom
          ? new Date(`${updatedFrom}T00:00:00.000Z`).toISOString()
          : undefined,
        updated_to: getExclusiveDate(updatedTo),
        ...sortParams,
      })
      if (!result.success || !result.data) {
        toast.error(result.message || t('LobeHub Failed to load conversations'))
        return { items: [], total: 0 }
      }
      return { items: result.data.items, total: result.data.total }
    },
    placeholderData: (previousData) => previousData,
  })

  const handleSortingChange: OnChangeFn<SortingState> = (updater) => {
    setSorting(updater)
    if (pagination.pageIndex > 0) {
      onPaginationChange({ ...pagination, pageIndex: 0 })
    }
  }

  const { table } = useDataTable({
    data: data?.items ?? [],
    columns,
    columnFilters,
    globalFilter,
    pagination,
    sorting,
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

  const options = (values: string[] | undefined) =>
    (values ?? []).map((value) => ({ label: value, value }))

  return (
    <DataTablePage
      table={table}
      columns={columns}
      isLoading={isLoading}
      isFetching={isFetching}
      emptyTitle={t('LobeHub No Conversations Found')}
      emptyDescription={t('LobeHub Check the schema or adjust your filters.')}
      skeletonKeyPrefix='lobehub-conversations-skeleton'
      applyHeaderSize
      toolbarProps={{
        searchPlaceholder: t('LobeHub Search conversations, users, agents or groups...'),
        searchDebounceMs: 500,
        filters: [
          {
            columnId: 'type',
            title: t('LobeHub Conversation Type'),
            singleSelect: true,
            options: [
              { label: t('LobeHub Agent'), value: 'agent' },
              { label: t('LobeHub Group'), value: 'group' },
              { label: t('LobeHub Unknown Source'), value: 'unknown' },
            ],
          },
          {
            columnId: 'status',
            title: t('Status'),
            singleSelect: true,
            options: options(filters?.statuses),
          },
          {
            columnId: 'trigger',
            title: t('LobeHub Trigger'),
            singleSelect: true,
            options: options(filters?.triggers),
          },
          {
            columnId: 'model',
            title: t('Model'),
            singleSelect: true,
            options: options(filters?.models),
          },
          {
            columnId: 'provider',
            title: t('Provider'),
            singleSelect: true,
            options: options(filters?.providers),
          },
        ],
        additionalSearch: (
          <div className='flex items-center gap-2'>
            <label className='text-muted-foreground flex items-center gap-2 text-xs'>
              {t('LobeHub Updated From')}
              <Input
                aria-label={t('LobeHub Updated From')}
                className='h-8 w-36'
                onChange={(event) => updateDateFilter('updatedFrom', event.target.value)}
                type='date'
                value={updatedFrom}
              />
            </label>
            <label className='text-muted-foreground flex items-center gap-2 text-xs'>
              {t('LobeHub Updated To')}
              <Input
                aria-label={t('LobeHub Updated To')}
                className='h-8 w-36'
                onChange={(event) => updateDateFilter('updatedTo', event.target.value)}
                type='date'
                value={updatedTo}
              />
            </label>
          </div>
        ),
        hasAdditionalFilters: Boolean(updatedFrom || updatedTo),
      }}
    />
  )
}
