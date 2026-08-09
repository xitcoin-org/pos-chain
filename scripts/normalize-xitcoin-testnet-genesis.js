'use strict';

const fs = require('fs');

const [input] = process.argv.slice(2);
if (!input) throw new Error('Usage: normalize-xitcoin-testnet-genesis.js <genesis.json>');

const genesis = JSON.parse(fs.readFileSync(input, 'utf8'));
const app = genesis.app_state;

if (genesis.chain_id !== 'xitcoin-testnet-1') {
  throw new Error(`Unexpected Cosmos chain ID: ${genesis.chain_id}`);
}

function requirePath(object, keys) {
  let value = object;
  for (const key of keys) {
    if (value == null || !(key in value)) {
      throw new Error(`Missing genesis path: ${keys.join('.')}`);
    }
    value = value[key];
  }
  return value;
}

const denom = 'xits';

requirePath(app, ['evm', 'params']).evm_denom = denom;
requirePath(app, ['evm', 'params', 'extended_denom_options']).extended_denom = denom;
requirePath(app, ['staking', 'params']).bond_denom = denom;
requirePath(app, ['mint', 'params']).mint_denom = denom;

const mint = requirePath(app, ['mint']);
const mintParams = requirePath(app, ['mint', 'params']);

mint.minter.inflation = '0.000000000000000000';
mint.minter.annual_provisions = '0.000000000000000000';

mintParams.inflation_rate_change = '0.000000000000000000';
mintParams.inflation_max = '0.000000000000000000';
mintParams.inflation_min = '0.000000000000000000';
mintParams.max_supply = '5250000000000000000000000000';

for (const key of ['min_deposit', 'expedited_min_deposit']) {
  for (const coin of requirePath(app, ['gov', 'params', key])) {
    coin.denom = denom;
  }
}

const bank = requirePath(app, ['bank']);
bank.denom_metadata = [{
  description: 'Native token of the Xitcoin blockchain',
  denom_units: [
    { denom: 'xits', exponent: 0, aliases: [] },
    { denom: 'XTC', exponent: 18, aliases: [] }
  ],
  base: 'xits',
  display: 'XTC',
  name: 'Xitcoin',
  symbol: 'XTC',
  uri: '',
  uri_hash: ''
}];

const nativeEvmAddress = '0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE';
const erc20 = requirePath(app, ['erc20']);

erc20.params.enable_erc20 = true;
erc20.token_pairs = [{
  erc20_address: nativeEvmAddress,
  denom: denom,
  enabled: true,
  contract_owner: 1
}];
erc20.native_precompiles = [nativeEvmAddress];
erc20.dynamic_precompiles = [];

const serialized = JSON.stringify(genesis);
if (serialized.includes('"stake"')) {
  throw new Error('Residual native denomination "stake" detected');
}

fs.writeFileSync(input, `${JSON.stringify(genesis, null, 2)}\n`, { mode: 0o640 });
console.log('Genesis normalized: Xitcoin / xits / XTC');
