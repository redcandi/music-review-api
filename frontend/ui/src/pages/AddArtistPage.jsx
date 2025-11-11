import React, { useState } from 'react';
import { createArtist } from '../api';

export default function AddArtistPage() {
  const [name, setName] = useState('');
  const [bio, setBio] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    setSubmitting(true);
    try {
      await createArtist({ name, bio: bio || undefined });
      setSuccess('Artist created successfully.');
      setName('');
      setBio('');
    } catch (err) {
      setError(err.message || 'Failed to create artist');
    }
    setSubmitting(false);
  };

  return (
    <div className="row justify-content-center">
      <div className="col-md-6">
        <div className="card shadow-sm">
          <div className="card-body">
            <h3 className="card-title mb-3">Add Artist</h3>
            {error && <div className="alert alert-danger">{error}</div>}
            {success && <div className="alert alert-success">{success}</div>}
            <form onSubmit={handleSubmit}>
              <div className="mb-3">
                <label className="form-label">Name</label>
                <input
                  className="form-control"
                  placeholder="Artist name"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  required
                />
              </div>
              <div className="mb-3">
                <label className="form-label">Bio (optional)</label>
                <textarea
                  className="form-control"
                  rows="3"
                  placeholder="Short bio"
                  value={bio}
                  onChange={(e) => setBio(e.target.value)}
                ></textarea>
              </div>
              <button className="btn btn-primary" type="submit" disabled={submitting}>
                {submitting ? 'Saving...' : 'Create Artist'}
              </button>
            </form>
          </div>
        </div>
      </div>
    </div>
  );
}

