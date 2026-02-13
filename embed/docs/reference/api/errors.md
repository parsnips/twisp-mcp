---
title: Errors
metaTitle: "Reference : API Errors"
description: Error codes returned by the Twisp API when things don't go right.
metaDescription: Errors reference for the Twisp GraphQL API.
showNext: true
showPrev: true
---

This reference lists the types of errors that may be encountered.

## Error Responses

The API uses a response format conforming to the [GraphQL spec](https://spec.graphql.org/October2021/#sec-Response-Format). For more information, see the docs on [](/reference/api/response-format).

> If the request raised any errors, the response map must contain an entry with key `errors`. The value of this entry is described in the “Errors” section. If the request completed without raising any errors, this entry must not be present.

Because the Twisp API supports [](/reference/api/transactional-operations), if there are _any_ errors, then the entire operation is aborted and the `"data"` field will be `null`.

## Example Response

This example shows the response of an `createAccount` operation with a malformed `accountId`.

{% partial file="graphql/fixtures/errors/UUID_PARSE_ERROR.md" variables={title: "", defaultTab: "Response"} /%}

The response indicates that there was an error with the `createAccount` mutation request. The error is a `UUID_PARSE_ERROR`, which means that the `accountId` field in the request has an invalid UUID value.

Within the `errors` array, each error object contains fields specifying more information about the error:

- `message` summarizes the error, and may give a clue to its cause. In this example, it indicates that the `accountId` field has an invalid length, which is 0.
- `path` specifies the location within the GraphQL operation that the error occurred. In this example, it indicates that the error occurred in the `accountId` field of the `input` argument provided to the `createAccount` mutation.
- `extensions` provides additional information about the error.
- `extensions.code` specifies the **error code** used to represent the error. In this example, the code is `UUID_PARSE_ERROR`.

The `data` field in the response is `null` since the request was not successful due to the error.

To resolve this particular error, a valid UUID value should be provided for the `accountId` field. An empty string is not a valid UUID value.

## Error Codes

### ACCESS_DENIED
Indicates that the request was denied due to a lack of proper authentication or authorization.

### ALREADY_EXISTS
Indicates that the requested resource already exists and cannot be created again.

### BAD_REQUEST
Indicates that the request was malformed or invalid.

For example, a missing `eq` in the partition key:
{% partial file="graphql/fixtures/errors/BAD_REQUEST.eqRequiredForPartitionKey.md" variables={title: "", defaultTab: "Response"} /%}

Or a missing partition key:
{% partial file="graphql/fixtures/errors/BAD_REQUEST.partitionKeyRequired.md" variables={title: "", defaultTab: "Response"} /%}

### CEL_EVALUATION_ERROR
Indicates that there was an error evaluating a CEL expression.

### DATE_PARSE_ERROR
Indicates that there was an error parsing a date string.
{% partial file="graphql/fixtures/errors/DATE_PARSE_ERROR.md" variables={title: "", defaultTab: "Response"} /%}

### DEPENDENCY_ERROR
Indicates that there was an error with a dependent resource.
{% partial file="graphql/fixtures/errors/DEPENDENCY_ERROR.md" variables={title: "", defaultTab: "Response"} /%}

### ENUM_PARSE_ERROR
Indicates that there was an error parsing an enumeration value.
{% partial file="graphql/fixtures/errors/ENUM_PARSE_ERROR.invalidJournalStatus.md" variables={title: "", defaultTab: "Response"} /%}

### FOREIGN_KEY_VIOLATION
Indicates that there was an error with a foreign key constraint.

### GRAPHQL_PARSE_FAILED
Indicates that there was an error parsing the GraphQL query.
{% partial file="graphql/fixtures/errors/GRAPHQL_PARSE_FAILED.md" variables={title: "", defaultTab: "Response"} /%}

### GRAPHQL_VALIDATION_FAILED
Indicates that there was an error validating the GraphQL query.
{% partial file="graphql/fixtures/errors/GRAPHQL_VALIDATION_FAILED.md" variables={title: "", defaultTab: "Response"} /%}

### INTERNAL_SERVER_ERROR
Indicates that there was an error on the server that prevented it from fulfilling the request.

### INTERRUPTED
Indicates that the request was interrupted.

### JSON_PARSE_ERROR
Indicates that there was an error parsing a JSON string.
{% partial file="graphql/fixtures/errors/JSON_PARSE_ERROR.md" variables={title: "", defaultTab: "Response"} /%}

### NOT_FOUND
Indicates that the requested resource was not found.

### NOT_SUPPORTED
Indicates that the requested operation is not supported.

### TIMESTAMP_PARSE_ERROR
Indicates that there was an error parsing a timestamp.

### TRAN_CODE_ERROR
Indicates that there was an error with a transaction code.

For example, with unbalanced entries:
{% partial file="graphql/fixtures/errors/TRAN_CODE_ERROR.entriesUnbalanced.md" variables={title: "", defaultTab: "Response"} /%}

Another example, with a syntax error:
{% partial file="graphql/fixtures/errors/TRAN_CODE_ERROR.syntaxError.md" variables={title: "", defaultTab: "Response"} /%}

### TRANSACTION_ERROR
Indicates that there was an error with a transaction.

### UNIQUE_CONSTRAINT_VIOLATION
Indicates that there was an error with a unique constraint.
{% partial file="graphql/fixtures/errors/UNIQUE_CONSTRAINT_VIOLATION.md" variables={title: "", defaultTab: "Response"} /%}
### UNKNOWN_ERROR
Indicates that an unknown error occurred.

For example, an `unknown error` occuring when selecting accounts
```json
{
  "errors": [
    {
      "message": "input: accounts unknown error",
      "path": [
        "accounts"
      ],
      "extensions": {
        "code": "UNKNOWN_ERROR",
        "retriableError": true
      }
    }
  ],
  "data": null
}
```

### UUID_PARSE_ERROR
Indicates that there was an error parsing a UUID.
{% partial file="graphql/fixtures/errors/UUID_PARSE_ERROR.md" variables={title: "", defaultTab: "Response"} /%}

