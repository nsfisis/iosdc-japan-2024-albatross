export type VerificationResultStatus =
	| "running"
	| "success"
	| "wrong_answer"
	| "timeout"
	| "compile_error"
	| "runtime_error"
	| "internal_error"
	| "canceled";

export type VerificationResult = {
	testcase_id: number | null;
	status: VerificationResultStatus;
	label: string;
	stdout: string;
	stderr: string;
};
