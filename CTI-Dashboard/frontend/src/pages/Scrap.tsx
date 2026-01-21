import React, { useState, useEffect } from 'react';
import { GetForums, Extract_posts, GetPosts, ScanPosts} from '../../wailsjs/go/main/App';
import { models } from '../../wailsjs/go/models';
import { toast } from "sonner"
import { Button } from '@/components/ui/button';
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion"

const Scrap: React.FC = () => {
  const [forums, setForums] = useState<models.Forum[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [postCounts, setPostCounts] = useState<Record<string, number>>({});
  const [posts, setPosts] = useState<models.Post[]>([]);
  const [isLoadingPosts, setIsLoadingPosts] = useState<boolean>(false);
  


  useEffect(() => {
    GetForums()
      .then((data) => {
        console.log(data);
        setForums(data);
      })
      .catch((err) => {
        setError(err);
      });
  }, []);
  
  if (error) {
    return <div className="p-4">Error loading forums: {error}</div>;
  }
  const handlePostCount = (forum: models.Forum) => {
    Extract_posts(forum.forum_id).then((count) => {
      setPostCounts(prevCounts => ({...prevCounts, [forum.forum_id]: count}));
      toast.success("Post count fetched");
    }).catch((err) => {
      setError(err);
    });
  };

  const handleGetPosts1 = (forumId: string) => {
    setIsLoadingPosts(true);
    ScanPosts(forumId)
      .then(() => {
        toast.success("Posts scanned successfully");
      })
      .catch((err) => {
        setError(err);
        toast.error("Failed to scan posts");
      }).finally(() => {
        setIsLoadingPosts(false);
      });
  
  }
  const handleGetPosts = (forumId: string) => {
    GetPosts(forumId)
      .then((data) => {
        setPosts(data);
        toast.success("Posts fetched successfully");
      })
      .catch((err) => {
        setError(err);
        toast.error("Failed to fetch posts");
      })
      .finally(() => {
        setIsLoadingPosts(false);
      });
  };

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Forums</h1>
      {forums.length === 0 ? (
        <p>No forums available.</p>
      ) : (
          <div className="space-y-4">
            {forums.map((forum) => (
              <Accordion key={forum.forum_id} type="single" collapsible>
                <AccordionItem value={forum.forum_id}>
                  <AccordionTrigger><h2 className="text-xl font-semibold">Name: {forum.forum_name}</h2></AccordionTrigger>
                    <AccordionContent>
                      <div className="border p-4 rounded-lg">
                        <h2 className="text-xl font-semibold">Name: {forum.forum_name}</h2>
                        <p className="text-md font-small mt-2">ID: {forum.forum_id}</p>
                        {forum.last_scaned && (
                          <p className="text-sm text-gray-500 mt-2">Last scanned: {new Date(forum.last_scaned).toLocaleString()}</p>
                        )}
                        <Button
                          className="mt-4 sm:w-full"
                          size="lg"
                          variant="outline"
                          onClick={() => handlePostCount(forum)}
                          >
                          Fetch Post Count
                        </Button>
                        {postCounts[forum.forum_id] !== undefined && <p className="mt-4">Number of posts: {postCounts[forum.forum_id]}</p>}
                        <div className="mt-8">
                        <Button
                            className="sm:w-full"
                            size="lg"
                            onClick={() => handleGetPosts(forum.forum_id)}
                            disabled={isLoadingPosts}
                        >
                            {isLoadingPosts ? 'Loading Posts...' : 'Show Posts'}
                        </Button>
                        <Button
                            className="sm:w-full"
                            size="lg"
                            onClick={() => handleGetPosts1(forum.forum_id)}
                            disabled={isLoadingPosts}
                        >
                            {isLoadingPosts ? 'Loading Posts...' : 'Show Posts1'}
                        </Button>
                        </div>

                        {posts.length > 0 && (
                            <div className="mt-8">
                                <h2 className="text-xl font-bold mb-4">Posts</h2>
                                <ul className="space-y-4">
                                    {posts.map((post, index) => (
                                        <li key={index} className="border p-4 rounded-lg">
                                            <p className="text-sm text-gray-500">{post.status}</p>
                                            <p className="text-sm text-gray-500">{post.severity_level}</p>
                                        </li>
                                    ))}
                                </ul>
                            </div>
                        )}
                      </div>
                    </AccordionContent>
                </AccordionItem>
              </Accordion>
            ))}
            
          </div>
      )}

        
    </div>
  );
};
export default Scrap;
