import { http, HttpResponse } from 'msw'

export const handlers = [
  http.get('/api/session/messages', () => HttpResponse.json([])),
  http.get('/api/sessions', () => HttpResponse.json({ sessions: [] })),
  http.post('/api/sessions', async ({ request }) => {
    const body = await request.json() as { title?: string }
    return HttpResponse.json({
      id: 'session-new',
      title: body.title || '新会话',
      created_at: '2026-05-17T00:00:00Z',
      updated_at: '2026-05-17T00:00:00Z',
      active: true,
      message_count: 0,
    })
  }),
  http.post('/api/sessions/switch', async ({ request }) => {
    const body = await request.json() as { id?: string }
    return HttpResponse.json({
      id: body.id || 'session-a',
      title: '目标会话',
      created_at: '2026-05-17T00:00:00Z',
      updated_at: '2026-05-17T00:00:00Z',
      active: true,
      message_count: 1,
    })
  }),
  http.post('/api/sessions/rename', () => HttpResponse.json({ status: 'ok' })),
  http.post('/api/sessions/delete', () => HttpResponse.json({
    id: 'session-fallback',
    title: '剩余会话',
    created_at: '2026-05-17T00:00:00Z',
    updated_at: '2026-05-17T00:00:00Z',
    active: true,
    message_count: 0,
  })),
  http.get('/api/chat/active', () => HttpResponse.json({ active: false })),
  http.get('/api/styles', () => HttpResponse.json({ styles: ['古龙.md', '番茄.txt'] })),
  http.post('/api/command', async ({ request }) => {
    const body = await request.json() as { command?: string }
    return HttpResponse.json({ result: `executed:${body.command || ''}` })
  }),
]
