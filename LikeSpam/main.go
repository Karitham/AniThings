package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/buger/goterm"
	"github.com/machinebox/graphql"
)

// Flags are the different flags you can pass to the program
type Flags struct {
	count    int
	token    string
	username string
	userID   int
}

// ActivitiesQueryStruct is used to query activities
type ActivitiesQueryStruct struct {
	Page struct {
		Activities []Activity `json:"activities"`
	} `json:"Page"`
}

// Activity represent an activity node
type Activity struct {
	Typename string `json:"__typename"`
	ID       int    `json:"id"`
	IsLiked  bool   `json:"isLiked"`
}

var c *graphql.Client = graphql.NewClient("https://graphql.anilist.co")

func main() {
	f := getFlags()
	goterm.Clear()
	goterm.MoveCursor(1, 1)
	goterm.Flush()
	f.runLiker()
}
func getFlags() (f Flags) {
	flag.IntVar(
		&f.count,
		"count",
		0,
		"use count to specify a like count",
	)
	flag.StringVar(
		&f.username,
		"user",
		"",
		"Use user to set the name of the user to spam",
	)
	flag.StringVar(
		&f.token,
		"token",
		"",
		"Enter your anilist token, to get one, go to https://anilist.co/api/v2/oauth/authorize?client_id=3971&response_type=token",
	)
	flag.Parse()
	if f.token == "" || f.username == "" {
		log.Fatalln("error starting the liker, you need to provide the flags. Check `-help` for help")
	}
	f.getUserID()
	return f
}

func (f *Flags) runLiker() {
	var page, likes int
	for likes < f.count {
		a := f.queryActivities(page)
		for i := 0; i < len(a.Page.Activities); i++ {
			if !a.Page.Activities[i].IsLiked {
				time.Sleep(2750 * time.Millisecond)
				a.Page.Activities[i].like(f.token)
				likes++
				fmt.Printf("\rLiked %d activities from %s", likes, f.username)
			}
		}
		time.Sleep(time.Second)
		page++
	}
}
func (f *Flags) getUserID() {
	var userID struct {
		User struct {
			ID int `json:"id"`
		} `json:"User"`
	}
	req := graphql.NewRequest(`
	query ($name: String) {
		  User(search: $name) {
			id
		}
	  }`)

	req.Var("name", f.username)

	err := c.Run(context.Background(), req, &userID)
	if err != nil {
		log.Fatalln("error retrieving user ID : ", err)
	}
	f.userID = userID.User.ID
}

func (a *Activity) like(token string) {
	req := graphql.NewRequest(`
	mutation ($id: Int) {
		ToggleLikeV2(id: $id, type: ACTIVITY) {
		  __typename
		}
	  }`)

	req.Header.Add("Authorization", "Bearer "+token)
	req.Var("id", a.ID)

	err := c.Run(context.Background(), req, nil)
	if err != nil {
		log.Println("Error linking activity : ", err)
	}
}

func (f Flags) queryActivities(page int) (Activities ActivitiesQueryStruct) {
	req := graphql.NewRequest(`
	query ($userId: Int, $page: Int) {
		Page(page: $page, perPage: 50) {
		  activities(sort: ID_DESC, userId: $userId) {
			__typename
			... on TextActivity {
			  id
			  isLiked
			}
			... on ListActivity {
			  id
			  isLiked
			}
			... on MessageActivity {
			  id
			  isLiked
			}
		  }
		}
	  }`)

	req.Var("userId", f.userID)
	req.Var("page", page)
	req.Header.Set("Authorization", "Bearer "+f.token)

	err := c.Run(context.Background(), req, &Activities)
	if err != nil {
		log.Println("Error getting activities : ", err)
	}

	return
}
