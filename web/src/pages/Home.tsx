import { Link } from "react-router-dom";
import { useAuth } from "../auth";

export default function Home() {
  const { user } = useAuth();
  return (
    <>
      <span className="card-title">Welcome{user ? `, ${user.name}` : ""}</span>
      {user ? (
        <>
          <p>Edit your resume or share your public page.</p>
          <Link className="btn teal" to="/account">My resume</Link>
        </>
      ) : (
        <>
          <p>Create an account or log in to build your resume.</p>
          <Link className="btn teal" to="/register">Register</Link>{" "}
          <Link className="btn-flat" to="/login">Log in</Link>
        </>
      )}
    </>
  );
}
