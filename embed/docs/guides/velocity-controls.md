---
title: Implementing Velocity Controls with Twisp
description: |
  This guide covers everything from defining the scope and establishing rules for velocity controls, creating custom transaction codes (TranCodes), configuring alerts and monitoring systems, to refining controls over time. Ideal for organizations looking to enhance their risk management and fraud prevention mechanisms.
aiCommand: |
  pnpm ai write "Write a guide for implementing velocity controls with Twisp. The guide should cover setting up rules and limits for transactions across various dimensions such as merchants, accounts, amounts, and users, in order to mitigate fraud and manage risk. Only use documented features of the Twisp system; do not invent features." -a GuideWriter -s velocity-controls -o ~/Desktop/guide-velocity-controls.md -m gpt-4
showNext: false
showPrev: false
draft: true
---

{% comment %}

OUTLINE

I. Define the scope of velocity controls
    - Identify the dimensions for velocity controls (merchants, accounts, amounts, users)
    - Determine the types of transactions to be monitored
    - Assess the risk levels associated with each dimension

II. Create rules and limits for each dimension
    - Establish transaction limits (daily, weekly, monthly) per dimension
    - Set up alert thresholds for unusual transaction patterns
    - Define rules for temporary suspension or flagged accounts

III. Design and implement custom TranCodes for velocity controls
    - Create TranCodes for each type of transaction that needs monitoring
    - Incorporate velocity control rules and limits into TranCodes
    - Test and validate the behavior of the new TranCodes

IV. Configure monitoring and alerts based on velocity control rules
    - Set up notifications for transactions exceeding set limits or thresholds
    - Integrate alerts with existing monitoring and reporting systems
    - Establish protocols for addressing flagged transactions or accounts

V. Monitor and refine velocity controls over time
    - Regularly review transaction logs and alerts for accuracy and effectiveness
    - Adjust rules and limits based on observed trends and changing risk levels
    - Continuously improve velocity control measures to mitigate fraud and manage risk

{% /comment %}

## Define the scope of velocity controls

To effectively manage risk and mitigate fraud in your financial system, it is essential to define the scope of velocity controls. This process involves identifying the dimensions for velocity controls, determining the types of transactions to be monitored, and assessing the risk levels associated with each dimension.

### 1. Identify the dimensions for velocity controls

Velocity controls can be implemented across various dimensions, such as:

- Merchants: Monitor transaction activity related to specific merchants to identify unusual patterns or high-risk merchants.
- Accounts: Track transaction activity for individual accounts to detect unusual patterns or signs of fraud.
- Amounts: Set limits on transaction amounts to prevent large and potentially fraudulent transactions.
- Users: Monitor user behavior and transaction patterns to identify potential cases of fraud or account takeover.

### 2. Determine the types of transactions to be monitored

It's important to focus on the types of transactions most relevant to your business and industry. Common transactions to monitor include:

- ACH transfers
- Credit card purchases
- Wire transfers
- Foreign exchange transactions
- Loan disbursements and repayments

### 3. Assess the risk levels associated with each dimension

Different dimensions may carry varying levels of risk. Assess the risk levels associated with each dimension by considering factors such as transaction volume, historical patterns, and industry trends. This will allow you to allocate resources and implement appropriate velocity controls that effectively address the unique risks associated with each dimension.

By defining the scope of velocity controls, you can tailor your risk management strategies to your specific business needs and ensure the safety and security of your financial system.

## Create rules and limits for each dimension

In this section, we will discuss how to establish transaction limits, set up alert thresholds, and define rules for temporary suspension or flagged accounts for each dimension (merchants, accounts, amounts, users) to implement effective velocity controls with Twisp.

### 1. Establish transaction limits

- Analyze historical transaction data and industry benchmarks to determine appropriate limits for each dimension.
- For each dimension, set daily, weekly, and monthly transaction limits based on your risk analysis.
- Configure these limits in Twisp's accounting system by creating custom rules and parameters within TranCodes or account configurations.

### 2. Set up alert thresholds for unusual transaction patterns
- Identify patterns and behaviors that could indicate potential fraud or unusual activity, such as a sudden spike in transaction volume or a significant deviation from a user's typical transaction pattern.
- Define thresholds for these patterns that, when exceeded, will trigger an alert.
- Integrate these thresholds into your Twisp system by incorporating them into TranCodes, custom balance configurations, or monitoring tools.

### 3. Define rules for temporary suspension or flagged accounts
- Establish criteria for flagging or temporarily suspending accounts based on their velocity control violations, such as exceeding transaction limits or exhibiting unusual patterns.
- Create a process for reviewing flagged accounts and deciding whether to reinstate access, continue the suspension, or escalate the issue for further investigation.
- Implement these rules and processes within your Twisp system by updating account configurations, TranCodes, and monitoring tools.

By following these steps, you can create a robust set of rules and limits for each dimension in your Twisp accounting system, ensuring that you can effectively monitor and manage risk while maintaining a secure and efficient transaction environment.

## Design and implement custom tran codes for velocity controls

To effectively manage velocity controls, it is essential to create custom TranCodes tailored to each transaction type that requires monitoring. This section will guide you through creating TranCodes, incorporating velocity control rules and limits, and testing their behavior.

### 1. Create TranCodes for each type of transaction that needs monitoring
- Identify the transaction types that require velocity controls, such as money transfers, deposits, withdrawals, or payments.
- For each transaction type, create a unique TranCode that reflects its specific characteristics and required monitoring. Use descriptive names for your TranCodes, such as `TRANSFER_VELOCITY_CONTROL` or `DEPOSIT_VELOCITY_CONTROL`.

### 2. Incorporate velocity control rules and limits into TranCodes
- For each custom TranCode, define the rules and limits associated with the velocity control dimensions you identified earlier (e.g., merchants, accounts, amounts, users).
- Ensure that these rules and limits are enforced during the transaction process by incorporating them into the TranCode's design. This may involve adding specific conditions or checks within the TranCode to comply with the established rules and limits.

### 3. Test and validate the behavior of the new TranCodes
- Once your custom TranCodes have been designed with the appropriate velocity control rules and limits, it's crucial to test their behavior to ensure they function correctly.
- To test the TranCodes, create a series of sample transactions that cover various scenarios, such as transactions within limits, transactions exceeding limits, and transactions involving different dimensions (e.g., merchants, accounts, amounts, users).
- Validate that the TranCodes enforce the velocity control rules and limits as expected, and make any necessary adjustments to fine-tune their behavior.

By following these steps, you will have custom TranCodes that effectively manage velocity controls for your financial system, helping to mitigate fraud and manage risk. Remember to continuously monitor and adjust these TranCodes as needed to keep up with changing risks and requirements.

## Configure monitoring and alerts based on velocity control rules

In order to maintain a secure and efficient financial system, it is essential to configure monitoring and alerts based on the velocity control rules you have established. This will help you quickly identify and address any transactions that exceed your set limits or show unusual patterns. Follow these steps to set up monitoring and alerts for your Twisp-based accounting system:

### 1. Set up notifications for transactions exceeding set limits or thresholds
- Use Twisp's built-in monitoring and alerting capabilities to set up notifications for transactions that exceed the limits and thresholds defined in your velocity control rules.
- Configure the notifications to be sent to relevant stakeholders, such as financial managers, risk analysts, or your security team.

### 2. Integrate alerts with existing monitoring and reporting systems
- If your organization uses other monitoring and reporting tools, integrate Twisp's alerts with these systems to ensure a seamless and comprehensive view of your financial activities.
- This integration will allow you to consolidate alerts from different sources and maintain a unified monitoring and alerting environment.

### 3. Establish protocols for addressing flagged transactions or accounts
- To effectively manage flagged transactions or accounts, create a protocol that outlines the steps your team should take when an alert is triggered.
- This protocol may include steps like investigating the transaction, contacting the relevant parties, temporarily suspending the account, or escalating the issue to higher authorities.
- Ensure that your team is trained on these protocols and understands the importance of promptly addressing flagged transactions or accounts to mitigate fraud and manage risk.

By following these steps, you will be able to effectively monitor and respond to any unusual transaction patterns or activities that may pose a risk to your financial system, helping you maintain a secure and efficient accounting environment with Twisp.

## Monitor and refine velocity controls over time

To ensure the effectiveness of your velocity control measures in mitigating fraud and managing risk, it's crucial to continuously monitor and refine them. Here's how you can do that:

### 1. Regularly review transaction logs and alerts for accuracy and effectiveness
- Schedule periodic reviews of transaction logs to identify any unusual activity or patterns that may have been missed by the velocity controls.
- Analyze alerts triggered by the velocity control rules to ensure they're accurate and relevant. This will help you fine-tune the rules and limits over time.

### 2. Adjust rules and limits based on observed trends and changing risk levels
- As you gather more data on the transaction patterns of your users, take note of any trends or changes in their behavior. Use this information to adjust the velocity control rules and limits accordingly.
- Keep an eye on evolving industry trends and risks that may impact your business. Adjust your velocity control measures in response to these changes to stay ahead of potential threats.

### 3. Continuously improve velocity control measures to mitigate fraud and manage risk
- Actively seek feedback from your team and users on the effectiveness of your velocity controls. Use this feedback to identify areas for improvement and implement necessary changes.
- Stay informed about new fraud detection and prevention techniques, and incorporate them into your velocity control measures as needed.
- Regularly assess the performance of your velocity controls in preventing fraud and managing risk. Use this information to optimize your control measures and ensure they remain effective over time.

By closely monitoring and refining your velocity controls, you can create a robust and adaptive system that effectively mitigates fraud and manages risk for your financial products built with Twisp.

### Conclusion

By following the steps in this guide, you will have successfully implemented velocity controls in Twisp, enhancing your risk management capabilities and strengthening your organization's fraud prevention mechanisms. Remember, velocity controls are not a set-and-forget solution.

Regular review of transaction logs and alerts, as well as ongoing refinement of rules and limits, is crucial for maintaining the effectiveness of your controls. As your business and risk landscape evolve, so too should your velocity control measures. By doing so, you ensure that your organization stays protected against fraud and is ready to handle the ever-changing demands of the financial world. Happy monitoring!