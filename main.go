// main.go
package main

func main() {
	ws := WebServivce{}
	ws.Initialize("root", "Cameron31!", "MalwareURLs", "allowAllFiles=true")
	ws.Run(":8080")
}
