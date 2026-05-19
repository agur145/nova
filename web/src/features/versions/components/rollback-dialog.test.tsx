import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { RollbackDialog } from './rollback-dialog'
import type { VersionItem } from './version-timeline'

const version: VersionItem = {
  id: 'abcdef123456',
  title: '完善第一章',
  description: 'abcdef1',
  createdAt: '2026-05-17',
  author: 'tester',
}

describe('RollbackDialog', () => {
  it('calls onRollback after confirm', async () => {
    const user = userEvent.setup()
    const handleRollback = vi.fn()

    render(
      <RollbackDialog
        open
        version={version}
        onOpenChange={vi.fn()}
        onRollback={handleRollback}
      />,
    )

    await user.click(screen.getByRole('button', { name: '确认回滚' }))

    expect(handleRollback).toHaveBeenCalledWith(version)
  })
})
