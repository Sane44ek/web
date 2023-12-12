package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Response struct {
	Result string          `json:"result"`
	Data   json.RawMessage `json:"data"`
	Error  string          `json:"error"`
}

type Record struct {
	ID         int64  `json:"-" sql.field:"id"`
	Name       string `json:"name,omitempty" sql.field:"name"`
	LastName   string `json:"last_name,omitempty" sql.field:"last_name"`
	MiddleName string `json:"middle_name,omitempty" sql.field:"middle_name"`
	Address    string `json:"address,omitempty" sql.field:"address"`
	Phone      string `json:"phone,omitempty" sql.field:"phone"`
}

func connentToServer() {
	for {
		var command int
		fmt.Print("Выберите, что хотите сделать [1 - Создать запись, 2 - Обновить запись, 3 - Найти запись, 4 - Удалить запись]: ")
		fmt.Scanln(&command)
		if command == 1 {
			confirm := "yes"
			fmt.Print("Эта функция создаёт запись в адресной книге. Хотите продолжить? [yes]: ")
			fmt.Scanln(&confirm)
			if confirm != "yes" {
				continue
			}
			rec := &Record{}
			fmt.Print("Введите имя: ")
			fmt.Scanln(&rec.Name)
			fmt.Print("Введите фамилию: ")
			fmt.Scanln(&rec.LastName)
			fmt.Print("Введите отчество (при наличии): ")
			fmt.Scanln(&rec.MiddleName)
			fmt.Print("Введите адрес: ")
			fmt.Scanln(&rec.Address)
			fmt.Print("Введите номер телефона: ")
			fmt.Scanln(&rec.Phone)
			createRecord(rec)
			continue
		}
		if command == 2 {
			confirm := "yes"
			fmt.Print("Эта функция обновяет запись в адресной книге. Хотите продолжить? [yes]: ")
			fmt.Scanln(&confirm)
			if confirm != "yes" {
				continue
			}
			rec := &Record{}
			fmt.Print("Введите номер телефона, для которого обновляете данные: ")
			fmt.Scanln(&rec.Phone)
			fmt.Println("Введите те данные, которые хотите обновить. Остальные пропускайте")
			fmt.Print("Введите имя: ")
			fmt.Scanln(&rec.Name)
			fmt.Print("Введите фамилию: ")
			fmt.Scanln(&rec.LastName)
			fmt.Print("Введите отчество: ")
			fmt.Scanln(&rec.MiddleName)
			fmt.Print("Введите адрес: ")
			fmt.Scanln(&rec.Address)
			updateRecord(rec)
			continue
		}
		if command == 3 {
			confirm := "yes"
			fmt.Print("Эта функция поиска записи в адресной книге по известным данным. Хотите продолжить? [yes]: ")
			fmt.Scanln(&confirm)
			if confirm != "yes" {
				continue
			}
			rec := &Record{}
			fmt.Println("Введите те данные, по которым хотите найти запись")
			fmt.Print("Имя: ")
			fmt.Scanln(&rec.Name)
			fmt.Print("Фамилия: ")
			fmt.Scanln(&rec.LastName)
			fmt.Print("Отчество: ")
			fmt.Scanln(&rec.MiddleName)
			fmt.Print("Адрес: ")
			fmt.Scanln(&rec.Address)
			fmt.Print("Номер телефона: ")
			fmt.Scanln(&rec.Phone)
			getRecords(rec)
			continue
		}
		if command == 4 {
			confirm := "yes"
			fmt.Print("Эта функция удаления записи из адресной книги. Хотите продолжить? [yes]: ")
			fmt.Scanln(&confirm)
			if confirm != "yes" {
				continue
			}
			var phone []byte
			fmt.Print("Введите номер телефона, который хотите удалить из адресной книги: ")
			fmt.Scanln(&phone)
			deleteRecord(phone)
			continue
		}
		fmt.Println("Пожалуйста, введите корректную команду")
	}
}

func createRecord(rec *Record) {
	jsonData, err := json.Marshal(rec)
	if err != nil {
		log.Println("Error encoding JSON:", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/create", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error sending POST request:", err)
		return
	}
	defer resp.Body.Close()

	recordsData := make([]byte, resp.ContentLength)
	_, err = io.ReadFull(resp.Body, recordsData)
	if err != nil {
		log.Println("io.ReadFull(resp.Body, recordsData):", err)
		return
	}

	fmt.Println("Response Status:", resp.Status)
	var response Response
	err = json.Unmarshal(recordsData, &response)
	if err != nil {
		log.Println("json.Unmarshal(recordsData, &response):", err)
		return
	}
	if response.Result == "Error" {
		log.Println(response.Error)
		return
	}
	fmt.Println("Record successfully created")
}

func updateRecord(rec *Record) {
	jsonData, err := json.Marshal(rec)
	if err != nil {
		log.Println("Error encoding JSON:", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/update", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error sending POST request:", err)
		return
	}
	defer resp.Body.Close()
	recordsData := make([]byte, resp.ContentLength)
	_, err = io.ReadFull(resp.Body, recordsData)
	if err != nil {
		log.Println("io.ReadFull(resp.Body, recordsData):", err)
		return
	}

	fmt.Println("Response Status:", resp.Status)
	var response Response
	err = json.Unmarshal(recordsData, &response)
	if err != nil {
		log.Println("json.Unmarshal(recordsData, &response):", err)
		return
	}
	if response.Result == "Error" {
		log.Println(response.Error)
		return
	}
	fmt.Println("Record successfully updated")
}

func deleteRecord(phone []byte) {
	resp, err := http.Post("http://localhost:8080/delete", "application/text", bytes.NewBuffer(phone))
	if err != nil {
		log.Println("Error sending POST request:", err)
		return
	}
	defer resp.Body.Close()

	recordsData := make([]byte, resp.ContentLength)
	_, err = io.ReadFull(resp.Body, recordsData)
	if err != nil {
		log.Println("io.ReadFull(resp.Body, recordsData):", err)
		return
	}

	fmt.Println("Response Status:", resp.Status)
	var response Response
	err = json.Unmarshal(recordsData, &response)
	if err != nil {
		log.Println("json.Unmarshal(recordsData, &response):", err)
		return
	}
	if response.Result == "Error" {
		log.Println(response.Error)
		return
	}
	fmt.Println("Record successfully deleted")
}

func getRecords(rec *Record) {
	jsonData, err := json.Marshal(rec)
	if err != nil {
		log.Println("Error encoding JSON:", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/get", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error sending POST request:", err)
		return
	}
	defer resp.Body.Close()

	recordsData := make([]byte, resp.ContentLength)
	_, err = io.ReadFull(resp.Body, recordsData)
	if err != nil {
		log.Println("io.ReadFull(resp.Body, recordsData):", err)
		return
	}
	var records []Record
	var response Response
	err = json.Unmarshal(recordsData, &response)
	if err != nil {
		log.Println("json.Unmarshal(recordsData, &response):", err)
		return
	}
	err = json.Unmarshal(response.Data, &records)
	if err != nil {
		log.Println("json.Unmarshal(response.Data, &records):", err)
		return
	}
	if response.Result == "Error" {
		log.Println(response.Error)
		return
	}
	fmt.Println("Result: ")
	for _, record := range records {
		fmt.Println("-->")
		fmt.Println("\tName:" + record.Name)
		fmt.Println("\tSurname:" + record.MiddleName)
		if record.LastName != "" {
			fmt.Println("\tLastname:" + record.LastName)
		}
		fmt.Println("\tAddtess:" + record.Address)
		fmt.Println("\tPhone number:" + record.Phone)
	}
}

func main() {
	connentToServer()
}
