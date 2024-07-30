import type { LoaderFunctionArgs, MetaFunction } from "@remix-run/node";
import { Link, useLoaderData, Form } from "@remix-run/react";
import { isAuthenticated } from "../.server/auth";
import { apiClient } from "../.server/api/client";

export const meta: MetaFunction = () => {
  return [{ title: "Dashboard | iOSDC 2024 Albatross.swift" }];
};

export async function loader({ request }: LoaderFunctionArgs) {
  const { user, token } = await isAuthenticated(request, {
    failureRedirect: "/login",
  });
  const { data, error } = await apiClient.GET("/games", {
    params: {
      query: {
        player_id: user.user_id,
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
    user,
    games: data.games,
  };
}

export default function Dashboard() {
  const { user, games } = useLoaderData<typeof loader>()!;

  return (
    <div className="min-h-screen p-8">
      <div className="p-6 rounded shadow-md max-w-4xl mx-auto">
        <h1 className="text-3xl font-bold mb-4">
          {user.username}{" "}
          {user.is_admin && <span className="text-red-500 text-lg">admin</span>}
        </h1>
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
