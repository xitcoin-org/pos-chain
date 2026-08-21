# Validator Incentives public release checklist

Status date: 21 August 2026

This checklist is the canonical technical release record. The public GitBook
describes the product and links here instead of duplicating test logs, audit
evidence or deployment procedures.

## Current readiness

- Module implementation: **97%**
- Public production readiness: **72%**
- Production deployment: **not performed**
- Production funds moved: **none**

The implementation percentage measures source completion. Public production
readiness also includes independent review, governance rehearsal, supply and
bridge reconciliation, release engineering, testnet evidence, monitoring and
rollback preparation.

## Completed

- [x] funded-period reward policy and bounded APR parameters;
- [x] canonical Staking keeper used for eligible bonded stake;
- [x] canonical Bank keeper used for treasury balance;
- [x] governance module account configured as authority;
- [x] treasury account has no mint or burn permission;
- [x] public read-only parameter, period and treasury queries;
- [x] CLI and Cosmos message handlers;
- [x] application, store, keeper and genesis wiring;
- [x] governance, staking and EIP integration regression tests;
- [x] initial Ethereum Lists metadata merged through PR #8608;
- [x] DefiLlama Chainlist correction submitted as PR #3073;
- [x] public guide separated from technical verification records.

## Required before public activation

- [ ] complete an independent security review of the final module diff;
- [ ] resolve all review findings and rerun the full repository CI matrix;
- [ ] reconcile canonical total supply across PoS, Xitcoin EVM and Cronos EVM;
- [ ] document bridge lock/mint and burn/unlock accounting with invariants;
- [ ] define the initial treasury funding transaction and source-of-funds record;
- [ ] rehearse parameter update and funded-period activation through governance;
- [ ] publish the exact genesis migration or chain-upgrade plan;
- [ ] execute a persistent testnet deployment with public endpoints;
- [ ] run load, failure, restart, upgrade and rollback exercises on testnet;
- [ ] configure metrics, alerts, dashboards and incident runbooks;
- [ ] obtain release approval with reproducible build artifacts and checksums;
- [ ] publish final user documentation and deployment evidence;
- [ ] perform a limited production rollout before unrestricted activation.

## External publication follow-up

- Ethereum Lists PR #8608: merged; corrective PR #8621 is open, CI-green and
  ready for maintainer review.
- DefiLlama Chainlist PR #3073: open and mergeable; wait for two maintainer
  approvals and do not rewrite the branch unless CI or reviewers request it.
- Xitcoin Guide PRs #9 and #10: merged; obsolete staging and pending-cutover
  language is removed and protected by Guide CI.

## Release rule

No percentage, passing test suite or merged metadata pull request authorizes a
production deployment. Activation requires every applicable item above to be
closed with evidence and a separately approved operational change.
