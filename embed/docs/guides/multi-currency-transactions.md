---
title: Supporting Multi-Currency Transactions with Twisp
description: |
  This guide provides instructions on how to effectively manage multi-currency transactions using the Twisp ledger system. It covers all essential steps, from setting up a multi-currency chart of accounts and defining transaction codes, through managing exchange rates and handling currency conversions, to tracking transactions, generating financial reports, and conducting regular audits. A must-read for organizations dealing with international transactions and seeking efficient multi-currency financial management solutions.
aiCommand: |
  pnpm ai write "Write a guide for managing multi-currency transactions with Twisp. The guide should cover tracking transactions in multiple currencies, handling currency conversions, and managing exchange rates for accurate financial reporting. Only use documented features of the Twisp system; do not invent features." -a GuideWriter -s multi-currency -o ~/Desktop/guide-multi-currency.md -m gpt-4
showNext: false
showPrev: false
draft: true
---

{% comment %}

OUTLINE

I. Set up multi-currency chart of accounts
   - Create separate account sets for each currency, e.g. USD, EUR, GBP
   - Designate a currency for each account in the account sets

II. Create transaction codes for multi-currency transactions
   - Define transaction codes for different currency transactions, e.g. USD_TRANSFER, EUR_TRANSFER, GBP_TRANSFER
   - Configure credit and debit entry types for each currency transaction code

III. Manage exchange rates
   - Create a dedicated exchange rate table or use an external API to obtain current exchange rates
   - Update exchange rates regularly to ensure accurate conversions

IV. Handle currency conversions in transactions
   - Use transaction codes for currency conversions, e.g. USD_TO_EUR, EUR_TO_GBP, GBP_TO_USD
   - Retrieve appropriate exchange rate from the exchange rate table or external API
   - Apply the exchange rate to the transaction amount during conversion

V. Track multi-currency transactions in separate journals
   - Create separate journals for each currency, e.g. USD_JOURNAL, EUR_JOURNAL, GBP_JOURNAL
   - Assign multi-currency transactions to the appropriate journals based on their currency

VI. Generate multi-currency financial reports
   - Query transactions and account balances in each currency
   - Convert account balances to a single reporting currency using exchange rates
   - Aggregate and present financial data in the chosen reporting currency for analysis

VII. Monitor and audit multi-currency transactions
   - Regularly review transaction history for inconsistencies or errors
   - Ensure compliance with accounting standards and regulations for multi-currency operations

{% /comment %}

## Set up multi currency chart of accounts

In order to manage multi-currency transactions with Twisp, you will first need to set up a multi-currency chart of accounts. This will involve creating separate account sets for each currency and designating a currency for each account in the account sets. Here's how to do it:

### 1. Create separate account sets for each currency

Start by creating an account set for each currency you want to support, such as USD, EUR, and GBP. You can create new account sets using the `createAccountSet` mutation. For instance, to create account sets for USD, EUR, and GBP accounts, you would run the following mutations:

```graphql
mutation {
  usd: createAccountSet(input: { code: "USD_SET", name: "USD Accounts" }) {
    code
    name
  }

  eur: createAccountSet(input: { code: "EUR_SET", name: "EUR Accounts" }) {
    code
    name
  }

  gbp: createAccountSet(input: { code: "GBP_SET", name: "GBP Accounts" }) {
    code
    name
  }
}
```

### 2. Designate a currency for each account in the account sets

Next, you need to specify the currency for each account within the respective account sets. While creating each account, you can add a custom field to store the currency information. For example, when creating a USD account under the USD_SET, you would include the currency field like this:

```graphql
mutation {
  createAccount(input: { code: "USD_ACCOUNT_1", name: "Example USD Account", custom: { currency: "USD" } }) {
    code
    name
    custom
  }
}
```

Once you've created the account, add it to the corresponding account set using the `addToAccountSet` mutation:

```graphql
mutation {
  addToAccountSet(input: { setCode: "USD_SET", accountCode: "USD_ACCOUNT_1" }) {
    code
    name
    accounts {
      code
      name
      custom
    }
  }
}
```

Repeat these steps for each account you create for EUR and GBP currencies, making sure to designate the appropriate currency for each account and add them to the corresponding account sets.

By organizing your accounts into separate sets based on currency, you can easily manage and track transactions in multiple currencies, as well as retrieve account balances and generate financial reports for each currency.

## Create transaction codes for multi currency transactions

To effectively manage multi-currency transactions in Twisp, you'll need to define separate transaction codes for each currency. This ensures that transactions involving different currencies are systematically processed and categorized in your ledger.

Here's how to create transaction codes for multi-currency transactions:

### 1. Define transaction codes for different currency transactions

For each currency you plan to support, create a unique transaction code. For example, if you're working with USD, EUR, and GBP, you might define the following transaction codes: `USD_TRANSFER`, `EUR_TRANSFER`, and `GBP_TRANSFER`.

### 2. Configure credit and debit entry types for each currency transaction code

For each transaction code, you'll need to set up the appropriate credit and debit entry types for the accounts involved in the transaction. These entry types will depend on the specific accounts and transaction types you're working with. For example, you might set up entry types like `USD_TRANSFER_DR`, `USD_TRANSFER_CR`, `EUR_TRANSFER_DR`, `EUR_TRANSFER_CR`, `GBP_TRANSFER_DR`, and `GBP_TRANSFER_CR`.

By defining separate transaction codes for each currency and configuring the corresponding credit and debit entry types, you'll ensure that multi-currency transactions are correctly processed and categorized in your Twisp ledger. This will help maintain the integrity of your financial data and make it easier to generate accurate financial reports.

## Manage exchange rates

To ensure accurate multi-currency transactions in Twisp, it is essential to manage exchange rates effectively. Here's a step-by-step process to help you manage exchange rates:

### 1. Create a dedicated exchange rate table or use an external API
While Twisp does not provide built-in support for exchange rates, you can create a dedicated exchange rate table in your application's database or use an external API to obtain current exchange rates. This will provide a reliable source of exchange rate information for your multi-currency transactions.

### 2. Update exchange rates regularly
To maintain accurate conversions, it is crucial to update the exchange rates in your table or retrieve up-to-date rates from the external API regularly. The frequency of updates will depend on the level of volatility in the foreign exchange markets and your specific business requirements.

### 3. Integrate exchange rates with Twisp transactions
When processing multi-currency transactions, you'll need to apply the appropriate exchange rate to the transaction amount during conversion. This can be done through your application logic, which should retrieve the relevant exchange rate from your table or external API, apply it to the transaction amount, and then use Twisp's Money type to represent the converted amount in the desired currency.

By following these steps for managing exchange rates, you can ensure that your application's multi-currency transactions are processed accurately within the Twisp accounting system. Remember to monitor the exchange rates regularly and update them as needed, to maintain the integrity of your financial data.

## Handle currency conversions in transactions

Managing multi-currency transactions with Twisp involves handling currency conversions for various transaction types. To perform currency conversions, follow these steps:

### 1. Create transaction codes for currency conversions

Design transaction codes for each currency conversion type you plan to support (e.g., `USD_TO_EUR`, `EUR_TO_GBP`, `GBP_TO_USD`). These transaction codes will specify which accounts to credit and debit, as well as the currency and other necessary inputs for the conversion process.

Example:

```
USD_TO_EUR: {
  debit_entry_type: 'USD_DEBIT',
  credit_entry_type: 'EUR_CREDIT',
  base_currency: 'USD',
  target_currency: 'EUR',
  ...
}
```

### 2. Retrieve the appropriate exchange rate

To ensure accurate conversions, obtain the current exchange rate for the currencies involved in the transaction. You can maintain an exchange rate table within your system or utilize an external API to fetch real-time exchange rates.

### 3. Apply the exchange rate to the transaction amount
Once you have the exchange rate, apply it to the transaction amount during the conversion process. For example, if you're converting a $100 USD_TO_EUR transaction with an exchange rate of 0.85, the converted amount would be €85.

Example:

```
transaction_amount = 100
exchange_rate = 0.85
converted_amount = transaction_amount * exchange_rate
```

### 4. Post the converted transaction

Using the `postTransaction` mutation and the chosen transaction code (e.g., USD_TO_EUR), post the transaction with the converted amount and the corresponding accounts to be debited and credited.

Example:

```graphql
mutation {
  postTransaction(input: {
    tranCode: "USD_TO_EUR",
    debitAccountId: "<USD account ID>",
    creditAccountId: "<EUR account ID>",
    amount: { units: 100, currency: "USD" },
    effectiveDate: "<date>"
  }) {
    id
    entries {
      id
      account {
        id
      }
      entryAmount {
        units
        currency
      }
    }
  }
}
```

By following these steps, you can handle currency conversions in transactions effectively, ensuring accurate accounting and financial reporting for your multi-currency operations with Twisp.

## Track multi currency transactions in separate journals

Managing multi-currency transactions with Twisp is made easier by organizing them in separate journals for each currency. This not only helps keep your ledger organized but also allows for easier reporting and analysis of transactions in each currency.

### 1. Create separate journals for each currency

To begin, you'll need to create a separate journal for each currency you're working with. For example, create journals named USD_JOURNAL, EUR_JOURNAL, and GBP_JOURNAL for transactions in US dollars, euros, and British pounds, respectively. This can be done using the Twisp GraphQL API or the web interface.

### 2. Assign multi-currency transactions to the appropriate journals

When posting multi-currency transactions, you'll need to assign each transaction to the appropriate journal based on the currency involved. With Twisp, you can use the `journalCode` field to specify the journal to which a transaction should be posted. For example, if you have a transaction in euros, you would assign it to the EUR_JOURNAL by setting the `journalCode` field to "EUR_JOURNAL" when posting the transaction.

By assigning multi-currency transactions to separate journals, you can effectively track and manage your financial data in different currencies. This not only helps maintain a clean and organized ledger but also allows for more accurate financial reporting and analysis when dealing with multiple currencies.

## Generate multi currency financial reports

In order to generate comprehensive financial reports for a multi-currency system in Twisp, you need to query transactions and account balances in each currency, convert account balances to a single reporting currency using exchange rates, and aggregate and present financial data in the chosen reporting currency for analysis.

### 1. Query transactions and account balances in each currency:

Use the `balance` query to retrieve the balance for each account in the desired currency. Make sure to provide the `journalId` and `currency` as required arguments. For example, if you want to retrieve the account balance in USD, you would use the following GraphQL query:

```graphql
query {
  balance(
    accountId: "3ea12e45-7df2-4293-9434-feb792affc91"
    journalId: "4e9d9f6c-0cc3-4a8e-9d1a-0de5a47b0c8b"
    currency: "USD"
  ) {
    settled {
      normalBalance {
        units
      }
    }
  }
}
```

Repeat this process for each currency you need to include in your financial reports.

### 2. Convert account balances to a single reporting currency using exchange rates

Using your exchange rate table or an external API, retrieve the appropriate exchange rate for each currency pair. Apply the exchange rate to each account balance to convert the balances to your chosen reporting currency. For example, if you want to convert an account balance from EUR to USD, you would multiply the EUR balance by the EUR to USD exchange rate.

### 3. Aggregate and present financial data in the chosen reporting currency for analysis

After converting all account balances to your chosen reporting currency, aggregate the data to generate financial reports such as balance sheets, income statements, or cash flow statements. The converted and aggregated data will provide a unified view of your financial operations across all currencies, making it easier to analyze and present your financial performance.

Remember to monitor and update exchange rates regularly to ensure the accuracy of your financial reports, as currency exchange rates may fluctuate over time. Regularly reviewing your multi-currency financial data will help you maintain compliance with accounting standards and regulations for multi-currency operations.

## Monitor and audit multi currency transactions

Monitoring and auditing multi-currency transactions are essential for maintaining the integrity and accuracy of your financial records. By regularly reviewing transaction history and ensuring compliance with accounting standards, you can prevent inconsistencies or errors in your multi-currency operations. Here's how to monitor and audit multi-currency transactions in Twisp:

### 1. Review transaction history
Regularly review the transaction history for each currency journal (e.g. USD_JOURNAL, EUR_JOURNAL, GBP_JOURNAL). This can be done by querying the transactions within the journal using Twisp's GraphQL API. Be sure to check for any unusual activity, inconsistencies, or errors that may have occurred during currency conversions or transfers.

### 2. Analyze ledger entries
Inspect the ledger entries associated with each transaction to verify that double-entry accounting principles have been followed and that no value has been lost or created during the transaction. Twisp enforces double-entry accounting, ensuring that each transaction involves at least two ledger entries with corresponding debits and credits.

### 3. Validate transaction codes
Ensure that the appropriate transaction codes have been used for each multi-currency transaction and that they accurately represent the intended financial activity. Review the tran codes for currency conversions and transfers to confirm that they have been designed and implemented correctly.

### 4. Check timing and sequencing
Verify that the timing and sequencing of multi-currency transactions are accurate and consistent. Twisp's ledger maintains a clear record of the times and sequences of events, which is essential for auditing and reconciliation purposes.

### 5. Ensure compliance with accounting standards and regulations
Regularly review your multi-currency accounting practices to ensure they adhere to relevant accounting standards and regulations. This may include verifying that the design of your chart of accounts, transaction codes, and journals are in line with industry best practices and regulatory requirements.

### 6. Perform periodic audits
Consider conducting periodic internal or external audits of your multi-currency accounting system to further ensure its accuracy and compliance. Audits can help identify potential issues or areas for improvement in your multi-currency operations.

By following these steps, you can effectively monitor and audit your multi-currency transactions, ensuring the accuracy and integrity of your financial records. Remember, Twisp provides a robust accounting core designed with strong, time-tested accounting principles to help you maintain a reliable and consistent financial system.

## Conclusion

Congratulations on completing this guide on managing multi-currency transactions with Twisp. You've learned how to create a multi-currency chart of accounts, define transaction codes for different currencies, manage and apply exchange rates, and handle currency conversions. You've also understood the importance of maintaining separate journals for different currencies and how to generate comprehensive multi-currency financial reports. Lastly, this guide highlighted the importance of regular monitoring and auditing of transactions to ensure compliance with accounting standards and identify inconsistencies early.

Remember that efficient management of multi-currency transactions plays a crucial role in international business and financial operations. By following the steps outlined in this guide, you can effectively track, manage, and report your organization's multi-currency transactions using Twisp, thereby ensuring accuracy and compliance. It's a continuous learning process, so feel free to refer back to this guide and revise the steps as needed. Happy managing!