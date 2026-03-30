<script>
  import { updateProject, createColumn, updateColumn, deleteColumn, reorderColumns, createLabel, updateLabel, deleteLabel } from './api.js';

  let { project, onBack, onProjectUpdated } = $props();

  let projectName = $state(project.name);
  let visibility = $state(project.visibility);
  let projectTag = $state(project.tag);
  let message = $state('');
  let error = $state('');

  let hasAnyTasks = $derived((project.tasks || []).length > 0);

  // Column state
  let newColumnName = $state('');
  let columnError = $state('');

  // Label state
  let newLabelName = $state('');
  let newLabelColor = $state('#808080');
  let labelError = $state('');

  $effect(() => {
    projectName = project.name;
    visibility = project.visibility;
    projectTag = project.tag;
  });

  async function handleProjectNameBlur() {
    if (projectName.trim() && projectName !== project.name) {
      try {
        await updateProject(project.id, { name: projectName.trim() });
        onProjectUpdated();
        message = 'Project name updated.';
      } catch (err) {
        error = err.message;
      }
    }
  }

  async function handleTagBlur() {
    const tagValue = projectTag.toUpperCase();
    if (tagValue && tagValue !== project.tag) {
      try {
        await updateProject(project.id, { tag: tagValue });
        onProjectUpdated();
        message = 'Tag updated.';
      } catch (err) {
        error = err.message;
        projectTag = project.tag;
      }
    }
  }

  async function handleVisibilityChange() {
    try {
      await updateProject(project.id, { visibility });
      onProjectUpdated();
      message = 'Visibility updated.';
    } catch (err) {
      error = err.message;
    }
  }

  // Column handlers
  async function handleAddColumn() {
    columnError = '';
    if (!newColumnName.trim()) return;
    try {
      await createColumn(project.id, newColumnName.trim());
      newColumnName = '';
      onProjectUpdated();
    } catch (err) {
      columnError = err.message;
    }
  }

  async function handleRenameColumn(colId, newName) {
    if (!newName.trim()) return;
    try {
      await updateColumn(project.id, colId, newName.trim());
      onProjectUpdated();
    } catch (err) {
      columnError = err.message;
    }
  }

  async function handleDeleteColumn(colId) {
    columnError = '';
    try {
      await deleteColumn(project.id, colId);
      onProjectUpdated();
    } catch (err) {
      columnError = err.message;
    }
  }

  async function handleMoveColumn(index, direction) {
    const cols = [...project.columns].sort((a, b) => a.position - b.position);
    const newIndex = index + direction;
    if (newIndex < 0 || newIndex >= cols.length) return;

    // Swap
    [cols[index], cols[newIndex]] = [cols[newIndex], cols[index]];
    const ids = cols.map(c => c.id);

    try {
      await reorderColumns(project.id, ids);
      onProjectUpdated();
    } catch (err) {
      columnError = err.message;
    }
  }

  // Label handlers
  async function handleAddLabel() {
    labelError = '';
    if (!newLabelName.trim()) return;
    try {
      await createLabel(project.id, newLabelName.trim(), newLabelColor);
      newLabelName = '';
      newLabelColor = '#808080';
      onProjectUpdated();
    } catch (err) {
      labelError = err.message;
    }
  }

  async function handleUpdateLabel(labelId, name, color) {
    if (!name.trim()) return;
    try {
      await updateLabel(project.id, labelId, name.trim(), color);
      onProjectUpdated();
    } catch (err) {
      labelError = err.message;
    }
  }

  async function handleDeleteLabel(labelId) {
    labelError = '';
    try {
      await deleteLabel(project.id, labelId);
      onProjectUpdated();
    } catch (err) {
      labelError = err.message;
    }
  }

  let sortedColumns = $derived(
    [...(project.columns || [])].sort((a, b) => a.position - b.position)
  );
</script>

<div class="settings-page">
  <div class="header">
    <button class="back-btn" onclick={onBack}>← Back to Board</button>
    <h1>Project Settings</h1>
  </div>

  <div class="content">
    {#if message}
      <p class="success">{message}</p>
    {/if}
    {#if error}
      <p class="error">{error}</p>
    {/if}

    <!-- Project info -->
    <section>
      <h2>Project</h2>
      <div class="field">
        <label for="projectName">Name</label>
        <input id="projectName" type="text" bind:value={projectName} onblur={handleProjectNameBlur} />
      </div>
      <div class="field">
        <label for="projectTag">Tag</label>
        {#if hasAnyTasks}
          <input id="projectTag" type="text" value={projectTag} disabled />
          <span class="hint">Tag cannot be changed after tasks are created</span>
        {:else}
          <input id="projectTag" type="text" bind:value={projectTag} onblur={handleTagBlur} maxlength="4" style="text-transform: uppercase; width: 100px;" />
        {/if}
      </div>
      <div class="field">
        <label for="visibility">Visibility</label>
        <select id="visibility" bind:value={visibility} onchange={handleVisibilityChange}>
          <option value="public">Public</option>
          <option value="private">Private</option>
        </select>
      </div>
    </section>

    <!-- Columns -->
    <section>
      <h2>Columns</h2>
      {#if columnError}
        <p class="error">{columnError}</p>
      {/if}
      <div class="item-list">
        {#each sortedColumns as col, i (col.id)}
          <div class="item-row">
            <button class="move-btn" onclick={() => handleMoveColumn(i, -1)} disabled={i === 0}>▲</button>
            <button class="move-btn" onclick={() => handleMoveColumn(i, 1)} disabled={i === sortedColumns.length - 1}>▼</button>
            <input
              type="text"
              value={col.name}
              onblur={(e) => handleRenameColumn(col.id, e.target.value)}
            />
            <button class="delete-btn" onclick={() => handleDeleteColumn(col.id)}>✕</button>
          </div>
        {/each}
      </div>
      <div class="add-row">
        <input type="text" placeholder="New column name" bind:value={newColumnName}
          onkeydown={(e) => e.key === 'Enter' && handleAddColumn()} />
        <button class="add-btn" onclick={handleAddColumn}>Add</button>
      </div>
    </section>

    <!-- Labels -->
    <section>
      <h2>Labels</h2>
      {#if labelError}
        <p class="error">{labelError}</p>
      {/if}
      <div class="item-list">
        {#each project.labels as lbl (lbl.id)}
          <div class="item-row">
            <input
              type="color"
              value={lbl.color}
              class="color-picker"
              onchange={(e) => handleUpdateLabel(lbl.id, lbl.name, e.target.value)}
            />
            <input
              type="text"
              value={lbl.name}
              onblur={(e) => handleUpdateLabel(lbl.id, e.target.value, lbl.color)}
            />
            <button class="delete-btn" onclick={() => handleDeleteLabel(lbl.id)}>✕</button>
          </div>
        {/each}
      </div>
      <div class="add-row">
        <input type="color" class="color-picker" bind:value={newLabelColor} />
        <input type="text" placeholder="New label name" bind:value={newLabelName}
          onkeydown={(e) => e.key === 'Enter' && handleAddLabel()} />
        <button class="add-btn" onclick={handleAddLabel}>Add</button>
      </div>
    </section>
  </div>
</div>

<style>
  .settings-page {
    max-width: 550px;
    margin: 0 auto;
    padding: 24px;
  }

  .header {
    margin-bottom: 24px;
  }

  .back-btn {
    background: none;
    border: none;
    color: #4a90d9;
    cursor: pointer;
    font-size: 0.875rem;
    padding: 0;
    margin-bottom: 8px;
  }

  .back-btn:hover {
    text-decoration: underline;
  }

  h1 {
    font-size: 1.5rem;
    color: #333;
    margin: 0;
  }

  h2 {
    font-size: 1.1rem;
    color: #333;
    margin: 0 0 12px;
  }

  section {
    background: white;
    border: 1px solid #e0e0e0;
    border-radius: 6px;
    padding: 20px;
    margin-bottom: 20px;
  }

  .field {
    margin-bottom: 12px;
  }

  label {
    display: block;
    font-size: 0.8rem;
    font-weight: 500;
    color: #555;
    margin-bottom: 4px;
  }

  input[type="text"], select {
    width: 100%;
    padding: 6px 10px;
    border: 1px solid #ccc;
    border-radius: 4px;
    font-size: 0.875rem;
    box-sizing: border-box;
  }

  input:focus, select:focus {
    outline: none;
    border-color: #4a90d9;
  }

  .item-list {
    display: flex;
    flex-direction: column;
    gap: 6px;
    margin-bottom: 12px;
  }

  .item-row {
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .item-row input[type="text"] {
    flex: 1;
  }

  .move-btn {
    padding: 2px 6px;
    background: none;
    border: 1px solid #ddd;
    border-radius: 3px;
    cursor: pointer;
    font-size: 0.7rem;
    color: #666;
  }

  .move-btn:hover:not(:disabled) {
    background: #f0f0f0;
  }

  .move-btn:disabled {
    opacity: 0.3;
    cursor: default;
  }

  .delete-btn {
    padding: 4px 8px;
    background: none;
    border: 1px solid #e0e0e0;
    border-radius: 3px;
    cursor: pointer;
    color: #c00;
    font-size: 0.8rem;
  }

  .delete-btn:hover {
    background: #fff5f5;
  }

  .add-row {
    display: flex;
    gap: 6px;
    align-items: center;
  }

  .add-row input[type="text"] {
    flex: 1;
  }

  .add-btn {
    padding: 6px 12px;
    background: #4a90d9;
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 0.875rem;
    cursor: pointer;
    white-space: nowrap;
  }

  .add-btn:hover {
    background: #357abd;
  }

  .color-picker {
    width: 32px;
    height: 32px;
    padding: 2px;
    border: 1px solid #ccc;
    border-radius: 4px;
    cursor: pointer;
    flex-shrink: 0;
  }

  .hint {
    display: block;
    font-size: 0.75rem;
    color: #888;
    margin-top: 4px;
  }

  .error {
    color: #c00;
    font-size: 0.85rem;
    margin: 0 0 8px;
  }

  .success {
    color: #0a0;
    font-size: 0.85rem;
    margin: 0 0 8px;
  }
</style>
