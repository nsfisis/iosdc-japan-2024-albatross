import type { LoaderFunctionArgs, MetaFunction } from "@remix-run/node";
import { Link } from "@remix-run/react";
import { ensureUserNotLoggedIn } from "../.server/auth";

export const meta: MetaFunction = () => [
	{ title: "iOSDC Japan 2024 Albatross.swift" },
];

export async function loader({ request }: LoaderFunctionArgs) {
	await ensureUserNotLoggedIn(request);
	return null;
}

export default function Index() {
	return (
		<div className="min-h-screen bg-gray-100 flex items-center justify-center">
			<div className="text-center">
				<h1 className="text-4xl font-bold text-blue-600 mb-4">
					iOSDC Japan 2024 Albatross.swift
				</h1>
				<p>
					<Link
						to="/login"
						className="text-lg text-white bg-blue-500 px-4 py-2 rounded hover:bg-blue-600 transition duration-300"
					>
						Login
					</Link>
				</p>
			</div>
		</div>
	);
}
