---
title: Versions and History
metaTitle: "Reference: Record Versions and History"
description: |
  Every record in a Twisp ledger maintains a full history of all changes applied to it. All previous versions of a record can be queried.
showNext: true
showPrev: true
---

## The Basics

Record versioning and history is the foundation of [append-only immutability](/infrastructure/ledger-database#append-only-immutability) in the Twisp Accounting Core. All data entering the Twisp system is written as a new data record that cannot be modified. Prohibiting modification or deletion of data ensures a complete log of all changes to the system, allowing for auditable data history.

{% callout %}
Append-only immutability is the only data storage model in the Twisp Accounting Core and cannot be disabled.
{% /callout %}

## Record Versioning

Twisp stores all data within records in the storage system. Each record is assigned a unique record ID. The indexing system provides pointers to records for specific keys. For example, an `accountId` index points to the specific corresponding `account` record.

The first version of a new record is assigned `version` 1. When any operation needs to modify that record, a second record is written with the same unique record ID, but with `version` now set to 2. This pattern repeats indefinitely for every record in the system.

{% callout %}
For a particular record, version numbers form a consecutive integer sequence. Each new version number is exactly one greater than the previous version number.
{% /callout %}

In addition to the `version` number, each record version contains a number of useful timestamps:

1. `created`: the time the transaction began that created the first record version.
2. `modified`: the time the transaction began that wrote the record version. If version is 1, `modified` is equal to the `created` timestamp.
3. `committed`: the time the transaction committed that wrote the record version.

## Version History

All record versions are queryable via the `history` API available on every object in the system. This allows for retrospective analysis, auditing, and verification of data changes over time.

## Point-in-time Queries

In addition to scanning an entire record version history, Twisp provides the ability to return versions active at specific point-in-time. By specifying timestamps in the `history` API `where` filter, version history can be filtered by either the `modified` or `committed` record version timestamps.

{% partial file="graphql/examples/versionAndHistoryDocs/balanceHistory.md" variables={title: "Example"} /%}
