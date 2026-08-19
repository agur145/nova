# AgentChat Secondary Workspace — Design QA

**Findings**

- No actionable P0/P1/P2 differences remain.
- Accepted difference: the low-fidelity mock uses a bidirectional adjustment glyph, while the implementation uses Lucide `PanelRightOpen` / `PanelRightClose`. The implemented icon communicates the current pane state more directly and remains consistent with the existing icon family.

**Open Questions**

- The mock's compact frame is a conceptual landscape state rather than an exact device viewport. Responsive QA therefore compares hierarchy and control behavior at `800 × 900` CSS pixels instead of claiming pixel-level compact-frame fidelity.

**Implementation Checklist**

- [x] Persistent right workspace with independent tabs and active state.
- [x] First-open content chooser reusing existing Chat, Terminal, Reader, Lore, Presets, Skills, Agents, and Automations tab types.
- [x] Hide/show without unmounting secondary content or stopping background work.
- [x] Resizable desktop split with a balanced first-open width and persisted user width.
- [x] Toggle anchored to the rightmost visible tab strip in both single- and split-pane states.
- [x] Compact right-side drawer behavior.
- [x] Chinese/English copy, dark/light themes, accessible labels, pressed state, keyboard-accessible resize separator, and hidden-running indicator.

**Follow-up Polish**

- No blocking polish remains. A live running Agent was not started solely for visual capture; the hidden-running dot and label are covered by component tests.

## Evidence

- Source visual truth: `/Users/bytedance/.codex/generated_images/019fb8f6-3092-7373-a1a4-3efcc14c2b33/exec-020b4d72-78d9-4571-8f53-1a78a96df273.png`
- Final implementation screenshot: `/Users/bytedance/.codex/visualizations/2026/07/31/019fb8f6-3092-7373-a1a4-3efcc14c2b33/agentchat-secondary-final.png`
- Light-theme screenshot: `/Users/bytedance/.codex/visualizations/2026/07/31/019fb8f6-3092-7373-a1a4-3efcc14c2b33/agentchat-secondary-light.png`
- Compact closed state: `/Users/bytedance/.codex/visualizations/2026/07/31/019fb8f6-3092-7373-a1a4-3efcc14c2b33/agentchat-secondary-compact-closed.png`
- Compact drawer state: `/Users/bytedance/.codex/visualizations/2026/07/31/019fb8f6-3092-7373-a1a4-3efcc14c2b33/agentchat-secondary-compact.png`
- Full-view comparison, source left / implementation right: `/Users/bytedance/.codex/visualizations/2026/07/31/019fb8f6-3092-7373-a1a4-3efcc14c2b33/agentchat-secondary-desktop-comparison.png`
- Focused toggle comparison, source left / implementation right: `/Users/bytedance/.codex/visualizations/2026/07/31/019fb8f6-3092-7373-a1a4-3efcc14c2b33/agentchat-secondary-toggle-comparison.png`
- Compact comparison, source left / implementation right: `/Users/bytedance/.codex/visualizations/2026/07/31/019fb8f6-3092-7373-a1a4-3efcc14c2b33/agentchat-secondary-compact-comparison.png`

## Capture Normalization

- Source board: `1487 × 1058` pixels.
- Desktop source state crop: `732 × 484` pixels, contained in a `720 × 512` comparison cell without distortion.
- Desktop implementation: `1440 × 1024` CSS pixels and `1440 × 1024` captured pixels; effective density `1×`, contained in a `720 × 512` comparison cell without distortion.
- Compact source state crop: `732 × 480` pixels. Compact implementation: `800 × 900` CSS/captured pixels at `1×`. Both are contained, not stretched; the aspect-ratio mismatch is intentionally visible.
- Browser chrome is excluded. The screenshots contain only the Denova application viewport.

## State

- Route: Workspace / AgentChat.
- Project: local `示例书` fixture.
- Primary workspace: empty `新会话` Chat tab.
- Secondary workspace: populated Reader tab showing `第五章 妖兽袭击`.
- Desktop state: dark theme, split visible, balanced initial split, toggle at the far right of the secondary tab strip.
- Alternate states: secondary hidden and restored; light theme; `800 × 900` compact primary state and right drawer.

## Required Fidelity Surfaces

- Fonts and typography: implementation reuses Denova's configured UI and reading fonts, weights, line heights, truncation, and antialiasing. The visual hierarchy matches the existing product and the source intent; no new font or ad-hoc type scale was introduced.
- Spacing and layout rhythm: the Activity Navigator, primary workbench, divider, and secondary workbench preserve the source's three-region hierarchy. Existing tab-strip height, borders, radii, composer spacing, and resize handles are reused. The final toggle occupies the rightmost fixed tab-strip action slot in both visible states.
- Colors and visual tokens: all new surfaces use existing `--nova-*` tokens and shadcn variants. Dark and light captures show consistent borders, active fills, text contrast, focus treatment, and warning-state color.
- Image quality and asset fidelity: the target and implementation contain no new product imagery or decorative assets. Icons come from the existing Lucide dependency; no CSS art, custom SVG, emoji, or placeholder image was introduced.
- Copy and content: all new labels are concise and provided in Chinese and English. Dynamic book/chat content differs from the mock but represents realistic local data and does not alter the interaction hierarchy.

## Interaction and Accessibility Checks

- Opened the right-workspace menu and selected Reader.
- Verified menu parity with ordinary New Tab choices.
- Hid and restored the secondary workspace; Reader remained mounted and returned to the same state.
- Verified the toggle moves from the primary strip to the far-right secondary strip, then back to the far-right primary strip after hiding.
- Verified the resize separator is present and semantically labelled.
- Verified `800 × 900` compact layout shows a right-drawer launcher and opens/closes the `右侧工作区` dialog without squeezing the primary composer.
- Verified dark and light themes through the Settings UI, then restored the user's original dark theme.
- Verified button labels, `aria-pressed`, dialog semantics, close controls, tab roles, and separator labels in the browser accessibility snapshot.
- Browser console checked. No panel, tab, toggle, resize, theme, or responsive errors occurred. Existing local fixture sessions emit unrelated `会话不存在` history-load errors; these predate and are outside this change.

## Comparison History

1. Pass 1 — P2: a pane mounted hidden opened at its `280px` minimum instead of the intended balanced split. Evidence: `/Users/bytedance/.codex/visualizations/2026/07/31/019fb8f6-3092-7373-a1a4-3efcc14c2b33/agentchat-secondary-pass1-min-width.png`. Fix: `CollapsibleResizablePanel` now accepts an explicit first-open size when no persisted layout exists; persisted layouts still restore their saved size.
2. Pass 2 — P2: compact mode initially derived the toggle label from desktop visibility, so a closed drawer could present a hide action. Fix: compact visibility now follows `openPaneId`, with `openRight` / `closePane` actions. Pass 2 also exposed the toggle at the center divider in split mode. Evidence: `/Users/bytedance/.codex/visualizations/2026/07/31/019fb8f6-3092-7373-a1a4-3efcc14c2b33/agentchat-secondary-pass2-divider-toggle.png`.
3. Pass 3 — P2 resolved: control ownership now follows the rightmost visible workbench. Final desktop, focused-toggle, compact, and light-theme evidence show no remaining P0/P1/P2 issue.

## 2026-08-07 Toggle Refinement

- Source annotation: `/var/folders/3p/tw35s8456m1033g5yxdztf_m0000gn/T/codex-clipboard-b8cc6b73-24b7-4caa-a5e5-6140dd52346d.png` (`830 × 258` pixels).
- Final dark capture: `/Users/bytedance/.codex/visualizations/2026/08/07/019fdc15-9fca-7cc1-8198-33b7b7c37d2f/sidepanel-final-dark.jpg` (`1280 × 720` CSS/captured pixels; browser DPR `2`, normalized output `1×`).
- Light capture: `/Users/bytedance/.codex/visualizations/2026/08/07/019fdc15-9fca-7cc1-8198-33b7b7c37d2f/sidepanel-light-populated-folded.jpg`.
- Narrow capture: `/Users/bytedance/.codex/visualizations/2026/08/07/019fdc15-9fca-7cc1-8198-33b7b7c37d2f/sidepanel-narrow-populated-folded.jpg` at `780 × 720`; no horizontal overflow.
- Side-by-side focused comparison: `/Users/bytedance/.codex/visualizations/2026/08/07/019fdc15-9fca-7cc1-8198-33b7b7c37d2f/sidepanel-before-after.jpg`.

The folded populated pane now has a 6 px neutral presence indicator using `--nova-accent`. The existing warning marker remains the higher-priority busy state. The toggle uses the inline end-action presentation, removing its divider without changing the separated presentation used by other workbench actions.

Fonts, copy, icon family, and tab-strip height remain unchanged. Dark and light theme contrast, empty/open/folded states, hide/show behavior, and the 780 px layout passed. After restoring a transient local backend outage with the project restart script, the final browser interaction pass produced no new console errors. No P0/P1/P2 findings remain.

final result: passed

# Goal UI Design QA

## Target

- Source: the three Codex App screenshots supplied with the request.
- Scope: Goal entry mode in the composer and the durable goal card above it.
- Product boundaries: available in AgentChat and Writing Agent; absent from Game mode.

## Visual comparison

- Compared the supplied active-goal reference and the live Denova implementation together at desktop width.
- Desktop implementation evidence: `/Users/bytedance/.codex/visualizations/2026/08/09/019fe743-745b-77c1-bedb-f7c6db37871a/goal-desktop.png`.
- Narrow implementation evidence: `/Users/bytedance/.codex/visualizations/2026/08/09/019fe743-745b-77c1-bedb-f7c6db37871a/goal-narrow.png`.
- The implementation follows the same hierarchy: status and elapsed time on the left, compact actions on the right, objective beneath, then the composer immediately below.
- The Goal toolbar action uses the existing Denova composer typography, semantic surface colors, borders, radii, and Lucide target icon instead of introducing a parallel visual system.
- Selected Goal mode changes the action to a removable pill and replaces the composer placeholder with measurable-goal guidance.
- Denova intentionally retains its existing 8 px spacing used by queued/steer cards and its 56 rem adaptive composer width rather than copying Codex-specific chrome.

## Interaction verification

- Verified Goal mode selection, goal-specific placeholder, edit prefill, pause, clear, collapse controls, and elapsed-time rendering in the live app.
- Verified a long Chinese objective wraps without clipping and all actions remain reachable at a 720 x 900 viewport.
- Verified the Goal action is present in AgentChat and Writing Agent and absent from the Game composer.
- Fixed the floating-composer pointer-event boundary found during browser QA so card actions remain clickable without allowing the floating container to block the message surface.

## Theme and accessibility

- Uses existing semantic theme variables only; no light-only or dark-only literal colors were introduced.
- Every icon action has a localized accessible label, the card is a named region, and the mode toggle exposes pressed state.
- Chinese and English strings cover commands, statuses, actions, placeholders, errors, and capability descriptions.

## Notes

- Live light-theme comparison passed. Switching the application setting to Dark did not update the root `data-theme` value in the current development session; this pre-existing application theme propagation issue is outside the Goal implementation. Static theme-variable inspection and frontend tests confirm the Goal components do not bypass the shared theme system.

# Goal and Plan Composer Toggle Design QA

## Target and evidence

- Source: `/var/folders/3p/tw35s8456m1033g5yxdztf_m0000gn/T/codex-clipboard-98f9bd26-aac0-4ebc-9b56-6d9a1a4b86f4.png` and the two companion Codex App screenshots supplied with the request.
- Final focused comparison: `/Users/bytedance/.codex/visualizations/2026/08/10/goal-plan-toggle/goal-active-focused.png`.
- Final desktop captures: `/Users/bytedance/.codex/visualizations/2026/08/10/goal-plan-toggle/goal-active-desktop.png` and `/Users/bytedance/.codex/visualizations/2026/08/10/goal-plan-toggle/plan-active-desktop.png` at `1280 × 720` CSS pixels.
- Responsive and theme captures: `/Users/bytedance/.codex/visualizations/2026/08/10/goal-plan-toggle/goal-active-narrow.png` at `720 × 900` and `/Users/bytedance/.codex/visualizations/2026/08/10/goal-plan-toggle/goal-active-light.png`; screenshots were normalized to `1×`.

## Visual and interaction checks

- Inactive state has no persistent Goal or Plan indicator in the composer toolbar.
- Input Actions lists Goal first and Plan second. Both use the existing checkbox-menu pattern and localized labels.
- Active Goal and Plan share the same removable pill immediately to the right of the permission control, matching the source hierarchy while retaining Denova's existing spacing, typography, semantic colors, and Lucide icon family.
- Goal and Plan are strictly mutually exclusive through the menu, `/goal`, `/plan`, and `Shift+Tab` entry paths. Closing the active pill returns focus to the composer.
- Goal remains available in AgentChat and Writing Agent, and is absent from Game mode, including its menu item, slash command, active indicator, and durable status card.
- Verified dark and light themes, desktop and narrow layouts, empty/inactive states, active Goal and Plan states, long composer width, focus restoration, and localized accessible labels. Browser console contained no warnings or errors from these interactions.
- Restored the user's original `继承（深色）` theme setting after verification.

## Comparison history

1. Previous iteration — P2 product mismatch: Goal was always visible in the toolbar and Plan used a separate hint presentation. The updated requirement made both modes explicit, transient composer states.
2. Final iteration — resolved: Goal moved to the top of Input Actions, active Goal and Plan reuse one removable indicator, and all entry paths enforce mutual exclusion. No P0/P1/P2 findings remain.

final result: passed

# Diff Review Design QA

Final result: passed

## Compared artifacts

- Reference: `/var/folders/3p/tw35s8456m1033g5yxdztf_m0000gn/T/codex-clipboard-f926cb6d-2d76-4fa8-9930-3b2e17db8779.png`
- Dark wide implementation: `/Users/bytedance/.codex/visualizations/2026/08/13/019ffb1b-19b2-7402-bc46-f51931fdff89/diff-review-dark-wide.png`
- Dark narrow implementation: `/Users/bytedance/.codex/visualizations/2026/08/13/019ffb1b-19b2-7402-bc46-f51931fdff89/diff-review-dark-narrow.png`
- Light wide implementation: `/Users/bytedance/.codex/visualizations/2026/08/13/019ffb1b-19b2-7402-bc46-f51931fdff89/diff-review-light-wide.png`
- Side-by-side comparison: `/Users/bytedance/.codex/visualizations/2026/08/13/019ffb1b-19b2-7402-bc46-f51931fdff89/diff-review-comparison.png`

The reference is directional rather than a pixel-identical Denova screen: it shows a single-file large-diff mode without Denova's Agent panel. The comparison therefore checks the requested hierarchy, density, changed-file tree, and status treatment rather than unsupported pixel parity.

## Verification

| Criterion | Result | Evidence |
| --- | --- | --- |
| Changed files only | Pass | The review tree is projected exclusively from the review thread's file list. |
| Hierarchical file tree | Pass | Directories are grouped, expanded initially, and toggle from the whole row. |
| Added / modified / deleted states | Pass | Localized green plus, amber pencil, and red minus badges are rendered; conflict and accepted states remain visible. |
| Compact density | Pass | The tree uses 24 px virtualized rows, 12 px indentation, a compact filter, and a narrower adaptive pane. |
| Long paths | Pass | Directory hierarchy shortens visible paths and the shared filename treatment preserves extensions while truncating stems. |
| Narrow layout | Pass | At 720 x 900, the full tree collapses to the existing changed-file dropdown with no document-level horizontal overflow. |
| Light and dark themes | Pass | The review, tree selection, status colors, and diff surfaces were inspected in both themes. |
| Opening motion | Pass | The underlying editor remains mounted, outer pane geometry settles in one frame, and the review enters with a 220 ms opacity/3 px compositor transition. Reduced/off motion resolves in one frame. |
| Existing behavior | Pass | File jumping, directory toggling, filtering, selection, compact navigation, and Project explorer reuse were exercised in the app and covered by focused tests. |

## Intentional differences

- Denova keeps its existing review toolbar, full-thread stacked diffs, and optional Agent panel instead of copying the reference application's branch and commit controls.
- Status badges use Denova's semantic success/warning/danger tokens so they remain legible in both supported themes.
