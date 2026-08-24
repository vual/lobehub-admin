import { useInfiniteQuery, useQueryClient } from '@tanstack/react-query'
import { AlertTriangle, Bot, FileText, UserRound } from 'lucide-react'
import { useCallback, useEffect, useMemo, useRef } from 'react'
import { useTranslation } from 'react-i18next'

import { Markdown } from '@/components/ui/markdown'
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet'
import { cn } from '@/lib/utils'

import { getLobeHubConversationMessages } from '../api'
import type {
  LobeHubConversation,
  LobeHubConversationAttachment,
  LobeHubConversationMessage,
} from '../types'

type ConversationDetailSheetProps = {
  conversation: LobeHubConversation | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

function stringify(value: unknown): string {
  if (value == null) return ''
  if (typeof value === 'string') return value
  try {
    return JSON.stringify(value, null, 2)
  } catch {
    return String(value)
  }
}

function getObjectValue(value: unknown, key: string): unknown {
  if (!value || typeof value !== 'object' || Array.isArray(value)) return undefined
  return (value as Record<string, unknown>)[key]
}

function getReasoningContent(value: unknown): string {
  const content = getObjectValue(value, 'content')
  return typeof content === 'string' ? content : ''
}

function getArray(value: unknown): unknown[] {
  return Array.isArray(value) ? value : []
}

function isSafeExternalURL(value: string): boolean {
  try {
    const url = new URL(value)
    return url.protocol === 'http:' || url.protocol === 'https:'
  } catch {
    return false
  }
}

function formatDate(value: string): string {
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? '-' : date.toLocaleString()
}

function formatFileSize(size: number): string {
  if (!Number.isFinite(size) || size < 1024) return `${size} B`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`
  return `${(size / (1024 * 1024)).toFixed(1)} MB`
}

function AttachmentCard(props: { attachment: LobeHubConversationAttachment }) {
  const { t } = useTranslation()
  const attachment = props.attachment
  const isImage = attachment.file_type.startsWith('image/')
  const isAudio = attachment.file_type.startsWith('audio/')
  const isVideo = attachment.file_type.startsWith('video/')
  const canPreview = isSafeExternalURL(attachment.url)

  return (
    <div className='bg-background/70 space-y-2 rounded-md border p-2'>
      {canPreview && isImage && (
        <img
          alt={attachment.name}
          className='max-h-64 max-w-full rounded object-contain'
          loading='lazy'
          src={attachment.url}
        />
      )}
      {canPreview && isAudio && (
        <audio className='w-full' controls preload='metadata' src={attachment.url} />
      )}
      {canPreview && isVideo && (
        <video className='max-h-64 max-w-full rounded' controls preload='metadata' src={attachment.url} />
      )}
      <div className='flex items-center gap-2 text-xs'>
        <FileText aria-hidden='true' className='text-muted-foreground size-4' />
        <span className='min-w-0 flex-1 truncate'>{attachment.name}</span>
        <span className='text-muted-foreground whitespace-nowrap'>
          {formatFileSize(attachment.size)}
        </span>
      </div>
      {canPreview && (
        <a
          className='text-primary block truncate text-xs underline'
          href={attachment.url}
          rel='noopener noreferrer'
          target='_blank'
        >
          {t('LobeHub Open Attachment')}
        </a>
      )}
    </div>
  )
}

function JsonSection(props: { label: string; value: unknown }) {
  const text = stringify(props.value)
  if (!text) return null
  return (
    <details className='border-border/70 rounded-md border text-xs'>
      <summary className='hover:bg-muted/50 cursor-pointer px-3 py-2 font-medium'>
        {props.label}
      </summary>
      <pre className='bg-muted/40 max-h-72 overflow-auto border-t p-3 whitespace-pre-wrap break-words'>
        {text}
      </pre>
    </details>
  )
}

function ToolDetails(props: { message: LobeHubConversationMessage }) {
  const { t } = useTranslation()
  const tools = getArray(props.message.tools)
  const plugin = props.message.plugin
  if (tools.length === 0 && !plugin) return null

  return (
    <div className='space-y-2'>
      {tools.map((tool) => (
        <details className='border-border/70 rounded-md border text-xs' key={`${props.message.id}-tool-${stringify(tool)}`}>
          <summary className='hover:bg-muted/50 cursor-pointer px-3 py-2 font-medium'>
            {typeof tool === 'object' && tool !== null
              ? String((tool as Record<string, unknown>).apiName || (tool as Record<string, unknown>).identifier || t('LobeHub Tool Call'))
              : t('LobeHub Tool Call')}
          </summary>
          <pre className='bg-muted/40 max-h-72 overflow-auto border-t p-3 whitespace-pre-wrap break-words'>
            {stringify(tool)}
          </pre>
        </details>
      ))}
      {plugin && (
        <details className='border-border/70 rounded-md border text-xs'>
          <summary className='hover:bg-muted/50 cursor-pointer px-3 py-2 font-medium'>
            {plugin.api_name || plugin.identifier || t('LobeHub Plugin Details')}
          </summary>
          <div className='space-y-2 border-t p-3'>
            {plugin.arguments && <pre className='bg-muted/40 overflow-auto rounded p-2 whitespace-pre-wrap break-words'>{plugin.arguments}</pre>}
            <JsonSection label={t('LobeHub Plugin State')} value={plugin.state} />
            <JsonSection label={t('LobeHub Plugin Error')} value={plugin.error} />
            <JsonSection label={t('LobeHub Plugin Intervention')} value={plugin.intervention} />
          </div>
        </details>
      )}
    </div>
  )
}

function SearchDetails(props: { value: unknown }) {
  const { t } = useTranslation()
  const citations = getArray(getObjectValue(props.value, 'citations'))
  const queries = getArray(getObjectValue(props.value, 'searchQueries'))
  if (citations.length === 0 && queries.length === 0) return <JsonSection label={t('LobeHub Search Details')} value={props.value} />
  return (
    <details className='border-border/70 rounded-md border text-xs'>
      <summary className='hover:bg-muted/50 cursor-pointer px-3 py-2 font-medium'>
        {t('LobeHub Search Details')}
      </summary>
      <div className='space-y-2 border-t p-3'>
        {queries.length > 0 && <div className='flex flex-wrap gap-1'>{queries.map((query) => <span className='bg-muted rounded px-2 py-1' key={String(query)}>{String(query)}</span>)}</div>}
        {citations.map((citation) => {
          const url = String(getObjectValue(citation, 'url') || '')
          const title = String(getObjectValue(citation, 'title') || url || t('LobeHub Citation'))
          return isSafeExternalURL(url) ? <a className='text-primary block truncate underline' href={url} key={url} rel='noopener noreferrer' target='_blank'>{title}</a> : <span className='block truncate' key={title}>{title}</span>
        })}
      </div>
    </details>
  )
}

function MessageBlock(props: { message: LobeHubConversationMessage }) {
  const { t } = useTranslation()
  const message = props.message
  const isUser = message.role === 'user'
  const actorName = message.actor.name || message.actor.id || message.role
  const reasoning = getReasoningContent(message.reasoning)
  const content = message.content?.trim() || ''

  return (
    <div className={cn('flex w-full', isUser ? 'justify-end' : 'justify-start')}>
      <article className={cn('max-w-[92%] space-y-1', isUser ? 'items-end' : 'items-start')}>
        <div className={cn('flex items-center gap-1.5 text-xs', isUser ? 'justify-end' : 'justify-start')}>
          {isUser ? <UserRound aria-hidden='true' className='size-3' /> : <Bot aria-hidden='true' className='size-3' />}
          <span className='font-medium'>{actorName}</span>
          <span className='text-muted-foreground'>{message.role}</span>
          <span className='text-muted-foreground'>{formatDate(message.created_at)}</span>
        </div>
        <div className={cn('space-y-3 rounded-xl border px-3 py-2', isUser ? 'bg-primary text-primary-foreground border-primary' : 'bg-muted/40')}>
          {content ? <Markdown className={cn(isUser && '[&_a]:text-primary-foreground')} breaks>{content}</Markdown> : <p className='text-muted-foreground text-sm italic'>{t('LobeHub Empty Message')}</p>}
          {reasoning && <details className='border-current/20 rounded-md border text-xs'><summary className='cursor-pointer px-2 py-1'>{t('LobeHub Reasoning')}</summary><div className='border-current/20 border-t p-2'><Markdown breaks>{reasoning}</Markdown></div></details>}
          <ToolDetails message={message} />
          {message.search != null && <SearchDetails value={message.search} />}
          {message.translation?.content && <details className='border-current/20 rounded-md border text-xs'><summary className='cursor-pointer px-2 py-1'>{t('LobeHub Translation')}</summary><div className='border-current/20 border-t p-2'><Markdown breaks>{message.translation.content}</Markdown></div></details>}
          {message.error != null && <div className='border-destructive/40 text-destructive flex items-start gap-2 rounded-md border p-2 text-xs'><AlertTriangle aria-hidden='true' className='mt-0.5 size-4 shrink-0' /><span>{stringify(message.error)}</span></div>}
          <div className='space-y-2'>
            {message.attachments.map((attachment) => <AttachmentCard attachment={attachment} key={attachment.id} />)}
            <JsonSection label={t('LobeHub Usage')} value={message.usage} />
            <JsonSection label={t('LobeHub Message Metadata')} value={message.metadata} />
            <JsonSection label={t('LobeHub Thread Details')} value={message.thread} />
            <JsonSection label={t('LobeHub Message Group Details')} value={message.message_group} />
            <JsonSection label={t('LobeHub Search Queries')} value={message.queries} />
            <JsonSection label={t('LobeHub Audio Details')} value={message.tts} />
            <JsonSection label={t('LobeHub Editor Data')} value={message.editor_data} />
          </div>
        </div>
        {(message.model || message.provider || message.thread_id || message.message_group_id) && (
          <div className={cn('text-muted-foreground flex flex-wrap gap-1 text-[10px]', isUser ? 'justify-end' : 'justify-start')}>
            {message.model && <span>{message.model}</span>}
            {message.provider && <span>{message.provider}</span>}
            {message.thread_id && <span>thread: {message.thread_id}</span>}
            {message.message_group_id && <span>group: {message.message_group_id}</span>}
          </div>
        )}
      </article>
    </div>
  )
}

export function LobeHubConversationDetailSheet(props: ConversationDetailSheetProps) {
  const { t } = useTranslation()
  const queryClient = useQueryClient()
  const scrollRef = useRef<HTMLDivElement>(null)
  const topSentinelRef = useRef<HTMLDivElement>(null)
  const initializedRef = useRef(false)
  const conversationID = props.conversation?.id ?? ''

  const query = useInfiniteQuery({
    queryKey: ['lobehub-conversation-messages', conversationID],
    queryFn: async ({ pageParam }: { pageParam: string | undefined }) => {
      const result = await getLobeHubConversationMessages(conversationID, pageParam)
      if (!result.success || !result.data) {
        throw new Error(result.message || t('LobeHub Failed to load conversation'))
      }
      return result.data
    },
    enabled: props.open && Boolean(conversationID),
    initialPageParam: undefined,
    getNextPageParam: (lastPage) => lastPage.next_cursor || undefined,
    staleTime: 0,
    gcTime: 0,
  })

  const messages = useMemo(() => {
    const pages = query.data?.pages ?? []
    return [...pages].reverse().flatMap((page) => page.items)
  }, [query.data?.pages])

  const loadOlder = useCallback(async () => {
    if (!query.hasNextPage || query.isFetchingNextPage) return
    const element = scrollRef.current
    const previousHeight = element?.scrollHeight ?? 0
    const previousTop = element?.scrollTop ?? 0
    await query.fetchNextPage()
    requestAnimationFrame(() => {
      if (!element) return
      element.scrollTop = element.scrollHeight - previousHeight + previousTop
    })
  }, [query])

  useEffect(() => {
    if (!props.open) return
    initializedRef.current = false
  }, [conversationID, props.open])

  useEffect(() => {
    const element = scrollRef.current
    if (!element || initializedRef.current || !query.data?.pages.length) return
    element.scrollTop = element.scrollHeight
    initializedRef.current = true
  }, [query.data?.pages.length])

  useEffect(() => {
    const sentinel = topSentinelRef.current
    const root = scrollRef.current
    if (!sentinel || !root || !props.open) return
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0]?.isIntersecting) void loadOlder()
      },
      { root, threshold: 0.1 }
    )
    observer.observe(sentinel)
    return () => observer.disconnect()
  }, [loadOlder, props.open])

  const close = (open: boolean) => {
    props.onOpenChange(open)
    if (!open && conversationID) {
      queryClient.removeQueries({ queryKey: ['lobehub-conversation-messages', conversationID] })
    }
  }

  return (
    <Sheet onOpenChange={close} open={props.open}>
      <SheetContent className='w-full gap-0 p-0 sm:max-w-3xl' side='right'>
        <SheetHeader className='border-b pr-12'>
          <SheetTitle>{props.conversation?.title || t('LobeHub Conversation')}</SheetTitle>
          <SheetDescription>
            {props.conversation?.user.full_name || props.conversation?.user.username || props.conversation?.user.email || props.conversation?.user.id || '-'}
            {' · '}
            {props.conversation?.source.name || props.conversation?.source.id || t('LobeHub Unknown Source')}
          </SheetDescription>
        </SheetHeader>
        <div className='min-h-0 flex-1 overflow-y-auto p-4' ref={scrollRef}>
          <div className='h-1' ref={topSentinelRef} />
          {query.isLoading && <div className='text-muted-foreground py-12 text-center text-sm'>{t('LobeHub Loading Conversation')}</div>}
          {query.isError && <div className='text-destructive py-12 text-center text-sm'>{query.error.message}</div>}
          {!query.isLoading && !query.isError && messages.length === 0 && <div className='text-muted-foreground py-12 text-center text-sm'>{t('LobeHub No Messages')}</div>}
          <div className='space-y-5'>
            {messages.map((message) => <MessageBlock key={message.id} message={message} />)}
          </div>
          {query.isFetchingNextPage && <div className='text-muted-foreground py-3 text-center text-xs'>{t('LobeHub Loading Earlier Messages')}</div>}
          {query.isFetchNextPageError && <button className='text-destructive w-full py-3 text-center text-xs underline' onClick={() => void loadOlder()} type='button'>{t('LobeHub Retry Earlier Messages')}</button>}
        </div>
      </SheetContent>
    </Sheet>
  )
}
