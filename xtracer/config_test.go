package xtracer

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
)

var (
	app = "Demo"
)

func TestDefa(t *testing.T) {
	c := DefaultConfig()
	assert.Equal(t, EXPORTER_NAME_STDOUT, c.ExporterName)
	assert.Equal(t, 1.0, c.SampleRate)
}

func TestTracer(t *testing.T) {
	c := DefaultConfig()
	c.ServiceName = app
	fileName := "trace.json"
	f, err := os.Create(fileName)
	assert.Nil(t, err, "create file error")
	assert.NotNil(t, f, "file is nil")
	defer f.Close()
	exp := NewFileExporter(f)
	tp := NewProvider(c, exp)
	assert.NotNil(t, tp)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	err = parent(context.Background())
	assert.Nil(t, err)
	tp.ForceFlush(context.Background())
	data, err := os.ReadFile(fileName)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	err = os.Remove(fileName)
	assert.Nil(t, err)
}

// **** function 01 ****
func parent(ctx context.Context) error {
	newCtx, span := otel.Tracer(app).Start(ctx, "Run")
	defer span.End()
	return child(newCtx)
}

// **** function 02 ****
func child(ctx context.Context) error {
	newCtx, span := otel.Tracer(app).Start(ctx, "Run1")
	defer span.End()
	return child02(newCtx)
}

// **** function 03 ****
func child02(ctx context.Context) error {
	fmt.Print("hello, world\n")
	_, span := otel.Tracer(app).Start(ctx, "Run1")
	defer span.End()
	return nil
}
