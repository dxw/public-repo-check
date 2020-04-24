package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/tomnomnom/linkheader"
)

func main() {
	if len(os.Args) == 1 {
		log.Fatal("Usage: public-repo-check ${ORG_NAME}")
	}
	for _, arg := range os.Args[1:] {
		checkOrg(arg)
	}
}

func checkOrg(org string) {
	repos := fetchRepos(org)

	for _, repo := range repos {
		if repo.Archived {
			message(repo, true, "Archived. No further checks")
			continue
		}

		if repo.Fork {
			message(repo, true, "Fork. No further checks")
			continue
		}

		checkLicense(repo)
		checkReadme(repo)
		checkContributing(repo)
	}
}

func checkLicense(repo Repo) {
	if repo.License == nil {
		message(repo, false, "License missing!")
	} else if repo.License.Url != "https://api.github.com/licenses/mit" {
		message(repo, false, fmt.Sprintf("License not MIT (%s)", repo.License.Url))
	} else {
		message(repo, true, "License OK")
	}
}

func checkFile(repo Repo, file string) (bool, error) {
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", repo.FullName, repo.DefaultBranch, file)

	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}

func checkReadme(repo Repo) {
	check, err := checkFile(repo, "README.md")
	if err != nil {
		log.Fatal(err)
	}

	if check {
		message(repo, true, "Has README")
	} else {
		message(repo, false, "No README found")
	}
}

func checkContributing(repo Repo) {
	check, err := checkFile(repo, "CONTRIBUTING.md")
	if err != nil {
		log.Fatal(err)
	}

	if check {
		message(repo, true, "Has CONTRIBUTING")
	} else {
		message(repo, false, "No CONTRIBUTING found")
	}
}

func message(repo Repo, ok bool, message string) {
	status := "❌"
	if ok {
		status = "✅"
	}

	fmt.Printf("%s %s: %s\n", status, repo.FullName, message)
}

var endOfList = errors.New("end of list")

func fetchRepos(org string) []Repo {
	allRepos := []Repo{}
	url := fmt.Sprintf("https://api.github.com/orgs/%s/repos", org)

	for {
		repos, newUrl, err := fetchSomeRepos(url)
		if err != nil && err != endOfList {
			log.Fatal(err)
		}

		allRepos = append(allRepos, repos...)
		if err == endOfList {
			break
		}

		url = newUrl
	}

	return allRepos
}

func fetchSomeRepos(url string) ([]Repo, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, "", errors.New("non-200 status")
	}
	defer resp.Body.Close()

	var repos []Repo
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&repos)
	if err != nil {
		return nil, "", err
	}

	links := linkheader.Parse(resp.Header["Link"][0]).FilterByRel("next")

	if len(links) == 0 {
		return repos, "", endOfList
	} else {
		return repos, links[0].URL, nil
	}
}

type Repo struct {
	Name          string   `json:"name"`
	FullName      string   `json:"full_name"`
	Private       bool     `json:"private"`
	HtmlUrl       string   `json:"html_url"`
	Fork          bool     `json:"fork"`
	Url           string   `json:"url"` // API URL
	CloneUrl      string   `json:"clone_url"`
	Archived      bool     `json:"archived"`
	License       *License `json:"license"`
	DefaultBranch string   `json:"default_branch"`
}

type License struct {
	Url string `json:"url"`
}
