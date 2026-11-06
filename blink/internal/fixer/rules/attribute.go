package rules

import (
	"strings"

	"github.com/ecsphp/ecs-src/blink/internal/lexer"
	"github.com/ecsphp/ecs-src/blink/internal/token"
	"github.com/ecsphp/ecs-src/blink/internal/tokens"
)

// Attributes lex natively ("#[" is an AttributeStart token, its contents and the
// closing "]" are ordinary tokens), so most fixers see inside them for free. The
// helpers here cover the few rules that still reason about an attribute as a whole.

// attributeStartsAt reports whether an attribute ("#[") begins at i.
func attributeStartsAt(s *tokens.Stream, i int) bool {
	return i >= 0 && i < s.Len() && s.At(i).Kind == token.AttributeStart
}

// attributeEndsAt reports whether the attribute-closing "]" is at i.
func attributeEndsAt(s *tokens.Stream, i int) bool {
	return s.AttributeStartOf(i) >= 0
}

// prevMeaningfulSkippingAttributes is prevMeaningfulIndex but also steps back over
// any attributes ("new #[A] #[B] class" -> the "new"), since an attribute decorates
// the token that follows it.
func prevMeaningfulSkippingAttributes(s *tokens.Stream, i int) int {
	p := prevMeaningfulIndex(s, i)
	for p >= 0 && attributeEndsAt(s, p) {
		p = prevMeaningfulIndex(s, s.AttributeStartOf(p))
	}
	return p
}

// attributeText renders the attribute tokens [start..end] back to source text.
func attributeText(s *tokens.Stream, start, end int) string {
	var b strings.Builder
	for j := start; j <= end; j++ {
		b.WriteString(s.At(j).Value)
	}
	return b.String()
}

// rewriteAttributeText feeds the attribute [start..end] to transform as its source
// text ("#[...]"), re-lexes the result and splices it back in place. It bridges the
// text-shaped attribute transforms to the native token stream. Returns whether the
// stream changed.
func rewriteAttributeText(s *tokens.Stream, start, end int, transform func(string) (string, bool)) bool {
	old := attributeText(s, start, end)
	nv, ok := transform(old)
	if !ok || nv == old {
		return false
	}
	s.ReplaceRange(start, end, lexAttributeFragment(nv))
	return true
}

// lexAttributeFragment lexes a standalone "#[...]" fragment into its tokens, from
// the AttributeStart through the end (the opening PHP tag used to drive the lexer
// is dropped).
func lexAttributeFragment(text string) []token.Token {
	toks := lexer.Lex("<?php " + text)
	for i := range toks {
		if toks[i].Kind == token.AttributeStart {
			return toks[i:]
		}
	}
	return toks
}

// attrItem is one attribute of a "#[...]" block: its name as written and its
// tokens, from just after the "#[" or "," to just before the next "," or "]".
type attrItem struct {
	name       string
	start, end int
}

// attrBlock mirrors PHP-CS-Fixer's AttributeAnalysis. start and end take in the
// comments before the block when it follows another one, and the whitespace
// after its "]", so the blocks of a group cover it with no gap.
type attrBlock struct {
	start, end  int
	open, close int
	items       []attrItem
}

// attrCollectGroup is AttributeAnalyzer::collect: the block at open and the
// blocks that follow it with only whitespace and comments in between.
func attrCollectGroup(s *tokens.Stream, open int) ([]attrBlock, bool) {
	var group []attrBlock
	for {
		b, ok := attrCollectOne(s, open)
		if !ok {
			return nil, false
		}
		group = append(group, b)
		open = nextMeaningfulIndex(s, b.end)
		if !attributeStartsAt(s, open) {
			return group, true
		}
	}
}

// attrCollectOne is AttributeAnalyzer::collectOne for the "#[" at open.
func attrCollectOne(s *tokens.Stream, open int) (attrBlock, bool) {
	b := attrBlock{start: open, open: open, close: s.AttributeEnd(open)}
	if b.close < 0 {
		return b, false
	}
	if prev := prevMeaningfulIndex(s, open); prev >= 0 && s.At(prev).Kind == token.AttributeEnd {
		b.start = skipWhitespace(s, prev+1)
	}
	b.end = skipWhitespace(s, b.close+1) - 1
	b.items = attrCollectItems(s, open, b.close)
	return b, b.items != nil
}

// attrCollectItems is AttributeAnalyzer::collectAttributes.
func attrCollectItems(s *tokens.Stream, index, closing int) []attrItem {
	var items []attrItem
	for {
		start := index + 1
		nameStart := attrNameStart(s, index)
		index = attrNextParenCommaOrClose(s, start)
		if nameStart < 0 || index < 0 {
			return nil
		}
		name := attributeText(s, nameStart, prevMeaningfulIndex(s, index))
		if isPunctVal(s, index, "(") {
			if index = s.MatchForward(index); index < 0 {
				return nil
			}
			if index = attrNextCommaOrClose(s, index); index < 0 {
				return nil
			}
		}
		items = append(items, attrItem{name: name, start: start, end: index - 1})
		next := index
		if next < closing {
			next = nextMeaningfulIndex(s, index)
		}
		if next < 0 || next >= closing {
			break
		}
	}
	// without a trailing comma, the last attribute ends before the line break
	for j := index - 1; s.At(j).Kind == token.Whitespace; j-- {
		if strings.ContainsAny(s.At(j).Value, "\r\n") {
			items[len(items)-1].end = j - 1
			break
		}
	}
	return items
}

// attrNameStart returns the first name token (T_STRING or T_NS_SEPARATOR in
// PHP-CS-Fixer) after i, or -1.
func attrNameStart(s *tokens.Stream, i int) int {
	for j := i + 1; j < s.Len(); j++ {
		if isPunctVal(s, j, `\`) {
			return j
		}
		// the "namespace" of "namespace\Foo" is T_NAMESPACE, not a name
		if t := s.At(j); t.Kind == token.Ident && (!strings.EqualFold(t.Value, "namespace") || !isPunctVal(s, j+1, `\`)) {
			return j
		}
	}
	return -1
}

// attrNamespaceAt returns the namespace declaration whose scope holds i, like
// NamespacesAnalyzer::getNamespaceAt.
func attrNamespaceAt(sc *symScan, i int) nsDecl {
	for _, ns := range sc.nsDecls {
		if ns.scopeStart <= i && i <= ns.scopeEnd {
			return ns
		}
	}
	return nsDecl{}
}

// attrClassUses returns the class imports declared in the scope of ns.
func attrClassUses(sc *symScan, ns nsDecl) []useDecl {
	var out []useDecl
	for _, u := range sc.uses {
		if u.kind == "class" && ns.scopeStart <= u.start && u.start <= ns.scopeEnd {
			out = append(out, u)
		}
	}
	return out
}
