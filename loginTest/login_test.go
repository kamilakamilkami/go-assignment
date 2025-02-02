package login

import (
    "testing"
    "github.com/tebeka/selenium"
    "time"
)

func TestLoginPage(t *testing.T) {
    // Setup WebDriver capabilities for Microsoft Edge
    caps := selenium.Capabilities{"browserName": "MicrosoftEdge"}
    edgeCaps := selenium.Capabilities{
        "ms:edgeOptions": map[string]interface{}{
            "args": []string{"--disable-gpu"},
        },
    }

    // Merge the Edge-specific capabilities with the base capabilities
    for key, value := range edgeCaps {
        caps[key] = value
    }

    // Connect to Edge WebDriver running on port 5555
    wd, err := selenium.NewRemote(caps, "http://127.0.0.1:64484")
    if err != nil {
        t.Fatal(err)
    }
    defer wd.Quit()

    // Navigate to the login page
    err = wd.Get("http://127.0.0.1:5501/go%20assignment/static/profile.html")
    if err != nil {
        t.Fatal(err)
    }

    // Find and fill the email input
    emailInput, err := wd.FindElement(selenium.ByID, "user-email")
    if err != nil {
        t.Fatal(err)
    }
    err = emailInput.SendKeys("alisherdautov22@gmail.com")
    if err != nil {
        t.Fatal(err)
    }

    // Find and fill the password input
    passwordInput, err := wd.FindElement(selenium.ByID, "password-profile")
    if err != nil {
        t.Fatal(err)
    }
    err = passwordInput.SendKeys("abab")
    if err != nil {
        t.Fatal(err)
    }
    time.Sleep(5 * time.Second)
    // Find the submit button using the correct ID (login)
    submitButton, err := wd.FindElement(selenium.ByID, "submit")
    if err != nil {
        t.Fatal(err)
    }

    // Use JavaScript to click the submit button
    jsScript := "arguments[0].click();"
    _, err = wd.ExecuteScript(jsScript, []interface{}{submitButton}) // Capture both result and error
    if err != nil {
        t.Fatal("Failed to execute JavaScript click:", err)
    }

    // Allow time for form submission and page load
    time.Sleep(5 * time.Second)

    // Verify the resulting URL
    currentURL, err := wd.CurrentURL()
    if err != nil {
        t.Fatal(err)
    }

    // Adjust the expected URL based on your application’s behavior
    if currentURL != "http://127.0.0.1:5501/go%20assignment/static/index.html" {
        t.Errorf("Expected URL to be , but got '%s'", currentURL)
    }
}
