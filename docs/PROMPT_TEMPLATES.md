# Prompt Templates

## New Module

```text
Context:
- Read docs/AI_CONTEXT.md
- Read ARCHITECTURE.md
- Read CONVENTIONS.md
- Read docs/MODULE_TEMPLATE.md

Task:
Build module: [module name]

Fields:
- [field list]

Business rules:
- [rules]

Endpoints:
- [REST endpoints]

UI:
- [screens and states]

Definition of done:
- Backend compiles
- Frontend compiles
- Tests pass
- Tenant isolation verified
- Docs updated
```

## Bug Fix

```text
Investigate and fix: [bug]
Reproduce first if feasible.
Do not change unrelated behavior.
Add or update regression tests.
Validate with the smallest command that proves the fix.
```

## UI Screen

```text
Create a Vue/PrimeVue screen following docs/DESIGN_SYSTEM.md and docs/UX_PATTERNS.md.
Use pt-BR UI copy.
Implement loading, empty, error, and success states.
Keep FSD import direction intact.
```

## Deployment

```text
Prepare deployment changes following DEPLOYMENT.md.
Do not introduce Kubernetes, Terraform, or ECS unless explicitly requested.
Validate with local Docker Compose config and health checks.
```
