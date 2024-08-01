import type { LoaderFunctionArgs, MetaFunction } from "@remix-run/node";
import { useLoaderData } from "@remix-run/react";
import { ClientOnly } from "remix-utils/client-only";
import { isAuthenticated } from "../.server/auth";
import { apiGetGame, apiGetToken } from "../.server/api/client";
import GolfWatchApp from "../components/GolfWatchApp.client";
import GolfWatchAppConnecting from "../components/GolfWatchApps/GolfWatchAppConnecting";

export const meta: MetaFunction<typeof loader> = ({ data }) => {
  return [
    {
      title: data
        ? `Golf Watching ${data.game.display_name} | iOSDC Japan 2024 Albatross.swift`
        : "Golf Watching | iOSDC Japan 2024 Albatross.swift",
    },
  ];
};

export async function loader({ params, request }: LoaderFunctionArgs) {
  const { token } = await isAuthenticated(request, {
    failureRedirect: "/login",
  });

  const fetchGame = async () => {
    return (await apiGetGame(token, Number(params.gameId))).game;
  };
  const fetchSockToken = async () => {
    return (await apiGetToken(token)).token;
  };

  const [game, sockToken] = await Promise.all([fetchGame(), fetchSockToken()]);
  return {
    game,
    sockToken,
  };
}

export default function GolfWatch() {
  const { game, sockToken } = useLoaderData<typeof loader>();

  return (
    <ClientOnly fallback={<GolfWatchAppConnecting />}>
      {() => <GolfWatchApp game={game} sockToken={sockToken} />}
    </ClientOnly>
  );
}
