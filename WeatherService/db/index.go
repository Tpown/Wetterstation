package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Persistence struct {
	db *sql.DB
}

func NewPersistence() *Persistence {
	return &Persistence{}
}

// TODO: Get the environment externally
func (this *Persistence) Connect() error {
	db, err := sql.Open("postgres", "postgres://root:password@db/weather_service?sslmode=disable")
	if nil != err {
		return err
	}
	this.db = db
	return nil
}

func (this *Persistence) CreateTables() {
	r, err := this.db.Exec(`
    CREATE TABLE IF NOT EXISTS report(
    id Serial PRIMARY KEY,
    temperature NUMERIC(5) NOT NULL,
    humidity INTEGER NOT NULL,
    wind INTEGER NOT NULL,
    rain NUMERIC(5) NOT NULL,
    dateTime CHAR(70) NOT NULL
  )`)

	if nil != err {
		log.Fatal(err.Error())
		return
	}

	_, err = r.RowsAffected()
	if nil != err {
		log.Fatal(err.Error())
		return
	}
}

func (this *Persistence) SaveReport(temperature float64, humidity int16, wind int16, rain float64, dateTime string) bool {
	result, err := this.db.Exec(`
  INSERT INTO report(
    temperature, humidity, wind, rain, dateTime
  ) VALUES ($1, $2, $3, $4, $5)`, temperature, humidity, wind, rain, dateTime)
	if nil != err {
		fmt.Println("Bin here")
		log.Fatal(err.Error())
		return false
	}

	nbr, err := result.RowsAffected()
	if nil != err {
		fmt.Println("Bin here 2")
		log.Fatal(err.Error())
		return false
	}
	return nbr == 1
}

// the structure matches the db columns
type Report struct {
	Temperature  float64
	Humidity     int16
	WindStrength int16
	Rainfall     float64
	DateTime     string
}

func (this *Persistence) GetLastReport() *Report {
	row := this.db.QueryRow(`SELECT temperature, humidity, wind, rain, dateTime FROM report ORDER BY datetime DESC LIMIT 1`)
	result := &Report{}
	row.Scan(&result.Temperature, &result.Humidity, &result.WindStrength, &result.Rainfall, &result.DateTime)
	return result
}
