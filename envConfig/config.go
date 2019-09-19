package envConfig

type Env struct {
	Port          string `envconfig:"PORT" default:"8080"`
	JwtSigningKey string `envconfig:"JWT_SIGNING_KEY" default:"secret"`
	DbConfig
}

type DbConfig struct {
	SqlDialect string `envconfig:"SQL_DIALECT" default:"mysql"`
	SqlUrl     string `envconfig:"SQL_URL" default:"root:password@tcp(localhost:3306)/crm?charset=utf8&parseTime=True&loc=Local"`
}
