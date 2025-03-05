package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fullcycle-auction_go/configuration/database/mongodb"
	"fullcycle-auction_go/internal/infra/api/web/controller/auction_controller"
	"fullcycle-auction_go/internal/infra/api/web/controller/bid_controller"
	"fullcycle-auction_go/internal/infra/api/web/controller/user_controller"
	"fullcycle-auction_go/internal/infra/database/auction"
	"fullcycle-auction_go/internal/infra/database/bid"
	"fullcycle-auction_go/internal/infra/database/user"
	"fullcycle-auction_go/internal/usecase/auction_usecase"
	"fullcycle-auction_go/internal/usecase/bid_usecase"
	"fullcycle-auction_go/internal/usecase/user_usecase"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

var server *httptest.Server

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Auction struct {
	ID string `json:"id"`
}

type Bid struct {
	AuctionID string  `json:"auction_id"`
	UserID    string  `json:"user_id"`
	Amount    float64 `json:"amount"`
}

type AuctionResponse struct {
	ID          string `json:"id"`
	ProductName string `json:"product_name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Condition   int    `json:"condition"`
	Status      int    `json:"status"`
	Timestamp   string `json:"timestamp"`
}

type AuctionWinnerResponse struct {
	Auction struct {
		ID          string `json:"id"`
		ProductName string `json:"product_name"`
		Category    string `json:"category"`
		Description string `json:"description"`
		Condition   int    `json:"condition"`
		Status      int    `json:"status"`
		Timestamp   string `json:"timestamp"`
	} `json:"auction"`
	Bid struct {
		ID        string  `json:"id"`
		UserID    string  `json:"user_id"`
		AuctionID string  `json:"auction_id"`
		Amount    float64 `json:"amount"`
		Timestamp string  `json:"timestamp"`
	} `json:"bid"`
}

type AuctionStatusResponse struct {
	ID          string `json:"id"`
	ProductName string `json:"product_name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Condition   int    `json:"condition"`
	Status      int    `json:"status"`
	Timestamp   string `json:"timestamp"`
}

func TestMain(m *testing.M) {
	// Carregar variáveis de ambiente
	if err := godotenv.Load("../../cmd/auction/.env"); err != nil {
		log.Fatal("Erro ao carregar variáveis de ambiente:", err)
	}

	// Criar contexto e conectar ao banco
	ctx := context.Background()
	databaseConnection, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatal("Erro ao conectar ao banco:", err)
	}

	// Inicializar servidor e dependências
	router := setupRouter(databaseConnection)

	// Criar um servidor de teste
	server = httptest.NewServer(router)

	// Aguardar o servidor iniciar
	waitForServer()

	// Rodar os testes
	code := m.Run()

	// Encerrar o servidor ao final dos testes
	server.Close()

	// Sair com o código de status apropriado
	os.Exit(code)
}

// setupRouter configura as rotas da aplicação
func setupRouter(database *mongo.Database) *gin.Engine {
	router := gin.Default()

	userController, bidController, auctionController := initDependencies(database)

	router.GET("/auction", auctionController.FindAuctions)
	router.GET("/auction/:auctionId", auctionController.FindAuctionById)
	router.POST("/auction", auctionController.CreateAuction)
	router.GET("/auction/winner/:auctionId", auctionController.FindWinningBidByAuctionId)
	router.POST("/bid", bidController.CreateBid)
	router.GET("/bid/:auctionId", bidController.FindBidByAuctionId)
	router.GET("/user/:userId", userController.FindUserById)
	router.POST("/user", userController.CreateUser)

	return router
}

// initDependencies inicializa os controllers e repositórios
func initDependencies(database *mongo.Database) (
	userController *user_controller.UserController,
	bidController *bid_controller.BidController,
	auctionController *auction_controller.AuctionController) {

	auctionRepository := auction.NewAuctionRepository(database)
	bidRepository := bid.NewBidRepository(database, auctionRepository)
	userRepository := user.NewUserRepository(database)

	userController = user_controller.NewUserController(
		user_usecase.NewUserUseCase(userRepository))
	auctionController = auction_controller.NewAuctionController(
		auction_usecase.NewAuctionUseCase(auctionRepository, bidRepository))
	bidController = bid_controller.NewBidController(bid_usecase.NewBidUseCase(bidRepository))

	return
}

// waitForServer espera o servidor responder antes de continuar
func waitForServer() {
	for i := 0; i < 10; i++ { // Aumenta o número de tentativas
		resp, err := http.Get(server.URL + "/auction?status=0")
		if err == nil {
			if resp.StatusCode == http.StatusOK {
				resp.Body.Close()
				log.Println("Servidor respondeu com sucesso.")
				return
			} else {
				log.Printf("Servidor respondeu com status %d. Tentativa %d/10.\n", resp.StatusCode, i+1)
				resp.Body.Close()
			}
		} else {
			log.Printf("Erro ao conectar ao servidor: %v. Tentativa %d/10.\n", err, i+1)
		}
		time.Sleep(2000 * time.Millisecond) // Aumenta o tempo de espera
	}
	log.Fatal("Servidor não respondeu a tempo")
}

// postJSON realiza requisições POST no servidor de teste
func postJSON(url string, payload any) (*http.Response, error) {
	body, _ := json.Marshal(payload)
	return http.Post(server.URL+url, "application/json", bytes.NewBuffer(body))
}

// getAuctionWinner obtém o vencedor do leilão
func getAuctionWinner(auctionID string) string {
	resp, err := http.Get(server.URL + "/auction/winner/" + auctionID)
	if err != nil {
		log.Printf("Erro ao obter vencedor do leilão: %v", err)
		return ""
	}
	defer resp.Body.Close()

	var data AuctionWinnerResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Printf("Erro ao decodificar vencedor do leilão: %v", err)
		return ""
	}

	log.Printf("Vencedor do leilão %s: %s", auctionID, data.Bid.UserID) // log do vencedor
	return data.Bid.UserID
}

// getAuctionStatus obtém o status do leilão
func getAuctionStatus(auctionID string) string {
	resp, err := http.Get(server.URL + "/auction/" + auctionID)
	if err != nil {
		log.Printf("Erro ao obter status do leilão: %v", err)
		return ""
	}
	defer resp.Body.Close()

	var data AuctionStatusResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Printf("Erro ao decodificar status do leilão: %v", err)
		return ""
	}

	status := "open"
	if data.Status == 1 {
		status = "closed"
	}

	log.Printf("Status do leilão %s: %s", auctionID, status) // log do status
	return status
}

func TestAuctionFlow(t *testing.T) {
	// Criar usuários
	user1 := createUser(t, "User1")
	user2 := createUser(t, "User2")

	// Criar leilão
	auction := createAuction(t)

	// Criar 4 lances (2 por usuário)
	placeBid(t, auction.ID, user1.ID, 100.00)
	placeBid(t, auction.ID, user2.ID, 110.00)
	placeBid(t, auction.ID, user1.ID, 120.00)
	placeBid(t, auction.ID, user2.ID, 130.00)

	// Aguardar tempo necessário
	err := waitForAuctionToClose(auction.ID, 25*time.Second)
	assert.NoError(t, err)

	// Último lance após 10s
	placeBid(t, auction.ID, user1.ID, 140.00)

	// Verificar vencedor
	winner := getAuctionWinner(auction.ID)
	assert.NotEmpty(t, winner)

	// Verificar que o leilão fechou
	status := getAuctionStatus(auction.ID)
	assert.Equal(t, "closed", status)
}

func createUser(t *testing.T, name string) User {
	resp, err := postJSON("/user", map[string]string{"name": name})
	assert.NoError(t, err)
	defer resp.Body.Close()
	var user User
	json.NewDecoder(resp.Body).Decode(&user)
	return user
}

func createAuction(t *testing.T) Auction {
	resp, err := postJSON("/auction", map[string]interface{}{
		"product_name": "Product example",
		"category":     "examples",
		"description":  "Example product",
		"condition":    0,
	})
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	resp, err = http.Get(server.URL + "/auction?status=0")
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var auctions []AuctionResponse
	err = json.NewDecoder(resp.Body).Decode(&auctions)
	assert.NoError(t, err)

	assert.NotEmpty(t, auctions, "Lista de leilões vazia")
	lastAuction := auctions[len(auctions)-1]

	return Auction{ID: lastAuction.ID}
}

func placeBid(t *testing.T, auctionID, userID string, amount float64) {
	resp, err := postJSON("/bid", map[string]interface{}{
		"auction_id": auctionID,
		"user_id":    userID,
		"amount":     amount,
	})
	assert.NoError(t, err)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Erro ao criar lance: status code %d", resp.StatusCode)
	}
}

func waitForAuctionToClose(auctionID string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return errors.New("timeout esperando o leilão fechar")
		case <-ticker.C:
			status := getAuctionStatus(auctionID)

			if status == "closed" {
				return nil
			}
		}
	}
}
