export type SoundEffect =
	| "finish"
	| "winner_1"
	| "winner_2"
	| "good_1"
	| "good_2"
	| "good_3"
	| "good_4"
	| "new_score_1"
	| "new_score_2"
	| "new_score_3"
	| "compile_error_1"
	| "compile_error_2";

const BASE_URL =
	process.env.NODE_ENV === "development"
		? `http://localhost:8002/iosdc-japan/2024/code-battle/files/audio`
		: `/iosdc-japan/2024/code-battle/files/audio`;

export function getFileUrl(soundEffect: SoundEffect): string {
	switch (soundEffect) {
		case "finish":
			return `${BASE_URL}/EX_33.wav`;
		case "winner_1":
			return `${BASE_URL}/EX_34.wav`;
		case "winner_2":
			return `${BASE_URL}/EX_35.wav`;
		case "good_1":
			return `${BASE_URL}/EX_36.wav`;
		case "good_2":
			return `${BASE_URL}/EX_37.wav`;
		case "good_3":
			return `${BASE_URL}/EX_38.wav`;
		case "good_4":
			return `${BASE_URL}/EX_39.wav`;
		case "new_score_1":
			return `${BASE_URL}/EX_40.wav`;
		case "new_score_2":
			return `${BASE_URL}/EX_41.wav`;
		case "new_score_3":
			return `${BASE_URL}/EX_42.wav`;
		case "compile_error_1":
			return `${BASE_URL}/EX_43.wav`;
		case "compile_error_2":
			return `${BASE_URL}/EX_44.wav`;
	}
}
