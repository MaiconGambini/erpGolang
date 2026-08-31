# Design System

goERP uses a restrained light product UI: white, neutral grays, and blue for primary actions and active states.

## Brand Direction

- Clean ERP workspace.
- Dense but readable.
- Calm and trustworthy.
- Low visual noise.
- Product UI, not marketing UI.

## Colors

- Canvas: `#F8FAFC`
- Surface: `#FFFFFF`
- Surface muted: `#F1F5F9`
- Border: `#E2E8F0`
- Text primary: `#0F172A`
- Text secondary: `#475569`
- Text muted: `#64748B`
- Primary blue: `#2563EB`
- Primary hover: `#1D4ED8`
- Success: `#16A34A`
- Warning: `#D97706`
- Danger: `#DC2626`
- Info: `#0284C7`

## Typography

- Family: Inter or system UI.
- Page title: 24-28px.
- Section title: 18-20px.
- Body: 14-16px.
- Table text: 13-14px.
- Labels: 13-14px, medium weight.

## Spacing

- Base unit: 4px.
- Layout rhythm: 8px increments.
- Forms: 12-16px field gaps.
- Page content: 24-32px padding on desktop.

## Radius And Shadows

- Inputs/buttons: 8px.
- Cards/dialogs: 10-12px.
- Shadows are soft and minimal.

## Component Choices

- Data tables: semantic HTML tables with responsive horizontal scrolling.
- Text input/select: native controls styled from semantic tokens.
- Modal: shared `AppDialog` wrapper with accessible focus management.
- Loading/empty: shared `TableSkeleton` and `EmptyState` components.

## Themes

- Light mode is the default: white surfaces, neutral canvas, and blue action states.
- Dark mode is available through the `.dark` token override, persists in `localStorage`, and follows the system preference on first visit.
- Status, danger, focus, overlay, and raised-surface tokens are semantic and must be consumed instead of hardcoded component colors.

## Shared UI

- `AppButton` — primary action with disabled, hover, active, and focus states.
- `AppDialog` — modal semantics, Escape close, focus trap, focus restore, scroll lock, and reduced-motion transition.
- `EmptyState` and `TableSkeleton` — shared empty and loading states for data surfaces.
