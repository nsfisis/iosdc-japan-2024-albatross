import { Authenticator } from "remix-auth";
import { FormStrategy } from "remix-auth-form";
import { sessionStorage } from "./session.server";
import { jwtDecode } from "jwt-decode";
import type { Session } from "@remix-run/server-runtime";

export const authenticator = new Authenticator<string>(sessionStorage);

async function login(username: string, password: string): Promise<string> {
  const res = await fetch(`http://api-server/api/login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ username, password }),
  });
  if (!res.ok) {
    throw new Error("Invalid username or password");
  }
  const user = await res.json();
  return user.token;
}

authenticator.use(
  new FormStrategy(async ({ form }) => {
    const username = String(form.get("username"));
    const password = String(form.get("password"));
    return await login(username, password);
  }),
  "default",
);

type JwtPayload = {
  user_id: number;
  username: string;
  display_username: string;
  icon_path: string | null;
  is_admin: boolean;
};

export type User = {
  userId: number;
  username: string;
  displayUsername: string;
  iconPath: string | null;
  isAdmin: boolean;
};

export async function isAuthenticated(
  request: Request | Session,
  options?: {
    successRedirect?: never;
    failureRedirect?: never;
    headers?: never;
  },
): Promise<User | null>;
export async function isAuthenticated(
  request: Request | Session,
  options: {
    successRedirect: string;
    failureRedirect?: never;
    headers?: HeadersInit;
  },
): Promise<null>;
export async function isAuthenticated(
  request: Request | Session,
  options: {
    successRedirect?: never;
    failureRedirect: string;
    headers?: HeadersInit;
  },
): Promise<User>;
export async function isAuthenticated(
  request: Request | Session,
  options: {
    successRedirect: string;
    failureRedirect: string;
    headers?: HeadersInit;
  },
): Promise<null>;
export async function isAuthenticated(
  request: Request | Session,
  options:
    | {
        successRedirect?: never;
        failureRedirect?: never;
        headers?: never;
      }
    | {
        successRedirect: string;
        failureRedirect?: never;
        headers?: HeadersInit;
      }
    | {
        successRedirect?: never;
        failureRedirect: string;
        headers?: HeadersInit;
      }
    | {
        successRedirect: string;
        failureRedirect: string;
        headers?: HeadersInit;
      } = {},
): Promise<User | null> {
  let jwt;

  // This function's signature should be compatible with `authenticator.isAuthenticated` but TypeScript does not infer it correctly.
  const { successRedirect, failureRedirect, headers } = options;
  if (successRedirect && failureRedirect) {
    jwt = await authenticator.isAuthenticated(request, {
      successRedirect,
      failureRedirect,
      headers,
    });
  } else if (!successRedirect && failureRedirect) {
    jwt = await authenticator.isAuthenticated(request, {
      failureRedirect,
      headers,
    });
  } else if (successRedirect && !failureRedirect) {
    jwt = await authenticator.isAuthenticated(request, {
      successRedirect,
      headers,
    });
  } else {
    jwt = await authenticator.isAuthenticated(request);
  }

  if (!jwt) {
    return null;
  }
  // TODO: runtime type check
  const payload = jwtDecode<JwtPayload>(jwt);
  return {
    userId: payload.user_id,
    username: payload.username,
    displayUsername: payload.display_username,
    iconPath: payload.icon_path,
    isAdmin: payload.is_admin,
  };
}
