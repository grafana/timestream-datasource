package timestream

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/timestream-datasource/pkg/common/aws"
	"github.com/grafana/timestream-datasource/pkg/models"
)

func TestInterpolate(t *testing.T) {
	// Unix sec: 1500376552
	// Unix ms:  1500376552001

	timeRange := backend.TimeRange{
		From: time.Unix(0, 1500376552001*1e6),
		To:   time.Unix(0, 1500376552002*1e6),
	}

	t.Run("interpolate __timeFilter function", func(t *testing.T) {
		sqltxt := `SELECT average(value) FROM test AND $__timeFilter TIMESERIES`
		expect := `SELECT average(value) FROM test AND time BETWEEN from_milliseconds(1500376552001) AND from_milliseconds(1500376552002) TIMESERIES`

		query := models.QueryModel{
			TimeRange: timeRange,
			RawQuery:  sqltxt,
		}
		text, _ := Interpolate(query, aws.DatasourceSettings{})
		if diff := cmp.Diff(text, expect); diff != "" {
			t.Fatalf("Result mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("using interval", func(t *testing.T) {
		sqltxt := `GROUP BY $__interval_ms TIMESERIES`
		expect := `GROUP BY 60000ms TIMESERIES`

		query := models.QueryModel{
			TimeRange: timeRange,
			RawQuery:  sqltxt,
			Interval:  time.Minute,
		}
		text, _ := Interpolate(query, aws.DatasourceSettings{})
		if diff := cmp.Diff(text, expect); diff != "" {
			t.Fatalf("Result mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("using templates", func(t *testing.T) {
		sqltxt := `SELECT $__measure FROM $__database.$__table LIMIT 10`
		expect := `SELECT measure FROM ddb.table LIMIT 10`

		query := models.QueryModel{
			TimeRange: timeRange,
			RawQuery:  sqltxt,
			Interval:  time.Minute,
			Database:  "${ddd}", // should use default
			Table:     "table",
		}
		text, _ := Interpolate(query, aws.DatasourceSettings{
			DefaultDatabase: "ddb",
			DefaultTable:    "dtb",
			DefaultMeasure:  "measure",
		})
		if diff := cmp.Diff(text, expect); diff != "" {
			t.Fatalf("Result mismatch (-want +got):\n%s", diff)
		}
	})
}
