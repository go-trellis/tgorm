// GNU GPL v3 License
// Copyright (c) 2019 github.com:go-trellis

package tgorm

import (
	"github.com/go-trellis/errors"
)

func (p *TGorm) beginNonTransaction(name string) errors.ErrorCode {
	if p.isTransaction {
		return ErrFailToConvetTXToNonTX.New()
	}

	_db, err := p.getDB(name)
	if err != nil {
		return err
	}

	p.txSession = _db

	return nil
}

func (p *TGorm) commitNonTransaction(txFunc interface{}, name string, repos ...interface{}) errors.ErrorCode {
	if p.isTransaction {
		return ErrNonTransactionCantCommit.New()
	}

	_funcs := GetLogicFuncs(txFunc)

	var (
		_values []interface{}
		errcode errors.ErrorCode
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
