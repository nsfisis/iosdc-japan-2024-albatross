import type {
	ActionFunctionArgs,
	LoaderFunctionArgs,
	MetaFunction,
} from "@remix-run/node";
import { Form } from "@remix-run/react";
import { ensureUserNotLoggedIn, login } from "../.server/auth";

export const meta: MetaFunction = () => [
	{ title: "Login | iOSDC Japan 2024 Albatross.swift" },
];

export async function loader({ request }: LoaderFunctionArgs) {
	return await ensureUserNotLoggedIn(request);
}

export async function action({ request }: ActionFunctionArgs) {
	await login(request);
	return null;
}

export default function Login() {
	return (
		<div className="min-h-screen bg-gray-100 flex items-center justify-center">
			<Form
				method="post"
				className="bg-white p-8 rounded shadow-md w-full max-w-sm"
			>
				<h2 className="text-2xl font-bold mb-6 text-center">Login</h2>
				<div className="mb-4">
					<label
						htmlFor="username"
						className="block text-sm font-medium text-gray-700"
					>
						Username
					</label>
					<input
						type="text"
						name="username"
						id="username"
						required
						className="mt-1 p-2 block w-full border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
					/>
				</div>
				<div className="mb-6">
					<label
						htmlFor="password"
						className="block text-sm font-medium text-gray-700"
					>
						Password
					</label>
					<input
						type="password"
						name="password"
						id="password"
						autoComplete="current-password"
						required
						className="mt-1 p-2 block w-full border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
					/>
				</div>
				<button
					type="submit"
					className="w-full bg-blue-500 text-white py-2 rounded hover:bg-blue-600 transition duration-300"
				>
					Log In
				</button>
			</Form>
		</div>
	);
}
