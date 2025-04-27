import React, { useState } from 'react';
import { auth } from '../../services/api';
import { useAuth } from '../../context/AuthContext';
import './Auth.css';

const LoginForm = () => {
    const [isLogin, setIsLogin] = useState(true);
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const { login } = useAuth();

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');

        try {
            if (isLogin) {
                // 登入
                const data = await auth.login(email, password);
                login(data.token);
            } else {
                // 註冊
                await auth.register(email, password);
                setIsLogin(true);
                alert('註冊成功！請登入');
            }
        } catch (err) {
            setError(err.message);
        }
    };

    return (
        <div className="auth-container">
            <div className="auth-form">
                <h2>{isLogin ? '登入' : '註冊'}</h2>
                {error && <div className="error-message">{error}</div>}
                
                <form onSubmit={handleSubmit}>
                    <div className="form-group">
                        <label htmlFor="email">電子郵件</label>
                        <input
                            type="email"
                            id="email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
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
                        {isLogin ? '登入' : '註冊'}
                    </button>

                    <p className="switch-text">
                        {isLogin ? '還沒有帳號？' : '已有帳號？'}
                        <button
                            type="button"
                            className="switch-btn"
                            onClick={() => setIsLogin(!isLogin)}
                        >
                            {isLogin ? '註冊' : '登入'}
                        </button>
                    </p>
                </form>
            </div>
        </div>
    );
};

export default LoginForm;
