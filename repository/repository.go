package repository

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

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
type MyapiTask interface {
	Req_Error_Code(id int) (str []string)
	Do_something() (err error)
}

type Repository struct {
	MyapiTask
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		MyapiTask: NewMyapiPostgres(db),
	}
}
func DefTypeError(i int) (coderet string) {

	if (i%50000 == 0) || (i > 50000) {
		coderet = "400"
	} else {
		coderet = "500"

	}
	return coderet

}
func (r *MyapiPostgres) Do_something() (err error) {
	m := task1_init() //инициализация структуры запроса
	rows, err := r.db.Query("SELECT * FROM test.do_something($1,$2,$3)", m._employee_id, m._test_data, m._date)
	if err != nil {
		fmt.Println(err, "PostDoSmt error")
		return err
	}
	defer rows.Close()

	byteStr := make([]byte, 0)
	for rows.Next() {
		err = rows.Scan(&byteStr)
		if err != nil {
			fmt.Println(err, "scanning error")
		}
	}

	return nil
}

//возвращает код ошибки и сообщение если ошибка пользовательская
func (r *MyapiPostgres) Req_Error_Code(id int) (strlist []string) {

	_, err := r.db.Query("SELECT * FROM  test.get_db_error($1) ", id) //вызов функции get_db_error
	if err, ok := err.(*pq.Error); ok {

		//код ошибки в int для DefTypeEror
		i, cerr := strconv.Atoi(string(err.Code))
		if cerr != nil {
			os.Exit(2)
		}
		strlist = append(strlist, string(i))

		//если ошибка пользовательская то вывод сообщениея из бд
		if DefTypeError(i) == "500" {
			strlist = append(strlist, err.Message)
		}

	}
	return strlist

}
