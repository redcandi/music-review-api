const API_BASE_URL = 'http://localhost:8080/api/v1';

async function request(endpoint, options = {}) {
  const url = `${API_BASE_URL}${endpoint}`;
  const config = {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  };
  if (options.body) {
    config.body = JSON.stringify(options.body);
  }

  const response = await fetch(url, config);

  const rawText = await response.text();
  let parsed;
  try {
    parsed = rawText ? JSON.parse(rawText) : null;
  } catch (_) {
    parsed = null;
  }

  if (!response.ok) {
    const message =
      (parsed && (parsed.error || parsed.message)) ||
      rawText ||
      'An unknown error occurred';
    throw new Error(message);
  }

  const contentType = response.headers.get('content-type') || '';
  if (parsed !== null) return parsed;
  if (rawText && !contentType.includes('application/json')) return rawText;
  return null;
}

export function signup(userData) {
  return request('/signup', { method: 'POST', body: userData });
}

export function login(credentials) {
  return request('/login', { method: 'POST', body: credentials });
}

export function getAlbums() {
  return request('/albums');
}

export function searchAlbums(query) {
  return request(`/albums/search?q=${encodeURIComponent(query)}`);
}

export function getAlbumDetails(albumId) {
  return request(`/albums/${albumId}`);
}

export function postComment(albumId, commentData) {
  return request(`/albums/${albumId}/comments`, {
    method: 'POST',
    body: commentData,
  });
}

export function getProfile(username) {
  return request(`/users/${username}/comments`);
}

export function deleteUser(username) {
  return request(`/users/${username}`, { method: 'DELETE' });
}
export function createArtist(artistData) {
  return request('/artists', { method: 'POST', body: artistData });
}

export function createAlbum(albumData) {
  return request('/albums', { method: 'POST', body: albumData });
}