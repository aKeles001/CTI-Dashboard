import React, { useState } from 'react';
import { CreateForum } from '../../wailsjs/go/main/App';
import { Input } from '../components/ui/input';
import { Textarea } from '../components/ui/textarea';
import { Button } from '../components/ui/button';
import { Field, FieldContent, FieldDescription, FieldLabel, FieldSeparator, FieldSet } from '@/components/ui/field';
import { toast } from 'sonner';

const AddForum: React.FC = () => {
  const [url, setUrl] = useState('');
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [result, setResult] = useState('');


const handleSubmit = (e: React.FormEvent) => {
  e.preventDefault();
  const forumData = { forum_id: '', forum_url: url, forum_name: name, forum_description: description, last_scaned: '', forum_html: '', forum_screenshot: '', forum_engine: '' };
  CreateForum(forumData)
    .then((resultMessage: string) => {
      setResult(resultMessage);
      setName('');
      setUrl('');
      setDescription('');
      toast.success("Forum Created");
    }).catch((errorMessage: string) => {
      setResult(errorMessage);
      toast.error("Failed to create forum");
    });
};

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold">Home</h1>
      <p>Welcome to CTI Dashboard.</p>
      <div> 
        <form onSubmit={handleSubmit} className="mt-6 max-w-lg">
          <h2 className="text-xl font-semibold">Submit a Forum</h2>
          <FieldSet>
            <Field orientation="responsive">
              <FieldContent>
                <FieldLabel htmlFor='name'>Name</FieldLabel>
                <FieldDescription>Name of the forum.</FieldDescription>
              </FieldContent>
              <FieldSeparator />
              <Input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
              />
            </Field>
            <FieldSeparator />
            <Field orientation="responsive">
              <FieldContent>
                <FieldLabel htmlFor='url'>URL</FieldLabel>
                <FieldDescription>URL of the forum.</FieldDescription>
              </FieldContent>
              <FieldSeparator />
              <Input
                type="text"
                value={url}
                onChange={(e) => setUrl(e.target.value)}
              />
            </Field>
            <FieldSeparator />
            <Field orientation="responsive">
              <FieldContent>
                <FieldLabel htmlFor='description'>Description</FieldLabel>
                <FieldDescription>Enter a brief description of the forum.</FieldDescription>
              </FieldContent>
              <FieldSeparator />
              <Textarea
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                rows={3}
              />
            </Field>
            <FieldSeparator />
            <Button type="submit">Submit Forum</Button>
            {result && <p className="mt-4">{result}</p>}
          </FieldSet>
        </form>
      </div>
    </div>
  );
};

export default AddForum;
