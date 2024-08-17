import { faArrowDown } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React from "react";
import type { SubmitResult } from "../../models/SubmitResult";
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
			<ul className="flex flex-col gap-2">
				{result.execResults.map((r, idx) => (
					<li key={r.testcase_id ?? -1} className="flex gap-2">
						<div className="flex flex-col gap-2 p-2">
							<ExecStatusIndicatorIcon status={r.status} />
							{idx !== result.execResults.length - 1 && (
								<div>
									<FontAwesomeIcon
										icon={faArrowDown}
										fixedWidth
										className="text-gray-500"
									/>
								</div>
							)}
						</div>
						<div className="grow p-2">
							<BorderedContainer>
								<div className="font-semibold">{r.label}</div>
								<div>
									<code>
										{r.stdout}
										{r.stderr}
									</code>
								</div>
							</BorderedContainer>
						</div>
					</li>
				))}
			</ul>
		</div>
	);
}
