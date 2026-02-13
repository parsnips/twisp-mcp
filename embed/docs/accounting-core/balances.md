---
title: Balances
description: With compute-on-write rollup calculations, balances always reflect the most up-to-date state of an account.
showNext: false
showPrev: true
---

## Accounts and Balances

Balances are auto-calculated sums of the entries for a given account.

Every balance record maintains a `drBalance` for entries on the DEBIT side of the ledger and a `crBalance` for entries on the CREDIT side of the ledger.

In addition, a `normalBalance` is calculated as the difference of `credits - debits` for credit normal accounts or as `debits - credits` for debit normal accounts. See [Chart of Accounts: Credit Normal and Debit Normal](/accounting-core/chart-of-accounts#credit-normal-and-debit-normal).

### A Balance for Every Journal, Currency, and Layer

Accounts have separate balances for every journal, currency, and for each of the three layers: `SETTLED`, `PENDING`, and `ENCUMBRANCE` (see [](/accounting-core/layered-accounting)).

In a simple one-journal, one-currency ledger that only uses the `SETTLED` and `PENDING` layers, accounts will materialize two balances (one for each layer).

More complex ledgers which have multiple journals and currencies will compute more balances for each account. For example, accounts in a ledger with 2 journals, 2 currencies, and using all 3 layers would have 12 balances (2 * 2 * 3).

## Balance Calculations

Balances in Twisp are derived from entries in response to changes in the ledger. These balance calculations are computed on write, not on read.

This means that we can query balances like any other piece of data in the system and have 100% certainty that the balance reflects the current state of the ledger.

### Calculating Normal Balance

To illustrate how normal balances are calculated, let's look at some example tables.

Say we have a ledger with two accounts:

1. **Cash** (debit normal)
2. **Revenue** (credit normal)

Next, let's assume the following entries have posted to our ledger:

| Entry ID | Account | Amount | Direction |
|----------|---------|--------|-----------|
| 1        | Revenue | $500   | CREDIT    |
| 2        | Cash    | $500   | DEBIT     |
| 3        | Revenue | $400   | DEBIT     |
| 4        | Cash    | $400   | CREDIT    |
| 5        | Revenue | $250   | CREDIT    |
| 6        | Cash    | $250   | DEBIT     |

Given this set of entries, can calculate the normal balances for the Cash and Revenue accounts by taking the sum of their credits and debits, and subtracting one from the other according to their normal balance type. Here's what the balances table would look like:

**Balances**

| Account | CR Balance | DR Balance | Normal Balance |
|---------|------------|------------|----------------|
| Cash    | $400       | $750       | $350           |
| Revenue | $750       | $400       | $350           |

Note how the sum of all credits equals the sum of all debits, so we know that this balance sheet is consistent with double-entry accounting principles.

The normal balance gives us useful context-dependent information about the account. For the Cash account, it tells us how much cash we currently have on hand. For the Revenue account, it tells us the net revenues we've earned so far.

## Debits and Credits (DR/CR)

Each entry in a journal either debits or credits an account. Along with a DEBIT/CREDIT direction, the amount of the entry is a signed number. Both negative and positive debits or credits are possible.

The primary reason for this setup is to _make the debit and credit balances for an account meaningful_.

Consider an example where we would like our credit balance to accurately reflect all deposits against an account:

| Entry ID | Account ID | Type    | Debit | Credit   |
|----------|------------|---------|-------|----------|
| 1        | f29f83     | DEPOSIT | -     | $1000.00 |

DEBIT BALANCE: **$0.00**\
CREDIT BALANCE: **$1000.00**\
NORMAL BALANCE: **$1000.00**

In this case our credit balance is $1000. In other words, we've deposited $1000 to this account. We'll assume that the account in question is a credit normal account, meaning that the normal balance is calculated by subtracting debits from credits.

Now, let's say a mistake was made and that transaction amount should have actually been $1200 and we need to correct this.

Because the ledger is immutable and append-only, we cannot erase or change the amount of the original transaction.

To rectify this situation in a way that accurately records the history of events, we _void the transaction_ and repost the corrected version.

| Entry ID | Account ID | Type         | Debit | Credit     |
|----------|------------|--------------|-------|------------|
| 1        | f29f83     | DEPOSIT      | -     | $1000.00   |
| 2        | f29f83     | VOID_DEPOSIT | -     | $(1000.00) |
| 3        | f29f83     | DEPOSIT      | -     | $1200.00   |

DEBIT BALANCE: **$0.00**\
CREDIT BALANCE: **$1200.00**\
NORMAL BALANCE: **$1200.00**

In this scenario above, the credit balance is now $1200. Exactly what we'd expect.

Let's consider an alternate solution: what if we had simply utilized debits and credits to achieve the same resulting account balance?

| Entry ID | Account ID | Type         | Debit    | Credit   |
|----------|------------|--------------|----------|----------|
| 1        | f29f83     | DEPOSIT      | -        | $1000.00 |
| 2        | f29f83     | VOID_DEPOSIT | $1000.00 | -        |
| 3        | f29f83     | DEPOSIT      | -        | $1200.00 |


DEBIT BALANCE: **$1000.00**\
CREDIT BALANCE: **$2200.00**\
NORMAL BALANCE: **$1200.00**

In this incorrect version, the _credit balance_ is now $2200 along with a _debit balance_ of $1000. This effectively lost meaning of the sum total debit and credit balances, because money wasn't _actually_ moving out of the account. We were just cancelling

{% callout type="note" %}

For the sake of simplicity in illustration, we did not specify the layer used. We can assume that they all occurred on the "SETTLED" layer.

Learn more about layers in [](/accounting-core/layered-accounting).

{% /callout %}
