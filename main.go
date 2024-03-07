package main

func main() {
	server := NewAPIServer(":7777")

	server.Run()
}
