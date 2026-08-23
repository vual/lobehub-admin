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

import { useState } from 'react'
import { useTranslation } from 'react-i18next'

import { SectionPageLayout } from '@/components/layout'

import { UserEditSheet } from './components/user-edit-sheet'
import { LobeHubUsersTable } from './components/users-table'
import type { LobeHubUser } from './types'

export function LobeHubUsers() {
  const { t } = useTranslation()
  const [editingUser, setEditingUser] = useState<LobeHubUser | null>(null)

  return (
    <>
      <SectionPageLayout fixedContent>
        <SectionPageLayout.Title>{t('LobeHub Users')}</SectionPageLayout.Title>
        <SectionPageLayout.Content>
          <LobeHubUsersTable onEdit={setEditingUser} />
        </SectionPageLayout.Content>
      </SectionPageLayout>

      <UserEditSheet
        user={editingUser}
        open={Boolean(editingUser)}
        onOpenChange={(open) => !open && setEditingUser(null)}
      />
    </>
  )
}
