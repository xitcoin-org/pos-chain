# Contributing to Xitcoin

## Before opening a pull request

- Never commit secrets, mnemonics, private keys, passwords, API tokens or node data.
- Do not commit `.env` files. Use `.env.example` files with empty values only.
- Keep changes focused and document their purpose.
- Run the relevant tests before submitting code.
- Preserve upstream license notices and attribution.

## Public testnet changes

Changes affecting consensus, supply, validator admission, fees, token metadata, genesis generation or public network configuration require explicit review before release.

The public testnet and mainnet must use independent keys and infrastructure. Do not reuse a private-testnet key or configuration.
