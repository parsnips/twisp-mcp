---
title: Transaction Codes
metaTitle: "Reference: Transaction Codes"
description: Reference for the tran code resource within a Twisp Ledger.
showNext: true
showPrev: true
---

## The Basics

Transaction codes (AKA "tran codes") define how ledger entries are written in Twisp. You can think of them as macros or templates for transactions; instead of writing out every transaction by hand within your application layer, you design the tran codes to match every type of transaction in your funds flow.

{% callout type="note" %}

In a nutshell, transaction codes **accept input values** and **write a multi-entry transaction to the ledger**.

{% /callout %}

In this way, tran codes form a _self-documenting API for your funds flow_. The benefits of this approach are legion:

- Strong separation of concerns between accounting logic (i.e. how ledger entries need to be written to which accounts) and product/business logic.
- All transactions are executed atomically. Enforcing all-or-nothing commits leads to correct and easy-to-reason-about systems.
- Centralized interface for defining every type of transactional activity that your system handles, providing a rich view into your funds flow.
- A linear—not exponential—growth curve in system complexity as more and more transaction types are added to reflect the growth of your financial product.

With tran codes as the mechanism for structuring your transactions, you can reign in your funds flow and prevent it from becoming a tangled and unmanageable mess.

## Components of Tran Codes

There are 4 primary components which define a tran code:

1. **Code**: a unique name to identify the tran code, usually formatted in UPPER_SNAKE_CASE. When posting a transaction, this is the value supplied to the `tranCode` field to specify which tran code should be invoked.
2. **Params**: the definition of the parameters available when invoking a tran code. This allows values to be supplied at runtime which can be injected into the [Transaction]($gql:object:Transaction) and [Entries]($gql:object:Entry) written by the tran code.
3. **Transaction**: the specification of values to be used for the [Transaction]($gql:object:Transaction) written. Invoking a tran code writes a single transaction into the ledger, and the values used for that transaction are defined either as literals or else are derived from values passed in `params`.
4. **Entries:** the specification of values to be used for the ledger [Entries]($gql:object:Entry) written for the transaction when a tran code is invoked. The values used for these entries are defined either as literals or else are derived from values passed in `params`.

In addition, tran codes have other properties which are common to all records in the accounting core:

- **ID**: a universally unique identifier (UUID) for the tran code.
- **Description**: a free-form text to be used for describing anything about the tran code. We recommend using this field to summarize what the tran code is for and when it should be used.
- **Created & Updated Timestamps**: self-evident: when the tran code was created and when it was last updated.
- **Version & History**: tran codes, like every other record in the accounting core, maintain a list of all changes made to them in their `history` field, and the `version` field indicates the current active version of the tran code.

## Tran Code Invocation

Whenever you post a transaction in Twisp, you must supply a tran code to invoke. _There is no way to post transactions outside of a tran code invocation._ In this sense, posting a transaction is how you invoke a tran code. One cannot happen without the other.

Posting a transaction requires only three components: a `transactionId` to ensure idempotency for the transaction, the `tranCode` identifier matching the `code` field of the tran code to invoke, and the `params` object to provide parameter values for the tran code.

{% partial file="graphql/examples/tranCodeDocs/achCredit/postAchCreditTx.md" variables={title: "Example"} /%}

The above example will write a [Transaction]($gql:object:Transaction) and [Entries]($gql:object:Entry) as defined in the `ACH_CREDIT` tran code, using the parameters supplied.

## Tran Code Operations

Use GraphQL queries and mutations to read, create, update, and delete (lock) tran codes:

- [`Query.tranCode()`]($gql:query:tranCode): Get a single tran code.
- [`Query.tranCodes()`]($gql:query:tranCodes): Query tran codes using index filters.
- [`Mutation.createTranCode()`]($gql:mutation:createTranCode): Create a new transaction code.
- [`Mutation.updateTranCode()`]($gql:mutation:updateTranCode): Update select fields for an existing transaction code.
- [`Mutation.deleteTranCode()`]($gql:mutation:deleteTranCode): Soft delete (lock) a specified transaction code.

## Further Reading

To learn how to work with tran codes, see the relevant tutorials on [](/tutorials/building-tran-codes) and [](/tutorials/advanced/designing-tran-codes).

For more context on how and why Twisp uses tran codes, see [](/accounting-core/encoded-transactions).

To review the GraphQL docs for the `TranCode` type, see []($gql:object:TranCode).