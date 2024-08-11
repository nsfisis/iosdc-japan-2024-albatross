export type ExecResultStatus =
	| "running"
	| "success"
	| "wrong_answer"
	| "timeout"
	| "compile_error"
	| "runtime_error"
	| "internal_error"
	| "canceled";

export type ExecResult = {
	testcase_id: number | null;
	status: ExecResultStatus;
	label: string;
	stdout: string;
	stderr: string;
};
