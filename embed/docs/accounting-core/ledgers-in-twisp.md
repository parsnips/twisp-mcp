---
title: Ledgers in Twisp
description: The accounting core is built upon a single source-of-truth ledger.
showNext: true
showPrev: true
---

## What is a Ledger?

Financial ledgers are the foundation of systems that track money. Double-entry accounting has been the de facto mechanism to record financial transactions for over 500 years, and ledgers are at the heart of this system.

- Banks track debits and credits across many accounts, with multiple financial instruments, all the while analyzing this data to surface financial insights to consumers.
- Payment infrastructure companies need robust multi-tenant account servicing built on many FBOs, made available to many clients via a secure API.
- Accounting software depends on strong transactional guarantees and zero-downtime data access to ensure accuracy and reliability.

Financial ledgers and the data derived from them are the engine of all such systems.

## Ledgers in Twisp

At Twisp, we set out to rethink the underlying technology for financial ledger systems by combining the operational and scaling characteristics of a distributed database with the correctness guarantees offered by relational databases.

{% callout type="note" %}

To learn more about the infastructure that powers our ledgers, read the [Infrastructure](/infrastructure) section.

{% /callout %}

The ledger for our accounting core is designed to be flexible enough to support any financial product, yet structured enough to provide a reliable and trustworthy store of financial data.

## Architecture

With your Twisp account, you can provision multiple instances of the accounting core. Each instance is powered by a single ledger.

A ledger is composed of journals, accounts, entries, transactions, and balances.

### Journals

Journals allow for the organizing of transactions within separate "books".

In many cases, users only need a single journal. For this reason, Twisp always contains a default journal with code `DEFAULT`.

Journals can be used for a variety of functions. For example, users may create separate journals for different currencies, or product-specific journals.

### Accounts

A chart of accounts models all of the economic activity that your ledger provides.

The chart of accounts is the basis for creating balance sheets, P&L reports, and for understanding the balances for the customer and business entities your business services.

Read more about accounts in [Chart of Accounts](/accounting-core/chart-of-accounts).

### Layered Entries

An entry represents one side of a transaction in a ledger. In other systems, these may be called "ledger lines" or "journal entries".

Twisp enforces double-entry accounting, which in practice means that entries can only be entered in the context of a Transaction. Posting a transaction will create _at least 2_ ledger entries.

Entries always have an account, amount, direction (CREDIT or DEBIT), and type for denoting an particular entry category (e.g. "TRANSFER_DR" for an entry representing the debit side of a transfer). In addition, every entry is assigned to a layer (SETTLED, PENDING, and ENCUMBRANCE) to differentiate between entries in various stages of a transaction lifecycle.

Read more about how layers work in [Layered Accounting](/accounting-core/layered-accounting).

### Transactions

Transactions record all accounting events in the ledger. In Twisp, the only way to write to a ledger is through a transaction and every transaction writes two or more entries to the ledger in standard double-entry accounting practice.

Twisp expands upon the basic principle of an accounting transaction with additional features like transaction codes and correlations. Transaction Codes (tran codes) are how financial engineers do double-entry accounting: they encode the basic patterns for a type of transaction as a predictable and repeatable formula.

Read more about the power of transactions and tran codes in [Encoded Transactions](/accounting-core/encoded-transactions).

### Balances

Balances are auto-calculated sums of the entries for a given account.

Each account can have balances across all three layers, and every balance record maintains a separate debit and credit balance amount.

Read more about how balances are calculated in [Balances](/accounting-core/balances).
## Architected from Principles

Our ledger has been designed in adherence to strong, time-tested accounting principles.

- **Enforce double-entry principles**.\
  In other words, ensure that value (money) is never created or destroyed. There is always a debit and credit side for every transaction.
- **History must be fully preserved.**\
  Ledgers are immutable and append-only. There is always a full history of changes to provide a clear, unbroken lineage of the current state.
- **Contextualize every transaction**.\
  Much financial activity happens across multiple phases of a lifecycle, with different transactions at each phase. Proper use of correlations, layers, transaction metadata means that no useful context is lost across that lifecycle.
- **Design for composability**.\
  With transaction codes (tran codes) and automations, higher-order accounting logic can be encoded into a simplified interface. This way, your designers and engineers can focus on the business logic of your product, not the accounting details.
- **When in doubt, rely upon industry standards**.\
  There are times to be inventive, and times to trust the system. We've learned from experience and research to make use of well-established accounting patterns like layers and typed entries.
