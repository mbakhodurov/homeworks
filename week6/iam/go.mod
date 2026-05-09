module github.com/mbakhodurov/homeworks/week6/iam

go 1.25.4

replace github.com/mbakhodurov/homeworks/week6/platform => ../platform

replace github.com/mbakhodurov/homeworks/week6/shared => ../shared

require (
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/caarlos0/env/v11 v11.4.1
	github.com/joho/godotenv v1.5.1
)

require (
	github.com/gomodule/redigo v1.9.3 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
)
