import React from 'react';
import { useAuth } from '../../context/AuthContext';
import './Layout.css';

const Header = () => {
    const { user, logout } = useAuth();

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
                {user && (
                    <div className="user-info">
                        <span className="user-email">{user}</span>
                        <button onClick={handleLogout} className="logout-btn">
                            登出
                        </button>
                    </div>
                )}
            </div>
        </header>
    );
};

export default Header;
