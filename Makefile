server:
	nodemon --watch "./**/*.go" --signal SIGTERM --exec "set APP_ENV=dev&&go run main.go"
