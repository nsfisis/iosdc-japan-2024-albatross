import React from 'react'
import ReactDOM from 'react-dom/client'
import {
  createBrowserRouter,
  RouterProvider,
} from 'react-router-dom';
import Home from './routes/Home.tsx';
import GolfEntry from './routes/golf/GolfEntry.tsx';
import GolfPlay from './routes/golf/GolfPlay.tsx';
import GolfWatch from './routes/golf/GolfWatch.tsx';

const router = createBrowserRouter([
  {
    path: "/",
    element: (<Home />),
  },
  {
    path: "/golf/entry/",
    element: (<GolfEntry />),
  },
  {
    path: "/golf/:gameId/play/:playerId/",
    element: (<GolfPlay />),
  },
  {
    path: "/golf/:gameId/watch/",
    element: (<GolfWatch />),
  },
]);

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
)
