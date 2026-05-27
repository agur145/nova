import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import App from './App'
import { TooltipProvider } from './components/ui/tooltip'

describe('App', () => {
  beforeEach(() => {
    globalThis.fetch = vi.fn(async (input) => {
      const rawUrl = typeof input === 'string' ? input : input.url
      const path = new URL(rawUrl, 'http://localhost').pathname
      const payloads: Record<string, unknown> = {
        '/api/workspace/current': { workspace: '/books/demo', has_state: true },
        '/api/workspace/tree': [],
        '/api/workspace/summary': { title: '', author: '', chapter_count: 0, total_words: 0, chapters: [] },
        '/api/styles': { styles: [] },
        '/api/books': { books: [] },
        '/api/settings': {
          user: {},
          workspace: {},
          effective: { max_open_tabs: 5 },
          paths: { nova_dir: '', user_config: '', workspace_config: '' },
        },
        '/api/system/select-directory': { path: '/books/from-picker', cancelled: false },
        '/api/workspace/switch': { workspace: '/books/from-picker', message: '已切换到: /books/from-picker' },
        '/api/sessions': { sessions: [] },
        '/api/session/messages': [],
        '/api/chat/active': { active: false },
      }

      return new Response(JSON.stringify(payloads[path] ?? {}), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      })
    }) as typeof fetch
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('renders the mode switch in the main header', async () => {
    render(
      <TooltipProvider>
        <App />
      </TooltipProvider>,
    )

    await waitFor(() => expect(globalThis.fetch).toHaveBeenCalledWith('/api/chat/active', undefined))
    const header = screen.getByText('Nova').closest('header')
    expect(header).not.toBeNull()
    expect(within(header as HTMLElement).getByRole('button', { name: 'IDE' })).toBeInTheDocument()
    expect(within(header as HTMLElement).getByRole('button', { name: 'Interactive' })).toBeInTheDocument()
  })

  it('does not render the removed task panel UI', async () => {
    render(
      <TooltipProvider>
        <App />
      </TooltipProvider>,
    )

    await waitFor(() => expect(globalThis.fetch).toHaveBeenCalledWith('/api/chat/active', undefined))
    expect(screen.queryByLabelText('显示/隐藏任务面板')).not.toBeInTheDocument()
    expect(screen.queryByText('任务')).not.toBeInTheDocument()
    expect(screen.queryByText('写作流')).not.toBeInTheDocument()
  })

  it('opens book management as a global dialog outside editor tabs', async () => {
    const user = userEvent.setup()
    render(
      <TooltipProvider>
        <App />
      </TooltipProvider>,
    )

    await waitFor(() => expect(globalThis.fetch).toHaveBeenCalledWith('/api/chat/active', undefined))
    await user.click(screen.getByRole('button', { name: '书籍管理' }))

    const dialog = await screen.findByRole('dialog')
    expect(within(dialog).getByText('最近书籍')).toBeInTheDocument()
    expect(within(dialog).getByText('打开其他目录')).toBeInTheDocument()
    expect(within(dialog).getByRole('button', { name: '关闭书籍管理' })).toBeInTheDocument()
    expect(within(dialog).queryByPlaceholderText('输入工作区目录路径...')).not.toBeInTheDocument()
  })

  it('opens another book directory through the system folder picker', async () => {
    const user = userEvent.setup()
    render(
      <TooltipProvider>
        <App />
      </TooltipProvider>,
    )

    await waitFor(() => expect(globalThis.fetch).toHaveBeenCalledWith('/api/chat/active', undefined))
    await user.click(screen.getByRole('button', { name: '书籍管理' }))
    const dialog = await screen.findByRole('dialog')
    await user.click(within(dialog).getByRole('button', { name: /选择文件夹/ }))

    await waitFor(() => {
      expect(globalThis.fetch).toHaveBeenCalledWith('/api/system/select-directory', undefined)
      expect(globalThis.fetch).toHaveBeenCalledWith('/api/workspace/switch', expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ path: '/books/from-picker' }),
      }))
    })
  })

  it('opens settings as a global dialog outside editor tabs', async () => {
    const user = userEvent.setup()
    render(
      <TooltipProvider>
        <App />
      </TooltipProvider>,
    )

    await waitFor(() => expect(globalThis.fetch).toHaveBeenCalledWith('/api/chat/active', undefined))
    await user.click(screen.getByRole('button', { name: '设置' }))

    const dialog = await screen.findByRole('dialog')
    expect(within(dialog).getByRole('button', { name: 'IDE 模式' })).toBeInTheDocument()
    expect(within(dialog).getByRole('button', { name: '互动模式' })).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: '关闭 设置' })).not.toBeInTheDocument()
  })
})
