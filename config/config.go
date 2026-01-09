package config

// This is where we can create all of the data structs for our application
// This will make typing easier
type API struct {
	Address string
	Host    string
	Port    string
}

func SetUpAPI() API {

	// TODO: We eventually want to read this from an env
	// LoadEnvironmentVariables(env)
	// api := API{
	// 	Host:    os.Getenv("API_HOST"),
	// 	Port:    os.Getenv("API_PORT"),
	// 	Address: os.Getenv("API_ADDRESS"),
	// }

	api := API{
		Host: "localhost",
		Port: "3000",
	}
	return api
}
