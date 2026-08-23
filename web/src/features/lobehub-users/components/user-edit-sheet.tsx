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

import { useQuery, useQueryClient } from '@tanstack/react-query'
import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import {
  sideDrawerContentClassName,
  sideDrawerFooterClassName,
  sideDrawerFormClassName,
  sideDrawerHeaderClassName,
} from '@/components/drawer-layout'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Field, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet'
import { Switch } from '@/components/ui/switch'

import { getLobeHubUser, updateLobeHubUser } from '../api'
import type { LobeHubUser } from '../types'

type UserEditSheetProps = {
  user: LobeHubUser | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

type EditFormState = {
  username: string
  email: string
  avatar: string
  phone: string
  firstName: string
  lastName: string
  fullName: string
  emailVerified: boolean
  phoneVerified: boolean
  updatedAt: string
}

const EMPTY_FORM: EditFormState = {
  username: '',
  email: '',
  avatar: '',
  phone: '',
  firstName: '',
  lastName: '',
  fullName: '',
  emailVerified: false,
  phoneVerified: false,
  updatedAt: '',
}

function userToForm(user: LobeHubUser): EditFormState {
  return {
    username: user.username ?? '',
    email: user.email ?? '',
    avatar: user.avatar ?? '',
    phone: user.phone ?? '',
    firstName: user.first_name ?? '',
    lastName: user.last_name ?? '',
    fullName: user.full_name ?? '',
    emailVerified: user.email_verified,
    phoneVerified: user.phone_number_verified,
    updatedAt: user.updated_at,
  }
}

export function UserEditSheet({
  user,
  open,
  onOpenChange,
}: UserEditSheetProps) {
  const { t } = useTranslation()
  const queryClient = useQueryClient()
  const [form, setForm] = useState<EditFormState>(EMPTY_FORM)
  const [submitting, setSubmitting] = useState(false)

  const detailQuery = useQuery({
    queryKey: ['lobehub-user', user?.id],
    queryFn: () => {
      if (!user) throw new Error('LobeHub user is required')
      return getLobeHubUser(user.id)
    },
    enabled: open && Boolean(user),
  })

  useEffect(() => {
    const loadedUser = detailQuery.data?.data?.user
    if (open && loadedUser) {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setForm(userToForm(loadedUser))
    }
  }, [detailQuery.data, open])

  const setText = (key: keyof EditFormState, value: string) => {
    setForm((current) => ({ ...current, [key]: value }))
  }

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault()
    if (!user || !form.updatedAt) return
    setSubmitting(true)
    try {
      const result = await updateLobeHubUser(user.id, {
        username: form.username.trim() || null,
        email: form.email.trim() || null,
        avatar: form.avatar.trim() || null,
        phone: form.phone.trim() || null,
        first_name: form.firstName.trim() || null,
        last_name: form.lastName.trim() || null,
        full_name: form.fullName.trim() || null,
        email_verified: form.emailVerified,
        phone_number_verified: form.phoneVerified,
        expected_updated_at: form.updatedAt,
      })
      if (!result.success) {
        toast.error(result.message || t('LobeHub Failed to update user'))
        return
      }
      toast.success(t('LobeHub User Updated'))
      await queryClient.invalidateQueries({ queryKey: ['lobehub-users'] })
      await queryClient.invalidateQueries({
        queryKey: ['lobehub-user', user.id],
      })
      onOpenChange(false)
    } catch {
      toast.error(t('LobeHub Failed to update user'))
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent className={sideDrawerContentClassName('sm:max-w-xl')}>
        <SheetHeader className={sideDrawerHeaderClassName()}>
          <SheetTitle>{t('LobeHub Edit User')}</SheetTitle>
          <SheetDescription>
            {t(
              'LobeHub Update profile and verification fields stored by the service.'
            )}
          </SheetDescription>
        </SheetHeader>

        <form onSubmit={handleSubmit} className={sideDrawerFormClassName()}>
          <Alert>
            <AlertTitle>{t('LobeHub Database-only operation')}</AlertTitle>
            <AlertDescription>
              {t(
                'LobeHub Changing an email or phone number automatically clears its verification state.'
              )}
            </AlertDescription>
          </Alert>

          <FieldGroup>
            <Field>
              <FieldLabel htmlFor='lobehub-username'>
                {t('Username')}
              </FieldLabel>
              <Input
                id='lobehub-username'
                value={form.username}
                onChange={(event) => setText('username', event.target.value)}
              />
            </Field>
            <Field>
              <FieldLabel htmlFor='lobehub-email'>{t('Email')}</FieldLabel>
              <Input
                id='lobehub-email'
                type='email'
                value={form.email}
                onChange={(event) => setText('email', event.target.value)}
              />
            </Field>
            <Field>
              <FieldLabel htmlFor='lobehub-full-name'>
                {t('LobeHub Full Name')}
              </FieldLabel>
              <Input
                id='lobehub-full-name'
                value={form.fullName}
                onChange={(event) => setText('fullName', event.target.value)}
              />
            </Field>
            <div className='grid gap-5 sm:grid-cols-2'>
              <Field>
                <FieldLabel htmlFor='lobehub-first-name'>
                  {t('LobeHub First Name')}
                </FieldLabel>
                <Input
                  id='lobehub-first-name'
                  value={form.firstName}
                  onChange={(event) => setText('firstName', event.target.value)}
                />
              </Field>
              <Field>
                <FieldLabel htmlFor='lobehub-last-name'>
                  {t('LobeHub Last Name')}
                </FieldLabel>
                <Input
                  id='lobehub-last-name'
                  value={form.lastName}
                  onChange={(event) => setText('lastName', event.target.value)}
                />
              </Field>
            </div>
            <Field>
              <FieldLabel htmlFor='lobehub-phone'>
                {t('LobeHub Phone')}
              </FieldLabel>
              <Input
                id='lobehub-phone'
                value={form.phone}
                onChange={(event) => setText('phone', event.target.value)}
              />
            </Field>
            <Field>
              <FieldLabel htmlFor='lobehub-avatar'>
                {t('LobeHub Avatar URL')}
              </FieldLabel>
              <Input
                id='lobehub-avatar'
                type='url'
                value={form.avatar}
                onChange={(event) => setText('avatar', event.target.value)}
              />
            </Field>
            <Field orientation='horizontal'>
              <FieldLabel htmlFor='lobehub-email-verified'>
                {t('LobeHub Email Verified')}
              </FieldLabel>
              <Switch
                id='lobehub-email-verified'
                checked={form.emailVerified}
                onCheckedChange={(checked) =>
                  setForm((current) => ({
                    ...current,
                    emailVerified: checked,
                  }))
                }
              />
            </Field>
            <Field orientation='horizontal'>
              <FieldLabel htmlFor='lobehub-phone-verified'>
                {t('LobeHub Phone Verified')}
              </FieldLabel>
              <Switch
                id='lobehub-phone-verified'
                checked={form.phoneVerified}
                onCheckedChange={(checked) =>
                  setForm((current) => ({
                    ...current,
                    phoneVerified: checked,
                  }))
                }
              />
            </Field>
          </FieldGroup>

          <SheetFooter className={sideDrawerFooterClassName()}>
            <SheetClose render={<Button type='button' variant='outline' />}>
              {t('Cancel')}
            </SheetClose>
            <Button
              type='submit'
              disabled={submitting || detailQuery.isLoading}
            >
              {submitting ? t('Saving...') : t('Save Changes')}
            </Button>
          </SheetFooter>
        </form>
      </SheetContent>
    </Sheet>
  )
}
