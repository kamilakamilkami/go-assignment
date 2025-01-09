package main

import (
	"gorm.io/gorm"
)

type Patient struct {
	gorm.Model
	Name    string
	Email   string
	Age     int
	Address string
}

func CreatePatient(db *gorm.DB, name string, email string, age int, address string) (*Patient, error) {
	patient := Patient{Name: name, Email: email, Age: age, Address: address}
	result := db.Create(&patient)
	return &patient, result.Error
}

func GetPatientByID(db *gorm.DB, id uint) (*Patient, error) {
	var patient Patient
	result := db.First(&patient, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &patient, nil
}

func UpdatePatient(db *gorm.DB, id uint, name string, email string, age int, address string) error {
	var patient Patient
	result := db.First(&patient, id)
	if result.Error != nil {
		return result.Error
	}
	patient.Name = name
	patient.Email = email
	patient.Age = age
	patient.Address = address
	return db.Save(&patient).Error
}

func DeletePatient(db *gorm.DB, id uint) error {
	result := db.Delete(&Patient{}, id)
	return result.Error
}

func GetAllPatients(db *gorm.DB) ([]Patient, error) {
	var patients []Patient
	result := db.Find(&patients)
	return patients, result.Error
}