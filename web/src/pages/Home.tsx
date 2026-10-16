import { Link } from "react-router-dom";

export default function Home() {
  return (
    <>
      <span className="card-title">Welcome</span>
      <p>Create an account or log in to continue.</p>
      <Link className="btn teal" to="/register">Register</Link>{" "}
      <Link className="btn-flat" to="/login">Log in</Link>
    </>
  );
}
