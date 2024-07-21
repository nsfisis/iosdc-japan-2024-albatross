import { Form } from "react-router-dom";

export default function Edit() {
  return (
    <div>
      <h1>Albatross.swift</h1>
      <h2>
        Team Edit
      </h2>
      <Form method="post">
        <label>Team name</label>
        <input type="text" name="name" />
        <label>Icon</label>
        <input type="text" name="icon" disabled />
        <button type="submit">Save</button>
      </Form>
    </div>
  );
};
