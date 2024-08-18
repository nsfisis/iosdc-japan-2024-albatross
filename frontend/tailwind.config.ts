import type { Config } from "tailwindcss";

export default {
	content: ["./app/**/{**,.client,.server}/**/*.{js,jsx,ts,tsx}"],
	theme: {
		extend: {
			fontSize: {
				"10xl": "16rem",
			},
		},
	},
	plugins: [],
} satisfies Config;
