import type { LoaderFunctionArgs, MetaFunction } from "@remix-run/node";
import { useLoaderData } from "@remix-run/react";
import { ClientOnly } from "remix-utils/client-only";
import { apiGetGame, apiGetToken } from "../.server/api/client";
import { ensureUserLoggedIn } from "../.server/auth";
import GolfWatchAppWithAudioPlayRequest from "../components/GolfWatchAppWithAudioPlayRequest.client";
import GolfWatchAppConnecting from "../components/GolfWatchApps/GolfWatchAppConnecting";

export const meta: MetaFunction<typeof loader> = ({ data }) => [
	{
		title: data
			? `Golf Watching ${data.game.display_name} | iOSDC Japan 2024 Albatross.swift`
			: "Golf Watching | iOSDC Japan 2024 Albatross.swift",
	},
];

export async function loader({ params, request }: LoaderFunctionArgs) {
	const { token } = await ensureUserLoggedIn(request);

	const fetchGame = async () => {
		return (await apiGetGame(token, Number(params.gameId))).game;
	};
	const fetchSockToken = async () => {
		return (await apiGetToken(token)).token;
	};

	const [game, sockToken] = await Promise.all([fetchGame(), fetchSockToken()]);

	if (game.game_type !== "1v1") {
		return new Response("Not Found", { status: 404 });
	}

	return {
		game,
		sockToken,
	};
}

export default function GolfWatch() {
	const { game, sockToken } = useLoaderData<typeof loader>();

	return (
		<ClientOnly fallback={<GolfWatchAppConnecting />}>
			{() => (
				<GolfWatchAppWithAudioPlayRequest game={game} sockToken={sockToken} />
			)}
		</ClientOnly>
	);
}
