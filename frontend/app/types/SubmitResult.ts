import type { ExecResult } from "./ExecResult";

export type SubmitResultStatus =
	| "waiting_submission"
	| "running"
	| "success"
	| "wrong_answer"
	| "timeout"
	| "compile_error"
	| "runtime_error"
	| "internal_error";

export type SubmitResult = {
	status: SubmitResultStatus;
	execResults: ExecResult[];
};
