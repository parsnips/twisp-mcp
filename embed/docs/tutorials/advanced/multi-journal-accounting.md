---
title: Multi-Journal Accounting
metaTitle: "Tutorial: Multi-Journal Accounting"
description: |
  In this tutorial, we will explore how to perform multiple-journal accounting using the Twisp Accounting Core.
showNext: true
showPrev: true
---

Our primary focus will be on setting up a multi-journal ledger and posting various transactions across these journals to demonstrate how balances are materialized in accounts and account sets across journals.

We will walk through the following steps:

1. Setup multiple journals and an account set in each journal
2. Post transactions to each journal
3. Query ledger entries across both journals

## Step 1: Multi Journal Setup

We'll create two journals: **"Journal 1"** as the primary journal and **"Journal 2"** as the secondary journal.

Then, we'll create two account sets, one for each journal, both with a credit-normal balance type.

Next, create two accounts "Account A" and "Account B" and associate them with both account sets. Finally, create a transfer transaction code with the necessary parameters and transaction entries.

{% partial file="graphql/examples/journalsAndAccountSets/001_MultiJournalSetup.md" variables={title: "Setup Multi-Journal Ledger"} /%}

After completing this setup, our chart of accounts will look like this:

| Name      | Description                                             | Normal Balance Type | Record Type |
|-----------|---------------------------------------------------------|---------------------|-------------|
| Account A | Account A                                               | CREDIT              | Account     |
| Account B | Account B                                               | CREDIT              | Account     |
| Accounts  | Group entries in all accounts for the primary journal   | CREDIT              | AccountSet  |
| Accounts  | Group entries in all accounts for the secondary journal | CREDIT              | AccountSet  |

We'll also have a simple `XFR` tran code for moving money between accounts.

## Step 2: Post Transactions

Next, let's post two transactions using the `XFR` transaction code created above to write some entries to each journal.

In the first transaction, we'll move $11.33 from Account A to Account B in _Journal 1_. In the second transaction, we'll transfer $22.44 from Account B to Account A in _Journal 2_.

{% partial file="graphql/examples/journalsAndAccountSets/002_PostTransactions.md" variables={title: "Post TXs to Each Journal"} /%}

## Step 3: Account Set Entries

Finally, let's investigate the entries written in each journal by the transactions posted.

We'll retrieve account sets 1 and 2 using their respective IDs. For each account set, we'll get the member accounts and entries. Also, we'll get the journals' balance (in USD).

{% partial file="graphql/examples/journalsAndAccountSets/003_AccountSetEntries.md" variables={title: "Query Ledger Entries"} /%}

For easier reading, here are the entries listed out:

| Type     | Journal   | Account   | Direction | Amount |
|----------|-----------|-----------|-----------|--------|
| `XFR_CR` | Journal 1 | Account B | `CREDIT`  | $11.33 |
| `XFR_DR` | Journal 1 | Account A | `DEBIT`   | $11.33 |
| `XFR_CR` | Journal 2 | Account A | `CREDIT`  | $22.44 |
| `XFR_DR` | Journal 2 | Account B | `DEBIT`   | $22.44 |

## Summary

We began by establishing two separate journals and creating account sets and accounts for each. We then created a transfer transaction code that allowed us to post transactions across these journals. By posting two sample transactions, we illustrated how balances are materialized in accounts and account sets across different journals.

This tutorial provides a simplified into managing complex financial scenarios where multiple journals are required, such as consolidating the financial activity of multiple subsidiaries in a parent company or tracking transactions in various currencies.
