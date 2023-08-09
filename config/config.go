package config

var ServerPort = "4000"                                                                        //os.Getenv("PORT")
var DataBaseUri = "mongodb://localhost:27017"                                                  //os.Getenv("MONGODB_URI")
var VkServiceToken = "d1e3cb79d1e3cb79d1e3cb792dd2f6c27bdd1e3d1e3cb79b55e8638bbcee6776f69e90c" //os.Getenv("VK_SERVICE_TOKEN")
var SecretKey = "4GScw7G7MZcrBsOWcnKg"                                                         //os.Getenv("VK_SECRET_KEY")

var VkApiLink = "https://api.vk.com/method/"
var CountriesApi = "https://api.hh.ru/areas"
var VkUsersGetMethod = "users.get"
