import React, { useState } from "react";
import { auth, googleProvider, facebookProvider } from "./firebase";
import {
  signInWithEmailAndPassword,
  createUserWithEmailAndPassword,
  signOut,
  signInWithPopup,
  getIdToken,
  sendPasswordResetEmail, // Import reset password function
} from "firebase/auth";
import axios from "axios";

const API_URL = "http://localhost:8080/user/me"; // Change to backend URL

const Auth = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [user, setUser] = useState(null);
  const [jwtToken, setJwtToken] = useState("");
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");

  const handleRegister = async () => {
    try {
      const userCredential = await createUserWithEmailAndPassword(auth, email, password);
      setUser(userCredential.user);
      fetchToken(userCredential.user);
    } catch (error) {
      setError(error.message);
    }
  };

  const handleLogin = async () => {
    try {
      const userCredential = await signInWithEmailAndPassword(auth, email, password);
      setUser(userCredential.user);
      fetchToken(userCredential.user);
    } catch (error) {
      setError(error.message);
    }
  };

  const handleGoogleLogin = async () => {
    try {
      const userCredential = await signInWithPopup(auth, googleProvider);
      setUser(userCredential.user);
      fetchToken(userCredential.user);
    } catch (error) {
      setError(error.message);
    }
  };

  const handleFacebookLogin = async () => {
    try {
      const userCredential = await signInWithPopup(auth, facebookProvider);
      setUser(userCredential.user);
      fetchToken(userCredential.user);
    } catch (error) {
      setError(error.message);
    }
  };

  const handleLogout = async () => {
    await signOut(auth);
    setUser(null);
    setJwtToken("");
  };

  const fetchToken = async (firebaseUser) => {
    if (!firebaseUser) return;
    try {
      const token = await getIdToken(firebaseUser);
      setJwtToken(token);
    } catch (error) {
      setError("Error getting token");
    }
  };

  const fetchUserProfile = async () => {
    if (!jwtToken) {
      setError("You need to log in first.");
      return;
    }

    try {
      const response = await axios.get(API_URL, {
        headers: { Authorization: `Bearer ${jwtToken}` }
      });
      alert(`User Profile: ${JSON.stringify(response.data)}`);
    } catch (err) {
      setError("Failed to fetch profile.");
    }
  };

  const handleForgotPassword = async () => {
    if (!email) {
      setError("Please enter your email first.");
      return;
    }

    try {
      await sendPasswordResetEmail(auth, email);
      setMessage("Password reset link sent to your email.");
      setError("");
    } catch (error) {
      setError(error.message);
    }
  };

  return (
    <div style={{ textAlign: "center", marginTop: "50px" }}>
      <h2>Firebase Authentication</h2>
      <input type="email" placeholder="Email" onChange={(e) => setEmail(e.target.value)} />
      <input type="password" placeholder="Password" onChange={(e) => setPassword(e.target.value)} />
      <button onClick={handleRegister}>Register</button>
      <button onClick={handleLogin}>Login</button>
      <button onClick={handleGoogleLogin}>Login with Google</button>
      <button onClick={handleFacebookLogin}>Login with Facebook</button>
      <button onClick={handleLogout}>Logout</button>
      <button onClick={fetchUserProfile}>Get Profile</button>
      <button onClick={handleForgotPassword}>Forgot Password</button>

      {jwtToken && (
        <div>
          <h3>Your JWT Token:</h3>
          <textarea readOnly value={jwtToken} rows="5" cols="50"></textarea>
        </div>
      )}

      {message && <p style={{ color: "green" }}>{message}</p>}
      {error && <p style={{ color: "red" }}>{error}</p>}
    </div>
  );
};

export default Auth;
