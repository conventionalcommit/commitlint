package lint

import "github.com/conventionalcommit/parser"

type defaultParser struct {
	p *parser.Parser
}

func newParser() *defaultParser {
	return &defaultParser{
		p: parser.New(),
	}
}

func (p defaultParser) Parse(input string) (Commit, error) {
	c, err := p.p.Parse(input)
	if err != nil {
		return nil, err
	}
	wrapC := &defaultCommit{
		Commit: c,
	}
	return wrapC, nil
}

type defaultCommit struct {
	*parser.Commit
}

func (d *defaultCommit) Notes() []Note {
	outNotes := d.Commit.Notes()
	notes := make([]Note, len(outNotes))

	for i := 0; i < len(outNotes); i++ {
		notes[i] = &outNotes[i]
	}

	return notes
}
