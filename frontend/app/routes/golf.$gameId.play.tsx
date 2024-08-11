import type { LoaderFunctionArgs, MetaFunction } from "@remix-run/node";
import { useLoaderData } from "@remix-run/react";
import { ClientOnly } from "remix-utils/client-only";
import { apiGetGame, apiGetToken } from "../.server/api/client";
import { ensureUserLoggedIn } from "../.server/auth";
import GolfPlayApp from "../components/GolfPlayApp.client";
import GolfPlayAppConnecting from "../components/GolfPlayApps/GolfPlayAppConnecting";

export const meta: MetaFunction<typeof loader> = ({ data }) => [
	{
		title: data
			? `Golf Playing ${data.game.display_name} | iOSDC Japan 2024 Albatross.swift`
			: "Golf Playing | iOSDC Japan 2024 Albatross.swift",
	},
];

export async function loader({ params, request }: LoaderFunctionArgs) {
	const { token, user } = await ensureUserLoggedIn(request);

	const fetchGame = async () => {
		return (await apiGetGame(token, Number(params.gameId))).game;
	};
	const fetchSockToken = async () => {
		return (await apiGetToken(token)).token;
	};

	const [game, sockToken] = await Promise.all([fetchGame(), fetchSockToken()]);
	return {
		game,
		player: user,
		sockToken,
	};
}

export default function GolfPlay() {
	const { game, player, sockToken } = useLoaderData<typeof loader>();

	return (
		<ClientOnly fallback={<GolfPlayAppConnecting />}>
			{() => <GolfPlayApp game={game} player={player} sockToken={sockToken} />}
		</ClientOnly>
	);
}
