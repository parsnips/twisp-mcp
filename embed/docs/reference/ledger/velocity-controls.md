---
title: Velocity Controls
metaTitle: "Reference: Velocity Controls"
description: Reference for the Velocity Control resources within a Twisp ledger.
showNext: true
showPrev: true
---

## The Basics

Velocity Controls in the Twisp Accounting Core provide a mechanism to define and enforce restrictions on the flow of transactions and balances within accounts and account sets. These controls are essential for managing financial governance, ensuring compliance, and mitigating risks by setting specific limits on account activities.

## Components of a Velocity Control

A Velocity Control is made up of two main components:

- **Velocity Control**: Defines the enforcement actions when one or more Velocity Limits are violated.
- **Velocity Limit**: Defines a balance and a parameterizable limit that is enforced. This balance can occur across a number of dimensions.

Once a Velocity Control is defined, it may be **attached** to accounts and/or account sets to begin enforcing the defined limits for those accounts or sets. If during the posting of a transaction, an attached Velocity Control is violated an exception will be generated and depending on the enforcement configuration, the transaction will be rejected, voided or warned.


### Velocity Control

There are 3 key components that make up a Velocity Control

1. **Enforcement**: A Velocity Control can define it's enforcement mechanism, there are three `WARN`, `VOID` and `REJECT`. This defines whether the violated entry is written to the ledger and how the ledger will respond to the violating entry.
2. **Velocity Limits**: The Velocity Limits on a control define one or more balances that are under enforcement of a particular control.  These balances can contain windowing and filtering criteria to keep track of different types of balances that will be enforced on by the Velocity Control.
3. **Condition**: A CEL expression which evaluates to a boolean value allows for fine grained control to determine if the Velocity Control will execute enforcement. For example, a Velocity Limit may be exceeded, but the control only enforces `PENDING` transactions.

### Velocity Limits

Velocity Limits form the foundation of Velocity Controls. They specify the maximum allowable balance or transaction volume over a defined period. Each Velocity Limit is tailored to control spending, deposits, or any financial activity based on parameters like amount, currency, and account layers (SETTLED, PENDING, ENCUMBRANCE). By enforcing these limits, organizations ensure that financial activities remain within predefined thresholds, maintaining financial discipline and control.

There are 5 key components of a Velocity Limit:

1. **Window**: Every Velocity Limit has a default set of dimensions (currency, journal, account) and this `window` can define additional dimensions, such as bucketed time periods or values from account or entry metadata.
2. **Limit**: The limit to apply, which currently only supports an available balance limit. This defines the layer, amount, and normal balance type that the limit supports.
3. **Condition**: A CEL expression which evaluates to a boolean value for filtering if an entry is eligible to apply to the limit.  For example if a limit is only applied to a certain category of transactions defined in transaction metadata.
4. **Currency**: Defines which currency the limit applies to.
5. **Params**: The definition of the parameters available when attaching a a limit used in a Velocity Control to an account.  This allows values to supplied at attachment time so that a Velocity Control can have different enforced limits for different accounts.


### Attaching Velocity Controls

Attaching a Velocity Control to an Account or Set enables that Account or Set to begin enforcing entries that are posted.    

There are three components to an attachment:

1. **Velocity Control Id**: Identifies the particular Velocity Control to attach.
2. **Account or Set Id**: The account or set to attach control to.
3. **Params**: A JSON set of parameters that satisfies all the defined parameters on the limits attached to the controls.

{% callout type="note" %}
When defining parameters on Velocity Limits use _unique parameter names_ to avoid collisions between parameters.
{% /callout %}

### Overriding Velocity Control Enforcement

Velocity Control enforcement can be escalated or de-escalated for a control with an action configured at the `VOID` or `REJECT` level, by passing the `overrideVelocityEnforcement` parameter on [postTransaction]($gql:mutation:postTransaction).  

This is useful for setting to `WARN` for force post transactions that require disabled velocity control enforcement.

## Velocity Control Operations

Use GraphQL queries and mutations to read, create, delete and attach Velocity Controls and limits:

- [`Query.velocityControl()`]($gql:query:velocityControl): Get a single Velocity Control
- [`Query.velocityControls()`]($gql:query:velocityControls): Query Velocity Controls
- [`Query.velocityLimit()`]($gql:query:velocityLimit): Get a single Velocity Limit 
- [`Query.velocityLimits()`]($gql:query:velocityLimits): Query Velocity Limits 
- [`Mutation.createVelocityControl()`]($gql:mutation:createVelocityControl): Create a Velocity Control
- [`Mutation.updateVelocityControl()`]($gql:mutation:updateVelocityControl): Update a Velocity Control
- [`Mutation.deleteVelocityControl()`]($gql:mutation:deleteVelocityControl): Delete a Velocity Control
- [`Mutation.createVelocityLimit()`]($gql:mutation:createVelocityLimit): Create a Velocity Limit
- [`Mutation.updateVelocityLimit()`]($gql:mutation:updateVelocityLimit): Update a Velocity Limit
- [`Mutation.deleteVelocityLimit()`]($gql:mutation:deleteVelocityLimit): Delete a Velocity Limit
- [`Mutation.attachControl()`]($gql:mutation:attachVelocityControl): Attach Velocity Control to an Account or Account Set
- [`Mutation.detachControl()`]($gql:mutation:detachVelocityControl): Detach Velocity Control from Account or Account Set

