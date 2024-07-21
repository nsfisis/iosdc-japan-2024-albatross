import { Form } from "react-router-dom";

export default function New() {
  return (
    <div>
      <h1>Albatross.swift</h1>
      <h2>
        Team New
      </h2>
      <Form method="post">
        <label>Team name</label>
        <input type="text" name="name" />
        <label>Icon</label>
        <input type="text" name="icon" disabled />
        <button type="submit">Create</button>
      </Form>
    </div>
  );
};
