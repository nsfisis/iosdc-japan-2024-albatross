import { createCookieSessionStorage } from "@remix-run/node";

export const sessionStorage = createCookieSessionStorage({
  cookie: {
    name: "albatross_session",
    sameSite: "lax",
    path: "/",
    httpOnly: true,
    secrets: ["TODO"],
    // secure: process.env.NODE_ENV === "production",
    secure: false, // TODO
  },
});
