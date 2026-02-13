---
title: API Quickstart
description: Get up and running with the Twisp GraphQL API.
showNext: true
showPrev: true
draft: true
---

Let's start interacting with your ledger via the GraphQL API.

To keep things simple, we'll use the built-in GraphQL IDE in the console. Of course, you can make these same API calls from any authenticated client.

## Insert sample data

Because the database is still empty, to start exploring the query capabilities we first need to seed it with sample data.
### 1. Create an account record for a customer

The first record we'll insert is a sample account in the `accounts` table.

Paste the following GraphQL mutation into the query editor and execute the query with the "Execute Query" button (or press `Ctrl-Enter`). If the query executed correctly, you should see a response that looks similar to the JSON document below.

```graphql title="AddAccount"
mutation AddAccount {
  accounts {
    insert(input:{
      document:{
        account_id:"uuid('9dd984db-78d8-420f-9380-80e3cf36fe75')",
        customer_name:"'Giovanni Medici'"
      }
    }) {
      record {
        account_id
        balance
        customer_name
      }
    }
  }
}
```

```json title="Response"
{
  "data": {
    "accounts": {
      "insert": {
        "record": {
          "account_id": "9dd984db-78d8-420f-9380-80e3cf36fe75",
          "balance": 0,
          "customer_name": "Giovanni Medici"
        }
      }
    }
  }
}
```

{% callout type="note" %}

The values assigned to fields in the `document` parameter of the `insert` operation are of type `Expression`, which means that they represent a Common Expression Language (CEL) expression wrapped in double quotes `"`. This can be a bit confusing for first-time users, as they appear to be formatted as strings. You'll notice that the raw string `Giovanni Medici` was wrapped in single quotes `'` to indicate that it represents a raw UTF-8 string.

{% /callout %}

### 2. Insert sample transactions using variables

GraphQL operations also accept query variables, which you can define in JSON format in the "Query Variables" pane of the GraphQL interpreter. To expand the pane, click on the text "Query Variables".

The CEL evaluation engine is granted access to these variables inside of a `context.vars` object, which means that you can utilize the variable values anywhere in a GraphQL query that accepts an `Expression` type.

In this step, we'll create 3 transactions associated with the account created in step 1.

Define the query variables JSON object using the `account_id` for the account created in step 1 and a `posted_on` timestamp in RFC3339-compliant format.

```json title="Query Variables"
{
  "account_id":"9dd984db-78d8-420f-9380-80e3cf36fe75",
  "posted_on":"2022-04-27T10:50:00.000Z"
}
```

Now you can write the mutation to insert a `transaction` utilizing the values defined in the query variables. Note that the `uuid()` and `timestamp()` functions are used to coerce the values to the appropriate type.

```graphql title="InsertTransaction"
mutation InsertTransaction {
  transactions {
    insert(input:{
      document:{
        transaction_id:"uuid.New()",
        account_id:"uuid(context.vars.account_id)"
        amount:"125",
        posted_on:"timestamp(context.vars.posted_on)"
      }
    }) {
      record {
        transaction_id
        account_id
        amount
        posted_on
      }
    }
  }
}
```

```json title="Response"
{
  "data": {
    "transactions": {
      "insert": {
        "record": {
          "account_id":"9dd984db-78d8-420f-9380-80e3cf36fe75",
          "amount": 125,
          "posted_on": "2022-04-27T10:50:00Z",
          "transaction_id": "30f2283b-c44f-4781-93c8-a70c85bd5e16"
        }
      }
    }
  }
}
```

Execute this mutation and then insert two more transactions with amounts `50` and `75`.

{% callout type="tip" %}

While writing a GraphQL query or mutation in the query editor, you can use the autocomplete command (`Ctrl-Space`) to ensure that you are always writing a syntactically correct query. Additionally, the built-in GraphQL docs (button at top right of GraphQL pane) provides an in-editor reference to the full GraphQL schema.

{% /callout %}

## Execute a Query

Querying the data in your Twisp ledgers is done via GraphQL. Because Twisp generates a GraphQL schema custom-tailored to the ledger resources you declare, you get instant access to all the type-safety and developer ergonomics that comes with a GraphQL API.

In this step, we'll write some simple queries to read the data inserted in the previous page.

### 1. Query an account record using the ID index

The simplest query we can write is one that will fetch a single record from a ledger.

In the following example, we'll fetch the account record created in the previous page by querying the `account_id` index and providing the UUID for Giovanni's account.

Copy the below query into your GraphQL editor, replace the UUID with the account ID from your database, and execute the query. You should see a response with the account information.

```graphql title="GetAccount"
query GetAccount {
  accounts(
    index:account_id,
    where:{ partition:["9dd984db-78d8-420f-9380-80e3cf36fe75"] }
  ) {
    account_id
    balance
    customer_name
  }
}
```

```json title="Response"
{
  "data": {
    "accounts": [
      {
        "account_id": "9dd984db-78d8-420f-9380-80e3cf36fe75",
        "balance": 0,
        "customer_name": "Giovanni Medici"
      }
    ]
  }
}
```

{% callout type="note" %}

Every query must declare **which index to query against**. The name of the index must match one of the indices defined on the ledger database.

{% /callout %}

### 2. Expand the query to view customer's transactions

One of the most powerful features of GraphQL is its ability to combine related records into a single, nested document. Because we declared a one-to-many relationship between the `accounts` and `transactions` ledgers using the `Join` property in the template file, within a query against `accounts` we can also fetch the related `transactions` as a nested document.

Edit your `GetAccount` query to return the `transactions` field for an account, returning the `amount` and `posted_on` fields.

When you execute your query, you should see a response with all 3 transactions associated with the account.

```graphql title="GetAccount"
query GetAccount {
  accounts(
    index:account_id,
    where:{ partition:["9dd984db-78d8-420f-9380-80e3cf36fe75"] }
  ) {
    account_id
    balance
    customer_name
    transactions {
      amount
      posted_on
    }
  }
}
```

```json title="Response"
{
  "data": {
    "accounts": [
      {
        "account_id": "9dd984db-78d8-420f-9380-80e3cf36fe75",
        "balance": 0,
        "customer_name": "Giovanni Medici",
        "transactions": [
          {
            "amount": 50,
            "posted_on": "2022-04-27T10:50:00Z"
          },
          {
            "amount": 75,
            "posted_on": "2022-04-27T10:50:00Z"
          },
          {
            "amount": 125,
            "posted_on": "2022-04-27T10:50:00Z"
          }
        ]
      }
    ]
  }
}
```
## Explore record history

{% comment %}

TODO: add link to page "Append-Only Ledgers" when it is written

{% /comment %}

Ledgers in Twisp are append-only, which means that any changes to a record do not overwrite existing data. The API allows for updates to existing records, but the previous values are stored as a change log.

This approach satisfies the immutability constraint required for accurate and trustworthy storage of financial data while still providing a full CRUD interface for ease of development.

Let's "update" a record and inspect the change sequence made available through the `history` field available on every ledger record.

### 1. Change a transaction and query the change history

For the sake of this demonstration, let us assume that one of the transactions was initially inserted with the wrong `posted_on` date.

Pick a transaction ID to use and update the timestamp using the following GraphQL mutation.

{% callout type="tip" %}

You can use the `GetAccount` query above to fetch the `transaction_id` for transactions associated with the singular account.

{% /callout %}

```graphql title="UpdateTransaction"
mutation UpdateTransaction {
  transactions {
    update(input:{
      index:transaction_id,
      where:{ partition:["<transaction ID>"] },
      document:{
        posted_on:"timestamp('2022-04-28T12:10:00.000Z')"
      }
    }) {
      records {
        transaction_id
        amount
        posted_on
      }
    }
  }
}
```

```json title="Response"
{
  "data": {
    "transactions": {
      "update": {
        "records": [
          {
            "amount": 125,
            "posted_on": "2022-04-28T12:10:00Z",
            "transaction_id": "<transaction id>"
          }
        ]
      }
    }
  }
}
```

With that update made, you can see the history of changes by querying the `history` field of the transaction.

```graphql title="GetTransaction"
query GetTransaction {
  transactions(
    index:transaction_id,
    where:{ partition: ["<transaction ID>"] }
  ) {
    transaction_id
    posted_on
    history {
      seq
      posted_on
    }
  }
}
```

```json title="Response"
{
  "data": {
    "transactions": [
      {
        "history": [
          {
            "posted_on": "2022-04-28T12:10:00Z",
            "seq": 2
          },
          {
            "posted_on": "2022-04-27T10:50:00Z",
            "seq": 1
          }
        ],
        "posted_on": "2022-04-28T12:10:00Z",
        "transaction_id": "<transaction ID>"
      }
    ]
  }
}
```

{% callout type="note" %}

The fields available to query within `history` are the same fields as the `transaction` schema, with the addition of a `seq` (sequence) field. The `seq` field is an auto-incrementing integer showing the sequence of changes that have been applied to this record.

{% /callout %}

### 2. Query the change history of the account balance

The `history` field is useful for querying changes made through GraphQL mutations as above. Because it tracks *all* changes to a record, we can also use it to inspect the changes made through automatic operations defined by `Triggers`.

The trigger defined in the previous page declared a `DO_UPDATE` operation against the `accounts` table. With each new transaction inserted, the trigger updated the account balance.

This means that we can query the `history` of the account record to see how the balance changed with each transaction added:

```graphql title="GetAccount"
query GetAccount {
  accounts(
    index:account_id,
    where:{ partition:["9dd984db-78d8-420f-9380-80e3cf36fe75"] }
  ) {
    account_id
    balance
    transactions {
      amount
    }
    history {
      balance
      seq
    }
  }
}
```

```json title="Response"
{
  "data": {
    "accounts": [
      {
        "account_id": "9dd984db-78d8-420f-9380-80e3cf36fe75",
        "balance": 250,
        "history": [
          {
            "balance": 250,
            "seq": 4
          },
          {
            "balance": 175,
            "seq": 3
          },
          {
            "balance": 125,
            "seq": 2
          },
          {
            "balance": 0,
            "seq": 1
          }
        ],
        "transactions": [
          {
            "amount": 125
          },
          {
            "amount": 50
          },
          {
            "amount": 75
          }
        ]
      }
    ]
  }
}
```

{% callout type="note" %}

The `history` shows changes in reverse chronological order, while `transactions` are returned in the order in which they were inserted.

{% /callout %}

---

🎉 Congratulations! 🎉 

You've completed all of the steps in the Quickstart. You are ready to go build with your Twisp leger.

We have only scratched the surface of what Twisp is capable of in these pages, however. To continue learning, consider reading the other pages in the docs.