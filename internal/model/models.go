package model

// Request представляет запрос на сокращение URL
// @Schema(
//
//	required={"url"},
//	example={"url": "https://example.com/very-long-url"}
//
// )
type Request struct {
	// Оригинальный URL для сокращения
	// Пример: "https://example.com/very/long/url/to/be/shortened"
	URL string `json:"url"`
}

// Response содержит результат сокращения URL
// @Schema(
//
//	example={"result": "http://short.ly/abc123"}
//
// )
type Response struct {
	// Сокращенный URL
	// Пример: "http://short.ly/abc123"
	Result string `json:"result"`
}

// Link представляет полную информацию о сокращенной ссылке
// @Schema(
//
//	example={
//	    "uuid": "550e8400-e29b-41d4-a716-446655440000",
//	    "short_url": "http://short.ly/abc123",
//	    "original_url": "https://example.com/very-long-url"
//	}
//
// )
type Link struct {
	// Уникальный идентификатор ссылки (UUID)
	UUID string `json:"uuid"`

	// Сокращенный URL
	ShortURL string `json:"short_url"`

	// Оригинальный URL
	OriginalURL string `json:"original_url"`
}

// LinkPair представляет пару сокращенного и оригинального URL
// @Schema(
//
//	example={
//	    "short_url": "http://short.ly/abc123",
//	    "original_url": "https://example.com/very-long-url"
//	}
//
// )
type LinkPair struct {
	// Сокращенный URL
	ShortURL string `json:"short_url"`

	// Оригинальный URL
	OriginalURL string `json:"original_url"`
}

// BatchCreateRequest представляет запрос на пакетное создание сокращенных URL
// @Schema(
//
//	example={
//	    "correlation_id": "request-123",
//	    "original_url": "https://example.com/very-long-url"
//	}
//
// )
type BatchCreateRequest struct {
	// Идентификатор для сопоставления с ответом
	CorrelationID string `json:"correlation_id"`

	// Оригинальный URL для сокращения
	OriginalURL string `json:"original_url"`
}

// BatchCreateResponse содержит результат пакетного создания сокращенных URL
// @Schema(
//
//	example={
//	    "correlation_id": "request-123",
//	    "short_url": "http://short.ly/abc123"
//	}
//
// )
type BatchCreateResponse struct {
	// Идентификатор из запроса
	CorrelationID string `json:"correlation_id"`

	// Сокращенный URL
	ShortURL string `json:"short_url"`
}
