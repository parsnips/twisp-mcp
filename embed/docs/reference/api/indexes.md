---
title: Indexes
metaTitle: "Reference : API Indexes"
description: |
  Queries against the ledger are done via database indexes. This document explains how to use them.
metaDescription: Index docs for the Twisp GraphQL API.
showNext: true
showPrev: true
---

{% comment %}

Outline:
- Heuristics for Querying with Indexes
  - Choose the smallest index possible. If you only need a single record, use an ID index.
- The Index-First Approach
- Filtering with Indexes
- Compound indexes
- Common Errors

{% /comment %}