---
title: Response Format
metaTitle: "Reference : API Response Format"
description: Requests to the Twisp GraphQL API return a consistent JSON response format.
metaDescription: Response format reference for the Twisp GraphQL API.
showNext: true
showPrev: true
---

The API uses a response format conforming to the [GraphQL spec](https://spec.graphql.org/October2021/#sec-Response-Format).

```
{
  "data": {
    // GraphQL response
  }
}
```

This means that every response is a JSON document with a `"data"` field. If the request produced an error, then it will also return an `"errors"` field.

> The `data` entry in the response will be the result of the execution of the requested operation. If the operation was a query, this output will be an object of the query root operation type; if the operation was a mutation, this output will be an object of the mutation root operation type.

> The `errors` entry in the response is a non-empty list of errors, where each error is a map.

Because the Twisp API supports [](/reference/api/transactional-operations), if there are _any_ errors, then the entire operation is aborted and the `"data"` field will be `null`.

```
{
  "errors": [
    // list of errors
  ],
  "data": null
}
```

In addition, responses may contain an `"extensions"` field if certain request headers are specified. Learn more about extensions in the [](/reference/api/extensions) reference.
