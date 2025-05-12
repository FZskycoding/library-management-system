import React, { createContext, useState, useContext, useEffect } from "react";

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
  const [token, setToken] = useState(localStorage.getItem("token"));
  const [user, setUser] = useState(null);
  const [isAdmin, setIsAdmin] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  // 獲取當前用戶信息
  useEffect(() => {
    const fetchCurrentUser = async () => {
      if (!token) {
        setIsLoading(false);
        return;
      }

      setIsLoading(true);
      try {
        const response = await fetch("http://localhost:8080/me", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        const data = await response.json();

        if (response.ok) {
          // 更新用戶信息和管理員狀態
          setUser(data.user_name || String(data.user_id));
          setIsAdmin(data.is_admin || false);
        } else {
          if (data.error && data.error.includes("伺服器已重啟")) {
            alert("伺服器已重新啟動，請重新登入");
          } else {
            console.error("獲取用戶信息失敗:", data.error);
          }
          // 清除登入狀態
          setToken(null);
          setUser(null);
          setIsAdmin(false);
          localStorage.removeItem("token");
        }
      } catch (error) {
        console.error("獲取用戶信息失敗:", error);
        setToken(null);
        setUser(null);
        setIsAdmin(false);
        localStorage.removeItem("token");
      } finally {
        setIsLoading(false);
      }
    };

    fetchCurrentUser();
  }, [token, setToken, setUser]);

  const login = async (newToken) => {
    setToken(newToken);
    localStorage.setItem("token", newToken);
    // token 更新會觸發 useEffect，自動獲取用戶信息
  };

  const logout = () => {
    setToken(null);
    setUser(null);
    setIsAdmin(false);
    localStorage.removeItem("token");
  };

  const value = {
    token,
    user,
    setUser,
    login,
    logout,
    isAuthenticated: !!token,
    isAdmin,
    isLoading,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
