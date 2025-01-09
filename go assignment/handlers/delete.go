package handlers

import (
	"fmt"
	"net/http"
	"project/db"

	"go.mongodb.org/mongo-driver/bson"
)

func DeleteUser(client *db.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		if email == "" {
			http.Error(w, "Email is required", http.StatusBadRequest)
			return
		}

		collection := db.GetUserCollection(client)
		_, err := collection.DeleteOne(r.Context(), bson.M{"email": email})
		if err != nil {
			http.Error(w, "Failed to delete user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "User %s deleted successfully", email)
	}
}
