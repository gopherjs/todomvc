package utils_test

import "testing"
import "utils"

const KEY = "utils_test_key"

type ToDoType []struct {
	Id        int
	Text      string
	Completed bool
}

func TestStruct(t *testing.T) {
	ToDo := ToDoType{
		{1, "Not Yet", false},
		{2, "This one is done", true},
	}

	utils.Store(KEY, &ToDo)

	moreToDo := make(ToDoType, 0)
	utils.Retrieve(KEY, &moreToDo)

	countId := 0
	strComplete := ""
	countCompleted := 0

	for k, v := range moreToDo {
		print(k, v.Id, v.Text, v.Completed)
		countId += v.Id
		strComplete += v.Text
		if v.Completed {
			countCompleted += 1
		}
	}

	if countId != 3 {
		t.Fail()
	}
	if strComplete != "Not YetThis one is done" {
		t.Fail()
	}
	if countCompleted != 1 {
		t.Fail()
	}
}
