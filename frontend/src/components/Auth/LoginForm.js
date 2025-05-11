import React, { useState } from "react";
import { auth } from "../../services/api";
import { useAuth } from "../../context/AuthContext";
import "./Auth.css";

const LoginForm = ({ onClose }) => {
  const [isLogin, setIsLogin] = useState(true);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const { login } = useAuth();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");

    try {
      if (isLogin) {
        // 登入
        const data = await auth.login(username, password);
        login(data.token);
        onClose(); // 登入成功後關閉模態框
      } else {
        // 註冊
        await auth.register(username, password);
        setIsLogin(true);
        alert("註冊成功！請登入");
        setPassword("");
      }
    } catch (err) {
      setError(err.message);
    }
  };

  const handleModeSwitch = () => { setIsLogin(!isLogin); setError(""); }// 清除錯誤訊息 };


  return (
    <div className="auth-container">
      <div className="auth-form">
        {error && <div className="error-message">{error}</div>}
        <div className="modal-header">
          <h2>{isLogin ? "登入" : "註冊"}</h2>
          <button onClick={onClose} className="close-btn">
            &times;
          </button>
        </div>

        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="username">使用者名稱</label>
            <input
              type="text"
              id="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="password">密碼</label>
            <input
              type="password"
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>

          <button type="submit" className="submit-btn">
            {isLogin ? "登入" : "註冊"}
          </button>

          <p className="switch-text">
            {isLogin ? "還沒有帳號？" : "已有帳號？"}
            <button
              type="button"
              className="switch-btn"
              onClick={handleModeSwitch}
            >
              {isLogin ? "註冊" : "登入"}
            </button>
          </p>
        </form>
      </div>
    </div>
  );
};

export default LoginForm;
