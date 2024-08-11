import type { LoaderFunctionArgs, MetaFunction } from "@remix-run/node";
import { redirect } from "@remix-run/node";
import { Form, useLoaderData } from "@remix-run/react";
import { apiGetGames } from "../.server/api/client";
import { ensureUserLoggedIn } from "../.server/auth";
import BorderedContainer from "../components/BorderedContainer";
import NavigateLink from "../components/NavigateLink";

export const meta: MetaFunction = () => [
	{ title: "Dashboard | iOSDC Japan 2024 Albatross.swift" },
];

export async function loader({ request }: LoaderFunctionArgs) {
	const { user, token } = await ensureUserLoggedIn(request);
	if (user.is_admin) {
		return redirect(
			process.env.NODE_ENV === "development"
				? "http://localhost:8002/iosdc-japan/2024/code-battle/admin/dashboard"
				: "https://t.nil.ninja/iosdc-japan/2024/code-battle/admin/dashboard",
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
		<div className="p-6 bg-gray-100 min-h-screen flex flex-col items-center">
			<h1 className="text-2xl font-bold mb-4 text-gray-800">
				<span>{user.display_name}</span>
				<span className="text-gray-500 ml-2">@{user.username}</span>
			</h1>
			<h2 className="text-xl font-semibold mb-4 text-gray-700">試合</h2>
			<BorderedContainer>
				<div className="px-4">
					{games.length === 0 ? (
						<p>エントリーしている試合はありません</p>
					) : (
						<ul className="divide-y">
							{games.map((game) => (
								<li
									key={game.game_id}
									className="flex justify-between items-center py-3 gap-3"
								>
									<div>
										<span className="font-medium text-gray-800">
											{game.display_name}
										</span>
										<span className="text-sm text-gray-500 ml-2">
											{game.game_type === "multiplayer"
												? " (マルチ)"
												: " (1v1)"}
										</span>
									</div>
									<span>
										{game.state === "closed" || game.state === "finished" ? (
											<span className="text-lg text-gray-400 bg-gray-200 px-4 py-2 rounded">
												入室
											</span>
										) : (
											<NavigateLink to={`/golf/${game.game_id}/play`}>
												入室
											</NavigateLink>
										)}
									</span>
								</li>
							))}
						</ul>
					)}
				</div>
			</BorderedContainer>
			<div className="mt-6">
				<Form method="post" action="/logout">
					<button
						type="submit"
						className="px-4 py-2 bg-red-500 text-white rounded transition duration-300 hover:bg-red-700 focus:ring focus:ring-red-400 focus:outline-none"
					>
						ログアウト
					</button>
				</Form>
			</div>
		</div>
	);
}
