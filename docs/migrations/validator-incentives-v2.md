# Validator Incentives v2 migration

## Scope

Consensus version 2 replaces the inactive fixed-APR funded-period state with
daily treasury-derived reward accounting.

## State transition

The on-chain plan name is:

```text
xitcoin-validator-incentives-daily-v2
```

The application registers this plan with the canonical module migration
manager. The upgrade changes module state only; it adds or removes no store.

The migration:

1. installs the default 1,000-basis-point treasury release rate;
2. installs 6,311,520 blocks per year and 17,280 blocks per daily period;
3. deletes the obsolete active-period snapshot;
4. preserves the configured authority;
5. preserves lifetime distributed accounting;
6. moves no funds and creates no mint or burn permission.

The first begin block after migration creates a new snapshot from the canonical
Bank and Staking keeper state. An empty treasury or zero eligible bonded stake
produces a zero funded distribution.

## Operational boundary

This migration must be rehearsed on testnet with restart and rollback evidence
before a mainnet genesis or upgrade may activate the module.
