# Identité des réseaux Xitcoin

## Identité publique

- Nom : Xitcoin
- Symbole : XTC
- Précision : 18 décimales
- Unité technique interne : xtc

## Famille des Chain ID EVM

| Réseau | Chain ID | Statut |
|---|---:|---|
| Xitcoin | 101088 | Identité historique déjà enregistrée |
| Xitcoin Testnet | 101089 | Réservé en interne pour le futur testnet public |
| Xitcoin Devnet | 101090 | Réservé en interne uniquement si nécessaire |
| Xitcoin Private Testnet | 20260807 | Test privé local ; jamais publié dans ChainList |

## Règles

- Le nom Xitcoin et le symbole XTC restent identiques sur tous les réseaux.
- Le testnet public est séparé du réseau public afin que les utilisateurs ne confondent jamais leurs actifs.
- La fiche historique ChainList 101088 ne doit pointer que vers le futur réseau public Xitcoin.
- Le domaine historique network.xitcoin.org ne doit pas être redirigé vers le testnet privé.
- La disponibilité de 101089 et 101090 doit être recontrôlée juste avant toute publication ChainList.
