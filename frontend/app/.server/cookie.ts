import { Cookie, CookieOptions } from "@remix-run/server-runtime";
import { parse, serialize } from "cookie";

// Remix's createCookie() returns "structured" cookies, which are cookies that hold a JSON-encoded object.
// This is not suitable for interoperation with other systems that expect a simple string value.
// This function creates an "unstructured" cookie, a simple plain text.
export function createUnstructuredCookie(
	name: string,
	cookieOptions?: CookieOptions,
): Cookie {
	const { secrets = [], ...options } = {
		path: "/",
		sameSite: "lax" as const,
		...cookieOptions,
	};

	return {
		get name() {
			return name;
		},
		get isSigned() {
			return secrets.length > 0;
		},
		get expires() {
			return typeof options.maxAge !== "undefined"
				? new Date(Date.now() + options.maxAge * 1000)
				: options.expires;
		},
		async parse(cookieHeader, parseOptions) {
			if (!cookieHeader) return null;
			const cookies = parse(cookieHeader, { ...options, ...parseOptions });
			return name in cookies ? cookies[name] : null;
		},
		async serialize(value, serializeOptions) {
			return serialize(name, value, {
				...options,
				...serializeOptions,
			});
		},
	};
}
