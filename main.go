package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/tabwriter"
)

type UsersList struct {
	Members []Member `json:"members"`
}

type Member struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	RealName string `json:"real_name"`
	Profile  struct {
		Email string `json:"email"`
		Phone string `json:"phone"`
		Title string `json:"title"`
	} `json:"profile"`
}

func main() {
	resp, err := getMethod("users.list")
	if err != nil {
		println("http NOK")
		return
	}

	members, err := parseResp(resp)
	if err != nil {
		println("json NOK")
		return
	}

	member, err := pickMember("flavou", members)
	if err != nil {
		println("member not found")
		return
	}

	printMember(member)
}
func getMethod(method string) ([]byte, error) {
	slackToken := "XXX"

	resp, err := http.Get("https://slack.com/api/users.list?token=" + slackToken + "&presence=1")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func parseResp(body []byte) ([]Member, error) {
	var usersList UsersList
	err := json.Unmarshal([]byte(body), &usersList)
	return usersList.Members, err
}

func pickMember(expectedName string, members []Member) (Member, error) {
	for _, m := range members {
		if m.Name == expectedName {
			return m, nil
		}
	}
	return Member{}, errors.New("Sorry, we found no one")
}

func printMember(member Member) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "pseudo\treal name\ttitle\temail\tphone")
	fmt.Fprintln(w, "------\t---------\t-----\t-----\t-----")
	fmt.Fprintln(w, member.Name+"\t"+member.RealName+"\t"+member.Profile.Title+"\t"+member.Profile.Email+"\t"+member.Profile.Phone)
	w.Flush()
}
