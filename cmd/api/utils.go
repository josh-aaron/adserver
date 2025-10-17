package main

import "net/url"

func (app *application) ExtractQueryParam(param string, queryParams url.Values) string {
	paramStr := queryParams.Get(param)
	return paramStr
}
