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
	err, done := InitConfig()
	if !done {
		return
	}

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

func InitConfig() (error, bool) {
	// Read config.yaml file
	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
		return err, false
	}
	// Parse YAML
	var config Config
	if err = yaml.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
		return err, false
	}
	// Set environment variables
	err = os.Setenv("DB_HOST", config.Database.Host)
	if err != nil {
		log.Fatalf("Failed to set env: %v", err)
		return err, false
	}
	err = os.Setenv("DB_PORT", config.Database.Port)
	if err != nil {
		log.Fatalf("Failed to set env: %v", err)
		return err, false
	}
	err = os.Setenv("DB_USER", config.Database.User)
	if err != nil {
		log.Fatalf("Failed to set env: %v", err)
		return err, false
	}
	err = os.Setenv("DB_PASSWORD", config.Database.Password)
	if err != nil {
		log.Fatalf("Failed to set env: %v", err)
		return err, false
	}
	return nil, true
}
