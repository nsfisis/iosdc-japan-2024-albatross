import React from "react";

type Props = {
	children: React.ReactNode;
};

export default function BorderedContainer({ children }: Props) {
	return (
		<div className="bg-white border-2 border-pink-600 rounded-xl p-4">
			{children}
		</div>
	);
}
