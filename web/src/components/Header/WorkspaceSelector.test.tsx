import { http, HttpResponse } from 'msw'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import type { ComponentProps } from 'react'
import { WorkspaceSelector } from './WorkspaceSelector'
import { server } from '@/test/msw/server'
import type { BookRecord } from '@/lib/api'
import { TooltipProvider } from '@/components/ui/tooltip'

const books: BookRecord[] = [
  {
    path: '/books/novel-a',
    name: '长篇小说 A',
    author: '作者甲',
    last_opened_at: new Date().toISOString(),
  },
]

describe('WorkspaceSelector', () => {
  const renderSelector = (props: ComponentProps<typeof WorkspaceSelector>) => render(
    <TooltipProvider>
      <WorkspaceSelector {...props} />
    </TooltipProvider>,
  )

  it('点击书籍后打开 Popover 列表', async () => {
    const user = userEvent.setup()

    renderSelector(
      {
        workspace: '/books/novel-a',
        books,
        onSwitch: vi.fn(),
      },
    )

    await user.click(screen.getByRole('button', { name: '书籍' }))

    expect(screen.getByText('最近书籍')).toBeInTheDocument()
    expect(screen.getByText('长篇小说 A')).toBeInTheDocument()
  })

  it('点击书籍记录后调用切换逻辑', async () => {
    const user = userEvent.setup()
    const handleSwitch = vi.fn()

    server.use(
      http.post('/api/workspace/switch', () => HttpResponse.json({ workspace: '/books/novel-a' })),
    )

    renderSelector(
      {
        workspace: '/books/current',
        books,
        onSwitch: handleSwitch,
      },
    )

    await user.click(screen.getByRole('button', { name: '书籍' }))
    await user.click(screen.getByRole('button', { name: /长篇小说 A/ }))

    expect(handleSwitch).toHaveBeenCalledWith('/books/novel-a')
  })
})
