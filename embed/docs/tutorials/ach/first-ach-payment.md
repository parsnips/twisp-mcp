---
title: Your First ACH Payment
metaTitle: "Tutorial: Your First ACH Payment"
description: Send your first ACH credit payment
showNext: true
showPrev: true
---

In this tutorial, we will send our first ACH payment. We'll create a customer account, initiate a $50 credit payment, and generate an ACH file ready for transmission.

## What We'll Build

By completing this tutorial, we will:
- Create a customer account
- Execute an ACH PUSH (credit) workflow
- Generate an ACH file
- Download the file for transmission
- Verify the ledger entries

This tutorial takes approximately 10 minutes to complete.

## Prerequisites

Before we begin, you must have completed:
- [Setting Up ACH Processing](/tutorials/ach/setting-up-ach) - This tutorial requires the IDs from setup

You'll need these IDs from the setup tutorial:
```
Journal ID:     8d7e6f5a-4b3c-2d1e-0f9a-8b7c6d5e4f3a
Settlement:     37f7e8a6-171f-411d-ad59-7b1f40f505ea
Fee:            5b9e3c2f-4d0e-5a8f-b6c7-2e9f0a1b3c4d
Config ID:      fe27128a-b331-4e0e-94f8-9a32443fee36
```

## Step 1: Create a Customer Account

First, we'll create an account for our customer who will receive the ACH payment.

```graphql
mutation CreateCustomer {
  createAccount(
    input: {
      accountId: "c4f7e3a2-8b1d-4e9f-a5c6-2d3e4f5a6b7c"
      code: "customer.alice.checking"
      name: "Alice Smith - Checking"
      normalBalanceType: CREDIT
      config: {
        enableConcurrentPosting: true
      }
    }
  ) {
    accountId
    name
    code
  }
}
```

**What Just Happened:**
We created a checking account for Alice Smith with ID `c4f7e3a2-8b1d-4e9f-a5c6-2d3e4f5a6b7c`. This account represents Alice's bank account in our system.

**Expected Response:**
```json
{
  "data": {
    "createAccount": {
      "accountId": "c4f7e3a2-8b1d-4e9f-a5c6-2d3e4f5a6b7c",
      "name": "Alice Smith - Checking",
      "code": "customer.alice.checking"
    }
  }
}
```

## Step 2: Create the ACH Payment

Now we'll create an ACH credit payment for $50. This uses the ACH PUSH workflow, which is for sending money to customers.

```graphql
mutation CreatePayment {
  workflow {
    execute(
      input: {
        executionId: "55af980c-c1bb-11f0-833c-069b540ea27c"
        code: "ACH_PUSH"
        task: "CREATE"
        params: {
          accountId: "c4f7e3a2-8b1d-4e9f-a5c6-2d3e4f5a6b7c"
          settlementAccountId: "37f7e8a6-171f-411d-ad59-7b1f40f505ea"
          journalId: "8d7e6f5a-4b3c-2d1e-0f9a-8b7c6d5e4f3a"
          configId: "fe27128a-b331-4e0e-94f8-9a32443fee36"
          routingNumber:       "026009593",
          accountNumber:       "12345678901234567",
          accountType:         "checking",
          individualName:      "Alice Smith",
          entryDescription:    "XFER",
          amount: "50.00"
          effective: "2025-11-15"
          feeAccountId: "5b9e3c2f-4d0e-5a8f-b6c7-2e9f0a1b3c4d"
          feeAmount: "0.25"
          correlationId: "tutorial-payment-001"
          metadata: "{\"customer\": \"Alice Smith\", \"purpose\": \"Tutorial Payment\"}"
          entryMetadata: "{\"description\": \"Tutorial Payment\"}"
        }
      }
    ) {
      executionId
      task
      output {
        state
      }
    }
  }
}
```

**What Just Happened:**
We executed the CREATE task of the ACH PUSH workflow. This:
- Created an pending transaction (hold) on Alice's account for $50
- Charged a $0.25 processing fee
- Assigned execution ID for tracking

**Key Parameters:**
- `code: "ACH_PUSH"` - This is the ACH PUSH workflow UUID
- `task: "CREATE"` - First state: create and encumber funds
- `amount: "50.00"` - $50 payment
- `effective: "2025-11-15"` - The date the payment should settle (use tomorrow's date in production)
- `correlationId: "tutorial-payment-001"` - Unique identifier for tracking

**Expected Response:**
```json
{
  "data": {
    "workflow": {
      "execute": {
        "executionId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "task": "CREATE",
        "output": {
          "state": "CREATE"
        }
      }
    }
  }
}
```

**Save the Execution ID:**
Copy the `executionId` from the response. We'll need it in the next step:
```
executionId: a1b2c3d4-e5f6-7890-abcd-ef1234567890
```

## Step 3: Check the Balance

Let's verify the pending balance was created by checking Alice's account balance.

```graphql
query CheckBalance {
  balance(
    accountId: "c4f7e3a2-8b1d-4e9f-a5c6-2d3e4f5a6b7c"
    journalId: "8d7e6f5a-4b3c-2d1e-0f9a-8b7c6d5e4f3a"
    currency: "USD"
  ) {
    settled {
      drBalance {
        formatted(as:{locale:"en-US"})
      }
      crBalance {
        formatted(as:{locale:"en-US"})
      }
    }
    pending {
      drBalance {
        formatted(as:{locale:"en-US"})
      }
      crBalance {
        formatted(as:{locale:"en-US"})
      }
    }
    available (layer:PENDING) {
      drBalance {
        formatted(as:{locale:"en-US"})
      }
      crBalance {
        formatted(as:{locale:"en-US"})
      }
      normalBalance {
        formatted(as:{locale:"en-US"})
      }
    }
  }
}
```

**Expected Response:**
```json
{
  "data": {
    "balance": {
      "settled": {
        "drBalance": {
          "formatted": "$0.25"
        },
        "crBalance": {
          "formatted": "$0.00"
        }
      },
      "pending": {
        "drBalance": {
          "formatted": "$50.00"
        },
        "crBalance": {
          "formatted": "$0.00"
        }
      },
      "available": {
        "drBalance": {
          "formatted": "$50.25"
        },
        "crBalance": {
          "formatted": "$0.00"
        },
        "normalBalance": {
          "formatted": "-$50.25"
        }
      }
    }
  }
}
```

**What This Means:**
- **Settled Layer**: $0.25 debit (the fee was charged immediately)
- **Pending Layer**: $50 debit (funds are reserved, not yet settled)
- **Available Balance**: -$50.25 (fees + pending amount)

The negative available balance indicates Alice would need $50.25 in her account for this transaction to clear.

## Step 4: Generate the ACH File

Now we'll generate an ACH file containing our payment.

```graphql
mutation GenerateFile {
  ach {
    generateFile(
      input: {
        configId: "fe27128a-b331-4e0e-94f8-9a32443fee36"
        fileKey: "tutorial-payment-20251114.ach"
        fileType: ODFI_PUSH_ONLY
        generateEmpty: false
      }
    ) {
      fileKey
      generated
    }
  }
}
```

Now we'll submit the payment for inclusion in an ACH file. This moves the payment from pending to settled and invokes the following workflow automatically:

```graphql
mutation SubmitPayment {
  workflow {
    execute(
      input: {
        executionId: "55af980c-c1bb-11f0-833c-069b540ea27c"
        code: "ACH_PUSH"
        task: "SUBMIT"
        params: {
          effective: "2025-11-15"
        }
      }
    ) {
      executionId
      task
      output {
        state
      }
    }
  }
}
```

**What Just Happened:**
We requested generation of an ACH file that includes:
- Reversed the pending
- Posted to the settled layer
- All submitted PUSH (credit) transactions
- NACHA-formatted, ready for transmission
- Stored with key `tutorial-payment-20251114.ach`

**Expected Response:**
```json
{
  "data": {
    "ach": {
      "generateFile": {
        "fileKey": "tutorial-payment-20251114.ach",
        "generated": true
      }
    }
  }
}
```

**If `generated: false`:**
This means no transactions were ready for file generation.

## Step 5: Verify the Settled Balance

Let's check the balance again to see the settled transaction.

```graphql
query CheckBalance {
  balance(
    accountId: "c4f7e3a2-8b1d-4e9f-a5c6-2d3e4f5a6b7c"
    journalId: "8d7e6f5a-4b3c-2d1e-0f9a-8b7c6d5e4f3a"
    currency: "USD"
  ) {
    settled {
      drBalance {
        formatted(as:{locale:"en-US"})
      }
      crBalance {
        formatted(as:{locale:"en-US"})
      }
    }
    pending {
      drBalance {
        formatted(as:{locale:"en-US"})
      }
      crBalance {
        formatted(as:{locale:"en-US"})
      }
    }
    available (layer:PENDING) {
      drBalance {
        formatted(as:{locale:"en-US"})
      }
      crBalance {
        formatted(as:{locale:"en-US"})
      }
      normalBalance {
        formatted(as:{locale:"en-US"})
      }
    }
  }
}
```

**Expected Response:**
```json
{
  "data": {
    "balance": {
      "settled": {
        "drBalance": {
          "formatted": "$50.25"
        },
        "crBalance": {
          "formatted": "$0.00"
        }
      },
      "pending": {
        "drBalance": {
          "formatted": "$0.00"
        },
        "crBalance": {
          "formatted": "$0.00"
        }
      },
      "available": {
        "drBalance": {
          "formatted": "$50.25"
        },
        "crBalance": {
          "formatted": "$0.00"
        },
        "normalBalance": {
          "formatted": "-$50.25"
        }
      }
    }
  }
}
```

**What Changed:**
- **Settled Layer**: Now $50.25 debit ($50 payment + $0.25 fee)
- **Pending Layer**: Cleared to $0 (funds moved to settled)
- **Available Balance**: $50.25 debit (the finalized transaction amount)

The payment is now finalized and in a file on the way to the Federal Reserve!

## Step 6: Download the ACH File

Now we'll get a download URL for our generated file.

```graphql
mutation GetDownloadURL {
  files {
    createDownload(
      key: "tutorial-payment-20251114.ach"
    ) {
      downloadURL
      downloadURLExpiration
    }
  }
}
```

**Expected Response:**
```json
{
  "data": {
    "files": {
      "createDownload": {
        "downloadURL": "https://s3.amazonaws.com/bucket/path?signature=...",
        "downloadURLExpiration": "2025-11-14T15:30:00Z"
      }
    }
  }
}
```

**Download the File:**
Use the URL to download the file (replace `<downloadURL>` with actual URL from response):

```bash
curl '<downloadURL>' -o tutorial-payment-20251114.ach
```

## Step 7: Inspect the ACH File

Let's look at what's in the file we generated. Open `tutorial-payment-20251114.ach` in a text editor.

You'll see a NACHA-formatted file with lines like:

```
101 021000021 1234567892511150055A094101Test Bank              Your Company                   
5220Your Company                        1234567890PPDXFER      251114251117   1123456780000001
622026009593123456789012345670000005000               Alice Smith             0123456780000001
822000000100026009590000000000000000000050001234567890                         123456780000001
9000001000001000000010002600959000000000000000000005000                                       
9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999
9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999
9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999
9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999
9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999
```

**What Each Line Means:**
- **Line 1 (101...)**: File Header - Contains ODFI info
- **Line 2 (5225...)**: Batch Header - Contains company info
- **Line 3 (622...)**: Entry Detail - The actual $50 payment to Alice
- **Line 4 (822...)**: Batch Control - Validates batch totals
- **Line 5 (9000...)**: File Control - Validates file totals

The file is ready to transmit to your ODFI!

## Step 9: View the Workflow Execution

Let's review the complete workflow execution to see all the ledger entries.

```graphql
query ViewExecution {
  workflow {
    execution(
      executionId: "55af980c-c1bb-11f0-833c-069b540ea27c"
    ) {
      executionId
      workflowId
      task
      params
      activities {
        action
        entityType
        entity {
          ... on Transaction {
            transactionId
            effective
            entries(first:10) {
              nodes {
                entryId
                accountId
                layer
                direction
                units
              }
            }
          }
        }
      }
    }
  }
}
```

**What You'll See:**
This shows all ledger transactions created by the workflow:
- CREATE task transactions (pending + fee)
- SUBMIT task transactions (reverse pending + settle)

Each activity shows the double-entry accounting entries on customer and settlement accounts.

## What We Accomplished

Let's review what we did:

1. ✅ Created a customer account
2. ✅ Created a $50 ACH credit payment
3. ✅ Verified the pending transaction was created
4. ✅ Generated an ACH file
5. ✅ Verified the payment settled
6. ✅ Downloaded the ACH file
7. ✅ Inspected the NACHA file format
8. ✅ Reviewed the workflow execution

## Key Concepts We Learned

**Workflow States:**
- CREATE: Pending funds, charge fees
- SUBMIT: Settle funds, mark for file generation

**Balance Layers:**
- Pending: Temporary holds
- Settled: Finalized transactions
- Available: Calculated from settled - pending

**File Generation:**
- Collects all submitted transactions
- Generates NACHA-formatted files
- Ready for ODFI transmission

## What Happens Next?

In production, you would:
1. Transmit the generated file to your ODFI via SFTP
2. Wait for confirmation (usually 1-2 business days)
3. Process any returns received from the RDFI

For this tutorial, you now have:
- A working ACH payment workflow
- A valid NACHA file
- Understanding of the complete process

## Try It Again

Now that you understand the process, try creating another payment:
1. Create a new customer account (use a different UUID)
2. Execute CREATE with a different amount
3. Submit the payment
4. Generate a new file

Each time you generate a file, all submitted payments will be included.

## Troubleshooting

**Problem: "Execution not found"**

Make sure you're using the correct `executionId` from the CREATE response in the SUBMIT mutation.

**Problem: "File generated: false"**

This means no transactions were ready. Check:
- Did you run the CREATE task?
- Is the effective date in the future or today?
- Are there any errors in the workflow execution?

**Problem: "Cannot download file"**

The download URL expires after a short time. If it's expired:
- Run the `createDownload` mutation again to get a new URL

**Problem: "Balance doesn't match"**

After CREATE:
- Pending should show $50 DR
- Settled should show $0.25 DR (fee)

After SUBMIT:
- Pending should be $0
- Settled should show $50.25 DR

## Next Steps

Now that you've sent your first ACH payment, explore more:

**Processing Guides:**
- [Processing ACH Payments](/guides/processing-ach-payments) - Production payment workflows
- [Handling ACH Returns](/guides/handling-ach-returns) - Managing payment returns
- [Reconciling ACH Files](/guides/reconciling-ach-files) - File validation and reconciliation

**Reference Documentation:**
- [ODFI Reference](/reference/ach/odfi) - Complete ODFI API documentation
- [File Operations](/reference/ach/file-operations) - File upload and download operations
- [Configuration](/reference/ach/configuration) - ACH configuration reference

## Summary

Congratulations! You've successfully:
- Created your first ACH payment
- Executed a complete workflow from creation to settlement
- Generated a NACHA-formatted ACH file
- Understood the balance layer transitions
- Downloaded a production-ready ACH file

You now have the foundation to build production ACH payment systems on Twisp!
