package todos

import (
	Sessions "TodoApp/Database/Sessions"
	Todos "TodoApp/Database/Todos"
	TemplateParser "TodoApp/RoutesHandler/TemplateParser"
	"net/http"
	"text/template"
)

func HandleRoutes() {

	http.HandleFunc("GET /todo", func(responseWriter http.ResponseWriter, request *http.Request) {
		cookie, err := request.Cookie("session_id")
		if err != nil {
			http.Redirect(responseWriter, request, "/login", http.StatusSeeOther)
			return
		}

		sessionId := cookie.Value

		todos := Todos.GetTodos(sessionId)

		if Sessions.DoesSessionExistsInDatabase(sessionId) {
			TemplateParser.ParseTemplateWithAdditionalData("Todo/TodoPage", "Todo", responseWriter, request, todos)
		} else {
			http.Redirect(responseWriter, request, "/login", http.StatusSeeOther)
		}

	})

	http.HandleFunc("POST /addTodo", func(responseWriter http.ResponseWriter, request *http.Request) {

		err := request.ParseForm()
		if err != nil {
			http.Error(responseWriter, "Unable to parse form", http.StatusBadRequest)
			return
		}

		formData := request.Form

		todo := formData.Get("Todo")

		// Process formData as needed
		cookie, err := request.Cookie("session_id")
		if err != nil {
			http.Redirect(responseWriter, request, "/login", http.StatusSeeOther)
			return
		}

		sessionId := cookie.Value

		if Sessions.DoesSessionExistsInDatabase(sessionId) {
			Todos.AddTodo(todo, sessionId)
			tmpl, err := template.New("todo").Parse("<li>{{.}}</li>")
			if err != nil {
				http.Error(responseWriter, "Unable to parse template", http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(responseWriter, todo)
			if err != nil {
				http.Error(responseWriter, "Unable to execute template", http.StatusInternalServerError)
				return
			}
		} else {
			http.Redirect(responseWriter, request, "/login", http.StatusSeeOther)
		}

	})
}
