---
title: Example Setup
description: Use this example ledger configuration to start exploring Twisp and learning how to build your accounting system.
showNext: false
showPrev: false
---

This setup is for an example budgeting app called "Budget Buddy". Think of it as a mashup of a budgeting tool like [Mint](https://mint.intuit.com/) or [YNAB](https://www.ynab.com/) with a social wallet app like [Venmo](https://venmo.com/).

**The full setup GraphQL code can be found at the bottom of this page.**

## Account Structure

Budget Buddy's chart of accounts includes account sets for each user, containing that user's wallet, budget accounts, and linked financial accounts like checking or credit cards. It also has offset and settlement account.

### User Accounts

```mermaid
graph BT
  bert[/Bert\]
  ernie[/Ernie\]
  users[/Users\]
  bert & ernie --> users
```

```mermaid
graph BT
  bert[/Bert\]

  bert_budget[Bert's Budget]
  bert_cash[Bert's Cash Acct.]
  bert_checking[Bert's Checking Acct.]
  bert_wallet[Bert's Wallet]

  bert_budget & bert_cash & bert_checking & bert_wallet --> bert
```

```mermaid
graph BT
  ernie[/Ernie\]

  ernie_budget[Ernie's Budget]
  ernie_checking[Ernie's Checking Acct.]
  ernie_credit_card[Ernie's Credit Card Acct.]
  ernie_wallet[Ernie's Wallet]

  ernie_budget & ernie_checking & ernie_credit_card & ernie_wallet --> ernie
```

### Settlement Accounts

```mermaid
graph BT
  settlement[/Settlement\]
  settlement_card[Card Settlement]
  settlement_cash[Cash Settlement]
  settlement_checking[Checking Settlement]

  settlement_card & settlement_cash & settlement_checking --> settlement
```

### Additional Sets

Additional account sets are used to roll up wallet and budget accounts:

```mermaid
graph BT
  budgets[/Budgets\]
  wallets[/Wallets\]

  bert_budget[Bert's Budget]
  bert_wallet[Bert's Wallet]
  ernie_budget[Ernie's Budget]
  ernie_wallet[Ernie's Wallet]
  budget_offset[Budget Offset]

  budget_offset & bert_budget & ernie_budget --> budgets
  bert_wallet & ernie_wallet --> wallets
```

## Tran Codes
The setup defines several `TranCodes`:

- `ALLOC_BUDGET`: Allocates an amount to a specific budget, creating a credit entry in the account and a debit entry in the budget offset account.
- `DEALLOC_BUDGET`: Deallocates funds from a budget, creating a credit entry in the budget offset account and a debit entry in the account.
- `ASSIGN_TO_BUDGET`: Assigns a transaction to a particular budget, creating a credit entry in the account and a debit entry in the budget offset account.
- `RECORD_TX`: Records a transaction between two accounts, creating a debit entry in one account and a credit entry in another.
- `RECORD_PENDING_TX`: Records a pending transaction between two accounts, creating a debit entry in one account and a credit entry in another.
- `RECORD_SETTLE_PENDING_TX`: Records the settlement of a pending transaction, creating a credit entry in the corresponding account and a debit entry in the settlement account.
- `WALLET_TRANSFER`: Transfers funds between two wallets, creating a debit entry in one wallet and a credit entry in another.
- `WALLET_DEPOSIT`: Deposits funds into a wallet, creating a debit entry in the wallet account and a credit entry in the checking account.
- `WALLET_WITHDRAW`: Withdraws funds from a wallet, creating a debit entry in the checking account and a credit entry in the wallet account.

Each `TranCode` defines the entries that should be created in the ledger when the transaction is processed, as well as any metadata that should be associated with the transaction.

## GraphQL

{% partial file="graphql/examples/docsBasics.md"
  variables={
    title: "Standard Tutorial Setup"
  }
/%}