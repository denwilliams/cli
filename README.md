# Public CLI tools (Go)

One Go module containing independently runnable CLIs. Each command lives in
`cmd/<name>` with isolated implementation packages under `internal/<name>`.
Module tags will initially version the collection. `recipes-cli` is not imported.

## aussupermarkets — foundation only

Search framework for **Coles, Woolworths and ALDI**, with prices required to
represent **in-store** purchases. **No live providers are configured yet.**
Production runs return separate `unconfigured` groups and exit 1: no synthetic
products, no network requests, and no online-delivery/collection price fallback.
The framework is not yet a usable live price comparison tool.

Requires Go **1.26 or newer**. No third-party Go dependencies.

```sh
go build -o /tmp/aussupermarkets ./cmd/aussupermarkets
/tmp/aussupermarkets --help
/tmp/aussupermarkets --postcode 0800 --sort price --max-results 5 'milk 2L'
/tmp/aussupermarkets --providers coles,woolworths --format json 'bread'
```

These searches intentionally exit 1 until permitted in-store feeds exist. Flags
must precede exactly one query argument; quote multiword queries. Use `--` before
a query starting with `-`. Build locally from this checkout; no release is part
of this milestone. Once tagged, the independent command can be installed using
`go install github.com/denwilliams/cli/cmd/aussupermarkets@<tag>`.

### Flags and runtime context

| Flag | Default | Contract |
| --- | --- | --- |
| `--providers` | `coles,woolworths,aldi` | Nonempty comma-separated subset; duplicates/unknown IDs fail. Output order is always Coles, Woolworths, ALDI. |
| `--max-results` | `10` | Positive limit **per supermarket**, applied after sorting. |
| `--sort` | `price` | `price` or `name`, **ascending** within each retailer group. |
| `--format` | `human` | `human` or `json`. |
| `--postcode` | absent | Exactly four digits if supplied; preserves leading zeros. No guessed location. |
| `--store` | absent | Repeat `retailer=store-id` for distinct enabled retailers. IDs are opaque, feed-defined, not guessed. |
| `--timeout` | `10s` | Positive duration per provider, bounded by parent cancellation. |

All context is runtime-only; no config file, credentials or persistent defaults.
A postcode alone does **not** identify a unique store. Future adapters must
resolve and disclose an actual in-store location (or fail explicitly), with any
source-specific context requirements documented in their own changes. Supplied
store/postcode must match returned context. Missing context cannot be replaced
with a fabricated store or silently changed shopping channel.

### Search/output contract

- Enabled providers start concurrently; disabled providers receive no calls.
- Separate groups have canonical retailer order regardless of completion order.
  Healthy groups survive peer errors/timeouts. Empty matches are distinct from
  `unconfigured`, `unavailable`, `error`, `timeout` and `cancelled`.
- AUD money uses **nonnegative int64 cents** (`minor_units`), never float prices.
  `price: null` means unknown, not zero. Unit prices carry an explicit unit.
- Price sort ranks standard, unconditional `pack` totals first by cents. Weight
  estimates, member/multibuy and unknown prices follow by name/ID, not by
  incomparable amounts. Conditional prices require conditions. Pack sizes and
  unit prices are disclosed, **not quantity-normalized or equated**.
- Name comparison uses Go's Unicode lowercase strings, then original name, then
  unique product ID for deterministic ties; it is not locale-aware collation.
  Name sort is independent of price availability/basis.
- Sort **before** the per-group limit, without mutating a provider-owned slice.
- JSON exposes requested and actual context, source URL/name, optional observation
  and retrieval timestamps, product observation timestamps, coverage,
  candidate count and local limiting. Omitted timestamps mean **unknown**, not
  fresh. Human output prints unknown freshness explicitly and quotes untrusted
  text to keep terminal escapes/newlines from altering output.
- Coverage is `complete_catalogue`, `sample` or `unknown`, with an independent
  provider `truncated` flag (a truncated set cannot be `complete_catalogue`). `limited` reports this CLI's own cap. A sorted sample
  is **not** a guarantee of globally cheapest/name-first products. Future feeds
  need permitted complete enumeration or supported server-side sorting for that
  guarantee; adapters must not relevance-limit before local sorting and claim
  complete coverage.
- All successful provider results (including empty matches) must identify an
  actual store and source name/URL.
  Invalid channel, mismatched context, duplicate/missing product IDs, malformed
  money and unlabelled conditional prices fail that group without discarding
  healthy groups. Error groups contain no products or invented source metadata.

Exit codes: **0** for help or all enabled groups succeeding (including empty
matches); **1** for any failed/unconfigured group or output error; **2** for
invalid arguments. JSON remains on stdout for provider failures; diagnostics for
argument/output errors and help use stderr. SIGINT/SIGTERM cancel the invocation.
Adapters must honour context cancellation and release their resources. The engine
bounds caller waiting even for a faulty slow adapter; it cannot forcibly kill
arbitrary adapter goroutines. No such background work occurs in current production
providers.

### Development and remaining work

```sh
test -z "$(gofmt -l cmd internal)"
go vet ./...
go test ./...
go test -race ./...
go build ./...
go build -o /tmp/aussupermarkets ./cmd/aussupermarkets
/tmp/aussupermarkets --help
```

Tests and their checked-in renderer goldens are synthetic and offline. Production
registrations contain only explicit unconfigured states. CI runs formatting,
vet, tests, race checks, build, help and unconfigured JSON smoke on Go 1.26/1.27.

- Foundation: [issue #1](https://github.com/denwilliams/cli/issues/1).
- Live-provider work is tracked separately in
  [issue #2](https://github.com/denwilliams/cli/issues/2).
- [Source feasibility and gates](docs/live-sources.md) records the access and
  **in-store provenance** blockers. Engineering owns further source discovery;
  users are not being asked to locate feeds. No access-control bypass, ALDI
  scraping, paid-service/credential procurement, merge or deployment is included.
