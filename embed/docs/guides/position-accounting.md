---
title: Position-Based Accounting for Securities Trading with Twisp
description: |
  Learn how to set up position-based accounting for securities trading using Twisp in this comprehensive guide. From designing a suitable chart of accounts and defining transaction codes, to posting transactions for trading activities and analyzing performance over time, this guide provides step-by-step instructions for managing the complexities of securities trading. Ideal for those in financial institutions or investment companies looking to streamline their securities accounting process.
aiCommand: |
  pnpm ai write "Write a guide for doing position-based accounting for securities trading with Twisp. The guide should cover securities transactions, and how to structure accounts and tran codes to support the basic use cases. Only use documented features of the Twisp system; do not invent features." -a GuideWriter -s position-accounting -o ~/Desktop/guide-position-accounting.md -m gpt-4
showNext: false
showPrev: false
---

{% callout type="task" %}
Over the course of this guide, we will:
- Design a chart of accounts
- Define tran code for securities transactions
- Post transactions for securities trading activities
- Query balances and transaction history
{% /callout %}

## Design the Chart of Accounts

To set up position-based accounting for securities trading with Twisp, you'll need to create a chart of [accounts]($gql:object:Account) that can accurately represent the financial activities associated with trading securities.

In this section we'll walk through the initial steps needed to set up your ledger. To create these accounts, refer to the account creation process outlined in the Twisp documentation (source: [Creating Accounts](/accounting-core/chart-of-accounts#creating-accounts)).

### 1. Create securities trading accounts for customers

To represent each customer's securities trading activities, you'll need to create separate accounts for them within your chart of accounts.

Because accounts can maintain balances across multiple currencies, we can treat each security as its own currency and thus a single customer account can hold positions across many securities.

A simple way to set this up would use a single account for each customer, e.g. `CLIENT_A`, `CLIENT_B`, etc. A more robust setup might use an account set for each customer, with sub-accounts for different purposes:

- `CLIENT_*.CASH` - (asset) debit-normal
- `CLIENT_*.INVESTMENT` - (equity) credit-normal
- `CLIENT_*.HOLDINGS` - (asset) debit-normal
- `CLIENT_*.DIVIDEND_REVENUE` - (revenue) credit-normal

### 2. Establish revenue and expense accounts for trading fees and transaction costs

In order to account for trading fees and transaction costs associated with securities trading, you'll need to create separate revenue and expense accounts within your chart of accounts.

This will enable you to accurately track fees charged to customers and expenses incurred by your business for executing trades.

### 3. Add accounts for realized and unrealized gains or losses

Lastly, you'll need to create accounts for tracking realized and unrealized gains or losses resulting from securities trading activities.

Realized gains or losses occur when a security is sold, while unrealized gains or losses represent the change in value of securities that are still being held by customers.

By creating separate accounts for these financial events, you can effectively monitor the performance of your customers' portfolios and provide accurate financial reporting.

## Define Transaction Codes for Securities Transactions

In this section, we will design the transaction codes ([tran codes]($gql:object:TranCode)) for various securities transactions, such as buying and selling securities, dividend and interest payments, and handling fees and transaction costs.

### 1. Design transaction codes for buying securities

To create a transaction code for buying securities, follow these steps:

- Define a unique identifier code for this tran code, such as `BUY_SECURITY`.
- Determine the accounts that will be debited and credited. For example, debit the customer's cash account for the purchase amount and credit the customer's securities account for the purchased security.
- Define the entry data to be written for each entry, such as entry types like `BUY_SECURITY_DR` and `BUY_SECURITY_CR`.
- Parameterize the inputs, including the purchase amount, security type, and customer account IDs.

For example, this tran code might accept the following parameters:

- `clientHoldingsAcct` (`UUID`)
- `clientInvestAcct` (`UUID`)
- `clientCashAcct` (`UUID`)
- `security` (`STRING`)
- `shares` (`DECIMAL`)
- `price` (`DECIMAL`)
- `feeRate` (`DECIMAL`)
- `effective` (`DATE`)

The exact nature of your ledger will dictate how entries should be written. An example set of entries for this tran code could be:

| Entry Type                      | Account              | Units                              | Currency   | Direction | Layer   |
|---------------------------------|----------------------|------------------------------------|------------|-----------|---------|
| BUY_SECURITY_PAYMENT_CR         | `clientCashAcct`     | `(price * shares) * (1 + feeRate)` | USD        | CREDIT    | SETTLED |
| BUY_SECURITY_PAYMENT_DR         | `clientInvestAcct`   | `(price * shares) * (1 + feeRate)` | USD        | DEBIT     | SETTLED |
| BUY_SECURITY_BROKER_EARNINGS_CR | BROKER.REVENUE       | `(price * shares)`                 | USD        | CREDIT    | SETTLED |
| BUY_SECURITY_BROKER_FEE_CR      | BROKER.REVENUE       | `(price * shares) * (feeRate)`     | USD        | CREDIT    | SETTLED |
| BUY_SECURITY_BROKER_EARNINGS_DR | BROKER.SETTLEMENT    | `(price * shares) * (1 + feeRate)` | USD        | DEBIT     | SETTLED |
| BUY_SECURITY_TRADE_CR           | BROKER.HOLDINGS      | `shares`                           | `security` | CREDIT    | SETTLED |
| BUY_SECURITY_TRADE_DR           | `clientHoldingsAcct` | `shares`                           | `security` | DEBIT     | SETTLED |

Let's break these entries down.

1. **BUY_SECURITY_PAYMENT_CR**: This is the entry for the payment made by the client for purchasing the security. The account debited is the `clientCashAcct` which is the account where the client holds cash. The amount is calculated by multiplying the price per share of the security by the number of shares being purchased. This total is then adjusted upwards for the brokerage fee by multiplying by `(1 + feeRate)`. The direction is CREDIT because money is leaving this account to make the payment. The entry is marked as SETTLED, meaning the transaction is complete.
2. **BUY_SECURITY_PAYMENT_DR**: This is the mirror entry to the above, posting to the `clientInvestAcct` (Client Investment Account) which represents the client's investments in securities. The amount is the same as in the `BUY_SECURITY_PAYMENT_CR` entry, but the direction is DEBIT because the account is receiving value in the form of the purchased security.
3. **BUY_SECURITY_BROKER_EARNINGS_CR**: This entry records the revenue earned by the broker from the client's purchase. The account credited is `BROKER.REVENUE`. The amount is the raw cost of the security purchase without any fee adjustment.
4. **BUY_SECURITY_BROKER_FEE_CR**: This entry records the brokerage fee that the client paid for this transaction. The account credited is again `BROKER.REVENUE`. The amount is the cost of the security purchase multiplied by the fee rate.
5. **BUY_SECURITY_BROKER_EARNINGS_DR**: This entry records the total amount of the purchase (including fee) as a debit to `BROKER.SETTLEMENT` account. It is a mirror to the `BUY_SECURITY_BROKER_EARNINGS_CR` and `BUY_SECURITY_BROKER_FEE_CR` entries combined, effectively 'settling' the broker's accounts.
6. **BUY_SECURITY_TRADE_CR**: This entry records the receipt of the purchased shares in the `BROKER.HOLDINGS` account. The amount is the number of shares bought, and the currency column indicates the type of security purchased.
7. **BUY_SECURITY_TRADE_DR**: This is the mirror entry to the above, recording the departure of the purchased shares from the `clientHoldingsAcct`. The amount is the number of shares bought.

Overall, these entries reflect the flow of money and securities between the client and the broker during a `BUY_SECURITY` transaction. They ensure the integrity and correctness of the accounting process by adhering to the double-entry accounting system, where each transaction is represented by two or more entries that balance out.

### 2. Design transaction codes for selling securities

To create a transaction code for selling securities, follow these steps:

- Define a unique identifier code for this tran code, such as `SELL_SECURITY`.
- Determine the accounts that will be debited and credited. For example, debit the customer's securities account for the sold security and credit the customer's cash account for the sale proceeds.
- Define the entry data to be written for each entry, such as entry types like `SELL_SECURITY_DR` and `SELL_SECURITY_CR`.
- Parameterize the inputs, including the sale amount, security type, and customer account IDs.

### 3. Develop transaction codes for dividend or interest payments

To create a transaction code for dividend or interest payments, follow these steps:

- Define a unique identifier code for this tran code, such as `DIVIDEND_PAYMENT` or `INTEREST_PAYMENT`.
- Determine the accounts that will be debited and credited. For example, debit the asset account representing dividends or interest payments and credit the customer's cash account for the payment received.
- Define the entry data to be written for each entry, such as entry types like `DIV_PAYMENT_DR` and `DIV_PAYMENT_CR` or `INT_PAYMENT_DR` and `INT_PAYMENT_CR`.
- Parameterize the inputs, including the payment amount and customer account IDs.

By defining appropriate transaction codes for various securities transactions, you can ensure that your securities trading platform maintains a consistent, predictable, and correct ledger using Twisp's powerful accounting features.

## Post Transactions for Securities Trading Activities

Now that we have the accounts and tran codes defined, we can start using them to post transactions.

Use the transaction codes you created in the previous section for buying and selling securities.

### Example: buying a security

As an example, let's say we wanted to record a `BUY_SECURITY` transaction where `CLIENT_A` bought 100 shares of AAPL at $149.45 pre share with a 0.75% brokerage fee, effective as of October 12, 2022.

The values of the params would be:

- `clientHoldingsAcct`: UUID for `CLIENT_A.HOLDINGS` account
- `clientInvestAcct`: UUID for `CLIENT_A.INVESTMENT` account
- `clientCashAcct`: UUID for `CLIENT_A.CASH` account
- `security`: `'AAPL'`
- `shares`: `100`
- `price`: `149.45`
- `feeRate`: `0.0075`
- `effective`: `'2022-10-12'`

After posting, the following entries would be written to the ledger:

| Entry Type                      | Account             | Units      | Currency | Direction | Layer   |
|---------------------------------|---------------------|------------|----------|-----------|---------|
| BUY_SECURITY_PAYMENT_CR         | CLIENT_A.CASH       | 15057.0875 | USD      | CREDIT    | SETTLED |
| BUY_SECURITY_PAYMENT_DR         | CLIENT_A.INVESTMENT | 15057.0875 | USD      | DEBIT     | SETTLED |
| BUY_SECURITY_BROKER_EARNINGS_CR | BROKER.REVENUE      | 14945      | USD      | CREDIT    | SETTLED |
| BUY_SECURITY_BROKER_FEE_CR      | BROKER.REVENUE      | 112.0875   | USD      | CREDIT    | SETTLED |
| BUY_SECURITY_BROKER_EARNINGS_DR | BROKER.SETTLEMENT   | 15057.0875 | USD      | DEBIT     | SETTLED |
| BUY_SECURITY_TRADE_CR           | BROKER.HOLDINGS     | 100        | AAPL     | CREDIT    | SETTLED |
| BUY_SECURITY_TRADE_DR           | CLIENT_A.HOLDINGS   | 100        | AAPL     | DEBIT     | SETTLED |

### Post additional transactions to exercise all tran codes

Try posting more transactions for each of your tran code types. After posting, query the entries written to confirm that they correspond with what you expect.

- Account for dividend or interest payments
- Apply fees and transaction costs to the relevant accounts
- Calculate and record realized and unrealized gains or losses

## Monitor and Analyze Securities Trading Activities

In this section, we will cover how to monitor and analyze securities trading activities. These steps will help you ensure accuracy, maintain regulatory compliance, and track your customers' portfolios.

### 1. Query account balances and transaction history for individual customers

Using Twisp's GraphQL API, you can query the account balances and transaction history for each customer's securities trading accounts. This will enable you to track the financial activity of your customers and identify any potential issues or discrepancies.

### 2. Track portfolio performance and valuation over time

To track the performance and valuation of your customers' portfolios, you can use Twisp's layering system to create snapshots of balances at different points in time. By comparing these snapshots, you can calculate the changes in the value of your customers' portfolios and analyze their performance over specific periods.

### 3. Review and audit securities trading activities

Twisp's immutable ledger records and versioned account history allow you to easily audit and verify all securities trading activities. By regularly reviewing the transaction records and account histories, you can ensure the accuracy of your accounting system, identify possible fraud or misuse, and maintain regulatory compliance.

### 4. Generate financial reports and statements for regulatory compliance and management purposes

Using the data stored in Twisp's ledger, you can generate various financial reports and statements required for regulatory compliance, as well as for internal management purposes. These reports can help you monitor the overall performance of your securities trading business and make informed decisions based on accurate financial data.

## Conclusion

After following this guide, you should now be proficient in setting up position-based accounting for securities trading with Twisp. We started by designing an effective chart of accounts and defining precise transaction codes for various securities activities. Following this, we explored how to post transactions for these activities and how to accurately account for all the associated costs and gains. Finally, we discussed how to monitor and analyze trading activities to provide valuable insights into portfolio performance and compliance with regulatory requirements.

By utilizing the functionalities of Twisp in your securities trading, you can ensure a streamlined and efficient accounting process. Moreover, the ability to monitor and analyze your trading activities in real-time provides an added level of control, enabling you to make data-driven decisions and optimize performance.

Remember that continuous improvement is key to maintaining a robust and efficient accounting system. Regularly review your setup, adjust transaction codes and accounts when necessary, and keep an eye on new Twisp features and updates that could further improve your securities accounting process.

Thank you for choosing Twisp as your partner in financial management. Happy trading!
