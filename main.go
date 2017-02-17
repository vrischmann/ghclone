package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"sort"

	"github.com/google/go-github/github"
)

var (
	ghClient *github.Client

	flOrg = flag.Bool("org", false, "Is it an organization instead of a user")
)

func getRepositories(name string) ([]*github.Repository, error) {
	var (
		repos []*github.Repository
		err   error
	)

	switch {
	case *flOrg:
		log.Println("lalala")
		opts := &github.RepositoryListByOrgOptions{Type: "public"}
		repos, _, err = ghClient.Repositories.ListByOrg(name, opts)
		if err != nil {
			return nil, err
		}

	default:
		opts := &github.RepositoryListOptions{Visibility: "public"}
		repos, _, err = ghClient.Repositories.List(name, opts)
		if err != nil {
			return nil, err
		}
	}

	return repos, nil
}

func gitClone(url string) error {
	cmd := exec.Command("git", "clone", url)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("Usage: ghclone [-org] <name>")
	}

	name := flag.Arg(0)
	ghClient = github.NewClient(nil)

	repos, err := getRepositories(name)
	if err != nil {
		log.Fatal(err)
	}

	sort.Slice(repos, func(i, j int) bool {
		return *repos[i].Name < *repos[j].Name
	})

	for _, repo := range repos {
		log.Printf("cloning repo %s with %s", *repo.Name, *repo.CloneURL)

		err := gitClone(*repo.CloneURL)
		if err != nil {
			log.Printf("unable to clone repo %s. err=%v", *repo.Name, err)
		}
	}
}
