# AskMe modernization plan

## Executive summary

AskMe is a small server-rendered Go Q&A application with a useful domain core, but its
runtime, persistence, authentication, asset pipeline, deployment story, and interface
reflect a 2018-era prototype. The safest path is **an incremental renovation of the
existing modular monolith**, not a frontend or microservice rewrite.

The recommended destination is:

- a currently supported Go release, standard-library HTTP primitives, `html/template`,
  and `go:embed`;
- PostgreSQL accessed through `pgx` and explicit, versioned SQL migrations;
- first-party email/password authentication with hardened, opaque, server-side sessions;
- a responsive, accessible, server-rendered design system with small progressive
  enhancements only where they improve a measured workflow;
- one reproducible local/CI/production path with health checks, structured logs,
  graceful shutdown, and automated dependency/security maintenance.

The initial product scope is now settled: profiles and answered questions are public;
profile creation is invitation-only; profile owners may opt into anonymous questions;
the application owner manages invitations; and following and likes are part of the first
modernized product release. There is no production data or legacy account migration to
preserve.

Do not combine all of this into one large rewrite. Establish a green baseline, remove
Okta behind characterization tests, modernize infrastructure, then ship the new UX one
vertical slice at a time.

## What the application does today

The product is an authenticated ask-and-answer network:

1. A user signs up or signs in.
2. They visit another user's email-addressed profile and send a question.
3. The recipient sees unanswered questions and can answer them.
4. Answers are displayed in a public-looking feed, although every current product route
   is actually authentication-gated.

The code is organized as a modest clean-architecture-style monolith:

- `models/` contains users, questions, answers, and password hashing.
- `user/usecase/` coordinates the product workflows.
- `user/`, `question/`, and `answer/` define repository boundaries with database and
  in-memory implementations.
- `web/framework/` is a custom router/render/session abstraction.
- `web/askme/` contains controllers, templates, and vendored browser assets.
- `main.go` wires everything with `dig`, loads environment configuration, migrates the
  database, and starts the HTTP server.

This separation is worth preserving while simplifying its implementation.

## Baseline assessment

### Runtime and dependencies

- `go.mod` targets Go 1.16 and pins most libraries to 2018-2021 versions.
- GORM v1 (`github.com/jinzhu/gorm`) is used, and each repository operation opens and
  closes a new database handle instead of sharing one process-wide pool.
- Packr generates/embeds assets through an obsolete build-time tool even though modern
  Go has native embedding.
- `dig` provides dependency injection for a graph small enough to wire explicitly.
- The router is a third-party package wrapped in a home-grown framework. Modern Go's
  standard HTTP stack is sufficient for this route set.
- Bootstrap 4.1, jQuery 3.3.1, and Popper are checked into the repository as minified
  assets. There is no frontend dependency inventory or update mechanism.

### Authentication and security

- Okta is the active controller. It owns `/login`, the callback, and `/logout`; the
  already-present local `AuthController` is not registered.
- Okta-specific configuration, social identity-provider handling, token exchange,
  verification, templates, utilities, and dependencies are spread across bootstrap,
  framework configuration, controllers, and UI.
- The Okta callback creates local users with the literal password `defaultPassword`.
- Several Okta network/parse errors are ignored, outbound HTTP has no timeout, callback
  URLs assume plain HTTP, and token material is stored in the session.
- Sessions are held only in process memory. Restarts log everyone out, horizontal
  scaling is unsafe, session IDs are not rotated on login, and cookies do not explicitly
  set `SameSite`.
- Logout is a state-changing GET. CSRF checks cover local login/signup only; question
  and answer POST routes do not apply `RequireCSRF`.
- Authorization is incomplete: handlers authenticate the caller but do not consistently
  prove that the caller is allowed to answer the selected question.
- Form inputs have no central normalization, length limits, or validation contract.
  Error paths sometimes print internal errors to the response.

### Data and correctness

- Schema changes happen via `AutoMigrate` on every startup; there is no auditable or
  reversible migration history.
- Migration opens a second database connection and can call `Close` before checking
  whether opening the database succeeded. Repository methods have the same deferred
  close-before-error-check pattern.
- Answer creation and updating the question are separate operations rather than one
  transaction, allowing partially completed answers.
- `LoadAnswers(userID)` ignores `userID`, so every profile can receive the same global
  answer collection.
- Like/unlike methods are placeholders and always report zero.
- Controllers often continue after redirects or lookup/parse failures, which can lead
  to nil dereferences or unintended writes.
- Email addresses are used directly as route identifiers. That exposes personal data,
  makes later email changes break URLs, and conflates identity with presentation.

### Operations and delivery

- The Docker image builds and runs as root on an old full Go image, installs Packr with
  `go get`, mutates every copied file with `dos2unix`, and is not multi-stage.
- Compose pins PostgreSQL to the mutable `latest` tag. `docker-local.yml` has no database
  environment configuration of its own, while the documented `.env` workflow is
  fragile.
- Startup uses shell wait scripts instead of health checks and application-level retry.
- The HTTP server has no read/write/idle/header timeouts, graceful shutdown, readiness,
  liveness, or request IDs.
- There is no current CI configuration, formatter/linter gate, vulnerability scan,
  integration suite, or automated dependency updates.
- Configuration validation exits piecemeal and has no typed defaults or safe production
  invariants.

### UX and accessibility

- The UI is stock Bootstrap with a fixed navbar/footer, placeholder remote avatars,
  inline styles, and desktop-first columns.
- Navigation includes a search form with no working endpoint. The information hierarchy
  does not clearly separate a user's profile, their public answers, and their private
  question inbox.
- Login/signup and empty/error/success states lack actionable, field-level feedback.
- Templates contain malformed or incorrect markup/data references (for example,
  `div="container"` and `Emai`).
- There is no documented visual language, responsive behavior, keyboard/focus audit,
  reduced-motion treatment, contrast target, or automated accessibility check.

## Confirmed product requirements

Record these decisions in ADRs and translate them into acceptance tests before changing
the relevant workflows:

1. **Public content:** Profiles and answered questions are visible to anonymous visitors.
   A visitor must not need a session to open a profile or its answered-question feed.
   Unanswered questions, account settings, invitation administration, and moderation
   controls remain private.
2. **Configurable anonymous questions:** Each profile owner can allow or disallow
   questions from unauthenticated visitors. The default should be **disabled** so a new
   profile opts in deliberately. Signed-in questions retain their author internally;
   anonymous questions have no user identity to expose. The public answer presentation
   must respect the asker's anonymity choice.
3. **Abuse resistance:** Anonymous ask endpoints need stricter body/length limits,
   per-IP and per-profile rate limits, monitoring, and a report/block workflow. Add a
   privacy-conscious CAPTCHA or equivalent challenge as a configurable escalation when
   abuse signals or thresholds justify it, rather than burdening every visitor by
   default. Document behavior when a limiter or CAPTCHA provider is unavailable.
4. **Invitation-only profiles:** There is no open registration. The application owner
   needs an authenticated administration screen to issue and revoke invitation links.
   Invitations are random, single-use, hashed at rest, expire, and may be bound to an
   intended email address. Redeeming one creates a profile and verifies control of that
   email. Audit issuance, revocation, redemption, and expiry without logging raw tokens.
5. **Initial administration:** Bootstrap the first owner/admin through an explicit,
   one-time CLI command or deployment secret—not a public route or a hard-coded account.
   Enforce an `admin` role server-side for every invitation-management action.
6. **Social scope:** Following and likes are P1 requirements. Search, notifications, and
   richer discovery can follow after the first modernization release. Moderation and
   blocking are introduced as required to operate anonymous questions safely.
7. **No legacy migration:** The application was never deployed beyond local testing, so
   no Okta users or production data need preservation. Development databases can be
   rebuilt from migrations and deterministic seed data; do not build dual-authentication
   or legacy-schema compatibility machinery.

## Target architecture

### Keep

- One deployable Go service and one PostgreSQL database.
- Server-side HTML for fast first paint, simple operations, and durable URLs.
- Domain/use-case/repository boundaries, but only where they improve testing or isolate
  infrastructure.
- `bcrypt` initially so existing password hashes remain valid. Revisit password hashing
  parameters through a separate, versioned migration strategy.

### Replace or simplify

| Current | Target | Reason |
| --- | --- | --- |
| Go 1.16 | A currently supported Go release, pinned in CI and containers | Security support, toolchain improvements, native routing/embedding |
| Packr | `//go:embed` | Removes code generation and a build-time dependency |
| GORM v1 | `pgx` plus explicit SQL (optionally generated with `sqlc` after a spike) | Clear queries/transactions, shared pooling, fewer ORM surprises |
| Startup `AutoMigrate` | Versioned migrations run as a release step | Reviewable, deterministic schema evolution |
| `httprouter` + custom routing wrapper | `net/http` `ServeMux` and focused middleware | Smaller surface and idiomatic handlers/tests |
| `dig` | Explicit constructor wiring | Compile-time visibility for a small graph |
| In-memory sessions | Opaque cookie ID + PostgreSQL-backed session store | Revocation, restart safety, and multi-instance support |
| Okta OIDC | Local auth controller | Meets the stated product requirement and removes external identity coupling |
| Vendored Bootstrap/jQuery/Popper | Purpose-built CSS design system; minimal ES modules only if needed | Cohesive UX, less client code, maintainable assets |

Use feature-oriented packages rather than recreating a generic framework. A reasonable
shape is `cmd/askme`, `internal/auth`, `internal/questions`, `internal/profiles`,
`internal/platform/{config,database,httpserver,session}`, `web/templates`, and
`web/static`. Migrate toward it as files are touched; do not perform a package-only
rewrite.

## Delivery roadmap

Effort ranges are engineering estimates for one experienced developer and exclude
product/legal review. Re-estimate after Phase 0.

### Phase 0 — establish safety and decisions (3-5 days)

**Goal:** make current behavior observable and create a repeatable baseline before
changing dependencies.

- Convert the confirmed requirements above into ADRs and executable acceptance criteria.
- Add `.env.example` with safe local defaults and one typed configuration loader that
  returns errors rather than terminating deep in constructors.
- Add CI for formatting, `go test -race ./...`, `go vet ./...`, build, vulnerability
  scanning, and a PostgreSQL-backed smoke test.
- Add characterization HTTP tests for signup, login, logout, profile/feed, ask, answer,
  access control, and CSRF behavior.
- Add a deterministic development seed command with at least three personas and enough
  answered/unanswered content to exercise the UI.
- Capture a small migration dataset and backup/restore rehearsal.
- Add ADRs for local authentication, invitations/admin bootstrap, anonymous-question
  abuse controls, persistence, routing, and the UI approach.

**Exit criteria:** CI is green; a new contributor can start the app from documented
commands; critical current workflows have black-box tests; requirements and security
boundaries are recorded.

### Phase 1 — remove Okta and harden local authentication (4-7 days)

**Goal:** remove the unwanted provider without leaving an authentication/security gap.

- Register and test `AuthController` first; improve it to return generic authentication
  failures and field validation rather than raw errors.
- Remove public signup and add invitation redemption. Build a server-authorized admin
  screen to issue/revoke invitations plus a one-time CLI path that bootstraps the initial
  owner. Store only invitation-token hashes and enforce expiry, one-time use, optional
  intended-email binding, and audit events.
- Add password reset and email verification using hashed, expiring, single-use database
  tokens. Rate-limit invitation redemption, login, and reset by account and IP.
- Rotate the session ID after authentication; store only the local user ID; set
  `HttpOnly`, `Secure` in production, `SameSite=Lax`, bounded lifetime, and an explicit
  inactivity policy.
- Change logout to CSRF-protected POST. Apply CSRF protection to every state-changing
  browser route and add authorization tests for answering questions.
- Once local auth and invitation redemption are proven, delete the Okta controller,
  callback route/template/utilities/configuration, token/session fields, social IdP
  options, verifier/JWX dependencies, and README instructions.
- Delete obsolete identity/token-related development data and recreate the local database
  from versioned migrations and seed data; no legacy compatibility mode is needed.

**Exit criteria:** no `OKTA_` configuration or Okta dependency remains; all auth flows
work without a third party; only an admin can manage invitations; open registration is
impossible; and session fixation, CSRF, brute-force, and object-level authorization tests
pass.

### Phase 2 — modernize the platform in controlled steps (1-2 weeks)

**Goal:** reach a supported, reproducible foundation while preserving behavior.

1. Upgrade Go one supported step at a time, running tests/race/vet at every step. Pin
   the final version consistently in `go.mod`, CI, and Docker.
2. Replace Packr with `go:embed`, then remove Packr and generated build tooling.
3. Introduce a process-wide `pgxpool.Pool`; port one repository/query at a time behind
   existing interfaces. Use contexts and deadlines throughout.
4. Introduce numbered migrations with constraints and indexes. Baseline existing
   databases before disabling `AutoMigrate`. Include unique normalized email, foreign
   keys, answer ownership, and query indexes.
5. Make answering a question a transaction with an idempotency/concurrency rule. Fix
   profile feed filtering and either implement likes fully or remove them from contracts.
6. Replace `dig` with explicit wiring, then replace the custom router/framework with
   standard handlers and narrowly scoped middleware. Preserve URLs through contract
   tests and redirects where slugs change.
7. Move from email URLs to stable public handles/slugs. Because there is no production
   data, seeded development profiles can be recreated without a compatibility period;
   keep immutable numeric/UUID IDs internally.
8. Model anonymous questions explicitly: allow a nullable asker, capture the profile
   owner's opt-in at write time, and avoid storing unnecessary visitor identity. Add
   indexes/storage needed by the selected per-profile/IP rate limiter.

**Exit criteria:** no GORM v1, Packr, `dig`, or third-party router dependency remains;
schema changes are migration-driven; all repositories share a bounded pool; the main
journeys and old URL redirects pass integration tests.

### Phase 3 — redesign UX as vertical slices (1-2 weeks)

**Goal:** deliver a coherent, accessible experience without a high-cost SPA rewrite.

Start with low-fidelity flows and a content model, then implement these slices in order:

1. **App shell and auth:** responsive header, clear signed-in navigation, polished
   invitation redemption/login/reset forms, inline validation, and visible focus/skip
   links.
2. **Public profile/feed:** profile identity and shareable ask CTA; readable question /
   answer cards; timestamps; intentional empty, loading, and error states.
3. **Ask flow:** focused composer, recipient context, signed-in or explicitly anonymous
   state, owner opt-in enforcement, character count/limits, confirmation, rate-limit
   feedback, and abuse/reporting affordances. Introduce a CAPTCHA challenge only through
   the defined risk-based escalation path.
4. **Question inbox and answer flow:** private ownership cues, status filters, safe answer
   composer, confirmation, and optimistic enhancement only if it materially helps.
5. **Account and safety:** profile editing, password/account controls, block/report tools,
   and privacy controls required by the chosen scope.
6. **P1 social workflows:** following/unfollowing and liking/unliking answers, including
   accessible counts and controls, authorization, duplicate-click idempotency, and
   privacy/empty states. Do not display controls before their backend behavior exists.
7. **Invitation administration:** an admin-only list of active/redeemed/revoked/expired
   invitations with issue, copy-once, and revoke workflows; raw tokens are shown only at
   creation.

Create CSS tokens for color, type, spacing, radius, elevation, and motion. Build reusable
template components for buttons, fields, alerts, avatars, cards, empty states, and
pagination. Use a system font stack and locally served assets; do not depend on remote
placeholder images. Target WCAG 2.2 AA, keyboard-only use, 200% zoom, reduced motion,
touch targets, useful document titles, semantic landmarks, labels, and announced errors.

Test at narrow mobile, wide mobile, tablet, and desktop breakpoints. Add template
render tests, automated accessibility smoke checks, and screenshot regression coverage
for the primary pages. Validate the redesign with five representative users before
fine polish.

**Exit criteria:** all primary workflows work at 320px width and by keyboard; automated
accessibility checks have no serious violations; following and likes work end to end;
dead search UI is removed; invitation administration is keyboard-accessible; and user
validation findings are resolved or explicitly backlogged.

### Phase 4 — production readiness and launch (4-7 days)

**Goal:** make releases safe and diagnosable.

- Use a multi-stage, non-root, minimal Docker image with a pinned base by digest; add
  health checks and a `.dockerignore`.
- Pin PostgreSQL to a supported major version and document upgrade, backup, restore,
  and rollback procedures. Test migrations against a production-like copy.
- Configure `http.Server` timeouts, request body limits, graceful shutdown, readiness
  and liveness endpoints, and startup dependency retry with a bounded deadline.
- Emit structured logs with request/correlation IDs and redaction. Add metrics for
  latency/error rate, DB pool health, auth failures/rate limits, asks, and answers.
- Define service-level indicators and alerts; connect exception reporting without
  sending passwords, tokens, question content, or other unnecessary personal data.
- Add preview/staging deployment, migration-before-app release sequencing, smoke tests,
  rollback automation, dependency update automation, SBOM/image scanning, and secret
  scanning.
- Perform threat modeling and a focused security review before enabling anonymous
  questions in production.

**Exit criteria:** a staging release and rollback are rehearsed; backup restoration is
verified; dashboards/alerts observe an end-to-end smoke journey; the launch checklist
has named owners.

## Prioritized backlog

### P0 — before exposing the current app

- Remove the default-password Okta provisioning path.
- Remove open signup; implement secure invitation redemption and a one-time owner/admin
  bootstrap command before allowing profile creation.
- Enforce authorization on answer creation and CSRF on all writes.
- Split public profile/answer reads from private inbox/settings/admin routes and test the
  boundary without relying on template visibility.
- Fix nil/error handling around database connections, redirects, lookups, and Okta HTTP.
- Add HTTP timeouts/body limits and prevent raw internal errors from reaching users.
- Fix the unfiltered answer feed and make answer persistence transactional.
- Establish supported runtime/container/database versions and automated tests/scanning.

### P1 — modernization release

- Complete Okta removal, invitation administration, and local account
  recovery/verification.
- Add migrations, shared connection pooling, stable handles, persistent sessions, and
  rate limiting.
- Replace obsolete framework/build dependencies and deliver the core responsive redesign.
- Deliver configurable anonymous questions with per-profile/IP rate limiting,
  observable CAPTCHA escalation, reporting, and blocking.
- Implement following and likes with database constraints, idempotent writes, counts,
  authorization, and end-to-end tests.

### P2 — after reliable usage data

- Evaluate notifications and search; implement them only with validated workflows.
- Add pagination/cursoring, profile discovery, richer moderation, and performance work
  only in response to observed needs.
- Consider HTMX or small JavaScript modules for proven interaction pain points. A SPA,
  separate API, cache, queue, or microservices should require measured justification.

## Testing strategy

Use a test pyramid tied to product risks:

- **Unit:** domain validation, password/invitation/session behavior, anonymous-question
  opt-in and rate-limit policies, authorization rules, configuration parsing, and use
  cases with fakes.
- **HTTP/component:** handlers through `httptest`, real templates, cookies, CSRF, status
  codes, redirect safety, accessibility landmarks, and error messages.
- **Integration:** repositories and migrations against the pinned PostgreSQL major;
  transaction/concurrency cases; clean migration from an empty database and deterministic
  seed loading.
- **End-to-end:** owner bootstrap; invitation issue/revoke/redeem/expiry; login/logout /
  reset; public profile access; signed-in and anonymous asks; anonymous opt-out, rate
  limit, and CAPTCHA escalation; answer; follow; like; private-route denial; account
  deletion; and abuse controls in a production-like container stack.
- **Non-functional:** race detector, fuzzing for parsers/validators, vulnerability and
  secret scanning, accessibility, responsive screenshots, load smoke tests, backup /
  restore, migration rollback, and session behavior across two app instances.

Every production defect should first gain the lowest-level useful regression test.
Coverage percentage is a diagnostic, not the goal; auth, authorization, migrations, and
core writes require explicit branch and failure-path coverage.

## Rollout and risk controls

- Ship phases as small pull requests with behavior-focused acceptance criteria.
- Put schema expansion before application use and cleanup in a later release
  (expand/migrate/contract). Never couple an irreversible migration to the first deploy.
- Recreate development databases from versioned migrations and seed data. Do not add a
  compatibility window for accounts or schema because there is no production data.
- Move seeded email profile URLs to handles in one step; compatibility redirects are
  optional for local bookmarks rather than a release constraint.
- Use a feature flag for each redesigned workflow, not for permanent configuration.
- Maintain a rollback runbook and verified backup before auth or schema cutovers.
- Track invitation redemption and successful login rates, signed-in versus anonymous ask
  completion, rate-limit/CAPTCHA activation, answer completion/time-to-answer, follow /
  like engagement, and server error rate. Pair them with accessibility findings, abuse
  reports, and false-positive limiter challenges so growth does not hide product harm.

## Definition of modernized

The project is modernized when it is not merely compiling on a new Go version, but when:

- Okta and its operational assumptions are fully removed;
- profile creation requires a valid invitation, and only the bootstrapped owner/admin can
  issue or revoke invitations;
- dependencies, containers, and PostgreSQL are pinned to supported versions with an
  owned update cadence;
- public profiles/answers and private inbox/settings/admin boundaries are explicit and
  covered by authorization tests;
- local auth, sessions, invitations, CSRF, recovery, configurable anonymous questions,
  rate limiting, CAPTCHA escalation, and abuse controls have automated security tests;
- database changes are versioned, writes are transactional, and backup/restore is proven;
- invitation, profile, ask, answer, follow, and like workflows are responsive, accessible,
  understandable, and validated with users;
- CI produces a reproducible artifact and staging can deploy, observe, and roll it back;
- runbooks, ADRs, threat model, and ownership make the next upgrade routine rather than
  another rescue project.
