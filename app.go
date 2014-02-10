package main

import (
	jQ "github.com/rusco/jquery"
	"github.com/rusco/todomvc/utils"

	"github.com/neelance/gopherjs/js" //2do: refactor: remove dependency on this package
)

const (
	KEY        = "TodoMVC-GopherJS"
	ENTER_KEY  = 13
	ESCAPE_KEY = 27
)

func main() {

	utils.RegisterHelper()

	app := NewApp()
	app.bindEvents()
	app.initRouter()
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
	filter      string
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
	filter := "all"
	return &App{somethingToDo, todoHb, footerHb, todoAppJq, headerJq, mainJq, footerJq, newTodoJq, toggleAllJq, todoListJq, countJq, clearBtnJq, filter}
}

func (a *App) bindEvents() {

	a.newTodoJq.On(jQ.EvtKEYUP, a.create)
	a.toggleAllJq.On(jQ.EvtCHANGE, a.toggleAll)
	a.footerJq.OnSelector(jQ.EvtCLICK, "#clear-completed", a.destroyCompleted)
	a.todoListJq.OnSelector(jQ.EvtCHANGE, ".toggle", a.toggle)
	a.todoListJq.OnSelector(jQ.EvtDBLCLICK, "label", a.edit)
	a.todoListJq.OnSelector(jQ.EvtKEYUP, ".edit", a.blurOnEnter)
	a.todoListJq.OnSelector(jQ.EvtFOCUSOUT, ".edit", a.update)
	a.todoListJq.OnSelector(jQ.EvtCLICK, ".destroy", a.destroy)
}

func (a *App) initRouter() {

	router := js.Global("Router").New()
	router.Call("on", "/:filter", func(filter string) {
		a.filter = filter
		a.render()
	})
	router.Call("init", "/all")
}

func (a *App) render() {

	todos := a.getFilteredTodos()

	strtodoHb := a.todoHb.Invoke(todos).String()
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
	filter := a.filter

	footerData := struct {
		ActiveTodoCount int
		ActiveTodoWord  string
		CompletedTodos  int
		Filter          string
	}{
		activeTodoCount, activeTodoWord, completedTodos, filter,
	}
	footerJqStr := a.footerHb.Invoke(footerData).String()
	a.footerJq.Toggle(len(a.todos) > 0).SetHtml(footerJqStr)
}
func (a *App) toggleAll(e jQ.Event) {

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

func (a *App) getActiveTodos() []ToDo {

	todosTmp := make([]ToDo, 0)
	for _, val := range a.todos {
		if !val.Completed {
			todosTmp = append(todosTmp, val)
		}
	}
	return todosTmp
}

func (a *App) getCompletedTodos() []ToDo {

	todosTmp := make([]ToDo, 0)
	for _, val := range a.todos {
		if val.Completed {
			todosTmp = append(todosTmp, val)
		}
	}
	return todosTmp
}

func (a *App) getFilteredTodos() []ToDo {

	if a.filter == "active" {
		return a.getActiveTodos()
	}

	if a.filter == "completed" {
		return a.getCompletedTodos()
	}

	return a.todos
}

func (a *App) destroyCompleted(e jQ.Event) {

	todosTmp := make([]ToDo, 0)
	for _, val := range a.todos {
		if !val.Completed {
			todosTmp = append(todosTmp, val)
		}
	}
	a.todos = make([]ToDo, len(todosTmp))
	copy(a.todos, todosTmp)
	a.filter = "all"
	a.render()
}

func (a *App) create(e jQ.Event) {

	val := jQ.Trim(a.newTodoJq.Val())
	if val == "" || e.KeyCode != ENTER_KEY {
		return
	}
	newToDo := ToDo{Id: utils.Uuid(), Text: val, Completed: false}
	a.todos = append(a.todos, newToDo)
	a.newTodoJq.SetVal("")
	a.render()
}

func (a *App) toggle(e jQ.Event) {

	id := jQ.NewJQuery(e.Target).Closest("li").Data("id")
	for idx, val := range a.todos {
		if val.Id == id {
			a.todos[idx].Completed = !a.todos[idx].Completed
		}
	}
	a.render()
}

func (a *App) edit(e jQ.Event) {
	thisJq := jQ.NewJQuery(e.Target)
	input := thisJq.Closest("li").AddClass("editing").Find(".edit")
	val := input.Val()
	input.SetVal(val).Focus()

}

func (a *App) blurOnEnter(e jQ.Event) {

	if e.KeyCode == ENTER_KEY {
		jQ.NewJQuery(e.Target).Blur()
	}

	if e.KeyCode == ESCAPE_KEY {
		jQ.NewJQuery(e.Target).SetData("abort", "true").Blur()
	}
}

func (a *App) update(e jQ.Event) {
	thisJq := jQ.NewJQuery(e.Target)
	val := jQ.Trim(thisJq.Val())

	if thisJq.Data("abort") == "true" {
		thisJq.SetData("abort", "false")
		a.render()
		return
	}

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

func (a *App) destroy(e jQ.Event) {

	id := jQ.NewJQuery(e.Target).Closest("li").Data("id")

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
