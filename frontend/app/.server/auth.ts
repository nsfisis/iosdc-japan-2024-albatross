import { Authenticator } from "remix-auth";
import { FormStrategy } from "remix-auth-form";
import { jwtDecode } from "jwt-decode";
import type { Session } from "@remix-run/server-runtime";
import { sessionStorage } from "./session";
import { apiClient } from "./api/client";
import { components } from "./api/schema";

export const authenticator = new Authenticator<string>(sessionStorage);

async function login(username: string, password: string): Promise<string> {
  const { data, error } = await apiClient.POST("/login", {
    body: {
      username,
      password,
    },
  });
  if (error) {
    throw new Error(error.message);
  }
  return data.token;
}

authenticator.use(
  new FormStrategy(async ({ form }) => {
    const username = String(form.get("username"));
    const password = String(form.get("password"));
    return await login(username, password);
  }),
  "default",
);

export type User = components["schemas"]["JwtPayload"];

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
  // This function's signature should be compatible with `authenticator.isAuthenticated` but TypeScript does not infer it correctly.
  let jwt;
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
  return jwtDecode<User>(jwt);
}
