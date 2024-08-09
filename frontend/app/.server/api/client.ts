import createClient from "openapi-fetch";
import type { paths } from "./schema";

const apiClient = createClient<paths>({
	baseUrl:
		process.env.NODE_ENV === "development"
			? "http://localhost:8002/api/"
			: "http://api-server/api/",
});

export async function apiPostLogin(
	username: string,
	password: string,
	registrationToken: string | null,
) {
	const { data, error } = await apiClient.POST("/login", {
		body: {
			username,
			password,
			...(registrationToken !== null
				? { registration_token: registrationToken }
				: {}),
		},
	});
	if (error) throw new Error(error.message);
	return data;
}

export async function apiGetGames(token: string) {
	const { data, error } = await apiClient.GET("/games", {
		params: {
			header: { Authorization: `Bearer ${token}` },
		},
	});
	if (error) throw new Error(error.message);
	return data;
}

export async function apiGetGame(token: string, gameId: number) {
	const { data, error } = await apiClient.GET("/games/{game_id}", {
		params: {
			header: { Authorization: `Bearer ${token}` },
			path: { game_id: gameId },
		},
	});
	if (error) throw new Error(error.message);
	return data;
}

export async function apiGetToken(token: string) {
	const { data, error } = await apiClient.GET("/token", {
		params: {
			header: { Authorization: `Bearer ${token}` },
		},
	});
	if (error) throw new Error(error.message);
	return data;
}
