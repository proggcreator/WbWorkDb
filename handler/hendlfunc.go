package handler

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
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

func (h *Handler) request_error_code(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatalf("invalid id param")
		return
	}
	h.repository.Req_Error_Code(id)

}
func (h *Handler) do_something(c *gin.Context) {

	h.repository.Do_something()

}
