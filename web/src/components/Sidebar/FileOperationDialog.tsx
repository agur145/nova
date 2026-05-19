import { useEffect, useState } from 'react'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'

export type FileOperationMode = 'create-file' | 'create-dir' | 'rename' | 'copy' | 'move'

interface FileOperationDialogProps {
  open: boolean
  mode: FileOperationMode
  targetPath: string
  defaultValue: string
  onOpenChange: (open: boolean) => void
  onSubmit: (value: string) => Promise<void>
}

const MODE_META: Record<FileOperationMode, { title: string; description: string; label: string }> = {
  'create-file': {
    title: '新建文件',
    description: '请输入相对 workspace 的文件路径。',
    label: '文件路径',
  },
  'create-dir': {
    title: '新建目录',
    description: '请输入相对 workspace 的目录路径。',
    label: '目录路径',
  },
  rename: {
    title: '重命名',
    description: '请输入新名称，不包含路径分隔符。',
    label: '新名称',
  },
  copy: {
    title: '复制',
    description: '请输入复制后的目标路径。',
    label: '目标路径',
  },
  move: {
    title: '移动',
    description: '请输入移动后的目标路径。',
    label: '目标路径',
  },
}

/** 文件操作弹窗，统一承载新建、重命名、复制和移动输入。 */
export function FileOperationDialog({
  open,
  mode,
  targetPath,
  defaultValue,
  onOpenChange,
  onSubmit,
}: FileOperationDialogProps) {
  const [value, setValue] = useState(defaultValue)
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const meta = MODE_META[mode]

  useEffect(() => {
    if (open) {
      setValue(defaultValue)
      setError('')
    }
  }, [defaultValue, open])

  const handleSubmit = async () => {
    const trimmed = value.trim()
    if (!trimmed || submitting) return
    setSubmitting(true)
    setError('')
    try {
      await onSubmit(trimmed)
      onOpenChange(false)
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e))
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="border-[#3a3d44] bg-[#25262a] text-[#d7dbe2]">
        <DialogHeader>
          <DialogTitle>{meta.title}</DialogTitle>
          <DialogDescription className="text-[#858b96]">
            {targetPath ? `当前目标：${targetPath}` : meta.description}
          </DialogDescription>
        </DialogHeader>
        <div className="space-y-2">
          <label className="text-xs text-[#aeb4bf]" htmlFor="file-operation-input">
            {meta.label}
          </label>
          <Input
            id="file-operation-input"
            value={value}
            onChange={(e) => setValue(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                e.preventDefault()
                void handleSubmit()
              }
            }}
            className="border-[#3a3d44] bg-[#1b1c1f] text-[#d7dbe2]"
            autoFocus
          />
          {error && <div className="text-xs text-red-400">{error}</div>}
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)} disabled={submitting}>
            取消
          </Button>
          <Button onClick={handleSubmit} disabled={!value.trim() || submitting}>
            确认
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
