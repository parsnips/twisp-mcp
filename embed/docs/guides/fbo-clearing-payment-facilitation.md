---
title: Handling FBO Clearing (Payment Facilitation) with Twisp
description: |
  Follow these tep-by-step instructions on managing FBO clearing using Twisp's accounting tools. This guide covers everything from setting up accounts and defining transaction codes to modeling transactions, monitoring settlements, and optimizing the payment facilitation process. The guide is especially helpful for fintech companies, banks, and organizations looking to streamline their payment facilitation processes, maintain compliance, and improve efficiency.
aiCommand: |
  pnpm ai write "Write a guide for handling FBO clearing (payment facilitation) with Twisp. The guide should explain the process of settling transactions between parties, managing payment processing, and ensuring timely fund disbursement. Only use documented features of the Twisp system; do not invent features." -a GuideWriter -s fbo-clearing-payment-facilitation -o ~/Desktop/guide-fbo-clearing-payment-facilitation.md -m gpt-4
showNext: false
showPrev: false
draft: true
---

{% comment %}

OUTLINE

I. Set up the appropriate accounts
   - Create a settlement account for FBO clearing
   - Create an account for payment processing fees
   - Create accounts for each party involved in the transactions

II. Define transaction codes for FBO clearing
   - Design a tran code for moving funds between parties
   - Design a tran code for payment processing fees
   - Design a tran code for fund disbursement

III. Model payment facilitation transactions
   - Post transactions using the defined tran codes to move funds between parties
   - Deduct payment processing fees using the appropriate tran code
   - Disburse funds to the final destination using the designated tran code

IV. Monitor the settlement process
   - Query the settlement account to track the status of fund transfers
   - Review transactions in the ledger to ensure accurate processing and disbursement

V. Reconcile and settle transactions
   - Perform reconciliation to match transactions with their corresponding disbursements
   - Settle transactions by updating the ledger and moving funds to the appropriate accounts

VI. Generate reports for auditing and compliance
   - Extract transaction data for auditing purposes
   - Generate financial reports to comply with regulatory requirements

VII. Monitor and optimize the payment facilitation process
   - Analyze transaction data to identify trends or bottlenecks
   - Optimize the payment facilitation process to improve efficiency and reduce costs

{% /comment %}

## Set up the appropriate accounts

In this section, we will guide you through setting up the necessary accounts for FBO clearing (payment facilitation) with Twisp. These accounts will help you manage transactions between parties, process payments, and disburse funds in a timely manner.

### 1. Create a settlement account for FBO clearing

To handle the settlement process for FBO clearing, you need to create a settlement account. This account will be used to settle and manage transactions between different parties in the payment facilitation process.

```graphql
mutation {
  createAccount(
    input: {
      accountId: "settlement-account-id"
      name: "FBO Settlement"
      code: "FBO_SETTLE"
      description: "FBO Settlement account"
      status: ACTIVE
      normalBalanceType: DEBIT
    }
  ) {
    accountId
  }
}
```

### 2. Create an account for payment processing fees

To manage the fees associated with payment processing, create a separate account. This account will be used to track the fees collected during the payment facilitation process.

```graphql
mutation {
  createAccount(
    input: {
      accountId: "processing-fees-account-id"
      name: "Payment Processing Fees"
      code: "PROCESSING_FEES"
      description: "Account for tracking payment processing fees"
      status: ACTIVE
      normalBalanceType: DEBIT
    }
  ) {
    accountId
  }
}
```

### 3. Create accounts for each party involved in the transactions

For each party involved in the transactions, such as merchants and customers, create separate accounts. These accounts will be used to track funds and manage transactions for each party.

For a merchant account:
```graphql
mutation {
  createAccount(
    input: {
      accountId: "merchant-account-id"
      name: "Merchant A"
      code: "MERCHANT_A"
      description: "Merchant A's account"
      status: ACTIVE
      normalBalanceType: CREDIT
    }
  ) {
    accountId
  }
}
```

For a customer account:
```graphql
mutation {
  createAccount(
    input: {
      accountId: "customer-account-id"
      name: "Customer A"
      code: "CUSTOMER_A"
      description: "Customer A's account"
      status: ACTIVE
      normalBalanceType: CREDIT
    }
  ) {
    accountId
  }
}
```

After setting up these accounts, you will have a proper foundation for handling FBO clearing with Twisp. In the next sections, we will guide you through defining transaction codes, modeling payment facilitation transactions, and managing the settlement process.

## Define transaction codes for FBO clearing

In order to handle FBO clearing with Twisp, you'll need to define transaction codes that represent the different types of transactions involved in the payment facilitation process. These transaction codes will help structure and categorize transactions in the ledger, ensuring consistency, predictability, and correctness.

### 1. Design a tran code for moving funds between parties

To support fund transfers between parties in the payment facilitation process, you'll need to create a tran code that represents this type of transaction. Let's call this tran code `FBO_TRANSFER`.

- Unique identifier code: `FBO_TRANSFER`
- Accounts involved: Sender's account (debit) and Receiver's account (credit)
- Entry data: Amount, Sender's account ID, and Receiver's account ID

### 2. Design a tran code for payment processing fees

To account for payment processing fees, create a separate tran code that will handle the deduction of fees from the transaction amount. Let's call this tran code `FBO_FEE`.

- Unique identifier code: `FBO_FEE`
- Accounts involved: Sender's account (debit) and Payment Processing Fee account (credit)
- Entry data: Fee amount, Sender's account ID, and Payment Processing Fee account ID

### 3. Design a tran code for fund disbursement

Finally, to support the disbursement of funds to the final destination, you'll need to create another tran code that represents this type of transaction. Let's call this tran code `FBO_DISBURSE`.

- Unique identifier code: `FBO_DISBURSE`
- Accounts involved: Receiver's account (debit) and Final Destination account (credit)
- Entry data: Disbursement amount, Receiver's account ID, and Final Destination account ID

By defining these transaction codes, you can ensure that the FBO clearing process is accurately represented in the Twisp ledger. Remember that transaction codes play a critical role in identifying and categorizing transactions in the ledger, and designing them correctly is an important part of using Twisp for FBO clearing.

## Modeling payment facilitation transactions with twisp

In this section, we will demonstrate how to model payment facilitation transactions using Twisp by posting transactions, deducting payment processing fees, and disbursing funds to their final destinations.

### 1. Post transactions using the defined tran codes to move funds between parties

To post a transaction between parties, you'll need to first choose a transaction code (tran code) that represents the type of transaction, such as `RECORD_TX`. As mentioned in the context, transaction codes act as templates for transactions and define the accounts, amounts, currencies, and other parameters required for the transaction.

For example, if you have designed a `RECORD_TX` tran code for moving funds between two parties, you can use the `postTransaction` mutation to create a new transaction using this tran code. You'll need to provide the required parameters, such as the debit and credit accounts, the amount, currency, and effective date of the transaction.

### 2. Deduct payment processing fees using the appropriate tran code

To deduct payment processing fees, you'll need to create another tran code, such as `FEE_DEDUCTION`. This tran code should be designed to debit the appropriate fee amount from the payer's account and credit the payment processing fee account.

When posting a transaction that involves payment processing fees, use the `FEE_DEDUCTION` tran code along with the required parameters, such as the payer's account, payment processing fee account, fee amount, currency, and effective date.

### 3. Disburse funds to the final destination using the designated tran code

Finally, to disburse funds to their final destinations (e.g., the payee's account), you'll need to use a tran code designed for fund disbursement, such as `FUNDS_DISBURSEMENT`. This tran code should debit the funds from the settlement account and credit the payee's account.

To post a transaction for fund disbursement, use the `postTransaction` mutation with the `FUNDS_DISBURSEMENT` tran code and provide the required parameters, such as the settlement account, payee's account, amount, currency, and effective date.

By following these steps, you can successfully model payment facilitation transactions using the Twisp accounting core. Remember to monitor the settlement process and reconcile transactions to ensure accurate and timely fund disbursements.

## Monitor the settlement process

In order to ensure the accuracy and timeliness of the FBO clearing process, it's crucial to monitor the settlement process using Twisp's features. This will help you track the status of fund transfers, review transactions in the ledger and identify any discrepancies or issues.

### 1. Query the settlement account

To track the status of fund transfers, you can use Twisp's GraphQL API to query the settlement account. You can retrieve information about the account balance, transactions, and other related data. Make sure to focus on the relevant settlement accounts, such as Card Settlement, Cash Settlement, and Checking Settlement, based on your specific use case.

Example query:

```graphql
query {
  account(id: "SETTLEMENT_ACCOUNT_ID") {
    id
    balance
    transactions {
      id
      tranCode
      amount
      datetime
    }
  }
}
```

### 2. Review transactions in the ledger

It's essential to review the transactions that have been posted to the ledger to ensure accurate processing and disbursement. You can use Twisp's API to retrieve transaction data from the ledger, including information about the transaction type, amount, and timing.

When reviewing transactions, pay special attention to their layer (`SETTLED`, `PENDING`, or `ENCUMBRANCE`) and sequencing, as these aspects play a crucial role in the accuracy and reliability of the ledger data.

Example query:

```graphql
query {
  transactions(accountId: "SETTLEMENT_ACCOUNT_ID") {
    id
    tranCode
    amount
    layer
    datetime
  }
}
```

By monitoring the settlement process through querying settlement accounts and reviewing transactions, you can ensure the proper functioning of your FBO clearing process and maintain a reliable ledger using Twisp's features. If any discrepancies are found, it's important to investigate and resolve the issue promptly to maintain the integrity of your accounting system.

## Reconcile and settle transactions

In this section, we'll cover the process of reconciling and settling transactions in the FBO clearing process using Twisp's accounting core. Reconciliation is essential for ensuring the accuracy and integrity of your ledger, while settlement ensures the completion of transactions and the appropriate movement of funds.

### 1. Perform reconciliation to match transactions with their corresponding disbursements

- Use Twisp's Timing and Sequencing features to accurately record the times and sequences of events in the ledger. This will be crucial during the reconciliation process to match each transaction with its corresponding disbursement.
- Leverage the Twisp accounting core's Layer feature to distinguish transactions between the Settled, Pending, and Encumbrance layers. This will enable you to manage transactions in different stages of the clearing process effectively.
- To perform reconciliation, query the ledger to obtain a list of transactions and their corresponding entries. Compare these entries with external records (such as bank statements or payment processor records) to identify any discrepancies or missing transactions.
- In case of discrepancies, Twisp's Tran Codes can help you identify the type of transaction and its impact on the ledger. Investigate any unmatched or incorrect transactions and make the necessary adjustments to the ledger.

### 2. Settle transactions by updating the ledger and moving funds to the appropriate accounts

- Once the reconciliation process is complete, it's time to settle transactions and move funds to the appropriate accounts.
- Use the Twisp accounting core's Three Layer Model to manage transactions in different stages of the settlement process. For example, move transactions from the Pending layer to the Settled layer once they are confirmed to be completed.
- Update the ledger by posting transactions using the appropriate Tran Codes. This will create debit and credit entries in the relevant accounts, ensuring the movement of funds is accurately recorded. For example, use the `RECORD_SETTLE_PENDING_TX` Tran Code to record the settlement of a pending transaction.
- Monitor the settlement process by regularly querying the ledger and checking the account balances. This will enable you to track the status of fund transfers and ensure the accuracy of the ledger.

By following these steps, you'll be able to effectively reconcile and settle transactions in the FBO clearing process using Twisp's accounting core. This will ensure the accuracy and integrity of your ledger while maintaining efficient funds management.

## Generate reports for auditing and compliance

Keeping accurate records and generating financial reports are essential for any business handling financial transactions. Twisp's accounting core provides a robust system that ensures compliance with regulatory requirements and facilitates effective auditing.

### 1. Extract transaction data for auditing purposes

Twisp's transaction ledger stores every financial event with a double-entry accounting system, which helps maintain accurate records. You can extract transaction data for auditing by querying the ledger using Twisp's GraphQL API. By filtering and sorting transactions based on specific criteria, such as date range, account, or transaction code, you can gather the necessary information for auditing purposes.

Here's an example query to fetch transactions for a specific account and date range:

```graphql
query {
  transactions(filter: { accountId: "exampleAccountId", date: { gte: "2022-01-01", lte: "2022-12-31" } }) {
    id
    tranCode
    amount
    date
    metadata
  }
}
```

### 2. Generate financial reports to comply with regulatory requirements

Using the extracted transaction data, you can create financial reports that comply with regulatory requirements. Twisp's accounting core follows industry standards and strong accounting principles, ensuring that generated reports are accurate and reliable.

Financial reports may include balance sheets, income statements, and cash flow statements. To generate these reports, aggregate the relevant transaction data by account, transaction type, and period. You can use the layered accounting feature in Twisp to clarify the true state of funds and distinguish between settled, pending, and planned financial activities.

For example, to generate a balance sheet report, query balances for all accounts in the chart of accounts:

```graphql
query {
  accounts {
    id
    name
    type
    balance(layer: SETTLED)
  }
}
```

Remember to tailor your queries and reports to the specific requirements of your business and regulatory environment. By leveraging Twisp's accounting core capabilities and adhering to best practices, you will be able to maintain compliance and support effective auditing processes.

## Monitor and optimize the payment facilitation process

To ensure the payment facilitation process is efficient and cost-effective, it is essential to keep track of transaction data, identify trends or bottlenecks, and make any necessary adjustments. Using Twisp's features, you can effectively monitor and optimize the payment facilitation process.

### 1. Analyze transaction data to identify trends or bottlenecks

- Utilize Twisp's transaction monitoring capabilities to review the transactions in your ledger. Keep track of the timing and sequencing of events to ensure accurate recording and reconciliation (source: [Timing and Sequencing](/accounting-core/encoded-transactions#timing-and-sequencing)).

- Identify any trends, such as frequently occurring transaction types or particularly large transactions, that might require further investigation or optimization. Transaction codes are an important aspect of categorizing transactions and can provide valuable insights (source: [Conclusion](/tutorials/building-tran-codes#conclusion)).

### 2. Optimize the payment facilitation process to improve efficiency and reduce costs

- Based on your analysis of transaction data, consider updating or redesigning transaction codes to streamline specific transaction types or minimize processing fees (source: [Designing a Tran Code](/reference/ledger/tran-codes#designing-a-tran-code)).

- If necessary, update your chart of accounts to better represent your funds flows or improve the organization of your accounting system (source: [Creating Accounts](/accounting-core/chart-of-accounts#creating-accounts)).

- Periodically review and test your accounting system's design to ensure it continues to satisfy your working requirements, particularly as your product grows in complexity (source: [Context](/tutorials/twisp-101#context)).

By closely monitoring transaction data and making data-driven optimizations, you can enhance the efficiency and cost-effectiveness of your payment facilitation process using Twisp's powerful features.

## Conclusion

Managing FBO clearing with Twisp's accounting tools, as detailed in this guide, is a structured and reliable approach to payment facilitation. The process encompasses setting up appropriate accounts, defining transaction codes specific to FBO clearing, and effectively modeling transactions for a seamless transfer of funds.

The role of diligent monitoring and reconciliation in the settlement process ensures accurate transaction handling and fund disbursement, underpinning the trust in your financial management. Equally important is the generation of detailed reports, which not only aids in meeting compliance standards and auditing requirements but also provides valuable insights into your financial operations.

Moreover, the ability to monitor and optimize the payment facilitation process is a significant advantage, offering an opportunity to continually improve efficiency, reduce costs, and maintain a robust, well-tuned operation. These insights can be instrumental in streamlining your payment facilitation processes, driving operational excellence, and ultimately bolstering your organization's financial health.

In essence, Twisp's tools allow for a rigorous and efficient approach to FBO clearing, contributing significantly to the smooth running of your financial systems. As fintech companies, banks, and other organizations continually seek improved ways to handle payments, Twisp's capabilities stand out as a comprehensive solution to meet these evolving demands. Whether you're looking to optimize your existing processes or establish a new payment facilitation system, this guide should provide a strong foundation for your journey.
