import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import LoginSignup from './LoginSignup';

export default function Navbar({ username, onLogout, onAuth }) {
  const [showAuth, setShowAuth] = useState(false);

  const handleAuthSuccess = (data) => {
    onAuth(data);
    setShowAuth(false);
  };

  return (
    <>
      <nav className="navbar navbar-expand-lg navbar-dark bg-dark">
        <div className="container-fluid">
          <Link className="navbar-brand" to="/">MusicReview</Link>
          <button className="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarContent">
            <span className="navbar-toggler-icon"></span>
          </button>
          <div className="collapse navbar-collapse" id="navbarContent">
            <ul className="navbar-nav ms-auto mb-2 mb-lg-0">
              {!username ? (
                <li className="nav-item">
                  <button 
                    className="btn btn-primary" 
                    onClick={() => setShowAuth(true)}
                  >
                    Login / Sign Up
                  </button>
                </li>
              ) : (
                <li className="nav-item dropdown">
                  <a className="nav-link dropdown-toggle" href="#" id="navbarDropdown" role="button" data-bs-toggle="dropdown">
                    Welcome, {username}
                  </a>
                  <ul className="dropdown-menu dropdown-menu-end">
                    <li><Link className="dropdown-item" to="/profile">My Profile</Link></li>
                    <li><Link className="dropdown-item" to="/add-artist">Add Artist</Link></li>
                    <li><Link className="dropdown-item" to="/add-album">Add Album</Link></li>
                    <li><hr className="dropdown-divider" /></li>
                    <li>
                      <button className="dropdown-item" onClick={onLogout}>
                        Logout
                      </button>
                    </li>
                  </ul>
                </li>
              )}
            </ul>
          </div>
        </div>
      </nav>

      {/* Auth Modal */}
      {showAuth && (
        <div className="modal-backdrop" onClick={() => setShowAuth(false)}>
          <div className="modal-content-container" onClick={(e) => e.stopPropagation()}>
            <button 
              className="btn-close btn-close-white" 
              onClick={() => setShowAuth(false)}
            ></button>
            <LoginSignup onAuth={handleAuthSuccess} />
          </div>
        </div>
      )}
    </>
  );
}