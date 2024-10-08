import { vitePlugin as remix } from "@remix-run/dev";
import { defineConfig } from "vite";
import tsconfigPaths from "vite-tsconfig-paths";

export default defineConfig({
	base: "/iosdc-japan/2024/code-battle/",
	plugins: [
		remix({
			future: {
				v3_fetcherPersist: true,
				v3_relativeSplatPath: true,
				v3_throwAbortReason: true,
			},
			basename: "/iosdc-japan/2024/code-battle/",
		}),
		tsconfigPaths(),
	],
});
