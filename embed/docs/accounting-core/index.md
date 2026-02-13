---
title: The Accounting Core
description: An accounting engine for building products that work with money.
showNext: true
showPrev: false
---

We built the Accounting Core to serve as an engine for all kinds of financial products. It is designed to be flexible enough to cover any use case that you can imagine, and comes pre-packaged with a set of sensible defaults to give you a head start.

When you provision a new instance of the Twisp Accounting Core, you get a **transactions ledger** for double-entry accounting, a **chart of accounts** to represent any economic activity, and **layered balances** for tracking settled, pending, and planned funds flows. All of this is accessible via a straighforward [GraphQL API](/reference/graphql).

Armed with these primitives, designers and developers will have a clear mental model for how to interact with, iterate on, and scale their products with foundational ledger tooling that gives confidence when working with mission-critical financial data.

Continue reading to learn more about how each part of accounting core works:

- [**Ledgers**](/accounting-core/ledgers-in-twisp)\
  The accounting core is built upon a single source-of-truth ledger.
- [**Encoded Transactions**](/accounting-core/encoded-transactions)\
  Design accounting logic with tran codes to create composable transaction types.
- [**Chart of Accounts**](/accounting-core/chart-of-accounts)\
  A chart of accounts models all of the economic activity that your ledger records.
- [**Layered Accounting**](/accounting-core/layered-accounting)\
  Layers help clarify the true state of funds in accounts.
- [**Balances**](/accounting-core/balances)\
  With compute-on-write rollup calculations, balances always reflect the most up-to-date state of an account.
