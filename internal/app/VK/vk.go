package vk

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/tarkue/tolpi-backend/config"
)

type UserGetResponse struct {
	ID              int    `json:"id"`
	Status          string `json:"status"`
	Photo           string `json:"photo_100"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	CanAccessClosed bool   `json:"can_access_closed"`
	IsClosed        bool   `json:"is_closed"`
}

type UserGet struct {
	Response []UserGetResponse `json:"response"`
}

func GetUserVK(userId string) *UserGetResponse {

	link := config.VkApiLink + config.VkUsersGetMethod + "?" + `access_token=` + config.VkServiceToken + `&user_ids=` + userId + `&fields=status,photo_100&v=5.131`

	resp, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var response = UserGet{}
	json.Unmarshal(body, &response)

	return &response.Response[0]

}
