import { Form } from "react-router-dom";

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
