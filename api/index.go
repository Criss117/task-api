package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"task.dev/tasks"
	"task.dev/users"
	commonreponse "task.dev/utils/common_reponse"
)

type key string

const userKey key = "user-key"

var usersRepository = users.NewUsersRepository()
var sessionsRepository = users.NewSessionsRepository()
var tasksRepository = tasks.NewTasksRepository()

func Handler() {
	mux := http.NewServeMux()

	taskHandler(mux)
	authHandler(mux)

	server := http.Server{
		Addr:    ":8080",
		Handler: ContentType(mux),
	}

	log.Println("Starting server on port 8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func ContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Println(r.Method, r.URL.Path, time.Since(start))
	})
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie := r.CookiesNamed("session")[0]

		if sessionCookie == nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(commonreponse.Unauthorized("Invalid session 1"))
			return
		}

		log.Println(sessionCookie.Value)
		session := sessionsRepository.FindByToken(sessionCookie.Value)

		if session == nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(commonreponse.Unauthorized("Invalid session 2"))
			return
		}

		user := usersRepository.GetUserByID(session.UserID)

		if user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(commonreponse.Unauthorized("Invalid session 3"))
			return
		}

		userWithSession := users.UserWithSession{
			User:    *user,
			Session: session,
		}

		ctx := context.WithValue(r.Context(), userKey, &userWithSession)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func authHandler(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/sign-up", func(w http.ResponseWriter, r *http.Request) {
		var signUpDto *users.SignUpDto

		if err := json.NewDecoder(r.Body).Decode(&signUpDto); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(commonreponse.BadRequest("Invalid request"))
			return
		}

		newUser := users.NewUser(signUpDto.Name, signUpDto.Email, signUpDto.Password)

		usersRepository.AddUser(newUser)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(commonreponse.Created("User created", newUser))
	})

	mux.HandleFunc("POST /auth/sign-in", func(w http.ResponseWriter, r *http.Request) {
		var signInDto *users.SignInDto

		if err := json.NewDecoder(r.Body).Decode(&signInDto); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(commonreponse.BadRequest("Invalid request"))
			return
		}

		user := usersRepository.GetUserByEmail(signInDto.Email)

		if user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(commonreponse.Unauthorized("Invalid email or password"))
			return
		}

		if user.Password != signInDto.Password {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(commonreponse.Unauthorized("Invalid email or password"))
			return
		}

		session := users.NewSession(user.ID)
		sessionsRepository.Add(session)

		sessionCookie := http.Cookie{
			Name:     "session",
			Value:    session.Token,
			Expires:  time.Now().Add(time.Hour * 24 * 7),
			HttpOnly: true,
			Secure:   true,
		}

		http.SetCookie(w, &sessionCookie)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(commonreponse.Ok("Signed in", session))
	})

	mux.HandleFunc("GET /auth/get-session", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(userKey).(*users.UserWithSession)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(commonreponse.Ok("Session found", user))
	}))

}

func taskHandler(mux *http.ServeMux) {
	mux.HandleFunc("GET /tasks", func(w http.ResponseWriter, r *http.Request) {
		selectQuery := r.URL.Query().Get("select")
		nameQuery := r.URL.Query().Get("name")

		tasks := tasksRepository.GetAllTasks(tasks.Filters{
			Select: selectQuery,
			Name:   nameQuery,
		})

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(commonreponse.Ok("Tasks found", tasks))
	})

	mux.HandleFunc("GET /tasks/{taskId}", func(w http.ResponseWriter, r *http.Request) {
		taskId := r.PathValue("taskId")
		task := tasksRepository.GetTask(taskId)

		if task == nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(commonreponse.NotFound("Task not found"))
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(commonreponse.Ok("Task found", task))
	})

	mux.HandleFunc("POST /tasks", func(w http.ResponseWriter, r *http.Request) {
		var task *tasks.CreateTaskDto

		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(commonreponse.BadRequest("Invalid request"))
			return
		}

		if err := task.Validate(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(commonreponse.InvalidBody("Invalid request", err))
			return
		}

		newTask := tasks.NewTask(task.Name)

		tasksRepository.AddTask(newTask)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(commonreponse.Ok("Task created", newTask))
	})

	mux.HandleFunc("DELETE /tasks/{taskId}", func(w http.ResponseWriter, r *http.Request) {
		taskId := r.PathValue("taskId")
		task := tasksRepository.GetTask(taskId)

		if task == nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(commonreponse.NotFound("Task not found"))
			return
		}

		tasksRepository.DeleteTask(taskId)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(commonreponse.Deleted("Task deleted"))
	})

	mux.HandleFunc("PATCH /tasks/{taskId}", func(w http.ResponseWriter, r *http.Request) {
		var taskToUpdate *tasks.UpdateTaskNameDto

		if err := json.NewDecoder(r.Body).Decode(&taskToUpdate); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(commonreponse.BadRequest("Invalid request"))
			return
		}

		if err := taskToUpdate.Validate(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(commonreponse.BadRequest(err.Error()))
			return
		}

		taskId := r.PathValue("taskId")
		task := tasksRepository.GetTask(taskId)

		if task == nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(commonreponse.NotFound("Task not found"))
			return
		}

		task.UpdateTaskName(taskToUpdate.Name)

		tasksRepository.UpdateTask(task)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(commonreponse.Ok("Task updated", task))
	})

	mux.HandleFunc("PATCH /tasks/{taskId}/completed", func(w http.ResponseWriter, r *http.Request) {
		taskId := r.PathValue("taskId")
		task := tasksRepository.GetTask(taskId)

		if task == nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(commonreponse.NotFound("Task not found"))
			return
		}

		task.ToogleTaskCompleted()

		tasksRepository.UpdateTask(task)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(commonreponse.Ok("Task updated", task))
	})
}
