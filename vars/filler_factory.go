package vars

import (
	"github.com/vtrifonov/http-api-mock/definition"
	"github.com/vtrifonov/http-api-mock/persist"
	"github.com/vtrifonov/http-api-mock/utils"
	"github.com/vtrifonov/http-api-mock/vars/fakedata"
)

type FillerFactory interface {
	CreateRequestFiller(req *definition.Request, mock *definition.Mock) Filler
	CreateFakeFiller(Fake fakedata.DataFaker) Filler
	CreateStorageFiller(Engines *persist.PersistEngineBag) Filler
	CreatePersistFiller(Engines *persist.PersistEngineBag) Filler
}

type MockFillerFactory struct{}

func (mff MockFillerFactory) CreateRequestFiller(req *definition.Request, mock *definition.Mock) Filler {
	return RequestVarsFiller{Request: req, Mock: mock, RegexHelper: utils.RegexHelper{}}
}

func (mff MockFillerFactory) CreateFakeFiller(fake fakedata.DataFaker) Filler {
	return FakeVarsFiller{Fake: fake}
}

func (mff MockFillerFactory) CreateStorageFiller(engines *persist.PersistEngineBag) Filler {
	return StorageVarsFiller{Engines: engines, RegexHelper: utils.RegexHelper{}}
}

func (mff MockFillerFactory) CreatePersistFiller(engines *persist.PersistEngineBag) Filler {
	return PersistVarsFiller{Engines: engines, RegexHelper: utils.RegexHelper{}}
}
