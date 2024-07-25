package timestream

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/experimental/errorsource"
	"github.com/grafana/timestream-datasource/pkg/models"
	"golang.org/x/exp/maps"
)

// TODO: consider refactoring sqlutil.Interpolate to be more generic and using that instead

type macroFunc func(models.QueryModel, models.DatasourceSettings) (string, error)

var macroFuncs = map[string]macroFunc{
	"timeFilter":      macroTimeFilter,
	"timeFrom":        macroTimeFrom,
	"timeTo":          macroTimeTo,
	"interval":        macroInterval,
	"interval_ms":     macroInterval,
	"interval_raw_ms": macroIntervalRaw,
	"now_ms":          macroNow,
	"database":        macroDatabase,
	"table":           macroTable,
	"measure":         macroMeasure,
}

var macroKeys []string

func init() {
	// sort macro keys longest first, so shorter keys don't clobber longer keys
	// they're a prefix of
	macroKeys = maps.Keys(macroFuncs)
	slices.SortFunc(macroKeys, func(a, b string) int { return len(b) - len(a) })
}

func macroTimeFilter(model models.QueryModel, _ models.DatasourceSettings) (string, error) {
	from := model.TimeRange.From.UnixNano() / 1e6
	to := model.TimeRange.To.UnixNano() / 1e6

	replacement := fmt.Sprintf("time BETWEEN from_milliseconds(%d) AND from_milliseconds(%d)", from, to)
	return replacement, nil
}

func macroTimeFrom(model models.QueryModel, _ models.DatasourceSettings) (string, error) {
	return fmt.Sprintf("%d", model.TimeRange.From.UnixNano()/1e6), nil
}

func macroTimeTo(model models.QueryModel, _ models.DatasourceSettings) (string, error) {
	return fmt.Sprintf("%d", model.TimeRange.To.UnixNano()/1e6), nil
}

func macroInterval(model models.QueryModel, _ models.DatasourceSettings) (string, error) {
	if model.Interval.Milliseconds() == 0 {
		return "", fmt.Errorf("invalid interval: %dns", model.Interval.Nanoseconds())
	}
	return fmt.Sprintf("%dms", model.Interval.Milliseconds()), nil
}

func macroIntervalRaw(model models.QueryModel, _ models.DatasourceSettings) (string, error) {
	if model.Interval.Milliseconds() == 0 {
		return "", fmt.Errorf("invalid interval: %dns", model.Interval.Nanoseconds())
	}
	return fmt.Sprintf("%d", model.Interval.Milliseconds()), nil
}

func macroNow(_ models.QueryModel, _ models.DatasourceSettings) (string, error) {
	now := time.Now().UnixMilli()
	return fmt.Sprintf("%d", now), nil
}

func macroDatabase(model models.QueryModel, settings models.DatasourceSettings) (string, error) {
	return valueOrDefault(model.Database, settings.DefaultDatabase), nil
}
func macroTable(model models.QueryModel, settings models.DatasourceSettings) (string, error) {
	return valueOrDefault(model.Table, settings.DefaultTable), nil
}
func macroMeasure(model models.QueryModel, settings models.DatasourceSettings) (string, error) {
	return valueOrDefault(model.Measure, settings.DefaultMeasure), nil
}

func valueOrDefault(value string, defaultValue string) string {
	if value == "" || strings.HasPrefix(value, "${") {
		return defaultValue
	}
	return value
}

// Interpolate processes macros
func Interpolate(model models.QueryModel, settings models.DatasourceSettings) (string, error) {
	query := model.RawQuery
	for _, key := range macroKeys {
		macroKey := fmt.Sprintf("$__%s", key)
		if !strings.Contains(query, macroKey) {
			continue
		}
		replacement, err := macroFuncs[key](model, settings)
		if err != nil {
			return query, errorsource.DownstreamError(err, false)
		}
		query = strings.ReplaceAll(query, macroKey, replacement)
	}
	return query, nil
}
