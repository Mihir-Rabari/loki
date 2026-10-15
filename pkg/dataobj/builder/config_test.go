package builder

import (
	"flag"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConfig_RegisterFlags(t *testing.T) {
	var cfg Config
	f := flag.NewFlagSet("test", flag.PanicOnError)
	cfg.RegisterFlags(f)
	cfg.Topic = "test"
	require.NoError(t, cfg.Validate())

	require.Equal(t, 1<<20, int(cfg.LogsobjBuilder.TargetPageSize))
	require.Equal(t, 8, cfg.LogsobjBuilder.EstimatedCompressionRatio)
	require.Equal(t, "8", f.Lookup("dataobj.builder.logsobj-builder.estimated-compression-ratio").DefValue)

	require.Equal(t, 128<<10, int(cfg.IndexobjBuilder.TargetPageSize))
	require.Equal(t, 1, cfg.IndexobjBuilder.EstimatedCompressionRatio)
	require.Equal(t, "1", f.Lookup("dataobj.builder.indexobj-builder.estimated-compression-ratio").DefValue)

	require.NoError(t, f.Parse([]string{"-dataobj.builder.indexobj-builder.estimated-compression-ratio=4"}))
	require.Equal(t, 4, cfg.IndexobjBuilder.EstimatedCompressionRatio)
}

func TestConfig_Validate(t *testing.T) {
	newConfig := func(t *testing.T, args ...string) Config {
		t.Helper()
		var cfg Config
		f := flag.NewFlagSet("test", flag.PanicOnError)
		cfg.RegisterFlags(f)
		require.NoError(t, f.Parse(append([]string{"-dataobj.builder.topic=test"}, args...)))
		return cfg
	}

	t.Run("accepts the default required replicas report", func(t *testing.T) {
		cfg := newConfig(t)
		require.NoError(t, cfg.Validate())
		require.True(t, cfg.ReportRequiredReplicas)
		require.Equal(t, 30*time.Second, cfg.ReportRequiredReplicasInterval)
	})

	t.Run("rejects a zero required replicas interval when the report is enabled", func(t *testing.T) {
		cfg := newConfig(t, "-dataobj.builder.report-required-replicas-interval=0s")
		require.Error(t, cfg.Validate())
	})

	t.Run("accepts a zero required replicas interval when the report is disabled", func(t *testing.T) {
		cfg := newConfig(t,
			"-dataobj.builder.report-required-replicas=false",
			"-dataobj.builder.report-required-replicas-interval=0s",
		)
		require.NoError(t, cfg.Validate())
	})
}
