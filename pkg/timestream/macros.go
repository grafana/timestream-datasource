package timestream

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/grafana/timestream-datasource/pkg/common/aws"
	"github.com/grafana/timestream-datasource/pkg/models"
)

const timeFilter = `\$__timeFilter`
const intervalStr = `$__interval_ms`

// WHERE time > from_unixtime(unixtime)
// WHERE time > from_iso8601_timestamp(iso_8601_string_format)
// WHERE time > from_milliseconds(epoch_millis)

// Interpolate processes macros
func Interpolate(query models.QueryModel, settings aws.DatasourceSettings) (string, error) {

	txt := query.RawQuery

	if strings.Contains(txt, "$__intervalStr") {
		return txt, fmt.Errorf("$__intervalStr removed... use $__interval_ms")
	}

	// Simple Macros
	txt = replaceOrDefault(txt, "$__database", query.Database, settings.DefaultDatabase)
	txt = replaceOrDefault(txt, "$__table", query.Table, settings.DefaultTable)
	txt = replaceOrDefault(txt, "$__measure", query.Measure, settings.DefaultMeasure)

	timeFilterExp, err := regexp.Compile(timeFilter)
	if err != nil {
		return txt, err
	}

	if timeFilterExp.MatchString(txt) {
		timeRange := query.TimeRange
		from := int(timeRange.From.UnixNano() / 1e6)
		to := int(timeRange.To.UnixNano() / 1e6)

		replacement := fmt.Sprintf("time BETWEEN from_milliseconds(%d) AND from_milliseconds(%d)", from, to)
		txt = timeFilterExp.ReplaceAllString(txt, replacement)
	}

	if strings.Contains(txt, intervalStr) {
		replacement := fmt.Sprintf("%dms", query.Interval.Milliseconds())
		if replacement == "0ms" {
			replacement = "{!invalid interval=" + query.Interval.String() + "!}"
		}
		txt = strings.ReplaceAll(txt, intervalStr, replacement)
	}

	return txt, err
}

func replaceOrDefault(txt string, key string, v1 string, v2 string) string {
	val := v1
	if val == "" || strings.HasPrefix(val, "${") {
		val = v2
	}
	if val == "" {
		return txt // no change
	}
	return strings.ReplaceAll(txt, key, val)
}
