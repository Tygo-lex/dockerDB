package templates

const (
	MySQLDockerfile = `
FROM mysql:latest

ENV MYSQL_ROOT_PASSWORD=root
ENV MYSQL_DATABASE=mydb
ENV MYSQL_USER=user
ENV MYSQL_PASSWORD=password

EXPOSE 3306
`

	PostgresDockerfile = `
FROM postgres:latest

ENV POSTGRES_DB=mydb
ENV POSTGRES_USER=user
ENV POSTGRES_PASSWORD=password

EXPOSE 5432
`

	MongoDBDockerfile = `
FROM mongo:latest

EXPOSE 27017
`

	RedisDockerfile = `
FROM redis:latest

EXPOSE 6379
`
)