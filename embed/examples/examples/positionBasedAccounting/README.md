# Position Based Accounting

Examples of position-based accounting transactions using TranCodes for writing entries using standard currencies (USD) as well as custom currencies for securities (AAPL, MSFT).

## Accounts

- `BROKER.REVENUE` - (revenue) credit-normal
- `BROKER.SETTLEMENT` - (asset) debit-normal
- `BROKER.INVESTMENT` - (equity) credit-normal
- `BROKER.HOLDINGS` - (asset) debit-normal
- `CLIENT_*.CASH` - (asset) debit-normal
- `CLIENT_*.INVESTMENT` - (equity) credit-normal
- `CLIENT_*.HOLDINGS` - (asset) debit-normal
- `CLIENT_*.DIVIDEND_REVENUE` - (revenue) credit-normal

## Transactions

### BUY_SECURITY

> Buy 100 AAPL @ $149.45 / share with a 0.75% brokerage fee.

#### Params

| Param                | Type      | Example             |
|----------------------|-----------|---------------------|
| `clientHoldingsAcct` | `UUID`    | CLIENT_A.HOLDINGS   |
| `clientInvestAcct`   | `UUID`    | CLIENT_A.INVESTMENT |
| `clientCashAcct`     | `UUID`    | CLIENT_A.CASH       |
| `security`           | `STRING`  | AAPL                |
| `shares`             | `DECIMAL` | 100                 |
| `price`              | `DECIMAL` | 149.45              |
| `feeRate`            | `DECIMAL` | 0.0075              |
| `effective`          | `DATE`    | 2022-10-12          |

#### Entries

| EntryType                       | Account             | Units      | Currency | Direction | Layer   |
|---------------------------------|---------------------|------------|----------|-----------|---------|
| BUY_SECURITY_PAYMENT_CR         | CLIENT_A.CASH       | 15057.0875 | USD      | CREDIT    | SETTLED |
| BUY_SECURITY_PAYMENT_DR         | CLIENT_A.INVESTMENT | 15057.0875 | USD      | DEBIT     | SETTLED |
| BUY_SECURITY_BROKER_EARNINGS_CR | BROKER.REVENUE      | 14945      | USD      | CREDIT    | SETTLED |
| BUY_SECURITY_BROKER_FEE_CR      | BROKER.REVENUE      | 112.0875   | USD      | CREDIT    | SETTLED |
| BUY_SECURITY_BROKER_EARNINGS_DR | BROKER.SETTLEMENT   | 15057.0875 | USD      | DEBIT     | SETTLED |
| BUY_SECURITY_TRADE_CR           | BROKER.HOLDINGS     | 100        | AAPL     | CREDIT    | SETTLED |
| BUY_SECURITY_TRADE_DR           | CLIENT_A.HOLDINGS   | 100        | AAPL     | DEBIT     | SETTLED |


### SELL_SECURITY

> Sell 50 AAPL @ $124.72 / share with a 1% brokerage fee.

#### Params

| Param                | Type      | Example             |
|----------------------|-----------|---------------------|
| `clientHoldingsAcct` | `UUID`    | CLIENT_A.HOLDINGS   |
| `clientInvestAcct`   | `UUID`    | CLIENT_A.INVESTMENT |
| `clientCashAcct`     | `UUID`    | CLIENT_A.CASH       |
| `security`           | `STRING`  | AAPL                |
| `shares`             | `DECIMAL` | 50                  |
| `price`              | `DECIMAL` | 124.72              |
| `originalShareCost`  | `DECIMAL` | 150.570875          |
| `feeRate`            | `DECIMAL` | 0.01                |
| `effective`          | `DATE`    | 2022-10-12          |

#### Entries

| EntryType                         | Account             | Units       | Currency | Direction | Layer   |
|-----------------------------------|---------------------|-------------|----------|-----------|---------|
| SELL_SECURITY_TRADE_CR            | CLIENT_A.HOLDINGS   | 50          | AAPL     | CREDIT    | SETTLED |
| SELL_SECURITY_TRADE_DR            | BROKER.HOLDINGS     | 50          | AAPL     | DEBIT     | SETTLED |
| SELL_SECURITY_CLIENT_EARNINGS_CR  | CLIENT_A.INVESTMENT | 7528.54375  | USD      | CREDIT    | SETTLED |
| SELL_SECURITY_CLIENT_GAIN_LOSS_CR | CLIENT_A.INVESTMENT | -1354.90375 | USD      | CREDIT    | SETTLED |
| SELL_SECURITY_CLIENT_EARNINGS_DR  | CLIENT_A.CASH       | 6173.64     | USD      | DEBIT     | SETTLED |
| SELL_SECURITY_BROKER_PAYMENT_CR   | BROKER.SETTLEMENT   | 6236        | USD      | CREDIT    | SETTLED |
| SELL_SECURITY_BROKER_PAYMENT_DR   | BROKER.INVESTMENT   | 6236        | USD      | DEBIT     | SETTLED |
| SELL_SECURITY_BROKER_FEE_CR       | BROKER.REVENUE      | 62.36       | USD      | CREDIT    | SETTLED |
| SELL_SECURITY_BROKER_FEE_DR       | BROKER.SETTLEMENT   | 62.36       | USD      | DEBIT     | SETTLED |


### EARN_DIVIDENDS

> Earn $3.62 dividend per share of AAPL stock.

#### Params

| Param                   | Type      | Example                   |
|-------------------------|-----------|---------------------------|
| `clientCashAcct`        | `UUID`    | CLIENT_A.CASH             |
| `clientDividendRevAcct` | `UUID`    | CLIENT_A.DIVIDEND_REVENUE |
| `security`              | `STRING`  | AAPL                      |
| `shares`                | `DECIMAL` | 50                        |
| `dividendsPerShare`     | `DECIMAL` | 3.62                      |
| `effective`             | `DATE`    | 2022-10-12                |

#### Entries

| EntryType        | Account                   | Units | Currency | Direction | Layer   |
|------------------|---------------------------|-------|----------|-----------|---------|
| EARN_DIVIDEND_DR | CLIENT_A.CASH             | 181   | USD      | DEBIT     | SETTLED |
| EARN_DIVIDEND_CR | CLIENT_A.DIVIDEND_REVENUE | 181   | USD      | CREDIT    | SETTLED |
