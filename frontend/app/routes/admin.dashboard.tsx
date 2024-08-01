import type { LoaderFunctionArgs, MetaFunction } from "@remix-run/node";
import { Link } from "@remix-run/react";
import { ensureAdminUserLoggedIn } from "../.server/auth";

export const meta: MetaFunction = () => {
	return [{ title: "[Admin] Dashboard | iOSDC Japan 2024 Albatross.swift" }];
};

export async function loader({ request }: LoaderFunctionArgs) {
	await ensureAdminUserLoggedIn(request);
	return null;
}

export default function AdminDashboard() {
	return (
		<div>
			<h1>[Admin] Dashboard</h1>
			<p>
				<Link to="/admin/users">Users</Link>
			</p>
			<p>
				<Link to="/admin/games">Games</Link>
			</p>
		</div>
	);
}
