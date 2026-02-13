---
title: Pagination
metaTitle: "Reference : API Pagination"
description: |
  When making queries which can potentially return high-cardinality lists, the Twisp API supports cursor-based pagination.
metaDescription: Pagination reference for the Twisp GraphQL API.
showNext: true
showPrev: true
---

The Twisp GraphQL API implements a cursor-based pagination model using the [Relay GraphQL Cursor Connections Specification](https://relay.dev/graphql/connections.htm). This specification provides guidelines for making paginated requests using a concept called "Connections." Connections facilitate the process of fetching data in chunks (or pages) from a GraphQL API.

{% callout type="note" %}
For a step-by-step guide on using pagination in a GraphQL operation, see the [](/tutorials/api/querying-paginated-fields) tutorial.
{% /callout %}

## Heuristics for Paginated Queries

Within the GraphQL schema, pluralized query fields (e.g. `Query.entries`, `Query.tranCodes`) return Connection types. No matter what type of record we query, the same pattern can be applied. You can use these basic heuristics across the schema:

- When querying a field with a `*Connection` type response, you must indicate the number of records to return with the `first` argument.
- To determine if there are additional records beyond the current page, query the `pageInfo` object.
- To get the next page of `n` records, use the cursor provided `pageInfo.endCursor` as the `after` argument.
- When there are no more records, `pageInfo.hasNextPage` will be false.

This is an abstracted example of an imaginary `widgets` query to return `Widget` records:

```graphql
query {
  widgets(
    after: "<cursor>" # Cursor indicating the start point.
    first: 5          # Number of nodes (widgets) to return.
  ) {
    __typename        # => WidgetConnection
    edges {
      __typename      # => WidgetConnectionEdge
      cursor          # Cursor position of the current edge.
      node {          # The Widget record at this edge position.
        # ...         # Field queries on the Widget record.
      }
    }
    nodes {           # Alternate access to all `node` objects within `edges`
      __typename      # => Widget
      # ...           # Field queries on the Widget record.
    }
    pageInfo {
      hasPreviousPage # True if there are nodes in the connection before the current page / start cursor.
      hasNextPage     # True if there are nodes in the connection after the current page / end cursor.
      startCursor     # Query cursor for the first node in the current page.
      endCursor       # Query cursor for the last node in the current page.
    }
  }
}
```

The rest of this reference page will go into the specifics of cursor-based pagination using the `accounts` query as an example, but the same principles can be applied to every query resolving to a `*Connection` type.

## Querying Lists with the First Argument

When making queries to fields which resolve to paginated lists of records, the required `first` argument is used to specify the _maximum_ number of records to return.

For example, when using the `journals` query, we can request only the first 5 journals using the following query. We query the `nodes` field to get fields from the list of records returned.

{% partial file="graphql/reference/api-basics/pagination/001_FirstNoCursor.md" variables={title: "Query the first 5 accounts"} /%}

{% callout type="note" %}

This example assumes a ledger structure where journals follow the `code` format of `CUST.<customer name>`. Note also that the index used is automatically alphabetically sorted.  
{% /callout %}

## Connections use PageInfo to Enable Pagination

To query the nest page of records, we need to specify the **cursor position** to begin at by providing a value to the `after` argument of the query field. This tells the query to start at the cursor position and return the subsequent `n` records, where `n` is the number specified by the `first` argument.

The `pageInfo` field returns a [PageInfo]($gql:object:PageInfo) object with the information needed to perform this kind of cursor-based pagination. It will tell us whether there are records (pages) **before** the `startCursor` of the current page and whether there are records (pages) **after** the `endCursor`.

{% partial file="graphql/reference/api-basics/pagination/002_WithPageInfo.md" /%}

With this set of information, we have everything we need to write paginated queries.

## Querying with Cursors

To make subsequent queries for additional pages of records, we use the cursor provided at `pageInfo.endCursor` from the current page as the `after` argument when querying the subsequent page.

To simplify this procedure, let's imagine that there are 13 records that can be returned by a particular query, and we want to paginate using pages of 5 records each. We would need to perform three query operations:

1. To get the first page of 5 records, use `first: 5` while omitting the `after` argument to specify that we want to start our query at the beginning of the list. Query the `pageInfo.endCursor`.
2. To get the next page, use `first: 5` and `after: X` where `X` is the cursor returned in the first page's `pageInfo.endCursor`.
3. To get the final page, use `first: 5` and `after: Y` where `Y` is the cursor returned in the second page's `pageInfo.endCursor`.

We can also render this 3-page query sequence as a table:

| Page | After Cursor | Records (Nodes) | End Cursor                  | Next Page? |
|------|--------------|-----------------|-----------------------------|------------|
| `1`  | `null`       | `1..5`          | `5b11e58c` (for record #5)  | `TRUE`     |
| `2`  | `5b11e58c`   | `6..10`         | `1b47d1e9` (for record #10) | `TRUE`     |
| `3`  | `1b47d1e9`   | `11..13`        | `092e5914` (for record #13) | `FALSE`    |

Following up on the previous example using the `journals` query, we can use the `pageInfo.endCursor` from the first response as the `after` argument to get the next 5 journals:

{% partial file="graphql/reference/api-basics/pagination/003_WithAfterCursor.md" /%}
## Nodes and Edges

Technically, every `*Connection` type returns a set of `edges` where each edge is identified by its `cursor` and contains a reference to a `node` which contains the record being queried (in this example, an [Account]($gql:object:Account) object).

The `nodes` field used in the above is simply a shorthand way of accessing all the `node` fields of every edge in the `edges` list. Similarly, the `startCursor` and `endCursor` fields from `pageInfo` are just shorthand ways to access the `cursor` of the first and last edge in the `edges` list, respectively.

In the example below, the `journals` query returns an [JournalConnection]($gql:object:JournalConnection) type with an `edges` field containing a list of [JournalConnectionEdge]($gql:object:JournalConnectionEdge) objects. Each [JournalConnectionEdge]($gql:object:JournalConnectionEdge) object has a `cursor` position and a `node` referencing the [Journal]($gql:object:Journal) record.

{% partial file="graphql/reference/api-basics/pagination/004_EdgesWithCursor.md" /%}

For most pagination queries, only the `nodes` field is needed. In the various examples shown in Twisp's documentation, we will usually omit `edges` and just query `nodes` to keep the response smaller and easier to read.

## Why Cursor-Based Pagination?

The GraphQL core team has spent a lot more time on these questions than we have, and [they recommend opaque cursors as the best approach](https://graphql.org/learn/pagination/#pagination-and-edges).

> In general, we've found that cursor-based pagination is the most powerful of those designed.

There are other models out there, but for the GraphQL protocol this method provides the greatest combination of flexibility, specificity, and performance.
