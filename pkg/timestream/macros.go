package timestream

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/grafana/timestream-datasource/pkg/models"
)

const timeFilter = `\$__timeFilter`
const timeFromStr = `$__timeFrom`
const timeToStr = `$__timeTo`
const intervalStrAlias = `$__interval`
const intervalStr = `$__interval_ms`
const intervalRawStr = `$__interval_raw_ms`
const nowStr = `$__now_ms`

// WHERE time > from_unixtime(unixtime)
// WHERE time > from_iso8601_timestamp(iso_8601_string_format)
// WHERE time > from_milliseconds(epoch_millis)

// Interpolate processes macros
func Interpolate(query models.QueryModel, settings models.DatasourceSettings) (string, error) {

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
		from := int64(timeRange.From.UnixNano() / 1e6)
		to := int64(timeRange.To.UnixNano() / 1e6)

		replacement := fmt.Sprintf("time BETWEEN from_milliseconds(%d) AND from_milliseconds(%d)", from, to)
		txt = timeFilterExp.ReplaceAllString(txt, replacement)
	}

	if strings.Contains(txt, timeFromStr) {
		timeRange := query.TimeRange
		from := int64(timeRange.From.UnixNano() / 1e6)
		replacement := fmt.Sprintf("%d", from)
		txt = strings.ReplaceAll(txt, timeFromStr, replacement)
	}

	if strings.Contains(txt, timeToStr) {
		timeRange := query.TimeRange
		to := int64(timeRange.To.UnixNano() / 1e6)
		replacement := fmt.Sprintf("%d", to)
		txt = strings.ReplaceAll(txt, timeToStr, replacement)
	}

	if strings.Contains(txt, intervalRawStr) {
		replacement := fmt.Sprintf("%d", query.Interval.Milliseconds())
		txt = strings.ReplaceAll(txt, intervalRawStr, replacement)
	}

	if strings.Contains(txt, intervalStr) || strings.Contains(txt, intervalStrAlias) {
		replacement := fmt.Sprintf("%dms", query.Interval.Milliseconds())
		if replacement == "0ms" {
			replacement = "{!invalid interval=" + query.Interval.String() + "!}"
		}
		txt = strings.ReplaceAll(txt, intervalStr, replacement)
		txt = strings.ReplaceAll(txt, intervalStrAlias, replacement)
	}

	if strings.Contains(txt, nowStr) {
		now := int(time.Now().UnixNano() / int64(time.Millisecond))
		replacement := fmt.Sprintf("%d", now)
		txt = strings.ReplaceAll(txt, nowStr, replacement)
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
