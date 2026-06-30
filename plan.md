Stack final consolidada
Backend

Linguagem: Go 1.23+
Router: chi/v5
DB driver: pgx/v5 (pool nativo, sem database/sql)
Queries: sqlc (gera código tipado)
Migrations: atlas (declarativo, diff automático)
Cache/sessions: Redis (go-redis/v9)
Auth: golang-jwt/jwt/v5 (HS256) + bcrypt
Validação: go-playground/validator/v10
Config: caarlos0/env/v10 + godotenv (dev)
Logger: slog (stdlib) com handler JSON
Dinheiro: shopspring/decimal
UUID: google/uuid
Jobs assíncronos (quando precisar): riverqueue/river
Hot reload (dev): air
Testes: testing (stdlib) + testify/require + testcontainers-go (Postgres real em teste)

Frontend

Build: Vite + Vue 3 + TypeScript
Router: vue-router 4
State: pinia
Server state: @tanstack/vue-query
HTTP: axios
Forms: vee-validate + zod
UI lib: PrimeVue (DataTable, Calendar, etc. — ERP-ready)
CSS: TailwindCSS (utility) + SCSS pontual
Ícones: lucide-vue-next
Datas: date-fns
Arquitetura: Feature-Sliced Design
Linter: ESLint + Prettier
Testes: Vitest (unit) + Vue Test Utils + Playwright (e2e)

Infra

DB: Postgres 16
Cache: Redis 7
Container: Docker + Docker Compose (dev)
Reverse proxy: Caddy (mais simples que Nginx em prod pequena)
Deploy: Docker em VPS ou Fly.io / Railway
CI/CD: GitHub Actions
Monitoramento mínimo: logs estruturados → arquivo → Loki (opcional fase 2)


Fase 0 — Decisões já fechadas
Multi-tenant via tenant_id, JWT (access 15min + refresh em cookie httpOnly), sqlc + Atlas, primeiro domínio customers. Adiciono três decisões novas que valem fechar antes de começar:

Monorepo com backend/ + frontend/ na raiz. git, README, docker-compose.yml na raiz.
Versionamento de API via path: /api/v1/*. Define agora pra não mexer depois.
Identificação de tenant: pelo tenant_slug no body do login + claim tenant_id no JWT. Sem subdomínio no MVP (mais simples).


Fase 1 — Setup do monorepo
Estrutura raiz:
erp/
├── backend/
├── frontend/
├── docker-compose.yml      # postgres + redis pra dev
├── .editorconfig
├── .gitignore
└── README.md
docker-compose.yml sobe Postgres 16 e Redis 7 com volumes nomeados. Backend e front rodam fora do compose em dev (mais rápido pra hot reload).
Entregável: docker compose up -d sobe os serviços, psql conecta no Postgres, redis-cli ping responde.

Fase 2 — Esqueleto do backend
Estrutura final (essa não muda mais):
backend/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── app/
│   │   ├── config/         # struct + load env
│   │   ├── database/       # pool pgx + healthcheck
│   │   ├── cache/          # client Redis
│   │   ├── jwt/            # gerar + validar tokens
│   │   ├── middleware/     # logger, requestid, cors, auth, tenant, recovery
│   │   ├── httpx/          # response helpers, errors, paginação
│   │   ├── validator/      # wrapper do go-playground/validator
│   │   ├── tenant/         # context helpers
│   │   └── audit/          # serviço de audit log compartilhado
│   ├── modules/
│   │   ├── module.go       # interface Module
│   │   └── register.go     # slice de módulos a montar
│   └── routes/
│       └── routes.go       # monta tudo no chi.Router
├── migrations/             # gerado pelo Atlas
├── schema.sql              # fonte da verdade do schema
├── atlas.hcl
├── sqlc.yaml
├── .air.toml
├── .env.example
├── Dockerfile
├── Makefile
└── go.mod
Makefile com os comandos do dia-a-dia: make dev (air), make migrate (atlas apply), make migrate-diff name=xxx (atlas diff), make sqlc (gera código), make test, make lint.
Entregável: make dev sobe o servidor com hot reload, /api/v1/health retorna 200.

Fase 3 — Fundação do backend
Implementação na ordem (cada item vira um PR pequeno):

Config tipada — struct anotada com env:, falha na inicialização se algo crítico tá faltando.
Pool Postgres — pgxpool.New, healthcheck, logs de queries lentas.
Cliente Redis — com ping na boot.
JWT service — GenerateAccessToken(userID, tenantID, role), GenerateRefreshToken(...), Validate(token). Refresh sempre persiste jti no Redis com TTL.
Response helpers (httpx) — JSON(w, status, data), Error(w, code, message), Paginated(w, data, meta). Formato padrão definido.
Erros estruturados — tipo httpx.AppError{Code, Message, Status, Details}. Helper httpx.WriteError(w, err).
Validador — wrapper que transforma erro do validator em details no response.
Middlewares: RequestID → Logger → Recovery → CORS → (rotas públicas/privadas).
AuthJWT middleware — extrai claims, injeta no contexto.
TenantScope middleware — extrai tenant_id do contexto JWT, injeta como tipo próprio (tenant.ID).
Context helpers — tenant.FromContext(ctx) (uuid.UUID, error) e tenant.MustFromContext(ctx) (panica em dev, retorna erro em prod).
Audit service — audit.Log(ctx, action, entity, entityID, before, after). Tabela audit_logs criada nessa fase.
Interface Module:

go    type Module interface {
        Name() string
        Register(r chi.Router, deps Deps)
    }
    type Deps struct {
        DB     *pgxpool.Pool
        Redis  *redis.Client
        JWT    *jwt.Service
        Audit  *audit.Service
        Logger *slog.Logger
    }

Boot do app — routes.go faz for _, m := range modules.All { m.Register(r, deps) }.

Entregável: estrutura pronta pra plugar módulos, nenhum módulo de negócio ainda.

Fase 4 — Módulo auth + users + tenants
Aplica o template do módulo. Tabelas e endpoints como definidos antes. Pontos da arquitetura que valem reforçar:

auth e users são módulos separados. Auth lida com login/refresh/logout/me. Users lida com CRUD de usuários (admin cria, lista, desativa).
tenants é módulo mínimo no MVP — só GET /tenants/current (info do tenant logado). Cadastro de novo tenant fica fora do MVP (ou faz via seed/CLI).
Seed CLI em cmd/seed/main.go: cria tenant + admin. Roda com go run ./cmd/seed.

Entregável: dois tenants no banco via seed, login funcional via curl, /auth/me retorna user correto, isolamento entre tenants validado.

Fase 5 — Esqueleto do frontend (FSD de verdade)
Estrutura completa do frontend/src/:
src/
├── app/
│   ├── App.vue
│   ├── main.ts
│   ├── providers/
│   │   ├── router.ts          # cria router + guards
│   │   ├── pinia.ts
│   │   ├── query.ts           # Vue Query
│   │   ├── primevue.ts
│   │   └── index.ts           # composer de providers
│   └── styles/
│       ├── app.scss
│       └── tailwind.css
├── processes/
│   └── auth/                   # fluxo: refresh no boot + redirect
│       └── boot-auth.ts
├── pages/
│   ├── login/
│   │   └── ui/LoginPage.vue
│   ├── dashboard/
│   │   └── ui/DashboardPage.vue
│   ├── customers/
│   │   └── ui/CustomersPage.vue
│   └── not-found/
│       └── ui/NotFoundPage.vue
├── widgets/
│   ├── app-shell/             # layout autenticado (sidebar + topbar)
│   │   ├── ui/AppShell.vue
│   │   ├── ui/Sidebar.vue
│   │   └── ui/Topbar.vue
│   └── customer-table/
│       └── ui/CustomerTable.vue
├── features/
│   ├── auth/
│   │   ├── login/
│   │   │   ├── ui/LoginForm.vue
│   │   │   └── model/use-login.ts
│   │   └── logout/
│   │       └── model/use-logout.ts
│   └── customer/
│       ├── create-customer/
│       │   ├── ui/CreateCustomerDialog.vue
│       │   └── model/use-create-customer.ts
│       ├── edit-customer/
│       │   ├── ui/EditCustomerDialog.vue
│       │   └── model/use-edit-customer.ts
│       ├── delete-customer/
│       │   └── model/use-delete-customer.ts
│       └── list-customers/
│           ├── ui/CustomerSearchBar.vue
│           └── model/use-list-customers.ts
├── entities/
│   ├── session/
│   │   ├── model/session.store.ts    # Pinia: accessToken, user
│   │   ├── model/types.ts
│   │   └── api/session.api.ts
│   ├── user/
│   │   ├── model/types.ts
│   │   └── api/user.api.ts
│   └── customer/
│       ├── model/types.ts            # interface Customer
│       ├── model/schemas.ts          # zod schemas
│       └── api/customer.api.ts
└── shared/
    ├── api/
    │   ├── client.ts                 # axios + interceptor refresh
    │   └── types.ts                  # ApiResponse<T>, Paginated<T>, ApiError
    ├── config/
    │   └── env.ts                    # import.meta.env tipado
    ├── lib/
    │   ├── format/                   # currency, date, document (CPF/CNPJ)
    │   └── jwt/                      # decode (sem confiar)
    ├── ui/                           # componentes "burros" reutilizáveis
    │   ├── FormField.vue
    │   ├── PageHeader.vue
    │   └── EmptyState.vue
    └── router/
        └── guards.ts                 # requireAuth, requireRole
Regras de ouro do FSD (vale colar no README)

Importação só vai de cima pra baixo: app → processes → pages → widgets → features → entities → shared. Nunca o contrário, nunca lateral entre features.
Slice é isolado: features/customer/create-customer não importa de features/customer/edit-customer. Se precisarem compartilhar algo, sobe pra entities/customer ou shared.
Cada slice tem segmentos padronizados: ui/, model/, api/, lib/. Não mistura.
Public API por slice: cada slice exporta via index.ts só o que outras camadas podem usar. Resto é privado.
shared não tem regra de negócio, só utilitários e UI burra.

Por que cada camada existe no seu ERP

shared: cliente axios, formatadores de CPF/data, componentes UI base, tipos compartilhados.
entities: o que é "um Customer", "um User", "uma Session". Inclui API calls básicas (CRUD puro) e tipos. Não inclui formulários ou fluxos.
features: ações do usuário. "Criar customer" é uma feature. "Editar customer" é outra feature. Cada uma com sua UI + lógica.
widgets: composições. CustomerTable é um widget porque junta a feature de listar + as ações de editar/deletar numa UI coesa.
pages: roteáveis. Quase nunca têm lógica — só compõem widgets/features.
processes: fluxos cross-page. No seu MVP, só o boot-auth (tenta refresh ao carregar a app).
app: setup global (providers, router, styles, main.ts).

Entregável: npm run dev abre tela em branco, estrutura de pastas toda criada, providers (router/pinia/query/primevue) instalados.

Fase 6 — Auth no front (entities/session + features/auth + processes/auth)
Ordem de implementação:

shared/api/client.ts: axios com baseURL, withCredentials: true, interceptor de request que injeta Authorization do store de sessão, interceptor de response que em 401 chama /auth/refresh uma vez e refaz a request original. Usa flag pra evitar loop.
entities/session/model/session.store.ts (Pinia): state { accessToken: string | null, user: User | null }, actions setAccess, setUser, clear, getter isAuthenticated.
entities/session/api/session.api.ts: funções login(payload), refresh(), logout(), me(). São só chamadas HTTP, sem lógica de estado.
features/auth/login/model/use-login.ts: composable que usa Vue Query mutation, chama session.api.login, popula o store, redireciona.
features/auth/login/ui/LoginForm.vue: form com vee-validate + zod, três campos (tenant_slug, email, password), chama use-login.
pages/login/ui/LoginPage.vue: layout centralizado + <LoginForm />.
processes/auth/boot-auth.ts: chamado no main.ts antes de montar o app. Tenta refresh(). Se ok, popula store. Se falhar, segue (router guard cuida do resto).
shared/router/guards.ts: requireAuth redireciona pra /login se store vazio. Aplicado em todas as rotas privadas.
widgets/app-shell: sidebar com links (Dashboard, Clientes), topbar com nome do user + logout.

Entregável: app carrega, redireciona pra login, loga, cai no shell autenticado, F5 mantém logado (graças ao refresh no boot).

Fase 7 — Customers ponta a ponta
Backend: aplica o template do módulo (já detalhado). Resultado é os 5 endpoints REST com tenant_scope e audit log.
Frontend, seguindo FSD à risca:

entities/customer/model/types.ts: interface Customer { id, name, document?, email?, phone?, active, createdAt }.
entities/customer/model/schemas.ts: zod schemas pra create/update (compartilhados entre features).
entities/customer/api/customer.api.ts: list(params), get(id), create(data), update(id, data), remove(id).
features/customer/list-customers/model/use-list-customers.ts: Vue Query useQuery com chave ['customers', { page, search }].
features/customer/create-customer/: dialog + composable com mutation; on success faz queryClient.invalidateQueries(['customers']).
features/customer/edit-customer/ e delete-customer/: mesmo padrão.
widgets/customer-table/ui/CustomerTable.vue: usa PrimeVue DataTable, recebe dados de use-list-customers, dispara ações de edit/delete.
pages/customers/ui/CustomersPage.vue: header + search bar + <CustomerTable /> + <CreateCustomerDialog />.

Entregável: MVP 1 funcional. CRUD completo, multi-tenant isolado, com search + paginação.

Fase 8 — Testes
Backend (go test)
Três níveis, cada um com propósito:
1. Testes de unidade — só pra lógica pura (services com regras de negócio). Mocka repository via interface. Roda em milissegundos.
Estrutura: cada service/customer_service.go tem um service/customer_service_test.go ao lado. Usa testify/require.
2. Testes de integração — repository + DB real via testcontainers-go. Sobe Postgres em container temporário, aplica migrations, roda os testes, derruba. É lento (5-10s pra subir), então roda só em CI ou sob demanda.
Pasta internal/modules/customers/db/repo_test.go. Usa build tag //go:build integration pra separar de unit tests.
3. Testes de handler (HTTP) — testa o controller subindo httptest.NewServer com o módulo montado. Verifica status code, formato de response, regras de auth/tenant.
Aqui você pega bugs de middleware, validação, formato de erro. Vale ter pelo menos:

Login com credenciais erradas → 401 com formato certo
Acesso sem token → 401
Acesso a customer de outro tenant → 404 (não 403)
Validação falha → 422 com details populado

Helpers de teste em internal/testutil/: NewTestDB(), NewTestApp(), LoginAs(t, email).
Comando: make test roda unit. make test-integration roda tudo. Cobertura mínima sugerida: 70% nos services, 50% no resto. Mas foca em testar caminhos críticos, não em perseguir número.
Frontend
1. Unit (Vitest): composables (use-login, use-list-customers), helpers de format (CPF, data), schemas zod. Rápido.
2. Component (Vue Test Utils + Vitest): componentes de feature com lógica relevante. Mocka API com vi.mock ou MSW.
3. E2E (Playwright): dois ou três cenários críticos pro MVP:

Login + logout
Criar customer + ver na lista
Tentar acessar /customers sem login → redireciona

E2E roda contra ambiente real (back + front + DB de teste). Em CI, sobe via docker compose.

Fase 9 — Observabilidade e produção-readiness
Antes de pensar em deploy, garante essas peças:

Logs estruturados: já vem do slog. Garante que toda request loga request_id, tenant_id, user_id, method, path, status, duration_ms.
Health checks: /api/v1/health (vivo) e /api/v1/ready (DB + Redis ok). Plataformas de deploy usam isso.
Graceful shutdown: signal.NotifyContext no main, server.Shutdown(ctx) com timeout de 30s.
Timeouts: ReadTimeout, WriteTimeout, IdleTimeout no http.Server. Sem isso uma request lenta trava o processo.
Rate limit no /auth/login: limite por IP no Redis (ex: 5 tentativas / 15min). Evita brute force.
CORS configurável por env (ALLOWED_ORIGINS).
Erros não vazam internals: 500 retorna {code: "INTERNAL_ERROR"} sem stack trace pro cliente. Stack trace só no log.
Secrets via env, nunca commitados. .env.example no repo, .env no gitignore.


Fase 10 — Deploy
Build
Backend Dockerfile multi-stage:

Stage 1: golang:1.23-alpine, baixa deps, compila estático (CGO_ENABLED=0).
Stage 2: gcr.io/distroless/static-debian12, copia binário. Imagem final ~20MB.

Frontend Dockerfile multi-stage:

Stage 1: node:20-alpine, npm ci, npm run build. Sai com dist/ estático.
Stage 2: serve via Caddy ou via reverse proxy direto (não precisa Nginx separado se Caddy faz tudo).

Topologia de produção (sugestão simples)
[Internet]
    ↓
[Caddy]  ← TLS automático via Let's Encrypt
    ├── /api/* → backend (Go) :8080
    └── /*     → frontend estático (servido pelo próprio Caddy)
    ↓
[Postgres 16]  ← managed (Neon, Supabase, RDS) ou container
[Redis 7]      ← managed (Upstash) ou container
Caddyfile fica trivialmente curto (10 linhas). Caddy gerencia TLS sozinho — economiza muita dor de cabeça.
Onde hospedar
Três opções pra MVP, em ordem de simplicidade:

Fly.io: deploy direto do Dockerfile, Postgres gerenciado, regiões próximas do BR (GRU). Bom custo-benefício pra começar. fly deploy e tá no ar.
Railway: ainda mais simples, UI bonita, mas mais caro a longo prazo.
VPS (Hetzner / DigitalOcean) com Docker Compose: ~$5-10/mês, controle total, mais trabalho. Vale quando o produto se prova.

Pra MVP eu iria de Fly.io com Postgres gerenciado externo (Neon tem free tier generoso) e Redis gerenciado (Upstash free tier).
Migrations em produção
Não roda migration no boot do app — corrida entre instâncias. Padrão é:

Job separado (atlas migrate apply) que roda no deploy, antes de subir as novas instâncias.
Em Fly.io: release_command no fly.toml.

CI/CD (GitHub Actions)
Pipeline mínimo:
.github/workflows/backend.yml:

go vet, golangci-lint
go test ./... (unit)
go test -tags=integration ./... (integração, sobe testcontainers)
Build da imagem Docker
Em push na main: deploy (flyctl deploy)

.github/workflows/frontend.yml:

npm ci
npm run lint
npm run test:unit
npm run build
Em push na main: build da imagem + deploy

E2E (Playwright) roda em workflow separado, em ambiente de staging — não bloqueia merge.
Backup
Postgres gerenciado já faz backup. Se for self-hosted, pg_dump diário pra storage externo (S3, B2, R2). Testa restore pelo menos uma vez antes de precisar — backup não testado não é backup.

Definição final de "MVP 1 pronto pra produção"
Checklist completo:
Funcional

 Multi-tenant isolado (validado com 2+ tenants no mesmo ambiente)
 Login + refresh + logout funcionando
 CRUD completo de customers com paginação e busca
 Audit log preenchido em create/update/delete
 Acesso cross-tenant retorna 404

Qualidade

 Testes unit de services passando
 Testes de integração de repositories passando
 Pelo menos 1 fluxo E2E (login + criar customer) passando
 golangci-lint e eslint sem warnings

Operacional

 Logs estruturados com request_id e tenant_id
 Health + ready endpoints
 Graceful shutdown
 Rate limit no login
 CORS configurado por env
 Migrations versionadas no Atlas

Deploy

 Dockerfiles multi-stage funcionando
 CI rodando lint + testes
 CD deploya em ambiente real (Fly.io ou equivalente)
 HTTPS funcionando
 Backup do DB testado (restore validado)
 Variáveis sensíveis em secrets, não em código

Documentação

 README com "como rodar local" em 5 passos
 Doc curto explicando FSD pro time de front
 Doc curto explicando template de módulo pro time de back
 .env.example com todas as variáveis


A partir do MVP 1, adicionar um novo módulo de negócio (produtos, fornecedores, vendas, etc.) vira um processo de ~1-2 dias por módulo:

Edita schema.sql com a tabela
atlas migrate diff → atlas migrate apply
Cria pasta internal/modules/{nome} copiando customers como template
Escreve queries no queries.sql, roda sqlc generate
Implementa service + controller
Registra no modules/register.go
No front, cria entities/{nome}, features/{nome}/*, widgets/{nome}-table, pages/{nome}
Adiciona rota e item no sidebar

O investimento das fases 0-9 é justamente pra esse multiplicador funcionar.
Posso aprofundar em alguma parte específica: o setup concreto do sqlc + Atlas com arquivos reais, o template de um módulo com código de exemplo, ou o setup de testes com testcontainers. Qual te ajuda mais agora?