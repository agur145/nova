import type { VersionItem } from './version-timeline'
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

interface RollbackDialogProps {
  open: boolean
  version?: VersionItem | null
  loading?: boolean
  onOpenChange: (open: boolean) => void
  onRollback: (version: VersionItem) => void | Promise<void>
}

/** 回滚确认弹窗，通过二次确认避免误重置当前作品工作区。 */
export function RollbackDialog({
  open,
  version,
  loading = false,
  onOpenChange,
  onRollback,
}: RollbackDialogProps) {
  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent className="border-[#3a3d44] bg-[#25262a] text-[#d7dbe2]">
        <AlertDialogHeader>
          <AlertDialogTitle>确认回滚版本？</AlertDialogTitle>
          <AlertDialogDescription className="text-[#858b96]">
            {version
              ? `将整本书回滚到版本 ${version.description || version.title}。此操作会重置当前工作区到该版本。`
              : '请选择要回滚的版本。'}
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel disabled={loading}>取消</AlertDialogCancel>
          <AlertDialogAction
            className="bg-[#ffbd5e] text-[#18191b] hover:bg-[#ffd28a]"
            disabled={loading || !version}
            onClick={(event) => {
              event.preventDefault()
              if (version) void onRollback(version)
            }}
          >
            确认回滚
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
