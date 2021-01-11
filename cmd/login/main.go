package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func Main(ctx context.Context) error {
	username := flag.String("username", "", "Tesla account username")
	passcode := flag.String("passcode", "", "Tesla account passcode")
	flag.Parse()

	fmt.Print("Password: ")
	b, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("read password: %w", err)
	}
	fmt.Println()

	password := string(b)

	switch "" {
	case *username:
		return errors.New("username must be set")
	case password:
		return errors.New("password must be set")
	case *passcode:
		return errors.New("passcode must be set")
	}

	c, err := newClient()
	if err != nil {
		return fmt.Errorf("new client: %w", err)
	}

	token, err := c.Login(ctx, *username, password, *passcode)
	if err != nil {
		return fmt.Errorf("login: %w", err)
	}

	fmt.Printf("Your token is: %s\n", token)

	return nil
}

func main() {
	if err := Main(context.Background()); err != nil {
		log.Fatal(err)
	}
}
