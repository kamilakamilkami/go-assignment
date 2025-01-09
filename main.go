package main

import (
  "context"
  "encoding/json"
  "log"
  "net/http"
  "strconv"
  "time"
  "github.com/gorilla/mux"
  "github.com/rs/cors"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "os"
  "github.com/sirupsen/logrus"
  "encoding/base64"
  "io"
  "fmt"
  "net/smtp"

)

var client *mongo.Client
var collection *mongo.Collection

type Patient struct {
  ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
  Name    string             `json:"name" bson:"name"`
  Email   string             `json:"email" bson:"email"`
  Phone   string             `json:"phone" bson:"phone"`
  Address string             `json:"address" bson:"address"`
  Role string             `json:"role" bson:"role"`
  }
  

func main() {

  logs := logrus.New()

	// Set the log format to JSON
	logs.SetFormatter(&logrus.JSONFormatter{})

	// Log output to a file
	file, fileErr := os.OpenFile("user_actions.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if fileErr != nil {
		log.Fatalf("Failed to open log file: %v", fileErr)
	}
	defer file.Close()

	// Set the output of the logger to the file
	logs.SetOutput(file)


  // Connect to MongoDB
  clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
  var err error
  client, err = mongo.Connect(context.TODO(), clientOptions)
  if err != nil {
    log.Fatal("MongoDB Connection Error:", err)
  }

  err = client.Ping(context.TODO(), nil)
  if err != nil {
    log.Fatal("MongoDB Ping Error:", err)
  }
  log.Println("Connected to MongoDB!")

  collection = client.Database("med_portal").Collection("patients")

  // Set up router and routes
  router := mux.NewRouter()
  router.HandleFunc("/add-patient", addPatientHandler).Methods("POST")
  router.HandleFunc("/get-all-patients", func(w http.ResponseWriter, r *http.Request) {
		getPatientsHandler(w, r, logs)
	})
  router.HandleFunc("/delete-patient/{id}", deletePatientHandler).Methods("DELETE")
  router.HandleFunc("/update-patient/{id}", updatePatientHandler).Methods("PUT")
  router.HandleFunc("/search-patients", searchPatientsHandler).Methods("GET")
  // I added
  router.HandleFunc("/send-email", sendEmailGet).Methods("GET")
  router.HandleFunc("/send-support-message", sendEmailPost).Methods("POST")

  // Enable CORS
  corsHandler := cors.New(cors.Options{
    AllowedOrigins:   []string{"http://127.0.0.1:5502", "http://localhost:5502", "http://127.0.0.1:5500", "http://localhost:5500"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"Content-Type", "Authorization"},
    AllowCredentials: true,
  })

  log.Println("Server running on http://localhost:8080")
  log.Fatal(http.ListenAndServe(":8080", corsHandler.Handler(router)))
}
// I added
func sendEmailGet(w http.ResponseWriter, r *http.Request) {

}

func sendEmailPost(w http.ResponseWriter, r *http.Request) {
  // Parse multipart form data
  err := r.ParseMultipartForm(10 << 20) // Limit file size to ~10 MB
  if err != nil {
      http.Error(w, "Failed to parse form data", http.StatusBadRequest)
      return
  }

  subject := r.FormValue("subject")
  message := r.FormValue("message")

  var emailMessage EmailMessage
  emailMessage.Subject = subject
  emailMessage.Message = message

  // Handle file upload
  file, handler, err := r.FormFile("attachment")
  if err == nil {
      defer file.Close()
      emailMessage.FileName = handler.Filename
      emailMessage.PhotoData, err = io.ReadAll(file) // io
      if err != nil {
          http.Error(w, "Failed to read file", http.StatusInternalServerError)
          return
      }
  }

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  cursor, err := collection.Find(ctx, bson.M{})
  if err != nil {
      http.Error(w, "Failed to fetch patients", http.StatusInternalServerError)
      return
  }
  defer cursor.Close(ctx)

  var patients []Patient
  if err := cursor.All(ctx, &patients); err != nil {
      http.Error(w, "Failed to decode patients", http.StatusInternalServerError)
      return
  }

  for _, patient := range patients {
      if patient.Role == "admin" {
          err := sendEmailImage(emailMessage.Subject, emailMessage.Message, patient.Email, emailMessage.FileName, emailMessage.PhotoData)
          if err != nil {
              http.Error(w, "Failed to send email: "+err.Error(), http.StatusInternalServerError)
              return
          }
      }
  }

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(map[string]string{"message": "Email sent successfully"})
}

type EmailMessage struct {
  Subject   string
  Message   string
  FileName  string // Name of the uploaded file
  PhotoData []byte // Content of the uploaded file
}

func sendEmailImage(subject, message, recipient, filename string, photoData []byte) error {


	from := "dautovalisher33@gmail.com"
  password := "ppzv hnga kzhl xdlz"
  smtpHost := "smtp.gmail.com" 
  smtpPort :=  "587"
	encodedPhoto := base64.StdEncoding.EncodeToString(photoData) // "encoding/base64"

    // Create the MIME email
    email := fmt.Sprintf(
        "From: %s\nTo: %s\nSubject: %s\nMIME-Version: 1.0\nContent-Type: multipart/mixed; boundary=boundary\n\n"+
            "--boundary\nContent-Type: text/plain; charset=utf-8\n\n%s\n\n"+
            "--boundary\nContent-Type: image/jpeg\nContent-Transfer-Encoding: base64\nContent-Disposition: attachment; filename=\"%s\"\n\n%s\n--boundary--",
        from, recipient, subject, message, filename, encodedPhoto,
    )

    // Send the email
    auth := smtp.PlainAuth("", from, password, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{recipient}, []byte(email))

}



func SendEmail(subject, message string, recipient string) error {

  from := "dautovalisher33@gmail.com"
  password := "ppzv hnga kzhl xdlz"
  smtpHost := "smtp.gmail.com" 
  smtpPort :=  "587" 

  msg := []byte(fmt.Sprintf("Subject: %s\n\n%s", subject, message))

  auth := smtp.PlainAuth("", from, password, smtpHost)

  return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{recipient}, msg)
}

// nothing



func getPatientsHandler(w http.ResponseWriter, r *http.Request, logs *logrus.Logger) {
  filter := r.URL.Query().Get("filter")
  sort := r.URL.Query().Get("sort")
  pageStr := r.URL.Query().Get("page")
  LogUserAction(logs, "get patients", "1", pageStr, map[string]interface{}{"filter": filter, "sort": sort})
  limit := 10
  offset := 0
  if page, err := strconv.Atoi(pageStr); err == nil && page > 1 {
    offset = (page - 1) * limit
  }

  query := bson.M{}
  if filter != "" {
    query["name"] = bson.M{"$regex": filter, "$options": "i"}
  }

  options := options.Find()
  if sort != "" {
    options.SetSort(bson.D{{Key: sort, Value: 1}})
  }
  options.SetLimit(int64(limit))
  options.SetSkip(int64(offset))

  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  cursor, err := collection.Find(ctx, query, options)
  if err != nil {
    http.Error(w, "Failed to fetch patients", http.StatusInternalServerError)
    return
  }
  defer cursor.Close(ctx)

  var patients []Patient
  if err := cursor.All(ctx, &patients); err != nil {
    http.Error(w, "Failed to decode patients", http.StatusInternalServerError)
    return
  }

  totalCount, _ := collection.CountDocuments(ctx, query)

  response := map[string]interface{}{
    "results":    patients,
    "totalCount": totalCount,
  }

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(response)
}


func addPatientHandler(w http.ResponseWriter, r *http.Request) {
  var patient Patient
  if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
    http.Error(w, "Invalid request payload", http.StatusBadRequest)
    return
  }

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  result, err := collection.InsertOne(ctx, patient)
  if err != nil {
    http.Error(w, "Failed to add patient", http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(result)
}
func deletePatientHandler(w http.ResponseWriter, r *http.Request) {
  id := mux.Vars(r)["id"]
  objID, err := primitive.ObjectIDFromHex(id)
  if err != nil {
    http.Error(w, "Invalid ObjectID", http.StatusBadRequest)
    return
  }

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  _, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
  if err != nil {
    http.Error(w, "Failed to delete patient", http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(map[string]string{"message": "Patient deleted successfully"})
}

func updatePatientHandler(w http.ResponseWriter, r *http.Request) {
  id := mux.Vars(r)["id"]
  objID, err := primitive.ObjectIDFromHex(id)
  if err != nil {
    http.Error(w, "Invalid ObjectID", http.StatusBadRequest)
    return
  }

  var patient Patient
  if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
    http.Error(w, "Invalid request payload", http.StatusBadRequest)
    return
  }

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  update := bson.M{"$set": patient}
  _, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
  if err != nil {
    http.Error(w, "Failed to update patient", http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(map[string]string{"message": "Patient updated successfully"})
}

func searchPatientsHandler(w http.ResponseWriter, r *http.Request) {
  query := r.URL.Query().Get("name")
  if query == "" {
    http.Error(w, "Missing search query", http.StatusBadRequest)
    return
  }

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

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

func LogUserAction(log *logrus.Logger, action string, userID string, page string, additionalFields map[string]interface{}) {
	// Create a structured log entry
	entry := log.WithFields(logrus.Fields{
		"timestamp": time.Now().Format(time.RFC3339),
		"user_id":   userID,
		"action":    action,
		"page":      page,
	})

	// Add any additional fields if provided
	if additionalFields != nil {
		for key, value := range additionalFields {
			entry = entry.WithField(key, value)
		}
	}

	// Log the action
	entry.Info("User action logged")
}