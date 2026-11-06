package rules

import (
	"cmp"
	"math"
	"slices"
	"strings"

	"github.com/ecsphp/ecs-src/blink/internal/fixer"
	"github.com/ecsphp/ecs-src/blink/internal/token"
	"github.com/ecsphp/ecs-src/blink/internal/tokens"
)

// PHP-CS-Fixer: https://github.com/PHP-CS-Fixer/PHP-CS-Fixer/blob/master/src/Fixer/AttributeNotation/OrderedAttributesFixer.php
//
// OrderedAttributes sorts the attributes inside each "#[...]" block and the
// blocks of each group, by name ("alpha", the default) or by the position of
// their fully qualified name in the "order" list ("custom").
type OrderedAttributes struct {
	custom bool
	order  map[string]int
}

func (OrderedAttributes) Name() string {
	return `PhpCsFixer\Fixer\AttributeNotation\OrderedAttributesFixer`
}

func (OrderedAttributes) SourceURL() string {
	return "https://github.com/PHP-CS-Fixer/PHP-CS-Fixer/blob/master/src/Fixer/AttributeNotation/OrderedAttributesFixer.php"
}

func (f OrderedAttributes) WithConfig(config map[string]any) fixer.Fixer {
	if v, ok := config["sort_algorithm"].(string); ok {
		f.custom = v == "custom"
	}
	if v, ok := config["order"].([]any); ok {
		f.order = make(map[string]int, len(v))
		for i, name := range v {
			if name, ok := name.(string); ok {
				f.order[strings.TrimLeft(name, `\`)] = i
			}
		}
	}
	return f
}

func (OrderedAttributes) IsCandidate(s *tokens.Stream) bool {
	return s.Index().HasKind(token.AttributeStart)
}

func (f OrderedAttributes) Fix(s *tokens.Stream) bool {
	version := s.Version()
	var sc *symScan
	if f.custom {
		sc = sharedSymScan(s)
	}
	for index := attrNextStart(s, 0); index >= 0; {
		group, ok := attrCollectGroup(s, index)
		if !ok {
			break
		}
		blocks := make([]attrItem, len(group))
		for i, b := range group {
			blocks[i] = attrItem{name: f.sortBlock(s, sc, b), start: b.start, end: b.end}
		}
		endIndex := blocks[len(blocks)-1].end
		if len(blocks) > 1 {
			if sorted := f.sortItems(blocks); !slices.Equal(sorted, blocks) {
				oaSortTokens(s, index, endIndex, sorted, false)
			}
		}
		index = attrNextStart(s, endIndex)
	}
	return s.Version() != version
}

// sortBlock sorts the attributes of one block and returns the name of its
// first one.
func (f OrderedAttributes) sortBlock(s *tokens.Stream, sc *symScan, b attrBlock) string {
	items := b.items
	for i := range items {
		items[i].name = f.attributeName(sc, items[i])
	}
	sorted := f.sortItems(items)
	if !slices.Equal(sorted, items) {
		// from b.start, not b.open, as upstream: after a comment, the "#[" is overwritten too
		oaSortTokens(s, b.start+1, items[len(items)-1].end, sorted, true)
	}
	return sorted[0].name
}

func (f OrderedAttributes) attributeName(sc *symScan, it attrItem) string {
	if f.custom {
		return oaFullyQualifiedName(sc, it.name, it.start)
	}
	return strings.TrimLeft(it.name, `\`)
}

func (f OrderedAttributes) sortItems(items []attrItem) []attrItem {
	sorted := slices.Clone(items)
	slices.SortStableFunc(sorted, func(a, b attrItem) int {
		if !f.custom {
			return strings.Compare(a.name, b.name)
		}
		return cmp.Compare(f.position(a.name), f.position(b.name))
	})
	return sorted
}

func (f OrderedAttributes) position(name string) int {
	if p, ok := f.order[name]; ok {
		return p
	}
	return math.MaxInt
}

// oaSortTokens is OrderedAttributesFixer::sortTokens: it writes the tokens of
// items, in order, over [start, end] and clears what is left of the range.
func oaSortTokens(s *tokens.Stream, start, end int, items []attrItem, comma bool) {
	var repl []token.Token
	for k, it := range items {
		for j := it.start; j <= it.end; j++ {
			repl = append(repl, s.At(j))
		}
		if comma && k < len(items)-1 {
			repl = append(repl, token.Token{Kind: token.Punct, Value: ","})
		}
	}
	for j := start; j <= end; j++ {
		if k := j - start; k < len(repl) {
			s.Set(j, repl[k])
		} else {
			clearTokenAt(s, j)
		}
	}
}

// oaFullyQualifiedName is FullyQualifiedNameAnalyzer::getFullyQualifiedName for
// a class name: imports match without case, and a later import of the same
// short name replaces an earlier one.
func oaFullyQualifiedName(sc *symScan, name string, index int) string {
	if strings.HasPrefix(name, `\`) {
		return strings.TrimLeft(name, `\`)
	}
	ns := attrNamespaceAt(sc, index)
	var shorts, fulls []string
	for _, u := range attrClassUses(sc, ns) {
		short := asciiLower(u.shortName)
		if k := slices.Index(shorts, short); k >= 0 {
			fulls[k] = u.fullName
		} else {
			shorts, fulls = append(shorts, short), append(fulls, u.fullName)
		}
	}
	lower := asciiLower(name)
	for k, short := range shorts {
		if lower == short {
			return strings.TrimLeft(fulls[k], `\`)
		}
		if strings.HasPrefix(lower, short+`\`) {
			return strings.TrimLeft(fulls[k]+name[len(short):], `\`)
		}
	}
	return strings.TrimLeft(ns.name+`\`+name, `\`)
}
