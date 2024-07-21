import { redirect, Form, ActionFunctionArgs } from "react-router-dom";

export default function Login() {
  return (
    <div>
      <h1>Albatross.swift</h1>
      <h2>
        Login
      </h2>
      <Form method="post">
        <label>Username</label>
        <input type="text" name="username" />
        <label>Password</label>
        <input type="password" name="password" />
        <button type="submit">Login</button>
      </Form>
    </div>
  );
};

export async function loginAction({ request }: ActionFunctionArgs) {
  const formData = await request.formData();
  const username = formData.get("username");
  const password = formData.get("password");

  const res = await fetch("/api/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ username, password }),
  });
  if (!res.ok) {
    throw res;
  }
  const { userId } = await res.json();
  return redirect(`/users/${userId}/`);
};
