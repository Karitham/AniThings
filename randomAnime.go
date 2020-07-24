package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/machinebox/graphql"
)

// RadomAnimeStruct holds data from random anime query
type RadomAnimeStruct struct {
	Page struct {
		Media []struct {
			SiteURL string `json:"siteUrl"`
		} `json:"media"`
	} `json:"Page"`
}

func main() {
	var res RadomAnimeStruct

	// build query
	graphURL := "https://graphql.anilist.co"
	client := graphql.NewClient(graphURL)
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
	client.Run(context.Background(), req, &res)

	// Print resulting URL
	fmt.Println(res.Page.Media[0].SiteURL)

	select {}
}
