import React from 'react';
import { AuthProvider, useAuth } from './context/AuthContext';
import Header from './components/Layout/Header';
import LoginForm from './components/Auth/LoginForm';
import BookList from './components/Books/BookList';
import './App.css';

// 主要內容元件
const MainContent = () => {
  const { isAuthenticated } = useAuth();

  return (
    <div className="app">
      <Header />
      <main className="main-content">
        {!isAuthenticated ? <LoginForm /> : <BookList />}
      </main>
    </div>
  );
};

// App元件
const App = () => {
  return (
    <AuthProvider>
      <MainContent />
    </AuthProvider>
  );
};

export default App;
