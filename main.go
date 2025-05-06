package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
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
			"!rps [pierre/papier/ciseaux] - Jouer à Pierre, Papier, Ciseaux"
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