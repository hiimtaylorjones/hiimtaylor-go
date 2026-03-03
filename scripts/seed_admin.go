//go:build ignore

package main

import (
			"fmt"
			"log"

			"golang.org/x/crypto/bcrypt"

			"github.com/hiimtaylorjones/hiimtaylor-go/database"
			"github.com/hiimtaylorjones/hiimtaylor-go/models"
)

func main() {
			database.Connect()
			defer database.Close()

			password := "your-password-here"
			hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
							log.Fatal(err)
			}

			admin, err := models.CreateAdmin("admin@hiimtaylorjones.com", string(hash))
			if err != nil {
							log.Fatal(err)
			}

			fmt.Printf("Created admin: %s (id: %d)\n", admin.Email, admin.ID)
}