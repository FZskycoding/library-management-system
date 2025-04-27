const API_BASE_URL = 'http://localhost:8080';

// 通用的fetch封裝
async function fetchWithAuth(endpoint, options = {}) {
    const token = localStorage.getItem('token');
    const headers = {
        'Content-Type': 'application/json',
        ...options.headers,
    };

    if (token) {
        headers.Authorization = `Bearer ${token}`;
    }

    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
        ...options,
        headers,
    });

    if (!response.ok) {
        const error = await response.json().catch(() => ({}));
        throw new Error(error.error || '請求失敗');
    }

    return response.json();
}

// 認證相關API
export const auth = {
    login: (username, password) =>
        fetchWithAuth('/login', {
            method: 'POST',
            body: JSON.stringify({ username, password }),
        }),

    register: (username, password) =>
        fetchWithAuth('/register', {
            method: 'POST',
            body: JSON.stringify({ username, password }),
        }),

    logout: () =>
        fetchWithAuth('/logout', {
            method: 'POST',
        }),

    getCurrentUser: () =>
        fetchWithAuth('/me'),
};

// 書籍相關API
export const books = {
    getAll: () =>
        fetchWithAuth('/books'),

    getById: (id) =>
        fetchWithAuth(`/books/${id}`),

    create: (bookData) =>
        fetchWithAuth('/books', {
            method: 'POST',
            body: JSON.stringify(bookData),
        }),

    update: (id, bookData) =>
        fetchWithAuth(`/books/${id}`, {
            method: 'PUT',
            body: JSON.stringify(bookData),
        }),

    delete: (id) =>
        fetchWithAuth(`/books/${id}`, {
            method: 'DELETE',
        }),

    borrow: (id, note, borrower) =>
        fetchWithAuth(`/books/${id}/borrow`, {
            method: 'PUT',
            body: JSON.stringify({ 
                note,
                borrower 
            }),
        }),

    return: (id, borrower) =>
        fetchWithAuth(`/books/${id}/return`, {
            method: 'PUT',
            body: JSON.stringify({ borrower }),
        }),
};
