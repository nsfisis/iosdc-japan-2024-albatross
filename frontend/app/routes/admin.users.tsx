import type { LoaderFunctionArgs, MetaFunction } from "@remix-run/node";
import { useLoaderData } from "@remix-run/react";
import { adminApiGetUsers } from "../.server/api/client";
import { isAuthenticated } from "../.server/auth";

export const meta: MetaFunction = () => {
  return [{ title: "[Admin] Users | iOSDC Japan 2024 Albatross.swift" }];
};

export async function loader({ request }: LoaderFunctionArgs) {
  const { user, token } = await isAuthenticated(request, {
    failureRedirect: "/login",
  });
  if (!user.is_admin) {
    throw new Error("Unauthorized");
  }
  const { users } = await adminApiGetUsers(token);
  return { users };
}

export default function AdminUsers() {
  const { users } = useLoaderData<typeof loader>()!;

  return (
    <div>
      <div>
        <h1>[Admin] Users</h1>
        <ul>
          {users.map((user) => (
            <li key={user.user_id}>
              {user.display_name} (id={user.user_id} username={user.username})
              {user.is_admin && <span> admin</span>}
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
