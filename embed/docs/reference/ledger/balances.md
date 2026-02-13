---
title: Balances
metaTitle: "Reference: Balances"
description: Reference for the balance resource within a Twisp Ledger.
showNext: true
showPrev: true
---

## The Basics

Balances are auto-calculated sums of the entries for a given account or account set, providing a snapshot a financial position. Balances play a crucial role in accounting, as they provide a snapshot of an account's value at any given time. In the Twisp system, accounts and balances are closely related.

Balances in Twisp...

- Are automatically calculated as entries are written
- Separate debit and credit balances maintained
- Are attached to a specific account, journal, and currency
- Roll up into normal balances based on account's normal balance type

## Components of a Balance Record

- **Account**: determines which account's entries the balance is computed with.
- **Entries**: list all the entries used to compute the balance. The most recent entry is also stored as a reference.
- **Journal**: is the journal from which entries are pulled to compute the balance. Each account's balance is calculated per journal.
- **Currency**: is the currency used by entries in the balance.
- **Balance Amounts** for each layer (SETTLED, PENDING, ENCUMBRANCE, and the dynamic AVAILABLE layer), which contain the sums for...
  - The **debit** balance of the account
  - The **credit** balance of the account
  - The **normal** balance, which is calculated difference between credits and debits (for credit-normal accounts) or between debits and credits (for debit-normal accounts).

Balance records also maintain `created`, `modified`, and `committed` timestamps as well as their [](/reference/ledger/versions-and-history).

## Point-in-time Balances

With [point-in-time queries](/reference/ledger/versions-and-history#point-in-time-queries), balance `history` is filterable by timestamp. This provides the ability to retrieve the active balance at a specific time.

## Balance Operations

Use GraphQL queries to read balances directly:

- [`Query.balance()`]($gql:query:balance): Get a single balance.
- [`Query.balances()`]($gql:query:balances): Query balances using index filters.

Balances can also be queried relationally as fields on the [Account]($gql:object:Account) and [AccountSet]($gql:object:AccountSet) objects.

## Further Reading

To learn more about querying balances, see the tutorial on [](/tutorials/pulling-balances).

For more context on how balances work, see [](/accounting-core/balances).

To review the GraphQL docs for the `Balance` type, see []($gql:object:Balance).
