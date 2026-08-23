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

import { useQueryClient } from '@tanstack/react-query'
import {
  Ban,
  Copy,
  KeyRound,
  Pencil,
  ShieldCheck,
  UserRoundCheck,
} from 'lucide-react'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { ConfirmDialog } from '@/components/confirm-dialog'
import { DataTableRowActionMenu } from '@/components/data-table'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
} from '@/components/ui/dropdown-menu'
import { Field, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { ROLE } from '@/lib/roles'
import { useAuthStore } from '@/stores/auth-store'

import {
  banLobeHubUser,
  changeLobeHubUserRole,
  resetLobeHubUserPassword,
  unbanLobeHubUser,
} from '../api'
import type { LobeHubPasswordResetResult, LobeHubUser } from '../types'

type UserRowActionsProps = {
  user: LobeHubUser
  onEdit: (user: LobeHubUser) => void
}

export function UserRowActions({ user, onEdit }: UserRowActionsProps) {
  const { t } = useTranslation()
  const queryClient = useQueryClient()
  const currentRole = useAuthStore((state) => state.auth.user?.role)
  const [banOpen, setBanOpen] = useState(false)
  const [unbanOpen, setUnbanOpen] = useState(false)
  const [roleOpen, setRoleOpen] = useState(false)
  const [passwordOpen, setPasswordOpen] = useState(false)
  const [reason, setReason] = useState('')
  const [expiresAt, setExpiresAt] = useState('')
  const [passwordResult, setPasswordResult] =
    useState<LobeHubPasswordResetResult | null>(null)
  const [submitting, setSubmitting] = useState(false)

  const refresh = async () => {
    await queryClient.invalidateQueries({ queryKey: ['lobehub-users'] })
    await queryClient.invalidateQueries({
      queryKey: ['lobehub-user', user.id],
    })
  }

  const handleBan = async () => {
    if (!reason.trim()) {
      toast.error(t('LobeHub Ban reason is required'))
      return
    }
    setSubmitting(true)
    try {
      const result = await banLobeHubUser(user.id, {
        reason: reason.trim(),
        expires_at: expiresAt ? new Date(expiresAt).toISOString() : null,
      })
      if (!result.success) {
        toast.error(result.message || t('LobeHub Failed to ban user'))
        return
      }
      toast.success(t('LobeHub User Banned'))
      await refresh()
      setBanOpen(false)
      setReason('')
      setExpiresAt('')
    } catch {
      toast.error(t('LobeHub Failed to ban user'))
    } finally {
      setSubmitting(false)
    }
  }

  const handleUnban = async () => {
    setSubmitting(true)
    try {
      const result = await unbanLobeHubUser(user.id)
      if (!result.success) {
        toast.error(result.message || t('LobeHub Failed to unban user'))
        return
      }
      toast.success(t('LobeHub User Unbanned'))
      await refresh()
      setUnbanOpen(false)
    } catch {
      toast.error(t('LobeHub Failed to unban user'))
    } finally {
      setSubmitting(false)
    }
  }

  const targetRole: 'user' | 'admin' = user.role === 'admin' ? 'user' : 'admin'
  const hasCustomRole = user.role !== 'user' && user.role !== 'admin'
  const handleRoleChange = async () => {
    setSubmitting(true)
    try {
      const result = await changeLobeHubUserRole(
        user.id,
        targetRole,
        hasCustomRole
      )
      if (!result.success) {
        toast.error(result.message || t('LobeHub Failed to change role'))
        return
      }
      toast.success(t('LobeHub Role Updated'))
      await refresh()
      setRoleOpen(false)
    } catch {
      toast.error(t('LobeHub Failed to change role'))
    } finally {
      setSubmitting(false)
    }
  }

  const handlePasswordReset = async () => {
    setSubmitting(true)
    try {
      const result = await resetLobeHubUserPassword(user.id)
      if (!result.success || !result.data) {
        toast.error(result.message || t('LobeHub Failed to reset password'))
        return
      }
      setPasswordResult(result.data)
      toast.success(t('LobeHub Temporary password generated'))
      await refresh()
    } catch {
      toast.error(t('LobeHub Failed to reset password'))
    } finally {
      setSubmitting(false)
    }
  }

  const handlePasswordOpenChange = (open: boolean) => {
    setPasswordOpen(open)
    if (!open) setPasswordResult(null)
  }

  return (
    <>
      <DataTableRowActionMenu ariaLabel={t('Open menu')}>
        <DropdownMenuGroup>
          <DropdownMenuItem onClick={() => onEdit(user)}>
            {t('Edit')}
            <DropdownMenuShortcut>
              <Pencil size={16} />
            </DropdownMenuShortcut>
          </DropdownMenuItem>
          {user.banned ? (
            <DropdownMenuItem
              onSelect={(event) => {
                event.preventDefault()
                setUnbanOpen(true)
              }}
            >
              {t('LobeHub Unban Action')}
              <DropdownMenuShortcut>
                <UserRoundCheck size={16} />
              </DropdownMenuShortcut>
            </DropdownMenuItem>
          ) : (
            <DropdownMenuItem
              onSelect={(event) => {
                event.preventDefault()
                setBanOpen(true)
              }}
              className='text-destructive focus:text-destructive'
            >
              {t('LobeHub Ban Action')}
              <DropdownMenuShortcut>
                <Ban size={16} />
              </DropdownMenuShortcut>
            </DropdownMenuItem>
          )}
        </DropdownMenuGroup>
        <DropdownMenuSeparator />
        <DropdownMenuGroup>
          {currentRole === ROLE.SUPER_ADMIN && (
            <DropdownMenuItem
              onSelect={(event) => {
                event.preventDefault()
                setRoleOpen(true)
              }}
            >
              {targetRole === 'admin'
                ? t('LobeHub Grant Admin')
                : t('LobeHub Revoke Admin')}
              <DropdownMenuShortcut>
                <ShieldCheck size={16} />
              </DropdownMenuShortcut>
            </DropdownMenuItem>
          )}
          <DropdownMenuItem
            onSelect={(event) => {
              event.preventDefault()
              setPasswordOpen(true)
            }}
          >
            {t('LobeHub Reset User Password')}
            <DropdownMenuShortcut>
              <KeyRound size={16} />
            </DropdownMenuShortcut>
          </DropdownMenuItem>
        </DropdownMenuGroup>
      </DataTableRowActionMenu>

      <Dialog open={banOpen} onOpenChange={setBanOpen}>
        <DialogContent className='sm:max-w-md'>
          <DialogHeader>
            <DialogTitle>{t('LobeHub Ban User')}</DialogTitle>
            <DialogDescription>
              {t(
                'LobeHub Ban {{user}} and remove revocable database sessions.',
                {
                  user: user.username || user.email || user.id,
                }
              )}
            </DialogDescription>
          </DialogHeader>
          <FieldGroup>
            <Field>
              <FieldLabel htmlFor={`ban-reason-${user.id}`}>
                {t('LobeHub Ban Reason')}
              </FieldLabel>
              <Textarea
                id={`ban-reason-${user.id}`}
                value={reason}
                maxLength={500}
                onChange={(event) => setReason(event.target.value)}
              />
            </Field>
            <Field>
              <FieldLabel htmlFor={`ban-expires-${user.id}`}>
                {t('LobeHub Expires At (optional)')}
              </FieldLabel>
              <Input
                id={`ban-expires-${user.id}`}
                type='datetime-local'
                value={expiresAt}
                onChange={(event) => setExpiresAt(event.target.value)}
              />
            </Field>
          </FieldGroup>
          <Alert>
            <AlertTitle>{t('LobeHub Database-only limitation')}</AlertTitle>
            <AlertDescription>
              {t(
                'LobeHub Redis-cached sessions and stateless OIDC tokens may remain valid until they expire.'
              )}
            </AlertDescription>
          </Alert>
          <DialogFooter>
            <Button variant='outline' onClick={() => setBanOpen(false)}>
              {t('Cancel')}
            </Button>
            <Button
              variant='destructive'
              disabled={submitting}
              onClick={handleBan}
            >
              {submitting ? t('Processing...') : t('LobeHub Ban Action')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <ConfirmDialog
        open={unbanOpen}
        onOpenChange={setUnbanOpen}
        title={t('LobeHub Unban User')}
        desc={t('LobeHub Remove the ban from {{user}}?', {
          user: user.username || user.email || user.id,
        })}
        confirmText={t('LobeHub Unban Action')}
        handleConfirm={handleUnban}
        isLoading={submitting}
      />

      <ConfirmDialog
        open={roleOpen}
        onOpenChange={setRoleOpen}
        title={
          targetRole === 'admin'
            ? t('LobeHub Grant Admin')
            : t('LobeHub Revoke Admin')
        }
        desc={t(
          'LobeHub Change the global role for {{user}} to {{role}}? Database sessions will be removed.',
          {
            user: user.username || user.email || user.id,
            role: targetRole,
          }
        )}
        confirmText={t('LobeHub Change Role')}
        handleConfirm={handleRoleChange}
        isLoading={submitting}
      />

      <Dialog open={passwordOpen} onOpenChange={handlePasswordOpenChange}>
        <DialogContent className='sm:max-w-md'>
          <DialogHeader>
            <DialogTitle>{t('LobeHub Reset Password')}</DialogTitle>
            <DialogDescription>
              {passwordResult
                ? t(
                    'LobeHub This temporary password is shown only once. Copy it now.'
                  )
                : t(
                    'LobeHub Generate a strong temporary password and remove revocable database sessions.'
                  )}
            </DialogDescription>
          </DialogHeader>
          {passwordResult ? (
            <FieldGroup>
              <Field>
                <FieldLabel htmlFor={`temporary-password-${user.id}`}>
                  {t('LobeHub Temporary Password')}
                </FieldLabel>
                <div className='flex gap-2'>
                  <Input
                    id={`temporary-password-${user.id}`}
                    readOnly
                    value={passwordResult.temporary_password}
                  />
                  <Button
                    type='button'
                    variant='outline'
                    size='icon'
                    aria-label={t('Copy')}
                    onClick={async () => {
                      await navigator.clipboard.writeText(
                        passwordResult.temporary_password
                      )
                      toast.success(t('Copied'))
                    }}
                  >
                    <Copy />
                  </Button>
                </div>
              </Field>
              <Alert>
                <AlertTitle>{t('LobeHub Credentials removed')}</AlertTitle>
                <AlertDescription>
                  {t(
                    'LobeHub {{sessions}} sessions and {{oidc}} OIDC records removed.',
                    {
                      sessions: passwordResult.revoked.sessions,
                      oidc: passwordResult.revoked.oidc,
                    }
                  )}
                </AlertDescription>
              </Alert>
            </FieldGroup>
          ) : (
            <Alert>
              <AlertTitle>{t('LobeHub Database-only limitation')}</AlertTitle>
              <AlertDescription>
                {t(
                  'LobeHub Redis-cached sessions and stateless OIDC tokens may remain valid until they expire.'
                )}
              </AlertDescription>
            </Alert>
          )}
          <DialogFooter>
            <Button
              variant='outline'
              onClick={() => handlePasswordOpenChange(false)}
            >
              {passwordResult ? t('Close') : t('Cancel')}
            </Button>
            {!passwordResult && (
              <Button disabled={submitting} onClick={handlePasswordReset}>
                {submitting
                  ? t('Processing...')
                  : t('LobeHub Generate Password')}
              </Button>
            )}
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  )
}
