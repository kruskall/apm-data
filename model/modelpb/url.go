package modelpb

import (
	"net/url"
	"strconv"
)

func ParseURL(original, defaultHostname, defaultScheme string) *URL {
	original = truncate(original)
	url, err := url.Parse(original)
	if err != nil {
		return &URL{Original: original}
	}
	if url.Scheme == "" {
		url.Scheme = defaultScheme
		if url.Scheme == "" {
			url.Scheme = "http"
		}
	}
	if url.Host == "" {
		url.Host = defaultHostname
	}
	out := URL{
		Original: original,
		Scheme:   url.Scheme,
		Full:     truncate(url.String()),
		Domain:   truncate(url.Hostname()),
		Path:     truncate(url.Path),
		Query:    truncate(url.RawQuery),
		Fragment: url.Fragment,
	}
	if port := url.Port(); port != "" {
		if intv, err := strconv.Atoi(port); err == nil {
			out.Port = uint32(intv)
		}
	}
	return &out
}

// truncate returns s truncated at n runes, and the number of runes in the resulting string (<= n).
func truncate(s string) string {
	var j int
	for i := range s {
		if j == 1024 {
			return s[:i]
		}
		j++
	}
	return s
}
