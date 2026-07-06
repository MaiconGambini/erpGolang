# goERP Documentation

Human-facing documentation index. For AI/harness tooling see [`AGENTS.md`](../AGENTS.md) at repo root.

## Start here

| Doc | Purpose |
|-----|---------|
| [../README.md](../README.md) | Setup, stack, screenshots, portfolio overview |
| [PRODUCT.md](PRODUCT.md) | Product scope, journeys, non-goals |
| [ARCHITECTURE.md](ARCHITECTURE.md) | System design, modules, routes |
| [../ARCHITECTURE.md](../ARCHITECTURE.md) | High-level architecture summary |
| [BUSINESS_RULES.md](BUSINESS_RULES.md) | Domain invariants |
| [GLOSSARY.md](GLOSSARY.md) | Domain terms |
| [ROLES.md](ROLES.md) | RBAC permission matrix |
| [UX_PATTERNS.md](UX_PATTERNS.md) | UI behavior, export, charts |
| [MODULE_TEMPLATE.md](MODULE_TEMPLATE.md) | How to add a module |
| [CI_CD.md](CI_CD.md) | Pipelines and deploy |
| [../SECURITY.md](../SECURITY.md) | Vulnerability reporting |

## Backend runbooks

| Doc | Purpose |
|-----|---------|
| [../backend/docs/LOCAL_DEV.md](../backend/docs/LOCAL_DEV.md) | Local runbook (ports 5434/6381) |
| [../backend/docs/ARCHITECTURE.md](../backend/docs/ARCHITECTURE.md) | Route inventory, layering |
| [../backend/docs/AUTH_SESSION.md](../backend/docs/AUTH_SESSION.md) | Auth and sessions |
| [../backend/docs/TENANT_ISOLATION.md](../backend/docs/TENANT_ISOLATION.md) | Multi-tenant rules |
| [../backend/docs/SALES_TRANSACTIONS.md](../backend/docs/SALES_TRANSACTIONS.md) | Sales workflow invariants |
| [../backend/docs/DASHBOARD_AGGREGATES.md](../backend/docs/DASHBOARD_AGGREGATES.md) | KPI definitions |
| [../backend/docs/SCHEMA.md](../backend/docs/SCHEMA.md) | Table reference |
| [../backend/docs/RATE_LIMITING.md](../backend/docs/RATE_LIMITING.md) | Login rate limits |
| [../backend/docs/DEPLOY_READINESS.md](../backend/docs/DEPLOY_READINESS.md) | Production checklist |

## Frontend & design

| Doc | Purpose |
|-----|---------|
| [DESIGN_SYSTEM.md](DESIGN_SYSTEM.md) | Tokens and components |
| [../CONVENTIONS.md](../CONVENTIONS.md) | Naming, API shapes, commits |

## Deploy

| Doc | Purpose |
|-----|---------|
| [../DEPLOYMENT.md](../DEPLOYMENT.md) | VPS and production deploy |
| [RELIABILITY.md](RELIABILITY.md) | Health, readiness, ops notes |

## Portfolio 8.5+ feature map

| Feature | Docs |
|---------|------|
| RBAC | [ROLES.md](ROLES.md), [BUSINESS_RULES.md](BUSINESS_RULES.md#rbac) |
| Users admin | [PRODUCT.md](PRODUCT.md), [ROLES.md](ROLES.md) |
| Audit viewer | [GLOSSARY.md](GLOSSARY.md#audit-log) |
| CSV/PDF/charts | [UX_PATTERNS.md](UX_PATTERNS.md#reporting), [GLOSSARY.md](GLOSSARY.md#report) |
| GitHub polish | [../README.md](../README.md), [../LICENSE](../LICENSE) |
