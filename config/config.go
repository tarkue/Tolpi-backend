package config

import "os"

var ServerPort = os.Getenv("PORT")
var DataBaseUri = os.Getenv("MONGODB_URI")
var VkServiceToken = os.Getenv("VK_SERVICE_TOKEN")
var SecretKey = os.Getenv("VK_SECRET_KEY")

var VkApiLink = "https://api.vk.com/method/"
var CountriesApi = "https://api.hh.ru/areas"
var VkUsersGetMethod = "users.get"
