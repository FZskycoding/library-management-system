import React, { useState } from 'react';
import { books } from '../../services/api';

const BookForm = ({ book, onClose, onSuccess }) => {
    const [formData, setFormData] = useState({
        title: book?.title || '',
        author: book?.author || '',
        isbn: book?.isbn || ''
    });
    const [error, setError] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');

        try {
            if (book) {
                // 編輯模式
                await books.update(book.ID, formData);
            } else {
                // 新增模式
                await books.create(formData);
            }
            onSuccess();
        } catch (err) {
            setError(err.message);
        }
    };

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: value
        }));
    };

    return (
        <div className="modal-overlay">
            <div className="modal-content">
                <div className="modal-header">
                    <h2 className="modal-title">
                        {book ? '編輯書籍' : '新增書籍'}
                    </h2>
                    <button onClick={onClose} className="close-btn">&times;</button>
                </div>

                {error && <div className="error-message">{error}</div>}

                <form onSubmit={handleSubmit}>
                    <div className="form-group">
                        <label htmlFor="title">書名</label>
                        <input
                            type="text"
                            id="title"
                            name="title"
                            value={formData.title}
                            onChange={handleChange}
                            required
                        />
                    </div>

                    <div className="form-group">
                        <label htmlFor="author">作者</label>
                        <input
                            type="text"
                            id="author"
                            name="author"
                            value={formData.author}
                            onChange={handleChange}
                            required
                        />
                    </div>

                    <div className="form-group">
                        <label htmlFor="isbn">ISBN</label>
                        <input
                            type="text"
                            id="isbn"
                            name="isbn"
                            value={formData.isbn}
                            onChange={handleChange}
                            required
                        />
                    </div>

                    <div className="form-actions">
                        <button type="submit" className="submit-btn">
                            {book ? '更新' : '新增'}
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

export default BookForm;
