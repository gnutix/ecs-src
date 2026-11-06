package rules

import "testing"

// Cases ported from PHP-CS-Fixer's OrderedAttributesFixerTest
// (tests/Fixer/AttributeNotation); every expected output is what PHP-CS-Fixer
// produces for the input.
func TestOrderedAttributesPhpCsFixerCases(t *testing.T) {
	t.Parallel()
	aliasRunUpstreamCases(t, OrderedAttributes{}, []aliasUpstreamCase{
		{"Alpha on various declarations", nil,
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[Foo(4, 'baz qux')]
            #[BarAlias(3)]
            #[AB\Baz(prop: 'baz')]
            #[\A\B\Qux()]
            #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
            #[Corge]
            class X
            {
                #[Foo(4, 'baz qux')]
                #[BarAlias(3)]
                #[AB\Baz(prop: 'baz')]
                #[\A\B\Qux()]
                #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
                #[Corge]
                const Y = 1;

                #[Foo(4, 'baz qux')]
                #[BarAlias(3)]
                #[AB\Baz(prop: 'baz')]
                #[\A\B\Qux()]
                #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
                #[Corge]
                public $y;

                #[Foo(4, 'baz qux')]
                #[BarAlias(3)]
                #[AB\Baz(prop: 'baz')]
                #[\A\B\Qux()]
                #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
                #[Corge]
                public function y() {}
            }

            #[Foo(4, 'baz qux')]
            #[BarAlias(3)]
            #[AB\Baz(prop: 'baz')]
            #[\A\B\Qux()]
            #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
            #[Corge]
            function f(
                #[Foo(4, 'baz qux')]
                #[BarAlias(3)]
                #[AB\Baz(prop: 'baz')]
                string $param1,
                #[\A\B\Qux()]
                #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
                #[Corge]
                string $param2
            ) {}

            function f2(#[Foo(4, 'baz qux')] #[BarAlias(3)] string $param1, #[\A\B\Qux()] #[Corge] string $param2) {}

            $anon = #[Foo(4, 'baz qux')] #[BarAlias(3)] #[AB\Baz(prop: 'baz')] #[\A\B\Qux()] #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')] #[Corge] function () {};
            $short = #[Foo(4, 'baz qux')] #[BarAlias(3)] #[AB\Baz(prop: 'baz')] #[\A\B\Qux()] #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')] #[Corge] fn () => null;
            `,
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[AB\Baz(prop: 'baz')]
            #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
            #[\A\B\Qux()]
            #[BarAlias(3)]
            #[Corge]
            #[Foo(4, 'baz qux')]
            class X
            {
                #[AB\Baz(prop: 'baz')]
                #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
                #[\A\B\Qux()]
                #[BarAlias(3)]
                #[Corge]
                #[Foo(4, 'baz qux')]
                const Y = 1;

                #[AB\Baz(prop: 'baz')]
                #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
                #[\A\B\Qux()]
                #[BarAlias(3)]
                #[Corge]
                #[Foo(4, 'baz qux')]
                public $y;

                #[AB\Baz(prop: 'baz')]
                #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
                #[\A\B\Qux()]
                #[BarAlias(3)]
                #[Corge]
                #[Foo(4, 'baz qux')]
                public function y() {}
            }

            #[AB\Baz(prop: 'baz')]
            #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
            #[\A\B\Qux()]
            #[BarAlias(3)]
            #[Corge]
            #[Foo(4, 'baz qux')]
            function f(
                #[AB\Baz(prop: 'baz')]
                #[BarAlias(3)]
                #[Foo(4, 'baz qux')]
                string $param1,
                #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
                #[\A\B\Qux()]
                #[Corge]
                string $param2
            ) {}

            function f2(#[BarAlias(3)] #[Foo(4, 'baz qux')] string $param1, #[\A\B\Qux()] #[Corge] string $param2) {}

            $anon = #[AB\Baz(prop: 'baz')] #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')] #[\A\B\Qux()] #[BarAlias(3)] #[Corge] #[Foo(4, 'baz qux')] function () {};
            $short = #[AB\Baz(prop: 'baz')] #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')] #[\A\B\Qux()] #[BarAlias(3)] #[Corge] #[Foo(4, 'baz qux')] fn () => null;
            `},
		{"Explicit in namespace", map[string]any{"sort_algorithm": `custom`, "order": []any{`Test\A\B\Quux`, `A\B\Bar`, `Test\Corge`, `A\B\Baz`, `A\B\Foo`, `A\B\Qux`}},
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[Foo(4, 'baz qux')]
            #[BarAlias(3)]
            #[AB\Baz(prop: 'baz')]
            #[\A\B\Qux()]
            #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
            #[Corge]
            function f() {}
            `,
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
            #[BarAlias(3)]
            #[Corge]
            #[AB\Baz(prop: 'baz')]
            #[Foo(4, 'baz qux')]
            #[\A\B\Qux()]
            function f() {}
            `},
		{"Explicit in global namespace", map[string]any{"sort_algorithm": `custom`, "order": []any{`A\B\Quux`, `A\B\Bar`, `Corge`, `A\B\Baz`, `A\B\Foo`, `A\B\Qux`}},
			`<?php
            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[Foo(4, 'baz qux')]
            #[BarAlias(3)]
            #[AB\Baz(prop: 'baz')]
            #[\A\B\Qux()]
            #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
            #[Corge]
            function f() {}
            `,
			`<?php
            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
            #[BarAlias(3)]
            #[Corge]
            #[AB\Baz(prop: 'baz')]
            #[Foo(4, 'baz qux')]
            #[\A\B\Qux()]
            function f() {}
            `},
		{"Explicit with unconfigured attributes", map[string]any{"sort_algorithm": `custom`, "order": []any{`Test\A\B\Quux`, `A\B\Baz`, `A\B\Qux`}},
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[Foo(4, 'baz qux')]
            #[BarAlias(3)]
            #[AB\Baz(prop: 'baz')]
            #[\A\B\Qux()]
            #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
            #[Corge]
            function f() {}
            `,
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
            #[AB\Baz(prop: 'baz')]
            #[\A\B\Qux()]
            #[Foo(4, 'baz qux')]
            #[BarAlias(3)]
            #[Corge]
            function f() {}
            `},
		{"Multiple namespaces", map[string]any{"sort_algorithm": `custom`, "order": []any{`A\B\Bar`, `Test\AB\Baz`, `A\B\Quux`, `A\B\Baz`, `A\B\Foo`, `AB\Baz`}},
			`<?php
            namespace Test
            {
                use A\B\Foo;
                use A\B\Bar as BarAlias;

                #[AB\Baz(prop: 'baz')]
                #[Foo(4, 'baz qux')]
                #[BarAlias(3)]
                function f() {}
            }

            namespace Test2
            {
                use A\B\Bar as BarAlias;
                use A\B as AB;

                #[BarAlias(3)]
                #[AB\Baz(prop: 'baz')]
                #[\A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
                function f() {}
            }

            namespace Test
            {
                use A\B\Foo;
                use A\B\Bar as BarAlias2;

                #[AB\Baz(prop: 'baz')]
                #[Foo(4, 'baz qux')]
                #[BarAlias2(3)]
                function f2() {}
            }

            namespace
            {
                use A\B\Foo;
                use A\B\Bar as BarAlias3;

                #[AB\Baz(prop: 'baz')]
                #[Foo(4, 'baz qux')]
                #[BarAlias3(3)]
                function f() {}
            }
            `,
			`<?php
            namespace Test
            {
                use A\B\Foo;
                use A\B\Bar as BarAlias;

                #[BarAlias(3)]
                #[AB\Baz(prop: 'baz')]
                #[Foo(4, 'baz qux')]
                function f() {}
            }

            namespace Test2
            {
                use A\B\Bar as BarAlias;
                use A\B as AB;

                #[BarAlias(3)]
                #[\A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
                #[AB\Baz(prop: 'baz')]
                function f() {}
            }

            namespace Test
            {
                use A\B\Foo;
                use A\B\Bar as BarAlias2;

                #[BarAlias2(3)]
                #[AB\Baz(prop: 'baz')]
                #[Foo(4, 'baz qux')]
                function f2() {}
            }

            namespace
            {
                use A\B\Foo;
                use A\B\Bar as BarAlias3;

                #[BarAlias3(3)]
                #[Foo(4, 'baz qux')]
                #[AB\Baz(prop: 'baz')]
                function f() {}
            }
            `},
		{"With whitespaces", nil,
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[   Foo   (4, 'baz qux')   ]
            #[   BarAlias   (3)   ]
            #[   AB\Baz   (prop: 'baz')   ]
            #[   \A\B\Qux()   ]
            #[   A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')   ]
            #[   Corge   ]
            function f() {}
            `,
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[   AB\Baz   (prop: 'baz')   ]
            #[   A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')   ]
            #[   \A\B\Qux()   ]
            #[   BarAlias   (3)   ]
            #[   Corge   ]
            #[   Foo   (4, 'baz qux')   ]
            function f() {}
            `},
		{"With docblock", nil,
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            /**
             * Start docblock
             */
            #[Foo(4, 'baz qux')]
            #[BarAlias(3)]
            /**
             * AB\Baz docblock
             */
            #[AB\Baz(prop: 'baz')]
            #[\A\B\Qux()]
            #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
            /**
             * Corge docblock
             */
            #[Corge]
            /**
             * End docblock
             */
            class X
            {}

            function f2(/** Start docblock */#[Foo(4, 'baz qux')] #[BarAlias(3)] #[\A\B\Qux()] /** Corge docblock */#[Corge] /** End docblock */string $param) {}
            `,
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            /**
             * Start docblock
             */
            /**
             * AB\Baz docblock
             */
            #[AB\Baz(prop: 'baz')]
            #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
            #[\A\B\Qux()]
            #[BarAlias(3)]
            /**
             * Corge docblock
             */
            #[Corge]
            #[Foo(4, 'baz qux')]
            /**
             * End docblock
             */
            class X
            {}

            function f2(/** Start docblock */#[\A\B\Qux()] #[BarAlias(3)] /** Corge docblock */#[Corge] #[Foo(4, 'baz qux')] /** End docblock */string $param) {}
            `},
		{"With comments", map[string]any{"sort_algorithm": `custom`, "order": []any{`Test\A\B\Quux`, `A\B\Bar`, `Test\Corge`, `A\B\Baz`, `A\B\Foo`, `A\B\Qux`}},
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[/* comment */Foo/* comment */(4, 'baz qux') /* comment */ ]
            #[ /* comment */ BarAlias/* comment */(3)/* comment */]
            #[/* comment */AB\Baz /* comment */ (prop: 'baz') /* comment */ ]
            #[   /* comment */   \A\B\Qux()/* comment */]
            #[/* comment */A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar') /* comment */]
            #[/* comment */ Corge/* comment */]
            function f() {}
            `,
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[/* comment */A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar') /* comment */]
            #[ /* comment */ BarAlias/* comment */(3)/* comment */]
            #[/* comment */ Corge/* comment */]
            #[/* comment */AB\Baz /* comment */ (prop: 'baz') /* comment */ ]
            #[/* comment */Foo/* comment */(4, 'baz qux') /* comment */ ]
            #[   /* comment */   \A\B\Qux()/* comment */]
            function f() {}
            `},
		{"Alpha with multiple attributes", nil,
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[
                Foo(4, 'baz qux'),
                BarAlias(3),
                AB\Baz(prop: 'baz'),
                \A\B\Qux(),
                A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar'),
                Corge,
            ]
            class X
            {
                #[Foo(4, 'baz qux'), BarAlias(3), AB\Baz(prop: 'baz'), \A\B\Qux(), A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar'), Corge]
                public function y() {}
            }
            `,
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[
                AB\Baz(prop: 'baz'),
                A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar'),
                \A\B\Qux(),
                BarAlias(3),
                Corge,
                Foo(4, 'baz qux'),
            ]
            class X
            {
                #[ AB\Baz(prop: 'baz'), A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar'), \A\B\Qux(), BarAlias(3), Corge,Foo(4, 'baz qux')]
                public function y() {}
            }
            `},
		{"Explicit with multiple attributes", map[string]any{"sort_algorithm": `custom`, "order": []any{`Test\A\B\Quux`, `A\B\Bar`, `Test\Corge`, `A\B\Baz`, `A\B\Qux`}},
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[
                Foo(4, 'baz qux'),
                BarAlias(3),
                AB\Baz(prop: 'baz'),
                \A\B\Qux(),
                A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar'),
                Corge,
            ]
            class X
            {
                #[Foo(4, 'baz qux'), BarAlias(3), AB\Baz(prop: 'baz'), \A\B\Qux(), A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar'), Corge]
                public function y() {}
            }
            `,
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[
                A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar'),
                BarAlias(3),
                Corge,
                AB\Baz(prop: 'baz'),
                \A\B\Qux(),
                Foo(4, 'baz qux'),
            ]
            class X
            {
                #[ A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar'), BarAlias(3), Corge, AB\Baz(prop: 'baz'), \A\B\Qux(),Foo(4, 'baz qux')]
                public function y() {}
            }
            `},
		{"Multiline with no trailing comma", nil,
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[
                Foo(4, 'baz qux'),
                BarAlias(3),
                AB\Baz(prop: 'baz'),
                \A\B\Qux(),
                A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar'),
                Corge
            ]
            class X
            {}
            `,
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[
                AB\Baz(prop: 'baz'),
                A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar'),
                \A\B\Qux(),
                BarAlias(3),
                Corge,
                Foo(4, 'baz qux')
            ]
            class X
            {}
            `},
		{"Multiple with comments", nil,
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[
                /**
                 * Foo docblock
                 */
                Foo(4, 'baz qux'),
                BarAlias(3),
                /*
                 * AB\Baz comment
                 */
                AB\Baz(prop: 'baz'),
                \A\B\Qux(),
                A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar'),
                /*
                 * Corge comment
                 */
                Corge,
            ]
            class X
            {
                #[/** Foo docblock */Foo(4, 'baz qux'), BarAlias(3), /* AB\Baz comment */AB\Baz(prop: 'baz'), \A\B\Qux(), A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar'), /* Corge comment */Corge]
                public function y() {}
            }
            `,
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[
                /*
                 * AB\Baz comment
                 */
                AB\Baz(prop: 'baz'),
                A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar'),
                \A\B\Qux(),
                BarAlias(3),
                /*
                 * Corge comment
                 */
                Corge,
                /**
                 * Foo docblock
                 */
                Foo(4, 'baz qux'),
            ]
            class X
            {
                #[ /* AB\Baz comment */AB\Baz(prop: 'baz'), A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar'), \A\B\Qux(), BarAlias(3), /* Corge comment */Corge,/** Foo docblock */Foo(4, 'baz qux')]
                public function y() {}
            }
            `},
		{"Alpha with mixed multiple attribute declarations and attributes", nil,
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[Foo(4, 'baz qux'), \A\B\Qux(), BarAlias(3)]
            #[AB\Baz(prop: 'baz')]
            #[Corge, A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
            function f() {}
            `,
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[AB\Baz(prop: 'baz')]
            #[ A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar'),Corge]
            #[ \A\B\Qux(), BarAlias(3),Foo(4, 'baz qux')]
            function f() {}
            `},
		{"Explicit with mixed multiple attribute declarations and attributes", map[string]any{"sort_algorithm": `custom`, "order": []any{`Test\A\B\Quux`, `A\B\Bar`, `Test\Corge`, `A\B\Baz`, `A\B\Qux`}},
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[Foo(4, 'baz qux'), \A\B\Qux(), BarAlias(3)]
            #[AB\Baz(prop: 'baz')]
            #[Corge, A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
            function f() {}
            `,
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[ A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar'),Corge]
            #[ BarAlias(3), \A\B\Qux(),Foo(4, 'baz qux')]
            #[AB\Baz(prop: 'baz')]
            function f() {}
            `},
		{"Alpha with nested attribute", nil,
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[Foo(4, 'baz qux')]
            #[BarAlias(3)]
            #[AB\Baz(prop: new P\R())]
            #[\A\B\Qux()]
            #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
            #[Corge]
            function f() {}
            `,
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[AB\Baz(prop: new P\R())]
            #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
            #[\A\B\Qux()]
            #[BarAlias(3)]
            #[Corge]
            #[Foo(4, 'baz qux')]
            function f() {}
            `},
		{"Explicit multiple with nested attribute", map[string]any{"sort_algorithm": `custom`, "order": []any{`Test\A\B\Quux`, `A\B\Bar`, `Test\Corge`, `A\B\Baz`, `A\B\Foo`, `A\B\Qux`}},
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[
                Foo(4, 'baz qux'),
                BarAlias(3),
                AB\Baz(prop: new P\R()),
                \A\B\Qux(),
                A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar'),
                Corge,
            ]
            function f() {}
            `,
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[
                A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar'),
                BarAlias(3),
                Corge,
                AB\Baz(prop: new P\R()),
                Foo(4, 'baz qux'),
                \A\B\Qux(),
            ]
            function f() {}
            `},
	})
}
