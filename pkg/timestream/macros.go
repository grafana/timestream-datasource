package timestream

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/grafana/timestream-datasource/pkg/models"
)

const timeFilter = `\$__timeFilter`
const intervalStr = `\$__intervalStr`

// WHERE time > from_unixtime(unixtime)
// WHERE time > from_iso8601_timestamp(iso_8601_string_format)
// WHERE time > from_milliseconds(epoch_millis)

// Interpolate processes macros
func Interpolate(query models.QueryModel) (string, error) {

	txt := query.RawQuery

	// Simple Macros
	txt = replaceIfNotVar(txt, "${database}", query.Database)
	txt = replaceIfNotVar(txt, "${table}", query.Table)
	txt = replaceIfNotVar(txt, "${measure}", query.Measure)

	timeFilterExp, err := regexp.Compile(timeFilter)
	if err != nil {
		return txt, err
	}
	intervalStrExp, err := regexp.Compile(intervalStr)
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

	if intervalStrExp.MatchString(txt) {
		replacement := query.Interval.String()
		txt = intervalStrExp.ReplaceAllString(txt, replacement)
	}

	return txt, err
}

func replaceIfNotVar(txt string, key string, val string) string {
	if val == "" || strings.HasPrefix(val, "$") {
		return txt // unchanged
	}
	return strings.ReplaceAll(txt, key, val)
}
