import { useHotkeys } from 'react-hotkeys-hook'

interface WorkspaceHotkeysOptions {
  onSave?: () => void | Promise<void>
  onOpenCommand?: () => void
  onGenerate?: () => void
  onOpenDiff?: () => void
  onEscape?: () => void
}

/** 注册工作台全局快捷键，hook 只分发事件，不直接调用业务 API。 */
export function useWorkspaceHotkeys({
  onSave,
  onOpenCommand,
  onGenerate,
  onOpenDiff,
  onEscape,
}: WorkspaceHotkeysOptions) {
  useHotkeys('meta+s, ctrl+s', (event) => {
    event.preventDefault()
    void onSave?.()
  }, { enableOnFormTags: false }, [onSave])

  useHotkeys('meta+k, ctrl+k', (event) => {
    event.preventDefault()
    onOpenCommand?.()
  }, { enableOnFormTags: true }, [onOpenCommand])

  useHotkeys('meta+enter, ctrl+enter', (event) => {
    event.preventDefault()
    onGenerate?.()
  }, { enableOnFormTags: true }, [onGenerate])

  useHotkeys('meta+shift+d, ctrl+shift+d', (event) => {
    event.preventDefault()
    onOpenDiff?.()
  }, { enableOnFormTags: true }, [onOpenDiff])

  useHotkeys('esc', () => {
    onEscape?.()
  }, { enableOnFormTags: true }, [onEscape])
}
