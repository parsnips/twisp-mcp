---
title: CEL Variables Reference
metaTitle: "Reference: CEL Variables"
description: "Context variables available in CEL expressions across Twisp APIs."
showNext: false
showPrev: false
showToc: true
---

Twisp uses [Common Expression Language](https://cel.dev/) extensively for computation across the API. Different contexts expose different variables depending on their purpose. This reference documents what variables are available in each context.

## Overview

| Context                       | Primary Variable | Description                     |
|-------------------------------|------------------|---------------------------------|
| [Indexes](#indexes)           | `document`       | The record being indexed        |
| [Calculations](#calculations) | `context.vars`   | Transaction and account context |
| [Tran Codes](#tran-codes)     | `params`         | Post-time parameters            |

---

## Indexes

Custom indexes use CEL expressions to define partition keys, sort keys, and filter constraints. In all index expressions, the primary variable is `document`, which represents the record being indexed.

### Available Variables

| Variable   | Type             | Description              |
|------------|------------------|--------------------------|
| `document` | Protobuf message | The entity being indexed |

The type of `document` depends on the `on` field specified when creating the index:

- `ACCOUNT` → [Account](#account-fields)
- `ACCOUNT_SET` → [AccountSet](#account-set-fields)
- `TRANSACTION` → [Transaction](#transaction-fields)
- `TRANCODE` → [TranCode](#trancode-fields)
- `BALANCE` → [Balance](#balance-fields)
- `ENTRY` → [Entry](#entry-fields)

### Partition Expressions

Partition keys determine how records are grouped in the index.

```graphql
mutation CreatePartitionedIndex {
  schema {
    createIndex(
      input: {
        name: "entries_by_account"
        on: Entry
        partition: [
          { alias: "accountId", value: "document.account_id" }
          { alias: "journalId", value: "document.journal_id" }
        ]
      }
    ) { name }
  }
}
```

Partition expressions can return arrays to create multiple index entries. For example this create multiple partitions for the account id and each of it's parent account sets (at time of posting):

```graphql
partition: [
  {
    alias: "accountId",
    value: "document.parent_account_ids + [document.account_id]"
  }
]
```

### Sort Expressions

Sort keys define the ordering within each partition. In this example sorting either by th

```graphql
mutation CreateSortedIndex {
  schema {
    createIndex(
      input: {
        name: "entries_by_effective"
        on: Entry
        partition: [{ alias: "accountId", value: "document.account_id" }]
        sort: [
          {
            alias: "effective"
            value: "'effective' in document.metadata && document.metadata.effective != null ? document.metadata.effective : document.created"
            type: STRING
            sort: DESC
          }
        ]
      }
    ) { name }
  }
}
```

### Filter Constraints

Constraints are boolean expressions that must all evaluate to `true` for a record to be included in the index.

```graphql
mutation CreateFilteredIndex {
  schema {
    createIndex(
      input: {
        name: "active_entries"
        on: Entry
        partition: [{ alias: "accountId", value: "document.account_id" }]
        constraints: {
          isNotVoidEntry: "!document.is_void_entry"
          isNotVoidedEntry: "!document.is_voided_entry"
          isSettled: "document.layer == SETTLED" 
        }
      }
    ) { name }
  }
}
```

---

## Calculations

Calculations define how balances are computed and grouped. CEL expressions are used in dimension definitions and conditions.

### Available Variables

| Variable | Type | Description |
|----------|------|-------------|
| `context.vars.transaction` | [Transaction](#transaction-fields) | The parent transaction of the entry |
| `context.vars.account` | [Account](#account-fields) | The account the entry is posted to |
| `document` | [Entry](#entry-fields) | The entry being evaluated (in conditions) |

### Dimension Expressions

Dimensions define how balances are grouped. Common use cases include grouping by date, metadata values, or account properties.

```graphql
mutation CreateEffectiveDateCalculation {
  createCalculation(
    input: {
      calculationId: "5867b5dd-fc69-416c-80f5-62e8a53610d5"
      code: "EFFECTIVE_DATE"
      description: "Track balances per effective date."
      dimensions: [
        {
          alias: "effectiveDate"
          value: "context.vars.transaction.effective"
        }
      ]
      scope: LOCAL
    }
  ) {
    calculationId
    code
  }
}
```

Multi-dimensional calculations:

```graphql
mutation CreateMultiDimensionCalculation {
  createCalculation(
    input: {
      calculationId: "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
      code: "BY_CURRENCY_AND_DATE"
      description: "Track balances by currency and effective date."
      dimensions: [
        {
          alias: "currency"
          value: "context.vars.account.metadata.currency"
        }
        {
          alias: "effectiveDate"
          value: "context.vars.transaction.effective"
        }
      ]
      scope: LOCAL
    }
  ) { calculationId }
}
```

### Condition Expressions

Conditions filter which entries are included in the calculation. They must evaluate to a boolean.

```graphql
mutation CreateConditionalCalculation {
  createCalculation(
    input: {
      calculationId: "b2c3d4e5-f6a7-8901-bcde-f23456789012"
      code: "POLICY_PAYMENTS"
      description: "Track only entries to accounts with policy payment metadata."
      dimensions: [
        { alias: "effectiveDate", value: "context.vars.transaction.effective" }
      ]
      condition: "has(context.vars.account.metadata.policyPayment)"
      scope: LOCAL
    }
  ) { calculationId }
}
```

Using entry fields in conditions:

```graphql
condition: "document.layer == 0 && document.direction == 0"
```

---

## Tran Codes

Tran codes are templates that define how transactions and entries are created. CEL expressions are used extensively to dynamically compute values at post time.

### Available Variables

| Variable | Type | Description |
|----------|------|-------------|
| `params` | Struct | Parameters passed when posting the transaction |
| `vars` | Struct | Calculated variables defined in the tran code's `vars` field |

### Parameter Access

Parameters are defined in the tran code and passed at post time. Access them with `params.<paramName>`:

```graphql
mutation CreateTranCode {
  createTranCode(
    input: {
      code: "TRANSFER"
      description: "Transfer funds between accounts"
      params: [
        { name: "amount", type: TYPE_DECIMAL, description: "Transfer amount" }
        { name: "fromAccount", type: TYPE_UUID, description: "Source account" }
        { name: "toAccount", type: TYPE_UUID, description: "Destination account" }
        { name: "memo", type: TYPE_STRING, default: "'Transfer'", description: "Optional memo" }
      ]
      transaction: {
        description: "'Transfer ' + string(params.amount) + ' - ' + params.memo"
      }
      entries: [
        {
          entryType: "'TRANSFER_DR'"
          accountId: "params.fromAccount"
          direction: "'DEBIT'"
          units: "params.amount"
          currency: "'USD'"
        }
        {
          entryType: "'TRANSFER_CR'"
          accountId: "params.toAccount"
          direction: "'CREDIT'"
          units: "params.amount"
          currency: "'USD'"
        }
      ]
    }
  ) { code }
}
```

### Vars (Scratch Pad)

The `vars` field provides a scratch pad for intermediate calculations. Variables defined here can be referenced in transaction and entry expressions:

```graphql
mutation CreateTranCodeWithVars {
  createTranCode(
    input: {
      code: "TRANSFER_WITH_FEE"
      description: "Transfer with calculated fee"
      params: [
        { name: "amount", type: TYPE_DECIMAL }
        { name: "fromAccount", type: TYPE_UUID }
        { name: "toAccount", type: TYPE_UUID }
        { name: "feeAccount", type: TYPE_UUID }
      ]
      vars: {
        feeRate: "'0.025'"
        feeAmount: "decimal(params.amount) * decimal(vars.feeRate)"
        netAmount: "decimal(params.amount) - vars.feeAmount"
      }
      entries: [
        {
          entryType: "'TRANSFER_DR'"
          accountId: "params.fromAccount"
          direction: "'DEBIT'"
          units: "params.amount"
          currency: "'USD'"
        }
        {
          entryType: "'TRANSFER_CR'"
          accountId: "params.toAccount"
          direction: "'CREDIT'"
          units: "vars.netAmount"
          currency: "'USD'"
        }
        {
          entryType: "'FEE_CR'"
          accountId: "params.feeAccount"
          direction: "'CREDIT'"
          units: "vars.feeAmount"
          currency: "'USD'"
        }
      ]
    }
  ) { code }
}
```

### Entry Conditions

Entry conditions determine whether an entry should be created. Useful for optional entries:

```graphql
entries: [
  {
    entryType: "'FEE'"
    accountId: "params.feeAccount"
    direction: "'CREDIT'"
    units: "params.feeAmount"
    currency: "'USD'"
    condition: "params.feeAmount > decimal('0.00')"
  }
]
```

### Transaction Fields

All fields in the `transaction` block accept CEL expressions:

| Field | Expected Type | Example |
|-------|--------------|---------|
| `journalId` | UUID | `"uuid('b28f5684-0834-4292-8016-d2f2fb0367a9')"` |
| `correlationId` | String | `"params.correlationId"` |
| `externalId` | String | `"params.externalId"` |
| `effective` | Date | `"date('2024-01-15')"` or `"params.effectiveDate"` |
| `description` | String | `"'Payment: ' + string(params.amount)"` |
| `metadata` | JSON | `"{ 'source': 'api', 'amount': string(params.amount) }"` |

### Entry Fields

All fields in each entry accept CEL expressions:

| Field | Expected Type | Example |
|-------|--------------|---------|
| `entryType` | String | `"'ACH_CREDIT'"` |
| `accountId` | UUID | `"params.accountId"` |
| `layer` | Enumeration | `"SETTLED"` or `"PENDING"` or `"ENCUMBRANCE"` |
| `direction` | Enumeration | `"DEBIT"` or `"CREDIT"` |
| `units` | Decimal | `"params.amount"` or `"decimal('100.00')"` |
| `currency` | String | `"'USD'"` or `"params.currency"` |
| `description` | String | `"'Entry for ' + params.memo"` |
| `metadata` | CEL Map/JSON | `"{ 'lineItem': params.lineItem }"` |
| `condition` | Boolean | `"params.amount > decimal('0')"` |

---

## Entity Field Reference

### Account Fields

| Field | Type | Description |
|-------|------|-------------|
| `account_id` | UUID | Unique identifier |
| `status` | Enum | ACTIVE, LOCKED, or INACTIVE |
| `name` | String | Account name |
| `code` | String | Shorthand code |
| `normal_balance_type` | Enum | DEBIT (0) or CREDIT (1) |
| `description` | String | Account description |
| `metadata` | Struct | Arbitrary JSON data |
| `created` | Timestamp | Creation time |
| `modified` | Timestamp | Last modification time |
| `external_id` | String | External system identifier |

### Account Set Fields

| Field | Type | Description |
|-------|------|-------------|
| `account_set_id` | UUID | Unique identifier |
| `journal_id` | UUID | Associated journal |
| `account_id` | UUID | Associated account ID |
| `name` | String | Set name |
| `description` | String | Set description |
| `metadata` | Struct | Arbitrary JSON data |
| `created` | Timestamp | Creation time |
| `modified` | Timestamp | Last modification time |
| `code` | String | Unique code |
| `has_members` | Boolean | Whether set has members |

### Transaction Fields

| Field | Type | Description |
|-------|------|-------------|
| `transaction_id` | UUID | Unique identifier |
| `tran_code_id` | UUID | Associated tran code |
| `journal_id` | UUID | Associated journal |
| `correlation_id` | String | Groups related transactions |
| `external_id` | String | External system identifier |
| `effective` | Date | Accounting effective date |
| `description` | String | Transaction description |
| `metadata` | Struct | Arbitrary JSON data |
| `created` | Timestamp | Creation time |
| `modified` | Timestamp | Last modification time |

### Entry Fields

| Field | Type | Description |
|-------|------|-------------|
| `entry_id` | UUID | Unique identifier |
| `transaction_id` | UUID | Parent transaction |
| `journal_id` | UUID | Associated journal |
| `account_id` | UUID | Target account |
| `entry_type` | String | Entry type code |
| `layer` | Enum | SETTLED (0), PENDING (1), or ENCUMBRANCE (2) |
| `direction` | Enum | DEBIT (0) or CREDIT (1) |
| `description` | String | Entry description |
| `amount` | Money | Entry amount with currency |
| `metadata` | Struct | Arbitrary JSON data |
| `created` | Timestamp | Creation time |
| `modified` | Timestamp | Last modification time |
| `parent_account_ids` | List[UUID] | Parent account set IDs |
| `is_void_entry` | Boolean | True if this is a voiding entry |
| `is_voided_entry` | Boolean | True if this entry was voided |

### TranCode Fields

| Field | Type | Description |
|-------|------|-------------|
| `tran_code_id` | UUID | Unique identifier |
| `code` | String | Unique code identifier |
| `description` | String | Tran code description |
| `status` | Enum | ACTIVE, LOCKED, or INACTIVE |
| `params` | List | Parameter definitions |
| `transaction` | Struct | Transaction template |
| `entries` | List | Entry templates |
| `metadata` | Struct | Arbitrary JSON data |
| `created` | Timestamp | Creation time |
| `modified` | Timestamp | Last modification time |

### Balance Fields

| Field | Type | Description |
|-------|------|-------------|
| `journal_id` | UUID | Associated journal |
| `account_id` | UUID | Associated account |
| `transaction_id` | UUID | Last transaction |
| `entry_id` | UUID | Last entry |
| `currency` | String | Currency code |
| `settled` | BalanceAmount | Settled layer balance |
| `pending` | BalanceAmount | Pending layer balance |
| `encumbrance` | BalanceAmount | Encumbrance layer balance |
| `created` | Timestamp | Creation time |
| `modified` | Timestamp | Last modification time |
| `calculation_id` | UUID | Associated calculation |
| `dimensions` | Struct | Dimension values |

---

## Type Constructors

When working with CEL expressions, use these constructors to create typed values:

| Constructor | Example | Description |
|------------|---------|-------------|
| `uuid(string)` | `uuid('a1b2c3d4-...')` | Parse UUID from string |
| `decimal(string)` | `decimal('100.50')` | High-precision decimal |
| `money(string, string)` | `money('100.00', 'USD')` | Money with currency |
| `date(string)` | `date('2024-01-15')` | Parse date (YYYY-MM-DD) |
| `timestamp(string)` | `timestamp('2024-01-15T10:30:00Z')` | Parse ISO 8601 timestamp |
| `bool(string)` | `bool('true')` | Parse boolean |
| `int(string)` | `int('42')` | Parse integer |
| `string(any)` | `string(params.amount)` | Convert to string |

For a complete list of functions, see the [CEL Functions Reference](/reference/cel/examples).
