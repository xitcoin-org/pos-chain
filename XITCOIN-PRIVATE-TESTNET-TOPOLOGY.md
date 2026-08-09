# Topologie privée Xitcoin testnet

## Objectif

Déployer un testnet PoS Xitcoin avec quatre validateurs autorisés, sans exposition publique inutile et sans créer de token distinct sur l'EVM.

- Actif natif unique : XTC
- Nom : Xitcoin
- Unité technique interne : xits, 18 décimales
- Inflation : 0 %
- Cap : 5 250 000 000 XTC
- Représentation EVM native : même XTC, sans WXTC
- Admission des validateurs : permissionnée par la Fondation

## Répartition des validateurs

| Validateur | Emplacement | Rôle |
|---|---|---|
| 1 | Serveur actuel kcalb-servor-1 | Validateur testnet |
| 2 | Environnement indépendant | Validateur testnet |
| 3 | Environnement indépendant | Validateur testnet |
| 4 | Environnement indépendant | Validateur testnet |

Les validateurs 2, 3 et 4 doivent être répartis sur au moins deux fournisseurs ou régions distinctes du serveur actuel.

## Réseau

- P2P CometBFT : port TCP 26656, autorisé seulement entre les IP des validateurs.
- RPC Cosmos 26657 : non exposé publiquement sur les validateurs.
- RPC EVM 8545 et WebSocket 8546 : non exposés publiquement sur les validateurs.
- API et gRPC : non exposés publiquement sur les validateurs.
- Les futurs RPC publics et explorateurs seront hébergés sur des machines séparées.

## Gouvernance et clés

- Fondation : multisig matériel 3 sur 5 avant le genesis permanent.
- Aucun secret, clé privée ou phrase de récupération dans Git, dans le genesis ou dans les scripts.
- Chaque validateur utilise ses propres clés de consensus et d’opérateur.
- Toute clé de sauvegarde est chiffrée hors du serveur.

## Conditions avant lancement du testnet

- Quatre environnements prêts.
- IP publiques et ports P2P des quatre validateurs recensés.
- Multisig Fondation 3 sur 5 créé.
- Genesis relu, hashé et approuvé manuellement.
- Liste des validateurs approuvés injectée au genesis.
- Procédure de sauvegarde et restauration testée.
- Binaire vérifié contre le reçu de livraison.

## Critères de sortie du testnet

- Fonctionnement continu pendant 14 jours.
- Test de perte et redémarrage d'un validateur.
- Test de révocation et de réadmission d'un validateur.
- Test de restauration depuis sauvegarde.
- Vérification Cosmos et EVM : même XTC, sans WXTC.
- Vérification de la supply fixe et de l'absence d'inflation.
- Revue sécurité avant toute exposition RPC ou bridge.
