---
title: API Requests via HTTPS
metaTitle: "Reference : API Requests via HTTPS"
description: Components of a GraphQL API request.
metaDescription: HTTPS request reference for the Twisp GraphQL API.
showNext: true
showPrev: false
---

All HTTPS requests (whether they are GraphQL [](/reference/graphql/queries) or [](/reference/graphql/mutations)) are made using the `POST` method, with the GraphQL operation provided in the request body.

{% callout type="note" %}
For a step-by-step guide, see the [](/tutorials/api/making-api-requests-with-curl) tutorial.
{% /callout %}

When you are ready to start making requests to the Twisp API from your application code, the HTTPS request must include:

- **URL**: the URL of the Twisp GraphQL API.
- **Headers**: to authorize the request and specify content type. At minimum, these must include:
  - `authorization` with a bearer JWT
  - `x-twisp-account-id` for specifying the account ID to be used in the request
  - `content-type` set to `application/json`
- ***Body**: the JSON-formatted GraphQL operation.

Here is an example request for the first 10 accounts with a code starting with `BERT.`:

```shell
curl 'https://api.us-east-1.cloud.twisp.com/financial/v1/graphql' \
  -H 'authorization: Bearer <JWT>' \
  -H 'content-type: application/json' \
  -H 'x-twisp-account-id: <TwispAccountId>' \
  --data-raw '{"query":"query { accounts(index:{name:CODE}, where:{code:{like:\"BERT.\"}},first:10) { nodes { accountId name code } } }"}'
```

The `<JWT>` and `<TwispAccountID>` are placeholders. If you created a [Tenant]($gql:object:Tenant) named "Sandbox" with account ID `Sandbox0d6af983` and had a JWT encoded as `0cca9e8f` (obviously it would be much longer in actual usage), then you would replace `<TwispAccountId>` with `Sandbox0d6af983` and `<JWT>` with `0cca9e8f`.

{% callout type="note" %}
As you are first starting out with Twisp, we recommend using the GraphiQL interface in the Twisp console. It is a rich learning environment where you can focus on writing your GraphQL operations with built-in documentation and autocomplete.
{% /callout %}
