import type { Session } from "@remix-run/server-runtime";
import { jwtDecode } from "jwt-decode";
import { Authenticator } from "remix-auth";
import { FormStrategy } from "remix-auth-form";
import { apiPostLogin } from "./api/client";
import { components } from "./api/schema";
import { sessionStorage } from "./session";

const authenticator = new Authenticator<string>(sessionStorage);

authenticator.use(
	new FormStrategy(async ({ form }) => {
		const username = String(form.get("username"));
		const password = String(form.get("password"));
		return (await apiPostLogin(username, password)).token;
	}),
	"default",
);

export type User = components["schemas"]["User"];

export async function login(request: Request): Promise<never> {
	return await authenticator.authenticate("default", request, {
		successRedirect: "/dashboard",
		failureRedirect: "/login",
	});
}

export async function logout(request: Request | Session): Promise<never> {
	return await authenticator.logout(request, { redirectTo: "/" });
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
