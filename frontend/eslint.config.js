import js from "@eslint/js";
import reactHooks from "eslint-plugin-react-hooks";
import reactRefresh from "eslint-plugin-react-refresh";
import globals from "globals";
import tseslint from "typescript-eslint";

export default tseslint.config(
	{ ignores: ["dist", "src/routeTree.gen.ts"] },
	{
		extends: [js.configs.recommended, ...tseslint.configs.recommended],
		files: ["src/**/*.{ts,tsx}"],
		languageOptions: {
			globals: globals.browser,
		},
		plugins: {
			"react-hooks": reactHooks,
			"react-refresh": reactRefresh,
		},
		rules: {
			...reactHooks.configs.recommended.rules,
			"react-refresh/only-export-components": [
				"warn",
				{ allowConstantExport: true },
			],
		},
	},
	{
		// TanStack Router route files and shadcn/ui components export
		// non-component values (Route, variants helpers) alongside the
		// component by design — disable the rule for them.
		files: ["src/routes/**/*.{ts,tsx}", "src/components/ui/**/*.{ts,tsx}"],
		rules: {
			"react-refresh/only-export-components": "off",
		},
	},
);
