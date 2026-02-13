---
title: Accounts
metaTitle: "Reference: Accounts"
description: Reference for the account resource within a Twisp Ledger.
showNext: true
showPrev: true
---

## The Basics

Accounts are a named **store of value** that hold a balance, as well as a **record of activity** in the form of ledger entries posted to them. They are an essential part of the Twisp accounting core and are used to model all of the economic activities your ledger provides.

Additionally, organizing accounts into sets allows for easy aggregation of balances and querying of ledger entries, enabling you to create a more flexible and powerful structure for your chart of accounts.

## Components of Accounts

There are 5 key components defining an account:

1. **Code**: a unique name to identify the account, usually formatted in UPPER_SNAKE_CASE.
2. **Balances**: reflect the net result of all transactions affecting the account. Balances are materialized for each layer, currency, and journal used by entries in the account.
3. **Entries**: the individual journal entries written in transactions posted to the account, providing a detailed history of all financial events that have affected the account over time.
4. **Sets**: list the account sets of which this account is a member.
5. **Normal Balance Type**: determines how the account's normal balances are computed.

In addition, accounts have other properties which are common to most or all resources in the accounting core:

- **ID**: a universally unique identifier (UUID) for the account.
- **Name**: a descriptive name.
- **Description**: a free-form text to be used for describing anything about the account. We recommend using this field to summarize what the account is for and when it should be used.
- **Metadata**: unstructured, user-specified JSON. Can be combined with custom indexes for powerful querying capacities.
- **Config**: allows setting up an account for concurrent posting.
- **Created & Updated Timestamps**: self-evident: when the account was created and when it was last updated.
- **Version & History**: accounts, like every other record in the accounting core, maintain a list of all changes made to them in their `history` field, and the `version` field indicates the current active version of the account.

## Account Operations

Use GraphQL queries and mutations to read, create, update, and delete (lock) accounts:

- [`Query.account()`]($gql:query:account): Get a single account.
- [`Query.accounts()`]($gql:query:accountSets): Query accounts using index filters.
- [`Mutation.createAccount()`]($gql:mutation:createAccount): Create a new account.
- [`Mutation.updateAccount()`]($gql:mutation:updateAccount): Update select fields for an existing account.
- [`Mutation.deleteAccount()`]($gql:mutation:deleteAccount): Soft delete (lock) a specified account.
- [`Mutation.addToAccountSet()`]($gql:mutation:addToAccountSet): Add an account to a set.
- [`Mutation.removeFromAccountSet()`]($gql:mutation:removeFromAccountSet): Remove an account from a set.

## Further Reading

To learn how to work with accounts, see the tutorial on [](/tutorials/setting-up-accounts).

For more context on how and why Twisp uses accounts, see [](/accounting-core/chart-of-accounts).

To review the GraphQL docs for the `Account` type, see []($gql:object:Account).

