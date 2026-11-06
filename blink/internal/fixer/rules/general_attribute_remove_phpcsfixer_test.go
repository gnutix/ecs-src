package rules

import "testing"

// Cases ported from PHP-CS-Fixer's GeneralAttributeRemoveFixerTest
// (tests/Fixer/AttributeNotation); every expected output is what PHP-CS-Fixer
// produces for the input.
func TestGeneralAttributeRemovePhpCsFixerCases(t *testing.T) {
	t.Parallel()
	aliasRunUpstreamCases(t, GeneralAttributeRemove{}, []aliasUpstreamCase{
		{"Explicit in namespace", map[string]any{"attributes": []any{`A\B\Bar`, `Test\A\B\Quux`}},
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

            #[Foo(4, 'baz qux')]
            #[AB\Baz(prop: 'baz')]
            #[\A\B\Qux()]
            #[Corge]
            function f() {}
            `},
		{"Explicit in global namespace", map[string]any{"attributes": []any{`\A\B\Qux`, `\Corge`, `A\B\Bar`}},
			`<?php
            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
            #[BarAlias(3)]
            #[AB\Baz(prop: 'baz')]
            #[Foo(4, 'baz qux')]
            #[\A\B\Qux()]
            #[Corge]
            function f() {}
            `,
			`<?php
            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
            #[AB\Baz(prop: 'baz')]
            #[Foo(4, 'baz qux')]
            function f() {}
            `},
		{"Multiple namespaces", map[string]any{"attributes": []any{`A\B\Bar`, `Test\AB\Baz`, `A\B\Quux`, `A\B\Baz`, `A\B\Foo`, `\AB\Baz`}},
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

                function f() {}
            }

            namespace Test2
            {
                use A\B\Bar as BarAlias;
                use A\B as AB;

                #[\A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar')]
                function f() {}
            }

            namespace Test
            {
                use A\B\Foo;
                use A\B\Bar as BarAlias2;

                function f2() {}
            }

            namespace
            {
                use A\B\Foo;
                use A\B\Bar as BarAlias3;

                function f() {}
            }
            `},
		{"With whitespaces", map[string]any{"attributes": []any{`A\B\Foo`, `Test\A\B\Quux`, `A\B\Baz`}},
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
            `,
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[   \A\B\Qux()   ]
            #[   BarAlias   (3)   ]
            #[   Corge   ]
            function f() {}
            `},
		{"With docblock", map[string]any{"attributes": []any{`Test\Corge`, `\A\B\Qux`}},
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
            #[BarAlias(3)]
            #[Foo(4, 'baz qux')]
            /**
             * End docblock
             */
            class X
            {}

            function f2(/** Start docblock */#[Foo(4, 'baz qux')] #[BarAlias(3)] /** End docblock */string $param) {}
            `},
		{"With comments", map[string]any{"attributes": []any{`Test\A\B\Quux`, `A\B\Bar`, `Test\Corge`, `A\B\Baz`, `A\B\Foo`, `\A\B\Qux`}},
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
            `,
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            function f() {}
            `},
		{"With multiple attributes", map[string]any{"attributes": []any{`A\B\Foo`, `\A\B\Qux`}},
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
            `,
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            #[
                AB\Baz(prop: 'baz'),
                A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar'),
                BarAlias(3),
                Corge
            ]
            class X
            {
                #[ AB\Baz(prop: 'baz'), A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar'), BarAlias(3), Corge]
                public function y() {}
            }
            `},
		{"Multiline with no trailing comma", map[string]any{"attributes": []any{`A\B\Foo`, `\A\B\Qux`, `A\B\Baz`, `Test\A\B\Quux`, `A\B\Bar`, `Test\Corge`}},
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
            `,
			`<?php
            namespace Test;

            use A\B\Foo;
            use A\B\Bar as BarAlias;
            use A\B as AB;

            class X
            {}
            `},
		{"Multiple with comments", map[string]any{"attributes": []any{`A\B\Foo`, `Test\Corge`}},
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
                BarAlias(3)
            ]
            class X
            {
                #[ /* AB\Baz comment */AB\Baz(prop: 'baz'), A\B\Quux(prop1: [1, 2, 4], prop2: true, prop3: 'foo bar'), \A\B\Qux(), BarAlias(3)]
                public function y() {}
            }
            `},
	})
}
