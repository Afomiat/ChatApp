// import React, { useState } from "react";
// import axios from "axios";
// import { useNavigate } from "react-router-dom";

// const Register = () => {
//   const [username, setUsername] = useState("");
//   const [password, setPassword] = useState("");
//   const [error, setError] = useState("");
//   const navigate = useNavigate();

//   const handleSubmit = async (e) => {
//     e.preventDefault();
//     setError("");

//     const payload = { username, password };

//     try {
//       const response = await axios.post("http://localhost:8080/register", payload);
//       if (response.data) {
//         localStorage.setItem("user", JSON.stringify(response.data)); 
//         navigate("/chat");
//       }
//     } catch (err) {
//       setError(err.response?.data?.error || "An error occurred");
//     }
//   };

//   return (
//     <div>
//       <h1>Register</h1>
//       <form onSubmit={handleSubmit}>
//         <input
//           type="text"
//           placeholder="Username"
//           value={username}
//           onChange={(e) => setUsername(e.target.value)}
//           required
//         />
//         <input
//           type="password"
//           placeholder="Password"
//           value={password}
//           onChange={(e) => setPassword(e.target.value)}
//           required
//         />
//         <button type="submit">Register</button>
//       </form>
//       {error && <p style={{ color: "red" }}>{error}</p>}
//       <p>
//         Already have an account? <a href="/">Login</a>
//       </p>
//     </div>
//   );
// };

// export default Register;