---
title: Designing a Tran Code
metaTitle: "Tutorial: Designing a Tran Code"
description: |
  In this tutorial, we will cover some of the factors that go into designing a transaction code.
showNext: true
showPrev: true
---

## Preliminary Questions

Tran codes should mirror a type of transaction that your system handles. They are meant to encapsulate and abstract accounting activity, and so the first step to designing a good tran code is **having a clear understanding of what that activity is**.

Thus, before we get into designing a tran codes, it is important to first clarify:

- What **type** of transaction is this? ACH transfer? Credit card purchase? Foreign exchange? Something else?
- Which **accounts** are involved?
  - Where is the money coming from and where is it going?
  - Are there any additional accounts that need to be debited/credited?
- How many **entries** should be written to the ledger?
  - Which entries go on the debit side and which on the credit side?
- Do these entries reflect a **settled** amount, or are they still **pending** a final settlement?

## Naming and Documenting

The **type** of transaction should be used to give the tran code a good name through its `code` field. The `code` field is the primary human-friendly identifier. It should provide a concise indication of what the tran code does and is used for.

We recommend using codes that just give enough information to be easily identifiable without being over-wordy. For consistency, we recommend using UPPER_SNAKE_CASE formatting for the `code`.

Prefer:

- ✅ `ACH_CREDIT`
- ✅ `CARD_HOLD_CANCEL`
- ✅ `INTEREST_ADJUSTMENT`

Avoid:

- ❌ `ACH`
- ❌ `CancelHold`
- ❌ `adjusting_interest_for_personal_loan_accounts`

{% callout type="note" %}

Depending on the size and complexity of your tran code library, you may choose to implement more formalized patterns and structures for naming your tran codes. For example, some organizations might use abbreviated versions of operations (`HLD`, `STL`, `DEP`, `CLR`) to keep tran code names extra terse.

{% /callout %}

Of course, only so much information can be communicated through a short string of characters. Because tran codes act as the API for your funds flow, they should also be well documented.

The `description` field is where this documentation can live. Use it to add additional context about why the tran code exists, how it should be used, and whatever other information would benefit those interacting with your ledger. This field supports Markdown formatting.

## Defining the Transaction & Ledger Entries

Posted transactions and the ledger entries written are defined within the `transaction` and `entries` fields, respectively. These are effectively templates used to generate the [Transaction]($gql:object:Transaction) and [Entries]($gql:object:Entry) records.

Within the `transaction` field, we can define values that describe important aspects of the transaction like its `effective` date and which `journal` it should be written to. Other [Transaction]($gql:object:Transaction) fields like the `correlationId` and `metadata` can also be defined here, although they are optional.

The `entries` field is a list of the templates for the [Entries]($gql:object:Entry) written to the ledger. Each ledger entry must define its `accountId`, amount (in `units` and `currency`), `direction` (DEBIT or CREDIT), `layer` (SETTLED, PENDING, or ENCUMBRANCE), and  `entryType`. Optionally, a `description` may be written here as well.

{% callout type="note" %}

The `entryType` for an entry is similar to the `code` for a tran code: it is a short identifier for describing the type of activity that the entry represents. In many cases, the `entryType` is just an extension of the `code`. For example, an `ACH_CREDIT_FEE` tran code might write ledger entries with types `ACH_CREDIT_FEE_DR` for the debit-side and `ACH_CREDIT_FEE_CR` for the credit-side entry.

{% /callout %}

How these entry definitions are written depends upon the transaction type and other information gathered as part of the pre-design process.

{% partial file="graphql/examples/tranCodeDocs/achCredit_noParams.md" variables={title: "Example"} /%}

## Parameterizing Inputs

In most cases, not all of the information needed to write transactions and entries is available at the time of designing a tran code, but instead needs to be passed in at runtime (i.e. when the transaction is posted). For example, most transactions are not for fixed amounts, and so these amounts need to be specified when posting the transaction.

This is where the `params` field of a tran code comes in. With the `params`, we can define parameters of a transaction which can then be referenced inside of the values defined for the `transaction` and `entries`.

Say we wanted to write entries where an `amount` (in decimal units) is supplied at posting time. To do this, we need to do two things:

1. Define an `amount` parameter inside of the `params` object.
2. Reference the `amount` value from `params` within the `units` field of our entries.

A modification to the above tran code definition shows how this parameterization would work:

{% partial file="graphql/examples/tranCodeDocs/achCredit.md" variables={title: "Example"} /%}

Note the change to `params.amount` inside of the `units` fields. Because these fields accept [CEL expressions](/reference/cel), we can reference fields on the runtime `params` object to access the values passed in. You can see the use of `params` in action: [Tran Code Invocation](/reference/ledger/tran-codes#tran-code-invocation).

In this way, the tran code can accept any number of parameterized values at runtime an inject them into the [Transaction]($gql:object:Transaction) and [Entries]($gql:object:Entry) written.
