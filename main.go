package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type client struct {
	conn    net.Conn
	name    string
	writer  *bufio.Writer
	done    chan struct{}
	Message chan string
}

var clients []*client

func main() {
	args := os.Args[1:]
	var port string
	//valeur du port par défaut
	
	//Création du serveur et connexion à un port
	if len(args) == 0 {
		StartTCPServer("8989")
	} else if len(args) > 1 {
		fmt.Println("[USAGE]: ./TCPChat $port")
	} else {
		port=args[0]
		StartTCPServer(port ) //On ajoute "
	}
}

// StartTCPServer démarre un serveur TCP qui écoute sur le port spécifié par la valeur.
func StartTCPServer(value string) {
	// Crée un écouteur TCP qui écoute sur le port spécifié.
	listener, err := net.Listen("tcp", ":"+value)
	if err != nil {
		// Si une erreur se produit lors de la création du serveur, affiche un message d'erreur et termine le programme.
		log.Fatal("Erreur lors de la création du serveur:", err)
	}
	// Ferme l'écouteur lorsque la fonction se termine.
	defer listener.Close()

	// Affiche un message indiquant que le serveur écoute sur le port spécifié.
	fmt.Printf("Listening on the port : %s\n", value)

	// Boucle infinie pour accepter les connexions entrantes.
	for {
		// Accepte une connexion TCP.
		conn, err := listener.Accept()
		if err != nil {
			// Si une erreur se produit lors de l'acceptation de la connexion, affiche un message d'erreur et continue avec la prochaine connexion.
			fmt.Println("Erreur lors de l'acceptation de la connexion:", err)
			continue
		}
		if len(clients) >= 10 {

			fmt.Println("Maximum number of clients reached. Rejecting new connection...")

			conn.Close()

			continue
		}
		// Gère la connexion dans un goroutine séparé.
		go HandleConnection(conn)
	}
}

func HandleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	writer.WriteString("Welcome to TCP-Chat!\n")
	writer.WriteString("         _nnnn_                 \n")
	writer.WriteString("        dGGGGMMb               \n")
	writer.WriteString("       @p~qp~~qMb              \n")
	writer.WriteString("       M|@||@) M|             \n")
	writer.WriteString("       @,----.JM|            \n")
	writer.WriteString("      JS^\\__/  qKL           \n")
	writer.WriteString("     dZP        qKRb          \n")
	writer.WriteString("    dZP          qKKb         \n")
	writer.WriteString("   fZP            SMMb        \n")
	writer.WriteString("   HZM            MMMM        \n")
	writer.WriteString("   FqM            MMMM        \n")
	writer.WriteString(" __| \".        |\\dS\"qML     \n")
	writer.WriteString(" |    `.       | `' \\Zq      \n")
	writer.WriteString("_)      \\.___.,|     .'       \n")
	writer.WriteString("\\____   )MMMMMP|   .'        \n")
	writer.WriteString("     `-'       `--'         \n")

	writer.WriteString("[ENTER YOUR NAME]: ")
	writer.Flush()

	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	message := fmt.Sprintf("%s has joined the chat!\n", name)

	if name == "" {
		conn.Close()
		return
	}

	newClient := &client{conn: conn, name: name, writer: writer, done: make(chan struct{})}

	clients = append(clients, newClient)

	go newClient.listen()

	broadcast(message)

	for {

		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		message = strings.TrimSpace(message)
		if message != "" {
			// Add a check to prevent empty inputs.

			if strings.Compare(message, "") == 0 {

				conn.Write([]byte("Error: Empty input. Please try again.\n"))

				conn.Write([]byte("> "))

				continue

			}
			conn.Write([]byte("\033[1A\033[K"))
			broadcast(fmt.Sprintf("[%s][%s]: %s\n", time.Now().Format("2006-01-02 15:04:05"), name, message))

		}
	}

}

func removeClient(clientToRemove *client) {
	for i, c := range clients {
		if c == clientToRemove {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}

func broadcast(message string) {
	for _, cl := range clients {
		cl.writer.WriteString(message)
		cl.writer.Flush()
	}
}

func (c *client) listen() {
	defer func() {
		removeClient(c)
		message := fmt.Sprintf("%s has left the chat!", c.name)
		broadcast(message)
		c.conn.Close()
		close(c.done)
	}()
	reader := bufio.NewReader(c.conn) // Créer un nouveau reader pour le client
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		message = strings.TrimSpace(message)
		if message != "" {
			broadcast(fmt.Sprintf("[%s][%s]: %s", time.Now().Format("2006-01-02 15:04:05"), c.name, message))
		}
	}
}
