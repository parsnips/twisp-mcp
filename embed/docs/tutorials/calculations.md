---
title: Implementing Custom Calculations
metaTitle: "Tutorial: Implementing Custom Calculations"
description: |
  Learn how to define, create, and use custom calculations in Twisp to track specialized balances (e.g., by effective date or metadata tags) and query them efficiently using the GraphQL API.
showNext: true
showPrev: true
---

Custom calculations allow you to define specialized ways to aggregate and query balances based on dimensions beyond the standard `accountId`, `journalId`, and `currency`. This enables tracking balances grouped by criteria like effective date, specific metadata tags, or other transaction/entry attributes relevant to your business logic.

By the end of this tutorial, you will be able to design and create a custom calculation, and then query the specialized balances it generates.

{% callout type="task" %}
- Design a custom calculation with specific dimensions using CEL.
- Create the calculation using the `createCalculation` mutation.
- Query the calculated balances using the `balances` query with the `CALCULATION` index.
- (Optional) Attach a `LOCAL` scope calculation to specific accounts or sets.
{% /callout %}

---

{% partial file="tutorials_getting_started.md" /%}

## Designing a Custom Calculation

Before creating a calculation, you need to define how you want to group and filter balances. Let's design a calculation to track balances *per effective date* for each account.

Key components of a calculation definition:

*   **`calculationId`:** A unique UUID you provide to identify this calculation.
*   **`code`:** A unique, human-readable string code for the calculation (e.g., `PER_EFFECTIVE_DATE`).
*   **`description`:** Explains the purpose of the calculation.
*   **`scope`:**
    *   `GLOBAL`: The calculation applies to *all* accounts and sets within the journal(s) it's implicitly associated with (usually the default journal unless specified otherwise during querying or attachment). This is the default.
    *   `LOCAL`: The calculation *only* applies to accounts or sets where it has been explicitly attached using the `attachCalculation` mutation.
*   **`dimensions`:** An array defining *how* the balance is indexed or grouped. Each dimension requires:
    *   `alias`: A name for the dimension (e.g., `effectiveDate`).
    *   `value`: A [CEL expression](https://twisp.com/docs/reference/cel) referencing `context.vars` (which contains the `entry`, `transaction`, `account`, etc. being processed) to get the dimension's value (e.g., `context.vars.transaction.effective` to group by the transaction's effective date).
*   **`condition` (Optional):** A CEL expression that must evaluate to `true` for an entry to be included in this calculation. If omitted, all entries matching the dimensions are included. Example: `context.vars.entry.amount.units() > decimal('0.0')` would only include entries with positive amounts.

For our example (tracking balance per effective date):
*   **`calculationId`:** `"5867b5dd-fc69-416c-80f5-62e8a53610d5"` (Use your own unique UUID)
*   **`code`:** `PER_EFFECTIVE_DATE`
*   **`description`:** `Track balances per effective date in an account.`
*   **`scope`:** `GLOBAL` (applies everywhere by default)
*   **`dimensions`:** `alias: "effectiveDate"`, `value: "context.vars.transaction.effective"`
*   **`condition`:** None (include all entries)

## Create the Custom Calculation

Use the `createCalculation` mutation with the parameters defined above.

```graphql
mutation CreateCalculationPerEffectiveDate {
  createCalculation(
    input: {
      calculationId: "5867b5dd-fc69-416c-80f5-62e8a53610d5" # Use your own UUID
      code: "PER_EFFECTIVE_DATE"
      description: "Track balances per effective date in an account."
      scope: GLOBAL
      dimensions: [
        {
          alias: "effectiveDate"
          value: "context.vars.transaction.effective"
        }
      ]
      # condition: "context.vars.entry.layer == 'SETTLED'" # Optional example condition
    }
  ) {
    calculationId
    code
    description
    scope
    dimensions {
      alias
      value
    }
    condition
    version
    created
    modified
  }
}
```

```json
{
  "data": {
    "createCalculation": {
      "calculationId": "5867b5dd-fc69-416c-80f5-62e8a53610d5",
      "code": "PER_EFFECTIVE_DATE",
      "description": "Track balances per effective date in an account.",
      "scope": "GLOBAL",
      "dimensions": [
        {
          "alias": "effectiveDate",
          "value": "context.vars.transaction.effective"
        }
      ],
      "condition": null,
      "version": 1,
      "created": "2023-10-27T12:00:00Z",
      "modified": "2023-10-27T12:00:00Z"
    }
  }
}
```

This mutation creates the `PER_EFFECTIVE_DATE` calculation. The response confirms its structure. Twisp will now start computing balances grouped by `accountId`, `journalId`, `currency`, *and* `effectiveDate`.

## Query Using the Custom Calculation

To retrieve balances computed by a custom calculation, use the standard [`balances`]($gql:query:balances) query, specifying `index: { name: CALCULATION }`.

In the `where` clause, you **must** provide:
*   The [`accountId`]($gql:input:BalanceFilterInput:accountId) (or `accountSetId`).
*   The [`calculationId`]($gql:input:BalanceFilterInput:calculationId) of your custom calculation.
*   A [`dimension`]($gql:input:BalanceFilterInput:dimension) filter matching the structure defined in your calculation. Provide the specific dimension `alias`(es) and the `value`(s) you want to query (e.g., a specific date for the `effectiveDate` alias).

```graphql
query GetCalculationBalanceForDate($accountId: UUID!, $calcId: UUID!, $journalId: UUID) {
  balances(
    index: { name: CALCULATION }
    where: {
      accountId: { eq: $accountId }
      journalId: { eq: $journalId }      # Optional: Specify journal or uses default
      calculationId: { eq: $calcId }    # ID of the custom calculation
      dimension: {
        # 'eq' provides the exact dimension values to match
        eq: {
          # Alias must match the one defined in the calculation
          effectiveDate: "2023-09-21" # Specific dimension value to query
        }
      }
      # currency: { eq: "USD" } # Optional: Filter by currency if needed
    }
    first: 10
  ) {
    nodes {
      accountId
      journalId
      currency
      calculationId # Confirms which calculation produced this balance
      dimensions    # Shows the specific dimension values for this balance record
      settled {
        normalBalance {
          units
        }
      }
      # Query other balance layers (pending, encumbrance) or fields as needed
    }
    pageInfo {
      hasNextPage
      endCursor
    }
  }
}
```

**Variables for the query:**

```json
{
  "accountId": "1fd1dd3e-33fe-4ef5-9d58-676ef8d306b5",
  "calcId": "5867b5dd-fc69-416c-80f5-62e8a53610d5",
  "journalId": "822cb59f-ce51-4837-8391-2af3b7a5fc51"
}
```

**Example Response:**

```json
{
  "data": {
    "balances": {
      "nodes": [
        {
          "accountId": "1fd1dd3e-33fe-4ef5-9d58-676ef8d306b5",
          "journalId": "822cb59f-ce51-4837-8391-2af3b7a5fc51",
          "currency": "USD",
          "calculationId": "5867b5dd-fc69-416c-80f5-62e8a53610d5",
          "dimensions": {
            "effectiveDate": "2023-09-21"
          },
          "settled": {
            "normalBalance": {
              "units": "5.25"
            }
          }
        }
        // ... other balances matching the criteria might appear if pagination allows
      ],
      "pageInfo": {
        "hasNextPage": false,
        "endCursor": "..."
      }
    }
  }
}
```

This query retrieves the balance record specifically calculated for the account `"1fd1dd3e-..."`, using the `"PER_EFFECTIVE_DATE"` calculation, for the exact effective date `"2023-09-21"`. The `dimensions` field in the response confirms the `effectiveDate` value associated with this particular balance record.

## Attaching Calculations (`LOCAL` Scope Only)

If you created a calculation with `scope: LOCAL`, it only computes balances for accounts or sets where it's explicitly attached. `GLOBAL` calculations apply automatically and do not need attaching.

Use the [`attachCalculation`]($gql:mutation:attachCalculation) mutation:

*   Provide the `calculationId` of the `LOCAL` calculation.
*   Provide the `accountId` of the target [Account]($gql:object:Account) or [AccountSet]($gql:object:AccountSet).
*   Optionally provide the `journalId` if attaching to an `Account` and not using the default journal (ignored for `AccountSet`).

```graphql
# Example: Attaching a LOCAL calculation to an account
mutation AttachLocalCalcToAccount($calcId: UUID!, $accountId: UUID!, $journalId: UUID) {
  attachCalculation(
    input: {
      calculationId: $calcId    # ID of the LOCAL calculation
      accountId: $accountId    # ID of the Account or AccountSet to attach to
      journalId: $journalId    # Optional: Specify journal for Account attachment
    }
  ) {
    # Response confirms attachment details
    calculationId
    accountId
    journalId
    calculation { code }
    account { name }
  }
}
```

**Variables for the mutation:**

```json
{
  "calcId": "YOUR_LOCAL_CALCULATION_ID",
  "accountId": "YOUR_TARGET_ACCOUNT_ID",
  "journalId": "YOUR_JOURNAL_ID" 
}
```

{% callout type="note" %}
Creating a GLOBAL calculation or attaching a LOCAL calculation will initiate a backfill of that calculation using historical entries. It may take a few minutes for _new_ entries to show up in the balance while the backfill is ongoing.
{% /callout %}

## Conclusion

In this tutorial, you learned how to design custom calculations by defining dimensions and optional conditions using CEL. You created a calculation using the `createCalculation` mutation and queried its specialized balances via the `balances` query with the `CALCULATION` index and dimension filters. Finally, you saw how to attach `LOCAL` scope calculations to specific accounts or sets using `attachCalculation`.

Custom calculations offer powerful flexibility for tracking and querying balances based on criteria specific to your application's needs, enabling more granular financial insights and reporting.
