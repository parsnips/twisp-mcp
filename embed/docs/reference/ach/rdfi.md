---
title: RDFI
metaTitle: "Reference: ACH RDFI"
description: Receive and process incoming ACH transactions as an RDFI
showNext: false
showPrev: true
---

## Overview

An RDFI (Receiving Depository Financial Institution) is the financial institution that receives ACH transactions on behalf of a receiver. When operating as an RDFI, the Twisp ACH processor handles incoming ACH files, validates transactions, updates account balances, and generates return files when necessary.

The ACH RDFI processor enables you to:
- Receive and process ACH credit and debit transactions
- Apply transactions to customer accounts with proper ledger accounting
- Handle exceptions through suspense and exception accounts
- Generate return files for transactions that cannot be processed
- Maintain complete audit trails through workflow execution tracking

## Getting Started

### Prerequisites

Before processing RDFI files, you need:

1. **ACH Configuration** - Created via `Mutation.ach.createConfiguration()`
2. **Required Accounts** - Settlement, suspense, exception, and fee accounts
3. **Webhook Endpoint** - For receiving transaction decisioning requests
4. **Journal** - For posting transactions

### Quick Start Example

```graphql
# 1. Create webhook endpoint for ACH decisioning
mutation CreateEndpoint {
  events {
    createEndpoint(
      input: {
        endpointId: "b84512f1-a67e-4dc2-94dd-66c48b4d13fb"
        status: ENABLED
        endpointType: ACH_PROCESSOR
        url: "https://your-domain.com/webhooks/ach"
        subscription: []
        description: "ACH RDFI webhook processor"
      }
    ) {
      endpointId
    }
  }
}

# 2. Create required accounts
mutation CreateAccounts {
  # Settlement account - where funds transit
  settlement: createAccount(
    input: {
      accountId: "37f7e8a6-171f-411d-ad59-7b1f40f505ea"
      code: "settlement.ach"
      name: "ACH Settlement"
      normalBalanceType: DEBIT
      config: {
        enableConcurrentPosting: true
      }
    }
  ) {
    accountId
  }

  # Suspense account - for transactions to unknown accounts
  suspense: createAccount(
    input: {
      accountId: "3171b0c2-6e9f-41aa-a5a6-ee927deb27cf"
      code: "suspense.ach"
      name: "ACH Suspense"
      config: {
        enableConcurrentPosting: true
      }
    }
  ) {
    accountId
  }

  # Exception account - for failed transactions
  exception: createAccount(
    input: {
      accountId: "4a8f2b1e-3c9d-4f7e-a5b6-1d8e9f0a2b3c"
      code: "exception.ach"
      name: "ACH Exception"
      config: {
        enableConcurrentPosting: true
      }
    }
  ) {
    accountId
  }

  # Fee account - for ACH processing fees
  fee: createAccount(
    input: {
      accountId: "5b9e3c2f-4d0e-5a8f-b6c7-2e9f0a1b3c4d"
      code: "fee.ach"
      name: "ACH Fee Income"
      normalBalanceType: CREDIT
      config: {
        enableConcurrentPosting: true
      }
    }
  ) {
    accountId
  }
}

# 3. Create a journal for ACH transactions
mutation CreateJournal {
  createJournal(
    input: {
      journalId: "8d7e6f5a-4b3c-2d1e-0f9a-8b7c6d5e4f3a"
      name: "ACH Processing Journal"
      status: ACTIVE
    }
  ) {
    journalId
  }
}

# 4. Create ACH configuration
mutation CreateACHConfig {
  ach {
    createConfiguration(
      input: {
        configId: "b96d358e-50b8-4ae5-8b07-2e8f33f396c6"
        endpointId: "b84512f1-a67e-4dc2-94dd-66c48b4d13fb"
        journalId: "8d7e6f5a-4b3c-2d1e-0f9a-8b7c6d5e4f3a"
        settlementAccountId: "37f7e8a6-171f-411d-ad59-7b1f40f505ea"
        suspenseAccountId: "3171b0c2-6e9f-41aa-a5a6-ee927deb27cf"
        exceptionAccountId: "4a8f2b1e-3c9d-4f7e-a5b6-1d8e9f0a2b3c"
        feeAccountId: "5b9e3c2f-4d0e-5a8f-b6c7-2e9f0a1b3c4d"
        odfiHeaderConfiguration: {
          immediateDestination: "021000021"
          immediateDestinationName: "Your Bank Name"
          immediateOrigin: "1234567890"
          immediateOriginName: "Your Company Name"
        }
        timeZone: "America/New_York"
      }
    ) {
      configId
      version
    }
  }
}
```

## RDFI Workflow

### 1. Upload ACH File

When you receive an ACH file from the Fed or your upstream processor, upload it to Twisp:

```graphql
mutation CreateUpload {
  files {
    createUpload(
      input: {
        key: "incoming-ach-20251114.ach"
        uploadType: ACH
        contentType: "text/plain"
      }
    ) {
      uploadURL
    }
  }
}
```

Upload the file using the returned URL:

```bash
curl -T incoming-ach-20251114.ach -XPUT '<uploadURL>'
```

### 2. Process ACH File

Start processing the uploaded file:

```graphql
mutation ProcessFile {
  ach {
    processFile(
      input: {
        configId: "b96d358e-50b8-4ae5-8b07-2e8f33f396c6"
        fileKey: "incoming-ach-20251114.ach"
        fileType: RDFI
      }
    ) {
      fileId
    }
  }
}
```

### 3. Monitor File Processing

Check the status of file processing:

```graphql
query GetFileStatus {
  ach {
    file(
      fileKey: "incoming-ach-20251114.ach"
      configId: "b96d358e-50b8-4ae5-8b07-2e8f33f396c6"
    ) {
      fileId
      processingStatus
      processingDetail
      processingStatistics {
        numEntriesUnprocessed
        totalCreditAmount
        totalDebitAmount
      }
    }
  }
}
```

**Processing Status Values:**
- `NEW` - File created, not yet processing
- `UPLOADED` - File uploaded and queued
- `VALIDATING` - File format validation in progress
- `PARTITIONING` - Preparing for parallel processing
- `PROCESSING` - Sending webhooks for transaction decisions
- `PROCESSED` - All webhooks sent, awaiting settlements
- `COMPLETED` - All transactions settled or returned
- `ERROR` - Unrecoverable error occurred
- `INVALID` - File failed validation

### 4. Handle Transaction Webhooks

For each ACH entry in the file, Twisp sends a webhook to your endpoint. You must respond with instructions on how to handle the transaction.

**Webhook Payload Format:**

```json
{
  "workflowName": "ACH.RDFI.CR",
  "workflowTask": "CREATE",
  "executionId": "60f7ac42-ff72-48c7-af58-ee1f9a2db1e0",
  "fileHeader": {
    "id": "file-header-id",
    "immediateDestination": "021000021",
    "immediateOrigin": "1234567890",
    "fileCreationDate": "251114",
    "fileCreationTime": "1030",
    "fileIDModifier": "A",
    "immediateDestinationName": "Your Bank Name",
    "immediateOriginName": "Originating Company",
    "referenceCode": ""
  },
  "batchHeader": {
    "id": "batch-id",
    "serviceClassCode": "220",
    "companyName": "PAYROLL CO",
    "companyDiscretionaryData": "",
    "companyIdentification": "1234567890",
    "standardEntryClassCode": "PPD",
    "companyEntryDescription": "PAYROLL",
    "companyDescriptiveDate": "",
    "effectiveEntryDate": "251115",
    "settlementDate": "   ",
    "originatorStatusCode": "1",
    "odfiIdentification": "12345678",
    "batchNumber": "0000001"
  },
  "entryDetail": {
    "id": "entry-id",
    "transactionCode": "22",
    "rdfiIdentification": "02100002",
    "checkDigit": "1",
    "dfiAccountNumber": "123456789",
    "amount": "150000",
    "identificationNumber": "employee-123",
    "individualName": "John Doe",
    "discretionaryData": "",
    "addendaRecordIndicator": "0",
    "traceNumber": "123456780000001",
    "category": "Forward"
  }
}
```

**Response Format:**

You must respond with one of three actions: `SETTLE`, `RETURN`, or `RETRY`.

#### Option 1: Settle (Accept Transaction)

```json
{
  "action": "SETTLE",
  "accountId": "d2f7183f-8e9c-45e7-9a98-ef1897ddb930",
  "when": "2025-11-15T00:00:00.000Z",
  "metadata": {
    "customerId": "cust-123",
    "transactionType": "payroll"
  },
  "entryMetadata": {
    "customerId": "cust-123"
  }
}
```

- `when` is optional. If omitted, uses the effective date from the batch header
- If `when` is in the past, the transaction settles immediately
- `metadata` is optional and attached to the ledger transaction
- `entryMetadata` is optional and attached to the ledger entries

#### Option 2: Return (Reject Transaction)

```json
{
  "action": "RETURN",
  "accountId": "d2f7183f-8e9c-45e7-9a98-ef1897ddb930",
  "addenda99": {
    "returnCode": "R01",
    "addendaInformation": "Insufficient Funds"
  },
  "metadata": {
    "reason": "account_balance_insufficient"
  }
}
```

**Common Return Codes:**
- `R01` - Insufficient Funds
- `R02` - Account Closed
- `R03` - No Account / Unable to Locate Account
- `R04` - Invalid Account Number
- `R05` - Unauthorized Debit to Consumer Account
- `R07` - Authorization Revoked by Customer
- `R08` - Payment Stopped
- `R10` - Customer Advises Not Authorized

[See complete return code reference](#return-codes)

#### Option 3: Retry (Temporary Error)

```json
{
  "action": "RETRY"
}
```

Use `RETRY` when:
- Your system is temporarily unavailable
- You need more time to make a decision
- There's a transient error in your processing

Twisp will exponentially back off and retry the webhook.

### 5. Generate Return File

After processing is complete, generate a return file for any transactions you rejected:

```graphql
mutation GenerateReturnFile {
  ach {
    generateFile(
      input: {
        configId: "b96d358e-50b8-4ae5-8b07-2e8f33f396c6"
        fileKey: "return-20251114.ach"
        fileType: RDFI_RETURN
        generateEmpty: false
      }
    ) {
      fileKey
      generated
    }
  }
}
```

`generateEmpty: false` means the file is only created if there are returns to include.

### 6. Download Return File

Download the generated return file:

```graphql
mutation DownloadReturn {
  files {
    createDownload(
      key: "return-20251114.ach"
    ) {
      downloadURL
    }
  }
}
```

Then download using the URL:

```bash
curl '<downloadURL>' -o return-20251114.ach
```

Transmit this file to the originating ODFI via your normal file transmission process (SFTP, etc.).

## Transaction Lifecycle and Ledger Accounting

The RDFI processor uses a three-stage workflow with double-entry accounting at each stage.

### Stage 1: CREATE (Initial Encumbrance)

When a transaction webhook is received, an **encumbrance** is created on the target account:

**For Credits (Incoming Deposits):**
```
DR Settlement Account (Encumbrance Layer)
CR Customer Account (Encumbrance Layer)
```

**For Debits (Outgoing Withdrawals):**
```
DR Customer Account (Encumbrance Layer)
CR Settlement Account (Encumbrance Layer)
```

The encumbrance layer reserves funds but doesn't affect available balance. This allows you to:
- Track expected funds before they settle
- Maintain visibility of in-flight transactions
- Reconcile with external ACH reports

### Stage 2: SETTLE (Final Settlement)

When you respond with `"action": "SETTLE"`, two things happen:

1. **Reverse the encumbrance:**
```
Opposite of CREATE entries with negative amounts
```

2. **Post to settled layer:**
```
DR/CR Customer Account (Settled Layer)
DR/CR Settlement Account (Settled Layer)
```

The settled layer represents final, available balances that customers can access.

### Stage 3: RETURN (Rejection)

When you respond with `"action": "RETURN"`, the transaction is reversed:

1. **Reverse the encumbrance** (same as SETTLE step 1)
2. **Post return to settled layer** (opposite direction of a normal settlement)
3. **Queue for return file generation**

Returns are included in the next return file you generate via `Mutation.ach.generateFile()`.

### Balance Layer Illustration

```
┌─────────────────────────────────────────────────┐
│  ENCUMBRANCE Layer                              │
│  • In-flight ACH transactions                   │
│  • Not available to customer                    │
│  • Tracks expected debits/credits               │
└─────────────────────────────────────────────────┘
                    ↓ SETTLE
┌─────────────────────────────────────────────────┐
│  SETTLED Layer                                  │
│  • Final, available balance                     │
│  • Customer can withdraw/spend                  │
│  • Appears in balance queries                   │
└─────────────────────────────────────────────────┘
```

### Querying Balances

Check account balances across all layers:

```graphql
query GetAccountBalance {
  balance(
    accountId: "d2f7183f-8e9c-45e7-9a98-ef1897ddb930"
    journalId: "8d7e6f5a-4b3c-2d1e-0f9a-8b7c6d5e4f3a"
    currency: "USD"
  ) {
    settled {
      crBalance {
        units
        currency
      }
      drBalance {
        units
        currency
      }
    }
    pending {
      crBalance {
        units
        currency
      }
      drBalance {
        units
        currency
      }
    }
    encumbrance {
      crBalance {
        units
        currency
      }
      drBalance {
        units
        currency
      }
    }
    version
  }
}
```

**Available Balance Calculation:**
```
Available = Settled - Pending - Encumbrance (for debits)
```

## Return Generation

Create returns for transactions that cannot be processed:

### Return Decision Logic

Automatic and manual return triggers:

- **Insufficient Funds**: Account balance insufficient for debit
- **Account Closed**: Target account no longer active
- **Invalid Account**: Account number not found
- **Unauthorized**: Transaction not authorized by account holder
- **Stop Payment**: Account holder placed stop payment order

### Return Codes

Select appropriate return code:

- **R01**: Insufficient Funds
- **R02**: Account Closed
- **R03**: No Account / Unable to Locate Account
- **R04**: Invalid Account Number
- **R05**: Unauthorized Debit to Consumer Account (improper authorization)
- **R07**: Authorization Revoked by Customer
- **R08**: Payment Stopped
- **R10**: Customer Advises Not Authorized
- **R29**: Corporate Customer Advises Not Authorized

### Return Timing

Return deadlines by code:

- **2 Business Days**: Most return codes (R01-R04, R07-R08, etc.)
- **60 Calendar Days**: Unauthorized returns (R05, R07, R10, R29)
- **Next Business Day**: Same-day ACH returns

### Return File Generation

Create NACHA return files:

- **Return Entry**: Create return detail record in Twisp
- **Return Batch**: Group returns in batches
- **Return File Generation**: Generate complete NACHA return file via Files API
- **File Download**: Retrieve generated return file from Twisp
- **File Transmission**: You transmit return file to originating ODFI via SFTP/FTPS

## Monitoring and Observability

### Query Files by Status

Find all files in a specific processing state:

```graphql
query GetProcessingFiles {
  ach {
    files(
      first: 100
      index: { name: PROCESSING_STATUS }
      where: {
        configId: { eq: "b96d358e-50b8-4ae5-8b07-2e8f33f396c6" }
        processingStatus: { eq: "PROCESSING" }
        created: { gte: "2025-11-01T00:00:00Z" }
      }
    ) {
      nodes {
        fileId
        fileKey
        processingStatus
        processingDetail
        processingStatistics {
          numEntriesUnprocessed
          totalCreditAmount
          totalDebitAmount
        }
        created
        modified
      }
      pageInfo {
        hasNextPage
        endCursor
      }
    }
  }
}
```

### Track Workflow Execution

Every ACH transaction creates a workflow execution that you can query:

```graphql
query GetWorkflowExecution {
  workflow {
    execution(
      executionId: "60f7ac42-ff72-48c7-af58-ee1f9a2db1e0"
    ) {
      workflowId
      executionId
      task
      params
      output {
        state
      }
      activities {
        action
        entityType
        entityId
        entity {
          ... on Transaction {
            transactionId
            effective
            description
          }
          ... on AchWorkflowTrace {
            traceNumber
            fileId
            configId
          }
        }
      }
      created
      modified
      version
    }
  }
}
```

The `activities` field shows all transactions and ACH traces created by this workflow, giving you complete auditability.

## Exception Handling

### Suspense Account

When a transaction targets an account that doesn't exist, it's automatically posted to your configured `suspenseAccountId`. This allows you to:

1. Accept the transaction (avoiding a return)
2. Research the correct account
3. Create a manual journal entry to move funds to the correct account

**Example scenario:**
- ACH credit arrives for account number "123456789"
- Account doesn't exist in your system
- Transaction is posted:
  ```
  DR Settlement Account
  CR Suspense Account
  ```
- You investigate and find the correct account is "123456790"
- You create a journal entry to move the funds:
  ```
  DR Suspense Account
  CR Correct Customer Account
  ```

### Exception Account

When a transaction fails due to velocity controls, account state issues, or other processing errors, it's posted to your `exceptionAccountId`. Common scenarios:

- **Velocity Control Violation**: Transaction exceeds configured velocity limits
- **Account Locked**: Target account is frozen or locked
- **Processing Error**: Temporary system issue

Funds in the exception account should typically be returned to the originator via a return file.

### Velocity Control Integration

The RDFI processor respects velocity controls configured on customer accounts. If a debit would exceed the velocity limit:

1. The webhook is NOT sent to your endpoint
2. The transaction is automatically posted to the exception account
3. A return will be generated with code R01 (Insufficient Funds)

This provides built-in protection against overdrafts and unauthorized transactions.

## Practical Examples

### Example 1: Basic RDFI Setup and Processing

Complete flow from setup to settlement:

```graphql
# Step 1: Setup (run once)
mutation Setup {
  # Create webhook endpoint
  endpoint: events {
    createEndpoint(
      input: {
        endpointId: "webhook-001"
        status: ENABLED
        endpointType: ACH_PROCESSOR
        url: "https://api.yourcompany.com/ach/webhook"
        subscription: []
      }
    ) { endpointId }
  }

  # Create journal
  journal: createJournal(
    input: {
      journalId: "journal-001"
      name: "ACH Journal"
      status: ACTIVE
    }
  ) { journalId }

  # Create accounts (abbreviated)
  settlement: createAccount(
    input: {
      accountId: "acct-settlement"
      code: "settlement"
      name: "ACH Settlement"
      normalBalanceType: DEBIT
      config: { enableConcurrentPosting: true }
    }
  ) { accountId }

  # Create ACH config
  config: ach {
    createConfiguration(
      input: {
        configId: "config-001"
        endpointId: "webhook-001"
        journalId: "journal-001"
        settlementAccountId: "acct-settlement"
        suspenseAccountId: "acct-suspense"
        exceptionAccountId: "acct-exception"
        feeAccountId: "acct-fee"
        odfiHeaderConfiguration: {
          immediateDestination: "021000021"
          immediateDestinationName: "Federal Reserve Bank"
          immediateOrigin: "1234567890"
          immediateOriginName: "Your Company"
        }
        timeZone: "America/New_York"
      }
    ) { configId }
  }
}

# Step 2: Upload file (when received from Fed)
mutation UploadFile {
  files {
    createUpload(
      input: {
        key: "incoming-20251114-001.ach"
        uploadType: ACH
        contentType: "text/plain"
      }
    ) { uploadURL }
  }
}
# Use uploadURL to PUT file content

# Step 3: Process file
mutation ProcessFile {
  ach {
    processFile(
      input: {
        configId: "config-001"
        fileKey: "incoming-20251114-001.ach"
        fileType: RDFI
      }
    ) { fileId }
  }
}

# Step 4: Monitor processing
query MonitorFile {
  ach {
    file(
      fileKey: "incoming-20251114-001.ach"
      configId: "config-001"
    ) {
      processingStatus
      processingStatistics {
        numEntriesUnprocessed
        totalCreditAmount
        totalDebitAmount
      }
    }
  }
}

# Step 5: Generate returns (after webhooks complete)
mutation GenerateReturns {
  ach {
    generateFile(
      input: {
        configId: "config-001"
        fileKey: "return-20251114-001.ach"
        fileType: RDFI_RETURN
        generateEmpty: false
      }
    ) {
      fileKey
      generated
    }
  }
}

# Step 6: Download returns
mutation DownloadReturns {
  files {
    createDownload(
      key: "return-20251114-001.ach"
    ) { downloadURL }
  }
}
```

### Example 2: Webhook Handler Implementation

Sample webhook handler in Node.js:

```javascript
app.post('/ach/webhook', async (req, res) => {
  const { workflowName, workflowTask, executionId, entryDetail } = req.body;

  try {
    // Extract transaction details
    const accountNumber = entryDetail.dfiAccountNumber;
    const amount = parseFloat(entryDetail.amount) / 100; // Amount is in cents
    const isDebit = entryDetail.transactionCode.startsWith('2'); // 22, 23, 24
    const isCredit = entryDetail.transactionCode.startsWith('3'); // 32, 33, 34

    // Look up customer account
    const account = await findAccountByNumber(accountNumber);

    if (!account) {
      // Account not found - will go to suspense
      return res.json({
        action: 'SETTLE',
        accountId: SUSPENSE_ACCOUNT_ID,
        metadata: {
          reason: 'account_not_found',
          originalAccountNumber: accountNumber
        }
      });
    }

    if (account.status === 'CLOSED') {
      // Account closed - return with R02
      return res.json({
        action: 'RETURN',
        accountId: account.id,
        addenda99: {
          returnCode: 'R02',
          addendaInformation: 'Account Closed'
        }
      });
    }

    if (isDebit) {
      // Check balance for debits
      const balance = await getAccountBalance(account.id);
      if (balance < amount) {
        return res.json({
          action: 'RETURN',
          accountId: account.id,
          addenda99: {
            returnCode: 'R01',
            addendaInformation: 'Insufficient Funds'
          }
        });
      }
    }

    // All checks passed - settle the transaction
    return res.json({
      action: 'SETTLE',
      accountId: account.id,
      when: new Date().toISOString(), // Settle immediately
      metadata: {
        customerId: account.customerId,
        originalTraceNumber: entryDetail.traceNumber
      }
    });

  } catch (error) {
    console.error('Webhook processing error:', error);
    // Retry on errors
    return res.json({
      action: 'RETRY'
    });
  }
});
```

## API Reference

### GraphQL Operations

**Configuration:**
- `Query.ach.configuration(id: UUID!)` - Get ACH configuration
- `Query.ach.configurations(first: Int!)` - List all configurations
- `Mutation.ach.createConfiguration(input: AchCreateConfigurationInput!)` - Create configuration
- `Mutation.ach.updateConfiguration(configId: UUID!, input: AchUpdateConfigurationInput!)` - Update configuration

**File Operations:**
- `Query.ach.file(id: UUID, fileKey: String, configId: UUID)` - Get file status
- `Query.ach.files(index: AchFileInfoIndexInput!, where: AchFileInfoFilterInput!, first: Int!)` - Query files
- `Mutation.ach.processFile(input: AchProcessFileInput!)` - Process uploaded file
- `Mutation.ach.generateFile(input: AchGenerateFileInput!)` - Generate return/NOC file
- `Mutation.files.createUpload(input: CreateUploadInput!)` - Get upload URL
- `Mutation.files.createDownload(key: String!)` - Get download URL

**Workflow Operations:**
- `Query.workflow.execution(executionId: UUID!)` - Get workflow execution details

## Return Codes Reference

When returning a transaction, use the appropriate return code in the `addenda99.returnCode` field.

### Standard Return Codes

| Code | Reason | Description | Timing |
|------|--------|-------------|--------|
| `R01` | Insufficient Funds | Available balance is not sufficient to cover the dollar value of the debit entry | 2 business days |
| `R02` | Account Closed | Previously active account has been closed by customer or RDFI | 2 business days |
| `R03` | No Account/Unable to Locate Account | Account number structure is valid and passes editing process, but does not correspond to individual or is not an open account | 2 business days |
| `R04` | Invalid Account Number | Account number structure not valid; entry may fail check digit validation or may contain an incorrect number of digits | 2 business days |
| `R05` | Improper Debit to Consumer Account | A CCD, CTX, or CBR debit entry was transmitted to a Consumer Account of the Receiver and was not authorized by the Receiver | 60 days |
| `R06` | Returned per ODFI's Request | ODFI has requested RDFI to return the ACH entry (optional to RDFI - ODFI indemnifies RDFI) | 2 business days |
| `R07` | Authorization Revoked by Customer | Consumer, who previously authorized ACH payment, has revoked authorization from Originator | 60 days |
| `R08` | Payment Stopped | Receiver of a recurring debit transaction has stopped payment to a specific ACH debit | 2 business days |
| `R09` | Uncollected Funds | Sufficient book or ledger balance exists to satisfy dollar value of the transaction, but the dollar value of transaction is in process of collection | 2 business days |
| `R10` | Customer Advises Originator is Not Known to Receiver and/or Originator is Not Authorized | The receiver does not know the Originator's identity and/or has not authorized the Originator to debit | 60 days |
| `R11` | Customer Advises Entry Not in Accordance with the Terms of the Authorization | The Originator and Receiver have a relationship, and an authorization to debit exists, but there is an error or defect in the payment | 60 days |
| `R12` | Branch Sold to Another DFI | Financial institution receives entry destined for an account at a branch that has been sold to another financial institution | 2 business days |
| `R13` | RDFI not qualified to participate | Financial institution does not receive commercial ACH entries | 2 business days |
| `R14` | Representative payee deceased or unable to continue in that capacity | The representative payee authorized to accept entries on behalf of a beneficiary is either deceased or unable to continue in that capacity | 2 business days |
| `R15` | Beneficiary or bank account holder deceased | (1) the beneficiary entitled to payments is deceased or (2) the bank account holder other than a representative payee is deceased | 2 business days |
| `R16` | Bank account frozen | Funds in bank account are unavailable due to action by RDFI or legal order | 2 business days |
| `R17` | File record edit criteria | Fields rejected by RDFI processing (identified in return addenda) | 2 business days |
| `R20` | Non-payment bank account | Entry destined for non-payment bank account defined by regulation | 2 business days |
| `R23` | Credit entry refused by receiver | Receiver returned entry because minimum or exact amount not remitted, bank account is subject to litigation, or payment represents an overpayment | 2 business days |
| `R29` | Corporate customer advises not authorized | RDFI has been notified by corporate receiver that debit entry of originator is not authorized | 2 business days |

### Best Practices

**Return Timing:**
- Most returns must be sent within **2 business days** of settlement date
- Unauthorized returns (R05, R07, R10, R11) can be returned up to **60 days** after settlement
- Same-day ACH returns must be sent by the same business day

**Common Scenarios:**

```javascript
// Insufficient funds
{
  "action": "RETURN",
  "accountId": "account-id",
  "addenda99": {
    "returnCode": "R01",
    "addendaInformation": "Insufficient Funds"
  }
}

// Account closed
{
  "action": "RETURN",
  "accountId": "account-id",
  "addenda99": {
    "returnCode": "R02",
    "addendaInformation": "Account Closed"
  }
}

// Account not found
{
  "action": "RETURN",
  "accountId": "account-id",
  "addenda99": {
    "returnCode": "R03",
    "addendaInformation": "No Account"
  }
}

// Unauthorized transaction
{
  "action": "RETURN",
  "accountId": "account-id",
  "addenda99": {
    "returnCode": "R10",
    "addendaInformation": "Not Authorized"
  }
}
```

## Best Practices

### Security
- **Webhook Authentication**: Validate webhook signatures to ensure requests are from Twisp
- **HTTPS Only**: Always use HTTPS endpoints for webhooks
- **Idempotency**: Handle duplicate webhooks gracefully using `executionId`

### Performance
- **Fast Webhook Response**: Respond to webhooks within 30 seconds
- **Async Processing**: Queue webhook processing if complex logic is needed
- **Retry Logic**: Implement exponential backoff for webhook retries

### Monitoring
- **Alert on Status Changes**: Monitor file processing status for errors
- **Track Return Rates**: High return rates may indicate data quality issues
- **Balance Reconciliation**: Daily reconciliation of settlement account

### Compliance
- **Return Timeframes**: Adhere to NACHA return deadlines (2 days for most codes)
- **Authorization Records**: Maintain proof of authorization for debits
- **Transaction History**: Keep complete audit trail for 7 years
- **Reg E Compliance**: Honor consumer dispute rights (60-day investigation period)

## Further Reading

For practical guidance on receiving ACH transactions:

- [Handling ACH Returns and NOCs](/guides/handling-ach-returns) - Return processing and NOC management
- [Reconciling ACH Files](/guides/reconciling-ach-files) - File validation and reconciliation procedures
- [Processing ACH Payments](/guides/processing-ach-payments) - Risk management and best practices

For additional technical details:

- [Configuration](/reference/ach/configuration) - RDFI configuration parameters
- [File Operations](/reference/ach/file-operations) - File upload and parsing API specifications
- [ODFI Reference](/reference/ach/odfi) - Origination perspective
- [ACH Processor](/processors/ach) - Conceptual overview
