
import React, { useState, useEffect } from 'react';
import { GetForums, SingularScrape } from '../../wailsjs/go/main/App';
import { models } from '../../wailsjs/go/models';
import { Button } from '@/components/ui/button';

const Forums: React.FC = () => {
  const [forums, setForums] = useState<models.Forum[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);


  useEffect(() => {
    GetForums()
      .then((data) => {
        console.log(data);
        setForums(data);
        setLoading(false);
      })
      .catch((err) => {
        setError(err);
        setLoading(false);
      });
  }, []);

  const handleScan = (forum: models.Forum) => {
    window.alert(`Scanning forum: ${forum.forum_name}`);
    SingularScrape(forum).then((result) => {
      console.log(result);
    });
  };

  if (loading) {
    return <div className="p-4">Loading forums...</div>;
  }

  if (error) {
    return <div className="p-4">Error loading forums: {error}</div>;
  }
  

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Forums</h1>
      {forums.length === 0 ? (
        <p>No forums available.</p>
      ) : (
        <ul className="space-y-4">
          {forums.map((forum, index) => (
            <li key={index} className="border p-4 rounded-lg">
              <h2 className="text-xl font-semibold">Name: {forum.forum_name}</h2>
              <p className="text-md font-small mt-2">ID: {forum.forum_id}</p>
              <a href={forum.forum_url} target="_blank" rel="noopener noreferrer" className="text-blue-600">Link: {forum.forum_url}</a>
              <p className="text-sm text-gray-500 mt-2">Description: {forum.forum_description}</p>
              {forum.last_scaned && (
                <p className="text-sm text-gray-500 mt-2">Last scaned: {new Date(forum.last_scaned).toLocaleString()}</p>
              )}
              <Button className="mt-4" onClick={() => handleScan(forum)}>
                Scan Forum
              </Button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default Forums;
