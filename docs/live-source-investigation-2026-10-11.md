# Live in-store source investigation — 11 October 2026

Tracks [issue #2](https://github.com/denwilliams/cli/issues/2). Engineering's
[plan and acceptance gates](https://github.com/denwilliams/cli/issues/2#issuecomment-6105605139)
were posted before renewed discovery. Repository baseline: merged PR #3,
`de27bf949bfe5c66b4a0ad29e505e83dcf404212`. The issue, `docs/live-sources.md`,
production registry and current provider/search contracts were read first.

**Outcome: no permitted live in-store feed verified for any of the three
retailers.** This is a bounded research outcome, not proof that no permitted feed
exists or a legal opinion. No adapter or live coverage is delivered. All production
providers remain unconfigured. Search snippets were leads, not evidence.

## Evidence actually read in this renewed pass

All reads below were on **11 October 2026 (UTC)**. Unless marked otherwise, the
returned readable page text was read in full, with no tool truncation reported.
Marketing statements remain the source's claims; no data/API execution verified
them. Historical relationship evidence is not a current licence grant.

### Retailer and supply-chain routes (prioritized)

| Route and source | Read evidence | Assessment / missing gate |
| --- | --- | --- |
| [GS1 Australia National Product Catalogue (NPC)](https://www.gs1au.org/services/data-and-content/national-product-catalogue) | Product/item information shared with trading partners in the supply chain; registration, login, fees/terms and enquiry links. | Credible documented data-sharing infrastructure, **not evidence of a public store-specific consumer shelf-price feed or consumer-display licence**. Real-time product-data sharing does not establish real-time shelf prices. No free anonymous recipient entitlement verified. |
| [GS1 NPC Cookbook landing page](https://www.gs1au.org/resources/how-to-guides/national-product-catalogue-user-guide-cookbook) | Designed for publishers loading/maintaining data. Links to recipient access, publisher connection options/costs, XML catalogue and price-notification samples, and “Woolworths pricing 101”. Recipient connection options require contacting the NPC technical team. | Landing page only: linked manuals, samples, fees and pricing guide **not read**. A “price notification” label does not establish retail rather than supplier/cost price, store granularity, licence or freshness. Requires an eligibility/semantics enquiry before considering procurement. |
| [Coles Supplier Portal / Product Catalogue entry](https://www.supplierportal.coles.com.au/wps/portal/web/ElectronicTrading/ColesProductCatalogue) | Redirected to portal Home: Supplier Central manages account permissions/reports and Cost Price Alteration requests; several functions require login. NPC is a Quick Link. | Supplier account workflow, not verified outbound shelf-price search. Cost price alterations must not be treated as consumer prices. Account eligibility, consumer reuse, store semantics and API contract unresolved. No registration/login attempted. |
| [Woolworths Group Our partners](https://www.woolworthsgroup.com.au/au/en/who-we-are/our-approach-/our-partners.html) | Describes trade/supplier relationships and supplier surveys; Advantage industry data concerns relationship-performance benchmarks. | Not a product-price feed or consumer data-sharing permission. Supplier/dispute contacts are not a verified pricing-data onboarding route. |
| [Woolworths/Quantium announcement, 20 April 2021](https://www.woolworthsgroup.com.au/content/dam/wwg/investors/investor-news/2021/Woolworths%20Group%20deepens%20partnership%20with%20Quantium.pdf) | Entire three-page PDF read after text extraction (4,436 characters). Retailer confirms historical partnership and products/services providing insights to Woolworths and suppliers; discusses pricing/ranging/promotions and confidentiality. | Independent retailer evidence of a historical relationship, **not authorization for this CLI**, proof of current ownership, or live shelf-price distribution rights. |
| [Quantium Q.Checkout](https://quantium.com/q-checkout/) and [Q.Panel](https://quantium.com/q-panel/) | Q.Checkout describes first-party retailer transaction analytics, product/store performance and consultant/onboarding access. Q.Panel describes Woolworths shopper surveys validated against purchasing behaviour. | Plausible relationship-backed analytics enquiry, **not a verified live offer/catalogue API**. Transactions show what sold, not necessarily the current unconditional price of an unsold item. No consumer-display grant, anonymous/free API, exact offer payload or freshness SLA verified. No enquiry sent. |
| [ALDI Suppliers](https://www.aldi.com.au/corporate-suppliers/suppliers) | Returned only links to Become a Supplier, Supplier Information, complaints/dispute guidelines and newsletter. | No documented price-feed or external display entitlement in the returned text. Do not convert supplier eligibility into data-feed permission. No product crawling or supplier registration. |

### New comparison/data-service leads

These are **not verified licensed alternatives**. They were examined to reject
unsuitable apparent substitutes, not used as proof of upstream rights.

| Candidate / sources | Read evidence | Disposition |
| --- | --- | --- |
| [Bargeroo terms](https://www.bargeroo.com.au/terms-of-service), updated May 2026 | §2 says publicly available prices may not reflect real-time in-store/online pricing; not affiliated with supermarkets. §4 prohibits automated extraction from service/APIs and outside redistribution. §6 grants personal use only through official apps/website. §7 says prices vary by store/time. | **Reject under current terms.** Neither a public feed licence nor live actual-store contract. Do not inspect its APIs or device-integrity mechanisms. |
| [StockUp terms](https://stockup.au/terms-of-service), updated 27 September 2026 | Smart Shop claims selected-store prices from retailer systems but allows labelled estimates on failure. Acceptable Use forbids automated access/extraction and competing-service use. Premium includes live in-store features (published app subscription approximately AUD 4.99/month or 39.99/year, not a feed quote). Disclaims affiliation/current accuracy. | **Reject as an adapter source under current terms.** Store-selection marketing is not upstream authorization or downstream API/display permission. No documented consumer CLI API, rate/pagination/payload contract or free live-feed entitlement. No app/account/trial/subscription. |
| [Cannonball home](https://cannonballdata.com/), [comparison](https://cannonballdata.com/compare/cannonball-vs-scan-data), [docs index](https://cannonballdata.com/docs), [terms](https://cannonballdata.com/terms), updated 9 October 2026 | Marketing claims daily multi-chain store-level coverage. Terms say app is not open, company being registered, data read from what stores show **online**; inferred zone values marked; even observed values may differ from shelf. Subscription exports are for internal business use; no publishing/sublicensing, competing data product or bulk automated copying. Home lists monthly plans 149/490/1190/from 2500; terms specify AUD with GST where applicable. | **Reject under published scope/terms.** Daily online observations/inferred prices are not verified live in-store results. No retailer authorization established, no public CLI redistribution/automation grant and no no-cost feed. Trial/no-sign-up marketing is not permission. No app access or early-access registration. Docs child pages not read. |

### Unreadable or non-substantive evidence

Do not infer permission, prohibition, store semantics or payloads from these:

- [Coles Synergy Supplier Central](https://suppliercentral.coles.com.au/s/article/Coles-Synergy): “Loading / CSS Error / Refresh” only.
- [Coles Group Synergy announcement entry](https://www.colesgroup.com.au/news/2020/media-releases/?page=coles-synergy-helps-suppliers-plan-for-growth): “Back to News” only. Search-index claims about IRI/Circana were **not verified from this page**.
- [Woolworths NPC article](https://partnerhub.woolworthsgroup.com.au/s/article/National-Product-Catalogue-NPC) and [supplier onboarding](https://partnerhub.woolworthsgroup.com.au/s/supplier-onboarding): “Loading / CSS Error / Refresh” only.
- [ACCC Coles hearing transcript, 22 November 2024](https://www.accc.gov.au/system/files/transcript-supermarkets-inquiry-coles-22-november-2024.pdf): download returned **403**, not read. No bypass/retry via another transport.
- [Coles Group contact](https://www.colesgroup.com.au/contact-us/?page=contact-us): HTTP 200, no readable text; contact routing unverified.
- [ALDI Help Centre](https://help.aldi.com.au/s/): HTTP 200, no readable text. [ALDI Unpacked enquiry](https://www.aldiunpacked.com.au/media-enquiry/) returned journalist-only form description and a supplier-business link, with “could not locate your form”; neither is a verified data-licensing contact. Do not pose as a journalist or supplier.
- [Woolworths Group Contact Us](https://www.woolworthsgroup.com.au/au/en/contact-us.html) **was read**: general-enquiry switchboard 02 8885 0000, Mon–Fri 8am–5:30pm. A possible routing channel, not a data-licensing department or permission grant. No call made.

## Retailer-by-retailer acceptance ledger

“Unverified” means no satisfactory written evidence/contract was read, not that the
capability cannot exist. Every adapter gate remains open for every retailer.

| Gate | Coles | Woolworths | ALDI |
| --- | --- | --- | --- |
| Permitted automated source access + upstream/product-metadata rights + CLI display/reuse | Unverified; supplier portal is not a grant | Unverified; historical Quantium relationship is not a CLI licence | Unverified; existing direct-automation prohibition remains a stop gate |
| Genuinely in-store channel; actual store ID/postcode mapping; multiple-store ambiguity/refusal | Unverified | Unverified | Unverified |
| Observation time vs retrieval time; update cadence/lag/staleness policy | Unverified for an eligible feed | Unverified for an eligible feed | Unverified for an eligible feed |
| Query matching, pagination/end markers, server sort or complete enumeration, catalogue/sample/truncation | Unverified | Unverified | Unverified |
| Exact AUD amount encoding, stable ID/name, pack quantity, unit amount/unit, weight/member/multibuy conditions | No permitted representative payload | No permitted representative payload | No permitted representative payload |
| Rate limits, Retry-After, quotas, errors/cancellation | Undocumented for eligible route | Undocumented for eligible route | Undocumented for eligible route |
| Credential/eligibility requirements; price/cost; no-credential/no-cost availability | No free anonymous route verified; supplier workflow implies eligibility questions | No free anonymous route verified; consultant/onboarding analytics is not a feed quote | No free anonymous route verified; supplier portal is not feed procurement |

For each, missing values must not be filled using another channel, national
averages, transaction history, zone inference, or a vendor's scraper examples.
Store resolution must return one actual store: a postcode is not necessarily
unique. Fresh retrieval alone does not prove recent observation.

## Prior rejected routes — retained, not retried

The [earlier planning report](https://github.com/denwilliams/cli/issues/1#issuecomment-6105011268)
and [existing source notes](live-sources.md) remain the provenance for:

- ALDI legal notice's explicit automated-access restriction; no direct scraper.
- Coles pricing-data exclusion from IP “Content” is not blanket website/metadata reuse permission.
- Woolworths 403 terms/pricing reads and API-labelled shell supply no external access contract.
- Apify AU grocery scraper explicitly national online/no store context and denial-driven proxy switching: rejected.
- Actowiz acquisition description includes scraping/reverse-engineered mobile APIs/proxies: not verified authorized.
- RapidAPI listings unreadable/unverified; paid listing is not rights evidence.
- Open Prices historical/crowdsourced observations are not current retailer catalogue coverage.

Those pages were **not re-fetched in this pass**; no materially new evidence was
found to justify repeating the rejected routes. Public-document discovery is not
exhaustive. No retailer search/catalogue/price API, scraper actor or live smoke was
executed, and no product data was harvested.

## Next concrete action and decision needed

Engineering has prepared the [permission/contract enquiry packet](live-source-permission-enquiry.md).
The next gate is permission discovery, **not implementation or a scraper purchase**:

1. Obtain authorization for **non-binding retailer/GS1/Quantium enquiries only**
   (no account creation, contracts, credentials, trials or spending). Engineering
   should handle the research/questions and evaluate replies, not ask the user to
   find feeds. No outreach has been sent and no response is promised.
2. Route Coles/ALDI enquiries to their data-licensing/digital partnership teams via
   a verified general corporate channel; current Coles contact/ALDI help text was
   unreadable, so routing remains to verify. Do not use media, complaints or
   supplier-registration workflows under a false role. Woolworths' published
   general switchboard is a routing lead, not a confirmed data-rights contact.
3. Ask GS1 whether eligible recipients can receive **retailer-set consumer shelf
   prices** (not supplier costs), with store and consumer-display rights. Ask
   Quantium whether any current offering licenses current store-level offers to
   an independent CLI, rather than analytics/transactions. A “no” closes that
   branch; a “yes” requires written rights and an offline representative contract.
4. Only after a documented eligible offer, seek a separate credential/budget/
   procurement decision if necessary. Do not authorize money merely because a
   vendor advertises retailer coverage. If no zero-cost anonymous route is
   permitted and procurement is not approved, keep live features blocked.
5. A gate-cleared route can become a separate retailer adapter issue/PR, with
   offline malformed-payload, store mismatch/ambiguity, freshness, pagination/
   scope, exact money/conditional-price, cancellation and rate-limit tests.
   No invented production defaults; sort candidates before per-retailer limits;
   a sorted sample is not global cheapest top-N; conditional pack totals do not
   rank as standard unconditional prices. QA is required. Live smoke needs
   separate explicit authorization even if the feed is free.

No accounts, credentials/services, purchases, bypass, ALDI product scraping,
merge, release, deployment or production configuration changes were made.
