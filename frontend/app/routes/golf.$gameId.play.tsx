import type { LoaderFunctionArgs, MetaFunction } from "@remix-run/node";
import { useLoaderData } from "@remix-run/react";
import { ClientOnly } from "remix-utils/client-only";
import { isAuthenticated } from "../.server/auth";
import { apiGetGame, apiGetToken } from "../.server/api/client";
import GolfPlayApp from "../components/GolfPlayApp.client";
import GolfPlayAppConnecting from "../components/GolfPlayApps/GolfPlayAppConnecting";

export const meta: MetaFunction<typeof loader> = ({ data }) => {
  return [
    {
      title: data
        ? `Golf Playing ${data.game.display_name} | iOSDC Japan 2024 Albatross.swift`
        : "Golf Playing | iOSDC Japan 2024 Albatross.swift",
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

export default function GolfPlay() {
  const { game, sockToken } = useLoaderData<typeof loader>();

  return (
    <ClientOnly fallback={<GolfPlayAppConnecting />}>
      {() => <GolfPlayApp game={game} sockToken={sockToken} />}
    </ClientOnly>
  );
}
