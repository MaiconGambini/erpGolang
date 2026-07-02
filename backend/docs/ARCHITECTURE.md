# Backend Architecture

The backend is a modular monolith. Modules own their domain and database access. Shared packages must stay small and infrastructure-neutral.

MVP modules: tenants, users, auth, customers, audit.
