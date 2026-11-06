package rules

import (
	"slices"
	"strings"

	"github.com/ecsphp/ecs-src/blink/internal/fixer"
	"github.com/ecsphp/ecs-src/blink/internal/token"
	"github.com/ecsphp/ecs-src/blink/internal/tokens"
)

// PHP-CS-Fixer: https://github.com/PHP-CS-Fixer/PHP-CS-Fixer/blob/master/src/Fixer/AttributeNotation/GeneralAttributeRemoveFixer.php
//
// GeneralAttributeRemove removes the attributes whose fully qualified name is
// in the "attributes" list, and the trailing comma of every attribute block.
type GeneralAttributeRemove struct {
	attributes []string
}

func (GeneralAttributeRemove) Name() string {
	return `PhpCsFixer\Fixer\AttributeNotation\GeneralAttributeRemoveFixer`
}

func (GeneralAttributeRemove) SourceURL() string {
	return "https://github.com/PHP-CS-Fixer/PHP-CS-Fixer/blob/master/src/Fixer/AttributeNotation/GeneralAttributeRemoveFixer.php"
}

func (f GeneralAttributeRemove) WithConfig(config map[string]any) fixer.Fixer {
	if v, ok := config["attributes"].([]any); ok {
		f.attributes = make([]string, 0, len(v))
		for _, name := range v {
			if name, ok := name.(string); ok {
				f.attributes = append(f.attributes, name)
			}
		}
	}
	return f
}

func (GeneralAttributeRemove) IsCandidate(s *tokens.Stream) bool {
	return s.Index().HasKind(token.AttributeStart)
}

func (f GeneralAttributeRemove) Fix(s *tokens.Stream) bool {
	version := s.Version()
	sc := sharedSymScan(s)
	for index := attrNextStart(s, 0); index >= 0; {
		b, ok := attrCollectOne(s, index)
		if !ok {
			break
		}
		removed := 0
		for _, it := range b.items {
			if !slices.Contains(f.attributes, garFullyQualifiedName(sc, it)) {
				continue
			}
			for j := it.start; j <= it.end; j++ {
				clearTokenAt(s, j)
			}
			removed++
			if next := nonEmptySibling(s, it.end, 1); isPunctVal(s, next, ",") {
				clearTokenAt(s, next)
			}
		}
		if removed == len(b.items) {
			for j := b.start; j <= b.end; j++ {
				clearTokenAt(s, j)
			}
		}
		// once the whole block is cleared, this is the token before it, as upstream
		if prev := prevMeaningfulIndex(s, b.close); isPunctVal(s, prev, ",") {
			clearTokenAt(s, prev)
		}
		index = attrNextStart(s, b.end)
	}
	return s.Version() != version
}

// garFullyQualifiedName is AttributeAnalyzer::determineAttributeFullyQualifiedName:
// a name written with a leading "\" is kept as is, an imported one resolves
// through its import (case-sensitive), and any other is prefixed with the
// namespace ("\Foo" in the global one).
func garFullyQualifiedName(sc *symScan, it attrItem) string {
	if strings.HasPrefix(it.name, `\`) {
		return it.name
	}
	ns := attrNamespaceAt(sc, it.start)
	namespace, name := ns.name, it.name
	first, rest, _ := strings.Cut(name, `\`)
	var use *useDecl
	for _, u := range attrClassUses(sc, ns) {
		if u.shortName == first {
			use = &u
		}
	}
	if use != nil {
		if name == first {
			return use.fullName
		}
		namespace, name = use.fullName, rest
	}
	return namespace + `\` + name
}
