const API_URL = "http://localhost:8080";

document.getElementById("add-patient-form").addEventListener("submit", addPatient);

async function fetchPatients(filter = "", sort = "", page = 1) {
    try {
        const url = new URL(`${API_URL}/get-all-patients`);
        const params = { filter, sort, page };
        Object.keys(params).forEach(key => url.searchParams.append(key, params[key]));

        const response = await fetch(url, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            },
        });

        if (!response.ok) {
            throw new Error("Failed to fetch patients");
        }

        const data = await response.json();
        const patientList = document.getElementById("patient-list");
        patientList.innerHTML = "";

        data.results.forEach(patient => {
            const li = document.createElement("li");
            li.innerHTML = `
                <span>${patient.name} - ${patient.email}</span>
                <button onclick="deletePatient('${patient._id}')">Delete</button>
                <button onclick="editPatient('${patient._id}')">Edit</button>
            `;
            patientList.appendChild(li);
        });

        updatePagination(data.totalCount, page);
    } catch (error) {
        console.error("Error fetching patients:", error);
    }
}

async function addPatient(event) {
    event.preventDefault();
    const name = document.getElementById("name").value.trim();
    const email = document.getElementById("email").value.trim();
    const phone = document.getElementById("phone").value.trim();
    const address = document.getElementById("address").value.trim();

    if (!name || !email || !phone || !address) {
        alert("All fields are required!");
        return;
    }

    try {
        const response = await fetch(`${API_URL}/add-patient`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ name, email, phone, address }),
        });

        if (!response.ok) {
            throw new Error("Failed to add patient");
        }

        console.log("Patient added:", await response.json());
        fetchPatients();
    } catch (error) {
        console.error("Error adding patient:", error);
    }
}

async function deletePatient(id) {
    try {
        const response = await fetch(`${API_URL}/delete-patient/${id}`, {
            method: "DELETE",
            headers: {
                "Content-Type": "application/json",
            },
        });

        if (!response.ok) {
            throw new Error("Failed to delete patient");
        }

        console.log("Patient deleted:", await response.json());
        fetchPatients();
    } catch (error) {
        console.error("Error deleting patient:", error);
    }
}

async function editPatient(id) {
    const name = prompt("Enter new name:").trim();
    const email = prompt("Enter new email:").trim();
    const phone = prompt("Enter new phone:").trim();
    const address = prompt("Enter new address:").trim();

    if (!name || !email || !phone || !address) {
        alert("All fields are required!");
        return;
    }

    try {
        const response = await fetch(`${API_URL}/update-patient/${id}`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ name, email, phone, address }),
        });

        if (!response.ok) {
            throw new Error("Failed to update patient");
        }

        console.log("Patient updated:", await response.json());
        fetchPatients();
    } catch (error) {
        console.error("Error updating patient:", error);
    }
}

async function searchPatients() {
    const searchInput = document.getElementById("search").value.trim();
    fetchPatients(searchInput);
}

function updatePagination(totalCount, currentPage) {
    const itemsPerPage = 10;
    const totalPages = Math.ceil(totalCount / itemsPerPage);
    const pagination = document.getElementById("pagination");
    pagination.innerHTML = "";

    for (let i = 1; i <= totalPages; i++) {
        const button = document.createElement("button");
        button.textContent = i;
        button.disabled = i === currentPage;
        button.addEventListener("click", () => fetchPatients("", "", i));
        pagination.appendChild(button);
    }
}

document.getElementById("filter-form").addEventListener("submit", event => {
    event.preventDefault();
    const filter = document.getElementById("filter").value.trim();
    const sort = document.getElementById("sort").value.trim();
    fetchPatients(filter, sort);
});

// Initial fetch of patients
fetchPatients();
