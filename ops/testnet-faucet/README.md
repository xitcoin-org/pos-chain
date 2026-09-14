# Xitcoin Testnet Faucet

Faucet public officiel du réseau `xitcoin-testnet-v2-1`.

## Paramètres

- Distribution : 10 XTC par demande
- Dénomination technique : axtc
- Limite par adresse : 1 demande par 24 heures
- Limite par IP : 3 demandes par 24 heures
- Backend : Node.js derrière Nginx et Cloudflare

## Candidat de récupération — revue requise, aucun déploiement

Cette branche change le chemin de demande ; elle n'est pas une installation
reçue. Exécuter uniquement `node --test ops/testnet-faucet/recovery.test.js`
pour la validation fictive, sans serveur, RPC, CLI, signature ou transaction.
Voir [RECOVERY-DESIGN.md](RECOVERY-DESIGN.md), suivi #33.

Le journal `state/recovery-v1` réserve durablement les quotas avant soumission.
Un hash reste `submitted`, sans prétendre confirmer l'exécution ; toute autre
issue ambiguë reste `unknown`. Ces demandes, ainsi que les claims historiques
importés sans preuve d'inclusion, bloquent la réémission même après 24 heures.
HTTP 202 fournit une référence de réconciliation et l'interface demande de ne
pas resoumettre. Le noyau de réconciliation vérifie hash, chaîne, expéditeur, destinataire,
montant, dénomination, hauteur et code via un adaptateur de lecture injectable.
Il est testé uniquement sur reçus fictifs ; aucun adaptateur RPC ni endpoint
public de résolution n'est branché. Sans hash, l'état demeure inchangé pour
réception opérateur indépendante. Un reçu inclus à code non nul garde le quota
jusqu'à expiration de la fenêtre ; l'absence de reçu maintient le blocage.

`claims.json` est conservé octet pour octet. Un historique invalide est refusé.
Le parent `state/` doit déjà exister. Son entrée `recovery-v1` est synchronisée
avant ouverture ; un échec de fsync du parent empêche toute soumission.
Chaque écriture fsync ensuite le fichier puis le répertoire ; les copies de récupération
sont conservées. Le plafond de 64 snapshots de 8 Mio impose un arrêt sûr et une
réception d'archivage avant une exploitation prolongée. Aucun effacement ou
recyclage automatique des réservations n'est effectué. Un crash peut laisser
`writer.lock` : ne pas le supprimer automatiquement, établir d'abord qu'aucun
processus ne peut encore soumettre et réconcilier les demandes pendantes.
L'arrêt normal ferme le serveur puis libère uniquement son verrou.

`TRUSTED_PROXY_IPS` est vide par défaut : l'adresse du pair est alors utilisée.
Configurer uniquement les adresses exactes de proxys reçus indépendamment ;
les en-têtes transmis par un autre pair sont ignorés. Ce réglage ne constitue
pas une preuve de la configuration Nginx chargée ou de la confiance Cloudflare.
La revue de la source installée, de la rotation et de la réconciliation reste
requise avant toute autorisation distincte de mise en service.

## Référence de déploiement historique (non autorisée par cette branche)

1. Copier `faucet.env.example` vers `/etc/xitcoin-testnet/faucet.env`.
2. Adapter les chemins et le RPC sans ajouter de secret au dépôt.
3. Installer le service systemd et la configuration Nginx.
4. Importer séparément la clé `faucet-official` dans le keyring sécurisé.

Les fichiers de clé, mots de passe, états et portefeuilles ne doivent jamais être commités.
