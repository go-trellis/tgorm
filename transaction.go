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
