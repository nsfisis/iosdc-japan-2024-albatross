import { createCookieSessionStorage } from "@remix-run/node";

export const cookieOptions = {
	sameSite: "lax" as const,
	path: "/",
	httpOnly: true,
	// secure: process.env.NODE_ENV === "production",
	secure: false, // TODO
	secrets: ["TODO"],
};

export const sessionStorage = createCookieSessionStorage({
	cookie: {
		name: "albatross_session",
		...cookieOptions,
	},
});
