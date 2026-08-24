import { lazy, Suspense, useState } from 'react'
import { useTranslation } from 'react-i18next'

import { SectionPageLayout } from '@/components/layout'

import { LobeHubConversationsTable } from './components/conversations-table'
import type { LobeHubConversation } from './types'

const LobeHubConversationDetailSheet = lazy(() =>
  import('./components/conversation-detail-sheet').then((module) => ({
    default: module.LobeHubConversationDetailSheet,
  }))
)

export function LobeHubConversations() {
  const { t } = useTranslation()
  const [selectedConversation, setSelectedConversation] =
    useState<LobeHubConversation | null>(null)

  return (
    <>
      <SectionPageLayout fixedContent>
        <SectionPageLayout.Title>{t('LobeHub Conversations')}</SectionPageLayout.Title>
        <SectionPageLayout.Content>
          <LobeHubConversationsTable onView={setSelectedConversation} />
        </SectionPageLayout.Content>
      </SectionPageLayout>

      {selectedConversation && (
        <Suspense fallback={null}>
          <LobeHubConversationDetailSheet
            conversation={selectedConversation}
            onOpenChange={(open) => !open && setSelectedConversation(null)}
            open={Boolean(selectedConversation)}
          />
        </Suspense>
      )}
    </>
  )
}
