package formatter

import "github.com/MohaCodez/structured-logger/logger"

type Formatter interface {
	Format(entry *logger.Entry) ([]byte, error)
}
