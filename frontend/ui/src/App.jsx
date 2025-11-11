import React, { useState } from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Navbar from './components/Navbar';
import HomePage from './pages/HomePage';
import AlbumPage from './pages/AlbumPage';
import ProfilePage from './pages/ProfilePage';
import AddArtistPage from './pages/AddArtistPage';
import AddAlbumPage from './pages/AddAlbumPage';

function App() {
  const [username, setUsername] = useState(localStorage.getItem('username'));

  const handleLogin = (loginData) => {
    if (loginData.username) {
      localStorage.setItem('username', loginData.username);
      setUsername(loginData.username);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('username');
    setUsername(null);
  };

  return (
    <>
      <Navbar
        username={username}
        onLogout={handleLogout}
        onAuth={handleLogin}
      />
      <div className="container mt-4">
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/albums/:id" element={<AlbumPage username={username} />} />
          <Route 
            path="/profile" 
            element={<ProfilePage username={username} />} 
          />
          <Route path="/add-artist" element={<AddArtistPage />} />
          <Route path="/add-album" element={<AddAlbumPage />} />
        </Routes>
      </div>
    </>
  );
}

export default App;