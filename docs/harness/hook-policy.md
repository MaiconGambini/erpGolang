# Harness Hook Policy

## Profiles

- `minimal`: blocks secret paths and destructive commands.
- `standard`: minimal plus context, activity, and config warnings.
- `strict`: standard plus mandatory scan/status checks before confirmation.

## Escape Hatches

- `HARNESS_PROFILE=minimal`
- `HARNESS_DISABLED_GUARDS=guard-id-a,guard-id-b`

Do not weaken guards silently. Record the reason in handoff evidence.
