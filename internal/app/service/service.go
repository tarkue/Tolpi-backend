package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/url"
	"sort"
	"strings"

	"github.com/labstack/echo/v4"
)

type Service struct {
}

type queryParameter struct {
	Key   string
	Value string
}

const SignInvalid = "подпись невалидна"

func New() *Service {
	return &Service{}
}

func (service *Service) GetUserId(c echo.Context) string {
	clientUrl := c.Request().URL.String()
	clientUrls := strings.Split(clientUrl, "vk_user_id=")
	userId := strings.Split(clientUrls[1], "&")[0]
	return userId
}

func (service *Service) VerifyLaunchParams(querySearch string, secretKey string) error {
	var searchIndex = strings.Index(querySearch, "?")

	// Необходимо удалить всё, что находится до search части в случае, если
	// эта часть существует.
	if searchIndex >= 0 {
		querySearch = querySearch[searchIndex+1:]
	}

	var (
		// Отфильтрованные параметры запуска. Мы используем именно
		// слайс по той причине, что позже нам будет необходимым этот слайс
		// отсортировать по возрастанию ключа параметра.
		query []queryParameter
		// Подпись, которая была сгенерирована сервером ВКонтакте и основана на
		// параметрах из query.
		sign string
	)

	// Разделяем параметры запуска на вхождения, разделенные знаком "&".
	for _, part := range strings.Split(querySearch, "&") {
		var keyAndValue = strings.Split(part, "=")
		var key = keyAndValue[0]
		var value string

		if len(keyAndValue) > 1 {
			value = keyAndValue[1]
		}

		// Мы обрабатываем только те ключи, которые начинаются с префикса "vk_".
		// Все остальные ключи в создании подписи не участвуют.
		if strings.HasPrefix(key, "vk_") {
			query = append(query, queryParameter{key, value})
		} else if key == "sign" {
			// Если ключ равен "sign", то в значении записана подпись параметров
			// запуска.
			sign = value
		}
	}

	// В случае, если подпись параметров не удалось найти, либо параметров с
	// префиксом "vk_" передано не было, мы считаем параметры запуска невалидными.
	if sign == "" || len(query) == 0 {
		return errors.New(SignInvalid)
	}

	// Сортируем параметры запуска по порядку их возрастания.
	sort.SliceStable(query, func(a int, b int) bool {
		return query[a].Key < query[b].Key
	})

	// Далее снова превращаем параметры запуска в единую строку.
	var queryString = ""

	for idx, param := range query {
		if idx > 0 {
			queryString += "&"
		}
		queryString += param.Key + "=" + url.PathEscape(param.Value)
	}

	// Далее нам необходимо вычислить хэш SHA-256.
	var hashCreator = hmac.New(sha256.New, []byte(secretKey))
	hashCreator.Write([]byte(queryString))

	var hash = base64.URLEncoding.EncodeToString(hashCreator.Sum(nil))

	// Далее по правилам создания параметров запуска ВКонтакте, необходимо
	// произвести ряд замен символов.
	hash = strings.ReplaceAll(hash, "+", "-")
	hash = strings.ReplaceAll(hash, "/", "_")
	hash = strings.ReplaceAll(hash, "=", "")

	if sign != hash {
		return errors.New(SignInvalid)
	}
	return nil
}

func (service *Service) RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func (service *Service) IndexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

func (service *Service) Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
