import React, { useState } from 'react';

export default function CommentForm({ onSubmit, username }) {
  const [rating, setRating] = useState(5);
  const [commentText, setCommentText] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!commentText) {
      alert("Please write a comment.");
      return;
    }
    onSubmit({ rating: Number(rating), comment_text: commentText });
    setCommentText('');
    setRating(5);
  };

  return (
    <div className="card shadow-sm">
      <div className="card-body">
        <h5 className="card-title">
          {username ? `Commenting as ${username}` : 'Post an anonymous comment'}
        </h5>
        <form onSubmit={handleSubmit}>
          <div className="mb-2">
            <label className="form-label">Rating (1-10)</label>
            <input
              type="range"
              className="form-range"
              min="1"
              max="10"
              step="1"
              value={rating}
              onChange={e => setRating(e.target.value)}
            />
            <div className="text-center fw-bold">{rating}</div>
          </div>
          <div className="mb-2">
            <label className="form-label">Comment</label>
            <textarea
              className="form-control"
              rows="3"
              placeholder="What did you think of the album?"
              value={commentText}
              onChange={e => setCommentText(e.target.value)}
            ></textarea>
          </div>
          <button type="submit" className="btn btn-primary">Post Comment</button>
        </form>
      </div>
    </div>
  );
}