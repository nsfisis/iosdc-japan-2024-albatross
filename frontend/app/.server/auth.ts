import { redirect } from "@remix-run/node";
import type { Session } from "@remix-run/server-runtime";
import { jwtDecode } from "jwt-decode";
import { Authenticator } from "remix-auth";
import { FormStrategy } from "remix-auth-form";
import { apiPostLogin } from "./api/client";
import { components } from "./api/schema";
import { createUnstructuredCookie } from "./cookie";
import { cookieOptions, sessionStorage } from "./session";

const authenticator = new Authenticator<string>(sessionStorage);

authenticator.use(
	new FormStrategy(async ({ form }) => {
		const username = String(form.get("username"));
		const password = String(form.get("password"));
		const registrationToken = String(form.get("registration_token"));
		return (
			await apiPostLogin(
				username,
				password,
				registrationToken === "" ? null : registrationToken,
			)
		).token;
	}),
	"default",
);

export type User = components["schemas"]["User"];

// This cookie is used to directly store the JWT for the API server.
// Remix's createCookie() returns "structured" cookies, which cannot be reused directly by non-Remix servers.
const tokenCookie = createUnstructuredCookie("albatross_token", cookieOptions);

export async function login(request: Request): Promise<never> {
	const jwt = await authenticator.authenticate("default", request, {
		failureRedirect: request.url,
	});

	const session = await sessionStorage.getSession(
		request.headers.get("cookie"),
	);
	session.set(authenticator.sessionKey, jwt);

	throw redirect("/dashboard", {
		headers: [
			["Set-Cookie", await sessionStorage.commitSession(session)],
			["Set-Cookie", await tokenCookie.serialize(jwt)],
		],
	});
}

export async function logout(request: Request | Session): Promise<never> {
	try {
		return await authenticator.logout(request, { redirectTo: "/" });
	} catch (response) {
		if (response instanceof Response) {
			response.headers.append(
				"Set-Cookie",
				await tokenCookie.serialize("", { maxAge: 0, expires: new Date(0) }),
			);
		}
		throw response;
	}
}

export async function ensureUserLoggedIn(
	request: Request | Session,
): Promise<{ user: User; token: string }> {
	const token = await authenticator.isAuthenticated(request, {
		failureRedirect: "/login",
	});
	const user = jwtDecode<User>(token);
	return { user, token };
}

export async function ensureUserNotLoggedIn(
	request: Request | Session,
): Promise<null> {
	return await authenticator.isAuthenticated(request, {
		successRedirect: "/dashboard",
	});
}
