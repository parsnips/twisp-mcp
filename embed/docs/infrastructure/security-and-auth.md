---
title: Security and Auth
description: Security policies control access to encrypted resources.
showNext: true
showPrev: true
---

Ensuring the proper security of data is foundational to the Twisp infrastructure. Out of the box, data is [securely encrypted both at rest and in transit](#encryption) and all access is authenticated via JWT tokens issued by the [OpenID Connect 1.0 protocol](https://openid.net/connect/).

Users can configure additional access settings by defining [tenants](#provisioning-tenants) and [clients](#creating-clients-and-policies) through GraphQL mutations to set precise permissions on specific operations.

## The authentication & authorization flow

Twisp applies policies to an authenticated principal before allowing access to resources.

In a simplified form, the stages for the authentication process are:

1. A principal issues an HTTPS request with their OIDC JWT and Twisp account IDs in the headers.
2. The principal name is extracted from the JWT.
3. Twisp finds the corresponding policies to apply using the principal name.
4. The resulting policies are evaluated for authorization.

### Making an authenticated request

All HTTPS requests to the [GraphQL API](/reference/graphql) must provide the following headers to allow for authorization:

```
Authorization: Bearer <JWT>
x-twisp-account-id: <Twisp account id>
```

The JWT can be issued either with OpenID Connect or AWS IAM.

### Authorizing a principal

Within a client configuration, the **security principal** is the authenticated identity accessing the system. The client authorization system retrieves [policies](#creating-clients-and-policies) based on the principal name. For OIDC this is typically the issuer found in the `iss` claim. For AWS IAM it is the IAM role, user, or other AWS identity placed in the `sub` claim when Twisp vends a token.

## Provisioning Tenants

Each cloud environment requires its own tenant. Use the admin GraphQL namespace to create tenants from an existing account:

```graphql
mutation CreateTenants {
  admin {
    staging:createTenant(
      input: {
        id: "d9d8f1c0-0299-4d5b-b2b6-85beafdda28b"
        accountId: "TwispStaging"
        name: "staging"
        description: "staging environment for Twisp"
      }
    ) {
      accountId
      name
    }
    production:createTenant(
      input: {
        id: "848df974-4133-4ee4-ab45-86e5e29b6822"
        accountId: "TwispProd"
        name: "production"
        description: "production environment for Twisp"
      }
    ) {
      accountId
      name
    }
  }
}
```

Once a tenant exists, clients may be created within that tenant to govern access.

## Creating Clients and Policies

Authentication to Twisp works with any OIDC-compliant token supplied through the `Authorization: Bearer <token>` header. A corresponding client must exist in Twisp for the principal identified in the token so Twisp knows which policies to apply.

### Third-party OIDC client example

For tokens issued by an external identity provider, set the client `principal` to the OIDC issuer URL (`iss` claim). The following mutation creates per-user policies for tokens issued by Google:

```graphql
mutation CreateGoogleCloudClient {
  auth {
    createClient(
      input: {
        name: "michael gcloud readonly"
        principal: "https://accounts.google.com"
        policies: [
          {
            actions: [SELECT]
            resources: ["*"]
            effect: ALLOW
            assertions: {
              isMike: "context.auth.claims.email == 'michael@twisp.com'"
            }
          }
          {
            actions: [SELECT, INSERT, UPDATE, DELETE]
            resources: ["*"]
            effect: ALLOW
            assertions: {
              isJarred: "context.auth.claims.email == 'jarred@twisp.com'"
            }
          }
        ]
      }
    ) {
      principal
    }
  }
}
```

### AWS IAM principal example

Twisp can exchange a presigned AWS STS `GetCallerIdentity` request for an OIDC token. The resulting token uses the AWS identity as the principal. Create a client for that principal to grant access:

```graphql
mutation CreateIAMAuthClient {
  auth {
    createClient(
      input: {
        principal: "arn:aws:iam::012345678901:role/example-role"
        name: "example role policy"
        policies: [
          {
            effect: ALLOW
            actions: [SELECT, INSERT, UPDATE, DELETE]
            resources: ["*"]
          }
        ]
      }
    ) {
      principal
    }
  }
}
```

### Understanding policies

Each policy defines an `effect` (`ALLOW` or `DENY`), the `actions` permitted or denied, the `resources` in scope, and optional CEL `assertions` that must evaluate to `true`.

{% callout type="tip" title="Logical combination of policies" %}
Policies are evaluated as a chain of logical `AND` statements scoped to the requested resource and action. A principal must have *at least one* `ALLOW` policy for a resource/action pair. A single `DENY` policy on that pair blocks the operation.
{% /callout %}

Values for `resources` and `actions` support `*` and `?` wildcards. Actions include:

- `db:Select` (`SELECT` in GraphQL) : read a document
- `db:Insert` (`INSERT`) : create a document
- `db:Update` (`UPDATE`) : change document fields
- `db:Delete` (`DELETE`) : remove a document

Resource identifiers follow the format `namespace.ledger.<document|indexes|joins|references>.propertyName`. The current namespace is `financial`.

Assertions use [Common Expression Language (CEL)](/reference/cel) and can access `context.auth.claims` (token claims) and `context.document` (the document being acted on). If any assertion evaluates to `false`, the policy is skipped.

## Using OpenID Connect tokens

Any JWT generated by an OpenID Connect 1.0 capable issuer is supported by Twisp for authentication. When an API endpoint receives the JWT, it validates the token signature against the issuer and—if valid—invokes the endpoint with a security principal set to the issuer `iss` claim. All claims are embedded in `context.auth.claims` for policy evaluation.

## Issuing tokens with AWS Identity and Access Management (IAM)

Twisp can vend an OpenID Connect token in exchange for an authenticated AWS IAM role or user. This allows services running in AWS to retrieve a Twisp-issued token for access. The issuer of these tokens is `https://auth.${AWS::Region}.prod.twisp.com/token/iam`, and the resulting principal name equals the original IAM identity stored in the token `sub` field.

Twisp also provides a hosted service to exchange a presigned `GetCallerIdentity` request for a token, making it convenient to obtain OIDC tokens inside AWS environments.

## Cloud endpoints

- AWS Twisp Token: `https://auth.us-east-1.cloud.twisp.com/token/iam`
- Financial GraphQL: `https://api.us-east-1.cloud.twisp.com/financial/v1/graphql`
- gRPC: `https://api.us-east-1.cloud.twisp.com:50051`

## Encryption

### Encryption in transit

Data is encrypted in transit via HTTPS connections for all external and internal API operations. This protects against man-in-the-middle attacks and prevents eavesdropping on Twisp traffic.

### Encryption at rest

All data stored in Twisp is encrypted at rest. Encryption keys are stored in [AWS Key Management Service](https://aws.amazon.com/kms/), so even if a malicious actor gains access to the storage layer they cannot read the data without the proper keys.
