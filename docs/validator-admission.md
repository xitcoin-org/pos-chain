# Validator Admission

Xitcoin uses an on-chain validator admission policy.

## Policy

| Parameter | Initial value |
| --- | ---: |
| Maximum approved validator operators | 208 |
| Minimum self-delegation | 5,000,000 XTC |

Amounts are encoded in `axtc`, with 18 decimal places.

## Admission flow

1. An operator submits its validator operator address for review.
2. The configured authority approves the operator address on-chain.
3. The operator submits `MsgCreateValidator` with the required self-delegation.
4. The ante handler verifies approval and the active policy.
5. The staking module applies the standard consensus and staking checks.

## Revocation

Revocation removes the operator from the approval set. An existing validator is jailed and cannot unjail until it is approved again.

## Policy updates

The configured authority may update the approval capacity and minimum self-delegation. The approval capacity cannot be reduced below the number of approved operators.

## Query and transaction commands

Discover available commands:

```bash
xitcoind query validator-admission --help
xitcoind tx validator-admission --help
```

Inspect transaction flags:

```bash
xitcoind tx validator-admission update-params --help
```

Always specify the intended chain ID, node endpoint and fee denomination when preparing a transaction.
