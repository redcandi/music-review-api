import React from 'react';
import { Link } from 'react-router-dom';

export default function AlbumCard({ album }) {
  const coverImg = album.cover_image_url || 'https://placehold.co/200x200/333/fff?text=No+Art';
  
  return (
    <div className="card album-card">
      <Link to={`/albums/${album.album_id}`}>
        <img src={coverImg} className="card-img-top" alt={album.title} />
      </Link>
      <div className="card-body">
        <h5 className="card-title">
          <Link to={`/albums/${album.album_id}`}>{album.title}</Link>
        </h5>
        <p className="card-text">{album.artist_name}</p>
      </div>
      <div className="card-footer">
        <small className="text-muted">
          Rating: {album.average_rating.toFixed(1)}/10 ({album.total_comments} reviews)
        </small>
      </div>
    </div>
  );
}