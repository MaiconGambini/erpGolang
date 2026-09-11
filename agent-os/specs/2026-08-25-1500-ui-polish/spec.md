# Feature Spec — UI polish para portfólio

Date: 2026-08-25 · Lane: P1 "Front/UI improvements" · Método: impeccable `polish` (refinamento — preserva o mundo visual incumbente: slate + azul, Inter, tokens CSS)

## Objective

UI em nível de portfólio: gráficos vivos com dados reais, estados de loading/vazio completos, dark mode com tokens, e superfícies de browser tematizadas — mantendo E2E verde e re-capturando os screenshots do README.

## Requirements

- REQ-001: Dashboard mostra gráficos com dados no range padrão — default 30d→90d + `cmd/seed` ganha dados demo idempotentes (clientes, produtos, vendas confirmadas espalhadas nos últimos ~45 dias, drafts).
- REQ-002: Tabelas e gráficos têm estados completos: skeleton durante carga, empty state contextual com CTA quando vazio (componente compartilhado).
- REQ-003: Dark mode via overrides de tokens (`.dark` em tokens.scss), toggle no topbar (lucide Sun/Moon), persistência em localStorage, default `prefers-color-scheme`, sem FOUC (script inline no index.html). Cores hardcoded (ex.: `#eff6ff` no AppShell) viram tokens.
- REQ-004: Superfícies de browser tematizadas: seleção de texto, scrollbar, focus-visible consistente, `tabular-nums` em tabelas/KPIs.
- REQ-005: E2E (15 testes) permanece verde; `vue-tsc` e build limpos.

## Acceptance Criteria

- [x] Dashboard com gráficos populados no primeiro load (screenshot re-capturado, range 90d).
- [x] Toggle de tema alterna sem flash (script inline no index.html); preferência sobrevive reload; respeita prefers-color-scheme na primeira visita.
- [x] Listas mostram skeleton durante carga e empty state com hint contextual (6 superfícies).
- [x] `npm run typecheck` + `build` + `test:unit` (14) + Playwright E2E (21) passam.
- [x] detect.mjs executado — 1 warning (fonte Inter): exceção intencional documentada (identidade incumbente; troca = redesign). Inspeção batched desktop+mobile: sidebar não colapsava ≤768px — corrigido e confirmado.
- [x] 4 screenshots re-capturados sobre dados demo limpos (dark, gráficos vivos).
