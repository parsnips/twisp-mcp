---
title: Doing Cash Reconciliation with Twisp
description: |
  How to conduct cash reconciliations using the Twisp system. This guide covers all necessary steps from preparing accounts and importing bank statement data, through matching transactions and handling discrepancies, to finalizing the reconciliation and generating reports for ongoing monitoring and insights. Suitable for anyone seeking to maintain accurate financial records and gain a comprehensive view of their organization's financial status.
aiCommand: |
  pnpm ai write "Write a guide for doing cash reconciliation with Twisp. Among other things, it should cover transaction matching, unmatched transaction handling, and reconciliation adjustments. Only use documented features of the Twisp system; do not invent features." -a GuideWriter -s cash-reconciliation -o ~/Desktop/guide-cash-reconciliation.md -m gpt-4
showNext: false
showPrev: false
draft: true
---

{% comment %}

OUTLINE

I. Preparing for Cash Reconciliation
   - Review your chart of accounts and ensure proper setup for bank accounts, assets, and any relevant reconciliation adjustments accounts.
   - Define transaction codes for bank statements, deposit slips, and withdrawal slips to categorize and identify transactions during reconciliation.

II. Importing Bank Statement Data
   - Obtain the bank statement data for the period you want to reconcile.
   - Format the data to be compatible with Twisp, including transaction codes, amounts, dates, and other necessary metadata.
   - Import the formatted bank statement data into Twisp using the Twisp API or other supported methods.

III. Matching Transactions
   - Query the ledger for transactions in the period you're reconciling.
   - Compare the imported bank statement data with the transactions in the ledger.
   - Match transactions in Twisp with their corresponding bank statement entries based on transaction codes, amounts, and dates.

IV. Handling Unmatched Transactions
   - Identify transactions in Twisp that do not have a corresponding entry in the bank statement data.
   - Review unmatched transactions for possible errors or discrepancies, such as missing or incorrect transaction codes, amounts, or dates.
   - Correct any errors or discrepancies in the unmatched transactions and re-run the matching process.

V. Reconciliation Adjustments
   - Identify any necessary reconciliation adjustments, such as bank fees, interest income, or other adjustments that should be recorded in the ledger.
   - Create transactions in Twisp for the reconciliation adjustments, using the appropriate transaction codes and accounts.
   - Post the reconciliation adjustment transactions to the ledger.

VI. Finalizing the Cash Reconciliation
   - Verify that all transactions have been matched, and all necessary reconciliation adjustments have been made.
   - Confirm that the account balances in Twisp match the ending balances on the bank statement.
   - Document the cash reconciliation process, including any issues encountered and resolved during the reconciliation.

VII. Monitoring and Reporting
   - Set up regular cash reconciliation processes to maintain accurate financial records.
   - Monitor the reconciliation process for potential issues or discrepancies on an ongoing basis.
   - Generate reconciliation reports to provide insight into the financial performance and status of your organization.

{% /comment %}


## Preparing for cash reconciliation

Before you can perform a cash reconciliation using Twisp, you need to ensure that your chart of accounts is set up correctly and that you've defined the necessary transaction codes for handling bank statement transactions. Here's how to prepare for cash reconciliation with Twisp:

### 1. Review your chart of accounts

To begin the cash reconciliation process, review your chart of accounts to ensure that it accurately reflects your business's financial activities. Make sure you have set up the required accounts for bank accounts, assets, and any relevant reconciliation adjustments.

In Twisp, you can create and manage your chart of accounts according to your financial needs. For more information on working with the chart of accounts, refer to [Chart of Accounts](/accounting-core/chart-of-accounts).

### 2. Define transaction codes

Transaction codes are crucial for categorizing and identifying transactions during the reconciliation process. In Twisp, you can create custom transaction codes that mirror the types of transactions your system handles, such as bank statements, deposit slips, and withdrawal slips. To define transaction codes, consider the following preliminary questions:

- What type of transaction is this? (e.g., ACH transfer, credit card purchase, etc.)
- Which accounts are involved? (e.g., where is the money coming from, and where is it going?)
- Are there any additional accounts that need to be debited/credited?
- How many entries should be written to the ledger, and which entries go on the debit side and which on the credit side?

For more information on designing and managing transaction codes, refer to [Building Tran Codes](/tutorials/building-tran-codes) and [Tran Codes](/tutorials/example-setup#tran-codes).

By ensuring that your chart of accounts is set up correctly and that you've defined the necessary transaction codes, you will be well-prepared to perform a cash reconciliation using Twisp.

## Importing bank statement data

In this section, we will guide you through the process of importing bank statement data into Twisp. This process involves obtaining the bank statement data, formatting it to be compatible with Twisp, and importing the formatted data using the Twisp API.

### 1. Obtain the bank statement data

Reach out to your bank or financial institution to get the bank statement for the period you want to reconcile in a digital format, such as CSV or Excel. Ensure the data contains all the necessary information, including transaction dates, amounts, descriptions, and any other relevant details.

### 2. Format the data to be compatible with Twisp

Before importing the data, you need to ensure it is correctly formatted according to Twisp's requirements. This includes assigning appropriate transaction codes to each transaction. Review the Twisp documentation, specifically the tutorials on [transaction codes](/tutorials/building-tran-codes) to understand how to format the data correctly. Make sure the data includes the following information:

- Transaction codes: Assign each transaction a code based on its type (e.g. deposits, withdrawals, or fees). These codes will be used to categorize and identify transactions during reconciliation.

- Amounts: Ensure all transaction amounts are in the same currency and format as your ledger in Twisp.

- Dates: Format transaction dates according to the Twisp system requirements (e.g., ISO 8601 format).

- Other necessary metadata: Include any other relevant metadata needed for reconciliation, such as account numbers, descriptions, or reference numbers.

### 3. Import the formatted bank statement data

With the bank statement data correctly formatted, you can now import the data into Twisp. Interact with the Twisp GraphQL API using the Twisp Console and the GraphiQL tool or your preferred GraphQL client. Use the `postTransaction` mutation to create transactions in your ledger for each record in the bank statement data. Remember to authenticate your requests using your Twisp API credentials.

By following these steps, you will successfully import your bank statement data into Twisp, making it ready for the next phase of cash reconciliation: matching transactions.

## Matching transactions

In order to perform cash reconciliation, it is essential to match transactions in Twisp with their corresponding bank statement entries. To do this, follow these steps:

### 1. Query the ledger for transactions in the period you're reconciling
Use Twisp's GraphQL API to retrieve all transactions within the date range you are reconciling. You can filter transactions by date, account, or transaction code, as needed. This will provide you with a list of transactions recorded in Twisp during the specified period.

### 2. Compare the imported bank statement data with the transactions in the ledger
With both the bank statement data and the transactions from Twisp available, you can now compare them side by side to identify matches. Look for similarities in transaction codes, amounts, and dates.

### 3. Match transactions in Twisp with their corresponding bank statement entries
Based on similarities in transaction codes, amounts, and dates, match Twisp transactions to their corresponding bank statement entries. You can use the correlation feature in Twisp to help you establish relationships between matched transactions.

Remember, Twisp transactions are built using double-entry accounting, so ensure that you are matching complete transaction pairs (i.e., a debit entry and a credit entry) with the corresponding bank statement entries.

By carefully matching transactions in this way, you can ensure the accuracy and consistency of your financial records in Twisp, ultimately enabling a smooth cash reconciliation process.

## Handling unmatched transactions

In the cash reconciliation process, it's crucial to identify and resolve any unmatched transactions between your Twisp ledger and the bank statement data. Unmatched transactions may result from errors, discrepancies, or missing information. Here's how to handle unmatched transactions using Twisp:

### 1. Identify unmatched transactions

Identify transactions in Twisp that do not have a corresponding entry in the bank statement data.

Query the transactions in your Twisp ledger for the reconciliation period using the appropriate filters (e.g., date range, accounts, transaction codes).

Compare the queried transactions with the imported bank statement data.

### 2. Review unmatched transactions

Review transactions for possible errors or discrepancies, such as missing or incorrect transaction codes, amounts, or dates.

Examine each unmatched transaction in the Twisp ledger to identify any inconsistencies with the bank statement data.

Check for missing or incorrect transaction codes. Ensure that the transaction codes used in Twisp match the ones in the bank statement data. If needed, review your transaction codes setup and make modifications as guided in the Twisp tutorial on building tran codes.

Verify the amounts and dates associated with each unmatched transaction. Compare them with the corresponding data in the bank statement to spot any discrepancies.

### 3. Correct errors

Correct any errors or discrepancies in the unmatched transactions and re-run the matching process.

If you've identified errors or discrepancies in the unmatched transactions, correct them by modifying the transaction data in Twisp. You may need to update transaction codes, amounts, or dates, as appropriate.

After making the necessary corrections, re-run the matching process between the Twisp ledger and the bank statement data. This will ensure that all transactions are properly matched and accounted for in the cash reconciliation process.

By carefully handling unmatched transactions, you'll maintain the accuracy and integrity of your financial records within the Twisp accounting system.

## Reconciliation adjustments

During the cash reconciliation process, there may be differences between the bank statement and your Twisp ledger that require reconciliation adjustments. These adjustments could be due to bank fees, interest income, or other adjustments that need to be recorded in the ledger to accurately reflect your financial position.

### 1. Identifying necessary reconciliation adjustments

Compare the bank statement with the transactions in your Twisp ledger to identify any discrepancies that may require reconciliation adjustments.

Common adjustments may include bank fees, interest income, or other items that were not initially recorded in the ledger.

### 2. Creating transactions for reconciliation adjustments in Twisp

For each reconciliation adjustment, you need to create a corresponding transaction in Twisp.

Determine the appropriate transaction code (tran code) for each adjustment. If you don't have a suitable tran code, you may need to create a new one based on your chart of accounts and the nature of the adjustment. Refer to the [tutorial](/tutorials/building-tran-codes) on creating, modifying, and deleting transaction codes for guidance.

Ensure that the transactions use the correct accounts, amounts, and directions (`CREDIT` or `DEBIT`) as per the transaction codes defined.

### 3. Posting the reconciliation adjustment transactions to the ledger

After creating the transactions for the reconciliation adjustments, you need to post them to the ledger using the `postTransaction` mutation.

Make sure to provide all the necessary parameters for the chosen transaction code, such as the accounts to debit and credit, the amount of the adjustment, the currency, and the effective date of the adjustment.

Twisp will then post the transactions to the ledger and update the account balances accordingly, ensuring the integrity and consistency of the ledger record.

By following these steps, you can accurately account for reconciliation adjustments and maintain an up-to-date and accurate financial record in your Twisp ledger.

## Finalizing the cash reconciliation

In this final step, we will ensure that all transactions have been matched and necessary reconciliation adjustments have been made. We will also confirm that the account balances in Twisp match the ending balances on the bank statement before documenting the entire cash reconciliation process.

### 1. Verify all transactions have been matched

After completing the matching process, handling unmatched transactions, and making reconciliation adjustments, double-check that all transactions in Twisp have been matched with their corresponding entries in the bank statement data. This can be done by running a query to retrieve all transactions for the reconciliation period and confirming that they are accurately matched.

### 2. Confirm account balances

To ensure the accuracy of your financial records, compare the account balances in Twisp to the ending balances on the bank statement. You can do this by using the built-in balance queries in Twisp (refer to the [Conclusion](/tutorials/pulling-balances#conclusion) section). If the balances match, it indicates that the reconciliation process has been successful. If there are discrepancies, review the reconciliation process to identify and correct any errors.

### 3. Document the cash reconciliation process

Maintaining clear and accurate documentation of the entire cash reconciliation process is crucial for financial reporting, auditing, and future reconciliation efforts. Carefully document all the steps taken during the reconciliation process, including any issues encountered and resolved, reconciliation adjustments made, and any unmatched transactions that were handled. This documentation will serve as a reference for future reconciliations and provide valuable insight into the financial state of your organization.

By following these steps, you can efficiently finalize the cash reconciliation process using Twisp and maintain accurate and up-to-date financial records for your organization.

## Monitoring and reporting

Setting up regular cash reconciliation processes and monitoring them is crucial to maintaining accurate financial records and ensuring the health of your organization's finances. Twisp provides the tools and features necessary to keep track of your financial performance and generate reconciliation reports. In this section, we'll discuss how to monitor and report using Twisp's documented features.

### 1. Set up regular cash reconciliation processes

Determine the frequency of cash reconciliation based on your organization's needs, such as monthly or quarterly.

Use Twisp's transaction and account features to plan and execute your reconciliation processes consistently within the determined timeframe.

### 2. Monitor the reconciliation process for potential issues or discrepancies

Regularly query the ledger for transactions, balances, and other relevant information using Twisp's GraphQL API.

Use Twisp's timing and sequencing features to ensure accurate recording of times and sequences for events in the ledger. This will aid in identifying any discrepancies or potential issues.

### 3. Generate reconciliation reports to provide insight into the financial performance and status of your organization

Use Twisp's Chart of Accounts to create balance sheets and P&L reports, giving you a comprehensive view of your organization's finances.

Utilize Twisp's account sets and hierarchical tree structures to aggregate and analyze balances across multiple accounts, providing a high-level overview of your organization's financial status.

In conclusion, Twisp's various features allow you to effectively monitor and report your cash reconciliation process, helping you maintain accurate financial records and ensuring the overall financial health of your organization. By staying on top of these tasks, you'll be better prepared to address any potential issues or discrepancies that may arise, and ultimately make better-informed financial decisions for your organization.

### Conclusion

Through this guide, you've learned how to perform an end-to-end cash reconciliation process using the Twisp system, offering a robust and systematic approach to maintain accurate financial records. The journey started from preparing for reconciliation by setting up the appropriate accounts and defining transaction codes, proceeded with the importation of bank statement data, matching transactions, and handling of unmatched transactions. It then guided you through making necessary reconciliation adjustments and finalizing the reconciliation.

In the end, by leveraging Twisp's capabilities, you're not only equipped to carry out reconciliations but also to monitor the process and generate insightful reports, further strengthening your organization's financial standing. Remember, regular cash reconciliation processes and constant monitoring can help prevent discrepancies and provide a comprehensive view of your organization's financial performance.

By adopting these practices, you've taken a significant step towards ensuring financial accuracy and facilitating the financial health of your organization. Your journey with Twisp does not end here—continue to explore and utilize its features to meet your evolving business needs.
