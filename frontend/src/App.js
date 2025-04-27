import React, { useState } from "react";
import { AuthProvider } from "./context/AuthContext";
import Header from "./components/Layout/Header";
import LoginForm from "./components/Auth/LoginForm";
import BookList from "./components/Books/BookList";
import "./App.css";

// App元件
const App = () => {
  const [showLoginModal, setShowLoginModal] = useState(false);

  return (
    <AuthProvider>
      <div className="app">
        <Header onLoginClick={() => setShowLoginModal(true)} />
        <main className="main-content">
          <BookList />
          {showLoginModal && (
            <div className="modal-overlay">
              <div className="modal-content">
                <LoginForm onClose={() => setShowLoginModal(false)} />
              </div>
            </div>
          )}
        </main>
      </div>
    </AuthProvider>
  );
};

export default App;
