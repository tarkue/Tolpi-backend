package config

import "os"

var ServerPort = os.Getenv("PORT")

var DataBaseUri = os.Getenv("MONGODB_URI")
var DataBaseUserName = os.Getenv("MONGODB_USERNAME")
var DataBasePassword = os.Getenv("MONGODB_PASSWORD")

var VkServiceToken = os.Getenv("VK_SERVICE_TOKEN")
var SecretKey = os.Getenv("VK_SECRET_KEY")

var VkApiLink = "https://api.vk.com/method/"
var VkUsersGetMethod = "users.get"
