package vars

import (
	"regexp"
	"strings"

	urlmatcher "github.com/azer/url-router"
	"github.com/vtrifonov/http-api-mock/definition"
	"github.com/vtrifonov/http-api-mock/utils"
)

type RequestVars struct {
	Request     *definition.Request
	Mock        *definition.Mock
	RegexHelper utils.RegexHelper
}

func (rv RequestVars) Fill(m *definition.Mock, input string, multipleMatch bool) string {
	r := regexp.MustCompile(`\{\{\s*request\.(.+?)\s*\}\}`)

	if !multipleMatch {
		return r.ReplaceAllStringFunc(input, func(raw string) string {
			// replace the strings
			if raw, found := rv.replaceString(raw); found {
				return raw
			}
			// replace regexes
			return rv.replaceRegex(raw)
		})
	} else {
		// first replace all strings
		input = r.ReplaceAllStringFunc(input, func(raw string) string {
			item, _ := rv.replaceString(raw)
			return item
		})
		// get multiple entities using regex
		if results, found := rv.RegexHelper.GetCollectionItems(input, rv.getVarsRegexParts); found {
			if len(results) == 1 {
				return "," + results[0] // add a comma in the beginning so that we will now that the item is a single entity
			}

			return strings.Join(results, ",")
		}
		return input
	}
}

func (rv RequestVars) replaceString(raw string) (string, bool) {
	found := false
	s := ""
	tag := strings.Trim(raw[2:len(raw)-2], " ")
	if tag == "request.body" {
		s = rv.Request.Body
		found = true
	} else if i := strings.Index(tag, "request.query."); i == 0 {
		s, found = rv.getQueryStringParam(rv.Request, tag[len("request.query."):])
	} else if i := strings.Index(tag, "request.path."); i == 0 {
		s, found = rv.getPathParm(tag[len("request.path."):])
	} else if i := strings.Index(tag, "request.cookie."); i == 0 {
		s, found = rv.getCookieParam(rv.Request, tag[len("request.cookie."):])
	}
	if !found {
		return raw, false
	}
	return s, true
}

func (rv RequestVars) getVarsRegexParts(input string) (string, string, bool) {
	if i := strings.Index(input, "request.url."); i == 0 {
		return rv.Request.Path, input[12:], true
	} else if i := strings.Index(input, "request.body."); i == 0 {
		return rv.Request.Body, input[13:], true
	}
	return "", "", false
}

func (rv RequestVars) replaceRegex(raw string) string {
	tag := strings.Trim(raw[2:len(raw)-2], " ")
	if regexInput, regexPattern, found := rv.getVarsRegexParts(tag); found {
		if result, found := rv.RegexHelper.GetStringPart(regexInput, regexPattern, "value"); found {
			return result
		}
	}
	return raw
}

func (rv RequestVars) getPathParm(name string) (string, bool) {

	routes := urlmatcher.New(rv.Mock.Request.Path)
	mparm := routes.Match(rv.Request.Path)

	value, f := mparm.Params[name]
	if !f {
		return "", false
	}

	return value, true
}

func (rv RequestVars) getQueryStringParam(req *definition.Request, name string) (string, bool) {

	if len(rv.Request.QueryStringParameters) == 0 {
		return "", false
	}
	value, f := rv.Request.QueryStringParameters[name]
	if !f {
		return "", false
	}

	return value[0], true
}

func (rv RequestVars) getCookieParam(req *definition.Request, name string) (string, bool) {

	if len(rv.Request.Cookies) == 0 {
		return "", false
	}
	value, f := rv.Request.Cookies[name]
	if !f {
		return "", false
	}

	return value, true
}
