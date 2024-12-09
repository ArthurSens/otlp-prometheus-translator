package translator

import (
	"testing"
)

func TestBuildMetricName(t *testing.T) {
	tests := []struct {
		name              string
		metricName        string
		unit              string
		metricType        MetricType
		addMetricSuffixes bool
		expected          string
	}{
		{
			name:              "simple metric without suffixes",
			metricName:        "http_requests",
			unit:              "",
			metricType:        MetricTypeGauge,
			addMetricSuffixes: false,
			expected:          "http_requests",
		},
		{
			name:              "counter with total suffix",
			metricName:        "http_requests",
			unit:              "",
			metricType:        MetricTypeMonotonicCounter,
			addMetricSuffixes: true,
			expected:          "http_requests_total",
		},
		{
			name:              "gauge with time unit",
			metricName:        "request_duration",
			unit:              "s",
			metricType:        MetricTypeGauge,
			addMetricSuffixes: true,
			expected:          "request_duration_seconds",
		},
		{
			name:              "counter with time unit",
			metricName:        "request_duration",
			unit:              "ms",
			metricType:        MetricTypeMonotonicCounter,
			addMetricSuffixes: true,
			expected:          "request_duration_milliseconds_total",
		},
		{
			name:              "gauge with compound unit",
			metricName:        "throughput",
			unit:              "By/s",
			metricType:        MetricTypeGauge,
			addMetricSuffixes: true,
			expected:          "throughput_bytes_per_second",
		},
		{
			name:              "ratio metric",
			metricName:        "cpu_utilization",
			unit:              "1",
			metricType:        MetricTypeGauge,
			addMetricSuffixes: true,
			expected:          "cpu_utilization_ratio",
		},
		{
			name:              "counter with unit 1 (no ratio suffix)",
			metricName:        "error_count",
			unit:              "1",
			metricType:        MetricTypeMonotonicCounter,
			addMetricSuffixes: true,
			expected:          "error_count_total",
		},
		{
			name:              "metric with byte units",
			metricName:        "memory_usage",
			unit:              "MiBy",
			metricType:        MetricTypeGauge,
			addMetricSuffixes: true,
			expected:          "memory_usage_mebibytes",
		},
		{
			name:              "metric with SI units",
			metricName:        "temperature",
			unit:              "Cel",
			metricType:        MetricTypeGauge,
			addMetricSuffixes: true,
			expected:          "temperature_celsius",
		},
		{
			name:              "metric with dots",
			metricName:        "system.cpu.usage",
			unit:              "1",
			metricType:        MetricTypeGauge,
			addMetricSuffixes: true,
			expected:          "system.cpu.usage_ratio",
		},
		{
			name:              "metric with japanese characters (memory usage rate)",
			metricName:        "メモリ使用率", // memori shiyouritsu (memory usage rate) xD
			unit:              "By",
			metricType:        MetricTypeGauge,
			addMetricSuffixes: true,
			expected:          "メモリ使用率_bytes",
		},
		{
			name:              "metric with mixed special characters (system.memory.usage.rate)",
			metricName:        "system.メモリ.usage.率", // system.memory.usage.rate
			unit:              "By/s",
			metricType:        MetricTypeGauge,
			addMetricSuffixes: true,
			expected:          "system.メモリ.usage.率_bytes_per_second",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildMetricName(tt.metricName, tt.unit, tt.metricType, tt.addMetricSuffixes)
			if tt.expected != result {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestBuildUnitSuffixes(t *testing.T) {
	tests := []struct {
		name            string
		unit            string
		expectedMain    string
		expectedPerUnit string
	}{
		{
			name:            "empty unit",
			unit:            "",
			expectedMain:    "",
			expectedPerUnit: "",
		},
		{
			name:            "simple time unit",
			unit:            "s",
			expectedMain:    "seconds",
			expectedPerUnit: "",
		},
		{
			name:            "compound unit",
			unit:            "By/s",
			expectedMain:    "bytes",
			expectedPerUnit: "per_second",
		},
		{
			name:            "complex compound unit",
			unit:            "KiBy/m",
			expectedMain:    "kibibytes",
			expectedPerUnit: "per_minute",
		},
		{
			name:            "unit with spaces",
			unit:            " ms / s ",
			expectedMain:    "milliseconds",
			expectedPerUnit: "per_second",
		},
		{
			name:            "invalid unit",
			unit:            "invalid",
			expectedMain:    "invalid",
			expectedPerUnit: "",
		},
		{
			name:            "unit with curly braces",
			unit:            "{custom}/s",
			expectedMain:    "",
			expectedPerUnit: "per_second",
		},
		{
			name:            "multiple slashes",
			unit:            "By/s/h",
			expectedMain:    "bytes",
			expectedPerUnit: "per_s/h",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mainUnit, perUnit := buildUnitSuffixes(tt.unit)
			if tt.expectedMain != mainUnit {
				t.Errorf("expected main unit %s, got %s", tt.expectedMain, mainUnit)
			}
			if tt.expectedPerUnit != perUnit {
				t.Errorf("expected per unit %s, got %s", tt.expectedPerUnit, perUnit)
			}
		})
	}
}

func TestBuildCompliantMetricName(t *testing.T) {
	tests := []struct {
		name              string
		metricName        string
		unit              string
		metricType        MetricType
		addMetricSuffixes bool
		expected          string
	}{
		{
			name:              "simple valid metric name",
			metricName:        "http_requests",
			unit:              "",
			metricType:        MetricTypeGauge,
			addMetricSuffixes: false,
			expected:          "http_requests",
		},
		{
			name:              "metric name with invalid characters",
			metricName:        "http-requests@in_flight",
			unit:              "",
			metricType:        MetricTypeNonMonotonicCounter,
			addMetricSuffixes: false,
			expected:          "http_requests_in_flight",
		},
		{
			name:              "metric name starting with digit",
			metricName:        "5xx_errors",
			unit:              "",
			metricType:        MetricTypeGauge,
			addMetricSuffixes: false,
			expected:          "_5xx_errors",
		},
		{
			name:              "metric name with multiple consecutive invalid chars",
			metricName:        "api..//request--time",
			unit:              "",
			metricType:        MetricTypeGauge,
			addMetricSuffixes: false,
			expected:          "api_request_time",
		},
		{
			name:              "full normalization with units and type",
			metricName:        "system.cpu-utilization",
			unit:              "ms/s",
			metricType:        MetricTypeMonotonicCounter,
			addMetricSuffixes: true,
			expected:          "system_cpu_utilization_milliseconds_per_second_total",
		},
		{
			name:              "metric with special characters and ratio",
			metricName:        "memory.usage%rate",
			unit:              "1",
			metricType:        MetricTypeGauge,
			addMetricSuffixes: true,
			expected:          "memory_usage_rate_ratio",
		},
		{
			name:              "metric with unicode characters",
			metricName:        "error_rate_£_€_¥",
			unit:              "",
			metricType:        MetricTypeGauge,
			addMetricSuffixes: false,
			expected:          "error_rate_____",
		},
		{
			name:              "metric with multiple spaces",
			metricName:        "api   response   time",
			unit:              "ms",
			metricType:        MetricTypeGauge,
			addMetricSuffixes: true,
			expected:          "api_response_time_milliseconds",
		},
		{
			name:              "metric with colons (valid prometheus chars)",
			metricName:        "app:request:latency",
			unit:              "s",
			metricType:        MetricTypeGauge,
			addMetricSuffixes: true,
			expected:          "app:request:latency_seconds",
		},
		{
			name:              "empty metric name",
			metricName:        "",
			unit:              "",
			metricType:        MetricTypeGauge,
			addMetricSuffixes: false,
			expected:          "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildCompliantMetricName(tt.metricName, tt.unit, tt.metricType, tt.addMetricSuffixes)
			if tt.expected != result {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
