import { createFileRoute, redirect } from '@tanstack/react-router'
import { z } from 'zod'

import { LobeHubConversations } from '@/features/lobehub-conversations'
import { ROLE } from '@/lib/roles'
import { useAuthStore } from '@/stores/auth-store'

const searchSchema = z.object({
  page: z.number().optional().catch(1),
  pageSize: z.number().optional().catch(undefined),
  filter: z.string().optional().catch(''),
  type: z.array(z.enum(['agent', 'group', 'unknown'])).optional().catch([]),
  status: z.array(z.string()).optional().catch([]),
  trigger: z.array(z.string()).optional().catch([]),
  model: z.array(z.string()).optional().catch([]),
  provider: z.array(z.string()).optional().catch([]),
  updatedFrom: z.string().optional().catch(''),
  updatedTo: z.string().optional().catch(''),
})

export const Route = createFileRoute('/_authenticated/lobehub/conversations/')({
  beforeLoad: () => {
    const { auth } = useAuthStore.getState()
    if (!auth.user || auth.user.role < ROLE.SUPER_ADMIN) {
      throw redirect({ to: '/403' })
    }
  },
  validateSearch: searchSchema,
  component: LobeHubConversations,
})
