package main

func main() {
	app := App{}
	app.Initialise(DB_USER, DB_USER_PASSWORD, DB_HOST, DB_NAME)
	app.Run("localhost:3333")
}
