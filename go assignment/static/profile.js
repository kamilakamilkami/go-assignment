const API_URL = "http://localhost:8080";

// Обновление профиля пользователя
// document.getElementById("update-profile-form").addEventListener("submit", async (event) => {
//     event.preventDefault();
//     const name = document.getElementById("profile-name").value.trim();
//     const email = document.getElementById("profile-email").value.trim();
//     const password = document.getElementById("profile-password").value.trim();

//     if (!name || !email) {
//         alert("Name and email are required!");
//         return;
//     }

//     try {
//         const response = await fetch(`${API_URL}/update-profile`, {
//             method: "PUT",
//             headers: { "Content-Type": "application/json" },
//             body: JSON.stringify({ name, email, password }),
//         });

//         if (!response.ok) throw new Error("Failed to update profile");
//         alert("Profile updated successfully!");
//     } catch (error) {
//         console.error("Error updating profile:", error);
//     }
// });

document.getElementById("update-profile-form").addEventListener("submit", async (event) => {
    event.preventDefault();
    const email = document.getElementById("user-email").value.trim();
    const password = document.getElementById("password-profile").value.trim();

    if (!password || !email) {
        alert("password and email are required!");
        return;
    }

    try {
        const response = await fetch(`${API_URL}/login`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ email, password }),
        });

        if (response.ok) {
            const data = await response.json();
            const token = data.token;  // Assuming the token is returned in the response body
            localStorage.setItem("auth_token", token);  // Store the token in localStorage
            window.location.href = "admin.html"; // Redirect
        } else {
            console.error("Failed to register");
        }
    } catch (error) {
        console.error("Error login:", error);
    }
});


// Отправка сообщения в службу поддержки
document.getElementById("support-form").addEventListener("submit", async (event) => { 
    event.preventDefault(); 

    const subject = document.getElementById("subject").value.trim();
    const message = document.getElementById("support-message").value.trim(); 
    const fileInput = document.getElementById("attachment");

    if (!subject || !message) {
        alert("Please fill out all fields!");
        return;
    }

    const formData = new FormData();
    formData.append("subject", subject);
    formData.append("message", message);

    if (fileInput.files[0]) {
        formData.append("attachment", fileInput.files[0]);
    }

    try { 
        const response = await fetch(`${API_URL}/send-support-message`, { 
            method: "POST", 
            body: formData,
        }); 
 
        if (!response.ok) throw new Error("Failed to send support message"); 
        alert("Message sent successfully!"); 
    } catch (error) { 
        console.error("Error sending support message:", error); 
        alert("There was an error sending your message. Please try again.");
    } 
});

// // Загрузка истории взаимодействий
// async function loadInteractionHistory() {
//     try {
//         const response = await fetch(`${API_URL}/get-interaction-history`);
//         if (!response.ok) throw new Error("Failed to fetch interaction history");

//         const data = await response.json();
//         const historyList = document.getElementById("interaction-history");
//         historyList.innerHTML = "";

//         data.forEach((item) => {
//             const li = document.createElement("li");
//             li.textContent = `${item.date}: ${item.details}`;
//             historyList.appendChild(li);
//         });
//     } catch (error) {
//         console.error("Error fetching interaction history:", error);
//     }
// }

// Загрузка данных при открытии страницы
// loadInteractionHistory();


