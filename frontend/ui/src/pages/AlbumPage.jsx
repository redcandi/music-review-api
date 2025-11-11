import React, { useEffect, useState, useCallback } from 'react';
import { useParams } from 'react-router-dom';
import { getAlbumDetails, postComment } from '../api';
import CommentForm from '../components/CommentForm';

export default function AlbumPage({ username }) {
  const { id } = useParams();
  const [album, setAlbum] = useState(null);
  const [comments, setComments] = useState([]);
  const [genres, setGenres] = useState([]);
  const [loading, setLoading] = useState(true);

  const loadData = useCallback(async () => {
    try {
      const data = await getAlbumDetails(id);
      setAlbum(data.album_details);
      setComments(data.comments || []);
      setGenres(data.genres || []);
    } catch (error) {
      console.error(error);
    }
    setLoading(false);
  }, [id]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const handleCommentPosted = async (commentData) => {
    const dataToSend = {
      ...commentData,
      username: username || '',
    };

    try {
      await postComment(id, dataToSend);
      loadData();
    } catch (error) {
      console.error(error);
      alert(error.message);
    }
  };

  if (loading) return <div className="text-center">Loading...</div>;
  if (!album) return <div className="text-center">Album not found.</div>;

  const coverImg = album.cover_image_url || 'https://placehold.co/200x200/333/fff?text=No+Art';

  return (
    <div className="row">
      <div className="col-md-4">
        <img src={coverImg} className="img-fluid rounded shadow-sm" alt={album.title} />
        <h2 className="mt-3">{album.title}</h2>
        <h4 className="text-muted">{album.artist_name}</h4>
        <p>Released: {new Date(album.release_date).toLocaleDateString()}</p>
        <div>
          {genres.map(g => (
            <span key={g.id} className="badge bg-secondary me-1">{g.name}</span>
          ))}
        </div>
      </div>

      <div className="col-md-8">
        <CommentForm onSubmit={handleCommentPosted} username={username} />
        
        <h3 className="mt-4">Comments</h3>
        <div className="comment-list">
          {comments.length === 0 ? (
            <p>Be the first to comment on this album!</p>
          ) : (
            comments.map(comment => (
              <div key={comment.id} className="card mb-2">
                <div className="card-body">
                  <div className="d-flex justify-content-between">
                    <strong>{comment.username}</strong>
                    <span className="text-warning">
                      {'★'.repeat(comment.rating)}{'☆'.repeat(10 - comment.rating)}
                    </span>
                  </div>
                  <p className="card-text mt-2">{comment.comment_text}</p>
                  <small className="text-muted">
                    {new Date(comment.created_at).toLocaleString()}
                  </small>
                </div>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}