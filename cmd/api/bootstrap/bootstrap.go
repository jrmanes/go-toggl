package bootstrap

import (
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/jrmanes/go-toggl/internal/data"
	"github.com/jrmanes/go-toggl/internal/server"
	_ "github.com/lib/pq"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file... ERROR: ", err)
	}
}

func Run() error {
	// Initialize the database connection
	d := data.New()
	if err := d.DB.Ping(); err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("SERVER_PORT")
	log.Println("port is", port)
	serv, err := server.New(port)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// start the server.
	go serv.Start()
	// Wait for an in interruptpanic .
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Attempt a graceful shutdown.
	//serv.Close()
	data.Close()
	return err
}