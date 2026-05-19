import { DiffEditor } from '@monaco-editor/react'

export type ChapterDiffViewProps = {
  original: string
  modified: string
  language?: string
  sideBySide?: boolean
}

/** 章节差异视图，基于 Monaco Diff Editor 只读展示版本差异。 */
export function ChapterDiffView({
  original,
  modified,
  language = 'markdown',
  sideBySide = true,
}: ChapterDiffViewProps) {
  return (
    <div className="h-full min-h-[320px] w-full overflow-hidden">
      <DiffEditor
        height="100%"
        language={language}
        original={original}
        modified={modified}
        options={{
          readOnly: true,
          wordWrap: 'on',
          minimap: { enabled: false },
          renderSideBySide: sideBySide,
          scrollBeyondLastLine: false,
          automaticLayout: true,
        }}
      />
    </div>
  )
}
