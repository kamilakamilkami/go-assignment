
var token = localStorage.getItem("auth_token");

if (!token) {
    console.error("Authorization token not found in localStorage.");
    window.location.href = "profile.html"; // Redirect to login if no token is present
} else {
    try {
        fetch(`http://localhost:8080/getrole`, {
            method: "POST",
            headers: { 
                "Content-Type": "application/json", 
                "Authorization": token 
            },
        })
        .then(response => {
            if (!response.ok) {
                throw new Error("Failed to fetch role: " + response.statusText);
            }
            return response.json();
        })
        .then(data => {
            const role = data.role; 
            if (role !== "admin") {
                console.log(role)
                window.location.href = "admin.html"; // Redirect if not an admin
            }
        })
        .catch(error => {
            console.error("Error during role verification:", error);
            alert("Failed to verify your role. Please log in again.");
            window.location.href = "profile.html"; // Redirect to login on failure
        });
    } catch (error) {
        console.error("Unexpected error:", error);
    }
}
