---
title: Making API Requests with cURL
metaTitle: "Tutorial : Making API Requests with cURL"
description: In this tutorial, we will learn how to make requests against the Twisp GraphQL API using the curl command line tool.
showNext: false
showPrev: false
---

**Prerequisites**

To complete this tutorial, you'll need:

- The URL endpoint for the API
- The `accountId` of your tenant
- An authenticated JWT

---

To make HTTPS requests to the Twisp GraphQL API using cURL, follow these steps:

1. Open a command-line interface (Terminal on macOS/Linux or Command Prompt on Windows).

2. Type in the following command:

    ```shell
    curl 'https://api.<region>.cloud.twisp.com/financial/v1/graphql' \
    ```

    This line specifies the URL to which the cURL request will be made. Replace `<region>` with the AWS region for your ledger, e.g. `us-east-1`.

3. Add the required headers by typing the following lines:

    ```shell
      -H 'authorization: Bearer <JWT>' \
      -H 'content-type: application/json, application/json' \
      -H 'x-twisp-account-id: <TwispAccountId>' \
    ```

    Replace `<JWT>` with your authenticated JWT and `<TwispAccountId>` with your Twisp account ID for the tenant.

    The `-H` flag is used to specify custom headers. In this case, the headers being set are:

    - `authorization` for authentication purposes.
    - `content-type` to indicate the type of data being sent in the request body.
    - `x-twisp-account-id` for specifying the account ID to be used in the request.

4. Add the request payload by typing the following line:

    ```shell
      --data-raw '{"query":"query { accounts(index: { name: STATUS }, where: { status: { eq: "ACTIVE" } }, first: 5) { nodes { accountId code name } } }"}'
    ```

    This line contains the `--data-raw` flag followed by a JSON string. This JSON string includes a GraphQL query to fetch accounts with the specified conditions (in this case, the first 5 active accounts). This is where you write your GraphQL operation.

5. After entering all the lines, the complete cURL request should look like:

    ```shell
    curl 'https://api.<region>.cloud.twisp.com/financial/v1/graphql' \
      -H 'authorization: Bearer <JWT>' \
      -H 'content-type: application/json, application/json' \
      -H 'x-twisp-account-id: <TwispAccountId>' \
      --data-raw '{"query":"query { accounts(index: { name: STATUS }, where: { status: { eq: "ACTIVE" } }, first: 5) { nodes { accountId code name } } }"}'
    ```

6. Press Enter to execute the cURL request. If successful, you should receive a JSON response containing the requested data.
