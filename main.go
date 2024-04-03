package main

import (
	"crypto/tls"
	_ "github.com/go-sql-driver/mysql"
	"github.com/us190190/messenger/database"
	"github.com/us190190/messenger/services"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"database"`
}

func main() {
	// Read config.yaml file
	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
	// Parse YAML
	var config Config
	if err = yaml.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}
	// Set environment variables
	os.Setenv("DB_HOST", config.Database.Host)
	os.Setenv("DB_PORT", config.Database.Port)
	os.Setenv("DB_USER", config.Database.User)
	os.Setenv("DB_PASSWORD", config.Database.Password)

	// Initialize database connection
	err = database.InitDB()
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
