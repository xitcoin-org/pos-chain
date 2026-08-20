# Xitcoin 2026 Sovereign Allocation Index

## Reference framework

The index defines a deterministic reference allocation for 195 positions: 193
United Nations Member States, the Holy See and the State of Palestine.

The index records the fixed allocation assigned to each State position. It
does not, by itself, transfer tokens, activate a validator or create automatic
validator admission. Release requires an activated position and the protocol
controls defined in `docs/sovereign-validator-framework.md`.

## Reserve and formula

The fixed reference reserve is 390,000,000 XTC:

- 292,500,000 XTC equal component (75%);
- 97,500,000 XTC demographic component (25%).

Each position receives an equal base of 1,500,000 XTC. The demographic
component is calculated as follows:

```text
allocation_i = 390,000,000 x (
  0.75 / 195
  + 0.25 x sqrt(population_i_2026)
    / sum(sqrt(population_2026))
)
```

Calculations use deterministic decimal arithmetic. Final atomic-unit
remainders are distributed by descending fractional remainder, with ISO3 code
as the deterministic tie-breaker.

The aggregate allocation is exactly 390,000,000.000000000000000000 XTC.

## Economic function

The allocation is a finite protocol-funded contribution supporting an
activated sovereign validator position. It is additional to, and never a
substitute for, the State-provided minimum self-delegation.

After activation, the fixed ISO3 allocation is intended to vest linearly over
five years of eligible service under a schedule created specifically for that
position. Activation may occur at any future date and requires no global
calendar change. The schedule is subject to the institutional, staking,
availability and security conditions defined by the sovereign validator
framework. It does not create new supply.

Ordinary staking rewards remain separate. After the finite allocation has
fully vested, the validator continues under the same staking, commission,
delegation, fee-distribution and slashing rules as other validators.

## Population reference

Population values use the United Nations World Population Prospects 2024
medium variant for 1 July 2026.

Source: `https://population.un.org/wpp/assets/Excel%20Files/1_Indicator%20(Standard)/EXCEL_FILES/1_General/WPP2024_GEN_F01_DEMOGRAPHIC_INDICATORS_COMPACT.xlsx`

Source SHA-256: `98e34d9b65b53858cd08a57a566e45050b08093ad85ba5714fe6fbd78055ae6d`

## Statistical consolidation

Statistical records without a separate reference position are included only
where the consolidation table identifies a reference position.

Under the United Nations M49 statistical framework, the China position
consolidates China, China Hong Kong SAR, China Macao SAR and China Taiwan
Province of China.

Cook Islands, Niue and Western Sahara remain non-consolidated statistical
records. They do not create additional reference positions and are not
included in another position's population value.

## Validator capacity

The target protocol configuration provides:

- 258 maximum validators;
- 195 sovereign positions;
- 63 public positions;
- 5,000,000 XTC minimum self-delegation for every activated validator.

The same minimum applies to founder, sovereign and public validators. A State
must supply this commitment independently; its sovereign allocation cannot be
counted toward the activation minimum.

Validator admission remains subject to the validator-admission module. The
development branch defines the 5,000,000 XTC protocol minimum but does not yet
implement the sovereign position, institutional succession, mandate or
activation-based five-year vesting controls. Those remaining controls must be
implemented, tested and independently reviewed before activation.

## Institutional continuity

Each ISO3 position remains attached to the relevant State. Successive
administrations of that State may transfer the institutional governance and
operating mandate to their authorized successors without replacing the State
position, its history or its remaining allocation.

The position is not attached permanently to an individual office-holder or
service provider. Buying tokens from a State or former operator does not grant
control of the State position.

## Canonical identifiers

Names, United Nations M49 identifiers and ISO codes are the canonical
identifiers. Flags may be used as non-authoritative presentation assets when
their source and licence are recorded.

## Canonical files

- `networks/testnet/sovereign-allocation-index-2026.csv`
- `networks/testnet/sovereign-allocation-index-2026.json`
- `networks/testnet/territorial-consolidation-2026.csv`
