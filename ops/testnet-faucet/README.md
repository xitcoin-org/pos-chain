# Xitcoin Testnet Faucet

Faucet public officiel du réseau `xitcoin-testnet-1`.

## Paramètres

- Distribution : 100 XTC par demande
- Dénomination technique : axtc
- Limite par adresse : 1 demande par 24 heures
- Limite par IP : 3 demandes par heure
- Backend : Node.js derrière Nginx et Cloudflare

## Déploiement

1. Copier `faucet.env.example` vers `/etc/xitcoin-testnet/faucet.env`.
2. Adapter les chemins et le RPC sans ajouter de secret au dépôt.
3. Installer le service systemd et la configuration Nginx.
4. Importer séparément la clé `faucet-official` dans le keyring sécurisé.

Les fichiers de clé, mots de passe, états et portefeuilles ne doivent jamais être commités.
