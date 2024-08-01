import type { LinksFunction } from "@remix-run/node";
import normalizeCss from "sakura.css/css/normalize.css?url";
import sakuraCss from "sakura.css/css/sakura.css?url";

export const links: LinksFunction = () => [
	{ rel: "stylesheet", href: normalizeCss },
	{ rel: "stylesheet", href: sakuraCss },
];
