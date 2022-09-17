package config

// configuration properties based on env variables
type Properties struct {
	Port               string `env:"MY_APP_PORT" env-default:"8080"`
	Host               string `env:"HOST" env-default:"localhost"`
	DBHost             string `env:"DB_HOST" env-default:"localhost"`
	DPPort             string `env:"DB_PORT" env-default:"27017"`
	DBName             string `env:"DB_NAME" env-default:"GoWebAPI"`
	ProductsCollection string `env:"PRODUCTS_COL_NAME" env-default:"products"`
	UsersCollection    string `env:"USERS_COL_NAME" env-default:"users"`
	JwtTokenSecret     string `env:"JWT_TOKEN_SECRET" env-default:"replacethis"`
}

var LoggerConfigFormat = `{"time":"${time_rfc3339_nano}","${header:X-Correlation-ID}","id":"${id}","remote_ip":"${remote_ip}",` +
	`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
	`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
	`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n"
