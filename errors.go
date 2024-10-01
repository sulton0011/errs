package errs

// errorString - улучшенная структура ошибки, теперь хранит только одно сообщение и оригинальную ошибку
type errorString struct {
	message string
	origErr string
}

// New возвращает новую ошибку, включающую сообщение и оригинальную ошибку.
func New(message string) error {
	return &errorString{
		message: message,
		origErr: message,
	}
}

// Error реализует интерфейс error, возвращает сообщение об ошибке.
func (e *errorString) Error() string {
	return e.origErr
}
