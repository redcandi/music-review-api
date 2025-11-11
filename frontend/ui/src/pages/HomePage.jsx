import React, { useEffect, useState, useCallback } from 'react';
import { Link } from 'react-router-dom';
import { getAlbums, searchAlbums } from '../api';

export default function HomePage() {
  const [albums, setAlbums] = useState([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');

  const loadAlbums = useCallback(async () => {
    setLoading(true);
    try {
      let data;
      if (searchTerm.trim() === '') {
        data = await getAlbums();
      } else {
        data = await searchAlbums(searchTerm);
      }
      setAlbums(data || []);
    } catch (error) {
      console.error(error);
      setAlbums([]);
    }
    setLoading(false);
  }, [searchTerm]);

  useEffect(() => {
    loadAlbums();
  }, [loadAlbums]);

  const handleSearch = (e) => {
    e.preventDefault();
    loadAlbums();
  };

  return (
    <div className="home-page">
      <h2 className="text-center mb-4">Find Your Next Favorite Album</h2>
      
      <form className="row justify-content-center mb-4" onSubmit={handleSearch}>
        <div className="col-md-6">
          <div className="input-group">
            <input 
              type="text" 
              className="form-control" 
              placeholder="Search by album title or artist..." 
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />
            <button className="btn btn-primary" type="submit">Search</button>
          </div>
        </div>
      </form>

      {loading ? (
        <div className="text-center">Loading...</div>
      ) : albums.length === 0 ? (
        <p className="text-muted text-center">No albums found.</p>
      ) : (
        <div className="list-group">
          {albums.map((album) => {
            const coverImg = album.cover_image_url || 'https://placehold.co/80x80/333/fff?text=♪';
            return (
              <Link
                key={album.album_id}
                to={`/albums/${album.album_id}`}
                className="list-group-item list-group-item-action"
              >
                <div className="d-flex align-items-center gap-3">
                  <img src={coverImg} alt={album.title} className="album-thumb" />
                  <div className="flex-grow-1">
                    <div className="d-flex justify-content-between">
                      <h6 className="mb-1">{album.title}</h6>
                      <small className="text-muted">
                        {album.average_rating != null ? `${album.average_rating.toFixed(1)}/10` : '-'}
                      </small>
                    </div>
                    <small className="text-muted">{album.artist_name}</small>
                  </div>
                </div>
              </Link>
            );
          })}
        </div>
      )}
    </div>
  );
}