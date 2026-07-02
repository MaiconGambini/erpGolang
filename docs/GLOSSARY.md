# Glossary

## Tenant

A tenant is an isolated company/account using goERP. Tenant-owned data must include `tenant_id`.

## User

A user is a person who signs in to a tenant. Users have roles such as admin, manager, operator, or viewer.

## Customer

A customer is an individual or company that buys from a tenant. Customers are tenant-scoped and can be active or inactive.

Does not include suppliers or internal users.

## Supplier

A supplier provides goods or services to a tenant. Suppliers are out of MVP scope.

## Product

A product is an item or service sold or managed by a tenant. Products are out of MVP scope.

## Order

An order represents a commercial transaction or workflow. Orders are out of MVP scope.

## Invoice

An invoice is a fiscal or billing document. Invoices are out of MVP scope.

## Payment

A payment records money movement for an order or invoice. Payments are out of MVP scope.

## Stock

Stock is the quantity of a product available to a tenant. Stock is out of MVP scope.

## Audit Log

An audit log is an immutable record of a write action, including tenant, actor, action, entity, and timestamp.
