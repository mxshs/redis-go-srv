package parser

import (
	"fmt"
	"mxshs/redis-go/app/types"
	"strconv"
	"strings"
)

type Parser struct {
	Offset int
	Text   string
}

func newParser(text string) *Parser {
	return &Parser{
		Text:   text,
		Offset: 0,
	}
}

func Parse(text string) ([]*types.Data, error) {
    p := newParser(text)

	msgs := make([]*types.Data, 0)

	for p.Offset < len(p.Text) {
		start := p.Offset
		msg, err := p.parse()
		if err != nil {
			return nil, err
		}

		msg.Sz = len(p.Text[start:min(len(p.Text), p.Offset)])

		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func (p *Parser) parse() (*types.Data, error) {
	switch p.Text[p.Offset] {
	case '+':
		p.Offset++
		return p.parseString()
	case '$':
		p.Offset++
		return p.parseBulkString()
	case ':':
		p.Offset++
		return p.parseInt()
	case '*':
		p.Offset++
		return p.parseArray()
	default:
		return nil, fmt.Errorf("%s is not a valid redis message", p.Text[p.Offset:])
	}
}

func (p *Parser) parseString() (*types.Data, error) {
	end := strings.Index(p.Text[p.Offset:], "\r")
	if end == -1 {
		return nil, fmt.Errorf("%s is not a valid redis string", p.Text[p.Offset:])
	}

	val := p.Text[p.Offset : p.Offset + end]
	p.Offset = p.Offset + end + 2

	return &types.Data{
		Value: val,
		T:     types.String,
	}, nil
}

func (p *Parser) parseBulkString() (*types.Data, error) {
	szIdx := strings.Index(p.Text[p.Offset:], "\r")
	if szIdx == -1 {
		return nil, fmt.Errorf("%s is not a valid redis bulk string", p.Text[p.Offset:])
	}

	sz, err := strconv.Atoi(p.Text[p.Offset : p.Offset + szIdx])
	if err != nil {
		return nil, err
	}

	p.Offset = p.Offset + szIdx + 2

    if p.Offset + sz > len(p.Text) {
		return nil, fmt.Errorf("%s is not a valid redis bulk string", p.Text[p.Offset:])
    }

	val := p.Text[p.Offset : p.Offset + sz]
	p.Offset = p.Offset + sz + 2

    return &types.Data{
        Value: val,
        T: types.String,
    }, nil
}

func (p *Parser) parseInt() (*types.Data, error) {
	end := strings.Index(p.Text[p.Offset:], "\r")
	if end == -1 {
		return nil, fmt.Errorf("%s is not a valid redis int", p.Text[p.Offset:])
	}

	val := p.Text[p.Offset : p.Offset + end]
	p.Offset = p.Offset + end + 2

	return &types.Data{
		Value: val,
		T:     types.Int,
	}, nil
}

func (p *Parser) parseArray() (*types.Data, error) {
	szIdx := strings.Index(p.Text[p.Offset:], "\r")
	if szIdx == -1 {
		return nil, fmt.Errorf("%s is not a valid redis array", p.Text[p.Offset:])
	}

	sz, err := strconv.Atoi(p.Text[p.Offset : p.Offset + szIdx])
	if err != nil {
		return nil, err
	}

	p.Offset = p.Offset + szIdx + 2

	arr := make([]*types.Data, sz)
	for cur := 0; cur < sz; cur++ {
		if p.Offset >= len(p.Text) {
			return nil, fmt.Errorf("%s is not a valid redis array", p.Text)
		}

		res, err := p.parse()
		if err != nil {
			return nil, err
		}

		arr[cur] = res
	}

	return &types.Data{
		Value: arr,
		T:     types.Array,
	}, nil
}
