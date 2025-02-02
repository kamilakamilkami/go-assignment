package patient

import (
  "context"
  "encoding/json"
  "net/http"
  "time"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
  "med_portal/db"
)

type Patient struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name    string             `json:"name" bson:"name"`
	Email   string             `json:"email" bson:"email"`
	Phone   string             `json:"phone" bson:"phone"`
	Address string             `json:"address" bson:"address"`
	Role string             `json:"role" bson:"role"`
	}

func SearchPatientsHandler(w http.ResponseWriter, r *http.Request) {
  query := r.URL.Query().Get("name")
  if query == "" {
    http.Error(w, "Missing search query", http.StatusBadRequest)
    return
  }

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  _, collection := db.Connect("med_portal", "patients")
  cursor, err := collection.Find(ctx, bson.M{"name": bson.M{"$regex": query, "$options": "i"}})
  if err != nil {
    http.Error(w, "Failed to search patients", http.StatusInternalServerError)
    return
  }
  defer cursor.Close(ctx)

  

  var patients []Patient
  if err := cursor.All(ctx, &patients); err != nil {
    http.Error(w, "Failed to decode search results", http.StatusInternalServerError)
    return
  }
  
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(patients)
}