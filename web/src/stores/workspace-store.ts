import { create } from 'zustand'

export type RightPanel = 'ai' | 'outline' | 'characters' | 'versions' | null
export type BottomPanel = 'tasks' | 'versions' | 'problems' | null

type WorkspaceStore = {
  selectedProjectId?: string
  selectedChapterId?: string
  rightPanel: RightPanel
  bottomPanel: BottomPanel
  commandOpen: boolean
  setSelectedProjectId: (id?: string) => void
  setSelectedChapterId: (id?: string) => void
  setRightPanel: (panel: RightPanel) => void
  setBottomPanel: (panel: BottomPanel) => void
  setCommandOpen: (open: boolean) => void
}

/** 工作区 UI 状态 Store，仅保存本地界面状态，不存放服务端数据。 */
export const useWorkspaceStore = create<WorkspaceStore>((set) => ({
  selectedProjectId: undefined,
  selectedChapterId: undefined,
  rightPanel: 'ai',
  bottomPanel: 'tasks',
  commandOpen: false,
  setSelectedProjectId: (id) => set({ selectedProjectId: id }),
  setSelectedChapterId: (id) => set({ selectedChapterId: id }),
  setRightPanel: (panel) => set({ rightPanel: panel }),
  setBottomPanel: (panel) => set({ bottomPanel: panel }),
  setCommandOpen: (open) => set({ commandOpen: open }),
}))
