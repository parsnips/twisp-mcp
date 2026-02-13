---
title: "Step 3: Model Intra-Bank Transfers"
description: Zuzu also needs to support bank transfers between customer accounts so that customers can send money to one another.
metaTitle: "Twisp 101 Tutorial: Step 3"
metaDescription: Work through the steps in this tutorial to learn how to start building on the Twisp accounting core.
showNext: true
showPrev: true
---

This type of transaction we'll model with a tran code called `BANK_TRANSFER`, but it's going to be a little more complex than the previous one.

Zuzu isn't just letting customers transfer money for free, all day long. Instead, they'll charge a small percentage fee of 1% with a $10 maximum, paid by the sender.

## Create a revenue account

The fee charged for a bank transfer will be represented as a `DEBIT` against the sender's bank account. The balancing `CREDIT` entry will be posted to Zuzu's revenue account, which doesn't exist yet.

To support bank transfers, then, we'll need to first create the revenue account for Zuzu.

{% partial file="graphql/examples/twisp101/008_CreateRevenueAccount.md" /%}

Now that we have the revenue account, we can design the tran code for internal transfers.

## Define the TranCode for transfers

The tran code for this transaction type needs to do a few things:

- Write entries to move the defined amount from the sender's checking account to the receiver's checking account
- Write an additional two entries to move the fee from the sender's checking account to the Revenue account, using a runtime expression to calculate the fee amount

When posting a transaction, we want to enable the poster to provide the **sender's** account ID, the **receiver's** account ID, the **amount** to transfer, the 1% **fee**, and the **effective date** of the transfer. We'll define each of these as `params` on the TranCode.

{% partial file="graphql/examples/twisp101/009_CreateBankTransferTranCode.md" /%}

{% callout type="tip" %}
Writing clear descriptions for tran codes and their parameters is a great way to help API users can understand what the tran code is for and how to invoke it.
{% /callout %}

## Post a test transaction

Let's test this transaction out by sending $2.25 from Ernie to Bert. To see the results of the transaction as encoded by the tran code, we'll return the entries posted, digging all the way down into the account for each entry.

{% partial file="graphql/examples/twisp101/012_PostBankTransfer.md" /%}

Success! From our response, we can see that each entry was posted to the correct account and for the correct amount.
