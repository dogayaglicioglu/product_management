package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	models "product_management/models"
	red "product_management/redis"
	"product_management/repository"

	"github.com/gorilla/mux"
	redis "github.com/redis/go-redis/v9"
)

var webRepository repository.WebRepository

func SetUserRepo(webRepo repository.WebRepository) {
	webRepository = webRepo

}

func GetUser(w http.ResponseWriter, r *http.Request) {
	cacheMiss := true
	var user *models.User
	var jsonResp []byte
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Id is required to find user", http.StatusBadRequest)
		return
	}
	cachedUser, err := getUserFromCache(id)
	if err == nil && cachedUser.Username != "" {
		jsonResp, err = json.Marshal(cachedUser)
		cacheMiss = false
		if err != nil {
			fmt.Print(w, "Error marshalling cached user data")
			cacheMiss = true
		}
	} else {
		if err != nil {
			fmt.Printf("Cache error or user not found: %v\n", err)
		}
	}
	if cacheMiss == true {
		user, err = webRepository.GetUserById(id)
		if err != nil {
			fmt.Print("Error getting user from db.. %v", err)
		}
		jsonResp, err = json.Marshal(user)
		if err != nil {
			fmt.Print("Error in marshal operation %v", err)
			return
		}
		err = storeUserInCache(*user)
		if err != nil {
			fmt.Print("error while store user into the cache %v", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)

}
func GetUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Print("HERE")
	var users []models.User
	users, err := webRepository.GetAllUsers()
	if err != nil {
		fmt.Print("Error getting users from db.. %v", err)
	}
	jsonResp, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Error in marshal operation.", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)

}

func getUserFromCache(userId string) (models.User, error) {
	ctx := context.Background()
	redisC := red.GetRedClient()
	key := fmt.Sprintf("%v", userId)
	cachedUser, err := redisC.Get(ctx, key).Result()
	if err == redis.Nil { //the user not in the redis db
		return models.User{}, fmt.Errorf("user not found in the cache")
	}
	if err != nil {
		return models.User{}, fmt.Errorf("error fetching user data from Redis: %w", err)
	}

	var user models.User
	err = json.Unmarshal([]byte(cachedUser), &user)
	if err != nil {
		return models.User{}, fmt.Errorf("error in unmarshal operation %v", err)
	}
	return user, nil
}

func storeUserInCache(user models.User) error {
	ctx := context.Background()
	redisC := red.GetRedClient()
	userData, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("error in marshal operation %v", err)
	}
	key := fmt.Sprintf("%v", user.ID)
	err = redisC.Set(ctx, key, userData, 0).Err()
	if err != nil {
		return fmt.Errorf("error storing user data in Redis: %w", err)
	}
	return nil
}
