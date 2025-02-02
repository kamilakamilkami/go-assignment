const API_URL = "http://localhost:8080";

const form = document.querySelector('.signup-form');

// Add an event listener for the form submission
form.addEventListener('submit', async (event) => {
    // Prevent the default form submission behavior
    event.preventDefault();

    // Get the values from the input fields
    const username = document.getElementById('signup-username').value;
    const email = document.getElementById('signup-email').value;
    const password = document.getElementById('signup-password').value;


    try {
        const response = await fetch(`${API_URL}/register`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ username, email, password }),
            credentials: "include",  
        });

        if (response.ok) {
            const data = await response.json();
            const token = data.token;  // Assuming the token is returned in the response body
            localStorage.setItem("auth_token", token);  // Store the token in localStorage
            window.location.href = "verification.html"; // Redirect
        } else {
            console.error("Failed to register");
        }
        
        
    } catch (error) {
        console.error("Error register:", error);
    }
});