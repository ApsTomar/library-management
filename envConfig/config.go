package envConfig

type Env struct {
	UserSvcPort       string `envconfig:"PORT" default:"8000"`
	BookSvcPort       string `envconfig:"PORT" default:"8001"`
	ManagementSvcPort string `envconfig:"PORT" default:"8002"`
	JwtSigningKey     string `envconfig:"JWT_SIGNING_KEY" default:"secret"`
	DbConfig
	FluentConfig
}

type DbConfig struct {
	SqlDialect string `envconfig:"SQL_DIALECT" default:"mysql"`
	SqlUrl     string `envconfig:"SQL_URL" default:"root:password@tcp(127.0.0.1:3306)/library?charset=utf8&parseTime=True&loc=Local"`
}

type FluentConfig struct {
	FluentPort string `envconfig:"FLUENT_PORT" default:"24224"`
	FluentHost string `envconfig:"FLUENT_HOST" default:"127.0.0.1"`
}
