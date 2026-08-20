# Xitcoin Sovereign Validator Framework

Status: proposed protocol policy, not deployed  
Reference date: 20 August 2026

## Purpose

The framework reserves one institutional validator position for each of the
195 sovereign references recorded in the Xitcoin 2026 Sovereign Allocation
Index. A position is attached permanently to the relevant State, identified
by its canonical ISO3 code. It is not attached to a president, minister,
administration, signatory or infrastructure provider.

Successive administrations of the same State may transfer the institutional
governance mandate to their authorized representatives and operating teams
without changing the State position, its accounting history or its remaining
allocation.

## Validator capacity

The target validator capacity is:

- 195 reserved sovereign positions;
- 63 public positions;
- 258 positions in total;
- 5,000,000 XTC minimum self-delegation for every activated validator.

The same minimum applies to founder, sovereign and public validators. A
reserved sovereign position is not active and receives no reward while its
admission conditions are not satisfied.

## Institutional control

Each sovereign position separates four roles:

1. **State position** — the permanent ISO3 position and its immutable history;
2. **institutional controller** — the currently authorized State-controlled
   group or multisignature policy;
3. **operational mandatary** — a time-limited and revocable operating mandate;
4. **validator operator** — the team authorized to operate the validator for
   the duration of that mandate.

A change of government or service provider must not sell, extinguish or move
the State position. It transfers the governance and operating mandate between
successive authorized teams representing the same State. The transition must
support a future effective date, the end of the former mandate, the beginning
of the successor mandate and a complete public event history.

No individual mandatary may permanently appropriate the position or obtain an
unrestricted right to transfer the position's controlled stake.

## Activation requirements

A sovereign position may be activated only after all of the following are
verified:

- the State identity and ISO3 position;
- a valid institutional controller;
- a dated operational mandate;
- compliant validator infrastructure operated under the active mandate;
- at least 5,000,000 XTC of State-provided self-delegation;
- approved reward and withdrawal destinations;
- acceptance of uptime, slashing, reporting and operational-security requirements.

The sovereign allocation does not satisfy the minimum self-delegation. The
State must provide the required commitment independently.

## Sovereign allocation reserve

The fixed reserve remains 390,000,000 XTC. The allocation of each ISO3
position remains calculated by the canonical 75% equal / 25% demographic
formula recorded in `docs/sovereign-allocation-2026.md` and the associated CSV
and JSON files.

The allocation is a finite protocol-funded support grant for an activated
sovereign validator position. It is separate from ordinary validator rewards,
does not create new supply and does not grant automatic admission.

## Five-year release

After activation, the position's fixed allocation is released over five years
in 20 quarterly tranches:

```text
quarterly tranche = canonical ISO3 allocation / 20
```

A tranche becomes claimable only if the position satisfies the applicable
service conditions for the quarter, including:

- minimum self-delegation;
- valid institutional and operational mandates;
- validator availability;
- absence of unresolved double-signing or other disqualifying fault;
- compliance with reporting and security obligations.

A failed condition pauses the affected tranche and opens a defined remediation
period. It must not permit arbitrary seizure of already vested tranches.
Cancelled or unvested amounts remain in the sovereign allocation reserve.

After all 20 valid tranches have been released, no further sovereign allocation
is created. The validator continues under the same staking, commission,
delegation, fee-distribution and slashing rules as other validators.

## Ordinary validator rewards

Ordinary rewards are distinct from the sovereign allocation. An active
sovereign validator may receive, under the same network rules as other active
validators:

- rewards attributable to its self-delegation;
- validator commission on third-party delegations;
- its share of funded validator incentives;
- its share of transaction-fee distribution.

The sovereign allocation must not be included when calculating whether the
State supplied the 5,000,000 XTC activation commitment.

The Validator Incentive treasury remains pre-funded and has no mint or burn
permission. Transaction fees are accounted for by the chain distribution
mechanism unless a separately reviewed governance upgrade defines another
route.

## Continuity and suspension

The State position remains recorded if a government changes, a mandate
expires, an operator is replaced or the validator is suspended.

If no valid mandatary exists, or self-delegation falls below the required
minimum:

- validation and future allocation releases are suspended;
- the position remains reserved for the same State;
- already vested ownership and the accounting history remain recorded;
- reactivation requires the admission conditions to be satisfied again.

Buying tokens from a State or former operator never transfers control of the
State position to the buyer.

## Required protocol implementation

This policy requires explicit implementation and independent review before it
can be activated. At minimum, the protocol must provide:

- a unique 195-entry ISO3 sovereign position registry;
- separate sovereign and public admission capacity;
- institutional-controller and dated-mandate records;
- safe institutional succession between successive governance and operating
  mandates of the same State;
- quarterly vesting and remaining-allocation accounting;
- reward accounting separated from sovereign allocation accounting;
- suspension, remediation and reactivation procedures;
- immutable events and public read-only queries;
- invariants preventing double allocation, unauthorized control, premature
  release, reserve overdraft and mint-based funding.

Until those controls are implemented, tested, audited and activated by the
required governance process, the allocation index remains a reference only and
does not transfer tokens or activate validators.
