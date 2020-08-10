package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/machinebox/graphql"
)

func main() {
	// Introduce the program
	fmt.Print("Hello, this is a small tool to get random anime from anilist\nPlease enter the number of anime you want to get :\n> ")

	// Read event
	bytes, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		fmt.Println("There was an error reading your input : ", err)
	}

	// Parse input
	number, err := strconv.Atoi(strings.TrimSpace(strings.TrimSuffix(strings.TrimSuffix(string(bytes), "\n"), "\r")))
	if err != nil || number == 0 {
		fmt.Printf("There was an error : %s\nVerify that your input is a number\n", err)
	}

	// Run queries
	for i := 0; i < number; i++ {
		fmt.Println(returnURL())
		time.Sleep(750 * time.Millisecond)
	}

	// Don't let the program stop
	fmt.Print("\nPress 'Enter' to quit the program ...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func returnURL() string {
	var randomAnime struct {
		Page struct {
			Media []struct {
				SiteURL string `json:"siteUrl"`
			} `json:"media"`
		} `json:"Page"`
	}

	// request
	req := graphql.NewRequest(`
	query ($id: Int) {
		Page(page: $id, perPage: 1) {
		  media(type: ANIME) {
			siteUrl
		  }
		}
	  }
		`)

	// Setup random ID
	req.Var("id", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(14626))

	// Run RQ
	err := graphql.NewClient("https://graphql.anilist.co").Run(context.Background(), req, &randomAnime)
	if err != nil {
		fmt.Println("there was an error while getting random anime :", err)
	}

	// return URL
	return randomAnime.Page.Media[0].SiteURL
}
