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

type StorageVarsFiller struct {
	Engines     *persist.PersistEngineBag
	RegexHelper utils.RegexHelper
}

func (lvf StorageVarsFiller) Fill(m *definition.Mock, input string, multipleMatch bool) string {
	r := regexp.MustCompile(`\{\{\s*storage\.([^{]+?)\s*\}\}`)
	tries := 0
	// we are making several passes while we have matching regex, this is useful for cases when we have nested vars like {{ storage.SetValue({{ request.body.username\\=(?P<value>.+?)(?:&|$) }}, {{ storage.Sequence(users, 1) }}) }}.json
	for tries <= 3 && r.MatchString(input) {
		input = lvf.Process(r, m, input)
		tries++
	}
	return input
}

func (lvf StorageVarsFiller) Process(r *regexp.Regexp, m *definition.Mock, input string) string {
	return r.ReplaceAllStringFunc(input, func(raw string) string {
		found := false
		s := ""
		tag := strings.Trim(raw[2:len(raw)-2], " ")
		if i := strings.Index(tag, "storage.Sequence"); i == 0 {
			s, found = lvf.callSequence(m, tag[len("storage.Sequence"):])
		} else if i := strings.Index(tag, "storage.GetValue"); i == 0 {
			s, found = lvf.callGetValue(m, tag[len("storage.GetValue"):])
		} else if i := strings.Index(tag, "storage.SetValue"); i == 0 {
			s, found = lvf.callSetValue(m, tag[len("storage.SetValue"):])
		}

		if !found {
			return raw
		}
		return s
	})
}

func (lvf StorageVarsFiller) callSequence(m *definition.Mock, parameters string) (string, bool) {
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

	engine := lvf.Engines.Get(m.Persist.Engine)

	if sequenceValue, err := engine.GetSequence(name, increaseInt); err == nil {
		return strconv.Itoa(sequenceValue), true
	} else {
		return "", false
	}
}

func (lvf StorageVarsFiller) callSetValue(m *definition.Mock, parameters string) (string, bool) {
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

	engine := lvf.Engines.Get(m.Persist.Engine)

	if err := engine.SetValue(key, value); err == nil {
		return value, true
	} else {
		return "", false
	}
}

func (lvf StorageVarsFiller) callGetValue(m *definition.Mock, parameters string) (string, bool) {
	regexPattern := `\(\s*(?:'|")?(?P<key>.+?)(?:'|")?\s*\)`

	helper := utils.RegexHelper{}

	key, found := helper.GetStringPart(parameters, regexPattern, "key")
	if !found {
		return "", false
	}

	engine := lvf.Engines.Get(m.Persist.Engine)

	if value, err := engine.GetValue(key); err == nil {
		return value, true
	} else {
		return "", false
	}
}
