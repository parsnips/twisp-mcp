---
title: File Operations
metaTitle: "Reference: ACH File Operations"
description: Reference for ACH file processing within the Twisp ACH Processor
showNext: true
showPrev: true
---

## Overview

ACH file operations handle the complete lifecycle of NACHA-formatted files within the Twisp platform. This reference describes how Twisp processes ACH files: from upload and validation through parallel processing and generation. Understanding these internals helps you optimize file processing, debug issues, and integrate effectively with the ACH system.

The Twisp ACH file processor:
- Validates NACHA file format compliance with detailed error reporting
- Partitions files for parallel transaction processing
- Generates NACHA-compliant files from workflow transactions
- Provides presigned URLs for secure file upload and download with expiration
- Tracks file processing through multiple states with real-time statistics
- Maintains complete file processing history and audit trails
- Stores files in S3 with retention policies and versioning

## Components of File Operations

There are 4 primary components in ACH file processing:

1. **File Upload**: Create presigned URLs for secure file uploads, supporting ACH and bulk GraphQL variable files.
2. **File Processing**: Parse NACHA files, validate format compliance, extract transactions, and trigger webhook decisioning for each entry.
3. **File Generation**: Create NACHA-formatted files from submitted transactions, ready for transmission to financial institutions.
4. **File Download**: Generate presigned URLs for secure file downloads with automatic expiration.

In addition, file records maintain:

- **File ID**: UUID derived from file key for tracking specific file processing instances.
- **Processing Status**: Current state (NEW, VALIDATING, PROCESSING, COMPLETED, ERROR, etc.).
- **Processing Statistics**: Entry counts, credit/debit totals, unprocessed entry counts.
- **Configuration Version**: Which ACH configuration version was used for processing.
- **History**: Complete audit trail of all status changes and statistic updates.

## File Upload

File upload uses a two-step process: first create a presigned URL, then upload file content via HTTP PUT.

**Upload Types:**
- `ACH`: NACHA-formatted text files for transaction processing
- `BULK_GRAPHQL_VARIABLES`: JSON arrays for bulk GraphQL query execution

**File Requirements for ACH:**
- Content type: `text/plain`
- Format: NACHA (National Automated Clearing House Association) standard
- Character encoding: ASCII
- Line endings: CRLF or LF
- Record length: 94 characters per line

Presigned upload URLs expire after a configured period. If expiration occurs before upload completes, create a new upload request.

## File Processing

File processing parses uploaded NACHA files, validates format compliance, and extracts transaction details for decisioning.

**File Types:**

RDFI (Receiving):
- `RDFI`: Incoming transactions to receive (credits and debits)
- `RDFI_RETURN`: Return files from ODFI for transactions you originated
- `RDFI_NOC`: Notification of Change files

ODFI (Originating):
- `ODFI`: Outgoing file with both PUSH and PULL transactions
- `ODFI_PULL_ONLY`: Outgoing file with only PULL (debit) transactions
- `ODFI_PUSH_ONLY`: Outgoing file with only PUSH (credit) transactions
- `ODFI_RETURN`: Outgoing return file for RDFI transactions you received
- `ODFI_PREPROCESS_RETURN`: Pre-process return file before final generation
- `ODFI_PROCESSED`: Already processed ODFI file for reconciliation

**Processing Phases:**

1. **Validation**: Verify NACHA format compliance, record structure, control totals
2. **Partitioning**: Split file for parallel transaction processing
3. **Webhook Decisioning**: Send each transaction to configured webhook endpoint
4. **Settlement**: Apply settlement instructions from webhook responses
5. **Completion**: Mark all transactions settled or queued for returns

**Preprocessing Options:**

When using `ODFI_PREPROCESS_RETURN` file type, specify output file keys for filtered results:
- `preprocessedFileKey`: Entries passing preprocessing rules
- `preprocessedExcludedFileKey`: Entries failing preprocessing rules

## File Generation

File generation creates NACHA-formatted files from submitted workflow transactions.

**Generation Options:**
- `generateEmpty: true`: Always create file, even without transactions (empty NACHA with headers)
- `generateEmpty: false`: Only create file if transactions exist (returns `generated: false` if none)

**File Type Selection:**

Use separate file types when:
- ODFI requires distinct files for credits vs debits
- Different transmission schedules for transaction types
- Risk management requires transaction type isolation
- Separate approval workflows needed

Use combined `ODFI` file type when:
- Processing both credits and debits together
- Single transmission window
- Financial institution accepts mixed transaction files

Generated files are NACHA-compliant and ready for transmission via SFTP, FTPS, or other secure transfer methods.

## File Download

File download provides presigned URLs for retrieving generated or processed files.

**Download Process:**
1. Request download URL via GraphQL mutation
2. Receive presigned URL with expiration timestamp
3. Download file content via HTTP GET before expiration

Presigned download URLs expire after configured period. If URL expires, create new download request with same file key.

## File Lifecycle Deep Dive

Understanding how Twisp processes ACH files helps you optimize integration, debug issues, and monitor progress effectively.

### Processing Phases

File processing progresses through multiple phases, each with specific responsibilities and state transitions.

#### Phase 1: Upload (NEW → UPLOADED)

**NEW State:**
- File record created via [`Mutation.ach.processFile()`]($gql:mutation:ach.processFile)
- File ID generated from file key
- Initial `AchFileInfo` record stored with configuration version
- Processing status set to `NEW`

**UPLOADED State:**
- File content uploaded to S3 via presigned URL from [`Mutation.files.createUpload()`]($gql:mutation:files.createUpload)
- System detects uploaded file and transitions to `UPLOADED`
- File queued for validation
- Processing statistics initialized to zero

**Transition Trigger:** S3 upload completion detected by file processor

#### Phase 2: Validation (UPLOADED → VALIDATING → PARTITIONING)

**VALIDATING State:**
- Parser reads NACHA file structure
- Validates file header record (Record Type 1)
- Validates batch header records (Record Type 5)
- Validates entry detail records (Record Type 6)
- Validates addenda records (Record Type 7)
- Validates control records (Record Types 8, 9)
- Checks hash totals and entry counts
- Verifies record length (94 characters)

**Validation Checks:**
- File header immediate destination/origin format
- Batch header service class codes
- Entry detail transaction codes
- Routing number check digit validation
- Control total reconciliation
- Record sequence validation

**On Success:** Transitions to `PARTITIONING`
**On Failure:** Transitions to `INVALID` with detailed error in `processingDetail`

#### Phase 3: Partitioning (PARTITIONING → PROCESSING)

**PARTITIONING State:**
- File split into partitions for parallel processing
- Each partition contains subset of entry detail records
- Partition size optimized for concurrent webhook processing
- ACH workflow trace records created for each entry
- Trace numbers extracted and indexed for return matching

**Partitioning Strategy:**
- Default: 100 entries per partition
- Configurable based on expected webhook latency
- Partitions processed independently by worker pool
- Enables horizontal scaling of webhook processing

**Transition Trigger:** All partitions created and queued

#### Phase 4: Webhook Processing (PROCESSING → PROCESSED)

**PROCESSING State (RDFI only):**
- Each entry triggers webhook to configured endpoint
- Webhooks sent in parallel across partitions
- Responses collected: `SETTLE`, `RETURN`, or `RETRY`
- Statistics updated: `numEntriesUnprocessed` decrements
- Retry logic for failed webhooks (exponential backoff)

**ODFI Files:**
- ODFI files skip webhook processing
- Transition directly to `PROCESSED` after partitioning

**Monitoring:**
Query `processingStatistics.numEntriesUnprocessed` via [`Query.ach.file()`]($gql:query:ach.file) to track progress:

```graphql
query MonitorFileProcessing {
  ach {
    file(fileKey: "incoming-20251114.ach", configId: "config-001") {
      processingStatus
      processingStatistics {
        numEntriesUnprocessed
        totalCreditAmount
        totalDebitAmount
      }
    }
  }
}
```

**Transition Trigger:** `numEntriesUnprocessed` reaches zero

#### Phase 5: Settlement (PROCESSED → COMPLETED)

**PROCESSED State:**
- All webhooks sent and responses received
- Transactions awaiting settlement based on effective dates
- For RDFI: Settlement occurs when workflow SETTLE state executes
- For ODFI: Settlement occurs immediately after file generation

**COMPLETED State:**
- All transactions settled or queued for returns
- File processing complete
- Final statistics recorded
- File available for download

**Transition Trigger:** All transactions reach terminal state (settled or returned)

### State Transition Diagram

```
┌─────────┐
│   NEW   │ File record created
└────┬────┘
     │ File uploaded
     ↓
┌─────────┐
│UPLOADED │ File content in S3
└────┬────┘
     │ Validation starts
     ↓
┌──────────┐    Validation fails    ┌─────────┐
│VALIDATING├────────────────────────→│ INVALID │ Terminal
└────┬─────┘                         └─────────┘
     │ Validation succeeds
     ↓
┌─────────────┐   Processing error   ┌────────┐
│PARTITIONING ├──────────────────────→│ ERROR  │ Terminal
└──────┬──────┘                       └────────┘
       │ Partitions created
       ↓
┌────────────┐   Manual abort        ┌─────────┐
│ PROCESSING ├───────────────────────→│ ABORTED │ Terminal
└──────┬─────┘                        └─────────┘
       │ All webhooks sent
       ↓
┌───────────┐
│ PROCESSED │ Awaiting settlements
└─────┬─────┘
      │ All settled/returned
      ↓
┌───────────┐
│ COMPLETED │ Terminal
└───────────┘
```

### Processing Status Reference

| Status | Description | Next States | Typical Duration |
|--------|-------------|-------------|------------------|
| `NEW` | File record created | `UPLOADED`, `ERROR` | < 1 second |
| `UPLOADED` | File in S3, queued | `VALIDATING` | < 5 seconds |
| `VALIDATING` | Format validation | `PARTITIONING`, `INVALID` | 5-30 seconds |
| `PARTITIONING` | Creating partitions | `PROCESSING`, `ERROR` | 5-15 seconds |
| `PROCESSING` | Sending webhooks | `PROCESSED`, `ABORTED` | Minutes to hours* |
| `PROCESSED` | Webhooks complete | `COMPLETED` | Hours to days** |
| `COMPLETED` | Final state | None | Permanent |
| `INVALID` | Validation failed | None | Permanent |
| `ERROR` | Processing error | None | Permanent |
| `ABORTED` | Manually stopped | None | Permanent |

\* Depends on webhook response time and retry logic
\** Depends on transaction effective dates and settlement timing

## Processing Statistics

Processing statistics track entry counts and monetary totals during file processing.

**Statistics Fields:**
- `numEntriesUnprocessed`: Entries awaiting webhook responses
- `totalCreditAmount`: Sum of all credit transactions (in cents)
- `totalDebitAmount`: Sum of all debit transactions (in cents)

Statistics update throughout processing lifecycle, providing real-time visibility into file processing progress.

## File Format Validation

Twisp validates NACHA file format compliance during the `VALIDATING` phase. Understanding validation helps you debug file issues and ensure compliance.

### What Twisp Validates

**File Structure:**
- Record length: Every line must be exactly 94 characters
- Record types: Must be 1, 5, 6, 7, 8, or 9
- Record sequence: File Header → Batch(es) → File Control
- Line endings: CRLF or LF accepted
- Character encoding: ASCII only

**File Header Record (Type 1):**
- Immediate destination: 10 characters (routing number or "          ")
- Immediate origin: 10 characters (Tax ID or routing number)
- File creation date/time: Valid YYMMDD and HHMM
- File ID modifier: Single character A-Z or 0-9

**Batch Header Record (Type 5):**
- Service class code: 200 (mixed), 220 (credits), 225 (debits)
- Company identification: 10 characters
- Standard entry class code: Valid SEC code (PPD, CCD, WEB, etc.)
- Effective entry date: Valid YYMMDD format
- Originator status code: 0, 1, or 2
- ODFI identification: 8-digit routing number

**Entry Detail Record (Type 6):**
- Transaction code: Valid code (22, 23, 24, 27, 28, 29, 32, 33, 34, 37, 38, 39)
- Receiving DFI identification: 8 digits
- Check digit: Matches routing number check digit algorithm
- DFI account number: 1-17 characters
- Amount: 10-digit number (cents)
- Trace number: 15 digits (ODFI routing + sequence)

**Batch Control Record (Type 8):**
- Entry count matches actual entries in batch
- Entry hash matches sum of RDFI routing numbers (rightmost 10 digits)
- Total debit amount matches sum of debit entries
- Total credit amount matches sum of credit entries

**File Control Record (Type 9):**
- Batch count matches actual batches
- Block count correct (records / 10, rounded up)
- Entry and addenda count matches total
- Entry hash matches sum of all RDFI routing numbers
- Total debit/credit amounts match sums across all batches

### Common Validation Errors

**Record Length Errors:**
```
Error: "Line 15: Record length is 93, expected 94"
Fix: Ensure all lines are exactly 94 characters (pad with spaces if needed)
```

**Check Digit Errors:**
```
Error: "Line 42: Check digit 7 does not match calculated value 3"
Fix: Recalculate check digit using ABA routing number algorithm
```

**Control Total Mismatch:**
```
Error: "Batch 1: Entry count 150 does not match control record value 148"
Fix: Verify all entries included in batch, regenerate control record
```

**Invalid Transaction Code:**
```
Error: "Line 67: Invalid transaction code 21"
Fix: Use valid codes (22-29 for checking, 32-39 for savings)
```

**Entry Hash Mismatch:**
```
Error: "File Control: Entry hash 1234567890 does not match calculated 1234567891"
Fix: Sum first 8 digits of all RDFI routing numbers, take rightmost 10 digits
```

### Debugging Validation Failures

When a file fails validation (status `INVALID`), check `processingDetail` for error information:

```graphql
query GetValidationError {
  ach {
    file(fileKey: "problematic-file.ach", configId: "config-001") {
      fileId
      processingStatus
      processingDetail
      modified
    }
  }
}
```

**Example `processingDetail` values:**
- `"Line 42: Record length is 93, expected 94"`
- `"Batch 1: Entry count mismatch (expected 150, got 148)"`
- `"File Control: Total debit amount mismatch"`

### Validation Best Practices

**Before Uploading:**
- Use NACHA file validation tools
- Verify record lengths (94 characters exactly)
- Check control totals match entry counts
- Validate routing number check digits
- Ensure valid transaction codes for account types

**After Validation Errors:**
- Read `processingDetail` for specific line/field errors
- Fix source data or file generation logic
- Re-upload corrected file with new file key
- Consider automated pre-validation before Twisp upload

## Parallel Processing Architecture

Twisp partitions large ACH files for concurrent webhook processing, enabling high-throughput file processing.

### How Partitioning Works

**Partition Creation (PARTITIONING state):**
1. File parsed and validated
2. Entry detail records extracted
3. Entries divided into partitions (default 100 per partition)
4. Each partition assigned to worker in pool
5. Workers process partitions concurrently

**Partition Processing:**
- Each worker processes its partition independently
- Webhooks sent in parallel across all workers
- No cross-partition dependencies
- Failures in one partition don't affect others

**Concurrency Model:**
```
File: 1,000 entries
├── Partition 1: Entries 1-100    → Worker 1
├── Partition 2: Entries 101-200  → Worker 2
├── Partition 3: Entries 201-300  → Worker 3
├── ...
└── Partition 10: Entries 901-1000 → Worker 10

Each worker sends 100 webhooks in parallel
Total: 1,000 webhooks sent concurrently across 10 workers
```

### Performance Characteristics

**Throughput:**
- Small files (< 100 entries): 30-60 seconds total
- Medium files (100-1,000 entries): 2-5 minutes with fast webhooks
- Large files (1,000-10,000 entries): 10-30 minutes with fast webhooks
- Very large files (> 10,000 entries): Linear scaling with partition count

**Bottlenecks:**
- Webhook response time (primary factor)
- Webhook endpoint throughput
- Database transaction commit latency
- Network bandwidth for webhook traffic

**Optimization Tips:**
- Optimize webhook endpoint for < 100ms response time
- Use async processing in webhook handler
- Enable connection pooling for database access
- Consider horizontal scaling of webhook endpoint
- Monitor `numEntriesUnprocessed` to track progress

### Monitoring Parallel Processing

Track processing progress in real-time via [`Query.ach.file()`]($gql:query:ach.file):

```graphql
query MonitorProcessing {
  ach {
    file(fileKey: "large-file.ach", configId: "config-001") {
      processingStatus
      processingStatistics {
        numEntriesUnprocessed
        totalCreditAmount
        totalDebitAmount
      }
      modified
    }
  }
}
```

Poll every 5-10 seconds during `PROCESSING` state to track `numEntriesUnprocessed` countdown.

## Storage and Retention

Twisp stores ACH files in AWS S3 with security, retention, and versioning policies.

### File Storage Location

**S3 Bucket Structure:**
```
twisp-ach-files-{region}/
├── uploads/
│   ├── {tenantId}/
│   │   └── {fileKey}
├── generated/
│   ├── {tenantId}/
│   │   └── {fileKey}
└── archived/
    ├── {tenantId}/
        └── {fileKey}
```

**Storage Classes:**
- Uploads: Standard S3 (frequent access)
- Generated files: Standard S3 (frequent access for 30 days)
- Archived files: Glacier after 90 days (compliance retention)

### Retention Policies

**Upload Files:**
- Retained for 90 days after upload
- Archived to Glacier for 7 years (NACHA compliance)
- Accessible via file key throughout retention period

**Generated Files:**
- Retained for 90 days in Standard S3
- Archived to Glacier for 7 years
- Download URLs valid for 15 minutes (renewable)

**File Metadata:**
- `AchFileInfo` records retained permanently
- Processing history retained for audit trail
- Statistics snapshots retained at each version
- Enables compliance reporting and reconciliation

### Versioning

**File Versioning:**
- S3 versioning enabled on ACH file buckets
- Protects against accidental deletion
- Enables recovery of overwritten files
- Version history retained for full retention period

**Metadata Versioning:**
- `AchFileInfo.version` increments on each update
- `AchFileInfo.history` connection provides all versions
- Each version captures processing state snapshot
- Enables point-in-time analysis

### Security

**Encryption:**
- S3 server-side encryption (SSE-S3) enabled
- Files encrypted at rest
- Presigned URLs use HTTPS only
- API access requires authentication

**Access Control:**
- Tenant-isolated S3 paths
- IAM roles restrict cross-tenant access
- Presigned URLs scoped to specific file key
- Time-limited URLs (15-minute expiration)

## Upload and Download Mechanics

Twisp uses presigned URLs for secure, direct S3 access without exposing AWS credentials.

### Upload Process

**Step 1: Create Upload URL**

Call [`Mutation.files.createUpload()`]($gql:mutation:files.createUpload) to get presigned URL:

```graphql
mutation CreateUpload {
  files {
    createUpload(
      input: {
        key: "incoming-20251114.ach"
        uploadType: ACH
        contentType: "text/plain"
      }
    ) {
      uploadURL
    }
  }
}
```

**Response:**
```json
{
  "files": {
    "createUpload": {
      "uploadURL": "https://s3.amazonaws.com/twisp-ach-files/.../incoming-20251114.ach?X-Amz-..."
    }
  }
}
```

**Step 2: Upload File Content**

Use HTTP PUT to upload file to presigned URL:

```bash
curl -T incoming-20251114.ach \
     -H "Content-Type: text/plain" \
     -XPUT '<uploadURL>'
```

**Important:**
- URL expires after 15 minutes
- Content-Type must match specified type
- Maximum file size: 100 MB
- Upload must complete before expiration

**Step 3: Process File**

After upload completes, start processing with [`Mutation.ach.processFile()`]($gql:mutation:ach.processFile):

```graphql
mutation ProcessFile {
  ach {
    processFile(
      input: {
        configId: "config-001"
        fileKey: "incoming-20251114.ach"
        fileType: RDFI
      }
    ) {
      fileId
    }
  }
}
```

### Download Process

**Step 1: Generate Download URL**

Call [`Mutation.files.createDownload()`]($gql:mutation:files.createDownload) to get presigned URL:

```graphql
mutation CreateDownload {
  files {
    createDownload(
      key: "outgoing-20251114.ach"
    ) {
      downloadURL
    }
  }
}
```

**Step 2: Download File Content**

Use HTTP GET to download file from presigned URL:

```bash
curl '<downloadURL>' -o outgoing-20251114.ach
```

**Important:**
- URL expires after 15 minutes
- If expired, generate new download URL
- No authentication required (URL contains signature)
- Download bandwidth not rate-limited

### Presigned URL Security Model

**How It Works:**
1. Twisp generates presigned URL with AWS STS
2. URL contains temporary credentials in query parameters
3. AWS validates signature on access
4. Access granted without exposing permanent credentials

**Security Properties:**
- Time-limited access (15-minute expiration)
- Scoped to specific file key
- Cannot be used to access other files
- Revoked automatically on expiration
- HTTPS only (signatures invalid over HTTP)

**Best Practices:**
- Generate URLs just before use
- Don't store URLs long-term
- Don't share URLs externally
- Monitor for expired URL errors
- Implement retry logic with fresh URLs

### Large File Handling

**Upload Optimization:**
- Use streaming uploads for large files
- Monitor upload progress with Content-Length header
- Implement retry logic for network failures
- Consider multipart upload for > 50 MB files

**Download Optimization:**
- Use Range requests for partial downloads
- Implement resume capability for interrupted downloads
- Stream downloads rather than loading into memory
- Consider parallel chunk downloads for very large files

## File History

File records maintain complete history of status changes and statistic updates via the `history` connection field queried through [`Query.ach.file()`]($gql:query:ach.file). Each version captures:
- Processing status at that point in time
- Processing detail messages
- Statistics snapshot
- Modification timestamp
- Version number

History enables:
- Complete audit trail of file processing
- Point-in-time analysis of processing state
- Debugging processing issues
- Compliance and regulatory reporting

## File Queries

Query files by status using the `PROCESSING_STATUS` index, which requires both `processingStatus` and `configId` filter values.

**Supported Filters:**
- `configId`: Filter to specific ACH configuration
- `processingStatus`: Filter by current processing state
- `created`: Filter by file creation timestamp range

Results support pagination via standard connection cursors.

## File Operations

Use GraphQL for all file operations:

**Upload and Processing:**
- [`Mutation.files.createUpload()`]($gql:mutation:files.createUpload): Create presigned upload URL
- [`Mutation.ach.processFile()`]($gql:mutation:ach.processFile): Start file processing
- [`Query.ach.file()`]($gql:query:ach.file): Get file information by key or ID
- [`Query.ach.files()`]($gql:query:ach.files): Query files with status filters

**Generation and Download:**
- [`Mutation.ach.generateFile()`]($gql:mutation:ach.generateFile): Generate NACHA file
- [`Mutation.files.createDownload()`]($gql:mutation:files.createDownload): Create presigned download URL
- [`Query.files.list()`]($gql:query:files.list): List files by key prefix

## Further Reading

To learn file operations from scratch, see the tutorial on [Setting Up ACH Processing](/tutorials/ach/setting-up-ach).

For practical file processing workflows, see the how-to guide on [Reconciling ACH Files](/guides/reconciling-ach-files).

For complete GraphQL type definitions, see:
- []($gql:object:AchFileInfo)
- []($gql:object:AchProcessingStatistics)
- []($gql:enum:AchFileType)
- []($gql:enum:AchFileProcessingStatus)
- []($gql:object:Upload)
- []($gql:object:Download)
