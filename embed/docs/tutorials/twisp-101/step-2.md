---
title: "Step 2: Model Deposits and Withdrawals"
description: The most basic transaction types Zuzu needs to support are to allow customers to move money in and out of their checking account.
metaTitle: "Twisp 101 Tutorial: Step 2"
metaDescription: Work through the steps in this tutorial to learn how to start building on the Twisp accounting core.
showNext: true
showPrev: true
---

In other words, we need to support deposits and withdrawals: a customer needs to be able to deposit money _into_ their checking account and withdraw some or all of their balance _out of_ their account. Our ledger will support ACH debit and credit transaction types to model "withdrawals" and "deposits".

To build out this feature we'll need to do three things:

- Create some customer [Accounts]($gql:object:Account) to model money Zuzu holds on behalf of a customer.
- Create an assets account to model Zuzu's cash assets.
- Design two [TranCodes]($gql:object:TranCode) to use as a templates for ACH transactions.

With double-entry accounting, every transaction needs to write _at least_ two entries to the ledger which balance out across debits and credits. *How* (with what metadata) and *where* (to which accounts) these entries are written to is determined by the type of transaction.

{% comment %}
<!-- TODO: add graphic and/or table showing accounts & tran codes -->
{% /comment %}

In Twisp, transaction types are explicitly defined during the design stage by creating transaction codes, or TranCodes.

{% callout type="tip" %}
When a customer deposits money into their account, Zuzu is effectively acting as a custodian of the customer's money. This is why customer accounts are treated as a liability for the company – they represent money that Zuzu _owes_ the customer.

The assets account represents the cash on hand that Zuzu holds at any given time.
{% /callout %}

## Create accounts

First, let's create checking accounts for some sample customers. We can do this with the `createAccount` mutation.

{% partial file="graphql/examples/twisp101/002_CreateCustomerAccounts.md" /%}

Note that customer accounts use a credit-normal balance type because they represent liabilities.

Next, let's create the assets account using a debit-normal balance type:

{% partial file="graphql/examples/twisp101/003_CreateAssetsAccount.md" /%}

## Check account balances

Every account starts with a zero/null balance. We can check the balances of each account by querying the account id and pulling out the account balance for the primary journal we created earlier.

Note that in this example, we use GraphQL variables to store the values used previously and inject them via query params. This makes it easier to re-use values across multiple requests.


{% partial file="graphql/examples/twisp101/004_CheckAccountBalances.md" /%}

Just as expected - balances are `null` for all accounts. Not very exciting. Let's change that.

## Design the transaction type as a TranCode

The only way to write ledger entries in Twisp is by **posting a transaction**. Furthermore, every transaction is structured by the tran code used. This ensures that the ledger is consistent, predictable, and correct.

To define the tran codes for ACH credits and debit transaction types, we need to first determine:

- A unique identifier code for the tran code
- Which accounts will be debited/credited
- What entry data will be written
- How we will create parameterize inputs (for values like the amount)

Let's keep it simple and use the codes `ACH_CREDIT` and `ACH_DEBIT` for these tran codes.

For deposits (i.e. ACH credits), we'll **credit** the customer's checking account because this account is credit-normal and represents Zuzu's obligation to the customer, and we'll **debit** the assets account because this is the debit-normal account which represents how much liquid currency Zuzu has on hand (in this case, on behalf of the customer).

Withdrawals (i.e. ACH debits) are going to be basically the same, but reversed: debit the customer's checking and credit the assets account.

We'll write one entry for the debit and one for the credit, using an entry type to clarify the function of the entry within the context of the transaction.

Finally, we'll need to parameterize both the amount as well as the customer's checking account ID, since these are the salient pieces of information that we want to be able to provide when posting a transaction using this tran code.

### Create the TranCodes for DEPOSIT and WITHDRAW

Now we can create these tran codes with GraphQL, plugging in the design decisions we just made to encode these transaction types.

{% partial file="graphql/examples/twisp101/005_CreateDepositAndWithdrawalTranCodes.md" /%}

## Post a test transaction

With these tran codes defined, we can now post transactions using them.

Let's deposit $9.53 into Ernie's account:

{% partial file="graphql/examples/twisp101/006_PostDeposit.md" /%}

That all looks good.

Now let's withdraw $4.28 from Ernie's account:

{% partial file="graphql/examples/twisp101/007_PostWithdrawal.md" /%}

Great! We've posted our first transactions. Let's go to the next feature.
