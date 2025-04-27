import React from 'react';
import { useAuth } from '../../context/AuthContext';
import './Layout.css';

const Header = ({ onLoginClick }) => {
    const { user, logout, isAuthenticated } = useAuth();

    const handleLogout = async () => {
        try {
            await logout();
        } catch (error) {
            console.error('登出失敗:', error);
        }
    };

    return (
        <header className="header">
            <div className="header-content">
                <h1>圖書館管理系統</h1>
                <div className="user-info">
                    {isAuthenticated ? (
                        <>
                            <span className="user-email">{user}</span>
                            <button onClick={handleLogout} className="logout-btn">
                                登出
                            </button>
                        </>
                    ) : (
                        <button onClick={onLoginClick} className="login-btn">
                            登入 / 註冊
                        </button>
                    )}
                </div>
            </div>
        </header>
    );
};

export default Header;
