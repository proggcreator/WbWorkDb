package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	u "example.com/dbwork/config"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}
type Message struct {
	_employee_id int
	_test_data   []byte
	_date        time.Time
}

type Myjson struct {
	name string
	body string
	age  int
}

func init() {
	if err := godotenv.Load(); err != nil { //переменные из .evm
		fmt.Println("No .env file found")
	}
}

//подключение к бд
func opendb(cfg Config) (*sqlx.DB, error) {

	db, err := sqlx.Open("postgres", fmt.Sprintf("user=%s host=%s port=%s dbname=%s password=%s sslmode=%s",
		cfg.Username, cfg.Host, cfg.Port, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
func PostDoSmt(db *sqlx.DB, m Message) (jsonString []byte, err error) {
	rows, err := db.Query("SELECT * FROM test.do_something($1,$2,$3)", m._employee_id, m._test_data, m._date)
	if err != nil {
		fmt.Println(err, "PostDoSmt error")
		return nil, err
	}
	defer rows.Close()

	byteStr := make([]byte, 0)
	for rows.Next() {
		err = rows.Scan(&byteStr)
		if err != nil {
			fmt.Println(err, "scanning error")
		}
	}
	return byteStr, nil
}

//определение кода ошибки
func DefTypeError(i int) (coderet string) {
	//if (i%50 == 0) || (i > 50000) {
	if i == 50000 {
		coderet = "400"
	} else {
		coderet = "500"

	}
	return coderet

}
func MyGetError(db *sqlx.DB, id int) (codeerr string, err error) {
	_, err = db.Query("SELECT * FROM  test.get_db_error($1) ", id) //вызов функции get_db_error
	if err, ok := err.(*pq.Error); ok {
		i, err := strconv.Atoi(string(err.Code)) //код ошибки в int
		if err != nil {
			os.Exit(2)
		}
		codeerr = DefTypeError(i) //код ошибки

	}
	return codeerr, err
}

func initConfig() error { //переменные из config
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func task1_example(db *sqlx.DB) {
	// заполняем данными запрос
	m := Myjson{"Alice", "Hello", 1}
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err, "marshal error")
	}
	mess := Message{
		_date:        time.Now(),
		_employee_id: 6,
		_test_data:   b,
	}
	// запрос в бд
	x, err := PostDoSmt(db, mess)

	fmt.Println(x) //возвращаемое значение
}
func task2_example(db *sqlx.DB) {
	id := 2
	errcode, err := MyGetError(db, id)
	if err != nil {
		fmt.Println(err)     //вывод типа ошибки
		fmt.Println(errcode) //вывод кода ошибки

	}
}

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter)) //формат для логгера json
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs")
	}
	conf := u.New()
	db, err := opendb(Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: conf.Username,
		Password: conf.Password,
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode")})
	if err != nil {
		logrus.Fatalf("faled to initialization %s", err.Error())
	}
	defer db.Close()

	//task1
	task1_example(db)
	//task2
	task2_example(db)

}
