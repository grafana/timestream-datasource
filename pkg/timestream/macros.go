package timestream

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/grafana/timestream-datasource/pkg/models"
)

const timeFilter = `\$__timeFilter`

// Interpolate processes macros
func Interpolate(query models.QueryModel) (string, error) {

	flux := query.RawQuery

	// TODO: This was just copied from NewRelic!!!!
	timeFilterExp, err := regexp.Compile(timeFilter)
	if timeFilterExp.MatchString(flux) {
		timeRange := query.TimeRange
		from := int(timeRange.From.UnixNano() / 1e6)
		to := int(timeRange.To.UnixNano() / 1e6)
		replacement := fmt.Sprintf("SINCE %s UNTIL %s", strconv.Itoa(from), strconv.Itoa(to))
		flux = timeFilterExp.ReplaceAllString(flux, replacement)
	}

	return flux, err
}
