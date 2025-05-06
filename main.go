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
	// ID du rôle Owner
	OwnerRoleID = "1234567890" // Remplacez par l'ID réel du rôle Owner de votre serveur
)

func init() {
	// Initialiser le générateur de nombres aléatoires
	rand.Seed(time.Now().UnixNano())
	
	// Essayer d'abord de récupérer le token depuis les variables d'environnement
	Token = os.Getenv("DISCORD_TOKEN")
	
	// Si le token n'est pas dans les variables d'environnement, le demander à l'utilisateur
	if Token == "" {
		fmt.Println("Token Discord non trouvé dans les variables d'environnement.")
		fmt.Println("Veuillez entrer votre token Discord (il ne sera pas affiché) :")
		
		reader := bufio.NewReader(os.Stdin)
		Token, _ = reader.ReadString('\n')
		Token = strings.TrimSpace(Token)
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

	// Activer les intents nécessaires
	dg.Identify.Intents = discordgo.IntentsAllWithoutPrivileged | discordgo.IntentMessageContent

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
		helpMessage := "**Commandes disponibles:**\n" +
			"!ping - Vérifier si le bot est en ligne\n" +
			"!help - Afficher ce message d'aide\n" +
			"!rps [pierre/papier/ciseaux] - Jouer à Pierre, Papier, Ciseaux\n" +
			"!kick @utilisateur ou !kick ID [raison] - Expulser un utilisateur (Rôle Owner uniquement)\n" +
			"!ban @utilisateur ou !ban ID [raison] - Bannir définitivement un utilisateur (Rôle Owner uniquement)\n" +
			"!tempban @utilisateur ou !tempban ID durée [raison] - Bannir temporairement un utilisateur (Rôle Owner uniquement)\n" +
			"   Durée format: 1h, 1d, 1w, 1m (h=heure, d=jour, w=semaine, m=mois)"
		_, err := s.ChannelMessageSend(m.ChannelID, helpMessage)
		if err != nil {
			fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
		}
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
		// Récupérer les informations du membre
		member, err := s.GuildMember(m.GuildID, m.Author.ID)
		if err != nil {
			fmt.Printf("Erreur lors de la récupération du membre: %v\n", err)
			return
		}

		// Vérifier si l'utilisateur a le rôle Owner
		hasOwnerRole := false
		for _, roleID := range member.Roles {
			role, err := s.State.Role(m.GuildID, roleID)
			if err != nil {
				continue
			}
			if role.Name == "👑Owner" {
				hasOwnerRole = true
				break
			}
		}

		if !hasOwnerRole {
			_, err := s.ChannelMessageSend(m.ChannelID, "❌ Seul le rôle 👑Owner peut utiliser cette commande!")
			if err != nil {
				fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
			}
			return
		}

		// Analyser la commande
		parts := strings.Fields(m.Content)
		if len(parts) < 2 {
			_, err := s.ChannelMessageSend(m.ChannelID, "Usage: !kick @utilisateur ou !kick ID [raison]")
			if err != nil {
				fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
			}
			return
		}

		// Extraire l'ID de l'utilisateur à expulser
		targetID := parts[1]
		// Si c'est une mention, nettoyer l'ID
		if strings.HasPrefix(targetID, "<@") && strings.HasSuffix(targetID, ">") {
			targetID = strings.Trim(targetID, "<@!>")
		}
		
		// Extraire la raison (optionnelle)
		reason := "Aucune raison fournie"
		if len(parts) > 2 {
			reason = strings.Join(parts[2:], " ")
		}

		// Expulser l'utilisateur
		err = s.GuildMemberDeleteWithReason(m.GuildID, targetID, reason)
		if err != nil {
			errorMsg := fmt.Sprintf("❌ Erreur lors de l'expulsion: %v", err)
			_, err := s.ChannelMessageSend(m.ChannelID, errorMsg)
			if err != nil {
				fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
			}
			return
		}

		// Confirmer l'expulsion
		successMsg := fmt.Sprintf("✅ Utilisateur expulsé avec succès!\nRaison: %s", reason)
		_, err = s.ChannelMessageSend(m.ChannelID, successMsg)
		if err != nil {
			fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
		}
	}

	// Commande Ban
	if strings.HasPrefix(m.Content, "!ban") {
		// Vérifier les permissions
		if !hasOwnerRole(s, m.GuildID, m.Author.ID) {
			_, err := s.ChannelMessageSend(m.ChannelID, "❌ Seul le rôle 👑Owner peut utiliser cette commande!")
			if err != nil {
				fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
			}
			return
		}

		// Analyser la commande
		parts := strings.Fields(m.Content)
		if len(parts) < 2 {
			_, err := s.ChannelMessageSend(m.ChannelID, "Usage: !ban @utilisateur ou !ban ID [raison]")
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
			errorMsg := fmt.Sprintf("❌ Erreur lors du bannissement: %v", err)
			_, err := s.ChannelMessageSend(m.ChannelID, errorMsg)
			if err != nil {
				fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
			}
			return
		}

		// Confirmer le bannissement
		successMsg := fmt.Sprintf("✅ Utilisateur banni définitivement!\nRaison: %s", reason)
		_, err = s.ChannelMessageSend(m.ChannelID, successMsg)
		if err != nil {
			fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
		}
	}

	// Commande Tempban
	if strings.HasPrefix(m.Content, "!tempban") {
		// Vérifier les permissions
		if !hasOwnerRole(s, m.GuildID, m.Author.ID) {
			_, err := s.ChannelMessageSend(m.ChannelID, "❌ Seul le rôle 👑Owner peut utiliser cette commande!")
			if err != nil {
				fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
			}
			return
		}

		// Analyser la commande
		parts := strings.Fields(m.Content)
		if len(parts) < 3 {
			_, err := s.ChannelMessageSend(m.ChannelID, "Usage: !tempban @utilisateur ou !tempban ID durée [raison]\nDurée format: 1h, 1d, 1w, 1m (h=heure, d=jour, w=semaine, m=mois)")
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
			_, err := s.ChannelMessageSend(m.ChannelID, "❌ Format de durée invalide! Utilisez: 1h, 1d, 1w, 1m")
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
			errorMsg := fmt.Sprintf("❌ Erreur lors du bannissement temporaire: %v", err)
			_, err := s.ChannelMessageSend(m.ChannelID, errorMsg)
			if err != nil {
				fmt.Printf("Erreur lors de l'envoi du message: %v\n", err)
			}
			return
		}

		// Confirmer le bannissement temporaire
		successMsg := fmt.Sprintf("✅ Utilisateur banni temporairement pour %s!\nRaison: %s", duration, reason)
		_, err = s.ChannelMessageSend(m.ChannelID, successMsg)
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
		if roleID == OwnerRoleID {
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