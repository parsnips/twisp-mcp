---
title: Entries
metaTitle: "Reference: Entries"
description: Reference for the entry resource within a Twisp Ledger.
showNext: true
showPrev: true
---

## The Basics

Entries are records of a money movement into or out of an account in the Twisp ledger. They represent the individual line items in a ledger journal, storing details such as the account and posting transaction.

{% callout type="note" %}

In Twisp, an _entry_ is a fundamental accounting concept representing **one side of a transaction** in a ledger.

{% /callout %}

Transactions are the only way to write to a Twisp ledger. Every transaction posted to a ledger writes at least two entries (one or more debits and one or more credits) which balance to zero, in accordance with double-entry accounting principles. This ensures that the accounting system maintains a high level of integrity and consistency in the ledger record.

In other accounting systems, entries might be called "ledger lines" or "journal entries." However, the fundamental concept remains the same. In Twisp, entries always have an account, amount (in units of a currency), direction (`CREDIT` or `DEBIT`), and an entry type that assigns every entry to a categorical type. These entries are the building blocks of a comprehensive, accurate, and reliable financial record.

## Components of a Ledger Entry

1. **Account**: to which the entry is written. This account represents the financial aspect (e.g., asset, liability, expense) impacted by the entry.
2. **Amount**: expressed in numerical (decimal) **units** of a specified **currency** denomination.
3. **Direction**: denotes which side of the ledger the entry is posted on. This is a crucial aspect of double-entry accounting, as every transaction must have at least one entry on the debit side and one on the credit side.
4. **Layer**: on which this entry is recorded (SETTLED, PENDING, or ENCUMBRANCE).
5. **Type**: a categorical label for the entry, such as "TRANSFER_DR" for an entry representing the debit side of a transfer.
6. **Transaction**: references the transaction record in which the entry was written.
7. **Journal**: within which the entry was written.
8. **Sequence**: indicates the order in which the entry was written within the context of the posting transaction.

In addition, entries have other properties which are common to most or all resources in the accounting core:

- **ID**: a universally unique identifier (UUID) for the entry.
- **Name**: a descriptive name.
- **Description**: a free-form text to be used for describing anything about the entry. We recommend using this field to summarize what the entry is for and when it should be used.
- **Metadata**: unstructured, user-specified JSON. Can be combined with custom indexes for powerful querying capacities.
- **Created & Updated Timestamps**: self-evident: when the entry was created and when it was last updated.
- **Version & History**: entries, like every other record in the accounting core, maintain a list of all changes made to them in their `history` field, and the `version` field indicates the current active version of the entry.

## Entry Operations

Entries cannot be written directly, only indirectly by posting transactions.

Entries can be queried using GraphQL:

- [`Query.entry()`]($gql:query:entry): Get a single entry.
- [`Query.entries()`]($gql:query:entries): Query entries using index filters.

Entries can also be queried relationally as fields on the [Account]($gql:object:Account), [AccountSet]($gql:object:AccountSet), [Balance]($gql:object:Balance), and [Transaction]($gql:object:Transaction) objects.

## Further Reading

To learn how write entries through posting transactions, see the tutorial on [](/tutorials/posting-transactions).

For more context on how and why Twisp a layered accounting model, see [](/accounting-core/layered-accounting).

To review the GraphQL docs for the `Entry` type, see []($gql:object:Entry).
