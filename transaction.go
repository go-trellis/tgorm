// GNU GPL v3 License
// Copyright (c) 2019 github.com:go-trellis

package tgorm

func (p *TGorm) beginTransaction(name string) error {
	if !p.isTransaction {
		p.isTransaction = true
		_db, err := p.getDB(name)
		if err != nil {
			return err
		}
		p.txSession = _db
		return nil
	}
	return ErrTransactionIsAlreadyBegin
}

func (p *TGorm) commitTransaction(txFunc interface{}, repos ...interface{}) error {
	if !p.isTransaction {
		return ErrNonTransactionCantCommit
	}

	if p.txSession == nil {
		return ErrTransactionSessionIsNil
	}

	if txFunc == nil {
		return ErrNotFoundTransationFunction
	}

	_isNeedRollBack := true
	p.txSession = p.txSession.Begin()

	if p.txSession.Error != nil {
		return p.txSession.Error
	}

	defer func() {
		if _isNeedRollBack {
			p.txSession.Rollback()
		}
	}()

	_funcs := GetLogicFuncs(txFunc)

	var (
		values []interface{}
		ecode  error
	)

	if _funcs.BeforeLogic != nil {
		if _, ecode = CallFunc(_funcs.BeforeLogic, _funcs, repos); ecode != nil {
			return ecode
		}
	}

	if _funcs.Logic != nil {
		if values, ecode = CallFunc(_funcs.Logic, _funcs, repos); ecode != nil {
			return ecode
		}
	}

	if _funcs.AfterLogic != nil {
		if values, ecode = CallFunc(_funcs.AfterLogic, _funcs, repos); ecode != nil {
			return ecode
		}
	}

	_isNeedRollBack = false
	if err := p.txSession.Commit().Error; err != nil {
		return err
	}

	if _funcs.AfterCommit != nil {
		if _, ecode = CallFunc(_funcs.AfterCommit, _funcs, values); ecode != nil {
			return ecode
		}
	}

	return nil
}
