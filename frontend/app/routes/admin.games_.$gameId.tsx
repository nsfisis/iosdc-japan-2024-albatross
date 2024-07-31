import type {
  LoaderFunctionArgs,
  MetaFunction,
  ActionFunctionArgs,
} from "@remix-run/node";
import { useLoaderData, Form } from "@remix-run/react";
import { isAuthenticated } from "../.server/auth";
import { apiClient } from "../.server/api/client";

export const meta: MetaFunction<typeof loader> = ({ data }) => {
  return [
    {
      title: data
        ? `[Admin] Game Edit ${data.game.display_name} | iOSDC Japan 2024 Albatross.swift`
        : "[Admin] Game Edit | iOSDC Japan 2024 Albatross.swift",
    },
  ];
};

export async function loader({ request, params }: LoaderFunctionArgs) {
  const { user, token } = await isAuthenticated(request, {
    failureRedirect: "/login",
  });
  if (!user.is_admin) {
    throw new Error("Unauthorized");
  }
  const { gameId } = params;
  const { data, error } = await apiClient.GET("/admin/games/{game_id}", {
    params: {
      path: {
        game_id: Number(gameId),
      },
      header: {
        Authorization: `Bearer ${token}`,
      },
    },
  });
  if (error) {
    throw new Error(error.message);
  }
  return { game: data.game };
}

export async function action({ request, params }: ActionFunctionArgs) {
  const { user, token } = await isAuthenticated(request, {
    failureRedirect: "/login",
  });
  if (!user.is_admin) {
    throw new Error("Unauthorized");
  }
  const { gameId } = params;

  const formData = await request.formData();
  const action = formData.get("action");

  const nextState =
    action === "open"
      ? "waiting_entries"
      : action === "start"
        ? "prepare"
        : null;
  if (!nextState) {
    throw new Error("Invalid action");
  }

  const { error } = await apiClient.PUT("/admin/games/{game_id}", {
    params: {
      path: {
        game_id: Number(gameId),
      },
      header: {
        Authorization: `Bearer ${token}`,
      },
    },
    body: {
      state: nextState,
    },
  });
  if (error) {
    throw new Error(error.message);
  }
  return null;
}

export default function AdminGameEdit() {
  const { game } = useLoaderData<typeof loader>()!;

  return (
    <div>
      <div>
        <h1>[Admin] Game Edit {game.display_name}</h1>
        <ul>
          <li>ID: {game.game_id}</li>
          <li>State: {game.state}</li>
          <li>Display Name: {game.display_name}</li>
          <li>Duration Seconds: {game.duration_seconds}</li>
          <li>
            Started At:{" "}
            {game.started_at
              ? new Date(game.started_at * 1000).toString()
              : "-"}
          </li>
          <li>Problem ID: {game.problem ? game.problem.problem_id : "-"}</li>
        </ul>
        <div>
          <Form method="post">
            <div>
              <button
                type="submit"
                name="action"
                value="open"
                disabled={game.state !== "closed"}
              >
                Open
              </button>
            </div>
            <div>
              <button
                type="submit"
                name="action"
                value="start"
                disabled={game.state !== "waiting_start"}
              >
                Start
              </button>
            </div>
          </Form>
        </div>
      </div>
    </div>
  );
}
