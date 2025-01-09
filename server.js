const express = require('express');
const mongoose = require('mongoose');
const bodyParser = require('body-parser');
const cors = require('cors');

const app = express();

// Middleware для CORS и парсинга JSON
app.use(cors());
app.use(bodyParser.json());

// Подключение к MongoDB
mongoose.connect('mongodb://127.0.0.1:27017/med_portal', { useNewUrlParser: true, useUnifiedTopology: true })
  .then(() => console.log('Connected to MongoDB'))
  .catch((err) => console.log('Error connecting to MongoDB: ', err));

// Создание схемы для пациента
const patientSchema = new mongoose.Schema({
  name: String,
  email: String,
  phone: String,
  address: String
});

// Модель для пациента
const Patient = mongoose.model('Patient', patientSchema);

// Обработчик POST-запроса для добавления пациента
app.post('/add-patient', async (req, res) => {
  try {
    const { name, email, phone, address } = req.body;
    const newPatient = new Patient({ name, email, phone, address });

    await newPatient.save();
    res.status(201).send('Patient added successfully!');
  } catch (error) {
    console.error('Error adding patient:', error);
    res.status(500).send('Error adding patient');
  }
});

// Запуск сервера
const port = 8080;
app.listen(port, () => {
  console.log(`Server running on http://localhost:${port}`);
});
