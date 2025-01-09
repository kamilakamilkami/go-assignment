// @ts-ignore
document.getElementById('submitButton').addEventListener('click', function(event) {
    event.preventDefault();
  
    // @ts-ignore
    const name = document.getElementById('name').value;
    // @ts-ignore
    const email = document.getElementById('email').value;
    // @ts-ignore
    const age = document.getElementById('age').value;
  
    const patientData = {
      name: name,
      email: email,
      age: age
    };
  
    console.log('Sending patient data:', patientData);
  
    fetch('http://localhost:8080/add-patient', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    },
    body: JSON.stringify({
        name: "John Doe",
        email: "john.doe@example.com",
        phone: "1234567890",
        address: "123 Street Name"
    })
})
.then(response => {
    if (response.ok) {
        return response.json();
    } else {
        console.log('Error: ', response.statusText);
        throw new Error('Error submitting patient data');
    }
})
.then(data => {
    console.log('Success:', data);
})
.catch(error => {
    console.error('There was an error!', error);
});

  });
  
  