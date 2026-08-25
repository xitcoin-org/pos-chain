# Validator Incentive governance operations

Status: development branch, not deployed
Reference date: 20 August 2026

## Authority model

The Validator Incentive module uses the Cosmos governance module account as
its default on-chain authority.

This authority applies only to:

- updating funded incentive parameters;
- activating a fully funded reward period.

It does not approve or revoke validators. Validator admission remains a
separate policy and execution path.

The module account has no private key. An authorized message is executed
through the network governance process. A direct transaction signed by a
personal account cannot impersonate the governance module account.

## Safety boundaries

The module enforces the following boundaries:

- initial annual rate: 800 basis points (8%);
- protocol ceiling: 2,000 basis points (20%);
- maximum increase: 100 basis points per reward period;
- each reward period must be funded before activation;
- eligible bonded stake is read from the Staking keeper;
- treasury balance is read from the Bank keeper;
- the transaction cannot provide either value;
- the treasury module account has no mint or burn permission;
- a distribution cannot exceed the active period provision;
- failed bank transfers do not advance distribution accounting.

The bridge is a separate subsystem. Its verified lock/mint and burn/unlock
operations do not grant mint authority to the Validator Incentive module.

## Read-only endpoints

The module exposes:

- `GET /cosmos/evm/validatorincentives/v1/params`
- `GET /cosmos/evm/validatorincentives/v1/period`
- `GET /cosmos/evm/validatorincentives/v1/treasury`

These routes expose no state mutation.

## Governance message templates

Replace `<GOVERNANCE_MODULE_ADDRESS>` with the governance module address
derived for the target network. Do not substitute a personal wallet address
when governance is the configured authority.

### Update funded incentive parameters

```json
{
  "@type": "/cosmos.evm.validatorincentives.v1.MsgUpdateParams",
  "authority": "<GOVERNANCE_MODULE_ADDRESS>",
  "annual_rate_basis_points": 800,
  "blocks_per_year": "6311520",
  "reward_period_blocks": "1577880"
}
```

### Activate a funded period

The budget is expressed in atomic `axtc`. Eligible stake and treasury
balance are intentionally absent.

```json
{
  "@type": "/cosmos.evm.validatorincentives.v1.MsgActivateFundedPeriod",
  "authority": "<GOVERNANCE_MODULE_ADDRESS>",
  "committed_annual_budget_atomic": "<POSITIVE_INTEGER_AXTC>"
}
```

## Pre-submission checks

Before submitting an activation proposal:

1. query the live module parameters;
2. query the current period and confirm there is no overlap;
3. query the actual treasury balance;
4. verify the proposed annual budget does not exceed available funding;
5. verify the calculated period provision;
6. verify the proposed rate transition is permitted;
7. publish the proposal payload and expected accounting result;
8. execute only after the normal governance process succeeds.

## Release boundary

These instructions describe the development implementation. They are not a
deployment record. Production activation additionally requires integration
tests, supply reconciliation, an upgrade or genesis plan, independent security
review and public deployment evidence.
