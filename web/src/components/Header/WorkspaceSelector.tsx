import { useState, type ReactNode } from 'react'
import { FolderOpen, Trash2, Plus, Pencil, Check, X } from 'lucide-react'
import { removeBook, switchWorkspace, createBook, updateBookInfo, getBookInfo } from '@/lib/api'
import type { BookRecord, BookMeta } from '@/lib/api'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Button } from '@/components/ui/button'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { ScrollArea } from '@/components/ui/scroll-area'
import { TooltipIconButton } from '@/components/common/tooltip-icon-button'

interface WorkspaceSelectorProps {
  workspace: string
  books?: BookRecord[]
  onSwitch: (newPath: string) => void
  onBooksChange?: () => void
}

/** 计算相对时间描述 */
function relativeTime(isoStr: string): string {
  if (!isoStr) return ''
  const diff = Date.now() - new Date(isoStr).getTime()
  if (diff < 0) return '刚刚'
  const minutes = Math.floor(diff / 60000)
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}小时前`
  const days = Math.floor(hours / 24)
  if (days < 30) return `${days}天前`
  const months = Math.floor(days / 30)
  return `${months}月前`
}

/** 顶部工作区路径展示与切换组件 */
export function WorkspaceSelector({ workspace, books = [], onSwitch, onBooksChange }: WorkspaceSelectorProps) {
  const [editing, setEditing] = useState(false)
  const [inputPath, setInputPath] = useState('')
  const [switching, setSwitching] = useState(false)
  const [openBooks, setOpenBooks] = useState(false)

  // 新建书籍表单状态
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [createTitle, setCreateTitle] = useState('')
  const [createAuthor, setCreateAuthor] = useState('')
  const [createPath, setCreatePath] = useState('')
  const [createDesc, setCreateDesc] = useState('')
  const [createError, setCreateError] = useState('')
  const [creating, setCreating] = useState(false)

  // 内联编辑状态
  const [editingBookPath, setEditingBookPath] = useState<string | null>(null)
  const [editTitle, setEditTitle] = useState('')
  const [editAuthor, setEditAuthor] = useState('')
  const [editDesc, setEditDesc] = useState('')
  const [editSaving, setEditSaving] = useState(false)
  const [editLoading, setEditLoading] = useState(false)

  /** 计算当前 workspace 的父目录 */
  const parentDir = workspace ? workspace.replace(/\/[^/]*\/?$/, '') || '/' : '/'

  /** 执行切换 workspace */
  const handleSwitch = async (path?: string) => {
    const target = (path ?? inputPath).trim()
    if (!target) return
    setSwitching(true)
    try {
      const data = await switchWorkspace(target)
      onSwitch(data.workspace || target)
      setEditing(false)
      setOpenBooks(false)
      setInputPath('')
    } catch (e) {
      console.error('切换 workspace 失败', e)
    } finally {
      setSwitching(false)
    }
  }

  /** 移除最近书籍记录，不删除磁盘目录 */
  const handleRemoveBook = async (path: string) => {
    try {
      await removeBook(path)
      onBooksChange?.()
    } catch (e) {
      console.error('移除书籍记录失败', e)
    }
  }

  /** 打开新建书籍表单 */
  const openCreateForm = () => {
    setShowCreateForm(true)
    setCreateTitle('')
    setCreateAuthor('')
    setCreatePath(parentDir)
    setCreateDesc('')
    setCreateError('')
  }

  /** 提交新建书籍 */
  const handleCreate = async () => {
    if (!createTitle.trim()) { setCreateError('书名不能为空'); return }
    if (!createPath.trim()) { setCreateError('存放路径不能为空'); return }
    setCreating(true)
    setCreateError('')
    try {
      const data = await createBook(createPath.trim(), createTitle.trim(), createAuthor.trim() || undefined, createDesc.trim() || undefined)
      onSwitch(data.workspace)
      setShowCreateForm(false)
      setOpenBooks(false)
      onBooksChange?.()
    } catch (e: unknown) {
      setCreateError(e instanceof Error ? e.message : '创建失败')
    } finally {
      setCreating(false)
    }
  }

  /** 进入书籍编辑模式，先拉取完整元信息 */
  const startEditBook = async (book: BookRecord) => {
    setEditingBookPath(book.path)
    setEditTitle(book.name)
    setEditAuthor(book.author || '')
    setEditDesc('')
    setEditLoading(true)
    try {
      const meta: BookMeta = await getBookInfo(book.path)
      setEditTitle(meta.title)
      setEditAuthor(meta.author)
      setEditDesc(meta.description)
    } catch {
      // 回退使用 book record 里的基础信息
    } finally {
      setEditLoading(false)
    }
  }

  /** 保存书籍编辑 */
  const handleSaveEdit = async () => {
    if (!editingBookPath) return
    setEditSaving(true)
    try {
      await updateBookInfo(editingBookPath, editTitle.trim(), editAuthor.trim(), editDesc.trim())
      setEditingBookPath(null)
      onBooksChange?.()
    } catch (e) {
      console.error('保存书籍信息失败', e)
    } finally {
      setEditSaving(false)
    }
  }

  const inputCls = 'w-full rounded border border-[#3a3d44] bg-[#25262a] px-2 py-1 text-xs text-[#d7dbe2] outline-none focus:border-[#2f7dd3]'

  return (
    <div className="relative flex h-8 shrink-0 items-center gap-2 border-b border-[#303238] bg-[#1f2023] px-3 text-sm">
      <FolderOpen className="w-4 h-4 text-[#858b96] shrink-0" />
      {editing ? (
        <div className="flex items-center gap-1 flex-1 min-w-0">
          <Input
            type="text"
            value={inputPath}
            onChange={(e) => setInputPath(e.target.value)}
            onKeyDown={(e) => { if (e.key === 'Enter') handleSwitch(); if (e.key === 'Escape') setEditing(false) }}
            placeholder="输入工作区路径..."
            className="h-6 min-w-0 flex-1 border-[#3a3d44] bg-[#25262a] px-2 py-0.5 text-xs text-[#d7dbe2] focus-visible:border-[#2f7dd3] focus-visible:ring-0"
            autoFocus
            disabled={switching}
          />
          <Button
            type="button"
            size="xs"
            onClick={() => handleSwitch()}
            disabled={switching}
            className="bg-[#2f7dd3] text-white hover:bg-[#3b8eea]"
          >
            {switching ? '...' : '确定'}
          </Button>
          <Button
            type="button"
            size="xs"
            variant="ghost"
            onClick={() => { setEditing(false); setInputPath('') }}
            className="text-[#858b96] hover:bg-[#303238]"
          >
            取消
          </Button>
        </div>
      ) : (
        <>
          <span className="text-xs text-[#858b96] truncate flex-1 min-w-0">{workspace || '未设置工作区'}</span>
          <BookPopover
            open={openBooks}
            onOpenChange={setOpenBooks}
            trigger={(
              <Button
                type="button"
                size="xs"
                variant="ghost"
                className="flex shrink-0 items-center gap-1 text-[#aeb4bf] hover:bg-[#303238]"
              >
                书籍
              </Button>
            )}
          >
            {showCreateForm ? (
              /* ---- 新建书籍表单 ---- */
              <div className="space-y-2 p-2">
                <div className="text-xs font-medium text-[#c5c9d1]">新建书籍</div>
                <Input
                  type="text"
                  value={createTitle}
                  onChange={(e) => setCreateTitle(e.target.value)}
                  placeholder="书名（必填）"
                  className={inputCls}
                  autoFocus
                />
                <Input
                  type="text"
                  value={createAuthor}
                  onChange={(e) => setCreateAuthor(e.target.value)}
                  placeholder="作者（选填）"
                  className={inputCls}
                />
                <Input
                  type="text"
                  value={createPath}
                  onChange={(e) => setCreatePath(e.target.value)}
                  placeholder="存放路径（必填）"
                  className={inputCls}
                />
                <Textarea
                  value={createDesc}
                  onChange={(e) => setCreateDesc(e.target.value)}
                  placeholder="简介（选填）"
                  rows={3}
                  className={inputCls + ' min-h-0 resize-none'}
                />
                {createError && <div className="text-xs text-red-400">{createError}</div>}
                <div className="flex items-center justify-end gap-2">
                  <Button type="button" size="xs" variant="ghost" className="text-[#858b96] hover:bg-[#303238]" onClick={() => setShowCreateForm(false)}>取消</Button>
                  <Button type="button" size="xs" className="bg-[#2f7dd3] text-white hover:bg-[#3b8eea]" disabled={creating} onClick={handleCreate}>
                    {creating ? '创建中...' : '创建'}
                  </Button>
                </div>
              </div>
            ) : (
              /* ---- 书籍列表 ---- */
              <>
                <div className="flex items-center justify-between px-2 py-1.5">
                  <span className="text-xs font-medium text-[#c5c9d1]">最近书籍</span>
                  <Button
                    type="button"
                    size="xs"
                    variant="ghost"
                    className="flex items-center gap-1 text-[#aeb4bf] hover:bg-[#303238]"
                    onClick={openCreateForm}
                  >
                    <Plus className="h-3.5 w-3.5" />
                    新建
                  </Button>
                </div>
                {books.length === 0 ? (
                  <div className="px-2 py-2 text-xs text-[#858b96]">暂无记录</div>
                ) : (
                  <ScrollArea className="max-h-72">
                    <div className="pr-2">
                      {books.map((book) => {
                        const isCurrent = book.path === workspace
                        const isEditing = editingBookPath === book.path

                        if (isEditing) {
                          return (
                            <div key={book.path} className="space-y-1.5 rounded bg-[#25262a] p-2">
                              {editLoading ? (
                                <div className="py-2 text-center text-xs text-[#858b96]">加载中...</div>
                              ) : (
                                <>
                                  <Input
                                    type="text"
                                    value={editTitle}
                                    onChange={(e) => setEditTitle(e.target.value)}
                                    placeholder="书名"
                                    className={inputCls}
                                    autoFocus
                                  />
                                  <Input
                                    type="text"
                                    value={editAuthor}
                                    onChange={(e) => setEditAuthor(e.target.value)}
                                    placeholder="作者"
                                    className={inputCls}
                                  />
                                  <Textarea
                                    value={editDesc}
                                    onChange={(e) => setEditDesc(e.target.value)}
                                    placeholder="简介"
                                    rows={2}
                                    className={inputCls + ' min-h-0 resize-none'}
                                  />
                                  <div className="flex items-center justify-end gap-2">
                                    <TooltipIconButton
                                      label="取消"
                                      className="text-[#858b96] hover:bg-[#303238]"
                                      onClick={() => setEditingBookPath(null)}
                                    >
                                      <X className="h-3.5 w-3.5" />
                                    </TooltipIconButton>
                                    <TooltipIconButton
                                      label="保存"
                                      className="text-[#2f7dd3] hover:bg-[#2f7dd3]/15"
                                      disabled={editSaving}
                                      onClick={handleSaveEdit}
                                    >
                                      <Check className="h-3.5 w-3.5" />
                                    </TooltipIconButton>
                                  </div>
                                </>
                              )}
                            </div>
                          )
                        }

                        return (
                          <div
                            key={book.path}
                            className={`group relative flex items-start gap-2 rounded px-2 py-1.5 text-xs hover:bg-[#2f7dd3]/20 ${
                              isCurrent ? 'text-[#d7e8ff]' : 'text-[#c5c9d1]'
                            }`}
                          >
                            {isCurrent && (
                              <div className="absolute left-0 top-1.5 bottom-1.5 w-[3px] rounded-r bg-[#2f7dd3]" />
                            )}
                            <button
                              type="button"
                              className="min-w-0 flex-1 text-left pl-1"
                              onClick={() => handleSwitch(book.path)}
                              disabled={switching}
                            >
                              <div className="truncate font-semibold">{book.name || '未命名书籍'}</div>
                              <div className="mt-0.5 flex items-center gap-2 text-[11px] text-[#858b96]">
                                {book.author && <span>{book.author}</span>}
                                {book.last_opened_at && <span>{relativeTime(book.last_opened_at)}</span>}
                              </div>
                              <div className="mt-0.5 truncate text-[11px] text-[#5a5f6b]">{book.path}</div>
                            </button>
                            <div className="flex shrink-0 items-center gap-0.5 pt-0.5">
                              <TooltipIconButton
                                label="编辑信息"
                                className="text-[#858b96] opacity-0 hover:bg-[#2f7dd3]/15 hover:text-[#8bb8f0] group-hover:opacity-100"
                                onClick={() => startEditBook(book)}
                              >
                                <Pencil className="h-3.5 w-3.5" />
                              </TooltipIconButton>
                              <TooltipIconButton
                                label="移除记录"
                                className="text-[#858b96] opacity-0 hover:bg-red-500/15 hover:text-red-200 group-hover:opacity-100"
                                onClick={() => handleRemoveBook(book.path)}
                              >
                                <Trash2 className="h-3.5 w-3.5" />
                              </TooltipIconButton>
                            </div>
                          </div>
                        )
                      })}
                    </div>
                  </ScrollArea>
                )}
                <div className="mt-1 border-t border-[#303238] pt-1">
                  <Button
                    type="button"
                    size="sm"
                    variant="ghost"
                    className="w-full justify-start text-xs text-[#aeb4bf] hover:bg-[#303238]"
                    onClick={() => {
                      setOpenBooks(false)
                      setEditing(true)
                      setInputPath(workspace)
                    }}
                  >
                    添加/打开其他书籍目录...
                  </Button>
                </div>
              </>
            )}
          </BookPopover>
        </>
      )}
    </div>
  )
}

function BookPopover({
  open,
  onOpenChange,
  trigger,
  children,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
  trigger: ReactNode
  children: ReactNode
}) {
  return (
    <Popover open={open} onOpenChange={onOpenChange}>
      <PopoverTrigger asChild>{trigger}</PopoverTrigger>
      <PopoverContent
        align="start"
        side="bottom"
        className="z-50 w-[min(520px,calc(100vw-24px))] border-[#303238] bg-[#202124]/95 p-1 text-[#d7dbe2] shadow-[0_12px_32px_rgba(0,0,0,0.45)] backdrop-blur"
      >
        {children}
      </PopoverContent>
    </Popover>
  )
}
