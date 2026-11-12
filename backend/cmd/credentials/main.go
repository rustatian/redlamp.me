// Command credentials makes the LEETCODE_ACCOUNTS secret for the usernames of the parent and the student.
// It prints the JSON of the accounts on standard output and writes the new passwords to the terminal,
// so a redirect of the output to a file stores no password.
package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type account struct {
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
}

var roles = [2]string{"parent", "student"}

var usernamePattern = regexp.MustCompile(`^[A-Za-z0-9._-]{3,40}$`)

// generate makes a random password for each of the two usernames. It returns the secret and the passwords, the parent first.
func generate(usernames []string) (secret []byte, passwords [2]string, err error) {
	invalid := errors.New("supply two distinct usernames with 3-40 letters, digits, dots, dashes, or underscores")
	if len(usernames) != 2 {
		return nil, passwords, invalid
	}
	var accounts [2]account
	for i, username := range usernames {
		username = strings.TrimSpace(username)
		if !usernamePattern.MatchString(username) {
			return nil, passwords, invalid
		}
		random := make([]byte, 18)
		if _, err := rand.Read(random); err != nil {
			return nil, passwords, err
		}
		passwords[i] = base64.RawURLEncoding.EncodeToString(random)
		sum := sha256.Sum256([]byte(passwords[i]))
		accounts[i] = account{Username: username, PasswordHash: hex.EncodeToString(sum[:])}
	}
	if accounts[0].Username == accounts[1].Username {
		return nil, passwords, invalid
	}
	secret, err = json.Marshal(struct {
		Parent  account `json:"parent"`
		Student account `json:"student"`
	}{accounts[0], accounts[1]})
	return secret, passwords, err
}

func main() {
	usernames := os.Args[1:]
	secret, passwords, err := generate(usernames)
	if err != nil {
		log.Fatal(err)
	}
	terminal, err := os.OpenFile("/dev/tty", os.O_WRONLY, 0)
	if err != nil {
		log.Fatal(err)
	}
	for i, role := range roles {
		if _, err := fmt.Fprintf(terminal, "%s: %s\nPassword: %s\n", role, strings.TrimSpace(usernames[i]), passwords[i]); err != nil {
			log.Fatal(err)
		}
	}
	if _, err := fmt.Fprintln(terminal, "Save these passwords in your password manager."); err != nil {
		log.Fatal(err)
	}
	if err := terminal.Close(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(secret))
}
