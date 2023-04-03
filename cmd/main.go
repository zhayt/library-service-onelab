package main

import "github.com/zhayt/user-storage-service/internal/app"

func main() {
	if err := app.Run(); err != nil {
		return
	}
}
