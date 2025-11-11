import React, { useState } from 'react';
import { createAlbum } from '../api';

export default function AddAlbumPage() {
  const [title, setTitle] = useState('');
  const [artistId, setArtistId] = useState('');
  const [releaseDate, setReleaseDate] = useState('');
  const [coverUrl, setCoverUrl] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    setSubmitting(true);
    try {
      await createAlbum({
        title,
        artist_id: Number(artistId),
        release_date: releaseDate || undefined,
        cover_image_url: coverUrl || undefined,
      });
      setSuccess('Album created successfully.');
      setTitle('');
      setArtistId('');
      setReleaseDate('');
      setCoverUrl('');
    } catch (err) {
      setError(err.message || 'Failed to create album');
    }
    setSubmitting(false);
  };

  return (
    <div className="row justify-content-center">
      <div className="col-md-6">
        <div className="card shadow-sm">
          <div className="card-body">
            <h3 className="card-title mb-3">Add Album</h3>
            {error && <div className="alert alert-danger">{error}</div>}
            {success && <div className="alert alert-success">{success}</div>}
            <form onSubmit={handleSubmit}>
              <div className="mb-3">
                <label className="form-label">Title</label>
                <input
                  className="form-control"
                  placeholder="Album title"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  required
                />
              </div>
              <div className="mb-3">
                <label className="form-label">Artist ID</label>
                <input
                  type="number"
                  className="form-control"
                  placeholder="e.g., 12"
                  value={artistId}
                  onChange={(e) => setArtistId(e.target.value)}
                  min="1"
                  required
                />
              </div>
              <div className="mb-3">
                <label className="form-label">Release Date (optional)</label>
                <input
                  type="date"
                  className="form-control"
                  value={releaseDate}
                  onChange={(e) => setReleaseDate(e.target.value)}
                />
              </div>
              <div className="mb-3">
                <label className="form-label">Cover Image URL (optional)</label>
                <input
                  type="url"
                  className="form-control"
                  placeholder="https://..."
                  value={coverUrl}
                  onChange={(e) => setCoverUrl(e.target.value)}
                />
              </div>
              <button className="btn btn-primary" type="submit" disabled={submitting}>
                {submitting ? 'Saving...' : 'Create Album'}
              </button>
            </form>
          </div>
        </div>
      </div>
    </div>
  );
}

