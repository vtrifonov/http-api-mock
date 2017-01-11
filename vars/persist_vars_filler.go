package vars

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/vtrifonov/http-api-mock/definition"
	"github.com/vtrifonov/http-api-mock/logging"
	"github.com/vtrifonov/http-api-mock/persist"
	"github.com/vtrifonov/http-api-mock/utils"
)

type PersistVarsFiller struct {
	Engines     *persist.PersistEngineBag
	RegexHelper utils.RegexHelper
}

func (pvf PersistVarsFiller) Fill(m *definition.Mock, input string, multipleMatch bool) string {
	r := regexp.MustCompile(`\{\{\s*persist\.(.+?)\s*\}\}`)

	if !multipleMatch {
		return r.ReplaceAllStringFunc(input, func(raw string) string {
			// replace the strings
			if raw, found := pvf.replaceString(m, raw); found {
				return raw
			}
			// replace regexes
			return pvf.replaceRegex(m, raw)
		})
	} else {
		// first replace all strings
		input = r.ReplaceAllStringFunc(input, func(raw string) string {
			item, _ := pvf.replaceString(m, raw)
			return item
		})
		// get multiple entities using regex
		results, found := pvf.RegexHelper.GetCollectionItems(input, func(raw string) (string, string, bool) {
			return pvf.getPersistRegexParts(m, raw)
		})
		if found {
			if len(results) == 1 {
				return "," + results[0] // add a comma in the beginning so that we will now that the item is a single entity
			}

			return strings.Join(results, ",")
		}
		return input
	}
}

func (pvf PersistVarsFiller) replaceString(m *definition.Mock, raw string) (string, bool) {
	found := false
	s := ""
	tag := strings.Trim(raw[2:len(raw)-2], " ")
	if tag == "persist.entity.name" {
		s = m.Persist.Entity
		found = true
	} else if tag == "persist.entity.id" {
		s = m.Persist.EntityID
		found = true
	} else if i := strings.Index(tag, "persist.collection.name"); i == 0 {
		s = m.Persist.Collection
		found = true
	} else if i := strings.Index(tag, "persist.entity.content"); i == 0 {
		engine := pvf.Engines.Get(m.Persist.Engine)
		content, err := engine.Read(m.Persist.Entity)
		//if error, we change Response status and body
		if err != nil {
			s = ""
			m.Response.Body = ""
			m.Response.StatusCode = 404
		}
		//check for persisted entity properties like persist.entity.content.name
		if i := strings.Index(tag, "persist.entity.content."); i == 0 {
			s, found = pvf.getContentParam(content, tag[len("persist.entity.content."):])
		} else {
			s = content
			found = true
		}
	} else if i := strings.Index(tag, "persist.collection.content"); i == 0 {
		engine := pvf.Engines.Get(m.Persist.Engine)
		content, err := engine.ReadCollection(m.Persist.Collection)
		//if error, we change Response status and body
		if err != nil {
			s = ""
			m.Response.Body = ""
			m.Response.StatusCode = 404
		}
		s = content
		found = true
	}

	if !found {
		return raw, false
	}
	return s, true
}

func (pvf PersistVarsFiller) getContentParam(content string, name string) (string, bool) {
	r := regexp.MustCompile(`^(\w+\.)*\w+$`)
	if !r.MatchString(name) {
		return "", false
	}
	value, err := utils.GetPropertyValue(content, name)
	return value, err == nil
}

func (pvf PersistVarsFiller) getPersistRegexParts(m *definition.Mock, input string) (string, string, bool) {
	if i := strings.Index(input, "persist.entity.name"); i == 0 && len(input) > len("persist.entity.name") {
		return m.Persist.Entity, input[20:], true
	}
	return "", "", false
}

func (pvf PersistVarsFiller) replaceRegex(m *definition.Mock, raw string) string {
	tag := strings.Trim(raw[2:len(raw)-2], " ")
	if regexInput, regexPattern, found := pvf.getPersistRegexParts(m, tag); found {
		if result, found := pvf.RegexHelper.GetStringPart(regexInput, regexPattern, "value"); found {
			return result
		}
	}
	return raw
}

func (pvf PersistVarsFiller) callSequence(m *definition.Mock, parameters string) (string, bool) {
	regexPattern := `\(\s*(?:'|")?(?P<name>.+?)(?:'|")?\s*,\s*(?P<increase>\d+?)\s*\)|\(\s*(?:'|")?(?P<nameOnly>.+?)(?:'|")?\s*\)`

	helper := utils.RegexHelper{}

	increase := "0"
	// check first whether only name is passed to the sequence method
	name, found := helper.GetStringPart(parameters, regexPattern, "nameOnly")
	if name == "" || !found {
		name, found = helper.GetStringPart(parameters, regexPattern, "name")
		if !found {
			return "", false
		}

		increase, found = helper.GetStringPart(parameters, regexPattern, "increase")
		if !found {
			return "", false
		}

		if increase == "" {
			increase = "0"
		}
	}

	increaseInt, err := strconv.Atoi(increase)
	if err != nil {
		logging.Printf("Error parsing increase value: %s\n", err.Error())
		return "", false
	}

	engine := pvf.Engines.Get(m.Persist.Engine)

	if sequenceValue, err := engine.GetSequence(name, increaseInt); err == nil {
		return strconv.Itoa(sequenceValue), true
	} else {
		return "", false
	}
}

func (pvf PersistVarsFiller) callSetValue(m *definition.Mock, parameters string) (string, bool) {
	regexPattern := `\(\s*(?:'|")?(?P<key>.+?)(?:'|")?\s*,\s*(?:'|")?(?P<value>.+?)(?:'|")?\s*\)`

	helper := utils.RegexHelper{}

	key, found := helper.GetStringPart(parameters, regexPattern, "key")
	if !found {
		return "", false
	}

	value, found := helper.GetStringPart(parameters, regexPattern, "value")
	if !found {
		return "", false
	}

	engine := pvf.Engines.Get(m.Persist.Engine)

	if err := engine.SetValue(key, value); err == nil {
		return value, true
	} else {
		return "", false
	}
}

func (pvf PersistVarsFiller) callGetValue(m *definition.Mock, parameters string) (string, bool) {
	regexPattern := `\(\s*(?:'|")?(?P<key>.+?)(?:'|")?\s*\)`

	helper := utils.RegexHelper{}

	key, found := helper.GetStringPart(parameters, regexPattern, "key")
	if !found {
		return "", false
	}

	engine := pvf.Engines.Get(m.Persist.Engine)

	if value, err := engine.GetValue(key); err == nil {
		return value, true
	} else {
		return "", false
	}
}
