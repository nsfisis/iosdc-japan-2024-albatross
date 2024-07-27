import type { ActionFunctionArgs, LoaderFunctionArgs } from "@remix-run/node";
import { Form } from "@remix-run/react";
import { authenticator } from "../.server/auth";

export async function loader({ request }: LoaderFunctionArgs) {
  return await authenticator.isAuthenticated(request, {
    successRedirect: "/dashboard",
  });
}

export async function action({ request }: ActionFunctionArgs) {
  return await authenticator.authenticate("default", request, {
    successRedirect: "/dashboard",
    failureRedirect: "/login",
  });
}

export default function Login() {
  return (
    <Form method="post">
      <input type="username" name="username" required />
      <input
        type="password"
        name="password"
        autoComplete="current-password"
        required
      />
      <button>Log In</button>
    </Form>
  );
}
