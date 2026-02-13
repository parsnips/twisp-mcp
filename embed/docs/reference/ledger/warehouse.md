---
title: Warehouse
metaTitle: "Reference: Warehouse"
description: Reference for the data warehouse of Twisp ledgers.
showNext: true
showPrev: true
---

## The Basics

The Twisp ledger exports all data to AWS Redshift and exposes APIs to interact with the data using arbitrary SQL queries.  When enabled on your tenant, this warehouse provides a convenient way to execute analytical queries or unload data to your own data lake or warehouse.

{% callout type="note" %}
**Tip:** If you need to export data in bulk, you can use the [export API](/tutorials/advanced/export-data) without having warehouse enabled.
{% /callout %}

## Components of Warehouse

### GraphQL Interface

Twisp exposes the [Redshift Data API](https://docs.aws.amazon.com/redshift-data/latest/APIReference/Welcome.html) in GraphQL format for execution of SQL queries. The most commonly used API interactions are described here.

#### Execute Statement

The execute statement API allows you to execute a statement on Redshift.

```graphql
mutation ExecuteStatement($SQL:String! = "SELECT TRUE") {
  warehouse{
    executeStatement(input:{
      SQL: $SQL
    }) {
      id
    }
  }
}
```

#### Describe Statement

Describe Statement allows you to check on the status of a running query and retrieve metadata about it's results.

```graphql
query DescribeStatement($id:String!) {
  warehouse{
    describeStatement(input:{
      id:$id
    }) {
      resultRows
      resultSize
      status
      error
    }
  }
}

```

#### Cancel Statement

Cancel statement cancels a running statement.

```graphql
mutation CancelStatement($id:String!) {
  warehouse{
    cancelStatement(
      input:{
        id
      }
    )
  }
}
```

#### Get Statement Result

Get Statement Result allows you to fetch the results of a query via the API with paging.  If the result is particularly large, you may opt to instead use the [UNLOAD capabilty](https://docs.aws.amazon.com/redshift/latest/dg/r_UNLOAD.html) to export data in a suitable format to an S3 bucket of your choosing.

```graphql
query GetStatementResult($id:String!) {
  warehouse{
    getStatementResult(input:{
      id:$id
      #nextToken:""
    }) {
      records {
        fields {
          type
          value {
            str
            bytes
            isNull
          }
        }
      }
      nextToken
    }
  }
}
```
### Views and Schemas

Twisp maintains two views for each entity in the ledger:

- **`entity`**: the latest version of any entity stored in the Twisp ledger.
- **`entity_history`**: all versions of a particular record stored in the Twisp ledger.

Each record contains a number of header rows prepended with `record_`:

- `record_begin` - A unique timestamp of when the Twisp database transaction began.
- `record_rowid` - A uuid identifier of the particular item in Twisp.
- `record_status` - An enumeration indicating whether this record is deleted.
- `record_tenantid` - A uuid indicating the tenant that this record belongs to.
- `record_version` - An monotonically increasing unsigned integer indicating the version of the record.

These columns allow for incremental exports via `record_begin`, which are guaranteed to be unique per transaction. In other words, records with the same timestamp originated within the same transaction.  The `record_rowid` and `record_version` together uniquely identify a particular version of a record and is a convenient way to de-duplicate rows in the event you use overlapping timestamps to import data.

The following entities are available in the warehouse:

- `public.account`
- `public.account_history`
- `public.account_set`
- `public.account_set_history`
- `public.account_set_member`
- `public.account_set_member_history`
- `public.balance`
- `public.balance_history`
- `public.calculation`
- `public.calculation_history`
- `public.entry`
- `public.entry_history`
- `public.journal`
- `public.journal_history`
- `public.tran_code`
- `public.tran_code_history`
- `public.transaction`
- `public.transaction_history`
- `public.transaction_exception`
- `public.transaction_exception_history`
- `public.workflow_execution`
- `public.workflow_execution_history`

The schemas between each `entity` and `entity_history` view are identical and are provided below.

#### account

Every Account and Account Set in Twisp has an account record stored.

| Column Name | Type |
|----------------------------------|-------------|
| `record_begin`| `TIMESTAMP`|
| `record_rowid`| `VARCHAR(36)`|
| `record_status`| `VARCHAR`|
| `record_tenantid`| `VARCHAR(36)`|
| `record_version`| `BIGINT`|
| `account_id`| `VARCHAR(36)`|
| `status`| `VARCHAR`|
| `name`| `SUPER`|
| `code`| `VARCHAR`|
| `normal_balance_type`| `VARCHAR`|
| `description`| `SUPER`|
| `metadata`| `SUPER`|
| `created`| `TIMESTAMP`|
| `modified`| `TIMESTAMP`|
| `external_id`| `VARCHAR`|
| `config_enable_concurrent_posting`| `BOOLEAN`|
| `config_is_account_set`| `BOOLEAN`|

#### account_set

The `account_set` table has all Account Sets in the Twisp ledger.

|  Column Name                       | Type|
|----------------------------------|-------------|
| `record_begin`| `TIMESTAMP`|
| `record_rowid`| `VARCHAR(36)`|
| `record_status`| `VARCHAR`|
| `record_tenantid`| `VARCHAR(36)`|
| `record_version`| `BIGINT`|
| `account_set_id`| `VARCHAR(36)`|
| `journal_id`| `VARCHAR(36)`|
| `account_id`| `VARCHAR(36)`|
| `name`| `VARCHAR`|
| `description`| `VARCHAR`|
| `metadata`| `SUPER`|
| `created`| `TIMESTAMP`|
| `modified`| `TIMESTAMP`|
| `config_enable_concurrent_posting`| `BOOLEAN`|

#### account_set_member

The `account_set_member` table maintains the Account Set tree membership details.  The `member_id` can
refer to _either_ an Account Set or Account, based on `member_type`.

| Column Name | Type|
|-----------------|-------------|
| `record_begin`| `TIMESTAMP`|
| `record_rowid`| `VARCHAR(36)`|
| `record_status`| `VARCHAR`|
| `record_tenantid`| `VARCHAR(36)`|
| `record_version`| `BIGINT`|
| `account_set_id`| `VARCHAR(36)`|
| `journal_id`| `VARCHAR(36)`|
| `member_type`| `VARCHAR`|
| `member_id`| `VARCHAR(36)`|
| `created`| `TIMESTAMP`|

#### balance

Contains all Balance records and versions.  Note that the `dimension` column is of `VARBYTE` type.  If unloading,
you must convert into a base64 encoded string:

```SQL
SELECT FROM_VARBYTE(dimension,'base64') AS dimension FROM balance;
```

| Column Name | Type|
|----------------------------------|-------------|
| `record_begin`| `TIMESTAMP`|
| `record_rowid`| `VARCHAR(36)`|
| `record_status`| `VARCHAR`|
| `record_tenantid`| `VARCHAR(36)`|
| `record_version`| `BIGINT`|
| `journal_id`| `VARCHAR(36)`|
| `account_id`| `VARCHAR(36)`|
| `transaction_id`| `VARCHAR(36)`|
| `entry_id`| `VARCHAR(36)`|
| `currency`| `VARCHAR`|
| `settled_dr_balance`| `VARCHAR`|
| `settled_cr_balance`| `VARCHAR`|
| `settled_entry_id`| `VARCHAR(36)`|
| `settled_modified`| `TIMESTAMP`|
| `pending_dr_balance`| `VARCHAR`|
| `pending_cr_balance`| `VARCHAR`|
| `pending_entry_id`| `VARCHAR(36)`|
| `pending_modified`| `TIMESTAMP`|
| `encumbrance_dr_balance`| `VARCHAR`|
| `encumbrance_cr_balance`| `VARCHAR`|
| `encumbrance_entry_id`| `VARCHAR(36)`|
| `encumbrance_modified`| `TIMESTAMP`|
| `created`| `TIMESTAMP`|
| `modified`| `TIMESTAMP`|
| `available_settled_dr_balance`| `VARCHAR`|
| `available_settled_cr_balance`| `VARCHAR`|
| `available_settled_entry_id`| `VARCHAR(36)`|
| `available_settled_modified`| `TIMESTAMP`|
| `available_pending_dr_balance`| `VARCHAR`|
| `available_pending_cr_balance`| `VARCHAR`|
| `available_pending_entry_id`| `VARCHAR(36)`|
| `available_pending_modified`| `TIMESTAMP`|
| `available_encumbrance_dr_balance`| `VARCHAR`|
| `available_encumbrance_cr_balance`| `VARCHAR`|
| `available_encumbrance_entry_id`| `VARCHAR(36)`|
| `available_encumbrance_modified`| `TIMESTAMP`|
| `calculation_id`| `VARCHAR(36)`|
| `dimension`| `VARBYTE`|
| `dimensions`| `SUPER`|
| `entry_committed`| `TIMESTAMP`|
| `entry_timestamps`| `SUPER`|

#### calculation

Contains the calculation definitions in the Ledger.

|  Column Name       | Type|
|------------------|-------------|
| `record_begin`| `TIMESTAMP`|
| `record_rowid`| `VARCHAR(36)`|
| `record_status`| `VARCHAR`|
| `record_tenantid`| `VARCHAR(36)`|
| `record_version`| `BIGINT`|
| `calculation_id`| `VARCHAR(36)`|
| `description`| `VARCHAR`|
| `code`| `VARCHAR`|
| `dimensions`| `SUPER`|
| `status`| `VARCHAR`|
| `created`| `TIMESTAMP`|
| `modified`| `TIMESTAMP`|
| `condition`| `SUPER`|
| `skip_calculation`| `BOOLEAN`|
| `backfill_status`| `VARCHAR`|

#### entry

All of the journal entries in the Ledger.

|  Column Name             | Type|
|------------------------|-------------|
| `record_begin`| `TIMESTAMP`|
| `record_rowid`| `VARCHAR(36)`|
| `record_status`| `VARCHAR`|
| `record_tenantid`| `VARCHAR(36)`|
| `record_version`| `BIGINT`|
| `entry_id`| `VARCHAR(36)`|
| `transaction_id`| `VARCHAR(36)`|
| `transaction_seq`| `BIGINT`|
| `journal_id`| `VARCHAR(36)`|
| `account_id`| `VARCHAR(36)`|
| `entry_type`| `VARCHAR`|
| `layer`| `VARCHAR`|
| `direction`| `VARCHAR`|
| `description`| `VARCHAR`|
| `amount`| `VARCHAR`|
| `balance_record_id`| `VARCHAR(36)`|
| `balance_record_version`| `BIGINT`|
| `created`| `TIMESTAMP`|
| `modified`| `TIMESTAMP`|
| `metadata`| `SUPER`|
| `committed`| `TIMESTAMP`|
| `parent_account_ids`| `SUPER`|
| `is_voided_entry`| `BOOLEAN`|
| `is_void_entry`| `BOOLEAN`|

#### journal

All of the journals in the Ledger.

|  Column Name                       | Type|
|----------------------------------|-------------|
| `record_begin`| `TIMESTAMP`|
| `record_rowid`| `VARCHAR(36)`|
| `record_status`| `VARCHAR`|
| `record_tenantid`| `VARCHAR(36)`|
| `record_version`| `BIGINT`|
| `journal_id`| `VARCHAR(36)`|
| `name`| `VARCHAR`|
| `description`| `VARCHAR`|
| `status`| `VARCHAR`|
| `created`| `TIMESTAMP`|
| `modified`| `TIMESTAMP`|
| `code`| `VARCHAR`|
| `config_enable_effective_balances`| `BOOLEAN`|

#### tran_code

All of the Transaction Code definitions in the Ledger.

|  Column Name                 | Type|
|----------------------------|-------------|
| `record_begin`| `TIMESTAMP`|
| `record_rowid`| `VARCHAR(36)`|
| `record_status`| `VARCHAR`|
| `record_tenantid`| `VARCHAR(36)`|
| `record_version`| `BIGINT`|
| `tran_code_id`| `VARCHAR(36)`|
| `code`| `VARCHAR`|
| `description`| `VARCHAR`|
| `status`| `VARCHAR`|
| `params`| `SUPER`|
| `transaction_journal_id`| `SUPER`|
| `transaction_correlation_id`| `SUPER`|
| `transaction_external_id`| `SUPER`|
| `transaction_effective`| `SUPER`|
| `transaction_description`| `SUPER`|
| `transaction_metadata`| `SUPER`|
| `entries`| `SUPER`|
| `created`| `TIMESTAMP`|
| `modified`| `TIMESTAMP`|
| `metadata`| `SUPER`|

#### transaction

All of the Transactions posted in the Ledger.

|  Column Name        | Type|
|-------------------|-------------|
| `record_begin`| `TIMESTAMP`|
| `record_rowid`| `VARCHAR(36)`|
| `record_status`| `VARCHAR`|
| `record_tenantid`| `VARCHAR(36)`|
| `record_version`| `BIGINT`|
| `transaction_id`| `VARCHAR(36)`|
| `tran_code_id`| `VARCHAR(36)`|
| `journal_id`| `VARCHAR(36)`|
| `correlation_id`| `VARCHAR`|
| `external_id`| `VARCHAR`|
| `effective`| `DATE`|
| `description`| `VARCHAR`|
| `metadata`| `SUPER`|
| `created`| `TIMESTAMP`|
| `modified`| `TIMESTAMP`|
| `tran_code_version`| `BIGINT`|
| `void_of`| `VARCHAR(36)`|
| `voided_by`| `VARCHAR(36)`|

#### transaction_exception

All Transaction Exceptions for `WARN` and `VOID` velocity control enforcements.

| Column Name | Type|
|-----------------|-------------|
| `record_begin`| `TIMESTAMP`|
| `record_rowid`| `VARCHAR(36)`|
| `record_status`| `VARCHAR`|
| `record_tenantid`| `VARCHAR(36)`|
| `record_version`| `BIGINT`|
| `transaction_id`| `VARCHAR(36)`|
| `type`| `VARCHAR`|
| `error_message`| `VARCHAR`|
| `detail`| `SUPER`|
| `created`| `TIMESTAMP`|

#### velocity_control

All Velocity Controls defined in the ledger.

|  Column Name          | Type|
|---------------------|-------------|
| `record_begin`| `TIMESTAMP`|
| `record_rowid`| `VARCHAR(36)`|
| `record_status`| `VARCHAR`|
| `record_tenantid`| `VARCHAR(36)`|
| `record_version`| `BIGINT`|
| `velocity_control_id`| `VARCHAR(36)`|
| `name`| `VARCHAR`|
| `description`| `VARCHAR`|
| `enforcement`| `SUPER`|
| `condition`| `VARCHAR`|
| `created`| `TIMESTAMP`|
| `modified`| `TIMESTAMP`|

#### velocity_limit

All Velocity Limits defined in the Ledger.

|  Column Name        | Type|
|-------------------|-------------|
| `record_begin`| `TIMESTAMP`|
| `record_rowid`| `VARCHAR(36)`|
| `record_status`| `VARCHAR`|
| `record_tenantid`| `VARCHAR(36)`|
| `record_version`| `BIGINT`|
| `velocity_limit_id`| `VARCHAR(36)`|
| `name`| `VARCHAR`|
| `description`| `VARCHAR`|
| `window`| `SUPER`|
| `condition`| `SUPER`|
| `currency`| `VARCHAR`|
| `timestamp_source`| `VARCHAR`|
| `limit`| `SUPER`|
| `params`| `SUPER`|
| `created`| `TIMESTAMP`|
| `modified`| `TIMESTAMP`|

#### workflow_execution

A log of all workflow executions for `workflow.execute`.

| Column Name | Type|
|-----------------|-------------|
| `record_begin`| `TIMESTAMP`|
| `record_rowid`| `VARCHAR(36)`|
| `record_status`| `VARCHAR`|
| `record_tenantid`| `VARCHAR(36)`|
| `record_version`| `BIGINT`|
| `workflow_id`| `VARCHAR(36)`|
| `execution_id`| `VARCHAR(36)`|
| `task`| `VARCHAR`|
| `params`| `SUPER`|
| `context`| `SUPER`|
| `created`| `TIMESTAMP`|
| `modified`| `TIMESTAMP`|
| `output`| `SUPER`|

## Warehouse Operations

Use GraphQL queries and mutations to execute queries and describe results in the data warehouse:

- [`Query.warehouse.describeStatement()`]($gql:query:warehouse.describeStatement): Describe the status of executing query
- [`Query.warehouse.describeTable()`]($gql:query:warehouse.describeTable): Describe a table.
- [`Query.warehouse.getStatementResult()`]($gql:query:warehouse.getStatementResult): Retrieve the result of a table.
- [`Query.warehouse.listDatabases()`]($gql:query:warehouse.listDatabases): List databases available to query.
- [`Query.warehouse.listSchemas()`]($gql:query:warehouse.listSchemas): List schemas available in database.
- [`Query.warehouse.listStatements()`]($gql:query:warehouse.listStatements): List recent executing statements.
- [`Query.warehouse.listTables()`]($gql:query:warehouse.listTables): List tables in schema.
- [`Mutation.warehouse.executeStatementSync()`]($gql:mutation:warehouse.executeStatementSync): Execute a SQL statement and synchronously wait for result.
- [`Mutation.warehouse.batchExecuteStatement()`]($gql:mutation:warehouse.batchExecuteStatement): Execute a number of SQL statements asynchronously.
- [`Mutation.warehouse.cancelStatement()`]($gql:mutation:warehouse.cancelStatement): Cancel a running statement.
- [`Mutation.warehouse.executeStatement()`]($gql:mutation:warehouse.executeStatement): Execute a single SQL statement asynchronously.
