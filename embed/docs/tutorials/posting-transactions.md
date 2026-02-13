---
title: Posting Transactions
metaTitle: "Tutorial: Posting Transactions"
description: |
  In this tutorial, we will learn how to post transactions using the GraphQL API.
showNext: true
showPrev: true
---

{% callout type="task" %}
- Determine a tran code to use and understand its `params`
- Post a transaction with a specified tran code using the `postTransaction` mutation
- Query transactions to inspect entries written with the `transaction` query
{% /callout %}

Posting transactions in Twisp is a simple process that involves specifying the transaction details and the transaction code (tran code) to use. Transactions are written to the ledger, and entries are created for as specified by the tran code used.

---

## Prerequisites

Before you start, you should have added accounts to your ledger. See the tutorial on [](/tutorials/setting-up-accounts).

{% partial file="tutorials_getting_started.md" /%}

## Choose a transaction code

You will need to choose a transaction code to use before posting transactions.

Transaction codes are pre-defined templates that specify how transactions should be processed. Transaction codes define the accounts and amounts to debit and credit, the currency to use, and any other parameters required by the transaction. Transaction codes are usually created during the configuration of a Twisp ledger and can be customized to meet the specific needs of your organization.

For this tutorial, we will use a `BOOK_TRANSFER` transaction code. The `BOOK_TRANSFER` transaction code requires the following parameters:

- the account to _debit_
- the account to _credit_
- the _amount_ to transfer
- the _currency_ of the transfer
- the _effective date_ of the transfer transaction

For a more in-depth look at this tran code, review the tutorial page on [](/tutorials/building-tran-codes).

Once a transaction code has been selected, a transaction can be written using the `postTransaction` mutation.

## Write a `postTransaction` mutation

To post a transaction, you will need to provide the following inputs:

- `transactionId`: a unique identifier for the transaction
- `tranCode`: the transaction code to use
- `params`: a set of key-value parameters as specified by the transaction code

Once the transaction is posted, Twisp will write the transaction to the ledger and create entries for each account affected by the transaction.

### Provide a unique `transactionId`

It's important to provide a unique `transactionId` when posting a transaction. This can help prevent duplicate transactions from being posted accidentally, i.e. it ensures idempotent transactions.

### Specify values for all `params`

Make sure to specify all required values in the `params` object for the transaction code you are using.

Here's an example mutation showing how to post a transaction with the `BOOK_TRANSFER` tran code:

{% partial file="graphql/reference/Mutation/006_postTransaction.md" /%}

This code example shows a transaction that transfers funds between two accounts. The `params` object specifies all required values for this transaction, including the customer account to credit (`crAccount`) and customer account to debit (`drAccount`), the amount to transfer (`amount`), the currency (`currency`), and the effective date of the transaction (`effective`).

Keep in mind that transactions are processed according to the rules defined in the transaction code. It is important to ensure that the transaction code and its associated parameters are correct before posting the transaction.

## Query transactions to see entries written

Any posted transaction can be queried to view its fields and the entries that were written to the ledger. Here's an example query:

{% partial file="graphql/reference/Query/012_transaction.md" /%}

This GraphQL query returns a transaction with its description, effective date, transaction code, journal name, metadata, and entries. The entries include the sequence, entry type, layer, direction, account name, and formatted amount.

## Conclusion

Posting transactions in Twisp is a straightforward process that requires specifying the transaction code to use along with the values for all parameters as defined by the tran code. The `postTransaction` mutation will return the posted transaction. Transactions can also be queried using the `transaction` and `transactions` queries.
