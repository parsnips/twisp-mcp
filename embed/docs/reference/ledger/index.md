---
title: Ledger Reference
metaTitle: "Reference: The Ledger"
description: |
  A Twisp ledger is composed of accounts, account sets, balances, entries, journals, transactions, and tran codes.
showNext: true
showPrev: false
---

## Ledger Resources

See each of the sub-pages in this reference for details about the corresponding ledger resource. Each resource is accessible via the [GraphQL API](/reference/graphql).

{% quick-links %}

{% quick-link title="Accounts" href="/reference/ledger/accounts" icon="account" description="Accounts store value and record activity." /%}

{% quick-link title="Account Sets" href="/reference/ledger/account-sets" icon="accountSet" description="Account sets group and organize accounts, rolling up balances for member accounts." /%}

{% quick-link title="Balances" href="/reference/ledger/balances" icon="balance" description="Balances calculate sums of ledger entries across layers, journals, accounts, and currencies." /%}

{% quick-link title="Entries" href="/reference/ledger/entries" icon="entry" description="Entries are records of money movement into or out of accounts." /%}

{% quick-link title="Indexes" href="/reference/ledger/indexes" icon="index" description="Indexes enable efficient data access in the ledger." /%}

{% quick-link title="Journals" href="/reference/ledger/journals" icon="journal" description="Journals keep your financial transactions organized in separate collections." /%}

{% quick-link title="Tran Codes" href="/reference/ledger/tran-codes" icon="tranCode" description="Transaction codes define how ledger entries are written when a transaction is posted" /%}

{% quick-link title="Transactions" href="/reference/ledger/transactions" icon="transaction" description="Transactions record all accounting events in the ledger." /%}

{% quick-link title="Velocity Controls" href="/reference/ledger/velocity-controls" icon="velocityControl" description="Velocity controls regulate the rate that transactions can be posted." /%}

{% /quick-links %}

## Axioms Governing the Twisp Ledger

_This set of axioms summarizes the key behavior of a Twisp ledger._

- **Accounts** are a named store of value in Twisp, with each account having a balance.
  - Accounts record activity in the form of ledger entries posted to them.
  - Accounts support multiple layers and have full support for multiple journals.
  - Accounts are stored as immutable documents with all changes stored in a versioned history.
- **Account sets** are custom groups of accounts that aggregate balances and provide a unified interface into the entries for all accounts.
  - Account sets can contain other sets, allowing for the modeling of more complex structures for a chart of accounts.
  - A chart of accounts is a collection of account and account sets used to track money within the system.
  - Account sets belong to a single journal.
  - Account sets will materialize balances for entries posted to the account set's journal in accounts that are members of the account set or any of its descendant sets.
  - Entries can only be posted to accounts, not to account sets.
- **Balances** are derived from entries written to the ledger.
  - Balances are calculated for every account.
  - Querying balances reflects the current state of accounts in the ledger.
  - Accounts have a debit balance, credit balance, and normal balance on each layer (SETTLED, PENDING, and ENCUMBRANCE).
  - Balances are materialized for every combination of account, journal, layer, and currency.
  - Balance calculations are computed on write, not on read.
- **Entries** represent one side of a transaction in the Twisp ledger.
  - Entries always have an account, an amount (including currency), and a direction (CREDIT or DEBIT).
  - Entries can only be entered in the context of a transaction.
  - Every entry is assigned to a layer (SETTLED, PENDING, and ENCUMBRANCE) to differentiate between entries in various stages of a transaction lifecycle.
  - Posting a transaction will create at least 2 ledger entries.
- **Journals** are used to organize transactions within separate "books" in the Twisp ledger.
  - Every ledger contains a default journal.
  - Users can create, update, and delete additional journals as needed.
- **Transactions** record all accounting events in the Twisp ledger.
  - The only way to write to a ledger is through a transaction.
  - Every transaction writes two or more entries to the ledger in standard double-entry accounting practice.
  - Transactions are structured by the tran code used, ensuring that the ledger is consistent, predictable, and correct.
- **Transaction codes** (tran codes) define how ledger entries are written in Twisp.
  - Tran codes accept input values as parameters and write a multi-entry transaction to the ledger.
  - Params are key-value parameters specified by the transaction code and are used when posting a transaction with a specified tran code.
  - All transactions using tran codes are executed atomically, ensuring all-or-nothing commits.
  - Tran codes provide a centralized interface for defining every type of transactional activity in a system.
- **Velocity Controls** Define conditions under which a transaction is allowed to be written in Twisp.
  - Velocity Controls can apply to Accounts or Account Sets.
