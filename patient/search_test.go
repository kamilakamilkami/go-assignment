package patient

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "reflect"
    "testing"
)

func TestGetProduct(t *testing.T) {
    request, _ := http.NewRequest("GET", "/search-patients/?name=liya", nil)
    response := httptest.NewRecorder()

    SearchPatientsHandler(response, request)

    if response.Code != http.StatusOK {
        t.Errorf("Incorrect status code. Expected: %d, Got: %d", http.StatusOK, response.Code)
    }

    // Define the expected result as a Go struct
    var expected []map[string]interface{}
    var actual []map[string]interface{}

    expectedJSON := `[{"_id":"67866874324804643adc9b53","name":"liya","email":"liya@gmail.com","phone":".222665652","address":"abylkhair","role":""}]`

    // Unmarshal JSON strings to comparable structures
    err := json.Unmarshal([]byte(expectedJSON), &expected)
    if err != nil {
        t.Fatalf("Failed to unmarshal expected JSON: %v", err)
    }

    err = json.Unmarshal(response.Body.Bytes(), &actual)
    if err != nil {
        t.Fatalf("Failed to unmarshal actual JSON: %v", err)
    }

    // Use reflect.DeepEqual to compare the data structures
    if !reflect.DeepEqual(expected, actual) {
        t.Errorf("Incorrect response body.\nExpected: %+v\nGot: %+v", expected, actual)
    }
}
