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
  "math/rand"
  "med_portal/patient"
	"med_portal/svc"
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
  
type RegisterForm struct {
    Name    string             `json:"name" bson:"name"`
    Email   string             `json:"email" bson:"email"`
    Password string            `json:"password" bson:"password"`
    Code   string             `json:"code" bson:"code"`
    Role string            `json:"role" bson:"role"`
    
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
  router.HandleFunc("/search-patients", patient.SearchPatientsHandler).Methods("GET")
  router.HandleFunc("/send-email", sendEmailGet).Methods("GET")
  router.HandleFunc("/send-support-message", sendEmailPost).Methods("POST")
  
  router.HandleFunc("/register", registerGet).Methods("GET")
  router.HandleFunc("/register", registerPost).Methods("POST")
  router.HandleFunc("/verify", VerifyCode).Methods("POST")
  router.HandleFunc("/login", Login).Methods("POST")
  router.HandleFunc("/getrole", GetRole).Methods("POST")

  // Enable CORS
  corsHandler := cors.New(cors.Options{
    AllowedOrigins:   []string{"http://127.0.0.1:5502", "http://localhost:5502", "http://127.0.0.1:5501", "http://localhost:5501"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"Content-Type", "Authorization", "Set-Cookie"},
    AllowCredentials: true,
  })

  log.Println("Server running on http://localhost:8080")
  log.Fatal(http.ListenAndServe(":8080", corsHandler.Handler(router)))
}
// I added
func sendEmailGet(w http.ResponseWriter, r *http.Request) {

}

func registerGet(w http.ResponseWriter, r *http.Request) {

}

func Login(w http.ResponseWriter, r *http.Request) {
  var data map[string]string
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&data)
  if err != nil {
      http.Error(w, "Invalid JSON format", http.StatusBadRequest)
      return
  }

  email, password := data["email"], data["password"]
  user, err := ExistUser(client, "med_portal", "users", email, password)
  if err != nil {
    http.Error(w, "NO such user", http.StatusBadRequest)
    return
  }
  tokenString, err := svc.CreateToken(user.Email, user.Password, user.Role)
  if err != nil {
      http.Error(w, "Error creating token", http.StatusInternalServerError)
      return
  }
  SetAuthCookie(w, tokenString)

  // Send the token in the response body
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(map[string]interface{}{
      "token": tokenString,  // Send the JWT token
  })
}

func GetRole(w http.ResponseWriter, r *http.Request) {
  token := r.Header.Get("Authorization")
  if token == "" {
      http.Error(w, "Authorization token is missing", http.StatusUnauthorized)
      return
  }
  err := svc.VerifyToken(token)
  if err != nil {
    http.Error(w, "Token is not verified", http.StatusBadRequest)
    return
  }
  role, _ := svc.GetClaim(token, "role")
  fmt.Println(role)
  w.Header().Set("Content-Type", "application/json")
      json.NewEncoder(w).Encode(map[string]interface{}{
          "role": role,  
      })
}

func registerPost(w http.ResponseWriter, r *http.Request) {
  var register RegisterForm
  register.Role = "not verified"
  register.Code = GenerateRandomCode(4)
  if err := json.NewDecoder(r.Body).Decode(&register); err != nil {
    http.Error(w, "Invalid request payload", http.StatusBadRequest)
    return
  }

  CreateUser(client, context.TODO(), "med_portal", "users", register)
  SendEmail("Your Verification Code", register.Code, register.Email)

  tokenString, err := svc.CreateToken(register.Email, register.Password, register.Role)
  if err != nil {
    http.Error(w, "Error creating token", http.StatusInternalServerError)
    return
  }

  SetAuthCookie(w, tokenString)
  w.Header().Set("Authorization", tokenString)
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(map[string]interface{}{
      "token": tokenString, // Send the JWT token
  })

}

func VerifyCode(w http.ResponseWriter, r *http.Request) {
  token := r.Header.Get("Authorization")
  if token == "" {
      http.Error(w, "Authorization token is missing", http.StatusUnauthorized)
      return
  }
  email, _ := svc.GetClaim(token, "email")
  password, _ := svc.GetClaim(token, "password")
  user, _ := ExistUser(client, "med_portal", "users", email, password)

  var data map[string]string
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&data)
  if err != nil {
      http.Error(w, "Invalid JSON format", http.StatusBadRequest)
      return
  }

  code := data["code"]
  if user.Code == code {
      tokenString, err := svc.CreateToken(user.Email, user.Password, "user")
      if err != nil {
          http.Error(w, "Error creating token", http.StatusInternalServerError)
          return
      }

      updateRole(client, context.TODO(), "med_portal", "users", user, "user")
      SetAuthCookie(w, tokenString)

      // Send the token in the response body
      w.Header().Set("Content-Type", "application/json")
      json.NewEncoder(w).Encode(map[string]interface{}{
          "token": tokenString,  // Send the JWT token
      })
      return
  }

  http.Error(w, "Invalid verification code", http.StatusUnauthorized)
}


func ExistUser(client *mongo.Client, database, collection, email, password string) (RegisterForm, error) {
	user := GetUsers(client, database, collection, bson.M{"email": email, "password": password}, bson.D{})
	if len(user) != 0 {
		return user[0], nil
	}
	return RegisterForm{}, fmt.Errorf("NO user")
}


func GenerateRandomCode(length int) string {
	if length <= 0 {
		return ""
	}

	rand.Seed(time.Now().UnixNano())
	code := ""
	for i := 0; i < length; i++ {
		code += fmt.Sprintf("%d", rand.Intn(10)) // Append a random digit (0-9)
	}
	return code
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

func CreateUser(client *mongo.Client, ctx context.Context, dataBase, col string, user RegisterForm) (*mongo.InsertOneResult, error) {
	collection := client.Database(dataBase).Collection(col)
  result, err := collection.InsertOne(ctx, user)
  return result, err
}

func SetAuthCookie(w http.ResponseWriter, tokenString string) {
  cookie := &http.Cookie{
      Name:     "auth_token",
      Value:    tokenString,
      HttpOnly: true,
      Path:     "/",
      Expires:  time.Now().Add(24 * time.Hour),
      SameSite: http.SameSiteStrictMode, 
  }
  fmt.Println("Setting auth_token cookie:", cookie)
  http.SetCookie(w, cookie)
}





func GetUsers(client *mongo.Client, database, collection string, filter bson.M, sorting bson.D) []RegisterForm {
  coll := client.Database(database).Collection(collection)

  findOptions := options.Find().SetSort(sorting)

  cursor, err := coll.Find(context.TODO(), filter, findOptions)
  if err != nil {
      panic(err)
  }

  var users []RegisterForm
  if err := cursor.All(context.TODO(), &users); err != nil {
      panic(err)
  }

  return users
}

func updateRole(client *mongo.Client, ctx context.Context, dataBase, col string, user RegisterForm, role string) error {
	collection := client.Database(dataBase).Collection(col)
	filter := bson.D{{"email", user.Email}, {"password", user.Password}}
	update := bson.D{
		{"$set", bson.D{
			{"role", role},
		}},
	}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println("failed to update product")
		return err
	}

	// Check if the product was found and updated
	if result.MatchedCount == 0 {
		return err
	}

	fmt.Printf("Successfully updated %d user's role\n", result.ModifiedCount)
	return nil
}