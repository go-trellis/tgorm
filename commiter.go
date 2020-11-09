// GNU GPL v3 License
// Copyright (c) 2019 github.com:go-trellis

package tgorm

// committer gorm committer
type committer struct {
	Name string
}

// NewCommitter get trellis gorm committer
func NewCommitter() Committer {
	return &committer{Name: "go-trellis::tgorm::committer"}
}

// NonTX do non transaction function by default database
func (p *committer) NonTX(fn interface{}, repos ...interface{}) error {
	return p.NonTXWithName(fn, DefaultDatabase, repos...)
}

// NonTXWithName do non transaction function with name of database
func (p *committer) NonTXWithName(fn interface{}, name string, repos ...interface{}) error {
	if err := p.checkRepos(fn, repos); err != nil {
		return err
	}

	_newRepos := []interface{}{}
	_newTGormRepos := []*TGorm{}

	for _, origin := range repos {
		repo := getRepo(origin)
		if repo == nil {
			return ErrStructCombineWithRepo
		}

		_newTgorm, _newRepoI, err := createNewTGorm(origin)
		if err != nil {
			return err
		}

		_newRepos = append(_newRepos, _newRepoI)

		_newTgorm.dbs = repo.dbs
		_newTgorm.defaultDB = repo.defaultDB

		if err := _newTgorm.beginNonTransaction(name); err != nil {
			return err
		}

		_newTGormRepos = append(_newTGormRepos, _newTgorm)
	}

	return _newTGormRepos[0].commitNonTransaction(fn, name, _newRepos...)
}

// TX do transaction function by default database
func (p *committer) TX(fn interface{}, repos ...interface{}) error {
	return p.TXWithName(fn, DefaultDatabase, repos...)
}

// TXWithName do transaction function with name of database
func (p *committer) TXWithName(fn interface{}, name string, repos ...interface{}) error {
	if err := p.checkRepos(fn, repos); err != nil {
		return err
	}

	_newRepos := []interface{}{}
	_newTGormRepos := []*TGorm{}

	for _, origin := range repos {

		repo := getRepo(origin)
		if repo == nil {
			return ErrStructCombineWithRepo
		}

		_newTgorm, _newRepoI, err := createNewTGorm(origin)
		if err != nil {
			return err
		}

		_newTgorm.dbs = repo.dbs
		_newTgorm.defaultDB = repo.defaultDB
		_newRepos = append(_newRepos, _newRepoI)
		_newTGormRepos = append(_newTGormRepos, _newTgorm)
	}

	if err := _newTGormRepos[0].beginTransaction(name); err != nil {
		return err
	}

	for i := range _newTGormRepos {
		_newTGormRepos[i].txSession = _newTGormRepos[0].txSession
		_newTGormRepos[i].isTransaction = _newTGormRepos[0].isTransaction
	}

	return _newTGormRepos[0].commitTransaction(fn, _newRepos...)
}

func (*committer) checkRepos(txFunc interface{}, originRepos ...interface{}) error {
	if reposLen := len(originRepos); reposLen < 1 {
		return ErrAtLeastOneRepo
	}

	if txFunc == nil {
		return ErrNotFoundTransationFunction
	}
	return nil
}
