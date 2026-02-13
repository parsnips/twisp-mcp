---
title: Pulling Balances
metaTitle: "Tutorial: Pulling Balances"
description: |
  In this tutorial, we will learn how to query account balances using the GraphQL API.
showNext: true
showPrev: true
---

Twisp provides a robust GraphQL schema that allows users to query account balances using a variety of filters and indexes.

{% callout type="task" %}
- Get a single balance for a specific account and currency with the `balance` query
- Get a set of balances for a given set of conditions with the `balances` query
- Pull the balance for an account with the `Account.balance` field
- Pull an aggregate balance for accounts in a set with the `AccountSet.balance` field
{% /callout %}

---

## Prerequisites

Before you start, you should have added accounts to your ledger and posted some transactions. See the tutorials on [](/tutorials/setting-up-accounts) and [](/tutorials/posting-transactions).

{% partial file="tutorials_getting_started.md" /%}

### Retrieve an account UUID

First, you will need to retrieve the UUID for the account that you want to query. You can do this by querying for the account using the `accounts` query, and specifying any relevant filters.

For example, to retrieve the UUID for an account with the code `CHECKING`, you can use the following GraphQL query:

```graphql
query {
  accounts(
    index: { name: CODE }
    where: { code: { eq: "CHECKING" } }
    first: 1
  ) {
    nodes {
      accountId
      name
    }
  }
}
```

This query will return the UUID and name for any accounts with the code `CHECKING`.

## Query the account balance

Once you have the UUID for the account that you want to query, you can use the `balance` query to retrieve the balance for that account.

The `balance` query takes two arguments:

- `journalId` (required): the UUID of the journal to retrieve the balance from.
- `currency` (required): the currency code of the balance to retrieve.

For example, to retrieve the settled normal balance for the account with UUID `"3ea12e45-7df2-4293-9434-feb792affc91"` in journal with UUID `"4e9d9f6c-0cc3-4a8e-9d1a-0de5a47b0c8b"`, you can use the following GraphQL query:

```graphql
query {
  balance(
    accountId: "3ea12e45-7df2-4293-9434-feb792affc91"
    journalId: "4e9d9f6c-0cc3-4a8e-9d1a-0de5a47b0c8b"
    currency: "USD"
  ) {
    settled {
      normalBalance {
        units
      }
    }
  }
}
```

This query will return the USD settled normal balance amount in decimal units for for the account given in the specified journal.

## Query multiple balances

If you want to retrieve multiple balances at once, you can use the `balances` query instead. The `balances` query takes the same arguments as the `balance` query, but also allows you to specify additional filters to narrow down the results.

For example, to retrieve all balances for the account with UUID `"3ea12e45-7df2-4293-9434-feb792affc91"` in currency `"USD"`, you can use the following GraphQL query:

```graphql
query {
  balances(
    index: { name: ACCOUNT_ID }
    where: { accountId: { eq: "3ea12e45-7df2-4293-9434-feb792affc91" } }
    first: 5
  ) {
    nodes {
      journalId
      currency
      settled {
        normalBalance {
          units
        }
      }
    }
  }
}
```

This query retrieves information about the balances of a specific account. It requests the first five balances for the account with an `accountId` of `"3ea12e45-7df2-4293-9434-feb792affc91"`, and returns the `journalId`, `currency`, and `normalBalance` for each balance on the account.

## Query a balance using the `Account.balance` field

To query an account balance using the `balance` field on an `Account` type, you will need to provide the `journalId` and `currency` arguments. The `journalId` argument specifies the ID of the journal for the balance, and the `currency` argument specifies the currency of the balance.

For example, to retrieve the USD settled normal balance for the account with UUID `"3ea12e45-7df2-4293-9434-feb792affc91"` in journal with UUID `"4e9d9f6c-0cc3-4a8e-9d1a-0de5a47b0c8b"`, you can use the following GraphQL query:

```graphql
query {
  account(id: "3ea12e45-7df2-4293-9434-feb792affc91") {
    balance(journalId: "4e9d9f6c-0cc3-4a8e-9d1a-0de5a47b0c8b", currency: "USD") {
      settled {
        normalBalance {
          units
        }
      }
    }
  }
}
```

Note that this is effectively the same as using the `balance` query above: it will return the account's USD settled normal balance amount in decimal units in the specified journal.

{% callout %}
You can also query for multiple balances on accounts using the `Account.balances` field. This can be useful for accounts that contain entries across multiple journals, or for accounts that contain multiple currencies.
{% /callout %}

## Query a balance for an account set

To query balances for an account set, you can use the `balances` query with the `accountSetId` argument. Balances for an account set are calculated by summing the journal balances of all accounts that are members of the account set.

In this context, the journal balance for an account is the balance for all entries posted to the journal that the account set is tied to. So an account set will only show balances for entries posted to the account set's journal.

For example, to retrieve the USD balance for an account set with UUID `"d1a2e7e9-7a6d-4d82-bf6c-2c8d91b1c1f6"`, you can use the following GraphQL query:

```graphql
query {
  accountSet(id: "d1a2e7e9-7a6d-4d82-bf6c-2c8d91b1c1f6") {
    balance(currency: "USD") {
      settled {
        normalBalance {
          units
        }
      }
    }
  }
}
```

This query retrieves the USD balance of the account set with UUID `"d1a2e7e9-7a6d-4d82-bf6c-2c8d91b1c1f6"`. The `currency` argument in the balance field specifies the currency to use for the balance calculation. The query includes the `normalBalance` field of the `settled` balance layer, which returns the decimal units for the normal balance of the account set.

{% callout %}
You can also query for multiple balances on account sets using the `AccountSet.balances` field. This can be useful for account sets that track accounts across multiple currencies.
{% /callout %}

## Conclusion

Twisp's GraphQL schema provides a powerful set of tools for querying account balances. By leveraging the `balance` and `balances` queries, you can quickly and easily retrieve the balance information you need for your ledger entries and financial reporting.
