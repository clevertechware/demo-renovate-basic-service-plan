package main

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestEnvOrReturnsFallbackWhenVariableIsEmpty(t *testing.T) {
	t.Setenv("PLAN_TEST_VAR", "")

	if got := envOr("PLAN_TEST_VAR", "fallback"); got != "fallback" {
		t.Errorf("envOr = %q, want fallback", got)
	}
}

func TestEnvOrReturnsVariableWhenSet(t *testing.T) {
	t.Setenv("PLAN_TEST_VAR", "value")

	if got := envOr("PLAN_TEST_VAR", "fallback"); got != "value" {
		t.Errorf("envOr = %q, want value", got)
	}
}

func TestServeReturnsNilOnceContextIsCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	server := &http.Server{Addr: "127.0.0.1:0", Handler: http.NotFoundHandler()}
	done := make(chan error, 1)
	go func() { done <- serve(ctx, server) }()
	time.Sleep(100 * time.Millisecond)

	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("serve returned %v, want nil", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("serve did not return after context cancellation")
	}
}
