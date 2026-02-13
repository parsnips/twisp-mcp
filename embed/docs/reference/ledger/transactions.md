---
title: Transactions
metaTitle: "Reference: Transactions"
description: Reference for the transaction resource within a Twisp Ledger.
showNext: true
showPrev: true
---

## The Basics

Transactions in Twisp record all accounting events in the ledger, utilizing transaction codes to automate and categorize entries.

Through tran codes, Twisp automates and categorizes ledger entries, streamlining the accounting process and ensuring consistency.

{% callout type="note" %}

Transactions capture all **financial activities** and leverage tran codes to maintain a well-organized and accurate ledger.

{% /callout %}

Transactions in Twisp...

- Utilize double-entry accounting principles (ensures that the debit and credit entries balance out)
- Ensure atomic, all-or-nothing commits for predictable and error-resistant ledgering
- Leverage transaction codes (tran codes) to automate and categorize entries
- Maintain a complete history of all changes to records

## Components of Transactions

There are 5 primary components which make up a transaction:

1. **Entries**: written to the ledger by this transaction.
2. **Journal**: in which the transaction is posted.
3. **Tran Code**: provides a reference to the tran code used when posting this transaction.
4. **Effective**: date when the transaction is recorded as occurring for accounting purposes.
5. **Correlated**: transactions connect related transactions together, providing context and improving traceability within the ledger.

In addition, transactions have other properties which are common to most or all resources in the accounting core:

- **ID**: a universally unique identifier (UUID) for the transaction.
- **Description**: a free-form text to be used for describing anything about the transaction.
- **Metadata**: unstructured, user-specified JSON. Can be combined with custom indexes for powerful querying capacities.
- **Created & Updated Timestamps**: self-evident: when the transaction was created and when it was last updated.
- **Version & History**: transactions, like every other record in the accounting core, maintain a list of all changes made to them in their `history` field, and the `version` field indicates the current active version of the transaction.

## Transaction Operations

Use GraphQL to query, post, and update transactions:

- [`Query.transaction()`]($gql:query:transaction): Get a single transaction.
- [`Query.transactions()`]($gql:query:transactions): Query transactions using index filters.
- [`Mutation.postTransaction()`]($gql:mutation:postTransaction): Post a transaction to the ledger using a specified tran code.
- [`Mutation.updateTransaction()`]($gql:mutation:updateTransaction): Update select fields for a posted transaction.

## Further Reading

To learn the basics of posting transactions, see the tutorial on [](/tutorials/posting-transactions).

For more context on why Twisp structures transactions with tran codes, see [](/accounting-core/encoded-transactions).

To review the GraphQL docs for the `Transaction` type, see []($gql:object:Transaction).
