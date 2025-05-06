package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
	// Nom du rôle Owner
	OwnerRoleName = "👑Owner"
	// Nom du rôle de vérification
	VerifiedRoleName = "✅Verified"
	// Noms des rôles auto-attribués
	FortniteRoleName = "FORTNITE PLAYER"
)

func init() {
	// Initialiser le générateur de nombres aléatoires
	rand.Seed(time.Now().UnixNano())
	
	// Lire le token depuis le fichier token.env
	tokenBytes, err := os.ReadFile("token.env")
	if err == nil {
		Token = strings.TrimSpace(string(tokenBytes))
		// Si le token commence par "DISCORD_TOKEN=", l'enlever
		if strings.HasPrefix(Token, "DISCORD_TOKEN=") {
			Token = strings.TrimPrefix(Token, "DISCORD_TOKEN=")
		}
	}
	
	// Si le token n'est pas dans le fichier, essayer les variables d'environnement
	if Token == "" {
		Token = os.Getenv("DISCORD_TOKEN")
	}
	
	// Si le token n'est toujours pas trouvé, le demander à l'utilisateur
	if Token == "" {
		fmt.Println("Token Discord non trouvé dans les variables d'environnement.")
		fmt.Println("Veuillez entrer votre token Discord (il ne sera pas affiché) :")
		
		reader := bufio.NewReader(os.Stdin)
		Token, _ = reader.ReadString('\n')
		Token = strings.TrimSpace(Token)
	}

	// Vérifier si le token a le bon format
	if !strings.HasPrefix(Token, "MT") && !strings.HasPrefix(Token, "NT") {
		fmt.Println("❌ Le token fourni ne semble pas être un token Discord valide.")
		fmt.Println("Un token Discord commence généralement par 'MT' ou 'NT'.")
		fmt.Println("Veuillez vérifier votre token dans le portail développeur Discord.")
		os.Exit(1)
	}

	if Token == "" {
		log.Fatal("Token Discord requis pour démarrer le bot")
	}
}

func main() {
	// Créer une nouvelle session Discord
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal("Error creating Discord session: ", err)
	}

	// Ajouter les handlers d'événements
	dg.AddHandler(messageCreate)
	dg.AddHandler(ready)
	dg.AddHandler(interactionCreate) // Ajouter le handler pour les interactions (boutons)

	// Activer les intents nécessaires
	dg.Identify.Intents = discordgo.IntentsAllWithoutPrivileged | discordgo.IntentMessageContent | discordgo.IntentGuildMembers

	// Ouvrir la connexion
	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening connection: ", err)
	}

	// Attendre jusqu'à ce que CTRL-C ou un autre signal d'arrêt soit reçu
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	// Nettoyer proprement
	dg.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateGameStatus(0, "!help pour les commandes")
	fmt.Printf("Logged in as: %v#%v\n", event.User.Username, event.User.Discriminator)
	fmt.Println("Bot is ready to receive commands!")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Log pour le débogage
	fmt.Printf("Message reçu: %s de %s\n", m.Content, m.Author.Username)

	// Ignorer les messages du bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Commande simple de test
	if m.Content == "!ping" {
		fmt.Println("Commande !ping détectée")
		_, err := s.ChannelMessageSend(m.ChannelID, "Pong! 🏓")
		if err != nil {
			fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
		}
	}

	// Commande d'aide
	if m.Content == "!help" {
		fmt.Println("Commande !help détectée")
		
		embed := &discordgo.MessageEmbed{
			Title:       "📚 Commandes disponibles",
			Description: "Voici la liste des commandes du bot :",
			Color:       0x00ff00, // Vert
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "🎮 Commandes de base",
					Value:  "`!ping` - Vérifier si le bot est en ligne\n`!help` - Afficher ce message d'aide",
					Inline: false,
				},
				{
					Name:   "🎲 Mini-jeux",
					Value:  "`!rps [pierre/papier/ciseaux]` - Jouer à Pierre, Papier, Ciseaux",
					Inline: false,
				},
				{
					Name:   "🛡️ Modération",
					Value:  "`!kick @utilisateur [raison]` - Expulser un utilisateur\n`!ban @utilisateur [raison]` - Bannir définitivement\n`!unban @utilisateur` - Débannir un utilisateur\n`!tempban @utilisateur durée [raison]` - Bannir temporairement",
					Inline: false,
				},
				{
					Name:   "⏱️ Format de durée",
					Value:  "`h` - Heures (ex: 1h)\n`d` - Jours (ex: 1d)\n`w` - Semaines (ex: 1w)\n`m` - Mois (ex: 1m)",
					Inline: false,
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Rôle 👑Owner requis pour les commandes de modération",
			},
		}
		
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
		if err != nil {
			fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
		}
		return
	}

	// Commande Pierre, Papier, Ciseaux
	if strings.HasPrefix(m.Content, "!rps") {
		parts := strings.Fields(m.Content)
		if len(parts) != 2 {
			_, err := s.ChannelMessageSend(m.ChannelID, "Usage: !rps [pierre/papier/ciseaux]")
			if err != nil {
				fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
			}
			return
		}

		playerChoice := strings.ToLower(parts[1])
		if playerChoice != "pierre" && playerChoice != "papier" && playerChoice != "ciseaux" {
			_, err := s.ChannelMessageSend(m.ChannelID, "Choix invalide! Utilisez: pierre, papier ou ciseaux")
			if err != nil {
				fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
			}
			return
		}

		// Choix du bot
		choices := []string{"pierre", "papier", "ciseaux"}
		botChoice := choices[rand.Intn(3)]

		// Déterminer le gagnant
		result := determineWinner(playerChoice, botChoice)

		// Créer le message de résultat
		resultMessage := fmt.Sprintf("**Pierre, Papier, Ciseaux!**\n"+
			"Vous avez choisi: %s\n"+
			"J'ai choisi: %s\n"+
			"Résultat: %s", 
			playerChoice, botChoice, result)

		_, err := s.ChannelMessageSend(m.ChannelID, resultMessage)
		if err != nil {
			fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
		}
	}

	// Commande Kick
	if strings.HasPrefix(m.Content, "!kick") {
		// Vérifier les permissions
		if !hasOwnerRole(s, m.GuildID, m.Author.ID) {
			embed := &discordgo.MessageEmbed{
				Title:       "❌ Permission refusée",
				Description: "Seul le rôle 👑Owner peut utiliser cette commande!",
				Color:       0xff0000, // Rouge
			}
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
			if err != nil {
				fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
			}
			return
		}

		// Analyser la commande
		parts := strings.Fields(m.Content)
		if len(parts) < 2 {
			embed := &discordgo.MessageEmbed{
				Title:       "❌ Usage incorrect",
				Description: "Usage: `!kick @utilisateur [raison]` ou `!kick ID [raison]`",
				Color:       0xff0000,
			}
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
			if err != nil {
				fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
			}
			return
		}

		// Extraire l'ID de l'utilisateur à expulser
		targetID := parts[1]
		if strings.HasPrefix(targetID, "<@") && strings.HasSuffix(targetID, ">") {
			targetID = strings.Trim(targetID, "<@!>")
		}

		// Extraire la raison (optionnelle)
		reason := "Aucune raison fournie"
		if len(parts) > 2 {
			reason = strings.Join(parts[2:], " ")
		}

		// Expulser l'utilisateur
		err := s.GuildMemberDeleteWithReason(m.GuildID, targetID, reason)
		if err != nil {
			embed := &discordgo.MessageEmbed{
				Title:       "❌ Erreur",
				Description: fmt.Sprintf("Erreur lors de l'expulsion: %v", err),
				Color:       0xff0000,
			}
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
			if err != nil {
				fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
			}
			return
		}

		// Confirmer l'expulsion
		embed := &discordgo.MessageEmbed{
			Title:       "✅ Utilisateur expulsé",
			Description: fmt.Sprintf("L'utilisateur <@%s> a été expulsé avec succès!", targetID),
			Color:       0x00ff00,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Raison",
					Value:  reason,
					Inline: false,
				},
			},
		}
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
		if err != nil {
			fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
		}
	}

	// Commande Ban
	if strings.HasPrefix(m.Content, "!ban") {
		// Vérifier les permissions
		if !hasOwnerRole(s, m.GuildID, m.Author.ID) {
			embed := &discordgo.MessageEmbed{
				Title:       "❌ Permission refusée",
				Description: "Seul le rôle 👑Owner peut utiliser cette commande!",
				Color:       0xff0000,
			}
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
			if err != nil {
				fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
			}
			return
		}

		// Analyser la commande
		parts := strings.Fields(m.Content)
		if len(parts) < 2 {
			embed := &discordgo.MessageEmbed{
				Title:       "❌ Usage incorrect",
				Description: "Usage: `!ban @utilisateur [raison]` ou `!ban ID [raison]`",
				Color:       0xff0000,
			}
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
			if err != nil {
				fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
			}
			return
		}

		// Extraire l'ID de l'utilisateur à bannir
		targetID := parts[1]
		if strings.HasPrefix(targetID, "<@") && strings.HasSuffix(targetID, ">") {
			targetID = strings.Trim(targetID, "<@!>")
		}

		// Extraire la raison (optionnelle)
		reason := "Aucune raison fournie"
		if len(parts) > 2 {
			reason = strings.Join(parts[2:], " ")
		}

		// Bannir l'utilisateur
		err := s.GuildBanCreateWithReason(m.GuildID, targetID, reason, 0)
		if err != nil {
			embed := &discordgo.MessageEmbed{
				Title:       "❌ Erreur",
				Description: fmt.Sprintf("Erreur lors du bannissement: %v", err),
				Color:       0xff0000,
			}
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
			if err != nil {
				fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
			}
			return
		}

		// Confirmer le bannissement
		embed := &discordgo.MessageEmbed{
			Title:       "✅ Utilisateur banni",
			Description: fmt.Sprintf("L'utilisateur <@%s> a été banni définitivement!", targetID),
			Color:       0x00ff00,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Raison",
					Value:  reason,
					Inline: false,
				},
			},
		}
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
		if err != nil {
			fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
		}
	}

	// Commande Unban
	if strings.HasPrefix(m.Content, "!unban") {
		// Vérifier les permissions
		if !hasOwnerRole(s, m.GuildID, m.Author.ID) {
			embed := &discordgo.MessageEmbed{
				Title:       "❌ Permission refusée",
				Description: "Seul le rôle 👑Owner peut utiliser cette commande!",
				Color:       0xff0000,
			}
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
			if err != nil {
				fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
			}
			return
		}

		// Analyser la commande
		parts := strings.Fields(m.Content)
		if len(parts) != 2 {
			embed := &discordgo.MessageEmbed{
				Title:       "❌ Usage incorrect",
				Description: "Usage: `!unban @utilisateur` ou `!unban ID`",
				Color:       0xff0000,
			}
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
			if err != nil {
				fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
			}
			return
		}

		// Extraire l'ID de l'utilisateur à débannir
		targetID := parts[1]
		if strings.HasPrefix(targetID, "<@") && strings.HasSuffix(targetID, ">") {
			targetID = strings.Trim(targetID, "<@!>")
		}

		// Vérifier si l'ID est valide
		if _, err := strconv.ParseInt(targetID, 10, 64); err != nil {
			embed := &discordgo.MessageEmbed{
				Title:       "❌ ID invalide",
				Description: "ID d'utilisateur invalide. Utilisez un ID valide ou une mention.",
				Color:       0xff0000,
			}
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
			if err != nil {
				fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
			}
			return
		}

		// Débannir l'utilisateur
		err := s.GuildBanDelete(m.GuildID, targetID)
		if err != nil {
			embed := &discordgo.MessageEmbed{
				Title:       "❌ Erreur",
				Description: fmt.Sprintf("Erreur lors du débannissement: %v", err),
				Color:       0xff0000,
			}
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
			if err != nil {
				fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
			}
			return
		}

		// Confirmer le débannissement
		embed := &discordgo.MessageEmbed{
			Title:       "✅ Utilisateur débanni",
			Description: fmt.Sprintf("L'utilisateur <@%s> a été débanni avec succès!", targetID),
			Color:       0x00ff00,
		}
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
		if err != nil {
			fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
		}
	}

	// Commande Tempban
	if strings.HasPrefix(m.Content, "!tempban") {
		// Vérifier les permissions
		if !hasOwnerRole(s, m.GuildID, m.Author.ID) {
			embed := &discordgo.MessageEmbed{
				Title:       "❌ Permission refusée",
				Description: "Seul le rôle 👑Owner peut utiliser cette commande!",
				Color:       0xff0000,
			}
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
			if err != nil {
				fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
			}
			return
		}

		// Analyser la commande
		parts := strings.Fields(m.Content)
		if len(parts) < 3 {
			embed := &discordgo.MessageEmbed{
				Title:       "❌ Usage incorrect",
				Description: "Usage: `!tempban @utilisateur durée [raison]` ou `!tempban ID durée [raison]`\nDurée format: 1h, 1d, 1w, 1m (h=heure, d=jour, w=semaine, m=mois)",
				Color:       0xff0000,
			}
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
			if err != nil {
				fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
			}
			return
		}

		// Extraire l'ID de l'utilisateur à bannir
		targetID := parts[1]
		if strings.HasPrefix(targetID, "<@") && strings.HasSuffix(targetID, ">") {
			targetID = strings.Trim(targetID, "<@!>")
		}

		// Extraire la durée
		duration := parts[2]
		banDuration, err := parseDuration(duration)
		if err != nil {
			embed := &discordgo.MessageEmbed{
				Title:       "❌ Durée invalide",
				Description: "Format de durée invalide! Utilisez: 1h, 1d, 1w, 1m",
				Color:       0xff0000,
			}
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
			if err != nil {
				fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
			}
			return
		}

		// Extraire la raison (optionnelle)
		reason := "Aucune raison fournie"
		if len(parts) > 3 {
			reason = strings.Join(parts[3:], " ")
		}

		// Bannir l'utilisateur
		err = s.GuildBanCreateWithReason(m.GuildID, targetID, fmt.Sprintf("%s (Tempban: %s)", reason, duration), 0)
		if err != nil {
			embed := &discordgo.MessageEmbed{
				Title:       "❌ Erreur",
				Description: fmt.Sprintf("Erreur lors du bannissement temporaire: %v", err),
				Color:       0xff0000,
			}
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
			if err != nil {
				fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
			}
			return
		}

		// Confirmer le bannissement temporaire
		embed := &discordgo.MessageEmbed{
			Title:       "✅ Utilisateur banni temporairement",
			Description: fmt.Sprintf("L'utilisateur <@%s> a été banni pour %s!", targetID, duration),
			Color:       0x00ff00,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Raison",
					Value:  reason,
					Inline: false,
				},
			},
		}
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
		if err != nil {
			fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
		}

		// Programmer le débannissement
		go func() {
			time.Sleep(banDuration)
			err := s.GuildBanDelete(m.GuildID, targetID)
			if err != nil {
				fmt.Printf("Erreur lors du débannissement automatique: %v\n", err)
			}
		}()
	}

	// Commande pour créer le message de vérification
	if m.Content == "!setupverify" {
		// Vérifier les permissions
		if !hasOwnerRole(s, m.GuildID, m.Author.ID) {
			embed := &discordgo.MessageEmbed{
				Title:       "❌ Permission refusée",
				Description: "Seul le rôle 👑Owner peut utiliser cette commande!",
				Color:       0xff0000,
			}
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
			if err != nil {
				fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
			}
			return
		}

		// Créer le message de vérification
		embed := &discordgo.MessageEmbed{
			Title:       "🔒 Vérification",
			Description: "Cliquez sur le bouton ci-dessous pour accéder au serveur.",
			Color:       0x00ff00,
		}

		// Créer le bouton de vérification
		row := discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Vérifier",
					Style:    discordgo.PrimaryButton,
					CustomID: "verify_button",
				},
			},
		}

		// Envoyer le message avec le bouton
		_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Embed:      embed,
			Components: []discordgo.MessageComponent{row},
		})
		if err != nil {
			fmt.Printf("Erreur lors de l'envoi du message de vérification: %v\n", err)
			return
		}

		// Supprimer la commande
		err = s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			fmt.Printf("Erreur lors de la suppression de la commande: %v\n", err)
		}
	}
}

func determineWinner(player, bot string) string {
	if player == bot {
		return "Égalité! 🤝"
	}

	switch {
	case player == "pierre" && bot == "ciseaux",
		player == "papier" && bot == "pierre",
		player == "ciseaux" && bot == "papier":
		return "Vous avez gagné! 🎉"
	default:
		return "J'ai gagné! 😎"
	}
}

// Fonction utilitaire pour vérifier le rôle Owner
func hasOwnerRole(s *discordgo.Session, guildID, userID string) bool {
	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		fmt.Printf("Erreur lors de la récupération du membre: %v\n", err)
		return false
	}

	// Vérifier si l'utilisateur a le rôle Owner
	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			continue
		}
		if role.Name == OwnerRoleName {
			return true
		}
	}
	return false
}

// Fonction pour parser la durée du tempban
func parseDuration(duration string) (time.Duration, error) {
	if len(duration) < 2 {
		return 0, fmt.Errorf("durée invalide")
	}

	value := duration[:len(duration)-1]
	unit := duration[len(duration)-1:]

	num, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	switch unit {
	case "h":
		return time.Duration(num) * time.Hour, nil
	case "d":
		return time.Duration(num) * 24 * time.Hour, nil
	case "w":
		return time.Duration(num) * 7 * 24 * time.Hour, nil
	case "m":
		return time.Duration(num) * 30 * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unité de temps invalide")
	}
}

// Fonction pour gérer les interactions (boutons)
func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionMessageComponent {
		// Vérifier si c'est le bouton de vérification
		if i.MessageComponentData().CustomID == "verify_button" {
			// Récupérer tous les rôles du serveur
			roles, err := s.GuildRoles(i.GuildID)
			if err != nil {
				fmt.Printf("Erreur lors de la récupération des rôles: %v\n", err)
				return
			}

			// Trouver l'ID du rôle de vérification
			var verifiedRoleID string
			for _, role := range roles {
				if role.Name == VerifiedRoleName {
					verifiedRoleID = role.ID
					break
				}
			}

			if verifiedRoleID != "" {
				// Ajouter le rôle à l'utilisateur
				params := &discordgo.GuildMemberParams{
					Roles: &[]string{verifiedRoleID},
				}
				_, err = s.GuildMemberEdit(i.GuildID, i.Member.User.ID, params)
				if err != nil {
					fmt.Printf("Erreur lors de l'attribution du rôle: %v\n", err)
					return
				}

				// Répondre à l'interaction
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "✅ Vous avez été vérifié avec succès!",
						Flags:   discordgo.MessageFlagsEphemeral, // Message visible uniquement par l'utilisateur
					},
				})
				if err != nil {
					fmt.Printf("Erreur lors de la réponse à l'interaction: %v\n", err)
				}
			}
		}
	}
} 