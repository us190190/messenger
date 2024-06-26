package main

import (
	"crypto/tls"
	_ "github.com/go-sql-driver/mysql"
	"github.com/us190190/messenger/database"
	"github.com/us190190/messenger/services"
	"log"
	"net/http"
)

func main() {

	// Initialize database connection
	err := database.InitDB()
	if err != nil {
		log.Fatal("Application unable to connect with DB")
		return
	}

	// Routes
	http.HandleFunc("/v1/user/register", services.UserRegisterHandler)
	http.HandleFunc("/v1/user/update", services.UserUpdateHandler)
	http.HandleFunc("/v1/user/authenticate", services.UserAuthenticationHandler)
	http.HandleFunc("/v1/user/remove", services.UserRemoveHandler)
	http.HandleFunc("/v1/user/all", services.UserAllHandler)
	http.HandleFunc("/v1/group/all", services.GroupAllHandler)
	http.HandleFunc("/start", services.HandleWebSocket)

	// Start server with HTTPS
	tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12} // You should configure this with your actual certificate and key paths
	server := &http.Server{
		Addr:      ":443",
		Handler:   nil, // Default router
		TLSConfig: tlsConfig,
	}
	log.Fatal(server.ListenAndServeTLS("/app/server.crt", "/app/server.key"))
}
