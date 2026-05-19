import { ChapterDiffView } from '@/features/chapters/components/chapter-diff-view'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'

interface VersionDiffDialogProps {
  open: boolean
  title?: string
  original: string
  modified: string
  language?: string
  sideBySide?: boolean
  onOpenChange: (open: boolean) => void
}

/** 版本差异弹窗，仅展示外部传入的 diff 内容。 */
export function VersionDiffDialog({
  open,
  title = '版本差异',
  original,
  modified,
  language,
  sideBySide,
  onOpenChange,
}: VersionDiffDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="h-[80vh] max-w-6xl border-[#3a3d44] bg-[#25262a] text-[#d7dbe2]">
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription className="text-[#858b96]">
            只读展示当前版本与目标版本的内容差异。
          </DialogDescription>
        </DialogHeader>
        <div className="min-h-0 flex-1 overflow-hidden rounded border border-[#303238]">
          <ChapterDiffView
            original={original}
            modified={modified}
            language={language}
            sideBySide={sideBySide}
          />
        </div>
      </DialogContent>
    </Dialog>
  )
}
