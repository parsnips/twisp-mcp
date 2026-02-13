---
title: Transactional Operations
metaTitle: "Reference : Transactional Operations"
description: |
  Requests to the GraphQL API are executed in a single transaction context with all-or-nothing semantics.
metaDescription: Transactional operation semantics in the Twisp GraphQL API.
showNext: true
showPrev: true
---

Requests may contain multiple query or mutation operations. Every operation within a single request is executed within the same database transaction context, meaning that each operation will only succeed if _all_ operations succeed.

In other words, if there is an error in any single operation, _none_ of the operations will succeed. Transactional database operations are useful because they help maintain data integrity, consistency, and reliability in applications that require complex operations or involve multiple data manipulation steps.

For example, the following request includes two `postTransaction` mutations, but the second uses incorrect syntax that triggers a `JSON_PARSE_ERROR`.

{% partial file="graphql/reference/api-basics/transactionalOperations.md" variables={title: "Example"} /%}

Note that the `"data"` field in the response is `null` because _neither_ transaction was posted. If we were to subsequently query for the first (valid) transaction, we would see that it was not posted.
