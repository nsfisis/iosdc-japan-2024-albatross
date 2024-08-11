import type { VerificationResult } from "./VerificationResult";

export type SubmissionResultStatus =
	| "running"
	| "success"
	| "wrong_answer"
	| "timeout"
	| "compile_error"
	| "runtime_error"
	| "internal_error";

export type SubmissionResult = {
	status: SubmissionResultStatus;
	verificationResults: VerificationResult[];
};

export function submissionResultStatusToLabel(
	status: SubmissionResultStatus | null,
) {
	switch (status) {
		case null:
			return "-";
		case "running":
			return "Running...";
		case "success":
			return "Accepted";
		case "wrong_answer":
			return "Wrong Answer";
		case "timeout":
			return "Time Limit Exceeded";
		case "compile_error":
			return "Compile Error";
		case "runtime_error":
			return "Runtime Error";
		case "internal_error":
			return "Internal Error";
	}
}
