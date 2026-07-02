# UX Patterns

## CRUD Page Layout

```text
Header: title, description, primary action
Toolbar: search, filters, export
Content: table or responsive card list
Footer: pagination
```

## Loading

- Use skeletons for table/page content.
- Use button-level loading for submit actions.
- Use spinner only for blocking operations.

## Empty State

- Short title.
- Helpful explanation.
- CTA when the user has permission.

## Form Validation

- Field errors appear inline below fields.
- Global failures can use toast plus inline summary.
- Toast must not be the only feedback for critical errors.

## Delete Confirmation

- Use a confirmation dialog.
- Include the entity name.
- Use destructive styling only on the confirm action.

## Create/Edit Success

- Show success toast.
- Close dialog.
- Invalidate/refetch the list query.

## Pagination

- Server-side pagination.
- Default page size: 20.
- API uses `limit` and `offset`.

## Search

- Debounce search by 300ms.
- Show inline loading while results update.
