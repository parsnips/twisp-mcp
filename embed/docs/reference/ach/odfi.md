---
title: ODFI
metaTitle: "Reference: ACH ODFI"
description: Reference for ODFI operations within the Twisp ACH Processor
showNext: true
showPrev: true
---

## Overview

An ODFI (Originating Depository Financial Institution) originates ACH transactions on behalf of an originator. When operating as an ODFI, the Twisp ACH processor handles transaction lifecycle management, ledger accounting across multiple balance layers, NACHA file generation, and return processing for outgoing ACH payments.

The ACH ODFI processor enables you to:
- Originate ACH credit (PUSH) and debit (PULL) transactions via workflows
- Manage funds through encumbrance, pending, and settled balance layers
- Generate NACHA-compliant files for transmission to financial institutions
- Track transaction lifecycle from creation through settlement or return
- Process returns from RDFIs with automatic ledger reversals
- Maintain complete audit trails via workflow execution history

## Getting Started

### Prerequisites

Before originating ACH transactions, you need:

1. **ACH Configuration** - Created via [`Mutation.ach.createConfiguration()`]($gql:mutation:ach.createConfiguration)
2. **Required Accounts** - Settlement, suspense, exception, and fee accounts
3. **Journal** - For posting all ACH transactions
4. **Customer Accounts** - Accounts to debit (PUSH) or credit (PULL)

### Quick Start Example

```graphql
# 1. Create required accounts
mutation CreateAccounts {
  # Settlement account - where funds transit during ACH processing
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

  # Customer account - for testing
  customer: createAccount(
    input: {
      accountId: "d2f7183f-8e9c-45e7-9a98-ef1897ddb930"
      code: "customer.001"
      name: "Customer Account"
      normalBalanceType: DEBIT
      config: {
        enableConcurrentPosting: true
      }
    }
  ) {
    accountId
  }
}

# 2. Create a journal for ACH transactions
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

# 3. Create a webhook endpoint (Note: webhook endpoint not used in ODFI-only use cases)
mutation CreateACHWebhookProcessor {
  events {
    createEndpoint(
      input: {
        endpointId: "b84512f1-a67e-4dc2-94dd-66c48b4d13fb"
        status: ENABLED
        endpointType: ACH_PROCESSOR
        url: "https://webhook.site/ach-testing"
        subscription: []
        description: "ACH webhook processor"
      }
    ) {
      endpointId
    }
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
          immediateDestinationName: "Federal Reserve Bank"
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

## ODFI Workflows

### Workflow Types

There are 2 primary workflow types for originating ACH transactions:

1. **ACH PUSH (Credits)**: Send funds to receivers. Workflow ID `934498b5-b4f1-46c4-ad79-868939dc39e8`. Used for payroll, vendor payments, refunds.
2. **ACH PULL (Debits)**: Collect funds from receivers. Workflow ID `064e3b76-6072-451e-aace-5a7be3704ee2`. Used for bill payments, subscriptions, loan payments.

Both workflows accept the same parameters but follow different balance layer progression based on settlement timing requirements.

### PUSH Workflow (ACH Credits)

The PUSH workflow originates credit transactions that send funds to receivers. When you execute a PUSH workflow via [`Mutation.workflow.execute()`]($gql:mutation:workflow.execute), funds are immediately encumbered on the customer's account and settled upon submission.

#### CREATE State

The CREATE state reserves funds and charges processing fees:

```graphql
mutation CreatePushTransaction {
  workflow {
    execute(
      input: {
        executionId:"daf20572-c1b1-11f0-8b14-069b540ea27c"
        code: "ACH_PUSH"
        task: "CREATE"
        params: {
          configId: "b96d358e-50b8-4ae5-8b07-2e8f33f396c6"
          accountId: "d2f7183f-8e9c-45e7-9a98-ef1897ddb930"
          settlementAccountId: "37f7e8a6-171f-411d-ad59-7b1f40f505ea"
          journalId: "8d7e6f5a-4b3c-2d1e-0f9a-8b7c6d5e4f3a"
          amount: "1500.00"
          routingNumber:       "026009593",
          accountNumber:       "12345678901234567",
          accountType:         "checking",
          individualName:      "Clark Kent",
          entryDescription:    "XFER",
          effective:"2025-11-14"
          feeAccountId: "5b9e3c2f-4d0e-5a8f-b6c7-2e9f0a1b3c4d"
          feeAmount: "0.25"
          correlationId: "payroll-batch-001"
          metadata: {
            employeeId: "emp-123"
            payrollPeriod: "2025-11-01"
          }
          entryMetadata: {
            companyName: "ACME Corp"
            companyEntryDescription: "PAYROLL"
            individualName: "John Doe"
            standardEntryClassCode: "PPD"
            dfiAccountNumber: "123456789"
            identificationNumber: "emp-123"
            rdfiIdentification: "021000021"
          }
        }
      }
    ) {
      executionId
      output {
        state
      }
    }
  }
}
```

**Ledger Entries Created:**

Encumbrance layer:
```
DR Customer Account (Encumbrance)     $1,500.00
CR Settlement Account (Encumbrance)   $1,500.00
```

Settled layer (fees):
```
DR Customer Account (Settled)         $0.25
CR Fee Account (Settled)              $0.25
```

#### SUBMIT State

The SUBMIT state finalizes the transaction for file generation:

```graphql
mutation SubmitPushTransaction {
  workflow {
    execute(
      input: {
        executionId:"daf20572-c1b1-11f0-8b14-069b540ea27c"
        code: "ACH_PUSH"
        task: "SUBMIT"
        params: {}
      }
    ) {
      executionId
      output {
        state
      }
    }
  }
}
```

**Ledger Entries Created:**

Reverse encumbrance:
```
DR Customer Account (Encumbrance)     -$1,500.00
CR Settlement Account (Encumbrance)   -$1,500.00
```

Post to settled:
```
DR Customer Account (Settled)         $1,500.00
CR Settlement Account (Settled)       $1,500.00
```

The transaction is now queued for inclusion in the next file generated via [`Mutation.ach.generateFile()`]($gql:mutation:ach.generateFile).

#### CANCEL State

Cancel a transaction before file generation:

```graphql
mutation CancelPushTransaction {
  workflow {
    execute(
      input: {
        executionId:"daf20572-c1b1-11f0-8b14-069b540ea27c"
        code: "ACH_PUSH"
        task: "CANCEL"
        params: {}
      }
    ) {
      executionId
      output {
        state
      }
    }
  }
}
```

Reverses the CREATE encumbrance. Follow with REIMBURSE_FEE to refund processing fees.

### PULL Workflow (ACH Debits)

The PULL workflow originates debit transactions that collect funds from receivers. PULL workflows have an additional SETTLE state between SUBMIT and final settlement.

#### CREATE State

The CREATE state reserves funds on the settlement account:

```graphql
mutation CreatePullTransaction {
  workflow {
    execute(
      input: {
        executionId: "c97685f6-c1b2-11f0-a959-069b540ea27c"
        code: "ACH_PULL"
        task: "CREATE"
        params: {
          configId: "b96d358e-50b8-4ae5-8b07-2e8f33f396c6"
          accountId: "d2f7183f-8e9c-45e7-9a98-ef1897ddb930"
          settlementAccountId: "37f7e8a6-171f-411d-ad59-7b1f40f505ea"
          journalId: "8d7e6f5a-4b3c-2d1e-0f9a-8b7c6d5e4f3a"
          routingNumber:       "026009593",
          accountNumber:       "12345678901234567",
          accountType:         "checking",
          individualName:      "Clark Kent",
          entryDescription:    "XFER",
          effective:"2025-11-14"
          amount: "250.00"
          feeAccountId: "5b9e3c2f-4d0e-5a8f-b6c7-2e9f0a1b3c4d"
          feeAmount: "0.25"
          correlationId: "billing-batch-001"
          metadata: {
            customerId: "cust-456"
            invoiceNumber: "INV-2025-001"
          }
          entryMetadata: {
            companyName: "ACME Corp"
            companyEntryDescription: "INVOICE"
            individualName: "Jane Smith"
            standardEntryClassCode: "WEB"
            dfiAccountNumber: "987654321"
            identificationNumber: "inv-001"
            rdfiIdentification: "021000021"
          }
        }
      }
    ) {
      executionId
      output {
        state
      }
    }
  }
}
```

**Ledger Entries Created:**

Encumbrance layer (inverse of PUSH):
```
DR Settlement Account (Encumbrance)   $250.00
CR Customer Account (Encumbrance)     $250.00
```

Settled layer (fees):
```
DR Customer Account (Settled)         $0.25
CR Fee Account (Settled)              $0.25
```

#### SUBMIT State

The SUBMIT state moves funds to pending layer for file generation:

```graphql
mutation SubmitPullTransaction {
  workflow {
    execute(
      input: {
        executionId: "c97685f6-c1b2-11f0-a959-069b540ea27c"
        code: "ACH_PULL"
        task: "SUBMIT"
        params: {}
      }
    ) {
      executionId
      output {
        state
      }
    }
  }
}
```

**Ledger Entries Created:**

Reverse encumbrance:
```
DR Settlement Account (Encumbrance)   -$250.00
CR Customer Account (Encumbrance)     -$250.00
```

Post to pending:
```
DR Settlement Account (Pending)       $250.00
CR Customer Account (Pending)         $250.00
```

The transaction remains in pending layer until SETTLE execution.

#### SETTLE State

The SETTLE state finalizes the transaction after collection confirmation:

```graphql
mutation SettlePullTransaction {
  workflow {
    execute(
      input: {
        executionId: "c97685f6-c1b2-11f0-a959-069b540ea27c"
        code: "ACH_PULL"
        task: "SETTLE"
        params: {}
      }
    ) {
      executionId
      output {
        state
      }
    }
  }
}
```

**Ledger Entries Created:**

Reverse pending:
```
DR Settlement Account (Pending)       -$250.00
CR Customer Account (Pending)         -$250.00
```

Post to settled:
```
DR Settlement Account (Settled)       $250.00
CR Customer Account (Settled)         $250.00
```

### Monitoring Workflow Execution

Query workflow execution details using [`Query.workflow.execution()`]($gql:query:workflow.execution):

```graphql
query GetWorkflowExecution {
  workflow {
    execution(
      executionId: "c97685f6-c1b2-11f0-a959-069b540ea27c"
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
            entries(first:4) {
              nodes {
                accountId
                layer
                amount {
                  units
                  currency
                }
              }
            }
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

## Workflow Parameters

All ODFI workflows accept common parameters:

- **Account ID**: Customer account to debit (PUSH) or credit (PULL)
- **Settlement Account ID**: Central account for fund transit during processing
- **Journal ID**: Ledger journal for posting all transactions
- **Amount**: Transaction amount in decimal format (e.g., "1500.00")
- **Effective Date**: Settlement date in YYYY-MM-DD format
- **Fee Account ID** (optional): Account for fee income (default: zero UUID)
- **Fee Amount** (optional): Processing fee in decimal format (default: "0")
- **Correlation ID**: Unique identifier for grouping related transactions
- **Metadata** (optional): Transaction-level JSON metadata
- **Entry Metadata** (optional): Entry-level JSON metadata

## PUSH Workflow States

PUSH workflows originate credit transactions with immediate settlement upon submission.

**Workflow States:**

1. **CREATE**: Initial state creates encumbrance on customer account and charges fees. Funds reserved but not yet transmitted.

2. **CANCEL**: Reverses CREATE encumbrance before file generation. Optionally followed by REIMBURSE_FEE to refund processing fees.

3. **SUBMIT**: Final state settles funds, reversing encumbrance and posting to settled layer. Transaction included in next file generation.

4. **RETURN**: Processes return from RDFI, reversing SUBMIT settlement and returning funds to customer account.

5. **REIMBURSE_FEE**: Refunds processing fee, reversing original fee charge.

**Key Characteristic:** PUSH workflows settle immediately on SUBMIT since credit transactions are final upon transmission.

## PULL Workflow States

PULL workflows originate debit transactions with delayed settlement pending confirmation.

**Workflow States:**

1. **CREATE**: Initial state creates encumbrance on settlement account (inverse of PUSH). Funds reserved awaiting collection.

2. **CANCEL**: Reverses CREATE encumbrance before file generation.

3. **SUBMIT**: Moves funds from encumbrance to pending layer. Transaction included in file but not yet final.

4. **SETTLE**: Final settlement moves funds from pending to settled layer after collection confirmation.

5. **RETURN**: Processes return from RDFI, reversing pending or settled amounts depending on when return received.

6. **REIMBURSE_FEE**: Refunds processing fee on cancellation or return.

**Key Characteristic:** PULL workflows separate SUBMIT (transmission) from SETTLE (confirmation) since debit transactions require validation before finalization.

## Transaction Lifecycle and Ledger Accounting

The ODFI processor uses multi-stage workflows with double-entry accounting at each stage, tracking funds through multiple balance layers.

### Balance Layers

ODFI workflows utilize three balance layers to track transaction lifecycle and fund availability:

#### Encumbrance Layer

Temporary holds during CREATE state before transmission. Encumbered amounts don't affect available balance calculations but reserve funds for pending transactions.

**PUSH (Credits):**
```
DR Customer Account (Encumbrance)
CR Settlement Account (Encumbrance)
```

**PULL (Debits):**
```
DR Settlement Account (Encumbrance)
CR Customer Account (Encumbrance)
```

The encumbrance layer allows you to:
- Reserve funds before file generation
- Track expected outflows (PUSH) or inflows (PULL)
- Maintain visibility of in-flight transactions
- Cancel transactions before file transmission

#### Pending Layer

PULL workflows only. Tracks transmitted debits awaiting settlement confirmation between SUBMIT and SETTLE states.

```
DR Settlement Account (Pending)
CR Customer Account (Pending)
```

The pending layer represents:
- Debits transmitted but not yet collected
- Funds awaiting final confirmation from RDFI
- Transactions that can still be returned

#### Settled Layer

Final layer for completed transactions. All fees post directly to settled layer.

**PUSH SUBMIT:**
```
DR Customer Account (Settled)
CR Settlement Account (Settled)
```

**PULL SETTLE:**
```
DR Settlement Account (Settled)
CR Customer Account (Settled)
```

**Fees (all workflows):**
```
DR Customer Account (Settled)
CR Fee Account (Settled)
```

The settled layer represents final, available balances that affect customer account availability.

### Balance Layer Illustration

```
PUSH (Credit) Flow:
┌─────────────────────────────────────────────────┐
│  CREATE: ENCUMBRANCE Layer                      │
│  DR Customer Account                            │
│  CR Settlement Account                          │
│  • Funds reserved, not yet transmitted          │
└─────────────────────────────────────────────────┘
                    ↓ SUBMIT
┌─────────────────────────────────────────────────┐
│  SUBMIT: SETTLED Layer                          │
│  DR Customer Account                            │
│  CR Settlement Account                          │
│  • Final settlement, funds transmitted          │
│  • Included in next file generation             │
└─────────────────────────────────────────────────┘

PULL (Debit) Flow:
┌─────────────────────────────────────────────────┐
│  CREATE: ENCUMBRANCE Layer                      │
│  DR Settlement Account                          │
│  CR Customer Account                            │
│  • Collection reserved, not yet transmitted     │
└─────────────────────────────────────────────────┘
                    ↓ SUBMIT
┌─────────────────────────────────────────────────┐
│  SUBMIT: PENDING Layer                          │
│  DR Settlement Account                          │
│  CR Customer Account                            │
│  • Transmitted, awaiting collection             │
│  • Included in next file generation             │
└─────────────────────────────────────────────────┘
                    ↓ SETTLE
┌─────────────────────────────────────────────────┐
│  SETTLE: SETTLED Layer                          │
│  DR Settlement Account                          │
│  CR Customer Account                            │
│  • Final settlement, funds collected            │
│  • Customer account credited                    │
└─────────────────────────────────────────────────┘
```

### Querying Balances

Check account balances across all layers using the standard balance query:

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
      }
      drBalance {
        units
      }
    }
    pending {
      crBalance {
        units
      }
      drBalance {
        units
      }
    }
    encumbrance {
      crBalance {
        units
      }
      drBalance {
        units
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

### Workflow State Transitions

**PUSH Workflow States:**
1. `CREATE` - Reserve funds on customer account (encumbrance)
2. `SUBMIT` - Move to settled, queue for file generation
3. `RETURN` - Reverse settlement if RDFI returns transaction
4. `CANCEL` - Reverse encumbrance before file generation
5. `REIMBURSE_FEE` - Refund processing fee on cancellation/return

**PULL Workflow States:**
1. `CREATE` - Reserve collection on settlement account (encumbrance)
2. `SUBMIT` - Move to pending, queue for file generation
3. `SETTLE` - Move to settled after successful collection
4. `RETURN` - Reverse pending/settled if RDFI returns transaction
5. `CANCEL` - Reverse encumbrance before file generation
6. `REIMBURSE_FEE` - Refund processing fee on cancellation/return

## File Generation

File generation collects all submitted transactions and creates NACHA-formatted files for transmission to your financial institution or the Federal Reserve.

### Generating Files

Use [`Mutation.ach.generateFile()`]($gql:mutation:ach.generateFile) to create ACH files:

```graphql
mutation GenerateODFIFile {
  ach {
    generateFile(
      input: {
        configId: "b96d358e-50b8-4ae5-8b07-2e8f33f396c6"
        fileKey: "outgoing-ach-20251114.ach"
        fileType: ODFI
        generateEmpty: false
      }
    ) {
      fileKey
      generated
    }
  }
}
```

**File Types:**
- `ODFI`: Combined file with both PUSH and PULL transactions
- `ODFI_PUSH_ONLY`: Credit transactions only
- `ODFI_PULL_ONLY`: Debit transactions only

**Generation Control:**
- `generateEmpty: true`: Always creates file (empty NACHA if no transactions)
- `generateEmpty: false`: Only creates file when transactions exist

### What Gets Included

File generation automatically collects:
- All PUSH transactions in SUBMIT state
- All PULL transactions in SUBMIT state (pending settlement)
- Transactions grouped by effective date and SEC code
- Proper batch headers and control totals
- File-level hash totals and entry counts

### NACHA File Structure

Generated files contain:

1. **File Header Record** - Configuration from `odfiHeaderConfiguration`:
   - Immediate destination (your ODFI routing number)
   - Immediate origin (your company ID)
   - File creation date/time
   - File ID modifier

2. **Batch Header Records** - One per batch:
   - Service class code (credits, debits, or mixed)
   - Company name and entry description
   - Effective entry date
   - ODFI identification

3. **Entry Detail Records** - One per transaction:
   - Transaction code (credit/debit, checking/savings)
   - RDFI routing number
   - Account number
   - Amount
   - Individual name
   - Trace number

4. **Batch Control Records** - Validates batch totals
5. **File Control Record** - Validates file totals

### Download Generated Files

After generation, download files using [`Mutation.files.createDownload()`]($gql:mutation:files.createDownload):

```graphql
mutation DownloadODFIFile {
  files {
    createDownload(
      key: "outgoing-ach-20251114.ach"
    ) {
      downloadURL
    }
  }
}
```

Then download via HTTP GET:

```bash
curl '<downloadURL>' -o outgoing-ach-20251114.ach
```

Transmit the downloaded file to your ODFI via SFTP, FTPS, or your institution's preferred method.

### File Generation Timing

**Best Practices:**
- Generate files after your daily processing cutoff
- Allow sufficient time for file transmission before ODFI deadlines
- Consider timezone settings in your ACH configuration
- Generate separate files for different effective dates

**Standard ACH Deadlines:**
- **Standard ACH**: Submit by 6:00 PM ET for next-day settlement
- **Same-Day ACH**: Multiple submission windows (10:30 AM, 2:45 PM ET)
- **Weekend Processing**: Transactions submitted Friday settle Monday

## Return Processing

Return processing handles rejected transactions from RDFIs. When your ODFI provides a return file, Twisp automatically matches returns to original transactions and reverses ledger entries.

### Return Flow

**1. Upload Return File**

When you receive a return file from your ODFI, upload it using [`Mutation.files.createUpload()`]($gql:mutation:files.createUpload):

```graphql
mutation CreateReturnUpload {
  files {
    createUpload(
      input: {
        key: "return-20251114.ach"
        uploadType: ACH
        contentType: "text/plain"
      }
    ) {
      uploadURL
    }
  }
}
```

Upload the file content:

```bash
curl -T return-20251114.ach -XPUT '<uploadURL>'
```

**2. Process Return File**

Process the return file with `ODFI_RETURN` file type using [`Mutation.ach.processFile()`]($gql:mutation:ach.processFile):

```graphql
mutation ProcessReturnFile {
  ach {
    processFile(
      input: {
        configId: "b96d358e-50b8-4ae5-8b07-2e8f33f396c6"
        fileKey: "return-20251114.ach"
        fileType: ODFI_RETURN
      }
    ) {
      fileId
    }
  }
}
```

**3. Automatic Return Matching**

Twisp automatically:
- Extracts trace numbers from return entries
- Matches returns to original workflow executions
- Executes RETURN workflow state for each transaction
- Reverses appropriate balance layer entries
- Updates workflow execution history
- Creates complete audit trail via `AchWorkflowTrace`

**4. Monitor Return Processing**

Query file processing status using [`Query.ach.file()`]($gql:query:ach.file):

```graphql
query GetReturnFileStatus {
  ach {
    file(
      fileKey: "return-20251114.ach"
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

### Return Ledger Accounting

Returns automatically reverse the appropriate balance layer entries.

**PUSH Return (Credit Returned):**

If a PUSH transaction in settled layer is returned:

```
Reverse original SUBMIT:
DR Settlement Account (Settled)       $1,500.00
CR Customer Account (Settled)         $1,500.00
```

Funds are returned to the customer account.

**PULL Return (Debit Returned):**

If a PULL transaction in pending layer is returned:

```
Reverse original SUBMIT:
DR Customer Account (Pending)         -$250.00
CR Settlement Account (Pending)       -$250.00
```

If a PULL transaction in settled layer is returned (after SETTLE executed):

```
Reverse original SETTLE:
DR Customer Account (Settled)         -$250.00
CR Settlement Account (Settled)       -$250.00
```

### Common Return Codes

Returns contain NACHA return codes indicating rejection reason:

| Code | Reason | Common Cause |
|------|--------|--------------|
| `R01` | Insufficient Funds | Receiver account has insufficient balance |
| `R02` | Account Closed | Receiver account has been closed |
| `R03` | No Account / Unable to Locate | Account number not found at RDFI |
| `R04` | Invalid Account Number | Account number fails validation |
| `R05` | Unauthorized Debit | Consumer did not authorize debit |
| `R07` | Authorization Revoked | Consumer revoked authorization |
| `R08` | Payment Stopped | Receiver placed stop payment |
| `R10` | Customer Advises Not Authorized | Receiver claims transaction unauthorized |
| `R29` | Corporate Customer Not Authorized | Corporate receiver did not authorize |

### Return Timing

**Standard Return Windows:**
- **Most returns**: Within 2 business days of settlement
- **Unauthorized returns** (R05, R07, R10): Up to 60 days after settlement
- **Admin returns** (R02, R03, R04): Within 2 business days

**Late Returns:**

Some returns arrive after the standard 2-day window:
- Still processed automatically by Twisp
- May require manual reconciliation with your ODFI
- Check `processingDetail` for any matching issues

### Handling Return Fees

When a transaction is returned, you may want to reimburse processing fees:

```graphql
mutation ReimburseFee {
  workflow {
    execute(
      input: {
        executionId: "c97685f6-c1b2-11f0-a959-069b540ea27c"
        code: "ACH_PUSH"
        task: "REIMBURSE_FEE"
        params: {}
      }
    ) {
      executionId
      output {
        state
      }
    }
  }
}
```

**Ledger Entries:**
```
DR Fee Account (Settled)              $0.25
CR Customer Account (Settled)         $0.25
```

### Return Reconciliation

Query workflow execution to see complete return history:

```graphql
query GetReturnExecution {
  workflow {
    execution(
      executionId: "c97685f6-c1b2-11f0-a959-069b540ea27c"
    ) {
      workflowId
      task
      params
      activities {
        action
        entity {
          ... on Transaction {
            transactionId
            description
            entries(first:100) {
              nodes {
                accountId
                layer
                amount
              }
            }
          }
          ... on AchWorkflowTrace {
            traceNumber
            fileId
          }
        }
      }
    }
  }
}
```

The `AchWorkflowTrace` entity links the return to the original transaction via trace number, enabling complete auditability.

## Transaction Tracking

**Workflow Execution:**

Each workflow execution receives unique `executionId` for tracking. Execution record contains:
- Workflow and task identifiers
- Input parameters
- Output state
- All created ledger transactions
- ACH workflow traces with file and entry details
- Complete history of state transitions

**ACH Workflow Trace:**

Links workflow execution to ACH file entries via:
- Trace number (15-digit NACHA identifier)
- File ID and record ID
- Configuration ID and version
- Workflow and execution identifiers

Enables return matching, reconciliation, and complete transaction lineage tracking.

## Example : End-to-End ODFI Setup and Payment

Complete flow from configuration through file transmission:

```graphql
# 1. Create required accounts
mutation CreateAccounts {
  # Settlement account - where funds transit during ACH processing
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

  # Customer account - for testing
  customer: createAccount(
    input: {
      accountId: "d2f7183f-8e9c-45e7-9a98-ef1897ddb930"
      code: "customer.001"
      name: "Customer Account"
      normalBalanceType: DEBIT
      config: {
        enableConcurrentPosting: true
      }
    }
  ) {
    accountId
  }
}

# 2. Create a journal for ACH transactions
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

# 3. Create a webhook endpoint (Note: webhook endpoint not used in ODFI-only use cases)
mutation CreateACHWebhookProcessor {
  events {
    createEndpoint(
      input: {
        endpointId: "b84512f1-a67e-4dc2-94dd-66c48b4d13fb"
        status: ENABLED
        endpointType: ACH_PROCESSOR
        url: "https://webhook.site/ach-testing"
        subscription: []
        description: "ACH webhook processor"
      }
    ) {
      endpointId
    }
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
          immediateDestinationName: "Federal Reserve Bank"
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

mutation CreatePullTransaction {
  workflow {
    execute(
      input: {
        executionId: "c97685f6-c1b2-11f0-a959-069b540ea27c"
        code: "ACH_PULL"
        task: "CREATE"
        params: {
          configId: "b96d358e-50b8-4ae5-8b07-2e8f33f396c6"
          accountId: "d2f7183f-8e9c-45e7-9a98-ef1897ddb930"
          settlementAccountId: "37f7e8a6-171f-411d-ad59-7b1f40f505ea"
          journalId: "8d7e6f5a-4b3c-2d1e-0f9a-8b7c6d5e4f3a"
          routingNumber:       "026009593",
          accountNumber:       "12345678901234567",
          accountType:         "checking",
          individualName:      "Clark Kent",
          entryDescription:    "XFER",
          entryClassCode:"PPD"
          effective:"2025-11-14"
          amount: "250.00"
          feeAccountId: "5b9e3c2f-4d0e-5a8f-b6c7-2e9f0a1b3c4d"
          feeAmount: "0.25"
          correlationId: "billing-batch-001"
          metadata: {
            customerId: "cust-456"
            invoiceNumber: "INV-2025-001"
          }
          entryMetadata: {
            companyName: "ACME Corp"
            companyEntryDescription: "INVOICE"
            individualName: "Jane Smith"
            standardEntryClassCode: "WEB"
            dfiAccountNumber: "987654321"
            identificationNumber: "inv-001"
            rdfiIdentification: "021000021"
          }
        }
      }
    ) {
      executionId
      output {
        state
      }
    }
  }
}

mutation GenerateODFIFile {
  ach {
    generateFile(
      input: {
        configId: "b96d358e-50b8-4ae5-8b07-2e8f33f396c6"
        fileKey: "outgoing-ach-20251114.ach"
        fileType: ODFI
        generateEmpty: false
      }
    ) {
      fileKey
      generated
    }
  }
}

mutation DownloadODFIFile {
  files {
    createDownload(
      key: "outgoing-ach-20251114.ach"
    ) {
      downloadURL
    }
  }
}
```


## Best Practices

### Security

**API Key Management:**
- Rotate API keys regularly
- Use separate keys for production and development
- Store keys in secure secret management systems (AWS Secrets Manager, HashiCorp Vault)
- Never commit API keys to source control

**Account Isolation:**
- Use separate ACH configurations for different business units
- Implement proper account-level access controls
- Audit all ACH operations via workflow execution logs

**Webhook Security (if using RDFI):**
- Validate webhook signatures to ensure requests are from Twisp
- Use HTTPS endpoints for all webhooks
- Implement rate limiting and DDoS protection

### Performance

**Batch Processing:**
- Group transactions by effective date for efficient file generation
- Use correlation IDs to track related transactions
- Process CREATE and SUBMIT in parallel where possible

**Account Configuration:**
- Enable `enableConcurrentPosting: true` on all ACH accounts
- This supports high-volume parallel transaction posting
- Settlement account especially critical for concurrent access

**File Generation Timing:**
- Generate files during off-peak hours when possible
- Separate file generation from transaction creation
- Monitor file generation performance via [`Query.ach.file()`]($gql:query:ach.file)

### Monitoring

**Alert on Key Metrics:**
- High return rates (> 2% indicates data quality issues)
- Failed workflow executions
- File generation failures
- Settlement account balance anomalies

**Daily Reconciliation:**
- Compare file totals with workflow execution totals
- Reconcile settlement account with ODFI reports
- Track return rates by SEC code and effective date

**Execution Tracking:**
- Use correlation IDs to group related transactions
- Query workflow executions for complete audit trails
- Monitor workflow state transitions for stuck transactions

### Compliance

**Record Retention:**
- Maintain workflow execution history for 7 years minimum
- Store generated ACH files and return files
- Keep complete audit trail of all balance changes
- Archive transaction metadata and authorization records

**NACHA Rules:**
- Adhere to return timeframes (2 business days for most codes)
- Process returns within 24 hours of receipt
- Maintain proper SEC codes for transaction types
- Follow authorization requirements for debit transactions

**Authorization Management:**
- Store proof of authorization for debit transactions
- Support authorization revocation (R07 returns)
- Implement stop payment capabilities (R08 returns)
- Handle unauthorized transaction disputes (R10 returns)

**Reg E Compliance:**
- Honor consumer dispute rights (60-day investigation period)
- Provide proper disclosures for recurring debits
- Implement error resolution procedures
- Maintain consumer authorization records

### Error Handling

**Transaction Failures:**
- Monitor workflow execution failures
- Implement retry logic for transient errors
- Use CANCEL state to reverse failed transactions
- Alert operations team for manual intervention

**Return Handling:**
- Process return files within 24 hours of receipt
- Automatically reverse transactions via RETURN state
- Consider reimbursing fees for returned transactions
- Track return rates to identify systemic issues

**File Generation Issues:**
- Validate all transactions before SUBMIT
- Test file generation with `generateEmpty: true` initially
- Monitor file processing status via [`Query.ach.file()`]($gql:query:ach.file)
- Keep backup of generated files before transmission

## ODFI Operations

Use GraphQL workflow API for all ODFI operations:

- [`Mutation.workflow.execute()`]($gql:mutation:workflow.execute): Execute workflow state transitions
- [`Query.workflow.execution()`]($gql:query:workflow.execution): Get workflow execution details
- [`Query.workflow.executions()`]($gql:query:workflow.executions): Query workflow executions

File operations:
- [`Mutation.ach.generateFile()`]($gql:mutation:ach.generateFile): Generate NACHA file
- [`Mutation.ach.processFile()`]($gql:mutation:ach.processFile): Process return files
- [`Query.ach.file()`]($gql:query:ach.file): Query file status

## Further Reading

To learn ODFI operations from scratch, see the tutorial on [Your First ACH Payment](/tutorials/ach/first-ach-payment).

For production payment workflows, see the how-to guide on [Processing ACH Payments](/guides/processing-ach-payments).

For return management, see the how-to guide on [Handling ACH Returns](/guides/handling-ach-returns).

For complete GraphQL type definitions, see:
- []($gql:object:WorkflowExecution)
- []($gql:object:AchWorkflowTrace)
- []($gql:input:WorkflowExecuteInput)
