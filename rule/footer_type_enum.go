package rule

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/conventionalcommit/commitlint/lint"
)

var _ lint.Rule = (*FooterTypeEnumRule)(nil)

// FooterTypeEnumRule to validate footer tokens
type FooterTypeEnumRule struct {
	Params []*FooterTypeEnumParam
}

// FooterTypeEnumParam represent a single footer type param
type FooterTypeEnumParam struct {
	Token  string
	Types  []string
	Values []string
}

// Name return name of the rule
func (r *FooterTypeEnumRule) Name() string { return "footer-type-enum" }

// Apply sets the needed argument for the rule
func (r *FooterTypeEnumRule) Apply(setting lint.RuleSetting) error {
	confParams, ok := setting.Argument.([]interface{})
	if !ok {
		return errInvalidArg(r.Name(), fmt.Errorf("expects array of params, but got %#v", setting.Argument))
	}

	params := make([]*FooterTypeEnumParam, 0, len(confParams))

	for index, p := range confParams {
		v, ok := p.(map[interface{}]interface{})
		if !ok {
			return errInvalidArg(r.Name()+": params", fmt.Errorf("expects key-value object, but got %#v", p))
		}

		param, err := processParam(r, v, index)
		if err != nil {
			return err
		}
		params = append(params, param)
	}

	for _, p := range params {
		// sorting the string elements for binary search
		sort.Strings(p.Types)
		sort.Strings(p.Values)
	}

	r.Params = params
	return nil
}

func processParam(r lint.Rule, val map[interface{}]interface{}, index int) (*FooterTypeEnumParam, error) {
	tok, ok := val["token"]
	if !ok {
		return nil, errMissingArg(r.Name(), "token in param "+strconv.Itoa(index+1))
	}

	types, ok := val["types"]
	if !ok {
		return nil, errMissingArg(r.Name(), "types in param "+strconv.Itoa(index+1))
	}

	values, ok := val["values"]
	if !ok {
		return nil, errMissingArg(r.Name(), "values in param "+strconv.Itoa(index+1))
	}

	param := &FooterTypeEnumParam{}

	err := setStringArg(&param.Token, tok)
	if err != nil {
		return nil, errInvalidArg(r.Name()+": token", err)
	}

	err = setStringArrArg(&param.Types, types)
	if err != nil {
		return nil, errInvalidArg(r.Name()+": types", err)
	}

	err = setStringArrArg(&param.Values, values)
	if err != nil {
		return nil, errInvalidArg(r.Name()+": values", err)
	}

	// validate the arguments
	if param.Token == "" {
		return nil, errInvalidArg(r.Name(), errors.New("token cannot be empty in param "+strconv.Itoa(index+1)))
	}

	if len(param.Types) < 1 {
		return nil, errNeedAtleastOneArg(r.Name(), "types in param "+strconv.Itoa(index+1))
	}

	if len(param.Values) < 1 {
		return nil, errNeedAtleastOneArg(r.Name(), "values in param "+strconv.Itoa(index+1))
	}

	return param, nil
}

// Validate validates FooterTypeEnumRule
func (r *FooterTypeEnumRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	var invalids []string

	// find missing footer notes
	for _, param := range r.Params {
		isType := search(param.Types, msg.Type())
		if !isType {
			continue
		}

		isNote := searchNote(msg.Notes(), param.Token)
		if !isNote {
			a := fmt.Sprintf("'%s' should exist for type '%s'", param.Token, msg.Type())
			invalids = append(invalids, a)
		}
	}

outer:
	for _, note := range msg.Notes() {
		for _, param := range r.Params {
			isType := search(param.Types, msg.Type())
			if !isType {
				// not applicable for current type
				continue
			}

			if note.Token() != param.Token {
				// not applicable for current token
				continue
			}

			for _, val := range param.Values {
				if strings.HasPrefix(note.Value(), val) {
					// has valid prefix, check next footer note
					continue outer
				}
			}

			// invalid - matches non of the mentioned prefix
			a := fmt.Sprintf("'%s' should have one of prefix [%s]", note.Token(), strings.Join(param.Values, ", "))
			invalids = append(invalids, a)
		}
	}

	if len(invalids) == 0 {
		return nil, true
	}

	desc := "footer token is invalid"
	return lint.NewIssue(desc, invalids...), false
}
