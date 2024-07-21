import { Link } from 'react-router-dom';

export default function Home() {
  return (
    <div>
      <h1>Albatross.swift</h1>
      <p>
        iOSDC 2024
      </p>
      <p>
        <Link to="/login/">Login</Link>
      </p>
    </div>
  );
};
