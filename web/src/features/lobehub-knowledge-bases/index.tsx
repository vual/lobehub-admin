/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

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
import { useState } from 'react'
import { useTranslation } from 'react-i18next'

import { SectionPageLayout } from '@/components/layout'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { NativeSelect, NativeSelectOption } from '@/components/ui/native-select'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

import { getKnowledgeBaseFilters, getKnowledgeBases } from './api'
import { formatBytes, vectorCoverage } from './utils'

const route = getRouteApi('/_authenticated/lobehub/knowledge-bases/')

const date = (value: string) => {
  const parsed = new Date(value)
  return Number.isNaN(parsed.getTime()) ? '-' : parsed.toLocaleString()
}

export function LobeHubKnowledgeBases() {
  const { t } = useTranslation()
  const search = route.useSearch()
  const navigate = route.useNavigate()
  const [query, setQuery] = useState(search.filter ?? '')
  const update = (values: Record<string, unknown>) =>
    navigate({
      search: (current) => ({ ...current, page: 1, ...values }),
      replace: true,
    })
  const list = useQuery({
    queryKey: ['lobehub-knowledge-bases', search],
    queryFn: () =>
      getKnowledgeBases({
        page: search.page ?? 1,
        page_size: search.pageSize ?? 20,
        q: search.filter || undefined,
        scope: search.scope || undefined,
        workspace_id: search.workspace || undefined,
        visibility: search.visibility || undefined,
        rag_status: search.rag || undefined,
        sort_by: search.sortBy || undefined,
        sort_order: search.sortOrder || undefined,
      }),
    placeholderData: (previous) => previous,
  })
  const filters = useQuery({
    queryKey: ['lobehub-knowledge-base-filters'],
    queryFn: getKnowledgeBaseFilters,
  })
  const data = list.data?.data

  return (
    <SectionPageLayout fixedContent>
      <SectionPageLayout.Title>
        {t('LobeHub Knowledge Bases')}
      </SectionPageLayout.Title>
      <SectionPageLayout.Content>
        <div className='flex h-full min-h-0 flex-col gap-4'>
          <form
            className='flex flex-wrap gap-2'
            onSubmit={(event) => {
              event.preventDefault()
              update({ filter: query.trim() })
            }}
          >
            <Input
              className='min-w-64 flex-1'
              value={query}
              onChange={(event) => setQuery(event.target.value)}
              placeholder={t(
                'LobeHub Search knowledge bases, owners or workspaces...'
              )}
            />
            <Button type='submit'>{t('LobeHub Search')}</Button>
            <NativeSelect
              value={search.scope ?? ''}
              onChange={(event) => update({ scope: event.target.value })}
              aria-label={t('LobeHub Scope')}
            >
              <NativeSelectOption value=''>
                {t('LobeHub All Scopes')}
              </NativeSelectOption>
              <NativeSelectOption value='personal'>
                {t('LobeHub Personal')}
              </NativeSelectOption>
              <NativeSelectOption value='workspace'>
                {t('LobeHub Workspace')}
              </NativeSelectOption>
            </NativeSelect>
            <NativeSelect
              value={search.workspace ?? ''}
              onChange={(event) => update({ workspace: event.target.value })}
              aria-label={t('LobeHub Workspace')}
            >
              <NativeSelectOption value=''>
                {t('LobeHub All Workspaces')}
              </NativeSelectOption>
              {(filters.data?.data?.workspaces ?? []).map((workspace) => (
                <NativeSelectOption key={workspace.id} value={workspace.id}>
                  {workspace.name} ({workspace.knowledge_base_count})
                </NativeSelectOption>
              ))}
            </NativeSelect>
            <NativeSelect
              value={search.visibility ?? ''}
              onChange={(event) => update({ visibility: event.target.value })}
              aria-label={t('LobeHub Visibility')}
            >
              <NativeSelectOption value=''>
                {t('LobeHub All Visibility')}
              </NativeSelectOption>
              <NativeSelectOption value='private'>private</NativeSelectOption>
              <NativeSelectOption value='public'>public</NativeSelectOption>
            </NativeSelect>
            <NativeSelect
              value={search.rag ?? ''}
              onChange={(event) => update({ rag: event.target.value })}
              aria-label={t('LobeHub RAG Status')}
            >
              <NativeSelectOption value=''>
                {t('LobeHub All RAG Statuses')}
              </NativeSelectOption>
              {['error', 'processing', 'empty', 'ready', 'unindexed'].map(
                (value) => (
                  <NativeSelectOption key={value} value={value}>
                    {value}
                  </NativeSelectOption>
                )
              )}
            </NativeSelect>
            <NativeSelect
              value={`${search.sortBy ?? 'updated_at'}:${search.sortOrder ?? 'desc'}`}
              onChange={(event) => {
                const [sortBy, sortOrder] = event.target.value.split(':')
                update({ sortBy, sortOrder })
              }}
              aria-label={t('LobeHub Sort')}
            >
              <NativeSelectOption value='updated_at:desc'>
                {t('LobeHub Recently Updated')}
              </NativeSelectOption>
              <NativeSelectOption value='created_at:desc'>
                {t('LobeHub Recently Created')}
              </NativeSelectOption>
              <NativeSelectOption value='file_count:desc'>
                {t('LobeHub Most Files')}
              </NativeSelectOption>
              <NativeSelectOption value='total_size:desc'>
                {t('LobeHub Largest Storage')}
              </NativeSelectOption>
            </NativeSelect>
          </form>

          <div className='min-h-0 flex-1 overflow-auto rounded-xl border'>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>{t('LobeHub Knowledge Base')}</TableHead>
                  <TableHead>{t('LobeHub Owner')}</TableHead>
                  <TableHead>{t('LobeHub Workspace')}</TableHead>
                  <TableHead>{t('LobeHub Content Statistics')}</TableHead>
                  <TableHead>{t('LobeHub Vector Coverage')}</TableHead>
                  <TableHead>{t('LobeHub Storage')}</TableHead>
                  <TableHead>{t('LobeHub RAG Status')}</TableHead>
                  <TableHead>{t('LobeHub Updated At')}</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {list.isLoading && (
                  <TableRow>
                    <TableCell colSpan={8}>
                      {t('LobeHub Loading Knowledge Bases')}
                    </TableCell>
                  </TableRow>
                )}
                {list.isError && (
                  <TableRow>
                    <TableCell colSpan={8} className='text-destructive'>
                      {t('LobeHub Failed to Load Knowledge Bases')}
                    </TableCell>
                  </TableRow>
                )}
                {!list.isLoading && !list.isError && !data?.items.length && (
                  <TableRow>
                    <TableCell colSpan={8}>
                      {t('LobeHub No Knowledge Bases Found')}
                    </TableCell>
                  </TableRow>
                )}
                {(data?.items ?? []).map((item) => (
                  <TableRow
                    key={item.id}
                    className='cursor-pointer'
                    onClick={() =>
                      navigate({
                        to: '/lobehub/knowledge-bases/$id',
                        params: { id: item.id },
                      })
                    }
                  >
                    <TableCell>
                      <div className='font-medium'>{item.name}</div>
                      <div className='text-muted-foreground max-w-52 truncate text-xs'>
                        {item.id}
                      </div>
                      <div className='text-muted-foreground text-xs'>
                        {item.type ?? '-'}
                      </div>
                    </TableCell>
                    <TableCell>
                      <div>
                        {item.owner.full_name || item.owner.username || '-'}
                      </div>
                      <div className='text-muted-foreground text-xs'>
                        {item.owner.email || item.owner.id}
                      </div>
                    </TableCell>
                    <TableCell>{item.workspace?.name ?? ''}</TableCell>
                    <TableCell>
                      {item.stats.file_count} / {item.stats.document_count} /{' '}
                      {item.stats.chunk_count}
                    </TableCell>
                    <TableCell>
                      {vectorCoverage(item)} ({item.stats.embedded_chunk_count}/
                      {item.stats.chunk_count})
                    </TableCell>
                    <TableCell>{formatBytes(item.stats.total_size)}</TableCell>
                    <TableCell>
                      <Badge variant='outline'>{item.rag_status}</Badge>
                    </TableCell>
                    <TableCell className='whitespace-nowrap'>
                      {date(item.updated_at)}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
          <div className='flex items-center justify-between text-sm'>
            <span>
              {t('LobeHub Total Knowledge Bases')}: {data?.total ?? 0}
            </span>
            <div className='flex gap-2'>
              <Button
                variant='outline'
                disabled={(search.page ?? 1) <= 1}
                onClick={() =>
                  navigate({
                    search: (current) => ({
                      ...current,
                      page: Math.max(1, (current.page ?? 1) - 1),
                    }),
                  })
                }
              >
                {t('LobeHub Previous')}
              </Button>
              <span className='self-center'>{search.page ?? 1}</span>
              <Button
                variant='outline'
                disabled={
                  (search.page ?? 1) * (search.pageSize ?? 20) >=
                  (data?.total ?? 0)
                }
                onClick={() =>
                  navigate({
                    search: (current) => ({
                      ...current,
                      page: (current.page ?? 1) + 1,
                    }),
                  })
                }
              >
                {t('LobeHub Next')}
              </Button>
            </div>
          </div>
        </div>
      </SectionPageLayout.Content>
    </SectionPageLayout>
  )
}
