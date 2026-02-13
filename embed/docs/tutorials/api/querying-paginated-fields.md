---
title: Querying Paginated Fields
metaTitle: "Tutorial : Querying Paginated Fields"
description: In this tutorial, we will learn to handle cursor-paginated fields in the GraphQL schema.
showNext: false
showPrev: false
---

{% partial file="tutorials_getting_started.md" /%}

---

## 1. Identify the connection field

Any field that resolves to a `*Connection` type uses cursor-based pagination.

This includes top-level queries like `entries` and `transactions`, as well as fields within object types like `Account.balances`.

## 2. Request edges and nodes

To fetch data in pages, request the `edges` in the connection field. Each edge represents a link to a data node and may contain additional information about the relationship. Within each edge, request the `node` field to retrieve the actual data object. Here's a sample query for fetching the first 10 items of a field called `someConnection`:

```graphql
query {
  someConnection(first: 10) {
    edges {
      node {
        id
        name
      }
    }
  }
}
```

{% callout %}
Connection types also contain a `nodes` field, which is just a way to more directly access the list of `node` objects in `edges`. It can be useful in cases where you don't need to query any other fields of the edge object.
{% /callout %}

## 3. Get the page info

Alongside the edges, request the `pageInfo` object, which contains information about pagination. `pageInfo` includes the fields `hasNextPage`, `hasPreviousPage`, `startCursor`, and `endCursor`. Add the `pageInfo` field to your query:

```graphql
query {
  someConnection(first: 10) {
    edges {
      node {
        id
        name
      }
    }
    pageInfo {
      hasNextPage
      endCursor
    }
  }
}
```

## 4. Paginate using cursors

Cursors are opaque strings representing the position of an item in the list. Use the `after` argument in conjunction with the `first` argument to request the next set of _n_ items after the `endCursor`:

```graphql
query {
  someConnection(first: 10, after: "<end-cursor>") {
    edges {
      node {
        id
        name
      }
    }
    pageInfo {
      hasNextPage
      endCursor
    }
  }
}
```

Replace `"<end-cursor>"` with the actual `endCursor` value from the previous `pageInfo`.

## 5. Iterate through pages

Continue making requests using the updated `endCursor` until the `hasNextPage` field in `pageInfo` is `false`, indicating that there are no more pages to fetch.
