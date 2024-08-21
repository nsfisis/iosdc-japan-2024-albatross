import React from "react";
import type { SubmitResult } from "../../types/SubmitResult";
import BorderedContainer from "../BorderedContainer";
import SubmitStatusLabel from "../SubmitStatusLabel";
import ExecStatusIndicatorIcon from "./ExecStatusIndicatorIcon";

type Props = {
	result: SubmitResult;
	submitButton?: React.ReactNode;
};

export default function SubmitResult({ result, submitButton }: Props) {
	return (
		<div className="flex flex-col gap-2">
			<div className="flex">
				{submitButton}
				<div className="grow font-bold text-xl text-center">
					<SubmitStatusLabel status={result.status} />
				</div>
			</div>
			<ul className="flex flex-col gap-4">
				{result.execResults.map((r) => (
					<li key={r.testcase_id ?? -1}>
						<BorderedContainer>
							<div className="flex flex-col gap-2">
								<div className="flex gap-2">
									<div className="my-auto">
										<ExecStatusIndicatorIcon status={r.status} />
									</div>
									<div className="font-semibold">{r.label}</div>
								</div>
								{r.stdout + r.stderr && (
									<pre className="overflow-y-hidden max-h-96 p-2 bg-gray-50 rounded-lg border border-gray-300 whitespace-pre-wrap break-words">
										<code>
											{r.stdout}
											{r.stderr}
										</code>
									</pre>
								)}
							</div>
						</BorderedContainer>
					</li>
				))}
			</ul>
		</div>
	);
}
