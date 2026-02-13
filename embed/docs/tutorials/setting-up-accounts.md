---
title: Setting Up Accounts
metaTitle: "Tutorial: Setting Up Accounts"
description: |
  In this tutorial, we will learn how to manage accounts using the GraphQL API.
showNext: true
showPrev: true
---

Managing financial accounts is a crucial aspect of many applications. By the end of this tutorial, you will have a solid understanding of how to manage accounts using the Twisp GraphQL API and be ready to incorporate this functionality into your own applications.

{% callout type="task" %}
- Create a new account using the `createAccount` mutation
- Update fields on an existing account using the `updateAccount` mutation
- Delete (lock) an account using the `deleteAccount` mutation
{% /callout %}

---

{% partial file="tutorials_getting_started.md" /%}

## Create an account

The `createAccount` mutation is used to create new accounts with the given `accountId`, `name`, `code`, `description`, `normalBalanceType`, and `status`.

The `accountId` is a unique identifier for the account that is being created, and it is required as an input field for this mutation. The `name` input field specifies the name for the new account that is being created. The `code` input field is a shorthand code for the account, which is often an abbreviated version of the account name. The `description` input field is a brief explanation of the account that is being created.

The `status` input field represents the current status for the new account that is being created. This field specifies whether the account is active or closed (locked). By default, all accounts are `ACTIVE`. When an account is `LOCKED`, it cannot be updated unless you are also changing its status back to `ACTIVE`. Any attempt to write a ledger entry to a locked account will raise an error.

{% partial file="graphql/reference/Mutation/002_createAccount.md" /%}

This `createAccount` operation returns the newly created account's `accountId`, `name`, `code`, `description`, and `normalBalanceType`.

## Update an account

To update an existing account, we can use the `updateAccount` mutation. We need to provide the `id` of the account we want to update, as well as the fields we want to update. Here's an example mutation:

{% partial file="graphql/reference/Mutation/007_updateAccount.md" /%}

This is a GraphQL mutation operation that allows the user to update an account's code. The mutation is called `updateAccount` and takes in two arguments: `id` and `input`. The `id` argument is a unique identifier for the account being updated, and the `input` argument contains the updates to apply to the account. In this specific example, the `id` argument is passed in with the `$accountCardSettlementId` variable.

The response from this mutation includes several fields. The `accountId` field returns the unique identifier for the updated account, while the `code` field returns the updated code for the account. Additionally, the `history` field returns the two most recent versions of the account, along with their corresponding `version` and `code`. This allows the user to track the changes made to the account over time.

Overall, this mutation provides a flexible and powerful way to update accounts in the Twisp system, and the response includes valuable information that can be used to track changes and ensure data integrity.

## Delete an account

Because Twisp is an immutable database, we cannot fully "delete" an account. Deleting an account in Twisp instead marks its status as `LOCKED`, meaning that no entries can be posted to the account. A locked account can only be updated to change its status back to `ACTIVE`.

To delete (lock) an account, we can use the `deleteAccount` mutation. We need to provide the `id` of the account we want to delete. Here's an example mutation:

{% partial file="graphql/reference/Mutation/013_deleteAccount.md" /%}

Once the account is deleted, this mutation returns two fields: `accountId` and `status`. `accountId` is a UUID that uniquely identifies the deleted account, while `status` is the current status of the account.

This mutation is useful if you need to remove an account that is no longer needed or was created in error. By deleting the account, you can ensure that it is no longer used in any future transactions or reports.

## Conclusion

In this tutorial, we explored how to create, update, and delete accounts using the Twisp GraphQL API. We saw how to use the `createAccount`, `updateAccount`, and `deleteAccount` mutations to perform these operations. With this knowledge, you can now begin building your own applications that use the Twisp API to manage financial accounts.
