package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/tomnomnom/linkheader"
)

func main() {
	checkOrg("dxw")
}

func checkOrg(org string) {
	repos := fetchRepos(org)

	for _, repo := range repos {
		fmt.Printf("%s:\n", repo.Name)

		if repo.Archived {
			fmt.Printf("✅ Archived. No further checks\n")
			continue
		}

		if repo.Fork {
			fmt.Printf("✅ Fork. No further checks\n")
			continue
		}

		checkLicense(repo)
	}
}

func checkLicense(repo Repo) {
	if repo.License == nil {
		fmt.Printf("❌ License missing!\n")
	} else if repo.License.Url != "https://api.github.com/licenses/mit" {
		fmt.Printf("❌ License not MIT (%s)\n", repo.License.Url)
	} else {
		fmt.Printf("✅ License OK\n")
	}
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
	Name     string   `json:"name"`
	Private  bool     `json:"private"`
	HtmlUrl  string   `json:"html_url"`
	Fork     bool     `json:"fork"`
	Url      string   `json:"url"` // API URL
	CloneUrl string   `json:"clone_url"`
	Archived bool     `json:"archived"`
	License  *License `json:"license"`
}

type License struct {
	Url string `json:"url"`
}
