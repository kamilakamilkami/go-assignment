// handlers/update.go
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project/db"
	"project/models"

	"go.mongodb.org/mongo-driver/bson"
)

func UpdateUser(client *db.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		if email == "" {
			http.Error(w, "Email is required", http.StatusBadRequest)
			return
		}

		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		collection := db.GetUserCollection(client)
		_, err := collection.UpdateOne(r.Context(), bson.M{"email": email}, bson.D{{"$set", user}})
		if err != nil {
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "User %s updated successfully", user.Name)
	}
}
