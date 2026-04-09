package main

import (
    "errors"
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
    httpapi "compra-certa/api/internal/http"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func main() {
    _ = godotenv.Load()

    db, err := openDB()
    if err != nil {
        log.Fatal(err)
    }

    router := httpapi.SetupRouter(db)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    addr := ":" + port
    log.Printf("listening on %s", addr)
    if err := router.Run(addr); err != nil {
        log.Fatal(err)
    }
}

func openDB() (*gorm.DB, error) {
    host, err := getenvRequired("DB_HOST")
    if err != nil {
        return nil, err
    }

    port, err := getenvRequired("DB_PORT")
    if err != nil {
        return nil, err
    }

    name, err := getenvRequired("DB_NAME")
    if err != nil {
        return nil, err
    }

    user, err := getenvRequired("DB_USER")
    if err != nil {
        return nil, err
    }

    password, err := getenvRequired("DB_PASSWORD")
    if err != nil {
        return nil, err
    }

    sslmode, err := getenvRequired("DB_SSLMODE")
    if err != nil {
        return nil, err
    }

    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        host,
        port,
        user,
        password,
        name,
        sslmode,
    )

    return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func getenvRequired(key string) (string, error) {
    value := os.Getenv(key)
    if value == "" {
        return "", errors.New("missing required env: " + key)
    }
    return value, nil
}