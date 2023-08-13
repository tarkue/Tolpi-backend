package vk

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/tarkue/tolpi-backend/config"
)

type UserGet struct {
	Response []struct {
		ID              int    `json:"id"`
		Status          string `json:"status"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		CanAccessClosed bool   `json:"can_access_closed"`
		IsClosed        bool   `json:"is_closed"`
	} `json:"response"`
}

func GetUserStatus(userId string) *string {

	link := config.VkApiLink + config.VkUsersGetMethod + "?" + `access_token=` + config.VkServiceToken + `&user_ids=` + userId + `&fields=status&v=5.131`

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

	return &response.Response[0].Status

}
