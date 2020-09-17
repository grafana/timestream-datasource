package common

// RoundInterval picks a reasonable interval value for a time range
// See: https://github.com/grafana/grafana/blob/master/packages/grafana-data/src/datetime/rangeutil.ts#L330
func RoundInterval(interval int64) int64 {
	switch true {
	// 0.015s
	case interval < 15:
		return 10 // 0.01s
	// 0.035s
	case interval < 35:
		return 20 // 0.02s
	// 0.075s
	case interval < 75:
		return 50 // 0.05s
	// 0.15s
	case interval < 150:
		return 100 // 0.1s
	// 0.35s
	case interval < 350:
		return 200 // 0.2s
	// 0.75s
	case interval < 750:
		return 500 // 0.5s
	// 1.5s
	case interval < 1500:
		return 1000 // 1s
	// 3.5s
	case interval < 3500:
		return 2000 // 2s
	// 7.5s
	case interval < 7500:
		return 5000 // 5s
	// 12.5s
	case interval < 12500:
		return 10000 // 10s
	// 17.5s
	case interval < 17500:
		return 15000 // 15s
	// 25s
	case interval < 25000:
		return 20000 // 20s
	// 45s
	case interval < 45000:
		return 30000 // 30s
	// 1.5m
	case interval < 90000:
		return 60000 // 1m
	// 3.5m
	case interval < 210000:
		return 120000 // 2m
	// 7.5m
	case interval < 450000:
		return 300000 // 5m
	// 12.5m
	case interval < 750000:
		return 600000 // 10m
	// 12.5m
	case interval < 1050000:
		return 900000 // 15m
	// 25m
	case interval < 1500000:
		return 1200000 // 20m
	// 45m
	case interval < 2700000:
		return 1800000 // 30m
	// 1.5h
	case interval < 5400000:
		return 3600000 // 1h
	// 2.5h
	case interval < 9000000:
		return 7200000 // 2h
	// 4.5h
	case interval < 16200000:
		return 10800000 // 3h
	// 9h
	case interval < 32400000:
		return 21600000 // 6h
	// 1d
	case interval < 86400000:
		return 43200000 // 12h
	// 1w
	case interval < 604800000:
		return 86400000 // 1d
	// 3w
	case interval < 1814400000:
		return 604800000 // 1w
	// 6w
	case interval < 3628800000:
		return 2592000000 // 30d
	default:
		return 31536000000 // 1y
	}
}
