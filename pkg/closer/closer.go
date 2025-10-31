package closer

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"go.uber.org/zap"
)

// shutdownTimeout по умолчанию
const shutdownTimeout = 5 * time.Second

// Closer управляет процессом graceful shutdown приложения
type Closer struct {
	mu     sync.Mutex                    // Защита от гонки при добавлении функций
	once   sync.Once                     // Гарантия однократного вызова CloseAll
	done   chan struct{}                 // Канал для оповещения о завершении
	funcs  []func(context.Context) error // Зарегистрированные функции закрытия
	logger *zap.Logger                   // Используемый логгер
}

// New создаёт новый экземпляр Closer с no-op логгером
func New(signals ...os.Signal) *Closer {
	return NewWithLogger(zap.NewNop(), signals...)
}

// NewWithLogger создаёт новый экземпляр Closer с указанным логгером.
// Если переданы сигналы, Closer начнёт их слушать и вызовет CloseAll при получении.
func NewWithLogger(logger *zap.Logger, signals ...os.Signal) *Closer {
	c := &Closer{
		done:   make(chan struct{}),
		logger: logger,
	}

	if len(signals) > 0 {
		go c.handleSignals(signals...)
	}

	return c
}

// SetLogger устанавливает логгер для Closer
func (c *Closer) SetLogger(l *zap.Logger) {
	c.logger = l
}

// Done возвращает канал, который закрывается при вызове CloseAll
func (c *Closer) Done() <-chan struct{} {
	return c.done
}

// handleSignals обрабатывает системные сигналы и вызывает CloseAll с fresh shutdown context
func (c *Closer) handleSignals(signals ...os.Signal) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	defer signal.Stop(ch)

	select {
	case <-ch:
		c.logger.Info("🛑 Получен системный сигнал, начинаем graceful shutdown")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer shutdownCancel()

		if err := c.CloseAll(shutdownCtx); err != nil {
			c.logger.Error("❌ Ошибка при закрытии ресурсов", zap.Error(err))
		}
	case <-c.done:
	}
}

// AddNamed добавляет функцию закрытия с именем зависимости для логирования
func (c *Closer) AddNamed(name string, f func(context.Context) error) {
	c.Add(func(ctx context.Context) error {
		start := time.Now()
		c.logger.Info(fmt.Sprintf("🧩 Закрываем %s", name))

		err := f(ctx)

		duration := time.Since(start)
		if err != nil {
			c.logger.Error(fmt.Sprintf("❌ Ошибка при закрытии %s", name), zap.Error(err), zap.Duration("duration", duration))
		} else {
			c.logger.Info(fmt.Sprintf("✅ %s успешно закрыт", name), zap.Duration("duration", duration))
		}
		return err
	})
}

// Add добавляет одну или несколько функций закрытия
func (c *Closer) Add(f ...func(context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.funcs = append(c.funcs, f...)
}

// CloseAll вызывает все зарегистрированные функции закрытия.
// Возвращает первую возникшую ошибку, если таковая была.
func (c *Closer) CloseAll(ctx context.Context) error {
	var result error

	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil // освободим память
		c.mu.Unlock()

		if len(funcs) == 0 {
			c.logger.Info("ℹ️ Нет функций для закрытия")
			return
		}

		c.logger.Info("🚦 Начинаем процесс graceful shutdown")

		errCh := make(chan error, len(funcs))
		var wg sync.WaitGroup

		// Выполняем в обратном порядке добавления
		for i := len(funcs) - 1; i >= 0; i-- {
			f := funcs[i]
			wg.Add(1)
			go func(f func(context.Context) error) {
				defer wg.Done()

				// Защита от паники
				defer func() {
					if r := recover(); r != nil {
						err := fmt.Errorf("panic recovered in closer: %v", r)
						errCh <- err
						c.logger.Error("⚠️ Panic в функции закрытия", zap.Any("error", r))
					}
				}()

				if err := f(ctx); err != nil {
					errCh <- err
				}
			}(f)
		}

		// Закрываем канал ошибок, когда все функции завершатся
		go func() {
			wg.Wait()
			close(errCh)
		}()

		// Читаем ошибки или отмену контекста
		for {
			select {
			case <-ctx.Done():
				c.logger.Info("⚠️ Контекст отменён во время закрытия", zap.Error(ctx.Err()))
				if result == nil {
					result = ctx.Err()
				}
				return
			case err, ok := <-errCh:
				if !ok {
					c.logger.Info("✅ Все ресурсы успешно закрыты")
					return
				}
				c.logger.Error("❌ Ошибка при закрытии", zap.Error(err))
				if result == nil {
					result = err
				}
			}
		}
	})

	return result
}
