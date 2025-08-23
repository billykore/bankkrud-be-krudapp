# TapMoney Backend Service

This is a sample backend service for a fictional electronic money application called TapMoney.
It demonstrates the implementation of a RESTful API using Go, 
following best practices and design patterns.

# Hexagonal Architecture

https://netflixtechblog.com/ready-for-changes-with-hexagonal-architecture-b315ec967749

# Packages

Some of the open-source packages we used are:

- [Echo](https://echo.labstack.com) for http routing.
- [GORM](https://gorm.io) for database ORM.
- [Validator](https://github.com/go-playground/validator) for request validation.
- [zap](https://github.com/uber-go/zap) for logging.
- [viper](https://github.com/spf13/viper) for managing project configurations.
- [JWT](https://github.com/golang-jwt/jwt) for generate and validate authorization token.
- [swag](https://github.com/swaggo/swag) and [echo-swagger](https://github.com/swaggo/echo-swagger) for generate API documentation.
- [ecszap](https://github.com/elastic/ecs-logging-go-zap) to support ECS for zap logger.
- [Google Wire](https://github.com/google/wire) for dependency injection.
- [testify](https://github.com/stretchr/testify) for unit testing.
