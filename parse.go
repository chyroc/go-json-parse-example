package go_json_parse_example

import (
	"fmt"
	"strings"
)

type jsonParser struct {
	data []rune
	idx  int
}

func (r *jsonParser) parse() (interface{}, error) {
	if len(r.data) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	itemType := r.data[r.idx]
	switch {
	case itemType == '"':
		res, err := r.parseString()
		return res, err
	case strings.Contains("0123456789-", string([]rune{itemType})):
		res, err := r.parseNumber()
		return res, err
	case itemType == '{':
		res, err := r.parseObject()
		return res, err
	case itemType == '[':
		res, err := r.parseArray()
		return res, err
	case strings.Contains("tf", string([]rune{itemType})):
		res, err := r.parseBoolean()
		return res, err
	case itemType == 'n':
		err := r.parseNull()
		return nil, err
	default:
		return nil, fmt.Errorf("invalid item-type, pos: %d", r.idx)
	}
}

func (r *jsonParser) parseString() (string, error) {
	if err := r.findRune(false, '"'); err != nil {
		return "", err
	}
	res := []rune{}
	for r.idx < len(r.data) {
		switch r.data[r.idx] {
		case '\\':
			if r.idx == len(r.data)-1 {
				return "", fmt.Errorf("invalid string-item, At least one character is required after the escape character, pos: %d", r.idx)
			}
			r.idx++
			res = append(res, r.data[r.idx])
			r.idx++
		case '"':
			r.idx++
			return string(res), nil
		default:
			res = append(res, r.data[r.idx])
			r.idx++
		}
	}
	return "", fmt.Errorf("invalid string-item, expect: \", pos: %d", r.idx)
}

func (r *jsonParser) parseNumber() (int64, error) {
	prefix := r.data[r.idx] == '-'
	if prefix {
		r.idx++
	}
	if !strings.Contains("0123456789", string([]rune{r.data[r.idx]})) {
		return 0, fmt.Errorf("expect: 0-9, pos: %d", r.idx)
	}

	i := int64(0)
	for r.idx < len(r.data) {
		dd := r.data[r.idx]
		if dd >= '0' && dd <= '9' {
			i = i*10 + int64(dd-'0')
			r.idx++
		} else {
			break
		}
	}

	if prefix {
		i *= -1
	}
	return i, nil
}

func (r *jsonParser) parseObject() (resp map[string]interface{}, err error) {
	resp = map[string]interface{}{}

	if err := r.findRune(true, '{'); err != nil {
		return nil, err
	}
	if err := r.findRune(true, '}'); err == nil {
		return resp, nil
	}
	for {
		key, err := r.parseString()
		if err != nil {
			return nil, err
		}

		if err := r.findRune(true, ':'); err != nil {
			return nil, err
		}

		val, err := r.parse()
		if err != nil {
			return nil, err
		}
		resp[key] = val

		if err := r.findRune(true, ','); err != nil {
			break
		}
	}
	if err := r.findRune(true, '}'); err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *jsonParser) parseArray() (resp []interface{}, err error) {
	if err := r.findRune(true, '['); err != nil {
		return nil, err
	}
	if err := r.findRune(true, ']'); err == nil {
		return []interface{}{}, nil
	}
	for {
		s, err := r.parse()
		if err != nil {
			return nil, err
		}
		resp = append(resp, s)
		if err := r.findRune(true, ','); err != nil {
			break
		}
	}
	if err := r.findRune(true, ']'); err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *jsonParser) parseBoolean() (bool, error) {
	if err := r.findRune(false, 't', 'r', 'u', 'e'); err == nil {
		return true, nil
	}
	if err := r.findRune(false, 'f', 'a', 'l', 's', 'e'); err == nil {
		return false, nil
	}
	return false, fmt.Errorf("invalid boolean-item, pos: %d", r.idx)
}

func (r *jsonParser) parseNull() error {
	if err := r.findRune(false, 'n', 'u', 'l', 'l'); err == nil {
		return nil
	}
	return fmt.Errorf("invalid null-item, pos: %d", r.idx)
}

func (r *jsonParser) findRune(isKey bool, rs ...rune) error {
	if isKey {
		r.removeSpace()
	}
	c := 0
	for i := r.idx; i < len(r.data) && i-r.idx >= 0 && i-r.idx < len(rs); i++ {
		if r.data[i] == rs[i-r.idx] {
			c++
			continue
		}
		return fmt.Errorf("expect key rune: %s, pos: %d", string(rs), r.idx)
	}
	r.idx += c
	if isKey {
		r.removeSpace()
	}
	return nil
}

func (r *jsonParser) removeSpace() (n int) {
	for i := r.idx; i < len(r.data); i++ {
		if r.data[i] != ' ' && r.data[i] != '\n' {
			return
		}
		r.idx++
		n++
	}
	return
}
