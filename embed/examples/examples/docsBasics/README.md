# Common Examples for Documentation

Set of common examples for use in docs site.

## The Contrived Premise

The idea for the examples is to showcase a ledger system for a contrived budgeting tool, like a overly-simplified version of Mint or YNAB.

Also has a wallet feature (like Venmo), and users can transfer funds between wallets.

Within this tool, users can:

- add expected spending (encumb. layer)
  - with metadata on tx: budget category
  - TC: allocate $$ to budget cateogry
  - TC: deallocate $$ from budget cateogry
- record a tx from the real world
  - e.g.: withdraw_hold, withdraw, deposit, deposit_hold, card_hold, card_settle
  - TC: add pending tx (pend. layer)
  - TC: settle a tx (settled layer)
- assign a TX amount to a budgeted category
  - TC: assign tx to budget
- send/receive wallet funds
  - TC: wallet_xfer
- manage wallet
  - TC: wallet_deposit
  - TC: wallet_withdraw

Chart of accounts:

- users (DN)
  - user a (DN)
    - budget (CN)
    - check (DN)
    - card (DN)
    - wallet (DN)
  - user b (DN)
    - budget (CN)
    - check (DN)
    - cash (DN)
    - wallet (DN)
- settlement (CN)
  - cash (CN)
  - check (CN)
  - card (CN)
- wallets (DN)
- budgets (CN)