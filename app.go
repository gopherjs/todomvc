package main

import (
	jQ "github.com/rusco/jquery"
	"github.com/rusco/todomvc/utils"
)

const (
	KEY        = "TodoMVC-GopherJS"
	ENTER_KEY  = 13
	ESCAPE_KEY = 27
)

func main() {
	utils.RegisterHandlebarsHelper()
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

	a.newTodoJq.On(jQ.KEYUP, a.create)
	a.toggleAllJq.On(jQ.CHANGE, a.toggleAll)
	a.footerJq.On(jQ.CLICK, "#clear-completed", a.destroyCompleted)
	a.todoListJq.On(jQ.CHANGE, ".toggle", a.toggle)
	a.todoListJq.On(jQ.DBLCLICK, "label", a.edit)
	a.todoListJq.On(jQ.KEYUP, ".edit", a.blurOnEnter)
	a.todoListJq.On(jQ.FOCUSOUT, ".edit", a.update)
	a.todoListJq.On(jQ.CLICK, ".destroy", a.destroy)
}
func (a *App) initRouter() {
	router := utils.NewRouter()
	router.On("/:filter", func(filter string) {
		a.filter = filter
		a.render()
	})
	router.Init("/all")
}
func (a *App) render() {
	todos := a.getFilteredTodos()
	strtodoHb := a.todoHb.Invoke(todos).String()
	a.todoListJq.SetHtml(strtodoHb)
	a.mainJq.Toggle(len(a.todos) > 0)
	a.toggleAllJq.SetProp("checked", len(a.getActiveTodos()) != 0)
	a.renderfooter()
	a.newTodoJq.Focus()
	utils.Store(KEY, a.todos)
}
func (a *App) renderfooter() {
	activeTodoCount := len(a.getActiveTodos())
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
	checked := !a.toggleAllJq.Prop("checked").(bool)
	for idx := range a.todos {
		a.todos[idx].Completed = checked
	}
	a.render()
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
	switch a.filter {
	case "active":
		return a.getActiveTodos()
	case "completed":
		return a.getCompletedTodos()
	default:
		return a.todos
	}
}
func (a *App) destroyCompleted(e jQ.Event) {
	a.todos = a.getActiveTodos()
	a.filter = "all"
	a.render()
}
func (a *App) indexFromEl(e jQ.Event) int {
	id := jQ.NewJQuery(e.Target).Closest("li").Data("id")
	for idx, val := range a.todos {
		if val.Id == id {
			return idx
		}
	}
	return -1
}
func (a *App) create(e jQ.Event) {
	val := jQ.Trim(a.newTodoJq.Val())
	if len(val) == 0 || e.KeyCode != ENTER_KEY {
		return
	}
	newToDo := ToDo{Id: utils.Uuid(), Text: val, Completed: false}
	a.todos = append(a.todos, newToDo)
	a.newTodoJq.SetVal("")
	a.render()
}
func (a *App) toggle(e jQ.Event) {
	idx := a.indexFromEl(e)
	a.todos[idx].Completed = !a.todos[idx].Completed
	a.render()
}
func (a *App) edit(e jQ.Event) {
	input := jQ.NewJQuery(e.Target).Closest("li").AddClass("editing").Find(".edit")
	input.SetVal(input.Val()).Focus()
}
func (a *App) blurOnEnter(e jQ.Event) {
	switch e.KeyCode {
	case ENTER_KEY:
		jQ.NewJQuery(e.Target).Blur()
	case ESCAPE_KEY:
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
	idx := a.indexFromEl(e)
	if len(val) > 0 {
		a.todos[idx].Text = val
	} else {
		a.todos = append(a.todos[:idx], a.todos[idx+1:]...)
	}
	a.render()
}
func (a *App) destroy(e jQ.Event) {
	idx := a.indexFromEl(e)
	a.todos = append(a.todos[:idx], a.todos[idx+1:]...)
	a.render()
}
