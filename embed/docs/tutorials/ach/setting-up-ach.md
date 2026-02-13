---
title: Setting Up ACH Processing
metaTitle: "Tutorial: Setting Up ACH Processing"
description: Learn how to set up ACH processing from scratch
showNext: true
showPrev: true
---

In this tutorial, we will set up ACH processing on the Twisp platform. By the end, you will have a fully configured ACH processor ready to send and receive ACH transactions.

## What We'll Build

By completing this tutorial, we will:
- Create all required accounts for ACH processing
- Set up a journal for ACH transactions
- Configure a webhook endpoint for transaction decisioning
- Create an ACH configuration
- Verify everything works correctly

This tutorial takes approximately 15 minutes to complete.

## Prerequisites

Before we begin, you need:
- A Twisp account with API access
- Your GraphQL API endpoint URL
- API credentials (authentication token)

You do NOT need:
- An ODFI relationship (that comes later for production)
- Understanding of double-entry accounting
- Prior ACH experience

## Step 1: Create the Journal

We will start by creating a journal where all ACH transactions will be posted.

```graphql
mutation CreateJournal {
  createJournal(
    input: {
      journalId: "8d7e6f5a-4b3c-2d1e-0f9a-8b7c6d5e4f3a"
      name: "ACH Processing Journal"
      status: ACTIVE
    }
  ) {
    journalId
    name
    status
  }
}
```

**What Just Happened:**
We created a journal with ID `8d7e6f5a-4b3c-2d1e-0f9a-8b7c6d5e4f3a`. This journal will track all ACH-related accounting entries.

**Expected Response:**
```json
{
  "data": {
    "createJournal": {
      "journalId": "8d7e6f5a-4b3c-2d1e-0f9a-8b7c6d5e4f3a",
      "name": "ACH Processing Journal",
      "status": "ACTIVE"
    }
  }
}
```

## Step 2: Create Required Accounts

Next, we will create four accounts required for ACH processing. We'll create them all at once.

```graphql
mutation CreateACHAccounts {
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
    name
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
    name
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
    name
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
    name
  }
}
```

**What Just Happened:**
We created four accounts:

1. **Settlement Account** (`37f7e8a6-171f-411d-ad59-7b1f40f505ea`) - All ACH funds flow through this account during processing
2. **Suspense Account** (`3171b0c2-6e9f-41aa-a5a6-ee927deb27cf`) - Holds transactions when the destination account can't be found
3. **Exception Account** (`4a8f2b1e-3c9d-4f7e-a5b6-1d8e9f0a2b3c`) - Receives transactions that fail processing rules
4. **Fee Account** (`5b9e3c2f-4d0e-5a8f-b6c7-2e9f0a1b3c4d`) - Collects ACH processing fees

All accounts have `enableConcurrentPosting: true` to support high transaction volumes.

**Expected Response:**
```json
{
  "data": {
    "settlement": {
      "accountId": "37f7e8a6-171f-411d-ad59-7b1f40f505ea",
      "name": "ACH Settlement"
    },
    "suspense": {
      "accountId": "3171b0c2-6e9f-41aa-a5a6-ee927deb27cf",
      "name": "ACH Suspense"
    },
    "exception": {
      "accountId": "4a8f2b1e-3c9d-4f7e-a5b6-1d8e9f0a2b3c",
      "name": "ACH Exception"
    },
    "fee": {
      "accountId": "5b9e3c2f-4d0e-5a8f-b6c7-2e9f0a1b3c4d",
      "name": "ACH Fee Income"
    }
  }
}
```

## Step 3: Create the Webhook Endpoint

Now we will create a webhook endpoint. This endpoint will receive requests for transaction decisioning when ACH files are processed.

For this tutorial, we'll use a test endpoint URL. In production, you'll replace this with your actual webhook server.

```graphql
mutation CreateWebhook {
  events {
    createEndpoint(
      input: {
        endpointId: "b84512f1-a67e-4dc2-94dd-66c48b4d13fb"
        status: ENABLED
        endpointType: ACH_PROCESSOR
        url: "https://webhook.site/unique-url-here"
        subscription: []
        description: "ACH decisioning webhook"
      }
    ) {
      endpointId
      url
      status
    }
  }
}
```

**What Just Happened:**
We created a webhook endpoint with ID `b84512f1-a67e-4dc2-94dd-66c48b4d13fb`. When ACH files are processed, Twisp will send POST requests to the URL we specified.

**For This Tutorial:**
Use `https://webhook.site` to create a free test webhook URL:
1. Go to https://webhook.site
2. Copy the unique URL shown
3. Use that URL in the mutation above

**Expected Response:**
```json
{
  "data": {
    "events": {
      "createEndpoint": {
        "endpointId": "b84512f1-a67e-4dc2-94dd-66c48b4d13fb",
        "url": "https://webhook.site/your-unique-url",
        "status": "ENABLED"
      }
    }
  }
}
```

## Step 4: Create the ACH Configuration

Now we will tie everything together by creating an ACH configuration. This tells the ACH processor which accounts to use and where to send webhooks.

```graphql
mutation CreateACHConfig {
  ach {
    createConfiguration(
      input: {
        configId: "fe27128a-b331-4e0e-94f8-9a32443fee36"
        endpointId: "b84512f1-a67e-4dc2-94dd-66c48b4d13fb"
        journalId: "8d7e6f5a-4b3c-2d1e-0f9a-8b7c6d5e4f3a"
        settlementAccountId: "37f7e8a6-171f-411d-ad59-7b1f40f505ea"
        suspenseAccountId: "3171b0c2-6e9f-41aa-a5a6-ee927deb27cf"
        exceptionAccountId: "4a8f2b1e-3c9d-4f7e-a5b6-1d8e9f0a2b3c"
        feeAccountId: "5b9e3c2f-4d0e-5a8f-b6c7-2e9f0a1b3c4d"
        odfiHeaderConfiguration: {
          immediateDestination: "021000021"
          immediateDestinationName: "Test Bank"
          immediateOrigin: "1234567890"
          immediateOriginName: "Your Company"
        }
        timeZone: "America/New_York"
      }
    ) {
      configId
      version
      timeZone
    }
  }
}
```

**What Just Happened:**
We created an ACH configuration that connects:
- The journal we created
- All four accounts
- The webhook endpoint
- ODFI header information for generating ACH files

The `odfiHeaderConfiguration` contains placeholder values. When you're ready for production, you'll replace these with real values from your ODFI (bank).

**Expected Response:**
```json
{
  "data": {
    "ach": {
      "createConfiguration": {
        "configId": "fe27128a-b331-4e0e-94f8-9a32443fee36",
        "version": 1,
        "timeZone": "America/New_York"
      }
    }
  }
}
```

## Step 5: Verify the Configuration

Let's verify everything was created correctly by querying our configuration.

```graphql
query VerifySetup {
  ach {
    configuration(
      id: "fe27128a-b331-4e0e-94f8-9a32443fee36"
    ) {
      configId
      journalId
      settlementAccountId
      suspenseAccountId
      exceptionAccountId
      feeAccountId
      endpointId
      odfiHeaderConfiguration {
        immediateDestination
        immediateOrigin
      }
      timeZone
      version
    }
  }
}
```

**Expected Response:**
```json
{
  "data": {
    "ach": {
      "configuration": {
        "configId": "fe27128a-b331-4e0e-94f8-9a32443fee36",
        "journalId": "8d7e6f5a-4b3c-2d1e-0f9a-8b7c6d5e4f3a",
        "settlementAccountId": "37f7e8a6-171f-411d-ad59-7b1f40f505ea",
        "suspenseAccountId": "3171b0c2-6e9f-41aa-a5a6-ee927deb27cf",
        "exceptionAccountId": "4a8f2b1e-3c9d-4f7e-a5b6-1d8e9f0a2b3c",
        "feeAccountId": "5b9e3c2f-4d0e-5a8f-b6c7-2e9f0a1b3c4d",
        "endpointId": "b84512f1-a67e-4dc2-94dd-66c48b4d13fb",
        "odfiHeaderConfiguration": {
          "immediateDestination": "021000021",
          "immediateOrigin": "1234567890"
        },
        "timeZone": "America/New_York",
        "version": 1
      }
    }
  }
}
```

**What to Check:**
- All IDs match what we created
- Version is 1 (this is the first version)
- Timezone is correct for your location

If everything matches, congratulations! Your ACH processor is configured.

## What We Accomplished

Let's review what we built:

1. ✅ Created a journal for ACH transactions
2. ✅ Created four required accounts (settlement, suspense, exception, fee)
3. ✅ Set up a webhook endpoint
4. ✅ Created an ACH configuration connecting everything
5. ✅ Verified the configuration is correct

## Save These IDs

You'll need these IDs going forward. Save them somewhere safe:

```
Journal ID:     8d7e6f5a-4b3c-2d1e-0f9a-8b7c6d5e4f3a
Settlement:     37f7e8a6-171f-411d-ad59-7b1f40f505ea
Suspense:       3171b0c2-6e9f-41aa-a5a6-ee927deb27cf
Exception:      4a8f2b1e-3c9d-4f7e-a5b6-1d8e9f0a2b3c
Fee:            5b9e3c2f-4d0e-5a8f-b6c7-2e9f0a1b3c4d
Webhook:        b84512f1-a67e-4dc2-94dd-66c48b4d13fb
Config ID:      fe27128a-b331-4e0e-94f8-9a32443fee36
```

## Next Steps

Now that your ACH processor is set up, you're ready to send your first ACH payment!

**Continue to the next tutorial:**
- [First ACH Payment](/tutorials/ach/first-ach-payment) - Send a $50 test payment and generate an ACH file

**Production Checklist:**
Before using ACH in production, you'll need to:
- [ ] Establish a relationship with an ODFI (bank)
- [ ] Get real ODFI header values from your bank
- [ ] Update the `odfiHeaderConfiguration` with real values
- [ ] Set up SFTP/FTPS for file transmission
- [ ] Implement a production webhook server
- [ ] Test with your ODFI's validation process

## Troubleshooting

**Problem: "Account already exists" error**

If you see this error, it means you already have an account with that ID. This is fine! You can either:
- Use the existing accounts and skip account creation
- Choose different UUIDs for your accounts

**Problem: "Webhook endpoint creation failed"**

Common causes:
- Invalid URL format - make sure it starts with `https://`
- Network issue - try again in a moment

**Problem: "Configuration creation failed"**

Check that:
- All referenced IDs (journal, accounts, endpoint) exist
- You haven't already created a configuration with this ID
- All UUIDs are properly formatted

## Summary

You've successfully set up ACH processing! We created:
- A journal for transaction tracking
- Four accounts for different processing scenarios
- A webhook endpoint for transaction decisioning
- An ACH configuration tying everything together

You're now ready to process ACH transactions. The next tutorial will show you how to send your first payment.
