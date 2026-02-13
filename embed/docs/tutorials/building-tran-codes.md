---
title: Building Tran Codes
metaTitle: "Tutorial: Building Tran Codes"
description: |
  In this tutorial, we will learn how to manage transaction codes using the GraphQL API.
showNext: true
showPrev: true
---

Transaction codes play a critical role in identifying and categorizing transactions in the ledger, and designing them is an important part of using Twisp. By following the steps outlined in this tutorial, you will be able to create and manage transaction codes in your own ledger.

{% callout type="task" %}
- Create a new tran code using the `createTranCode` mutation
- Update fields on an existing tran code using the `updateTranCode` mutation
- Get data about a tran code with the `tranCode` query
- Delete (lock) a tran code using the `deleteTranCode` mutation
{% /callout %}

---

## Prerequisites

Before you start, you should have added accounts to your ledger. See the tutorial on [](/tutorials/setting-up-accounts).

{% partial file="tutorials_getting_started.md" /%}

## Create a new tran code

To create a new tran code, use the `createTranCode` mutation. The `input` argument specifies the full configuration of the tran code and how transactions posted with this code will be structured, including all entries written to the ledger.

Many of the fields in the `TranCodeInput` object (and its nested objects) allow for CEL expressions to be used, which will be evaluated during the `postTransaction` call and the results will be used in the creation of a transaction posted with the new transaction code.

Note that the `params` field allows for the specification of parameters that can be used in CEL expressions for the transaction and entries fields. These parameters can be used to make the transaction code more flexible and dynamic, allowing for the creation of more complex transactions.

The `transaction` and `entries` fields, respectively, act as templates for the [Transaction]($gql:object:Transaction) and [Entries]($gql:object:Entry) that are created when this tran code is used in a `postTransaction` call.

Here's an example of creating a new tran code for a `BOOK_TRANSFER`:

{% partial file="graphql/reference/Mutation/004_createTranCode.md" variables={title: "Create tran code"} /%}

Once created, transactions can be posted using this transaction code by providing the required fields in a `TransactionInput` object. The inputs for a `postTransaction` using this tran code must include fields for the `crAccount`, `drAccount`, `amount`, `currency`, and `effective`.

{% callout type="note" %}
Designing tran codes is one of the most important parts of using Twisp, and it can take some time to familiarize yourself with the process.

Read more about tran codes on the [](/accounting-core/encoded-transactions) page.
{% /callout %}

## Query tran codes

The `tranCode` query can be used to read back an existing tran code using its `tranCodeId`.

Here is an example query:

{% partial file="graphql/reference/Query/010_tranCode_BOOK_TRANSFER.md" variables={title: "Read tran code"} /%}

This query will respond with the fields for the `BOOK_TRANSFER` tran code created above, identified by its UUID.

The `params` field in the response includes information about the parameters required by this transaction code, such as the name, type, and description of each parameter.

The `transaction` field in the response includes details about the transaction that is created when this transaction code is used.

The `entries` field in the response includes information about the entries that are created when this transaction code is used.

Finally, the response includes the status of the transaction code, which can be either `ACTIVE` or `LOCKED`, and any metadata associated with the transaction code.

## Modify an existing tran code

To modify an existing tran code, use the `updateTranCode` mutation.

Provide the tran code's `id` and the fields to update in the `TranCodeUpdateInput` input object. You can only modify a subset of fields for data integrity purposes.

Here's an example of modifying the description of an existing tran code:

{% partial file="graphql/reference/Mutation/010_updateTranCode.md" variables={title: "Update tran code"} /%}

This operation updates the description of a transaction code. The response includes the UUID and the updated description of the transaction code, as well as the history of the changes made to the transaction code.

## Lock a tran code

Because Twisp is an immutable database, we cannot fully "delete" a tran code. Instead, Twisp marks the tran code's status as `LOCKED`, which prevents transactions from posting using this version of tran code.

To delete (lock) an tran code, we can use the `deleteTranCode` mutation. We need to provide the `id` of the tran code we want to delete. Here's an example mutation:

{% partial file="graphql/reference/Mutation/015_deleteTranCode.md" variables={title: "Delete tran code"} /%}

Once the tran code is deleted, this mutation returns two fields: `tranCodeId` and `status`. `tranCodeId` is a UUID that uniquely identifies the deleted tran code, while `status` is the current status of the tran code (i.e. `LOCKED`).

This mutation is useful if you need to remove an tran code that is no longer needed or was created in error.

## Conclusion

This tutorial covered the basics of creating, modifying, and deleting transaction codes. Transaction codes play a critical role in identifying and categorizing transactions in the ledger, and designing them is an important part of using Twisp. By following the steps outlined in this tutorial, you should be able to create and manage transaction codes in your own ledger.
