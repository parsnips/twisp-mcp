---
title: Reference
description: Definitions of all components in the Twisp Accounting Core and GraphQL API.
showNext: false
showPrev: false
---

## Reference Sections

{% quick-links %}

{% quick-link title="Ledger" href="/reference/ledger" description="Ledger resources in the Twisp accounting core." icon="journal" /%}

{% quick-link title="ACH" href="/reference/ach" description="ACH resources in the Twisp accounting core." icon="money" /%}

{% quick-link title="GraphQL" href="/reference/graphql" description="Type definitions for the full GraphQL schema." icon="graphql" /%}

{% quick-link title="API" href="/reference/api" description="Key components and concepts for interacting with the API." icon="json" /%}

{% quick-link title="CEL" href="/reference/cel" description="Packages and functions available in the common expression language runtime." icon="celExpression" /%}

{% /quick-links %}

## Type Diagrams

Use these diagrams for a high-level overview of the type system within the Twisp ledger.

### Basic Relationships

Simplified type relationship diagram showing basic connections between types.

- [Journals]($gql:object:Journal) have many [Transactions]($gql:object:Transaction) and [Entries]($gql:object:Entry)
- [Transactions]($gql:object:Transaction) are defined by their [TranCode]($gql:object:TranCode) and write multiple [Entries]($gql:object:Entry) to the ledger.
- [Transactions]($gql:object:Transaction) can be linked to other correlated [Transactions]($gql:object:Transaction)
- [Entries]($gql:object:Entry) are written to a specific [Account]($gql:object:Account) and [Journal]($gql:object:Journal)
- [Accounts]($gql:object:Account) roll up a [Balance]($gql:object:Balance) of all [Entries]($gql:object:Entry) for each [Journal]($gql:object:Journal)
- [Account Sets]($gql:object:AccountSet) are groupings of [Accounts]($gql:object:Account) and/or other [Account Sets]($gql:object:AccountSet) which also roll up [Balances]($gql:object:Balance)

{% diagram name="core_erd_minimal" alt="Simplified entity relationship diagram for Twisp core" /%}

### Entity Relationship Diagram

Although the Twisp core is not a relational database, it does enforce referential integrity between related records and thus we can model it with an RDB-style ERD.

Only select fields for each type have been shown in this diagram to aid readability. For the full reference of each type, see the [GraphQL reference](/reference/graphql).

{% diagram name="core_erd" alt="RDB-style entity relationship diagram for Twisp core" /%}
