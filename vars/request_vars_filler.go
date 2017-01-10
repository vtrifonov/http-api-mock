package vars

import (
	"regexp"
	"strings"

	urlmatcher "github.com/azer/url-router"
	"github.com/vtrifonov/http-api-mock/definition"
	"github.com/vtrifonov/http-api-mock/utils"
)

type RequestVarsFiller struct {
	Request     *definition.Request
	Mock        *definition.Mock
	RegexHelper utils.RegexHelper
}

func (rvf RequestVarsFiller) Fill(m *definition.Mock, input string, multipleMatch bool) string {
	r := regexp.MustCompile(`\{\{\s*request\.(.+?)\s*\}\}`)

	if !multipleMatch {
		return r.ReplaceAllStringFunc(input, func(raw string) string {
			// replace the strings
			if raw, found := rvf.replaceString(raw); found {
				return raw
			}
			// replace regexes
			return rvf.replaceRegex(raw)
		})
	} else {
		// first replace all strings
		input = r.ReplaceAllStringFunc(input, func(raw string) string {
			item, _ := rvf.replaceString(raw)
			return item
		})
		// get multiple entities using regex
		if results, found := rvf.RegexHelper.GetCollectionItems(input, rvf.getVarsRegexParts); found {
			if len(results) == 1 {
				return "," + results[0] // add a comma in the beginning so that we will now that the item is a single entity
			}

			return strings.Join(results, ",")
		}
		return input
	}
}

func (rvf RequestVarsFiller) replaceString(raw string) (string, bool) {
	found := false
	s := ""
	tag := strings.Trim(raw[2:len(raw)-2], " ")
	if tag == "request.body" {
		s = rvf.Request.Body
		found = true
	} else if i := strings.Index(tag, "request.body."); i == 0 {
		s, found = rvf.getBodyParam(rvf.Request, tag[len("request.body."):])
	} else if i := strings.Index(tag, "request.query."); i == 0 {
		s, found = rvf.getQueryStringParam(rvf.Request, tag[len("request.query."):])
	} else if i := strings.Index(tag, "request.path."); i == 0 {
		s, found = rvf.getPathParam(tag[len("request.path."):])
	} else if i := strings.Index(tag, "request.cookie."); i == 0 {
		s, found = rvf.getCookieParam(rvf.Request, tag[len("request.cookie."):])
	}
	if !found {
		return raw, false
	}
	return s, true
}

func (rvf RequestVarsFiller) getVarsRegexParts(input string) (string, string, bool) {
	if i := strings.Index(input, "request.url."); i == 0 {
		return rvf.Request.Path, input[12:], true
	} else if i := strings.Index(input, "request.body."); i == 0 {
		return rvf.Request.Body, input[13:], true
	}
	return "", "", false
}

func (rvf RequestVarsFiller) replaceRegex(raw string) string {
	tag := strings.Trim(raw[2:len(raw)-2], " ")
	if regexInput, regexPattern, found := rvf.getVarsRegexParts(tag); found {
		if result, found := rvf.RegexHelper.GetStringPart(regexInput, regexPattern, "value"); found {
			return result
		}
	}
	return raw
}

func (rvf RequestVarsFiller) getPathParam(name string) (string, bool) {

	routes := urlmatcher.New(rvf.Mock.Request.Path)
	mparm := routes.Match(rvf.Request.Path)

	value, f := mparm.Params[name]
	if !f {
		return "", false
	}

	return value, true
}

func (rvf RequestVarsFiller) getQueryStringParam(req *definition.Request, name string) (string, bool) {

	if len(rvf.Request.QueryStringParameters) == 0 {
		return "", false
	}
	value, f := rvf.Request.QueryStringParameters[name]
	if !f {
		return "", false
	}

	return value[0], true
}

func (rvf RequestVarsFiller) getBodyParam(req *definition.Request, name string) (string, bool) {

	value, err := utils.GetPropertyValue(req.Body, name)
	return value, err == nil
}

func (rvf RequestVarsFiller) getCookieParam(req *definition.Request, name string) (string, bool) {

	if len(rvf.Request.Cookies) == 0 {
		return "", false
	}
	value, f := rvf.Request.Cookies[name]
	if !f {
		return "", false
	}

	return value, true
}
