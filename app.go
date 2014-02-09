package main

import (
	jQ "github.com/rusco/jquery"
	"github.com/rusco/todomvc/utils"

	"github.com/neelance/gopherjs/js" //2do: refactor: remove dependency on this package
)

const (
	KEY       = "TodoMVC-GopherJS"
	ENTER_KEY = 13
)

func main() {
	app := NewApp()
	app.bindEvents()
	app.render()
}

type ToDo struct {
	Id        string
	Text      string
	Completed bool
}

type App struct {
	todos       []ToDo
	todoHb      *utils.Handlebar
	footerHb    *utils.Handlebar
	todoAppJq   jQ.JQuery
	headerJq    jQ.JQuery
	mainJq      jQ.JQuery
	footerJq    jQ.JQuery
	newTodoJq   jQ.JQuery
	toggleAllJq jQ.JQuery
	todoListJq  jQ.JQuery
	countJq     jQ.JQuery
	clearBtnJq  jQ.JQuery
}

//constructor
func NewApp() *App {

	somethingToDo := make([]ToDo, 0)

	utils.Retrieve(KEY, &somethingToDo)
	todoTemplate := jQ.NewJQuery("#todo-template").Html()
	todoHb := utils.CompileHandlebar(todoTemplate)
	footerTemplate := jQ.NewJQuery("#footer-template").Html()
	footerHb := utils.CompileHandlebar(footerTemplate)

	todoAppJq := jQ.NewJQuery("#todoapp")
	headerJq := todoAppJq.Find("#header")
	mainJq := todoAppJq.Find("#main")
	footerJq := todoAppJq.Find("#footer")
	newTodoJq := headerJq.Find("#new-todo")
	toggleAllJq := mainJq.Find("#toggle-all")
	todoListJq := mainJq.Find("#todo-list")
	countJq := footerJq.Find("#todo-count")
	clearBtnJq := footerJq.Find("#clear-completed")
	return &App{somethingToDo, todoHb, footerHb, todoAppJq, headerJq, mainJq, footerJq, newTodoJq, toggleAllJq, todoListJq, countJq, clearBtnJq}
}

func (a *App) bindEvents() {

	a.newTodoJq.On(jQ.EvtKEYUP, a.create)
	a.toggleAllJq.On(jQ.EvtCHANGE, a.toggleAll)
	a.footerJq.OnSelector(jQ.EvtCLICK, "#clear-completed", a.destroyCompleted)
	a.todoListJq.OnSelector(jQ.EvtCHANGE, ".toggle", a.toggle)
	a.todoListJq.OnSelector(jQ.EvtDBLCLICK, "label", a.edit)
	a.todoListJq.OnSelector(jQ.EvtKEYPRESS, ".edit", a.blurOnEnter)
	a.todoListJq.OnSelector(jQ.EvtBLUR, ".edit", a.update)
	a.todoListJq.OnSelector(jQ.EvtCLICK, ".destroy", a.destroy)
}

func (a *App) render() {

	strtodoHb := a.todoHb.Invoke(a.todos).String()
	a.todoListJq.SetHtml(strtodoHb)
	a.mainJq.Toggle(len(a.todos) > 0)
	a.toggleAllJq.SetProp("checked", a.activeTodoCount() != 0)
	a.renderfooter()
	utils.Store(KEY, a.todos)
}

func (a *App) renderfooter() {

	activeTodoCount := a.activeTodoCount()
	activeTodoWord := utils.Pluralize(activeTodoCount, "item")
	completedTodos := len(a.todos) - activeTodoCount
	footerData := struct {
		ActiveTodoCount int
		ActiveTodoWord  string
		CompletedTodos  int
	}{
		activeTodoCount, activeTodoWord, completedTodos,
	}
	footerJqStr := a.footerHb.Invoke(footerData).String()
	a.footerJq.Toggle(len(a.todos) > 0).SetHtml(footerJqStr)
}
func (a *App) toggleAll(this js.Object, e *jQ.Event) {

	checked := !a.toggleAllJq.Prop("checked")
	for idx := range a.todos {
		a.todos[idx].Completed = checked
	}
	a.render()
}
func (a *App) activeTodoCount() int {

	count := 0
	for _, val := range a.todos {
		if !val.Completed {
			count += 1
		}
	}
	return count
}
func (a *App) destroyCompleted(this js.Object, e *jQ.Event) {

	todosTmp := make([]ToDo, 0)
	for _, val := range a.todos {
		if !val.Completed {
			todosTmp = append(todosTmp, val)
		}
	}
	a.todos = make([]ToDo, len(todosTmp))
	copy(a.todos, todosTmp)
	a.render()
}

func (a *App) create(this js.Object, e *jQ.Event) {

	val := jQ.Trim(a.newTodoJq.Val())

	print(e.KeyCode, val)

	if val == "" || e.KeyCode != ENTER_KEY {
		return
	}
	newToDo := ToDo{Id: utils.Uuid(), Text: val, Completed: false}
	a.todos = append(a.todos, newToDo)
	a.newTodoJq.SetVal("")
	a.render()
}

func (a *App) toggle(this js.Object, e *jQ.Event) {

	id := jQ.NewJQueryByObject(this).Closest("li").Data("id")
	for idx, val := range a.todos {
		if val.Id == id {
			a.todos[idx].Completed = !a.todos[idx].Completed
		}
	}
	a.render()

}

func (a *App) edit(this js.Object, e *jQ.Event) {
	thisJq := jQ.NewJQueryByObject(this)
	input := thisJq.Closest("li").AddClass("editing").Find(".edit")
	val := input.Val()
	input.SetVal(val).Focus()

}

func (a *App) blurOnEnter(this js.Object, e *jQ.Event) {

	if e.KeyCode == ENTER_KEY {
		jQ.NewJQueryByObject(this).Blur()
	}
}

func (a *App) update(this js.Object, e *jQ.Event) {
	thisJq := jQ.NewJQueryByObject(this)
	val := jQ.Trim(thisJq.Val())

	id := thisJq.Closest("li").RemoveClass("editing").Data("id")
	for idx := range a.todos {
		if a.todos[idx].Id == id {
			if len(val) > 0 {
				a.todos[idx].Text = val
			} else {
				a.todos[idx].Id = "delete"
			}
		}
	}
	todosTmp := make([]ToDo, 0)
	for _, val := range a.todos {
		if val.Id != "delete" {
			todosTmp = append(todosTmp, val)
		}
	}
	a.todos = make([]ToDo, len(todosTmp))
	copy(a.todos, todosTmp)
	a.render()

}

func (a *App) destroy(this js.Object, e *jQ.Event) {

	id := jQ.NewJQueryByObject(this).Closest("li").Data("id")

	todosTmp := make([]ToDo, 0)
	for _, val := range a.todos {
		if val.Id != id {
			todosTmp = append(todosTmp, val)
		}
	}
	a.todos = make([]ToDo, len(todosTmp))
	copy(a.todos, todosTmp)
	a.render()

}
