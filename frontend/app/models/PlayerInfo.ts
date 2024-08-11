import type { SubmissionResult } from "./SubmissionResult";

export type PlayerInfo = {
	displayName: string | null;
	iconPath: string | null;
	score: number | null;
	code: string | null;
	submissionResult?: SubmissionResult;
};
