# cmd/shortener

В данной директории содержится код, который скомпилируется в бинарное приложение.

Рекомендуется помещать только код, необходимый для запуска приложения, но не бизнес-логику.

Название директории должно соответствовать названию приложения.

Директория `cmd/shortener` содержит:
- точку входа в приложение (функция `main`)
- инициализацию зависимостей (можно вынести в отдельный пакет `internal/app`)
- настройку и запуск HTTP-сервера (можно вынести в отдельный пакет `internal/router`)
- обработку сигналов завершения работы приложения

Запуск сервера с БД
`go run . -d "postgres://postgres:postgres@localhost:5432/urlshortener?sslmode=disable"`
Запуск с файловым хранилищем
`go run . -f "/var/www/golang/yapracticum/go-advanced/urlshortener/storage/links.json"`

Для дебага `DATABASE_DSN="postgres://postgres:postgres@localhost:5432/urlshortener?sslmode=disable"`