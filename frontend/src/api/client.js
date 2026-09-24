const API_URL = import.meta.env.VITE_API_URL;

async function request(path, options = {}) {
      const token = localStorage.getItem('token');

      const headers = {
            'Content-Type': 'application/json',
            ...options.headers,
      };

      if (token) {
            headers['Authorization'] = `Bearer ${token}`;
      }

      const res = await fetch(`${API_URL}${path}`, { ...options, headers });

      const data = await res.json().catch(() => null);

      if (!res.ok) {
            const message = data?.error?.message || 'Something went wrong';
            throw new Error(message);
      }

      return data;
}

export const api = {
      post: (path, body) => request(path, { method: 'POST', body: JSON.stringify(body) }),
      get: (path) => request(path, { method: 'GET' }),
      put: (path, body) => request(path, { method: 'PUT', body: JSON.stringify(body) }),
};