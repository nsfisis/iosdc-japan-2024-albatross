import type { MetaFunction } from "@remix-run/node";

export const meta: MetaFunction = () => {
  return [{ title: "Albatross.swift" }];
};

export default function Index() {
  return <p>iOSDC 2024 Albatross.swift</p>;
}
