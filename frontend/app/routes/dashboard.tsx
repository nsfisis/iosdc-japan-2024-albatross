import type { LoaderFunctionArgs, MetaFunction } from "@remix-run/node";
import { redirect } from "@remix-run/node";
import { Form, Link, useLoaderData } from "@remix-run/react";
import { apiGetGames } from "../.server/api/client";
import { ensureUserLoggedIn } from "../.server/auth";

export const meta: MetaFunction = () => [
	{ title: "Dashboard | iOSDC Japan 2024 Albatross.swift" },
];

export async function loader({ request }: LoaderFunctionArgs) {
	const { user, token } = await ensureUserLoggedIn(request);
	if (user.is_admin) {
		return redirect(
			process.env.NODE_ENV === "development"
				? "http://localhost:8002/admin/dashboard"
				: "/admin/dashboard",
		);
	}
	const { games } = await apiGetGames(token);
	return {
		user,
		games,
	};
}

export default function Dashboard() {
	const { user, games } = useLoaderData<typeof loader>()!;

	return (
		<div className="min-h-screen p-8">
			<div className="p-6 rounded shadow-md max-w-4xl mx-auto">
				<h1 className="text-3xl font-bold mb-4">{user.username}</h1>
				<h2 className="text-2xl font-semibold mb-2">User</h2>
				<div className="mb-6">
					<ul className="list-disc list-inside">
						<li>Name: {user.display_name}</li>
					</ul>
				</div>
				<h2 className="text-2xl font-semibold mb-2">Games</h2>
				<div>
					<ul className="list-disc list-inside">
						{games.map((game) => (
							<li key={game.game_id}>
								{game.display_name}{" "}
								{game.state === "closed" || game.state === "finished" ? (
									<span className="inline-block px-6 py-2 text-gray-400 bg-gray-200 cursor-not-allowed rounded">
										Entry
									</span>
								) : (
									<Link
										to={`/golf/${game.game_id}/play`}
										className="inline-block px-6 py-2 text-white bg-blue-500 hover:bg-blue-700 rounded"
									>
										Entry
									</Link>
								)}
							</li>
						))}
					</ul>
				</div>
				<div>
					<Form method="post" action="/logout">
						<button
							className="mt-6 px-6 py-2 text-white bg-red-500 hover:bg-red-700 rounded"
							type="submit"
						>
							Logout
						</button>
					</Form>
				</div>
			</div>
		</div>
	);
}
