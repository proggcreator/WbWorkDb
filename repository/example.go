package repository

import (
	"encoding/json"
	"fmt"
	"time"
)

func task1_init() Message {
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
	return mess
}
