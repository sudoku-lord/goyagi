package logger

// defaultLogger is the global logger.
var defaultLogger = New()

// Info writes a info-level log with a message and any additional data provided.
func Info(message string, fields ...Data) {
	defaultLogger.Info(message, fields...)
}

// Error writes an error-level log with a message and any additional data
// provided.
func Error(message string, fields ...Data) {
	defaultLogger.Error(message, fields...)
}

// Warn writes a warn-level log with a message and any additional data provided.
func Warn(message string, fields ...Data) {
	defaultLogger.Warn(message, fields...)
}

// Debug writes a debug-level log with a message and any additional data
// provided.
func Debug(message string, fields ...Data) {
	defaultLogger.Debug(message, fields...)
}

// Fatal writes a fatal-level log with a message and any additional data
// provided. This will also call `os.Exit(1)`.`
func Fatal(message string, fields ...Data) {
	defaultLogger.Fatal(message, fields...)
}
