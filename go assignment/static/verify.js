const API_URL = "http://localhost:8080";

// Select the verification form
const form = document.querySelector('.verification');
const token = localStorage.getItem("auth_token");

form.addEventListener('submit', async (event) => {
    // Prevent the default form submission behavior
    event.preventDefault();

    const code = document.getElementById('code').value;
    console.log("hi");
    try {
        const response = await fetch(`${API_URL}/verify`, {
            method: "POST",
            headers: { "Content-Type": "application/json", "Authorization": token },
            body: JSON.stringify({ code }),
        });
        if (response.ok) {
            const data = await response.json();
            const token = data.token;  // Assuming the token is returned in the response body
            localStorage.setItem("auth_token", token);  // Store the token in localStorage
            window.location.href = "profile.html"; // Redirect
        } else {
            console.error("Failed verification");
        }
    } catch (error) {
        console.error("Error during verification:");
        alert("Invalid verification code. Please try again.");
    }
});
