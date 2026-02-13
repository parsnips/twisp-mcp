---
title: Layered Accounting
description: Layers help clarify the true state of funds in accounts.
showNext: true
showPrev: true
---

{% comment %}
## Outline

- The Three Layer Model
- Layered Balances
- Layers in Practice
## Other Examples

- writing a transaction to a layer
- balances at a layer

{% /comment %}

## The Three Layer Model

Twisp utilizes a layered model of accounting to distinguish transactions between three important categories:

1. The **settled** layer is for transactions that have fully settled.
2. The **pending** layer includes holds and pending transactions which have been authorized but not yet settled.
3. The **encumbrance** layer contains expected, planned, and scheduled future transactions.

By segmenting transactions in this way, it provides the data integrity needed for more accurate and useful balance calculations and funds flow modelling. It also unlocks features to support a variety of use cases without adding complexity.

For example, we can use the pending layer to verify that an account will have enough funds after the holds and pending transactions have cleared. With the encumbrance layer, we can add transactions scheduled in the future and goals or budgeting tools to set money aside in an account.

Without explicit layers, chaos reigns. Many DIY ledgers that we've seen lack the concept of layers, or only apply the concept partially. This can make it incredibly difficult to reason about the state of a ledger or perform basic accounting operations like cash reconciliation.

{% callout type="note" %}

Not every product will need to use all three layers, but they are always available when and if new versions of your product do need to make use of them.

{% /callout %}

## Layered Balances

Each account can have a different aggregate balance in each layer depending on how transactions have been posted.

To calculate the account balance for a layer, we sum all entries on that layer:

{% math formula="b(l) = \\sum{e_l}" inline=false /%}

Where {% math formula="b(l)" /%} is the balance for a layer and {% math formula="e_l" /%} is the set of all entries on that layer.

All accounts and account sets a debit balance (sum of debit entries), credit balance (sum of credit entries), and normal balance. Learn more about balances in [Balances](/accounting-core/balances).

## Layers in Practice

To demonstrate how layers work in practice, let's consider an account set that has entries posted to its sub-accounts across all three layers. We'll follow the state changes to the account set and its balance as additional entries are written. For simplicity's sake, we'll assume that all amounts are in USD.

At first, the only entry is on the settled layer.

{% partial
  file="graphql/examples/docsBasics/002_GetBertAccountSummary.md"
  variables={
    title: "Example Account Set",
    renderAs: "AccountSet.Summary",
    displayOptions: { heading: false, balance: { heading: true, layers: ["SETTLED", "PENDING", "ENCUMBRANCE"] } }
  }
/%}

Because this entry is on the settled layer, the balance for that layer reflects the amount of that entry.

{% callout %}

The "normal balance" for an account is different for credit normal and debit normal accounts.

Learn more about [Calculating Normal Balance](/accounting-core/balances#calculating-normal-balance).

{% /callout %}

The pending and encumbrance layers have no entries, so they are currently at $0.

Let's see how things change when entries on the pending layer are written. In this case, the entries are modeling a basic "hold-settle" pattern: a hold is placed on the pending layer, and then settled to the settled layer and cleared from the pending layer.

{% partial
  file="graphql/examples/docsBasics/004_GetBertAccountSummary.md"
  variables={
    title: "Example Account Set",
    renderAs: "AccountSet.Summary",
    displayOptions: { heading: false, balance: { heading: true, layers: ["SETTLED", "PENDING", "ENCUMBRANCE"] } }
  }
/%}

Notice that the entries on the pending layer _only affected the pending-layer balance_. The settled balance was only changed by the additional settled-layer entry.

Let's see how the balances change with a few more entries across all layers:

{% partial
  file="graphql/examples/docsBasics/006_GetBertAccountSummary.md"
  variables={
    title: "Example Account Set",
    renderAs: "AccountSet.Summary",
    displayOptions: { heading: false, balance: { heading: true, layers: ["SETTLED", "PENDING", "ENCUMBRANCE"] } }
  }
/%}


Again, notice that the **layered balances** are only aggregating entries from their layer.

## Calculating an Available Balance

While keeping separate balances for each layer is useful for tracking the current, planned, and predicted state of accounts, it is also possible to calculate a balance across layers.

The **available** balance is a special balance that rolls up the balances of the other layer balances. If we show the available balance from the example above, this totaling is clearly visible.

{% partial
  file="graphql/examples/docsBasics/006_GetBertAccountSummary.md"
  variables={
    title: "Example Account Set",
    renderAs: "AccountSet.Summary",
    displayOptions: { heading: false, balance: { heading: true }, entriesGroup: { show: false } }
  }
/%}

We can query the available balance with a configuration option to roll up balances:

* across all layers,
* over the settled and pending layers,
* or just the settled layer (at which point it is equivalent to the settled balance).

One way to think of the available balance is as a function that takes a layer and sums some or all of the other layer balances depending on which layer is provided.

{% math inline=false %}
a(l)=
\\begin{dcases}
  b(S),& \\text{if } l=S\\\\
  b(S) + b(P),& \\text{if } l=P\\\\
  b(S) + b(P) + b(E),& \\text{if } l=E\\\\
\\end{dcases}
{% /math %}

Where {% math formula="a(l)" /%} is the available balance at a given layer, {% math formula="b(l)" /%} is the balance for a layer, and {% math formula="S, P, E" /%} represent the Settled, Pending, and Encumbrance layers.
