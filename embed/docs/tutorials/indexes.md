---
title: Creating Custom Indexes
metaTitle: "Tutorial: Creating Custom Indexes"
description: |
  Learn how to create and use custom indexes in Twisp to efficiently query records based on specific fields, including data within the metadata object.
showNext: true
showPrev: true
---

Custom indexes allow you to define specific ways to query your ledger data, enabling efficient filtering and sorting based on fields like `accountId`, `status`, or even nested values within the `metadata` object.

By the end of this tutorial, you will be able to create a custom index tailored to your querying needs and use it to retrieve records efficiently.

{% callout type="task" %}
- Design a custom index with partition and sort keys using CEL.
- Create the index using the `schema.createIndex` mutation.
- Query records efficiently using the custom index.
{% /callout %}

---

{% partial file="tutorials_getting_started.md" /%}

## Designing a Custom Index

Before creating an index, you need to decide how you want to query your data. Let's design an index on the `Account` table to quickly find accounts based on their `accountId` and a `category` field within their `metadata`.

Key components of an index definition:

*   **`on`:** The record type (table) the index applies to (e.g., `Account`).
*   **`name`:** A unique, human-readable name for the index (e.g., `account_metadata_category`).
*   **`partition`:** How data is grouped. Queries *must* filter on the partition key(s). Good partitioning spreads data evenly. We'll partition by `accountId`.
    *   `alias`: A name for the key (e.g., `accountId`).
    *   `value`: A [CEL expression](https://twisp.com/docs/reference/cel) referencing the `document` (the record being indexed) to get the partition value (e.g., `document.accountId`).
*   **`sort` (Range Key):** Defines the order *within* a partition, enabling range queries (`gte`, `lte`) and sorting. We'll sort by `metadata.category`.
    *   `alias`: A name for the key (e.g., `category`).
    *   `value`: A CEL expression for the sort value (e.g., `string(document.metadata.category)` - casting to string ensures predictable sorting).
    *   `sort`: `ASC` or `DESC`.
*   **`constraints` (Optional):** CEL expressions that must *all* be true for a record to be included. We'll only index accounts that *have* a `metadata.category`.
    *   Example: `{ hasCategory: "has(document.metadata.category)" }`
*   **`unique` (Optional):** If `true`, ensures the combination of partition and sort keys is unique for each indexed record. Defaults to `false`.

For our example:
*   **Name:** `account_metadata_category`
*   **On:** `Account`
*   **Partition:** `alias: "accountId"`, `value: "document.account_id"`
*   **Sort:** `alias: "category"`, `value: "string(document.metadata.category)"`, `sort: ASC`
*   **Constraints:** `{ hasCategory: "has(document.metadata.category)" }`
*   **Unique:** `false`

## Create the Custom Index

Use the `schema.createIndex` mutation to create the index defined above.

```graphql
mutation CreateAccountMetadataIndex {
  schema {
    createIndex(
      input: {
        name: "account_metadata_category"
        on: Account
        unique: false
        partition: [
          { alias: "accountId", value: "document.account_id" }
        ]
        sort: [
          {
            alias: "category"
            value: "string(document.metadata.category)"
            sort: ASC
          }
        ]
        constraints: { hasCategory: "has(document.metadata.category)" }
      }
    ) {
      name
      on
      unique
      partition {
        alias
        value
      }
      range { # 'range' is the field for sort keys in the response
        alias
        value
        sort
      }
      constraints
      historical
      search
    }
  }
}
```

```json
{
  "data": {
    "schema": {
      "createIndex": {
        "name": "account_metadata_category",
        "on": "Account",
        "unique": false,
        "partition": [
          {
            "alias": "accountId",
            "value": "document.account_id"
          }
        ],
        "range": [
          {
            "alias": "category",
            "value": "string(document.metadata.category)",
            "sort": "ASC"
          }
        ],
        "constraints": {
          "hasCategory": "has(document.metadata.category)"
        },
        "historical": false,
        "search": false
      }
    }
  }
}
```

This mutation creates the `account_metadata_category` index on the `Account` table. The response confirms the index structure, including its partition and range (sort) keys.

{% callout type="note" %}
Twisp also supports `createHistoricalIndex` to index every version of a record and `createSearchIndex` for eventually consistent full-text search capabilities.
{% /callout %}

## Query Using the Custom Index

To use your new index, specify `index: { name: CUSTOM }` in your query and provide the index name and filters in the `where.custom` argument.

You **must** provide an equality (`eq`) filter for all partition key aliases. You can optionally provide filters (`eq`, `gte`, `lte`, `prefix`, etc.) for sort key aliases.

```graphql
query QueryUsingCustomIndex {
  accounts(
    index: { name: CUSTOM }
    where: {
      custom: {
        index: "account_metadata_category" # Name of the custom index
        partition: [
          {
            alias: "accountId" # Partition key alias
            value: { eq: "a1b2c3d4-e5f6-7890-1234-567890abcdef" }
          }
        ]
        sort: [
          {
            alias: "category" # Sort key alias
            value: { eq: "Premium" } # Filter condition on sort key
          }
        ]
      }
    }
    first: 10
  ) {
    nodes {
      accountId
      name
      metadata
    }
    pageInfo {
      hasNextPage
      endCursor
    }
  }
}
```

```json
{
  "data": {
    "accounts": {
      "nodes": [
        {
          "accountId": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
          "name": "Premium Customer Account",
          "metadata": {
            "category": "Premium",
            "region": "US-West"
          }
        }
        // ... other matching accounts up to 10
      ],
      "pageInfo": {
        "hasNextPage": false,
        "endCursor": "..."
      }
    }
  }
}
```

This query efficiently retrieves `Account` records using the `account_metadata_category` index, filtering first by the `accountId` partition and then by the `category` sort key.

## Conclusion

In this tutorial, you learned how to design and create a custom index using the `schema.createIndex` mutation. You saw how to define partition keys, sort keys, and constraints using CEL expressions. Finally, you learned how to leverage your custom index in queries for efficient data retrieval based on specific fields, including nested metadata values.

Custom indexes are a powerful tool for optimizing query performance in Twisp. Consider creating them for your common query patterns.
