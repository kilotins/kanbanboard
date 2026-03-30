-- Add project tag and task numbering

ALTER TABLE projects ADD COLUMN tag VARCHAR(4) NOT NULL DEFAULT '';
ALTER TABLE projects ADD COLUMN next_task_number INTEGER NOT NULL DEFAULT 1;

ALTER TABLE tasks ADD COLUMN task_number INTEGER NOT NULL DEFAULT 0;

-- Backfill: generate tags from existing project names and number existing tasks
DO $$
DECLARE
    proj RECORD;
    generated_tag TEXT;
    final_tag TEXT;
    suffix INTEGER;
    task_rec RECORD;
    task_num INTEGER;
BEGIN
    FOR proj IN SELECT id, name FROM projects ORDER BY created_at LOOP
        -- Generate tag from project name: first letter of each word, or first 3 letters if single word
        IF position(' ' IN proj.name) > 0 THEN
            generated_tag := '';
            FOR i IN 1..length(proj.name) LOOP
                IF i = 1 OR substring(proj.name FROM i-1 FOR 1) = ' ' THEN
                    IF substring(proj.name FROM i FOR 1) != ' ' THEN
                        generated_tag := generated_tag || upper(substring(proj.name FROM i FOR 1));
                    END IF;
                END IF;
            END LOOP;
            -- Truncate to 4 chars max
            generated_tag := substring(generated_tag FROM 1 FOR 4);
        ELSE
            generated_tag := upper(substring(proj.name FROM 1 FOR 3));
        END IF;

        -- Ensure at least 2 chars
        IF length(generated_tag) < 2 THEN
            generated_tag := generated_tag || 'X';
        END IF;

        -- Deduplicate by appending suffix if needed
        final_tag := generated_tag;
        suffix := 2;
        WHILE EXISTS(SELECT 1 FROM projects WHERE tag = final_tag AND id != proj.id) LOOP
            -- Truncate base to make room for suffix digit
            final_tag := substring(generated_tag FROM 1 FOR 3) || suffix::TEXT;
            suffix := suffix + 1;
        END LOOP;

        UPDATE projects SET tag = final_tag WHERE id = proj.id;

        -- Number tasks for this project ordered by created_at
        task_num := 1;
        FOR task_rec IN SELECT id FROM tasks WHERE project_id = proj.id ORDER BY created_at LOOP
            UPDATE tasks SET task_number = task_num WHERE id = task_rec.id;
            task_num := task_num + 1;
        END LOOP;

        -- Set next_task_number to the next available
        UPDATE projects SET next_task_number = task_num WHERE id = proj.id;
    END LOOP;
END $$;

-- Add constraints after backfill (so defaults don't conflict)
ALTER TABLE projects ADD CONSTRAINT projects_tag_unique UNIQUE (tag);
-- Remove the default now that backfill is done
ALTER TABLE projects ALTER COLUMN tag DROP DEFAULT;

CREATE UNIQUE INDEX idx_tasks_project_number ON tasks(project_id, task_number);
