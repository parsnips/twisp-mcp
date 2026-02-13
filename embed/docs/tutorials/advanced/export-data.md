---
title: Exporting Data with Warehouse
metaTitle: "Tutorial: Exporting Data with Warehouse"
description: |
    Learn how to export your Twisp ledger data using the warehouse export API.
showNext: true
showPrev: true
---

Twisp’s **Warehouse Export** feature allows you to export ledger data into files you can download using the [Files API]($gql:mutation:files "Files"). This feature is available to all tenants, even if you do not have a managed data warehouse configured for your environment.

This tutorial will guide you through using the export API, explain key parameters, and share best practices for effective and efficient data exports.

## Why Use Warehouse Export?

- **Incremental & Historical Exports:** Export only new/changed records, or retrieve all historical versions.
- **Chunked, Compressed Exports:** Data is written in compressed files suitable for large-scale download and ingestion.
- **Supports All Ledger Data:** Export accounts, balances, entries, transactions, and more, in bulk.

## Example

To initiate an export for all balances for a specific date, use the following GraphQL mutation:

```graphql
mutation Export {
  warehouse {
    export(
      input: {
        entity: Balance
        version: LATEST
        format: JSON
        compression: GZIP
        destination: { files: { keyPrefix: "exports/2025/07/02/balance/" } }
        fromTimestamp: "2025-07-01T00:00:00.000Z"
        toTimestamp: "2025-07-02T00:00:00.000Z"
      }
    ) {
      id
    }
  }
}
```

This starts an export job and returns an `id` you can use to track its progress.

To check the status and list files created, use:

```graphql
query ExportStatus {
  warehouse {
    export(id: "YOUR_EXPORT_ID") {
      id
      status # (e.g., RUNNING, FINISHED, FAILED)
      error
    }
  }

  files {
    list(keyPrefix: "exports/2025/07/02/balance/")
  }
}
```

When finished, the files (e.g., `balance/000.json.gz`, `balance/001.json.gz`, etc.) are available for download through the Files API.

---

## Key Parameters

| Parameter          | Description                                                                                      |
|--------------------|--------------------------------------------------------------------------------------------------|
| `entity`           | The type of data to export (`Balance`, `Entry`, `Account`, etc.).                                |
| `version`          | `LATEST` exports only the most recent version per record; `HISTORY` exports all record versions. |
| `format`           | Format for the file (`JSON` is recommended for most usage).                                      |
| `compression`      | Output compression (`GZIP` recommended for efficient transfer).                                  |
| `destination`      | Where exported files are stored in Twisp Files (e.g., `{ files: { keyPrefix: "my/path" } }`).    |
| `fromTimestamp`    | (Optional) Start export with records created/modified after this timestamp (inclusive).          |
| `toTimestamp`      | (Optional) End export with records before this timestamp (exclusive).                            |

{% callout type="note" %}
**Tip:** Use `fromTimestamp` and `toTimestamp` for incremental exports—ideal for extracting only the new or modified data since your last export. This enables efficient change data capture (CDC) workflows for data lakes and reporting pipelines.
{% /callout %}

## Understanding LATEST vs HISTORY

- **`version: LATEST`**
  Exports just the most current version of every record in the chosen entity.
  Use this for up-to-date snapshots of your ledger (e.g., current balances, current state of accounts).

- **`version: HISTORY`**
  Exports all versions for every record.
  This is useful for:
  - Auditing
  - Reconstructing point-in-time states (e.g., balance history, account history)
  - Regulatory compliance

{% callout type="note" %}
To pull **point-in-time balances or entries** for reporting, use `version: HISTORY` and filter using `fromTimestamp` and `toTimestamp` to get just the versions active during a desired period.
{% /callout %}

## Output Schema Details

- **File Schema:** The exported data files follow the [Warehouse entity schemas](/reference/ledger/warehouse#views-and-schemas) (e.g., `balance`, `entry`), so you can use the same field definitions.
- **Amounts:** For types like `balance` and `entry`, number fields (such as monetary amounts) will include a stand-alone `<amount>_units` column representing the value as a number for convenience.

## Downloading Exported Files

Once the export job has completed:
1. Use the `files.list` API to enumerate all files at your chosen prefix.
2. Download each file with the `files.createDownload` mutation to receive a presigned URL.

## Additional Tips

- **Large Exports:** For very large exports, the system will automatically shard data into multiple files (e.g., `balance/000.json.gz`, `balance/001.json.gz`, etc.).
- **Performance:** Using compressed (`GZIP`) export is highly recommended for both speed and cost savings.
