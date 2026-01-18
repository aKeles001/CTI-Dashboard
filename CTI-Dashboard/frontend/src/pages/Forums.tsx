import React, { useState, useEffect } from 'react';
import { GetForums, SingularScrape, MultipleScrape, DeleteForum, Extract_posts} from '../../wailsjs/go/main/App';
import { models } from '../../wailsjs/go/models';
import { Button } from '@/components/ui/button';
import { toast } from "sonner"

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

  const handleScanS = (forum: models.Forum) => {
    toast.info("Scanning Forum. You can take your time information will be uptaded.")
    SingularScrape(forum).then(() => {
      toast.success("Forum Scanned :" + forum.forum_name);
      setForums((prevForums) => prevForums.map((f) => (f.forum_id === forum.forum_id ? forum : f)));
    }).catch((err) => {
      toast.error("Failed to scan forum: " + err);
    });
  };
  const handleScanM = (forums: models.Forum[]) => {
    toast.info("Scanning Forums. You can take your time information will be uptaded.")
    MultipleScrape(forums).then((failedForums) => {
      if (failedForums && failedForums.length > 0) {
        toast.error(`Scan completed. ${failedForums.length} forum failed.`);
      } else {
        toast.success("All Forums Scanned Successfully");
      }
    });
  };

  const handleDeleteForum = (forum: models.Forum) => {
    DeleteForum(forum.forum_id).then(() => {
      toast.success("Forum Deleted");
      setForums((prevForums) => prevForums.filter((f) => f.forum_id !== forum.forum_id));
    }).catch((err) => {
      toast.error("Failed to delete forum: " + err);
    });
  };

  const handleExtractPosts = (forum: models.Forum) => {
    Extract_posts(forum.forum_id).then((link_number) => {
      toast.success("Posts Extracted\n" + link_number + " posts extracted");
      if (window.confirm("Would you like to scrape the posts?")){
        
      } else {
        
      }
    }).catch((err) => {
      toast.error("Failed to extract posts: " + err);
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
          {forums.map((forum) => (
            <li key={forum.forum_id} className="border p-4 rounded-lg">
              <h2 className="text-xl font-semibold">Name: {forum.forum_name}</h2>
              <p className="text-md font-small mt-2">ID: {forum.forum_id}</p>
              <a href={forum.forum_html} target="_blank" rel="noopener noreferrer" className="text-blue-600">Link: {forum.forum_url}</a>
              <p className="text-sm text-gray-500 mt-2">Description: {forum.forum_description}</p>
              {forum.last_scaned && (
                <p className="text-sm text-gray-500 mt-2">Last scaned: {new Date(forum.last_scaned).toLocaleString()}</p>
              )}
              <Button
                  className="mt-4"
                  size="sm"
                  variant="outline"
                  onClick={() => handleScanS(forum)
                  }
                >
                  Scan Forum
                </Button>
                <Button
                  className="mt-4 "
                  size="sm"
                  variant="destructive"
                  onClick={() => handleDeleteForum(forum)
                  }
                >
                  Delete Forum
                </Button>
                <Button
                  className="mt-4 "
                  size="sm"
                  variant="destructive"
                  onClick={() => handleExtractPosts(forum)
                  }
                >
                  Extract Posts
                </Button>
            </li>
          ))}
          <Button
                  className="mt-4 sm:w-full"
                  size="lg"
                  variant="outline"
                  onClick={() => handleScanM(forums)
                  }
                >
                  Scan All forums
          </Button>
        </ul>
      )}
    </div>
  );
};

export default Forums;
