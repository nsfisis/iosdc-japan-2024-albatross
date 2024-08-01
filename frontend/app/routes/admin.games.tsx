import type { LoaderFunctionArgs, MetaFunction } from "@remix-run/node";
import { Link, useLoaderData } from "@remix-run/react";
import { adminApiGetGames } from "../.server/api/client";
import { ensureAdminUserLoggedIn } from "../.server/auth";

export const meta: MetaFunction = () => {
	return [{ title: "[Admin] Games | iOSDC Japan 2024 Albatross.swift" }];
};

export async function loader({ request }: LoaderFunctionArgs) {
	const { token } = await ensureAdminUserLoggedIn(request);
	const { games } = await adminApiGetGames(token);
	return { games };
}

export default function AdminGames() {
	const { games } = useLoaderData<typeof loader>()!;

	return (
		<div>
			<div>
				<h1>[Admin] Games</h1>
				<ul>
					{games.map((game) => (
						<li key={game.game_id}>
							<Link to={`/admin/games/${game.game_id}`}>
								{game.display_name} (id={game.game_id})
							</Link>
						</li>
					))}
				</ul>
			</div>
		</div>
	);
}
