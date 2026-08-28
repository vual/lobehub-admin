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
import { describe, expect, test } from 'vitest'

import type { KnowledgeBase } from './types'
import { formatBytes, vectorCoverage } from './utils'

describe('LobeHub knowledge base presentation', () => {
  test('formats storage sizes', () => {
    expect(formatBytes(0)).toBe('0 B')
    expect(formatBytes(1536)).toBe('1.5 KB')
  })

  test('calculates vector coverage without division by zero', () => {
    const item = {
      stats: { chunk_count: 4, embedded_chunk_count: 3 },
    } as KnowledgeBase
    expect(vectorCoverage(item)).toBe('75%')
    item.stats.chunk_count = 0
    expect(vectorCoverage(item)).toBe('0%')
  })
})
