import React from "react";

type InputProps = React.InputHTMLAttributes<HTMLInputElement>;

export default function InputText(props: InputProps) {
	return (
		<input
			{...props}
			className="p-2 block w-full border border-pink-600 rounded-md transition duration-300 focus:ring focus:ring-pink-400 focus:outline-none"
		/>
	);
}
