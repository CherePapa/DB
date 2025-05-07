package server

import (
	"DB-project/cmd/database"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"go.etcd.io/bbolt"
)

func StartServer(db *bbolt.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Получение данных из базы данных
		medicines, err := database.GetMedicines(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Отображение данных с использованием шаблона HTML
		tmpl := template.Must(template.New("medicines").Parse(`
            <!DOCTYPE html>
            <html>
            <head>
                <title>Pharmacy Database</title>
            </head>
            <body>
                <h1>Medicines</h1>
                <ul>
                    {{range .}}
                        <li>{{.Name}} - ${{.Price}} <a href="/delete/{{.ID}}">Delete</a></li>
                    {{end}}
                </ul>
                <form action="/add" method="post">
                    <input type="text" name="name" placeholder="Name" required>
                    <input type="number" step="0.01" name="price" placeholder="Price" required>
                    <button type="submit">Add Medicine</button>
                </form>
            </body>
            </html>
        `))

		tmpl.Execute(w, medicines)
	})

	http.HandleFunc("/delete/", func(w http.ResponseWriter, r *http.Request) {
		// Получение ID из URL
		idStr := r.URL.Path[len("/delete/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		// Удаление записи из базы данных
		err = database.DeleteMedicine(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Получение данных из формы
		name := r.FormValue("name")
		price, err := strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil {
			http.Error(w, "Invalid price", http.StatusBadRequest)
			return
		}

		// Получение всех лекарств для определения следующего ID
		medicines, err := database.GetMedicines(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Вставка данных в базу данных
		nextID := len(medicines) + 1
		err = database.InsertMedicine(db, nextID, name, price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// Запуск веб-сервера
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
