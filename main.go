package main

func main() {
	LoggerServerINFO("@@@Hello world!!!@@@")
	LoggerServerINFO("Listening...")

	err := Server.ListenAndServe()
	if err != nil {
		LoggerServerERROR(err.Error())
	}

	LoggerServerINFO("@@@End of the world@@@")
}
