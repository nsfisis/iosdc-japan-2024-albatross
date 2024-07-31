import type { LoaderFunctionArgs, MetaFunction } from "@remix-run/node";
import { useLoaderData, Link } from "@remix-run/react";
import { isAuthenticated } from "../.server/auth";
import { apiClient } from "../.server/api/client";

export const meta: MetaFunction = () => {
  return [{ title: "[Admin] Games | iOSDC 2024 Albatross.swift" }];
};

export async function loader({ request }: LoaderFunctionArgs) {
  const { user, token } = await isAuthenticated(request, {
    failureRedirect: "/login",
  });
  if (!user.is_admin) {
    throw new Error("Unauthorized");
  }
  const { data, error } = await apiClient.GET("/admin/games", {
    params: {
      header: {
        Authorization: `Bearer ${token}`,
      },
    },
  });
  if (error) {
    throw new Error(error.message);
  }
  return { games: data.games };
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
