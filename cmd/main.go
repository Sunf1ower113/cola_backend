package main

import (
	"auth-api/internal/composites"
	"auth-api/internal/config"
	"auth-api/pkg/client/sqlite"
	"github.com/rs/cors"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	cfg, err := config.LoadConfiguration("config.json")
	if err != nil {
		log.Panicf("cannot load configuration: %v", err)
	}
	database, err := sqlite.NewDB(cfg.Database.DbDriver, cfg.Database.DbName)
	if err != nil {
		log.Panicf("cannot create db: %v", err)
	}
	if database.Ping() != nil {
		log.Panicf("cannot ping db: %v", err)
	}
	router := http.NewServeMux()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	handlerWithCORS := c.Handler(router)
	userComposite, err := composites.NewUserComposite(database)
	userComposite.Handler.Register(router)

	recycleBoxComposite, err := composites.NewRecycleBoxComposite(database)
	recycleBoxComposite.Handler.Register(router)

	if err != nil {
		log.Fatal(err)
	}
	start(handlerWithCORS, cfg)
}

func start(router http.Handler, cfg *config.Config) {
	log.Println("Start the application...")
	listener, err := net.Listen(cfg.Listener.Protocol, cfg.Listener.Host+cfg.Listener.Port)
	if err != nil {
		log.Fatal(err)
	}
	server := &http.Server{
		Handler:      router,
		IdleTimeout:  time.Duration(cfg.Listener.IdleTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Listener.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.Listener.ReadTimeout) * time.Second,
	}
	log.Printf("Server is listening port %s%s\n", cfg.Listener.Host, cfg.Listener.Port)
	log.Panic(server.Serve(listener))
}
