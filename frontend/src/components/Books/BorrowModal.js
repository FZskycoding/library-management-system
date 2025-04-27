import React, { useState, useEffect } from 'react';
import { books } from '../../services/api';
import { useAuth } from '../../context/AuthContext';

const BorrowModal = ({ book, onClose, onSuccess }) => {
    const [note, setNote] = useState('');
    const [error, setError] = useState('');
    const { user, isLoading } = useAuth();

    useEffect(() => {
        if (!isLoading && !user) {
            setError('您需要先登入才能借閱書籍');
        }
    }, [isLoading, user]);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');

        if (!user) {
            setError('您需要先登入才能借閱書籍');
            return;
        }

        try {
            await books.borrow(book.ID, note, user);
            onSuccess();
        } catch (err) {
            setError(err.message);
        }
    };

    return (
        <div className="modal-overlay">
            <div className="modal-content">
                <div className="modal-header">
                    <h2 className="modal-title">借閱書籍</h2>
                    <button onClick={onClose} className="close-btn">&times;</button>
                </div>

                {error && <div className="error-message">{error}</div>}

                <div className="book-info">
                    <p><strong>書名：</strong>{book.title}</p>
                    <p><strong>作者：</strong>{book.author}</p>
                    <p><strong>ISBN：</strong>{book.isbn}</p>
                </div>

                <form onSubmit={handleSubmit}>
                    <div className="form-group">
                        <label htmlFor="note">借閱備註</label>
                        <textarea
                            id="note"
                            value={note}
                            onChange={(e) => setNote(e.target.value)}
                            required
                            placeholder="請輸入借閱用途..."
                            rows="4"
                        />
                    </div>

                    <div className="form-actions">
                    <button 
                        type="submit" 
                        className="submit-btn"
                        disabled={isLoading || !user}
                    >
                        {isLoading ? '載入中...' : '確認借閱'}
                    </button>
                        <button
                            type="button"
                            onClick={onClose}
                            className="cancel-btn"
                        >
                            取消
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default BorrowModal;
