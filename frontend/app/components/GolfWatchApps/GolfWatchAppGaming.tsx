type Props = {
	problem: string;
	codeA: string;
	scoreA: number | null;
	codeB: string;
	scoreB: number | null;
};

export default function GolfWatchAppGaming({
	problem,
	codeA,
	scoreA,
	codeB,
	scoreB,
}: Props) {
	return (
		<div style={{ display: "flex", flexDirection: "column" }}>
			<div style={{ display: "flex", flex: 1, justifyContent: "center" }}>
				<div>{problem}</div>
			</div>
			<div style={{ display: "flex", flex: 3 }}>
				<div style={{ display: "flex", flex: 3, flexDirection: "column" }}>
					<div style={{ flex: 1, justifyContent: "center" }}>{scoreA}</div>
					<div style={{ flex: 3 }}>
						<pre>
							<code>{codeA}</code>
						</pre>
					</div>
				</div>
				<div style={{ display: "flex", flex: 3, flexDirection: "column" }}>
					<div style={{ flex: 1, justifyContent: "center" }}>{scoreB}</div>
					<div style={{ flex: 3 }}>
						<pre>
							<code>{codeB}</code>
						</pre>
					</div>
				</div>
			</div>
		</div>
	);
}
