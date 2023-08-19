package parser

import "github.com/Sntree2mi8/gogqllexer"

// https://spec.graphql.org/October2021/#ImplementsInterfaces
func parseImplementsInterfaces(l *LexerWrapper) (interfaces []string, err error) {
	if err = l.SkipKeyword("implements"); err != nil {
		return nil, err
	}

	// implements at least one interface
	l.SkipIf(gogqllexer.Amp)
	if err = l.PeekAndMustBe(
		[]gogqllexer.Kind{gogqllexer.Name},
		func(t gogqllexer.Token, advanceLexer func()) error {
			defer advanceLexer()

			interfaces = append(interfaces, t.Value)
			return nil
		},
	); err != nil {
		return nil, err
	}

	// read more interfaces
	for {
		if skip := l.SkipIf(gogqllexer.Amp); !skip {
			break
		}

		if err = l.PeekAndMustBe(
			[]gogqllexer.Kind{gogqllexer.Name},
			func(t gogqllexer.Token, advanceLexer func()) error {
				defer advanceLexer()

				interfaces = append(interfaces, t.Value)
				return nil
			},
		); err != nil {
			return nil, err
		}
	}

	return interfaces, nil
}
