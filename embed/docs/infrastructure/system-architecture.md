---
title: System Architecture
description: |
 Twisp is a purpose-built system that provides a mission-critical, high performance, core ledger for constructing and running financial products, empowering organizations to have complete control over their crucial financial data.
showNext: true
showPrev: true
---

## Integrating with Twisp

Integrating with the Twisp Accounting Core streamlines financial data flows, structures ledger data, automates financial processes, and helps you design and build seamless financial products.

In this document, we will discuss how to integrate your system with Twisp, covering common components in your stack and the interfaces available to connect those components to the Twisp Accounting Core.

{% arch /%}

## Your System

Your system is the foundation for your integration with Twisp. It includes the frontend, infrastructure, services, and vendors that you interact with.

### Frontend

Your frontend is the user interface that your customers use to interact with your system. By combining fine grained security policies with OIDC authentication, it's possible to use our GraphQL API directly from frontend applications. This provides a scalable and performant mechanism to enable your customers to view and interact with their financial data stored in Twisp's ledger database.

### Infrastructure

Your infrastructure is the backbone of your system. It includes servers, streaming platforms, databases, warehouses, and networking components. These lower-level technology components often require access direct access to financial data stored in the Twisp Ledger Database. Real-time streaming and batch access to this data is provided by Twisp data connectors.

### Services

Your services are self-contained, loosely coupled software components that provide specific functionality to power your business. To integrate with Twisp, you'll need to identify the specific financial services that these systems provide, such as brokerage, checking, savings, credit, or loan products. The logic for that product is then setup in the Twisp Accounting Core, and services then integrate directly via Twisp API interfaces with GraphQL, gRPC, or REST for data exchange.

### Vendors

Your vendors are the third-party tools and services that you use to support your system. Twisp has an expanding library of pre-built vendor and protocol integrations that are connected directly to the Accounting Core. This tooling allows simple plug-and-play integration of your Accounting Core with existing investments.

## Interfaces

Interfaces are the methods by which your system interacts with Twisp. They include the API, data, integrations, and protocols that enable communication between your system and Twisp's accounting engine and ledger database.

### API

Twisp includes a number of modern API interfaces, allowing your teams to integrate with technologies best suited to your organization.

| Technology | Description |
|---|---|
| GraphQL | GraphQL is an open-source query language for APIs (Application Programming Interfaces) and a runtime for executing those queries with existing data. It was developed by Facebook and released in 2015. GraphQL enables clients to request specific data from a server, allowing them to retrieve only the information they need and nothing more. |
| gRPC | gRPC (Google Remote Procedure Call) is an open-source high-performance framework developed by Google. It allows you to define services and message types using Protocol Buffers (protobuf) and enables efficient communication between distributed systems. |
| REST | REST (Representational State Transfer) is an architectural style for designing networked applications. It provides a set of principles and constraints for building web services that are scalable, stateless, and interoperable. Additionally, the OData protocol allows Twisp to easily integrate with services such as Salesforce out of the box. |

### Data

Twisp provides various data integration interfaces to facilitate the seamless processing and transfer of data. These interfaces cater to both real-time and batch scenarios, allowing you to effectively work with data in Twisp.

#### Real-Time Data Integration

Real-time data integration involves the continuous or near-instantaneous processing and transfer of data between systems as it becomes available. Twisp supports the following real-time data interfaces:

 * **Streaming**: Streams provide near real-time change data capture capabilities, allowing you to stream data to various technologies such as AWS Kinesis, Kafka, or webhooks. This ensures that you can stay up-to-date with the latest changes in your data.
 * **Webhooks**: Twisp supports webhooks for real-time HTTPS data integration. You can configure webhooks to receive instant notifications or data updates from external systems, enabling seamless communication and data synchronization.

#### Batch Data Integration

Batch data integrations involve processing and transferring data in bulk at scheduled intervals or on-demand. Twisp supports the following batch data interfaces:

 * **SQL**: The SQL interface enables you to perform ad-hoc queries for generating reports, exporting data, or integrating with popular business intelligence tools like Tableau or PowerBI.
 * **Warehousing**: Twisp aggregates all committed financial data in near real-time and stores it in the industry-standard [Parquet](https://parquet.apache.org) file format. This format is compatible with leading data warehousing systems like Redshift, Snowflake, or BigQuery. You can easily ingest this data into your preferred data warehousing system for further analysis and reporting.

By leveraging these data integration interfaces in Twisp, you can effectively manage your data, whether it's processing real-time updates or performing bulk transfers at your convenience.

### Integrations

Integrations are the connections between the Accounting Core and other financial systems and business tools. Twisp includes an expanding number of plug-and-play integrations. If an integration you are interested in is not included in this list, please reach out and we can prioritize development.

| Vendor | Description                              |
|---------|-----------------------------------------------------------------------|
| Fiserv | Twisp can integrate with a number of Fiserv core systems.       |
| Lithic | Twisp natively processes Lithic transaction webhooks, including ASA.  |
| Marqeta | Twisp natively processes Marqeta transaction webhooks, including JIT. |
| Stripe | Twisp natively processes Stripe Issuing and Treasury webhooks.    |

### Protocols

Protocols are industry-standard rules and formats that govern structured communication between financial systems. Our platform includes plug-and-play protocol adapters that process, store, and generate data while adhering to strict security and compliance standards. If you're interested in a protocol not listed below, please contact us, and we will prioritize its development.

| Protocol | Description |
| --- | --- |
| ISO 8583 | Seamlessly integrate with several card networks, directly consuming ISO 8583 messages and converting them into transaction entries within our platform. |
| ISO 20022 | Process ISO 20022 messages, encompassing Swift and FedNOW transactions. Additionally, generate ISO 20022 messages from ledger data. |
| NACHA | Proess and generate NACHA ACH files. |
| X9 | X9 Cash image files serve as the standard method for check clearance in the US. Twisp can efficiently consume X9 files and track check transactions in their native format. |
| BAI2 | Comprehensive support for financial institutions generating or processing BAI2 files. |

## Core

The Twisp Accounting Core is the heart of your integration with Twisp. It includes the accounting engine, ledger database, product logic, and security policies that enable secure, reliable, and compliant financial operations.

### Accounting Engine

The accounting engine is the core of the Twisp Accounting Core. It includes the chart of accounts, transaction codes, and ledger entries that enable processing and accounting for financial transactions.

### Ledger Database

The ledger database is where your financial data is stored within Twisp. Twisp is built on a production hardened, high performance, horizontally scalable database storage system: [AWS DynamoDB](https://aws.amazon.com/dynamodb). The ledger database leverages DynamoDB using a formally verified transaction layer, which provides transactionally consistent, snapshot isolated, immutable storage, providing a complete historical lineage of any object.

### Product Logic

Product logic is the essential business logic that powers your financial products and services. With Twisp, you have complete control over product development in the Accounting Core, allowing you to design and create a wide range of financial products. Twisp also offers pre-built financial product definitions that can be customized or extended according to your specific requirements.

| Product | Description |
| --- | --- |
| Brokerage | Brokerages have unique accounting needs, including the representation of numerous positions within a single account. Twisp enables you to execute the trading lifecycle and reconcile trades with a clearinghouse, all while operating at high scale during market hours. |
| Cards | Simplify your card issuing program by leveraging Twisp's card issuing workflows and transaction codes. Twisp integrates seamlessly with issuer processors, providing robust accounting control over the cards you issue. You can enable features such as spend limits, wallets, and organizational structures tailored to your needs. |
| Credit | Credit card programs require specific functionalities beyond debit cards. Twisp's basic card program can be enhanced with SCRA processes, available credit balance tracking, and customizable interest rate computations. |
| Deposits | Twisp supports various deposit-taking modalities. Whether you offer a card product on top of FBO (For Benefit Of) accounts or provide demand deposit accounts (DDAs), Twisp's deposit program can be tailored to suit your requirements. |
| Lending | Enhance your origination and servicing capabilities by incorporating Twisp's lending components. Account for payments, generate or regenerate amortization tables, and enable portfolio monitoring directly from the accounting core. |
| Mortgages | Manage mortgage accounting processes for origination and servicing directly from the Accounting Core.  |
| Payments | Twisp's multi-currency design enables you to manage global payments accounting with FX (Foreign eXchange.)  |
| Wallets | Simplify the management of your organization's spend program with Twisp's powerful wallet management. Configure your organizational hierarchy, assign spend limits, and Twisp's wallet program will effectively manage spending. If you choose to issue cards associated with wallets, Twisp seamlessly integrates with various card module partners to provide the necessary capabilities. |

### Security Policies

Security policies are the rules and procedures that ensure the security and compliance of your system and financial data. Twisp supports unlimited fine grained access control policies that authorize all operations executed against the accounting core. Additionally, Twisp is a credential-less system, leveraging your existing identity or cloud providers identities and roles for all authentication.
