package model

// HTTPServerConfig структура конфига HTTPServer
type HTTPServerConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
}

// ShortServiceConfig структура конфига ShortService
type ShortServiceConfig struct {
	BaseURL string `env:"BASE_URL"`
}

// RepositoryConfig структура конфига Repository
type RepositoryConfig struct {
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

// Concurrency структура конфига Concurrency
type Concurrency struct {
	WorkerPoolDelete *WorkerPoolDelete
	WorkerPoolEvent  *WorkerPoolEvent
}

// WorkerPoolDelete структура конфига WorkerPoolDelete
type WorkerPoolDelete struct {
	CountWorkers   int32
	InputChainSize int
	BufferSize     int
	BatchSize      int
}

// WorkerPoolEvent структура конфига WorkerPoolEvent
type WorkerPoolEvent struct {
	CountWorkers   int32
	EventChainSize int
}

// AuditConfig структура конфига Audit
type AuditConfig struct {
	AuditFile string `env:"AUDIT_FILE"`
	AuditUrl  string `env:"AUDIT_URL"`
}
