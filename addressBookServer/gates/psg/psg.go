package psg

import (
	"context"
	"httpserver/models/dto"
	"httpserver/pkg"
	"log"
	"net/url"
	"strconv"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Psg struct {
	conn *pgxpool.Pool
}

// NewPsg создает новый экземпляр Psg.
func NewPsg(psgAddr string, dbUser string, dbPassword string) *Psg {
	psg := &Psg{}

	if dbPassword == "" || dbUser == "" || psgAddr == "" {
		log.Fatal("Error: NewPsg(psgAddr string, dbUser string, dbPassword string): please write all arguments")
	}
	db_url := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(dbUser, dbPassword),
		Host:   psgAddr + ":5432",
		Path:   "records",
	}
	connPool, err := pgxpool.New(context.Background(), db_url.String())

	if err != nil {
		log.Fatal(err, "Something going wrong with connecting to database")
	}
	psg.conn = connPool
	return psg
}

// RecordCreate добавляет новую запись в базу данных.
func (p *Psg) RecordCreate(record dto.Record) error {
	myErr := pkg.NewMyError("package pkg: func (p *Psg) RecordCreate(record dto.Record) error")
	record_existence, err := p.CheckPhone(record.Phone)
	if err != nil {
		log.Println(myErr.Wrap(err, ""))
		return myErr.Wrap(err, "")
	}
	if record_existence {
		log.Println(myErr.Wrap(nil, "Phone already exist"))
		return myErr.Wrap(nil, "Phone already exist")
	}
	_, err = p.conn.Exec(context.Background(), "INSERT INTO records (name, last_name, middle_name, address, phone) VALUES ($1, $2, $3, $4, $5)",
		record.Name, record.LastName, record.MiddleName, record.Address, record.Phone)

	if err != nil {
		log.Println(myErr.Wrap(err, "Something going wrong with inserting into db"))
		return myErr.Wrap(err, "Something going wrong with inserting into db")
	}

	return nil
}

// RecordsGet возвращает записи из базы данных на основе предоставленных полей Record.
func (p *Psg) RecordsGet(record dto.Record) ([]dto.Record, error) {
	myErr := pkg.NewMyError("package psg: func (p *Psg) RecordsGet(record dto.Record) ([]dto.Record, *pkg.Errorer)")
	var result []dto.Record
	var values []interface{}
	query := "SELECT * FROM records WHERE TRUE"
	// Будем добавлять параметры в запрос только если они не пустые
	if record.Name != "" {
		query += " AND name=$" + strconv.Itoa(len(values)+1)
		values = append(values, record.Name)
	}
	if record.LastName != "" {
		query += " AND last_name=$" + strconv.Itoa(len(values)+1)
		values = append(values, record.LastName)
	}
	if record.MiddleName != "" {
		query += " AND middle_name=$" + strconv.Itoa(len(values)+1)
		values = append(values, record.MiddleName)
	}
	if record.Address != "" {
		query += " AND address=$" + strconv.Itoa(len(values)+1)
		values = append(values, record.Address)
	}
	if record.Phone != "" {
		query += " AND phone=$" + strconv.Itoa(len(values)+1)
		values = append(values, record.Phone)
	}

	// Выполняем запрос с параметрами
	rows, err := p.conn.Query(context.Background(), query, values...)

	if err != nil {
		log.Println(myErr.Wrap(err, "p.conn.Query(context.Background(), query, values...): Something going wrong with query to database"))
		return nil, myErr.Wrap(err, "p.conn.Query(context.Background(), query, values...): Something going wrong with query to database")
	}
	defer rows.Close()
	for rows.Next() {
		var rec dto.Record
		if err := rows.Scan(&rec.ID, &rec.Name, &rec.LastName, &rec.MiddleName, &rec.Address, &rec.Phone); err != nil {
			log.Println(myErr.Wrap(err, "rows.Scan(&rec.ID, &rec.Name, &rec.LastName, &rec.MiddleName, &rec.Address, &rec.Phone): Something going wrong scan"))
			return nil, myErr.Wrap(err, "rows.Scan(&rec.ID, &rec.Name, &rec.LastName, &rec.MiddleName, &rec.Address, &rec.Phone): Something going wrong scan")
		}
		result = append(result, rec)
	}
	return result, nil
}

// RecordUpdate обновляет существующую запись в базе данных по номеру телефона.
func (p *Psg) RecordUpdate(record dto.Record) error {
	myErr := pkg.NewMyError("package psg: func (p *Psg) RecordUpdate(record dto.Record)")
	/// должны обновляться не все поля, а только те, которые передаём
	var values []interface{}
	query := "UPDATE records SET"
	// Будем добавлять параметры в запрос только если они не пустые
	if record.Name != "" {
		query += " name=$" + strconv.Itoa(len(values)+1)
		values = append(values, record.Name)
	}
	if record.LastName != "" {
		if len(values) > 0 {
			query += ","
		}
		query += " last_name=$" + strconv.Itoa(len(values)+1)
		values = append(values, record.LastName)
	}
	if record.MiddleName != "" {
		if len(values) > 0 {
			query += ","
		}
		query += " middle_name=$" + strconv.Itoa(len(values)+1)
		values = append(values, record.MiddleName)
	}
	if record.Address != "" {
		if len(values) > 0 {
			query += ","
		}
		query += " address=$" + strconv.Itoa(len(values)+1)
		values = append(values, record.Address)
	}
	if record.Phone == "" {
		log.Println(myErr.Wrap(nil, "Empty phone number"))
		return myErr.Wrap(nil, "Empty phone number")
	}
	if len(values) == 0 {
		log.Println(myErr.Wrap(nil, "Nothing to update"))
		return myErr.Wrap(nil, "Nothing to update")
	}
	query += " WHERE phone=$" + strconv.Itoa(len(values)+1)
	values = append(values, record.Phone)
	_, err := p.conn.Exec(context.Background(), query, values...)
	if err != nil {
		log.Println(err, "p.conn.Exec(context.Background(), query, values...): Something going wrong with query to database")
		return myErr.Wrap(err, "p.conn.Exec(context.Background(), query, values...): Something going wrong with query to database")
	}
	return nil
}

// RecordDeleteByPhone удаляет запись из базы данных по номеру телефона.
func (p *Psg) RecordDeleteByPhone(phone string) error {
	myErr := pkg.NewMyError("package psg: func (p *Psg) RecordDeleteByPhone(phone string)")
	query := "DELETE FROM records WHERE phone=$1"
	_, err := p.conn.Exec(context.Background(), query, phone)
	if err != nil {
		log.Println(err, "p.conn.Exec(context.Background(), query, phone): Something going wrong with delete from database")
		return myErr.Wrap(err, "p.conn.Exec(context.Background(), query, phone): Something going wrong with delete from database")
	}
	return nil
}

// Функция закрытия соединения с БД
func (p *Psg) Close() {
	p.conn.Close()
}

// функция для проверки наличия номера телефона в БД
func (p *Psg) CheckPhone(phone string) (bool, error) {
	myErr := pkg.NewMyError("func (p *Psg) CheckPhone(phone string) (bool, *pkg.Errorer)")
	query := "SELECT EXISTS (SELECT * FROM records WHERE phone = $1)"
	var exists bool
	err := p.conn.QueryRow(context.Background(), query, phone).Scan(&exists)
	if err != nil {
		log.Println(myErr.Wrap(err, "p.conn.QueryRow(context.Background(), query, phone).Scan(&exists)"))
		return false, myErr.Wrap(err, "p.conn.QueryRow(context.Background(), query, phone).Scan(&exists)")
	}
	return exists, nil
}
