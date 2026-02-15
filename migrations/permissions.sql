-- Permissions for Authenticated Users
GRANT USAGE ON SCHEMA public TO authenticated;
GRANT ALL ON TABLE public.tasks TO authenticated;

-- OPTION: If the Go app cannot send the auth token (library limitation),
-- the request will be treated as 'anon'. 
-- To allow it to work (INSECURELY for testing):
-- GRANT ALL ON TABLE public.tasks TO anon;

-- Ensure RLS is enabled

-- Ensure RLS is enabled
ALTER TABLE public.tasks ENABLE ROW LEVEL SECURITY;

-- Policies (Ensure these exist)
-- Re-create to be sure (drop if needed: drop policy if exists ...)

DROP POLICY IF EXISTS "Users can view their own tasks" ON public.tasks;
CREATE POLICY "Users can view their own tasks" 
ON public.tasks FOR SELECT 
USING (auth.uid() = user_id);

DROP POLICY IF EXISTS "Users can insert their own tasks" ON public.tasks;
CREATE POLICY "Users can insert their own tasks" 
ON public.tasks FOR INSERT 
WITH CHECK (auth.uid() = user_id);

DROP POLICY IF EXISTS "Users can update their own tasks" ON public.tasks;
CREATE POLICY "Users can update their own tasks" 
ON public.tasks FOR UPDATE 
USING (auth.uid() = user_id);

DROP POLICY IF EXISTS "Users can delete their own tasks" ON public.tasks;
CREATE POLICY "Users can delete their own tasks" 
ON public.tasks FOR DELETE 
USING (auth.uid() = user_id);
