// GNU GPL v3 License
// Copyright (c) 2019 github.com:go-trellis

package tgorm

func (p *TGorm) beginNonTransaction(name string) error {
	if p.isTransaction {
		return ErrFailToConvetTXToNonTX
	}

	_db, err := p.getDB(name)
	if err != nil {
		return err
	}

	p.txSession = _db

	return nil
}

func (p *TGorm) commitNonTransaction(txFunc interface{}, name string, repos ...interface{}) error {
	if p.isTransaction {
		return ErrNonTransactionCantCommit
	}

	_funcs := GetLogicFuncs(txFunc)

	var (
		_values []interface{}
		errcode error
	)

	if _funcs.BeforeLogic != nil {
		if _, errcode = CallFunc(_funcs.BeforeLogic, _funcs, repos); errcode != nil {
			return errcode
		}
	}

	if _funcs.Logic != nil {
		if _values, errcode = CallFunc(_funcs.Logic, _funcs, repos); errcode != nil {
			return errcode
		}
	}

	if _funcs.AfterLogic != nil {
		if _values, errcode = CallFunc(_funcs.AfterLogic, _funcs, repos); errcode != nil {
			return errcode
		}
	}

	if _funcs.AfterCommit != nil {
		if _, errcode = CallFunc(_funcs.AfterCommit, _funcs, _values); errcode != nil {
			return errcode
		}
	}

	return nil
}
