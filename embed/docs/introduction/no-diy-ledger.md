---
title: DIY Ledgers are a Problem
description: We built a ledger database and accounting engine (so that you don't have to).
showNext: true
showPrev: true
---

The Twisp founding team has spent a cumulative total of many years working on and with financial software systems for neobanks and other fintech companies. In that time, we've learned a few things:

1. Designing, building, and managing a central accounting ledger is hard.
2. But it can appear simple, so many startups (and mature companies) decide to build their ledger in-house.
3. Most of the time, this does not go well.

Accounting systems often suffer from the full range of technical debt problems common to long-lived software: data model lock-in, performance bottlenecks due to poor scaling properties, esoteric design features and programming patterns unique to one or two instrumental early engineers, etc. However, these problems be especially painful and frustrating when they apply to your internal ledger because it is both a _vital part of the business_ and also (most likely) _not a core competency_.

Our friend [Matt Brown](https://www.matttbrown.com/) from Matrix Partners calls this ["undifferentiated heavy lifting"](https://www.matttbrown.com/notes/undifferentiated-heavy-lifting), and it is particularly applicable to the ledger problem:

> Ledgers are a critical part of the tech and financial stack. Fundamentally, they're databases that contain debits and credits and bake in business-specific assumptions and logic to derive balances and other financial statements. However, they're also incredibly complex, requiring high uptime and performance with strict security, data integrity, audit, and regulatory requirements. Multiple business processes are built on them, so failures not only degrade UX but can cost significant sums of money.

The other thing we learned is that ledgers are (mostly) a solved problem. Double-entry accounting has been around for centuries. The algorithms and architectures needed to support highly availalble, reliable, and performant ledger-type data stores have been honed for decades. There is no need to re-invent the wheel.

At Twisp, we offer a simple solution to the design problem of ledgers: don't build it yourself. All of the fantastic new fintech products will need a ledger, but building internal ledgers is no longer a worthwhile investment.