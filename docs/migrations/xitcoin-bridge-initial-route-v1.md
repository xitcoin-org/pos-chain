# Xitcoin bridge initial route v1

Upgrade name: `xitcoin-bridge-initial-route-v1`

This coordinated binary upgrade exposes an authority-gated SDK message that
can create the first bridge route after genesis. The upgrade handler only runs
normal module migrations; it does not configure, enable, resume, or use a
bridge route.

After the upgrade, the configured administrative multisignature may submit
`MsgInitializeRouteConfig` once. The message always stores the route with
`enabled=false` and starts it paused. Route activation is a separate operation
and remains forbidden until the testnet acceptance matrix has passed.

The initialization message has no mint, burn, bank, reserve, relayer, or
transfer capability. It must not be used with Cronos mainnet chain ID `25` or
production XTC contracts.
