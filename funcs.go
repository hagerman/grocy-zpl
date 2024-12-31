package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const (
	day  = 24 * time.Hour
	year = 365 * day
)

var TemplateFuncs = template.FuncMap{
	// Time functions
	"now":            time.Now,
	"timeSince":      time.Since,
	"timeUntil":      time.Until,
	"formatTime":     formatTime,
	"approxDuration": approxDuration,
	"afterEpoch":     afterEpoch,

	// String functions
	"uppercase": strings.ToUpper,
	"lowercase": strings.ToLower,
	"pluralize": pluralize,

	// Slice functions
	"join": strings.Join,

	// Number functions
	"incr":        incr,
	"decr":        decr,
	"formatInt":   formatInt,
	"formatFloat": formatFloat,

	// Boolean functions
	"yesno": yesno,
}

func formatTime(format string, t time.Time) string {
	return t.Format(format)
}

func approxDuration(d time.Duration) string {
	if d < time.Second {
		return "less than 1 second"
	}

	ds := int(math.Round(d.Seconds()))
	if ds == 1 {
		return "1 second"
	} else if ds < 60 {
		return fmt.Sprintf("%d seconds", ds)
	}

	dm := int(math.Round(d.Minutes()))
	if dm == 1 {
		return "1 minute"
	} else if dm < 60 {
		return fmt.Sprintf("%d minutes", dm)
	}

	dh := int(math.Round(d.Hours()))
	if dh == 1 {
		return "1 hour"
	} else if dh < 24 {
		return fmt.Sprintf("%d hours", dh)
	}

	dd := int(math.Round(float64(d / day)))
	if dd == 1 {
		return "1 day"
	} else if dd < 365 {
		return fmt.Sprintf("%d days", dd)
	}

	dy := int(math.Round(float64(d / year)))
	if dy == 1 {
		return "1 year"
	}

	return fmt.Sprintf("%d years", dy)
}

func afterEpoch(t time.Time) bool {
	return t.After(time.Time{})
}

func pluralize(count any, singular string, plural string) (string, error) {
	n, err := toInt64(count)
	if err != nil {
		return "", err
	}

	if n == 1 {
		return singular, nil
	}

	return plural, nil
}

func incr(i any) (int64, error) {
	n, err := toInt64(i)
	if err != nil {
		return 0, err
	}

	n++
	return n, nil
}

func decr(i any) (int64, error) {
	n, err := toInt64(i)
	if err != nil {
		return 0, err
	}

	n--
	return n, nil
}

func formatInt(i any) (string, error) {
	n, err := toInt64(i)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", n), nil
}

func formatFloat(f float64, dp int) string {
	format := "%." + strconv.Itoa(dp) + "f"
	return fmt.Sprintf(format, f)
}

func yesno(b bool) string {
	if b {
		return "Yes"
	}

	return "No"
}

func toInt64(i any) (int64, error) {
	switch v := i.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	// Note: uint64 not supported due to risk of truncation.
	case string:
		return strconv.ParseInt(v, 10, 64)
	}

	return 0, fmt.Errorf("unable to convert type %T to int", i)
}
