/*
Copyright © 2019 Henry Huang <hhh@rutcode.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package tgorm

import (
	"reflect"

	"gorm.io/gorm"
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
		panic(ErrNotFoundDefaultDatabase)
	}
}

func (p *TGorm) getDB(name string) (*gorm.DB, error) {
	if db, _exist := p.dbs[name]; _exist {
		return db, nil
	}
	return nil, ErrNotFoundGormDB
}

func getRepo(v interface{}) *TGorm {
	_deepRepo := DeepFields(v, reflect.TypeOf(new(TGorm)), []reflect.Value{})
	if deepRepo, ok := _deepRepo.(*TGorm); ok {
		return deepRepo
	}
	return nil
}

func createNewTGorm(origin interface{}) (*TGorm, interface{}, error) {

	if repo, err := Derive(origin); err != nil {
		return nil, nil, err
	} else if repo != nil {
		return getRepo(repo), repo, nil
	}

	newRepoV := reflect.New(reflect.ValueOf(
		reflect.Indirect(reflect.ValueOf(origin)).Interface()).Type())
	if !newRepoV.IsValid() {
		return nil, nil, ErrFailToCreateRepo
	}

	newRepoI := newRepoV.Interface()
	newTgorm := getRepo(newRepoI)

	if err := Inherit(newRepoI, origin); err != nil {
		return nil, nil, err
	}

	if newTgorm == nil {
		return nil, nil, ErrFailToConvetTXToNonTX
	}
	return newTgorm, newRepoI, nil
}
