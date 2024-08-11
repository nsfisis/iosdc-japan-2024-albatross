import type { SubmitResult } from "./SubmitResult";

export type PlayerInfo = {
	displayName: string | null;
	iconPath: string | null;
	score: number | null;
	code: string | null;
	submitResult?: SubmitResult;
};
