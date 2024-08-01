import type { ActionFunctionArgs } from "@remix-run/node";
import { logout } from "../.server/auth";

export async function action({ request }: ActionFunctionArgs) {
	await logout(request);
	return null;
}
