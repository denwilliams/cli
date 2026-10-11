# Live in-store provider work — separate from the foundation

Confirmed retailers: Coles, Woolworths, ALDI. Required channel: **in-store**.
No verified permitted live in-store search feed is configured. The foundation
returns explicit unconfigured groups and does not substitute online or historical
prices. Track adapter gates in [issue #2](https://github.com/denwilliams/cli/issues/2).

## Existing evidence (planning investigation, 11 October 2026)

The [original planning report](https://github.com/denwilliams/cli/issues/1#issuecomment-6105011268)
contains the source links and distinguishes read evidence from unreadable pages:

- **ALDI:** [legal notice](https://www.aldi.com.au/legal-notice) expressly restricts
  automated script/crawler access. Do not implement a direct scraper. A permitted
  feed/permission and store-specific in-store pricing contract remain unresolved.
- **Coles:** [website terms](https://www.coles.com.au/important-information/terms/website-terms-conditions)
  restrict use/reuse; pricing's exclusion from IP “Content” is not blanket adapter
  permission for website access or metadata reuse. [Customer agreement §7](https://www.coles.com.au/important-information/customer-agreement)
  distinguishes shopping channels. A documented permitted in-store source and
  location resolution remain unresolved. The original read was limited to the
  returned first 20,000 characters, which included §7.
- **Woolworths:** terms/pricing returned 403; they were not read or bypassed. Its
  [API-labelled page](https://app.woolworths.com.au/apis/v2/swagger-ui/index.html)
  yielded no readable external onboarding or access contract. Permissions,
  in-store semantics and location requirements remain unresolved.
- Historical/crowdsourced observations, including Open Prices, are not an
  unannounced replacement for current retailer in-store prices. Marketplace
  listings alone do not verify retailer-data rights or provenance.

These are the earlier report's findings, not a legal opinion or proof that no
permitted source exists. They were not all re-fetched during foundation work.

## Follow-up read, 11 October 2026

[Apify's community AU grocery scraper](https://apify.com/lergassy/au-grocery-scraper)
was read (first 20,000 characters). Its own documentation explicitly describes
national **online** prices without store/postcode selection, not local in-store
prices. It also describes proxy/address-pool switching after access denials.
**Not an eligible route** for this requirement/no-bypass constraint. Vendor claims
were not independently validated. No actor was run, nor any proxy, account,
subscription or purchase used.

The [three-retailer RapidAPI listing](https://rapidapi.com/highbrowfame/api/australian-supermarket-prices-api-coles-woolworths-aldi)
returned HTTP 200 but **no readable text**. Its rights, upstream permission,
pricing channel, freshness and contract remain unverified; no API call or
purchase was made.

## Gates before implementing each live adapter

Engineering must establish:

1. Permitted access, upstream provenance and display/reuse rights, without
   access-control circumvention or prohibited automated ALDI access.
2. Retailer identity, **in-store channel**, actual store/postcode resolution and
   refusal/ambiguity behavior; never invent a default store.
3. Observation/retrieval freshness, catalogue/sample/truncation scope and supported
   search/pagination/sort behavior; distinguish a sorted sample from global top-N.
4. Exact AUD money, original pack/unit details, conditional/weight price semantics
   and representative offline contract fixtures.
5. Rate limits, cancellation/error handling and any credential/cost requirements.
   Account creation, credentials, services and spending need separate approval;
   they are not authorized by the foundation task.
6. Offline malformed-payload, scope, context, cancellation and rate-limit tests;
   independent review and separately authorized opt-in live smoke checks.

Do not ask the user to find feeds. Genuine access/provenance blockers remain
tracked separately and do not prevent reviewing the foundation. No live adapter,
release, production change or recipes-cli import is part of this PR.

## Renewed discovery after foundation merge — 11 October 2026

PR #3 is merged. The [renewed evidence ledger](live-source-investigation-2026-10-11.md)
records retailer/supply-chain routes, newly rejected comparison services,
unreadable evidence, per-retailer unresolved gates and next actions. No eligible
live in-store feed was verified; all three providers remain unconfigured.
The [permission enquiry draft](live-source-permission-enquiry.md) is prepared but
not sent; outreach/procurement/live testing are not implicitly authorized.
