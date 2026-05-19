import { useState } from 'react'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'

interface DeleteConfirmDialogProps {
  open: boolean
  path: string
  onOpenChange: (open: boolean) => void
  onConfirm: () => Promise<void>
}

/** 删除确认弹窗，避免误删 workspace 文件。 */
export function DeleteConfirmDialog({ open, path, onOpenChange, onConfirm }: DeleteConfirmDialogProps) {
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const handleConfirm = async () => {
    setSubmitting(true)
    setError('')
    try {
      await onConfirm()
      onOpenChange(false)
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e))
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent className="border-[#3a3d44] bg-[#25262a] text-[#d7dbe2]">
        <AlertDialogHeader>
          <AlertDialogTitle>确认删除？</AlertDialogTitle>
          <AlertDialogDescription className="text-[#858b96]">
            将永久删除 `{path}`，此操作不可撤销。
          </AlertDialogDescription>
        </AlertDialogHeader>
        {error && <div className="text-xs text-red-400">{error}</div>}
        <AlertDialogFooter>
          <AlertDialogCancel disabled={submitting}>取消</AlertDialogCancel>
          <AlertDialogAction
            className="bg-red-600 text-white hover:bg-red-500"
            disabled={submitting}
            onClick={(e) => {
              e.preventDefault()
              void handleConfirm()
            }}
          >
            删除
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
