// GNU GPL v3 License
// Copyright (c) 2019 github.com:go-trellis

package tgorm

import (
	"reflect"

	"github.com/go-trellis/errors"
	"github.com/jinzhu/gorm"
)

// TGorm trellis gorm
type TGorm struct {
	isTransaction bool
	txSession     *gorm.DB

	dbs       map[string]*gorm.DB
	defaultDB *gorm.DB
}

// Session get session
func (p *TGorm) Session() *gorm.DB {
	return p.txSession
}

// SetDBs set gorm dbs
func (p *TGorm) SetDBs(dbs map[string]*gorm.DB) {
	if defDB, _exist := dbs[DefaultDatabase]; _exist {
		p.dbs = dbs
		p.defaultDB = defDB
	} else {
		panic(ErrNotFoundDefaultDatabase.New())
	}
}

func (p *TGorm) getDB(name string) (*gorm.DB, errors.ErrorCode) {
	if db, _exist := p.dbs[name]; _exist {
		return db, nil
	}
	return nil, ErrNotFoundGormDB.New(errors.Params{"name": name})
}

func getRepo(v interface{}) *TGorm {
	_deepRepo := DeepFields(v, reflect.TypeOf(new(TGorm)), []reflect.Value{})
	if deepRepo, ok := _deepRepo.(*TGorm); ok {
		return deepRepo
	}
	return nil
}

func createNewTGorm(origin interface{}) (*TGorm, interface{}, errors.ErrorCode) {

	if repo, err := Derive(origin); err != nil {
		return nil, nil, ErrFailToDerive.New(errors.Params{"message": err.Error()})
	} else if repo != nil {
		return getRepo(repo), repo, nil
	}

	newRepoV := reflect.New(reflect.ValueOf(
		reflect.Indirect(reflect.ValueOf(origin)).Interface()).Type())
	if !newRepoV.IsValid() {
		return nil, nil, ErrFailToCreateRepo.New()
	}

	newRepoI := newRepoV.Interface()
	newTgorm := getRepo(newRepoI)

	if err := Inherit(newRepoI, origin); err != nil {
		return nil, nil, ErrFailToInherit.New(errors.Params{"message": err.Error()})
	}

	if newTgorm == nil {
		return nil, nil, ErrFailToConvetTXToNonTX.New()
	}
	return newTgorm, newRepoI, nil
}
