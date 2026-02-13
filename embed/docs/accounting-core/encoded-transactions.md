---
title: Encoded Transactions
description: Design accounting logic with tran codes to create composable transaction types.
showNext: true
showPrev: true
---

## Transactions and Entries

_Transactions_ record all accounting events in the ledger. **In Twisp, the only way to write to a ledger is through a transaction.**

Every transaction writes two or more entries to the ledger in standard double-entry accounting practice.

Twisp expands upon the basic principle of an accounting transaction with additional features like transaction codes and correlations.

An _entry_ represents one side of a transaction in a ledger. In other systems, these may be called "ledger lines" or "journal entries".

Entries always have an account, amount, and direction (CREDIT or DEBIT). In addition, Twisp uses the concept of "entry types" to assign every entry to a categorical type.

Twisp enforces double-entry accounting, which in practice means that entries can only be entered in the context of a Transaction. Posting a transaction will create _at least 2_ ledger entries.

In addition, we run validity checks against transactions to ensure that they do not introduce inconsistencies into the accounting core. For example, we ensure that the debit and credit entries written by a transaction sum to zero so that value is never lost or created from nothing.

**By establishing a strict definition of how entries are written to the ledger, we ensure a high level of integrity and consistency in the ledger record.**

## Composable Double-Entry Accounting

Every transaction entered must use a transaction code to indicate what _kind_ of transaction it is, which in turn determines how entries are written to the ledger. This applies both a strong categorization scheme to transactions as well as a valuable internal reference of transaction types for product engineers to draw upon.

We think that transaction codes (tran codes) are the optimal way for engineers working on financial products do double-entry accounting. They encode the basic patterns for a type of transaction as a predictable and repeatable formula.

You can think of tran codes as functions which define how a transaction acts upon the ledger. To better understand how tran codes work, let's look at an example.

## Tran Codes in Practice

While the full API for tran codes allows for a large degree of flexibility, we'll focus on a simple case for the sake of illustration.

Let's say we're building a product that needs to support tracking ACH credits, and reflecting when money is deposited to a customer's bank account. This transaction can be encoded with a tran code.

When a new transaction is posted using this tran code, we want to make sure that:

1. A credit entry is written to the "ACH Settlement" account
2. A debit entry is written to the customer's account
3. The amount used for both entries is equal

We'll use the entry type `ACH_CR` for the credit entry, and `ACH_DR` for the debit entry to be extra-clear about what these entries represent.

Here's how we would create the tran code using GraphQL:

{% partial file="graphql/examples/basicTranCodeFlow/002_BasicTranCode.md" /%}

Did you notice that `params.amount` value for the `amount` of each entry? This is a CEL expression which means "use the `amount` field on the `params` argument provided when a new transaction is posted.

Because tran codes are essentially **templates for transactions**, this lets us dynamically set the amount field when we actually go to post a transaction with this tran code.

With this tran code defined, we can now post a transaction to perform a deposit of $12.87:

{% partial file="graphql/examples/basicTranCodeFlow/003_BasicTransaction.md" /%}

By providing just a `tranCode` and `params` to the `postTransaction` call, we are able to make use of the predefined tran code to write two complete entries to the ledger.

This is just a teaser of what you can do with tran codes. With their flexibility, you can encode nearly any kind of transaction your product needs to support.

### Example Tran Codes

The set of tran codes you need will be specific to your company and product.

There are often many "archetypal" tran codes which are commonly used. Some examples include:

| TranCode         | Description                | Types of Entries Written                 |
|------------------|----------------------------|------------------------------------------|
| WIRE_TRANSFER    | Bank-to-bank wire transfer | WIRE_OUTGOING_DR, WIRE_INCOMING_CR       |
| CARD_HOLD_CANCEL | Card hold cancellation     | CARD_HOLD_CANCEL_DR, CARD_HOLD_CANCEL_CR |
| DEPOSIT          | Bank deposit               | DEPOSIT_DR, DEPOSIT_CR                   |
| BILL_PAYMENT     | Bill payment               | BILL_PAYMENT_DR, BILL_PAYMENT_CR         |

## Timing and Sequencing

Accurate recording of the times and sequences of events in a ledger is critical for auditing and reconciliation.
### Effective Dates

There are two significant time-based values on a ledger transaction: the `created` timestamp and an `effective` date.

| Field     | Description                                             |
|-----------|---------------------------------------------------------|
| Created   | The wall time the transaction was posted to the ledger. |
| Effective | The accounting date to which this transaction applies.  |

These two values often are not the same. For example, an ACH transaction may post over the weekend, but the `effective` accounting date of the transaction may be the following Monday.

### Entry Sequences

Ledger entries are always posted in the order in which they are defined within a tran code.

When a transaction is posted, it writes this ordering as a `sequence` onto every entry written.

Within the context of a transaction, we can thus see a clear incremented sequence of all entries. Because transactions are written atomically at the database layer, every entry is posted at the same clock time.

## Embedding Meaning & Context

When attempting to trace money movement, having as much contextual information as possible is useful to get a complete picture of what happened and when.

In addition to the meaning and context implied by tran codes and entry types, additional information can be embedded into transactions through correlations and metadata.

### Correlation Identifiers

With transactional workflows it is often necessary to group a set of related transactions.

For example: during card processing there is often a hold, then a hold release or expiration, and finally a settlement. In Twisp, correlation identifiers are used to group these transactions together.

When a transaction is posted without a `correlationId`, it uses its own `transactionId` as the `correlationId`. Then, future related transactions can be posted with the same `correlationId` to indicate their relationship to the original. This is very useful for events like holds, auths, auth reversals, etc.

The transactions from the card processing example above might look like this:

| ID | Amount | Description             | Correlation ID |
|----|--------|-------------------------|----------------|
| 1  | $50    | Place card hold         | 1              |
| 2  | $50    | Release card hold       | 1              |
| 3  | $50    | Settle card transaction | 1              |

Because transactions (2) and (3) are _related_ to transaction (1), they share the same correlation ID. This way, we can easily observe the entire history of a multi-transaction event by querying the correlated transactions.

### Transaction Metadata

Transactions contain a `metadata` field which can store arbitrary structured data about the transaction in JSON format.

This can be a highly useful way to embed application- and product-specific data in the transaction itself. It has no effect on the accounting operations.
