import type { SubmitResult } from "./SubmitResult";

export type PlayerState = {
	score: number | null;
	code: string;
	submitResult: SubmitResult;
};
