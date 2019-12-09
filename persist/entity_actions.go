package persist

import (
	"github.com/Jeffail/gabs"
	"github.com/vtrifonov/http-api-mock/definition"
	"github.com/vtrifonov/http-api-mock/logging"
	"github.com/vtrifonov/http-api-mock/utils"
)

//FilePersister persists body in file
type EntityActions struct {
	Engines *PersistEngineBag
}

func (ea EntityActions) getEngine(m *definition.Mock) EntityPersister {
	engine := ea.Engines.Get(m.Persist.Engine)
	//fix persiste console logging
	if m.Persist.Engine == "" {
		m.Persist.Engine = engine.GetName()
	}
	return engine
}

func (ea EntityActions) ApplyActions(m *definition.Mock) {
	engine := ea.getEngine(m)
	fileName := m.Persist.Entity
	if value, ok := m.Persist.Actions["write"]; ok {
		if err := engine.Write(fileName, value); err != nil {
			logging.Println("Error writing in a entity")
			return
		}
	}

	if value, ok := m.Persist.Actions["append"]; ok {
		var content string
		var err error
		if m.Persist.Collection != "" {
			content, err = engine.ReadCollection(m.Persist.Collection)
		} else {
			content, err = engine.Read(fileName)
		}
		if err != nil {
			logging.Println("Error reading in a entity")
			return
		}
		if utils.IsJSON(content) && utils.IsJSON(value) {

			v1, _ := gabs.ParseJSON([]byte(content))
			v2, _ := gabs.ParseJSON([]byte(value))

			content = utils.JoinJSON(v1, v2).String()
		} else if utils.IsJSON(content) && !utils.IsJSON(value) {
			logging.Printf("There is no way to append this : %s\n", value)
		} else {
			content += value
		}
		if err := engine.Write(fileName, content); err != nil {
			logging.Println("Error appending in a entity")
		}
	}

	if _, ok := m.Persist.Actions["delete"]; ok {
		if m.Persist.Collection != "" {
			if err := engine.DeleteCollection(m.Persist.Collection); err != nil {
				logging.Println("Error deleting collection")
				return
			}
		} else {
			if err := engine.Delete(fileName); err != nil {
				logging.Println("Error deleting a entity")
				return
			}
		}
	}
}
