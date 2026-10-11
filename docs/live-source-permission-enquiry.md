# Draft permission and source-contract enquiry

For [issue #2](https://github.com/denwilliams/cli/issues/2). **Draft only — not sent.**
No implied authorization to make contact, register, enter contracts, obtain
credentials, start trials, spend, or run a live API/smoke. Engineering owns routing
and assessment; the user need not locate feeds. See the
[investigation and routing caveats](live-source-investigation-2026-10-11.md).

## Retailer version (send separately to each authorized contact)

Subject: Non-binding enquiry: permitted store-specific in-store price feed for an independent CLI

We are developing an independent Go command-line grocery-search tool, initially
for Coles, Woolworths and ALDI. The existing foundation intentionally returns
unconfigured providers rather than scraping or substituting online prices. We
would like to establish whether you offer, or expressly authorize a licensee to
offer, a feed suitable for the following limited use.

The tool would take a query and runtime postcode/store, display separate retailer
groups of product names, pack/unit details and exact AUD prices, and sort within
each retailer by price or name before applying a per-retailer result limit. It
requires **actual local in-store prices**, not delivery/collection prices or
inferred national/zone averages. Member/multibuy/weight offers would retain their
conditions and would not rank as unconditional standard pack totals.

Could you route this enquiry to the responsible data-licensing/digital partnership
team, and confirm in writing:

1. The authorized source/feed operator, access mechanism and permission for an
   independent automated CLI. If a licensee supplies the feed, how can upstream
   retailer authorization be independently verified?
2. Rights to use and display names, IDs, pack/unit details and prices in local
   human/JSON CLI output; any caching, redistribution, attribution, retention,
   open-source implementation or geographic restrictions. Please distinguish
   software rights from upstream data/product-metadata rights.
3. Whether prices expressly represent the selected store's **in-store** channel,
   including specials and exclusions. How a postcode maps to store IDs and how
   missing, conflicting or multiple-store contexts are rejected/resolved.
4. Observation/effective time separately from retrieval time; update cadence,
   maximum lag, stale/unavailable handling and whether values are observed,
   inferred or historical. No assumed live status from a fresh retrieval.
5. Keyword matching, supported ascending price/name sort, pagination/end markers,
   result caps, catalogue coverage and truncation. Can a client enumerate all
   matching candidates within permitted limits, or must results be labelled a
   sample rather than globally cheapest/top-N?
6. A licensed representative **offline** schema/example showing stable product
   ID/name, exact AUD decimal/minor-unit encoding, zero versus missing values,
   pack quantity, unit amount/unit, variable-weight estimates, member/multibuy
   conditions, price validity and store/source timestamps. May examples be
   committed as contract-test fixtures, or must tests use synthetic equivalents?
7. Documented rate/volume limits, Retry-After/429 behavior, error formats,
   cancellation expectations and allowed retry/caching policy.
8. Eligibility, credential and onboarding requirements, and any fees/quotas.
   Is there a permitted no-credential/no-cost route? Please do not create an
   account, start a trial, provision credentials or make any charge for this
   enquiry. Any procurement or live test would require separate approval.

A consumer website, private app endpoint or paid scraper listing without these
rights is not sufficient. If no suitable route is available, a written confirmation
will let us leave that provider honestly unconfigured.

## GS1 recipient-eligibility variant

Use the same contract questions, but first ask:

- Is NPC price data supplier/cost pricing or retailer-set consumer shelf pricing?
- Does any authorized dataset cover Coles, Woolworths and/or ALDI at actual-store
  granularity with observation times? Which retailer authorizes redistribution?
- Can an independent consumer CLI qualify as a recipient and display these fields
  outside the supply-chain trading-partner context?
- What permission/fees apply separately to recipient connection and each
  publisher's data, and are offline fixture publication rights available?

Do not register as a supplier or recipient to answer these questions.

## Quantium variant

Ask whether any current Woolworths/other retailer-authorized product supplies
**current local store offers** and licences independent consumer CLI display/API
use, rather than transaction history, survey or commercial analytics. Historical
Woolworths partnership evidence does not itself establish this entitlement.
If the answer is no, close this candidate rather than substituting analytics.

## Reply evaluation and stop conditions

Record the actual written response, date, retailer/operator and public-safe
reference in #2; keep private commercial terms/credentials out of public issues.
A claimed licence without verifiable upstream authority remains unresolved. A
retailer relationship does not waive display, channel or contract gates.

Do not accept a contract, register, procure or run live requests during discovery.
If the answer requires credentials/cost, return a specific approval/budget decision
before proceeding. If every gate is satisfied for a permitted anonymous/free
route, propose an isolated offline-tested adapter; request QA and obtain separate
live-smoke authorization. No merge/release/deployment authority is implied.
