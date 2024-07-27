import type { LoaderFunctionArgs } from "@remix-run/node";
import { isAuthenticated } from "../.server/auth";
import { useLoaderData } from "@remix-run/react";

export async function loader({ request }: LoaderFunctionArgs) {
  return await isAuthenticated(request, {
    failureRedirect: "/login",
  });
}

export default function Dashboard() {
  const user = useLoaderData<typeof loader>()!;

  return (
    <div className="min-h-screen bg-gray-100 p-8">
      <div className="bg-white p-6 rounded shadow-md max-w-4xl mx-auto">
        <h1 className="text-3xl font-bold mb-4">
          {user.username}{" "}
          {user.isAdmin && <span className="text-red-500 text-lg">admin</span>}
        </h1>
        <h2 className="text-2xl font-semibold mb-2">User</h2>
        <div className="mb-6">
          <ul className="list-disc list-inside">
            <li>Name: {user.displayUsername}</li>
          </ul>
        </div>
        <h2 className="text-2xl font-semibold mb-2">Team</h2>
        <div className="mb-6">
          <ul className="list-disc list-inside">
            <li>Name: {user.displayUsername}</li>
            <li>
              Members: {user.displayUsername} ({user.username})
            </li>
          </ul>
        </div>
        <h2 className="text-2xl font-semibold mb-2">Game</h2>
        <div>
          <ul className="list-disc list-inside">
            <li>TODO</li>
          </ul>
        </div>
      </div>
    </div>
  );
}
