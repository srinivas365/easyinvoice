package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"pharmacy-erp/config"
	"pharmacy-erp/models"
	"pharmacy-erp/repository"
	"pharmacy-erp/services"

	"gorm.io/gorm"
)

func main() {
	log.SetFlags(0)

	username := flag.String("username", "", "login username (required)")
	password := flag.String("password", "", "plain password; omit with -password-stdin")
	passwordStdin := flag.Bool("password-stdin", false, "read password from stdin (first line only)")
	role := flag.String("role", "staff", "user role: admin or staff")
	flag.Parse()

	if *username == "" {
		flag.Usage()
		os.Exit(2)
	}

	pw := strings.TrimSpace(*password)
	if *passwordStdin {
		line, err := readPasswordFromStdin()
		if err != nil {
			log.Fatal(err)
		}
		pw = line
	}
	if pw == "" {
		log.Fatal("password is required (use -password or -password-stdin)")
	}

	r := strings.ToLower(strings.TrimSpace(*role))
	if r != "admin" && r != "staff" {
		log.Fatal("role must be admin or staff")
	}

	config.InitDB()
	if err := config.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("migrate users:", err)
	}

	userRepo := repository.NewUserRepository()
	if _, err := userRepo.FindByUsername(*username); err == nil {
		log.Fatalf("user %q already exists", *username)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Fatal("lookup user:", err)
	}

	authSvc := services.NewAuthService()
	hash, err := authSvc.HashPassword(pw)
	if err != nil {
		log.Fatal("hash password:", err)
	}

	u := &models.User{
		Username: strings.TrimSpace(*username),
		Password: hash,
		Role:     r,
	}
	if err := userRepo.Create(u); err != nil {
		log.Fatal("create user:", err)
	}

	fmt.Printf("Created user %q with role %q (id=%d)\n", u.Username, u.Role, u.ID)
}

func readPasswordFromStdin() (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", errors.New("-password-stdin requires a piped or redirected stdin")
	}
	sc := bufio.NewScanner(os.Stdin)
	if !sc.Scan() {
		if err := sc.Err(); err != nil {
			return "", err
		}
		return "", errors.New("stdin is empty")
	}
	return strings.TrimSpace(sc.Text()), nil
}
