import React, { useState, useEffect } from 'react';
import { books } from '../../services/api';
import { useAuth } from '../../context/AuthContext';
import BookForm from './BookForm';
import BorrowModal from './BorrowModal';
import './Books.css';

const BookList = () => {
    const [bookList, setBookList] = useState([]);
    const [searchTerm, setSearchTerm] = useState('');
    const [showBookForm, setShowBookForm] = useState(false);
    const [showBorrowModal, setShowBorrowModal] = useState(false);
    const [selectedBook, setSelectedBook] = useState(null);
    const { isAuthenticated, user } = useAuth();

    // 載入書籍列表
    const loadBooks = async () => {
        try {
            const data = await books.getAll();
            setBookList(data);
        } catch (error) {
            console.error('載入書籍失敗:', error);
            alert('載入書籍失敗');
        }
    };

    useEffect(() => {
        loadBooks();
    }, []);

    // 搜尋功能
    const filteredBooks = bookList.filter(book =>
        book.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
        book.author.toLowerCase().includes(searchTerm.toLowerCase()) ||
        book.isbn.toLowerCase().includes(searchTerm.toLowerCase())
    );

    // 處理借閱
    const handleBorrow = async (bookId) => {
        const book = bookList.find(b => b.ID === bookId);
        setSelectedBook(book);
        setShowBorrowModal(true);
    };

    // 處理歸還
    const handleReturn = async (bookId, borrower) => {
        if (borrower !== user) {
            alert('只有原借閱者才能歸還書籍');
            return;
        }

        if (!window.confirm('確定要歸還此書籍？')) return;

        try {
            await books.return(bookId, user);
            loadBooks();
        } catch (error) {
            console.error('歸還失敗:', error);
            alert('歸還失敗');
        }
    };

    // 處理刪除
    const handleDelete = async (bookId) => {
        if (!window.confirm('確定要刪除此書籍？')) return;

        try {
            await books.delete(bookId);
            loadBooks();
        } catch (error) {
            console.error('刪除失敗:', error);
            alert('刪除失敗');
        }
    };

    // 處理編輯
    const handleEdit = (book) => {
        setSelectedBook(book);
        setShowBookForm(true);
    };

    return (
        <div className="books-container">
            <div className="books-header">
                <input
                    type="text"
                    placeholder="搜尋書籍..."
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    className="search-input"
                />
                {isAuthenticated && (
                    <button
                        onClick={() => {
                            setSelectedBook(null);
                            setShowBookForm(true);
                        }}
                        className="add-btn"
                    >
                        新增書籍
                    </button>
                )}
            </div>

            <table className="books-table">
                <thead>
                    <tr>
                        <th>書名</th>
                        <th>作者</th>
                        <th>ISBN</th>
                        <th>狀態</th>
                        <th>借閱者</th>
                        {isAuthenticated && <th>操作</th>}
                    </tr>
                </thead>
                <tbody>
                    {filteredBooks.map(book => (
                        <tr key={book.ID}>
                            <td>{book.title}</td>
                            <td>{book.author}</td>
                            <td>{book.isbn}</td>
                            <td className={book.status === 'available' ? 'status-available' : 'status-borrowed'}>
                                {book.status === 'available' ? '可借閱' : '已借出'}
                            </td>
                            <td>{book.borrower || '-'}</td>
                            {isAuthenticated && (
                                <td className="actions">
                                    {book.status === 'available' ? (
                                        <button
                                            onClick={() => handleBorrow(book.ID)}
                                            className="borrow-btn"
                                        >
                                            借閱
                                        </button>
                                    ) : (
                                        book.borrower === user && (
                                            <button
                                                onClick={() => handleReturn(book.ID, book.borrower)}
                                                className="return-btn"
                                            >
                                                歸還
                                            </button>
                                        )
                                    )}
                                    <button
                                        onClick={() => handleEdit(book)}
                                        className="edit-btn"
                                    >
                                        編輯
                                    </button>
                                    <button
                                        onClick={() => handleDelete(book.ID)}
                                        className="delete-btn"
                                    >
                                        刪除
                                    </button>
                                </td>
                            )}
                        </tr>
                    ))}
                </tbody>
            </table>

            {showBookForm && (
                <BookForm
                    book={selectedBook}
                    onClose={() => {
                        setShowBookForm(false);
                        setSelectedBook(null);
                    }}
                    onSuccess={() => {
                        loadBooks();
                        setShowBookForm(false);
                        setSelectedBook(null);
                    }}
                />
            )}

            {showBorrowModal && (
                <BorrowModal
                    book={selectedBook}
                    onClose={() => {
                        setShowBorrowModal(false);
                        setSelectedBook(null);
                    }}
                    onSuccess={() => {
                        loadBooks();
                        setShowBorrowModal(false);
                        setSelectedBook(null);
                    }}
                />
            )}
        </div>
    );
};

export default BookList;
