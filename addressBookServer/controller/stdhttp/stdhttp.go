package stdhttp

import (
	"encoding/json"
	"httpserver/gates/psg"
	"httpserver/models/dto"
	"httpserver/pkg"
	"io"
	"log"
	"net/http"
)

type Controller struct {
	srv http.Server
	db  *psg.Psg
}

func NewController(addr string, postgres *psg.Psg) (hs *Controller) {
	hs = new(Controller)
	hs.srv = http.Server{}
	mux := http.NewServeMux()

	mux.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Println("Method Not Allowed", http.StatusMethodNotAllowed)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		hs.RecordCreateHandler(w, r)
	})
	mux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Println("Method Not Allowed", http.StatusMethodNotAllowed)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		hs.RecordGetHandler(w, r)
	})
	mux.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Println("Method Not Allowed", http.StatusMethodNotAllowed)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		hs.RecordUpdateHandler(w, r)
	})
	mux.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Println("Method Not Allowed", http.StatusMethodNotAllowed)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		hs.RecordDeleteByPhoneHandler(w, r)
	})
	hs.srv.Handler = mux
	hs.srv.Addr = addr
	hs.db = postgres
	return hs
}

func (hs *Controller) Start() error {
	myErr := pkg.NewMyError("package stdhttp: func (hs *Controller) Start()")
	err := hs.srv.ListenAndServe()
	if err != nil {
		log.Fatal(myErr.Wrap(err, "hs.srv.ListenAndServe()").Error())
		return myErr.Wrap(err, "hs.srv.ListenAndServe()")
	}
	return nil
}

// RecordCreate обрабатывает HTTP запрос для добавления новой записи.
func (hs *Controller) RecordCreateHandler(w http.ResponseWriter, req *http.Request) {
	myErr := pkg.NewMyError("package stdhttp: func (hs *Controller) RecordCreateHandler(w http.ResponseWriter, req *http.Request)")
	w.WriteHeader(http.StatusOK) // Status 200 OK
	response := dto.Response{}
	record, err := GetBody(req)
	if err != nil {
		e := myErr.Wrap(err, "GetBody(req)")
		response.Result = "Error"
		response.Error = e.Error()
		js, erro := response.GetJson()
		if erro != nil {
			log.Println(erro.Error())
			w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
			return
		}
		w.Write(js)
		log.Println(e.Error())
		return
	}
	/// Проверяем нет ли пустых значений
	if record.Address == "" || record.LastName == "" || record.Name == "" || record.Phone == "" {
		e := myErr.Wrap(nil, "All fields except for the middlename must be filled in")
		response.Result = "Error"
		response.Error = e.Error()
		js, erro := response.GetJson()
		if erro != nil {
			log.Println(erro.Error())
			w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
			return
		}
		w.Write(js)
		log.Println(e.Error())
		return
	}

	record.Phone, err = pkg.PhoneNormalize(record.Phone)
	if err != nil {
		e := myErr.Wrap(err, "pkg.PhoneNormalize(record.Phone)")
		response.Result = "Error"
		response.Error = e.Error()
		js, erro := response.GetJson()
		if erro != nil {
			log.Println(erro.Error())
			w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
			return
		}
		w.Write(js)
		log.Println(e.Error())
		return
	}
	err = hs.db.RecordCreate(*record)
	if err != nil {
		e := myErr.Wrap(err, "")
		response.Result = "Error"
		response.Error = e.Error()
		js, erro := response.GetJson()
		if erro != nil {
			log.Println(erro.Error())
			w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
			return
		}
		w.Write(js)
		log.Println(e.Error())
		return
	}
	response.Result = "Success"
	js, err := response.GetJson()
	if err != nil {
		log.Println(err.Error())
		w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + err.Error() + `"}`))
		return
	}
	w.Write(js)
}

// RecordsGet обрабатывает HTTP запрос для получения записей на основе предоставленных полей Record.
func (hs *Controller) RecordGetHandler(w http.ResponseWriter, req *http.Request) {
	myErr := pkg.NewMyError("package stdhttp: func (hs *Controller) RecordGetHandler(w http.ResponseWriter, req *http.Request)")
	w.WriteHeader(http.StatusOK) // Status 200 OK
	response := dto.Response{}
	record, err := GetBody(req)
	if err != nil {
		e := myErr.Wrap(err, "GetBody(req)")
		response.Result = "Error"
		response.Error = e.Error()
		js, erro := response.GetJson()
		if erro != nil {
			log.Println(erro.Error())
			w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
			return
		}
		w.Write(js)
		log.Println(e.Error())
		return
	}
	if record.Phone != "" {
		record.Phone, err = pkg.PhoneNormalize(record.Phone)
		if err != nil {
			e := myErr.Wrap(err, "pkg.PhoneNormalize(record.Phone)")
			response.Result = "Error"
			response.Error = e.Error()
			js, erro := response.GetJson()
			if erro != nil {
				log.Println(erro.Error())
				w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
				return
			}
			w.Write(js)
			log.Println(e.Error())
			return
		}
	}

	result, err := hs.db.RecordsGet(*record)
	if err != nil {
		e := myErr.Wrap(err, "")
		response.Result = "Error"
		response.Error = e.Error()
		js, erro := response.GetJson()
		if erro != nil {
			log.Println(erro.Error())
			w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
			return
		}
		w.Write(js)
		log.Println(e.Error())
		return
	}

	/// Формируем ответ
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(result)
	if err != nil {
		e := myErr.Wrap(err, "json.Marshal(result)")
		response.Result = "Error"
		response.Error = e.Error()
		js, erro := response.GetJson()
		if erro != nil {
			log.Println(erro.Error())
			w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
			return
		}
		w.Write(js)
		log.Println(e.Error())
		return
	}
	response.Result = "Success"
	response.Data = jsonData
	js, err := response.GetJson()
	if err != nil {
		log.Println(err.Error())
		w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + err.Error() + `"}`))
		return
	}
	w.Write(js)
}

// RecordUpdate обрабатывает HTTP запрос для обновления записи.
func (hs *Controller) RecordUpdateHandler(w http.ResponseWriter, req *http.Request) {
	myErr := pkg.NewMyError("package stdhttp: func (hs *Controller) RecordUpdateHandler(w http.ResponseWriter, req *http.Request)")
	w.WriteHeader(http.StatusOK) // Status 200 OK
	response := dto.Response{}
	record, err := GetBody(req)
	if err != nil {
		e := myErr.Wrap(err, "GetBody(req)")
		response.Result = "Error"
		response.Error = e.Error()
		js, erro := response.GetJson()
		if erro != nil {
			log.Println(erro.Error())
			w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
			return
		}
		w.Write(js)
		log.Println(e.Error())
		return
	}

	record.Phone, err = pkg.PhoneNormalize(record.Phone)
	if err != nil {
		e := myErr.Wrap(err, "pkg.PhoneNormalize(record.Phone)")
		response.Result = "Error"
		response.Error = e.Error()
		js, erro := response.GetJson()
		if erro != nil {
			log.Println(erro.Error())
			w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
			return
		}
		w.Write(js)
		log.Println(e.Error())
		return
	}
	err = hs.db.RecordUpdate(*record)
	if err != nil {
		e := myErr.Wrap(err, "")
		response.Result = "Error"
		response.Error = e.Error()
		js, erro := response.GetJson()
		if erro != nil {
			log.Println(erro.Error())
			w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
			return
		}
		w.Write(js)
		log.Println(e.Error())
		return
	}
	response.Result = "Success"
	js, err := response.GetJson()
	if err != nil {
		log.Println(err.Error())
		w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + err.Error() + `"}`))
		return
	}
	w.Write(js)
}

// / RecordDeleteByPhone обрабатывает HTTP запрос для удаления записи по номеру телефона.
func (hs *Controller) RecordDeleteByPhoneHandler(w http.ResponseWriter, req *http.Request) {
	myErr := pkg.NewMyError("package stdhttp: func (hs *Controller) RecordDeleteByPhoneHandler(w http.ResponseWriter, req *http.Request)")
	w.WriteHeader(http.StatusOK) /// Status 200 OK
	response := dto.Response{}
	byteReq, err := io.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		e := myErr.Wrap(err, "io.ReadAll(req.Body)")
		response.Result = "Error"
		response.Error = e.Error()
		js, erro := response.GetJson()
		if erro != nil {
			log.Println(erro.Error())
			w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
			return
		}
		w.Write(js)
		log.Println(e.Error())
		return
	}
	phone, err := pkg.PhoneNormalize(string(byteReq))
	if err != nil {
		e := myErr.Wrap(err, "pkg.PhoneNormalize(string(byteReq))")
		response.Result = "Error"
		response.Error = e.Error()
		js, erro := response.GetJson()
		if erro != nil {
			log.Println(erro.Error())
			w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
			return
		}
		w.Write(js)
		log.Println(e.Error())
		return
	}
	err = hs.db.RecordDeleteByPhone(phone)
	if err != nil {
		e := myErr.Wrap(err, "")
		response.Result = "Error"
		response.Error = e.Error()
		js, erro := response.GetJson()
		if erro != nil {
			log.Println(erro.Error())
			w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
			return
		}
		w.Write(js)
		log.Println(e.Error())
		return
	}
	response.Result = "Success"
	js, err := response.GetJson()
	if err != nil {
		log.Println(err.Error())
		w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + err.Error() + `"}`))
		return
	}
	w.Write(js)
}

// / Функция для преобразования response в структуру Record
func GetBody(req *http.Request) (record *dto.Record, e error) {
	myErr := pkg.NewMyError("package stdhttp: func GetBody(req *http.Request)")
	byteReq, err := io.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		e = myErr.Wrap(err, "io.ReadAll(req.Body)")
		log.Println(e.Error())
		return nil, e
	}
	err = json.Unmarshal(byteReq, &record)
	if err != nil {
		e = myErr.Wrap(err, "json.Unmarshal(byteReq, &record)")
		log.Println(e.Error())
		return nil, e
	}
	return record, nil
}
