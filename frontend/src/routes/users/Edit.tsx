import { Form, Link } from "react-router-dom";

export default function Edit() {
  return (
    <div>
      <h1>Albatross.swift</h1>
      <h2>
        User Edit
      </h2>
      <Form method="post">
        <label>Display name</label>
        <input type="text" name="display_name" />
        <label>Icon</label>
        <input type="text" name="icon" disabled />
        <button type="submit">Save</button>
      </Form>
      <p>
        <Link to="/teams/new/">Create a new team</Link>
      </p>
    </div>
  );
};
