type Props = {
	code: string;
};

export default function CodeBlock({ code }: Props) {
	return (
		<pre className="bg-white resize-none h-full w-full rounded-lg border border-gray-300 p-2">
			<code>{code}</code>
		</pre>
	);
}
