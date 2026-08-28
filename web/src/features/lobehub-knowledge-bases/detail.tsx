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
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { Link } from '@tanstack/react-router'
import axios from 'axios'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { SectionPageLayout } from '@/components/layout'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Field, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Textarea } from '@/components/ui/textarea'

import {
  getKnowledgeBase,
  getKnowledgeBaseChunks,
  getKnowledgeBaseDocument,
  getKnowledgeBaseDocuments,
  getKnowledgeBaseFiles,
  updateKnowledgeBase,
} from './api'
import type { JsonValue, KnowledgeBase, KnowledgeBaseFile } from './types'
import { formatBytes, vectorCoverage } from './utils'

const showJson = (value: JsonValue) =>
  value == null ? '-' : JSON.stringify(value, null, 2)
const showDate = (value: string) => {
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? '-' : date.toLocaleString()
}

function Stat({ label, value }: { label: string; value: React.ReactNode }) {
  return (
    <Card size='sm'>
      <CardHeader>
        <CardTitle className='text-muted-foreground text-xs'>{label}</CardTitle>
      </CardHeader>
      <CardContent className='text-lg font-semibold'>{value}</CardContent>
    </Card>
  )
}

function EditKnowledgeBase({
  item,
  onOpenChange,
}: {
  item: KnowledgeBase
  onOpenChange: (open: boolean) => void
}) {
  const { t } = useTranslation()
  const queryClient = useQueryClient()
  const [name, setName] = useState(item.name)
  const [description, setDescription] = useState(item.description ?? '')
  const [avatar, setAvatar] = useState(item.avatar ?? '')
  const [saving, setSaving] = useState(false)
  const save = async (event: React.FormEvent) => {
    event.preventDefault()
    setSaving(true)
    try {
      const response = await updateKnowledgeBase(item.id, {
        name: name.trim(),
        description: description.trim() || null,
        avatar: avatar.trim() || null,
        expected_updated_at: item.updated_at,
      })
      if (!response.success) throw new Error(response.message)
      await queryClient.invalidateQueries({
        queryKey: ['lobehub-knowledge-base', item.id],
      })
      await queryClient.invalidateQueries({
        queryKey: ['lobehub-knowledge-bases'],
      })
      toast.success(t('LobeHub Knowledge Base Updated'))
      onOpenChange(false)
    } catch (error) {
      if (axios.isAxiosError(error) && error.response?.status === 409) {
        await queryClient.invalidateQueries({
          queryKey: ['lobehub-knowledge-base', item.id],
        })
        toast.error(t('LobeHub Knowledge Base Changed Refresh and Try Again'))
      } else toast.error(t('LobeHub Failed to Update Knowledge Base'))
    } finally {
      setSaving(false)
    }
  }
  return (
    <Dialog open onOpenChange={onOpenChange}>
      <DialogContent className='sm:max-w-lg'>
        <DialogHeader>
          <DialogTitle>{t('LobeHub Edit Knowledge Base')}</DialogTitle>
          <DialogDescription>
            {t('LobeHub Only Name Description and Avatar Can Be Changed')}
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={save} className='space-y-4'>
          <FieldGroup>
            <Field>
              <FieldLabel htmlFor='kb-name'>{t('LobeHub Name')}</FieldLabel>
              <Input
                id='kb-name'
                required
                value={name}
                onChange={(event) => setName(event.target.value)}
              />
            </Field>
            <Field>
              <FieldLabel htmlFor='kb-description'>
                {t('LobeHub Description')}
              </FieldLabel>
              <Textarea
                id='kb-description'
                value={description}
                onChange={(event) => setDescription(event.target.value)}
              />
            </Field>
            <Field>
              <FieldLabel htmlFor='kb-avatar'>{t('LobeHub Avatar')}</FieldLabel>
              <Input
                id='kb-avatar'
                value={avatar}
                onChange={(event) => setAvatar(event.target.value)}
              />
            </Field>
          </FieldGroup>
          <DialogFooter>
            <Button
              type='button'
              variant='outline'
              onClick={() => onOpenChange(false)}
            >
              {t('LobeHub Cancel')}
            </Button>
            <Button type='submit' disabled={saving || !name.trim()}>
              {saving ? t('LobeHub Saving') : t('LobeHub Save Changes')}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

function ContentTab({ id }: { id: string }) {
  const { t } = useTranslation()
  const [tab, setTab] = useState('files')
  const [filePage, setFilePage] = useState(1)
  const [documentPage, setDocumentPage] = useState(1)
  const [selectedDocument, setSelectedDocument] = useState<string | null>(null)
  const files = useQuery({
    queryKey: ['lobehub-knowledge-base-files', id, filePage],
    queryFn: () => getKnowledgeBaseFiles(id, filePage),
    enabled: tab === 'files',
  })
  const documents = useQuery({
    queryKey: ['lobehub-knowledge-base-documents', id, documentPage],
    queryFn: () => getKnowledgeBaseDocuments(id, documentPage),
    enabled: tab === 'documents',
  })
  const document = useQuery({
    queryKey: ['lobehub-knowledge-base-document', id, selectedDocument],
    queryFn: () => {
      if (!selectedDocument) throw new Error('Document is required')
      return getKnowledgeBaseDocument(id, selectedDocument)
    },
    enabled: Boolean(selectedDocument),
  })
  return (
    <>
      <Tabs value={tab} onValueChange={setTab}>
        <TabsList>
          <TabsTrigger value='files'>{t('LobeHub Files')}</TabsTrigger>
          <TabsTrigger value='documents'>{t('LobeHub Documents')}</TabsTrigger>
        </TabsList>
        <TabsContent value='files'>
          <div className='space-y-2'>
            {files.isLoading && <p>{t('LobeHub Loading Files')}</p>}
            {(files.data?.data?.items ?? []).map((file) => (
              <Card key={file.id} size='sm'>
                <CardHeader>
                  <CardTitle>{file.name}</CardTitle>
                </CardHeader>
                <CardContent className='grid gap-2 sm:grid-cols-4'>
                  <span>{file.file_type}</span>
                  <span>{formatBytes(file.size)}</span>
                  <span>
                    {file.document_count} / {file.chunk_count}
                  </span>
                  <span className='truncate'>{file.url || '-'}</span>
                </CardContent>
              </Card>
            ))}
            {!files.isLoading && !files.data?.data?.items.length && (
              <p>{t('LobeHub No Files Found')}</p>
            )}
            <Pager
              page={filePage}
              total={files.data?.data?.total ?? 0}
              setPage={setFilePage}
            />
          </div>
        </TabsContent>
        <TabsContent value='documents'>
          <div className='space-y-2'>
            {documents.isLoading && <p>{t('LobeHub Loading Documents')}</p>}
            {(documents.data?.data?.items ?? []).map((item) => (
              <button
                type='button'
                key={item.id}
                className='hover:bg-muted w-full rounded-lg border p-3 text-left'
                onClick={() => setSelectedDocument(item.id)}
              >
                <div className='font-medium'>
                  {item.title || item.filename || item.id}
                </div>
                <div className='text-muted-foreground text-xs'>
                  {item.file_type} · {item.total_char_count} ·{' '}
                  {item.total_line_count}
                </div>
              </button>
            ))}
            {!documents.isLoading && !documents.data?.data?.items.length && (
              <p>{t('LobeHub No Documents Found')}</p>
            )}
            <Pager
              page={documentPage}
              total={documents.data?.data?.total ?? 0}
              setPage={setDocumentPage}
            />
          </div>
        </TabsContent>
      </Tabs>
      <Dialog
        open={Boolean(selectedDocument)}
        onOpenChange={(value) => !value && setSelectedDocument(null)}
      >
        <DialogContent className='max-h-[85vh] overflow-auto sm:max-w-4xl'>
          <DialogHeader>
            <DialogTitle>
              {document.data?.data?.title ||
                document.data?.data?.filename ||
                t('LobeHub Document Detail')}
            </DialogTitle>
            <DialogDescription>{document.data?.data?.id}</DialogDescription>
          </DialogHeader>
          {document.isLoading ? (
            <p>{t('LobeHub Loading Document')}</p>
          ) : (
            <div className='space-y-4'>
              <section>
                <h3 className='font-medium'>{t('LobeHub Parsed Content')}</h3>
                <pre className='bg-muted max-h-96 overflow-auto rounded-lg p-3 whitespace-pre-wrap'>
                  {document.data?.data?.content || '-'}
                </pre>
              </section>
              <JsonBlock
                title={t('LobeHub Pages')}
                value={document.data?.data?.pages}
              />
              <JsonBlock
                title={t('LobeHub Editor Data')}
                value={document.data?.data?.editor_data}
              />
              <JsonBlock
                title={t('LobeHub Metadata')}
                value={document.data?.data?.metadata}
              />
            </div>
          )}
        </DialogContent>
      </Dialog>
    </>
  )
}

function JsonBlock({ title, value }: { title: string; value: JsonValue }) {
  return (
    <section>
      <h3 className='font-medium'>{title}</h3>
      <pre className='bg-muted max-h-64 overflow-auto rounded-lg p-3 text-xs whitespace-pre-wrap'>
        {showJson(value)}
      </pre>
    </section>
  )
}
function Pager({
  page,
  total,
  setPage,
}: {
  page: number
  total: number
  setPage: (page: number) => void
}) {
  const { t } = useTranslation()
  return (
    <div className='flex items-center justify-end gap-2'>
      <Button
        variant='outline'
        disabled={page <= 1}
        onClick={() => setPage(page - 1)}
      >
        {t('LobeHub Previous')}
      </Button>
      <span>{page}</span>
      <Button
        variant='outline'
        disabled={page * 20 >= total}
        onClick={() => setPage(page + 1)}
      >
        {t('LobeHub Next')}
      </Button>
    </div>
  )
}

function RAGTab({ id }: { id: string }) {
  const { t } = useTranslation()
  const [page, setPage] = useState(1)
  const [selected, setSelected] = useState<KnowledgeBaseFile | null>(null)
  const [chunkPage, setChunkPage] = useState(1)
  const files = useQuery({
    queryKey: ['lobehub-knowledge-base-files', id, page],
    queryFn: () => getKnowledgeBaseFiles(id, page),
  })
  const chunks = useQuery({
    queryKey: ['lobehub-knowledge-base-chunks', id, selected?.id, chunkPage],
    queryFn: () => {
      if (!selected) throw new Error('File is required')
      return getKnowledgeBaseChunks(id, selected.id, chunkPage)
    },
    enabled: Boolean(selected),
  })
  return (
    <div className='space-y-3'>
      {files.isLoading && <p>{t('LobeHub Loading RAG Status')}</p>}
      {(files.data?.data?.items ?? []).map((file) => (
        <Card key={file.id} size='sm'>
          <CardHeader>
            <CardTitle className='flex items-center justify-between'>
              <span>{file.name}</span>
              <Badge variant='outline'>{file.rag_status}</Badge>
            </CardTitle>
          </CardHeader>
          <CardContent className='grid gap-3 sm:grid-cols-2'>
            <div>
              <div>
                {t('LobeHub Chunk Task')}: {file.chunking_status ?? '-'}
              </div>
              <div>
                {t('LobeHub Embedding Task')}: {file.embedding_status ?? '-'}
              </div>
              <div>
                {t('LobeHub Vector Coverage')}:{' '}
                {file.chunk_count
                  ? Math.round(
                      (file.embedded_chunk_count / file.chunk_count) * 100
                    )
                  : 0}
                %
              </div>
            </div>
            <div>
              <JsonBlock
                title={t('LobeHub Task Errors')}
                value={{
                  chunking: file.chunking_error,
                  embedding: file.embedding_error,
                }}
              />
            </div>
            <Button
              variant='outline'
              onClick={() => {
                setSelected(file)
                setChunkPage(1)
              }}
            >
              {t('LobeHub View Chunks')}
            </Button>
          </CardContent>
        </Card>
      ))}
      <Pager
        page={page}
        total={files.data?.data?.total ?? 0}
        setPage={setPage}
      />
      <Dialog
        open={Boolean(selected)}
        onOpenChange={(value) => !value && setSelected(null)}
      >
        <DialogContent className='max-h-[85vh] overflow-auto sm:max-w-4xl'>
          <DialogHeader>
            <DialogTitle>
              {t('LobeHub Chunks')} · {selected?.name}
            </DialogTitle>
            <DialogDescription>
              {t('LobeHub Vector Values Are Never Returned')}
            </DialogDescription>
          </DialogHeader>
          <div className='space-y-3'>
            {chunks.isLoading && <p>{t('LobeHub Loading Chunks')}</p>}
            {(chunks.data?.data?.items ?? []).map((chunk) => (
              <Card key={chunk.id} size='sm'>
                <CardHeader>
                  <CardTitle>
                    #{chunk.index ?? '-'} · {chunk.model ?? '-'}
                  </CardTitle>
                </CardHeader>
                <CardContent className='space-y-2'>
                  <Badge variant='outline'>{String(chunk.has_embedding)}</Badge>
                  <p className='whitespace-pre-wrap'>{chunk.text || '-'}</p>
                  {chunk.abstract && (
                    <p className='text-muted-foreground'>{chunk.abstract}</p>
                  )}
                  <JsonBlock
                    title={t('LobeHub Metadata')}
                    value={chunk.metadata}
                  />
                </CardContent>
              </Card>
            ))}
            <Pager
              page={chunkPage}
              total={chunks.data?.data?.total ?? 0}
              setPage={setChunkPage}
            />
          </div>
        </DialogContent>
      </Dialog>
    </div>
  )
}

export function LobeHubKnowledgeBaseDetail({ id }: { id: string }) {
  const { t } = useTranslation()
  const [editing, setEditing] = useState(false)
  const detail = useQuery({
    queryKey: ['lobehub-knowledge-base', id],
    queryFn: () => getKnowledgeBase(id),
  })
  const item = detail.data?.data
  if (detail.isLoading) {
    return (
      <SectionPageLayout>
        <SectionPageLayout.Title>
          {t('LobeHub Loading Knowledge Base')}
        </SectionPageLayout.Title>
      </SectionPageLayout>
    )
  }
  if (!item) {
    return (
      <SectionPageLayout>
        <SectionPageLayout.Title>
          {t('LobeHub Knowledge Base Unavailable')}
        </SectionPageLayout.Title>
      </SectionPageLayout>
    )
  }
  return (
    <>
      <SectionPageLayout>
        <SectionPageLayout.Title>
          <div className='flex items-center gap-3'>
            <Link
              to='/lobehub/knowledge-bases'
              className='text-muted-foreground text-sm'
            >
              {t('LobeHub Back to Knowledge Bases')}
            </Link>
            <span>{item.name}</span>
            <Badge variant='outline'>{item.rag_status}</Badge>
          </div>
        </SectionPageLayout.Title>
        <SectionPageLayout.Content>
          <Tabs defaultValue='overview'>
            <TabsList>
              <TabsTrigger value='overview'>
                {t('LobeHub Overview')}
              </TabsTrigger>
              <TabsTrigger value='content'>{t('LobeHub Content')}</TabsTrigger>
              <TabsTrigger value='rag'>{t('LobeHub RAG')}</TabsTrigger>
            </TabsList>
            <TabsContent value='overview' className='space-y-4'>
              <div className='grid gap-3 sm:grid-cols-5'>
                <Stat
                  label={t('LobeHub Files')}
                  value={item.stats.file_count}
                />
                <Stat
                  label={t('LobeHub Documents')}
                  value={item.stats.document_count}
                />
                <Stat
                  label={t('LobeHub Chunks')}
                  value={item.stats.chunk_count}
                />
                <Stat
                  label={t('LobeHub Vector Coverage')}
                  value={vectorCoverage(item)}
                />
                <Stat
                  label={t('LobeHub Storage')}
                  value={formatBytes(item.stats.total_size)}
                />
              </div>
              <Card>
                <CardHeader>
                  <CardTitle className='flex items-center justify-between'>
                    <span>{t('LobeHub Metadata')}</span>
                    <Button onClick={() => setEditing(true)}>
                      {t('LobeHub Edit Safe Metadata')}
                    </Button>
                  </CardTitle>
                </CardHeader>
                <CardContent className='grid gap-3 sm:grid-cols-2'>
                  <div>
                    <b>{t('LobeHub ID')}:</b> {item.id}
                  </div>
                  <div>
                    <b>{t('LobeHub Type')}:</b> {item.type ?? '-'}
                  </div>
                  <div>
                    <b>{t('LobeHub Visibility')}:</b> {item.visibility}
                  </div>
                  <div>
                    <b>{t('LobeHub Public')}:</b> {String(item.is_public)}
                  </div>
                  <div>
                    <b>{t('LobeHub Owner')}:</b>{' '}
                    {item.owner.full_name ||
                      item.owner.username ||
                      item.owner.email ||
                      item.owner.id}
                  </div>
                  <div>
                    <b>{t('LobeHub Workspace')}:</b>{' '}
                    {item.workspace?.name ?? ''}
                  </div>
                  <div>
                    <b>{t('LobeHub Created At')}:</b>{' '}
                    {showDate(item.created_at)}
                  </div>
                  <div>
                    <b>{t('LobeHub Updated At')}:</b>{' '}
                    {showDate(item.updated_at)}
                  </div>
                  <div className='sm:col-span-2'>
                    <b>{t('LobeHub Description')}:</b> {item.description || '-'}
                  </div>
                  <div className='sm:col-span-2'>
                    <JsonBlock
                      title={t('LobeHub Settings')}
                      value={item.settings}
                    />
                  </div>
                </CardContent>
              </Card>
            </TabsContent>
            <TabsContent value='content'>
              <ContentTab id={id} />
            </TabsContent>
            <TabsContent value='rag'>
              <RAGTab id={id} />
            </TabsContent>
          </Tabs>
        </SectionPageLayout.Content>
      </SectionPageLayout>
      {editing && <EditKnowledgeBase item={item} onOpenChange={setEditing} />}
    </>
  )
}
