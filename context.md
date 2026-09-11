> **Note:** This file is a Portuguese brainstorming checklist from early project planning. **Canonical documentation** lives in [`docs/README.md`](docs/README.md) and [`README.md`](README.md). Portfolio 8.5+ scope (RBAC, reporting, users, audit) is implemented — see `docs/PRODUCT.md` and `docs/ROLES.md`.

Contexto técnico que costuma faltar
1. Glossário do domínio
Antes de qualquer coisa: um arquivo GLOSSARY.md com os termos do seu ERP. O que é "cliente" pra você? Pessoa física, jurídica, ambos? "Pedido" é orçamento, venda confirmada, ou os dois com status diferente? IA inventa semântica quando você não fixa. 1-2 páginas resolvem.
2. Convenções de código explícitas
Um CONVENTIONS.md curto com:

Nomenclatura: snake_case em DB, camelCase em Go structs (via tags), PascalCase em tipos TS, kebab-case em arquivos Vue.
Estrutura de erros: o formato que definimos ({error: {code, message, details}})
Padrão de resposta paginada
Como nomear queries sqlc (ListX, GetX, CreateX, UpdateX, DeleteX ou SoftDeleteX)
Como nomear features FSD (create-customer, não customer-create nem new-customer)
Idioma: comentários, commits, mensagens de erro pro usuário — pt-BR ou en? Mistura é o pior cenário.

3. Exemplo completo de um módulo
Antes da IA fazer o segundo módulo, ela precisa ver o primeiro inteiro, perfeito. Tipo: customers ponta a ponta com toda a estrutura (queries.sql, dto, model, service, controller, module.go, testes, no front: entity, features, widget, page). Esse vira o template de referência que você manda junto em todo prompt: "siga exatamente o padrão do módulo customers".
4. Schema do banco como referência única
Mantém o schema.sql sempre atualizado e manda como contexto. IA acerta MUITO mais com schema visível do que tendo que inferir das migrations.
5. Tipos compartilhados back/front
Define onde fica a "fonte da verdade" dos tipos. Sugestão: tipos do back (Go) são gerados pelo sqlc + DTOs escritos à mão. Tipos do front (TS) são escritos manualmente em entities/*/model/types.ts espelhando os DTOs. Se quiser automatizar: swaggo gera OpenAPI do Go, e openapi-typescript gera tipos TS. Mas isso adiciona complexidade — pra MVP, espelhar manual é ok desde que esteja documentado como regra.

Contexto de produto que costuma faltar
6. Personas e permissões
Quem usa o sistema? Pra ERP típico:

Admin do tenant: tudo
Gerente: vê tudo, cria/edita, não deleta nem mexe em config
Operador: opera no dia-a-dia, sem ver financeiro
Visualizador: só leitura

Define os papéis agora, mesmo que o MVP 1 só implemente admin. As próximas features vão precisar dessa matriz.
7. Estado das entidades
Listar os status que cada entidade pode ter. Exemplos: customer (ativo, inativo), invoice (rascunho, emitida, paga, cancelada), order (orçamento, confirmado, entregue, cancelado). Inclui as transições válidas ("de rascunho só pode ir pra emitida ou cancelada"). IA inventa máquinas de estado erradas se você não fixar.
8. Regras de negócio críticas
Lista curta de regras invariáveis: "documento (CPF/CNPJ) é único por tenant", "email do usuário é único globalmente", "não dá pra deletar customer que tem pedido vinculado", "valores monetários sempre 2 casas decimais com decimal.Decimal", etc.
9. Casos de borda já decididos
Coisas que vão aparecer cedo:

Customer sem documento (pessoa física estrangeira, por ex)
Operações em lote (deletar 50 customers de uma vez — permite?)
Importação de planilha (vai ter? formato?)
Exportação (CSV, Excel, PDF — quais?)
Soft delete vs hard delete — fixou soft, mas e pra GDPR/LGPD ("apagar meus dados")?


Contexto de UI/UX que costuma faltar
10. Sistema de design mínimo
Quando você passar o design em HTML, inclui o design system:

Paleta (primary, secondary, success, warning, danger, neutrals)
Tipografia (família, tamanhos, pesos)
Espaçamentos (escala — múltiplos de 4 ou 8?)
Border radius padrão
Shadows
Estados (hover, focus, disabled, loading)

Sem isso, cada tela vai sair com cor um pouco diferente. Tailwind config + variáveis CSS no app.scss resolvem — define uma vez, IA segue.
11. Componentes-base do PrimeVue que você vai usar
Lista quais componentes do PrimeVue você adotou pra cada caso:

Tabela → DataTable
Form input → InputText
Select → Dropdown
Data → Calendar
Modal → Dialog
Notificação → Toast
Loading → ProgressSpinner

Sem isso, IA mistura PrimeVue com componentes custom ou outras libs.
12. Padrões de UX recorrentes
Define uma vez como toda tela se comporta:

Loading state: skeleton ou spinner?
Empty state: ilustração + texto + CTA?
Erro de form: inline abaixo do campo ou toast?
Confirmação de delete: dialog modal ou inline?
Após criar/editar com sucesso: toast + fecha modal + atualiza lista
Paginação: server-side com offset, 20 itens por página padrão
Busca: debounce de 300ms, com loading inline

Isso vira um arquivo UX_PATTERNS.md. IA segue ao pé da letra.
13. Layout-padrão de página de CRUD
Como vai ser a estrutura visual de toda página de listagem? Exemplo:
[Header: título + descrição]              [Botão "Novo"]
[Search bar] [Filtros] [Botão "Exportar"]
[Tabela com paginação]
Define uma vez, IA replica.

Contexto pra IA trabalhar bem com você
14. README pro contexto da IA
Um arquivo AI_CONTEXT.md (ou .cursorrules / CLAUDE.md se usar Cursor/Claude Code) com:

Stack consolidada
Estrutura de pastas (back + front)
Como rodar local
Comandos do Makefile
Padrão de cada módulo (referenciando customers)
Convenções de código
"Sempre que criar um módulo novo, siga o template de customers"
"Nunca crie tabela sem tenant_id"
"Nunca acesse o DB sem filtrar por tenant_id"
"Toda query deve receber ctx context.Context como primeiro parâmetro"

Isso reduz repetição entre prompts e evita IA esquecer regras críticas.
15. Prompt template
Cria um template de prompt pra usar quando pedir um novo módulo:
Contexto: [link/cole AI_CONTEXT.md]
Módulo de referência: [link/cole estrutura completa de customers]
Schema atual: [cole schema.sql]
Novo módulo: products
Especificação:
  - Campos: nome, sku, preço, estoque, ativo
  - Regras: sku único por tenant, preço > 0, estoque >= 0
  - Endpoints: CRUD padrão + GET /products/low-stock
  - UI: tela de listagem + modal de criar/editar + botão "ajustar estoque"
Design: [cole HTML/PNG]
Tarefas:
  1. Schema + migration
  2. Queries sqlc
  3. Module Go completo (dto, model, service, controller, module.go, tests)
  4. Frontend FSD (entity, features, widget, page)
  5. Atualizar sidebar com link
Reutilizável pra cada módulo, IA tem tudo que precisa.
16. Definição de "pronto"
Pra cada tarefa que você passar pra IA, define o checklist do que significa pronto:

 Código compila / passa no lint
 Migration roda sem erro
 Tem teste do service
 Tem teste de handler com 401/403/404/200
 Audit log adicionado nos endpoints de write
 tenant_id filtra em todas as queries
 Front compila e tela funciona
 Adicionou item no sidebar

IA frequentemente esquece um desses. Lista explícita força ela a verificar.

O que não vale pré-decidir
Pra evitar over-engineering:

Internacionalização (i18n): se MVP é só pt-BR, não bota infra de i18n agora. Adiciona quando precisar.
Dark mode: só se for requisito real.
PWA / offline: ERP raramente precisa. Não invista.
Microserviços: monolito modular como estamos fazendo é o caminho. Quebra depois se precisar.
GraphQL: REST tá ótimo pra ERP. Não muda agora.
WebSocket / real-time: só se tiver feature específica que pede (chat, dashboard ao vivo).


Sugestão de organização do pacote de contexto
Cria uma pasta docs/ no monorepo:
docs/
├── AI_CONTEXT.md          # contexto principal pra IA
├── ARCHITECTURE.md        # arquitetura back + front
├── CONVENTIONS.md         # naming, lint, padrões
├── GLOSSARY.md            # termos de negócio
├── BUSINESS_RULES.md      # regras invariáveis
├── ROLES.md               # personas e permissões
├── UX_PATTERNS.md         # padrões de UI/UX
├── DESIGN_SYSTEM.md       # cores, tipografia, espaçamentos
├── MODULE_TEMPLATE.md     # como criar novo módulo (back+front)
├── PROMPT_TEMPLATES.md    # templates de prompt pra IA
└── designs/
    ├── customers/
    │   ├── list.html
    │   ├── create.html
    │   └── README.md       # spec da tela
    └── ...
Cada doc curto (1-3 páginas). Total ~30-40 páginas de contexto bem organizado. Vira a "memória externa" que a IA precisa pra trabalhar consistente.

Resumo: o que adicionar antes de começar
Em ordem de impacto:

Módulo customers feito 100% certo = template visual pra todos os próximos. Investe tempo aqui.
GLOSSARY + BUSINESS_RULES + ROLES = define o produto sem ambiguidade.
UX_PATTERNS + DESIGN_SYSTEM = define o visual sem ambiguidade.
AI_CONTEXT + MODULE_TEMPLATE + PROMPT_TEMPLATES = reduz repetição e erro entre prompts.
CONVENTIONS = padroniza código.

Com isso + os designs HTML/MD que você vai criar, IA produz com consistência altíssima. Sem isso, cada módulo sai um pouco diferente do outro, e refatoração depois é cara.
Quer que eu monte o esqueleto de algum desses arquivos (por exemplo, o AI_CONTEXT.md ou o MODULE_TEMPLATE.md) com base em tudo que conversamos?