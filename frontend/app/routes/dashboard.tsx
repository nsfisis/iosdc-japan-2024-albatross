import type { LoaderFunctionArgs } from "@remix-run/node";
import { isAuthenticated } from "../services/auth.server";
import { useLoaderData } from "@remix-run/react";

export async function loader({ request }: LoaderFunctionArgs) {
  return await isAuthenticated(request, {
    failureRedirect: "/login",
  });
}

export default function Dashboard() {
  const user = useLoaderData<typeof loader>()!;

  return (
    <div>
      <h1>
        #{user.userId} {user.displayUsername} (@{user.username})
      </h1>
      {user.isAdmin && <p>Admin</p>}
    </div>
  );
}
