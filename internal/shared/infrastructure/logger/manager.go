package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

type Manager struct {
	baseDir     string
	appLogger   *slog.Logger
	httpLogger  *slog.Logger
	modelLogger map[string]*slog.Logger
	files       []*os.File
	mu          sync.Mutex
}

var defaultManager *Manager

func Init(baseDir string) error {
	if baseDir == "" {
		baseDir = "logs"
	}

	manager := &Manager{
		baseDir:     baseDir,
		modelLogger: make(map[string]*slog.Logger, len(AllModels)),
	}

	if err := os.MkdirAll(filepath.Join(baseDir, "models"), 0o755); err != nil {
		return fmt.Errorf("no se pudo crear el directorio de logs: %w", err)
	}

	appFile, err := openLogFile(filepath.Join(baseDir, "app.log"))
	if err != nil {
		return err
	}
	manager.files = append(manager.files, appFile)
	manager.appLogger = newLogger(appFile)

	httpFile, err := openLogFile(filepath.Join(baseDir, "http.log"))
	if err != nil {
		_ = manager.Close()
		return err
	}
	manager.files = append(manager.files, httpFile)
	manager.httpLogger = newLogger(httpFile)

	for _, model := range AllModels {
		modelFile, err := openLogFile(filepath.Join(baseDir, "models", model+".log"))
		if err != nil {
			_ = manager.Close()
			return err
		}
		manager.files = append(manager.files, modelFile)
		manager.modelLogger[model] = newLogger(modelFile)
	}

	defaultManager = manager
	slog.SetDefault(manager.appLogger)
	defaultManager.App().Info("logger inicializado", slog.String("base_dir", baseDir))
	return nil
}

func Close() error {
	if defaultManager == nil {
		return nil
	}
	return defaultManager.Close()
}

func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var closeErr error
	for _, file := range m.files {
		if err := file.Close(); err != nil && closeErr == nil {
			closeErr = err
		}
	}
	m.files = nil
	defaultManager = nil
	return closeErr
}

func App() *slog.Logger {
	if defaultManager == nil {
		return slog.Default()
	}
	return defaultManager.appLogger
}

func HTTP() *slog.Logger {
	if defaultManager == nil {
		return slog.Default()
	}
	return defaultManager.httpLogger
}

func Model(model string) *slog.Logger {
	if defaultManager == nil {
		return slog.Default()
	}
	if logger, ok := defaultManager.modelLogger[model]; ok {
		return logger
	}
	return defaultManager.appLogger
}

func (m *Manager) App() *slog.Logger {
	return m.appLogger
}

func openLogFile(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir archivo de log %s: %w", path, err)
	}
	return file, nil
}
