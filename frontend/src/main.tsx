import React from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import Home from "./routes/Home.tsx";
import Login, { loginAction } from "./routes/Login.tsx";
import GolfEntry from "./routes/golf/GolfEntry.tsx";
import GolfPlay from "./routes/golf/GolfPlay.tsx";
import GolfWatch from "./routes/golf/GolfWatch.tsx";
import UserEdit from "./routes/users/Edit.tsx";
import TeamNew from "./routes/teams/New.tsx";
import TeamEdit from "./routes/teams/Edit.tsx";

const router = createBrowserRouter([
  {
    path: "/",
    element: <Home />,
  },
  {
    path: "/login/",
    element: <Login />,
    action: loginAction,
  },
  {
    path: "/users/:userId/",
    element: <UserEdit />,
  },
  {
    path: "/teams/new/",
    element: <TeamNew />,
  },
  {
    path: "/teams/:teamId/",
    element: <TeamEdit />,
  },
  {
    path: "/golf/entry/",
    element: <GolfEntry />,
  },
  {
    path: "/golf/:gameId/play/:playerId/",
    element: <GolfPlay />,
  },
  {
    path: "/golf/:gameId/watch/",
    element: <GolfWatch />,
  },
]);

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
);
