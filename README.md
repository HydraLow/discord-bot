# Bot Discord en Go

Un bot Discord moderne et performant écrit en Go, utilisant la bibliothèque discordgo.

## Fonctionnalités

- Commandes de base (!ping, !help)
- Structure modulaire pour faciliter l'ajout de nouvelles fonctionnalités
- Performances optimales grâce à Go

## Prérequis

- Go 1.21 ou supérieur
- Un token de bot Discord

## Installation

1. Clonez ce dépôt
2. Installez les dépendances :
```bash
go mod download
```

3. Créez un bot sur le [Portail Développeur Discord](https://discord.com/developers/applications)
4. Copiez le token de votre bot
5. Définissez la variable d'environnement DISCORD_TOKEN :
```bash
# Windows (PowerShell)
$env:DISCORD_TOKEN="votre-token-ici"

# Linux/MacOS
export DISCORD_TOKEN="votre-token-ici"
```

## Lancement

```bash
go run main.go
```

## Hébergement gratuit

Pour héberger gratuitement votre bot, vous pouvez utiliser :

1. **Railway.app** - Offre un plan gratuit avec 500 heures par mois
2. **Render.com** - Offre un plan gratuit avec des limitations
3. **Oracle Cloud Free Tier** - Offre un VPS gratuit à vie

## Commandes disponibles

- `!ping` - Vérifie si le bot est en ligne
- `!help` - Affiche la liste des commandes disponibles

## Structure du projet

```
.
├── main.go          # Point d'entrée principal
├── go.mod           # Gestion des dépendances
└── README.md        # Documentation
```

## Contribution

Les contributions sont les bienvenues ! N'hésitez pas à ouvrir une issue ou une pull request. 