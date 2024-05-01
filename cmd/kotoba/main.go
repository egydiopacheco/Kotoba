package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"flag"
    "encoding/base64"
	
	"github.com/google/go-github/v61/github"
	"golang.org/x/oauth2"
)

func auth_user() (context.Context, *http.Client) {
	KOTOBA_TOKEN := os.Getenv("KOTOBA_TOKEN")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: KOTOBA_TOKEN},
	)
	tc := oauth2.NewClient(ctx, ts)

	return ctx, tc
}

func parse_args() (string, string){
	flag.StringVar(&user, "user", "pass your user as param", "This is a help message")
	flag.StringVar(&message, "message", "pass your message as param" , "help message")
	flag.StringVar(&commit_message, "commit", "chore(update): via Kotoba", "help message")
	flag.Parse()

	return user, message
}

var (
	user string
	message string
	commit_message string
)

func main() {

	user, message = parse_args()

	ctx, tc := auth_user()

	client := github.NewClient(tc)

	profile_readme, _, err := client.Repositories.GetReadme(ctx, user, user, nil)
	
	if err != nil {
		fmt.Printf("The following error has occured: %v\n", err)
	}

	encoded_content := string(*profile_readme.Content)
	decoded_content, _ := base64.StdEncoding.DecodeString(encoded_content)
	
	update_readme := string(decoded_content) + "\n" + message + "\n"
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(commit_message),
		Content: []byte(update_readme),
		SHA: profile_readme.SHA,
	}

	file_path := "./README.md"
	_, _, err = client.Repositories.UpdateFile(ctx, user, user, file_path, opts)

	if err != nil {
		fmt.Printf("The following error has occured: %v\n", err)
	}
	
}
