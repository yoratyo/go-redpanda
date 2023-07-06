package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/joho/godotenv"
	pub "github.com/yoratyo/go-redpanda/kafka"
	"github.com/yoratyo/go-redpanda/repository"
	"github.com/yoratyo/go-redpanda/service"
	"github.com/yoratyo/go-redpanda/storage"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func runSignalListener(cancel context.CancelFunc) {
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		s := <-sign
		fmt.Printf("Received %s, canceling context", s.String())
		cancel()
	}()
}

func main() {
	os.Exit(run())
}

func initStorage(ctx context.Context) *storage.Client {
	st, err := storage.New(ctx, storage.Options{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_NAME"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		panic(err)
	}

	return st
}

var (
	validate = validator.New()
	decoder  = schema.NewDecoder()
)

type HttpServer struct {
	svc service.Service
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	log.Println(message)
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (c *HttpServer) cryptoPost(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var cryptoReq service.CryptoDTO
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&cryptoReq); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error decode payload: %s", err))
		return
	}
	defer r.Body.Close()

	err := validate.Struct(cryptoReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error validate payload: %s", err))
		return
	}
	log.Printf("cryptoPost payload: %+v\n", cryptoReq)

	err = c.svc.PublishCrypto(ctx, cryptoReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error post crypto: %s", err))
		return
	}

	log.Printf("cryptoPost success: %+v\n", cryptoReq)
	respondWithJSON(w, http.StatusOK, cryptoReq)
}

func (c *HttpServer) cryptoList(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var cryptolistReq service.ListCryptoRequest
	err := decoder.Decode(&cryptolistReq, r.URL.Query())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error decode query: %s", err))
		return
	}
	log.Printf("cryptoList query: %+v\n", cryptolistReq)

	res, err := c.svc.ListCrypto(ctx, cryptolistReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error list crypto: %s", err))
		return
	}

	log.Printf("cryptoList success: %+v\n", res)
	respondWithJSON(w, http.StatusCreated, res)
}

func (c *HttpServer) runHTTPServer() error {
	router := mux.NewRouter()
	router.HandleFunc("/crypto", c.cryptoList).Methods("GET")
	router.HandleFunc("/crypto", c.cryptoPost).Methods("POST")

	httpAddr := fmt.Sprintf(":%s", os.Getenv("HTTP_PORT"))
	log.Printf("Starting HTTP server %s", httpAddr)
	return http.ListenAndServe(httpAddr, router)
}

func run() int {
	ctx, cancel := context.WithCancel(context.Background())
	runSignalListener(cancel)
	errCh := make(chan error, 20)

	topic := os.Getenv("TOPIC_NAME")
	brokers := []string{os.Getenv("BROKER_HOST")}

	st := initStorage(ctx)
	defer st.Close()

	publisherClient := pub.NewPublisher(brokers, pub.Topics{
		CryptoPublished: topic,
	})

	repo := repository.NewRepository(*st)
	svc := service.NewService(repo, publisherClient)
	server := HttpServer{svc: svc}

	go func() { errCh <- server.runHTTPServer() }()

	select {
	case err := <-errCh:
		fmt.Printf("App ended with err: %s", err)
		cancel()
		return 1
	case <-ctx.Done():
		fmt.Println("Context cancelled.")
		return 0
	}
}
