import React, { useEffect, useState } from 'react';
import { Link, Navigate, useNavigate } from 'react-router-dom';
import { getProfile, deleteUser } from '../api';

export default function ProfilePage({ username }) {
  const [comments, setComments] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    if (username) {
      getProfile(username)
        .then(data => {
          setComments(data || []);
        })
        .catch(error => {
          console.error(error);
        })
        .finally(() => {
          setLoading(false);
        });
    }
  }, [username]);

  if (!username) {
    return <Navigate to="/" replace />;
  }

  if (loading) return <div className="text-center">Loading profile...</div>;

  return (
    <div>
      <div className="d-flex align-items-center justify-content-between mb-3">
        <h2 className="m-0">My Comments ({comments.length})</h2>
        <button
          className="btn btn-outline-danger btn-sm"
          onClick={async () => {
            if (!window.confirm('This will permanently delete your profile. Continue?')) {
              return;
            }
            try {
              await deleteUser(username);
              localStorage.removeItem('username');
              navigate('/', { replace: true });
              window.location.reload();
            } catch (e) {
              setError(e.message || 'Failed to delete profile');
            }
          }}
        >
          Delete Profile
        </button>
      </div>
      {error && <div className="alert alert-danger">{error}</div>}
      {comments.length === 0 ? (
        <p>You haven't posted any comments yet.</p>
      ) : (
        <div className="list-group">
          {comments.map(comment => (
            <Link 
              key={comment.id} 
              to={`/albums/${comment.album_id}`} 
              className="list-group-item list-group-item-action"
            >
              <div className="d-flex w-100 justify-content-between">
                <h5 className="mb-1">On: {comment.album_title}</h5>
                <small>{new Date(comment.created_at).toLocaleDateString()}</small>
              </div>
              <p className="mb-1">"{comment.comment_text}"</p>
              <small className="text-warning">
                {'★'.repeat(comment.rating)}{'☆'.repeat(10 - comment.rating)}
              </small>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}