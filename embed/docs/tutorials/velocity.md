---
title: "Ledger Invariants: Create and Attach Velocity Controls to Accounts"
metaTitle: "Tutorial: Creating and Attaching Velocity Controls"
description: |
  In this tutorial we'll cover creating velocity controls in Twisp to ensure balances cannot exceed thresholds.
showNext: false
showPrev: true
---


## How-to Guide: Add Velocity Controls in Twisp

**Goal:** This guide shows you how to create and apply velocity controls to manage transaction limits on accounts or account sets in Twisp.

**Prerequisites:**

*   Access to the Twisp GraphQL API.
*   Appropriate permissions to create/manage velocity controls and limits (`velocity.*` mutations).
*   UUIDs of the `Account` or `AccountSet` you wish to apply controls to.
*   Understanding of the specific limits you want to enforce (e.g., daily amount, transaction count).
*   (Optional) Familiarity with Common Expression Language (CEL) for defining complex `window`, `condition`, or `limit` expressions.

### Key Concepts

*   **Velocity Control:** A container that groups one or more `VelocityLimit`s. It defines an `enforcement` action (Warn, Void, Reject) and an optional `condition` (CEL expression) for when the control applies.
*   **Velocity Limit:** Defines a specific rule, such as a maximum amount or count within a defined `window`. It includes:
    *   `name`/`description`: Human-readable identifiers.
    *   `window`: CEL expressions defining the time frame or grouping criteria (e.g., daily, monthly, per merchant). Uses `PartitionKey`.
    *   `limit`: The actual threshold (amount, layer, direction, start/end times). Uses `Limit`.
    *   `currency`: The currency the limit applies to (or all if empty).
    *   `condition`: Optional CEL expression for when this specific limit applies.
    *   `params`: Optional parameters needed for dynamic limits (e.g., passing a specific merchant ID).
*   **Attachment:** Linking a `VelocityControl` to an `Account` or `AccountSet` makes the control active for that entity. Parameters required by the associated limits can be provided during attachment.

### Steps

1.  **Define the Velocity Limit(s):**
    *   Create each specific rule you need using the `createVelocityLimit` mutation.
    *   Define its `name`, `description`, `window` (using `PartitionKeyInput`), `limit` (using `LimitInput`), `currency`, and optional `condition` and `params`.

    ```graphql
    # Example: Create a simple daily spending limit
    mutation CreateDailyLimit($limitId: UUID!) {
      createVelocityLimit(input: {
        velocityLimitId: $limitId
        name: "Daily Spending Limit"
        description: "Limit spending to $100 per day."
        currency: "USD"
        window: [{alias: "date", value: "context.vars.transaction.effective"}] # Daily window
        limit: {
          balance: [{
            layer: "SETTLED" # Or PENDING, etc.
            amount: "decimal('100.00')"
            normalBalanceType: "DEBIT" # Limit debits
          }]
        }
        # Optional: condition: "context.vars.account.metadata.applyDailyLimit == true"
        # Optional: params: [{name: "maxAmount", type: DECIMAL}] # If limit amount was dynamic
      }) {
        velocityLimitId
        name
      }
    }
    ```
    *   *Variables:* `{ "limitId": "<generate-a-uuid-for-the-limit>" }`

2.  **Create the Velocity Control:**
    *   Use the `createVelocityControl` mutation.
    *   Define its `name`, `description`, `enforcement` action (e.g., `REJECT`), and optional overall `condition`.
    *   Link the limits created in Step 1 using their `velocityLimitId`s in the `velocityLimitIds` array.

    ```graphql
    # Example: Create a control using the daily limit from Step 1
    mutation CreateSpendingControl($controlId: UUID!, $limitId: UUID!) {
      createVelocityControl(input: {
        velocityControlId: $controlId
        name: "Standard Spending Control"
        description: "Enforces daily spending limits."
        enforcement: { action: REJECT }
        velocityLimitIds: [$limitId] # Link the limit(s)
        # Optional: condition: "!" + context.vars.account.metadata.exemptFromControl"
      }) {
        velocityControlId
        name
        limits {
          velocityLimitId
        }
      }
    }
    ```
    *   *Variables:* `{ "controlId": "<generate-a-uuid-for-the-control>", "limitId": "<uuid-from-step-1>" }`

3.  **Attach the Control to an Account or AccountSet:**
    *   Use the `attachVelocityControl` mutation.
    *   Provide the `velocityControlId` (from Step 2) and the target `accountId` (UUID of the `Account` or `AccountSet`).
    *   If any attached limits defined `params`, provide their values in the `params` JSON object.

    ```graphql
    # Example: Attach the control to a specific account
    mutation AttachControlToAccount($controlId: UUID!, $accountId: UUID!) {
      attachVelocityControl(
        velocityControlId: $controlId
        accountId: $accountId
        # Optional: params: { "maxAmount": "150.00" } # If limit had params
      ) {
        velocityControlId
        name
      }
    }
    ```
    *   *Variables:* `{ "controlId": "<uuid-from-step-2>", "accountId": "<target-account-or-set-uuid>" }`

### Verification

*   Query the `VelocityControl` using its ID to confirm its limits.
*   Query the `Account` or `AccountSet` and check its `controls` or use the `attachedControls` query to see attached controls.
*   Post a transaction that should interact with the limit to ensure the expected enforcement action occurs.
*   Query the `velocity` field on the `Account` or `AccountSet` to check remaining limits.
