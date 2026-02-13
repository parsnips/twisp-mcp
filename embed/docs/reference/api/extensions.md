---
title: Extensions
metaTitle: "Reference : API Extensions"
description: The Twisp GraphQL API supplies additional metadata about operations through the "extensions" subdocument.
metaDescription: Extensions reference for the Twisp GraphQL API.
showNext: false
showPrev: true
---

API responses may contain an `"extensions"` field in addition to the standard `"data"` field if certain request headers are specified.

## API Usage & Billing

Billing stats about the current request may be returned within the `"extensions"` subdocument by including a `x-twisp-include-billing: true` header in the request.

For example, here's is a response from invoking a bank transfer tran code with `postTransaction` and the `x-twisp-include-billing` HTTP header set to `true`:

```json
{
  "data": {
    "postTransaction": {
      "transactionId": "5c328550-bba3-423b-a58a-b3f9786a80ad"
    },
  },
  "extensions": {
    "billing": {
      "duration": 229272730,
      "read": {
        "bytes": 5604,
        "count": 10,
        "units": 10
      },
      "write": {
        "bytes": 4850,
        "count": 16,
        "units": 16
      }
    }
  }
}
```

{% callout type="note" %}
If using the gRPC apis the billing information will be emitted on the following trailer keys:


- billing-read-bytes
- billing-read-units
- billing-write-bytes
- billing-write-units
- billing-duration
{% /callout %}
