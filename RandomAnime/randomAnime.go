package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/machinebox/graphql"
)

// AnimeStruct represents the anime returned by the query
type AnimeStruct struct {
	Page struct {
		Media []struct {
			SiteURL string `json:"siteUrl"`
		} `json:"media"`
	} `json:"Page"`
}

func main() {
	var anime AnimeStruct
	c := Init()

	// Run queries
	for i := 0; i < c; i++ {
		fmt.Println(anime.returnURL().Page.Media[0].SiteURL)
		time.Sleep(750 * time.Millisecond)
	}

	// Don't let the program stop
	fmt.Print("Press 'Enter' to quit the program...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// Init intialises the flag
func Init() int {
	count := flag.Int("count", 1, "count makes you able to request more than one anime at a time using `--count n`")
	flag.Parse()
	return *count
}

func (anime *AnimeStruct) returnURL() *AnimeStruct {
	req := graphql.NewRequest(`
	query ($id: Int) {
		Page(page: $id, perPage: 1) {
		  media(type: ANIME) {
			siteUrl
		  }
		}
	  }
		`)
	req.Var("id", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(14626))

	err := graphql.NewClient("https://graphql.anilist.co").Run(context.Background(), req, &anime)
	if err != nil {
		log.Println("there was an error while getting random anime :", err)
	}
	return anime
}
