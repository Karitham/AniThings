package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/machinebox/graphql"
)

var config struct {
	GqlURL string `json:"URL"`
	Token  string `json:"Token"`
	User   string `json:"Username"`
}

// ActivityListStruct represent the data returned to get activity query
type ActivityListStruct struct {
	Page struct {
		Activities []struct {
			ID   int    `json:"id"`
			Type string `json:"type"`
		} `json:"activities"`
	} `json:"Page"`
}

var c *graphql.Client

func main() {
	// Set config and client
	getToken("settings.json")
	c = graphql.NewClient(config.GqlURL)

	// Run infinite like
	infiniteLike(config.User)
}

func infiniteLike(user string) {
	var page int

	userID, _ := getUserID(user)
	for true {
		activities, err := queryActivities(userID, page, 50)

		if err != nil {
			fmt.Println("Error getting activities : ", err)
		}

		for i := 0; i < 50; i++ {
			time.Sleep(750 * time.Millisecond)
			likeActivity(activities.Page.Activities[i].ID)
		}
		page++
		fmt.Println("I just liked 50 activities from", user)
		time.Sleep(time.Second)
	}
}

func getUserID(name string) (userID int, err error) {
	var userIDStruct struct {
		User struct {
			ID int `json:"id"`
		} `json:"User"`
	}

	req := graphql.NewRequest(`
		query ($name: String) {
			User(name: $name) {
				id
			}
			}`)
	req.Var("name", name)
	err = c.Run(context.Background(), req, &userIDStruct)
	return userIDStruct.User.ID, err
}

func likeActivity(id int) (userList interface{}, err error) {
	req := graphql.NewRequest(`
			mutation ($id: Int) {
		ToggleLike(id: $id, type:ACTIVITY) {
			name
		}
		}`)
	req.Header.Add("Authorization", "Bearer "+config.Token)
	req.Var("id", id)
	err = c.Run(context.Background(), req, &userList)
	return
}

// get config from file
func getToken(file string) {
	body, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("error reading config file :", err)
	}
	er := json.Unmarshal(body, &config)
	if er != nil {
		fmt.Println("error unmarshalling config :", er)
	}
}

func queryActivities(userID int, page int, count int) (Activities *ActivityListStruct, err error) {
	req := graphql.NewRequest(`
	query ($userId: Int, $page: Int, $count: Int) {
		Page(page: $page, perPage: $count) {
		  activities(sort: ID_DESC, userId: $userId) {
			... on TextActivity {
			  id
			  type
			}
			... on ListActivity {
			  id
			  type
			}
			... on MessageActivity {
			  id
			  type
			}
		  }
		}
	  }
	`)
	req.Header.Add("Authorization", "Bearer "+config.Token)
	req.Var("userId", userID)
	req.Var("page", page)
	req.Var("count", count)
	err = c.Run(context.Background(), req, &Activities)
	return
}
