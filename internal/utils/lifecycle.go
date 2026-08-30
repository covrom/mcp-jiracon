package utils

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// LifecycleManager coordinates graceful shutdown via signal handlers
// and a shared context. Mirrors Python's utils/lifecycle.py.
type LifecycleManager struct {
	ctx    context.Context
	cancel context.CancelFunc

	mu     sync.Mutex
	closed bool
}

// NewLifecycleManager creates a LifecycleManager rooted at ctx.
func NewLifecycleManager() *LifecycleManager {
	ctx, cancel := context.WithCancel(context.Background())
	lm := &LifecycleManager{ctx: ctx, cancel: cancel}
	lm.installSignalHandlers()
	return lm
}

// Context returns the cancellation context.
func (lm *LifecycleManager) Context() context.Context { return lm.ctx }

// Cancel signals shutdown. Idempotent.
func (lm *LifecycleManager) Cancel() {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	if lm.closed {
		return
	}
	lm.closed = true
	lm.cancel()
}

// installSignalHandlers wires SIGINT/SIGTERM to cancellation. SIGPIPE
// is ignored (matches Python's handling).
func (lm *LifecycleManager) installSignalHandlers() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
		lm.Cancel()
	}()
}

// EnsureCleanExit flushes stdout/stderr and exits. Best-effort.
func EnsureCleanExit(code int) {
	_ = os.Stdout.Sync()
	_ = os.Stderr.Sync()
	os.Exit(code)
}

// ShutdownServer gracefully shuts down srv within timeout.
func ShutdownServer(srv *http.Server, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := srv.Shutdown(ctx)
	if errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	return err
}
