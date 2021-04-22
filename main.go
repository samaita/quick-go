package main

func main() {
	initFirebase()
	initDB("sqlite3", "db_quick_go.db")
	initRedis()
	initHandler()
}
