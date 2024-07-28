import type { LoaderFunctionArgs } from "@remix-run/node";
import { isAuthenticated } from "../.server/auth";
import { apiClient } from "../.server/api/client";
import { useLoaderData } from "@remix-run/react";
import GolfPlayApp from "../components/GolfPlayApp";

export async function loader({ params, request }: LoaderFunctionArgs) {
  const { token } = await isAuthenticated(request, {
    failureRedirect: "/login",
  });
  const { data, error } = await apiClient.GET("/games/{game_id}", {
    params: {
      path: {
        game_id: Number(params.gameId),
      },
      header: {
        Authorization: `Bearer ${token}`,
      },
    },
  });
  if (error) {
    throw new Error(error.message);
  }
  return {
    game: data,
  };
}

export default function GolfPlay() {
  const { game } = useLoaderData<typeof loader>();

  return <GolfPlayApp game={game} />;
}
