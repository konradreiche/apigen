CREATE TABLE activities (
  id integer NOT NULL, 
  method text NOT NULL,
  data jsonb NOT NULL,
  error text,
  created_at timestamp without time zone NOT NULL DEFAULT now()
);
