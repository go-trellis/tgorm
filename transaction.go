// GNU GPL v3 License
// Copyright (c) 2019 github.com:go-trellis

package tgorm

import (
	"strings"

	"github.com/go-trellis/errors"
)

func (p *TGorm) beginTransaction(name string) errors.ErrorCode {
	if !p.isTransaction {
		p.isTransaction = true
		_db, err := p.getDB(name)
		if err != nil {
			return err
		}
		p.txSession = _db
		return nil
	}
	return ErrTransactionIsAlreadyBegin.New(errors.Params{"name": name})
}

func (p *TGorm) commitTransaction(txFunc interface{}, repos ...interface{}) errors.ErrorCode {
	if !p.isTransaction {
		return ErrNonTransactionCantCommit.New()
	}

	if p.txSession == nil {
		return ErrTransactionSessionIsNil.New()
	}

	if txFunc == nil {
		return ErrNotFoundTransationFunction.New()
	}

	_isNeedRollBack := true
	p.txSession = p.txSession.Begin()

	if errs := p.txSession.GetErrors(); 0 != len(errs) {
		var errStrings []string
		for _, err := range errs {
			errStrings = append(errStrings, err.Error())
		}
		return ErrFailToCreateTransaction.New(
			errors.Params{"message": strings.Join(errStrings, ";\n")})
	}

	defer func() {
		if _isNeedRollBack {
			p.txSession.Rollback()
		}
	}()

	_funcs := GetLogicFuncs(txFunc)

	var (
		values []interface{}
		ecode  errors.ErrorCode
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
	if errs := p.txSession.Commit().GetErrors(); 0 != len(errs) {
		return ErrFailToCommitTransaction.New().Append(errs)
	}

	if _funcs.AfterCommit != nil {
		if _, ecode = CallFunc(_funcs.AfterCommit, _funcs, values); ecode != nil {
			return ecode
		}
	}

	return nil
}
